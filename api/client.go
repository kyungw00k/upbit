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

// Client Upbit API HTTP 클라이언트
type Client struct {
	baseURL    string
	accessKey  string
	secretKey  string
	httpClient *http.Client
	limiter    *ratelimit.Limiter
}

// NewClient 기본 클라이언트 생성
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

// NewClientWithURL 커스텀 Base URL로 클라이언트 생성 (테스트/목업용)
func NewClientWithURL(baseURL, accessKey, secretKey string) *Client {
	c := NewClient(accessKey, secretKey)
	c.baseURL = baseURL
	return c
}

// Request HTTP 요청 수행 (Rate Limit, 재시도, 에러 파싱 통합)
func (c *Client) Request(ctx context.Context, method, path string, query map[string]string, body interface{}, result interface{}) error {
	group := ratelimit.GroupFromPath(path)

	return retry.Retry(ctx, func() error {
		// Rate limit 대기
		if err := c.limiter.Wait(ctx, group); err != nil {
			return err
		}

		// URL 구성
		rawURL := c.baseURL + path
		u, err := url.Parse(rawURL)
		if err != nil {
			return fmt.Errorf("URL 파싱 실패: %w", err)
		}

		// 쿼리 파라미터 설정
		if len(query) > 0 {
			q := u.Query()
			for k, v := range query {
				q.Set(k, v)
			}
			u.RawQuery = q.Encode()
		}

		// 요청 본문 직렬화
		var bodyReader io.Reader
		var bodyParams map[string]string
		if body != nil {
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("요청 본문 직렬화 실패: %w", err)
			}
			bodyReader = bytes.NewReader(bodyBytes)

			// POST 인증용: JSON body에서 키 순서를 보존한 query string 생성
			// Upbit 인증 문서: JSON body의 Key-Value 순서를 유지하여 query string으로 변환
			if method == http.MethodPost {
				bodyQueryString := jsonBodyToQueryString(bodyBytes)
				if bodyQueryString != "" {
					bodyParams = map[string]string{"__raw_query__": bodyQueryString}
				}
			}
		}

		// HTTP 요청 생성
		req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
		if err != nil {
			return fmt.Errorf("HTTP 요청 생성 실패: %w", err)
		}

		// 공통 헤더 설정
		req.Header.Set("Accept", "application/json")
		// gzip 압축은 Quotation API (시세 조회)에만 적용
		if isQuotationPath(path) {
			req.Header.Set("Accept-Encoding", "gzip")
		}
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		// 인증 토큰 추가 (키가 설정된 경우)
		if c.accessKey != "" && c.secretKey != "" {
			var token string
			var tokenErr error

			if method == http.MethodPost && bodyParams != nil {
				// POST: JSON body 키 순서를 보존한 raw query string으로 해시
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

		// 요청 수행
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("HTTP 요청 실패: %w", err)
		}
		defer resp.Body.Close()

		// Remaining-Req 헤더로 Rate Limiter 동적 갱신
		if remainingReq := resp.Header.Get("Remaining-Req"); remainingReq != "" {
			c.limiter.UpdateFromHeader(remainingReq, group)
		}

		// 응답 본문 읽기 (gzip 압축 해제 지원)
		respBody, err := readResponseBody(resp)
		if err != nil {
			return fmt.Errorf("응답 본문 읽기 실패: %w", err)
		}

		// HTTP 에러 처리
		if resp.StatusCode >= 400 {
			apiErr := ParseAPIError(respBody, resp.StatusCode)
			return apiErr
		}

		// 결과 역직렬화
		if result != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("응답 파싱 실패: %w", err)
			}
		}

		return nil
	})
}

// GET GET 요청 헬퍼
func (c *Client) GET(ctx context.Context, path string, query map[string]string, result interface{}) error {
	return c.Request(ctx, http.MethodGet, path, query, nil, result)
}

// POST POST 요청 헬퍼
func (c *Client) POST(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, path, nil, body, result)
}

// DELETE DELETE 요청 헬퍼
func (c *Client) DELETE(ctx context.Context, path string, query map[string]string, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, path, query, nil, result)
}

// GETWithRawQuery 배열 쿼리 파라미터를 지원하는 GET 요청 헬퍼
// rawQuery: URL 인코딩되지 않은 쿼리 문자열 (예: "uuids[]=uuid1&uuids[]=uuid2")
func (c *Client) GETWithRawQuery(ctx context.Context, path string, rawQuery string, result interface{}) error {
	return c.RequestWithRawQuery(ctx, http.MethodGet, path, rawQuery, result)
}

