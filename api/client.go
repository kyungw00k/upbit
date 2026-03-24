// Package api provides an HTTP client for the Upbit exchange API.
// See https://docs.upbit.com/reference/ for API documentation.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kyungw00k/upbit/ratelimit"
	"github.com/kyungw00k/upbit/retry"
)

const (
	BaseURL        = "https://api.upbit.com/v1"
	defaultTimeout = 30 * time.Second
)

// Client is an HTTP client for the Upbit API.
type Client struct {
	baseURL    string
	accessKey  string
	secretKey  string
	httpClient *http.Client
	limiter    *ratelimit.Limiter
}

// NewClient creates a new client with default settings.
func NewClient(accessKey, secretKey string) *Client {
	return &Client{
		baseURL:   BaseURL,
		accessKey: accessKey,
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		limiter: ratelimit.NewLimiter(),
	}
}

// NewClientWithURL creates a client with a custom base URL (for testing/mocking).
func NewClientWithURL(baseURL, accessKey, secretKey string) *Client {
	c := NewClient(accessKey, secretKey)
	c.baseURL = baseURL
	return c
}

// Request executes an HTTP request with integrated rate limiting, retry, and error parsing.
func (c *Client) Request(ctx context.Context, method, path string, query map[string]string, body interface{}, result interface{}) error {
	group := ratelimit.GroupFromPath(path)

	return retry.Retry(ctx, func() error {
		// Wait for rate limit token.
		if err := c.limiter.Wait(ctx, group); err != nil {
			return err
		}

		// Build URL.
		rawURL := c.baseURL + path
		u, err := url.Parse(rawURL)
		if err != nil {
			return fmt.Errorf("URL 파싱 실패: %w", err)
		}

		// Set query parameters.
		if len(query) > 0 {
			q := u.Query()
			for k, v := range query {
				q.Set(k, v)
			}
			u.RawQuery = q.Encode()
		}

		// Serialize request body.
		var bodyReader io.Reader
		var bodyParams map[string]string
		if body != nil {
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("요청 본문 직렬화 실패: %w", err)
			}
			bodyReader = bytes.NewReader(bodyBytes)

			// For POST auth: build a query string preserving JSON body key order.
			// Per Upbit auth docs: convert JSON body Key-Value pairs to query string maintaining order.
			if method == http.MethodPost {
				bodyQueryString := jsonBodyToQueryString(bodyBytes)
				if bodyQueryString != "" {
					bodyParams = map[string]string{"__raw_query__": bodyQueryString}
				}
			}
		}

		// Create HTTP request.
		req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
		if err != nil {
			return fmt.Errorf("HTTP 요청 생성 실패: %w", err)
		}

		// Set common headers.
		req.Header.Set("Accept", "application/json")
		// gzip encoding is applied only for the Quotation API (market data).
		if isQuotationPath(path) {
			req.Header.Set("Accept-Encoding", "gzip")
		}
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		// Attach auth token if credentials are set.
		if c.accessKey != "" && c.secretKey != "" {
			var token string
			var tokenErr error

			if method == http.MethodPost && bodyParams != nil {
				// POST: hash using the raw query string with preserved JSON body key order.
				rawQuery := bodyParams["__raw_query__"]
				token, tokenErr = GenerateTokenFromRawQuery(c.accessKey, c.secretKey, rawQuery)
			} else {
				token, tokenErr = GenerateToken(c.accessKey, c.secretKey, query)
			}

			if tokenErr != nil {
				return fmt.Errorf("JWT 토큰 생성 실패: %w", tokenErr)
			}
			req.Header.Set("Authorization", "Bearer "+token)
		}

		// Execute request.
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("HTTP 요청 실패: %w", err)
		}
		defer resp.Body.Close()

		// Dynamically update rate limiter from Remaining-Req header.
		if remainingReq := resp.Header.Get("Remaining-Req"); remainingReq != "" {
			c.limiter.UpdateFromHeader(remainingReq, group)
		}

		// Read response body (supports gzip decompression).
		respBody, err := readResponseBody(resp)
		if err != nil {
			return fmt.Errorf("응답 본문 읽기 실패: %w", err)
		}

		// Handle HTTP errors.
		if resp.StatusCode >= 400 {
			apiErr := ParseAPIError(respBody, resp.StatusCode)
			return apiErr
		}

		// Deserialize result.
		if result != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("응답 파싱 실패: %w", err)
			}
		}

		return nil
	})
}

// GET is a helper for GET requests.
func (c *Client) GET(ctx context.Context, path string, query map[string]string, result interface{}) error {
	return c.Request(ctx, http.MethodGet, path, query, nil, result)
}

// POST is a helper for POST requests.
func (c *Client) POST(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, path, nil, body, result)
}

// DELETE is a helper for DELETE requests.
func (c *Client) DELETE(ctx context.Context, path string, query map[string]string, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, path, query, nil, result)
}