// DELETEWithRawQuery 배열 쿼리 파라미터를 지원하는 DELETE 요청 헬퍼
// rawQuery: URL 인코딩되지 않은 쿼리 문자열 (예: "uuids[]=uuid1&uuids[]=uuid2")
func (c *Client) DELETEWithRawQuery(ctx context.Context, path string, rawQuery string, result interface{}) error {
	return c.RequestWithRawQuery(ctx, http.MethodDelete, path, rawQuery, result)
}

// RequestWithRawQuery 배열 쿼리 파라미터를 지원하는 HTTP 요청
// map[string]string으로 표현할 수 없는 배열 파라미터 (uuids[]=a&uuids[]=b 등) 지원
func (c *Client) RequestWithRawQuery(ctx context.Context, method, path string, rawQuery string, result interface{}) error {
	group := ratelimit.GroupFromPath(path)

	return retry.Retry(ctx, func() error {
		// Rate limit 대기
		if err := c.limiter.Wait(ctx, group); err != nil {
			return err
		}

		// URL 구성
		rawURL := c.baseURL + path
		if rawQuery != "" {
			rawURL += "?" + rawQuery
		}

		// HTTP 요청 생성
		req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
		if err != nil {
			return fmt.Errorf("HTTP 요청 생성 실패: %w", err)
		}

		// 공통 헤더 설정
		req.Header.Set("Accept", "application/json")

		// 인증 토큰 추가 (raw 쿼리 문자열로 해시 생성)
		if c.accessKey != "" && c.secretKey != "" {
			token, tokenErr := GenerateTokenFromRawQuery(c.accessKey, c.secretKey, rawQuery)
			if tokenErr != nil {
				return fmt.Errorf("JWT 토큰 생성 실패: %w", tokenErr)
			}
			req.Header.Set("Authorization", "Bearer "+token)
		}

		// 요청 수행
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("HTTP 요청 실패: %w", err)
		}
		defer resp.Body.Close()

		// Remaining-Req 헤더로 Rate Limiter 동적 갱신
		if remainingReq := resp.Header.Get("Remaining-Req"); remainingReq != "" {
			c.limiter.UpdateFromHeader(remainingReq, group)
		}

		// 응답 본문 읽기
		respBody, err := readResponseBody(resp)
		if err != nil {
			return fmt.Errorf("응답 본문 읽기 실패: %w", err)
		}

		// HTTP 에러 처리
		if resp.StatusCode >= 400 {
			apiErr := ParseAPIError(respBody, resp.StatusCode)
			return apiErr
		}

		// 결과 역직렬화
		if result != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("응답 파싱 실패: %w", err)
			}
		}

		return nil
	})
}

// jsonBodyToQueryString JSON body를 키 순서를 보존하여 query string으로 변환
// Upbit POST 인증: JSON body의 Key-Value 순서를 유지하여 hash 생성 필요
// 예: {"market":"KRW-BTC","side":"bid"} → "market=KRW-BTC&side=bid"
func jsonBodyToQueryString(bodyBytes []byte) string {
	dec := json.NewDecoder(bytes.NewReader(bodyBytes))

	// '{' 토큰 소비
	t, err := dec.Token()
	if err != nil {
		return ""
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return ""
	}

	var pairs []string
	for dec.More() {
		// 키 읽기
		keyToken, err := dec.Token()
		if err != nil {
			break
		}
		key, ok := keyToken.(string)
		if !ok {
			break
		}

		// 값 읽기
		valToken, err := dec.Token()
		if err != nil {
			break
		}

		// 값을 문자열로 변환
		var val string
		switch v := valToken.(type) {
		case string:
			val = v
		case float64:
			// 정수면 정수로, 아니면 소수점 유지
			if v == float64(int64(v)) {
				val = fmt.Sprintf("%d", int64(v))
			} else {
				val = fmt.Sprintf("%v", v)
			}
		case bool:
			val = fmt.Sprintf("%v", v)
		case nil:
			continue // null 값은 스킵
		default:
			val = fmt.Sprintf("%v", v)
		}

		pairs = append(pairs, key+"="+val)
	}

	return strings.Join(pairs, "&")
}

// isQuotationPath Quotation API (시세 조회) 경로 여부 판단
// Upbit Quotation API: /trading_pairs, /tickers, /ticker, /candles/, /orderbooks, /trades/
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

// readResponseBody 응답 본문 읽기 (gzip 압축 해제 지원)
func readResponseBody(resp *http.Response) ([]byte, error) {
	var reader io.Reader = resp.Body

	// Content-Encoding: gzip 처리
	if strings.EqualFold(resp.Header.Get("Content-Encoding"), "gzip") {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip 압축 해제 실패: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	return io.ReadAll(reader)
}