// GETWithRawQuery is a GET helper that supports array query parameters.
// rawQuery: unencoded query string (e.g. "uuids[]=uuid1&uuids[]=uuid2").
func (c *Client) GETWithRawQuery(ctx context.Context, path string, rawQuery string, result interface{}) error {
	return c.RequestWithRawQuery(ctx, http.MethodGet, path, rawQuery, result)
}

// DELETEWithRawQuery is a DELETE helper that supports array query parameters.
// rawQuery: unencoded query string (e.g. "uuids[]=uuid1&uuids[]=uuid2").
func (c *Client) DELETEWithRawQuery(ctx context.Context, path string, rawQuery string, result interface{}) error {
	return c.RequestWithRawQuery(ctx, http.MethodDelete, path, rawQuery, result)
}

// RequestWithRawQuery executes an HTTP request supporting array query parameters.
// Supports array parameters (e.g. uuids[]=a&uuids[]=b) that cannot be expressed as map[string]string.
func (c *Client) RequestWithRawQuery(ctx context.Context, method, path string, rawQuery string, result interface{}) error {
	group := ratelimit.GroupFromPath(path)

	return retry.Retry(ctx, func() error {
		// Wait for rate limit token.
		if err := c.limiter.Wait(ctx, group); err != nil {
			return err
		}

		// Build URL.
		rawURL := c.baseURL + path
		if rawQuery != "" {
			rawURL += "?" + rawQuery
		}

		// Create HTTP request.
		req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
		if err != nil {
			return fmt.Errorf("HTTP 요청 생성 실패: %w", err)
		}

		// Set common headers.
		req.Header.Set("Accept", "application/json")

		// Attach auth token using the raw query string for hashing.
		if c.accessKey != "" && c.secretKey != "" {
			token, tokenErr := GenerateTokenFromRawQuery(c.accessKey, c.secretKey, rawQuery)
			if tokenErr != nil {
				return fmt.Errorf("JWT 토큰 생성 실패: %w", tokenErr)
			}
			req.Header.Set("Authorization", "Bearer "+token)
		}

		// Execute request.
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("HTTP 요청 실패: %w", err)
		}
		defer resp.Body.Close()

		// Dynamically update rate limiter from Remaining-Req header.
		if remainingReq := resp.Header.Get("Remaining-Req"); remainingReq != "" {
			c.limiter.UpdateFromHeader(remainingReq, group)
		}

		// Read response body.
		respBody, err := readResponseBody(resp)
		if err != nil {
			return fmt.Errorf("응답 본문 읽기 실패: %w", err)
		}

		// Handle HTTP errors.
		if resp.StatusCode >= 400 {
			apiErr := ParseAPIError(respBody, resp.StatusCode)
			return apiErr
		}

		// Deserialize result.
		if result != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("응답 파싱 실패: %w", err)
			}
		}

		return nil
	})
}

// jsonBodyToQueryString converts a JSON body to a query string preserving key order.
// Upbit POST auth requires hashing the body Key-Value pairs in their original order.
// e.g. {"market":"KRW-BTC","side":"bid"} → "market=KRW-BTC&side=bid"
func jsonBodyToQueryString(bodyBytes []byte) string {
	dec := json.NewDecoder(bytes.NewReader(bodyBytes))

	// Consume the '{' token.
	t, err := dec.Token()
	if err != nil {
		return ""
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return ""
	}

	var pairs []string
	for dec.More() {
		// Read key.
		keyToken, err := dec.Token()
		if err != nil {
			break
		}
		key, ok := keyToken.(string)
		if !ok {
			break
		}

		// Read value.
		valToken, err := dec.Token()
		if err != nil {
			break
		}

		// Convert value to string.
		var val string
		switch v := valToken.(type) {
		case string:
			val = v
		case float64:
			// Use integer representation if the value is a whole number.
			if v == float64(int64(v)) {
				val = fmt.Sprintf("%d", int64(v))
			} else {
				val = fmt.Sprintf("%v", v)
			}
		case bool:
			val = fmt.Sprintf("%v", v)
		case nil:
			continue // skip null values
		default:
			val = fmt.Sprintf("%v", v)
		}

		pairs = append(pairs, key+"="+val)
	}

	return strings.Join(pairs, "&")
}

// isQuotationPath reports whether the path belongs to the Quotation API (market data).
// Upbit Quotation API paths: /trading_pairs, /tickers, /ticker, /candles/, /orderbooks, /trades/
func isQuotationPath(path string) bool {
	quotationPrefixes := []string{
		"/trading_pairs",
		"/tickers",
		"/ticker",
		"/candles/",
		"/orderbooks",
		"/trades/",
	}
	for _, prefix := range quotationPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// readResponseBody reads the response body, transparently decompressing gzip if needed.
func readResponseBody(resp *http.Response) ([]byte, error) {
	var reader io.Reader = resp.Body

	// Handle Content-Encoding: gzip.
	if strings.EqualFold(resp.Header.Get("Content-Encoding"), "gzip") {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip decompression failed: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	return io.ReadAll(reader)
}
