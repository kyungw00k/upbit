package sentiment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const defaultBaseURL = "https://api.alternative.me/fng/"

// FearGreedData represents a single Fear & Greed Index data point.
type FearGreedData struct {
	Value          int    `json:"value"`
	Classification string `json:"value_classification"`
	Timestamp      string `json:"timestamp"`
}

// apiResponse mirrors the JSON structure returned by the Alternative.me API.
type apiResponse struct {
	Data []struct {
		Value          string `json:"value"`
		Classification string `json:"value_classification"`
		Timestamp      string `json:"timestamp"`
	} `json:"data"`
}

// FetchFearGreedIndex fetches the latest `limit` Fear & Greed Index entries
// from the Alternative.me API.
func FetchFearGreedIndex(ctx context.Context, limit int) ([]FearGreedData, error) {
	url := fmt.Sprintf("%s?limit=%d&format=json", defaultBaseURL, limit)
	return fetch(ctx, url)
}

// FetchLatestFearGreedIndex fetches only the most recent Fear & Greed Index entry.
func FetchLatestFearGreedIndex(ctx context.Context) (*FearGreedData, error) {
	results, err := FetchFearGreedIndex(ctx, 1)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no fear & greed data returned")
	}
	return &results[0], nil
}

// fetch performs the HTTP GET and parses the API response.
func fetch(ctx context.Context, url string) ([]FearGreedData, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if len(apiResp.Data) == 0 {
		return nil, fmt.Errorf("empty data array in response")
	}

	results := make([]FearGreedData, 0, len(apiResp.Data))
	for i, entry := range apiResp.Data {
		value, err := strconv.Atoi(entry.Value)
		if err != nil {
			return nil, fmt.Errorf("parsing value at index %d: %q is not a valid integer: %w", i, entry.Value, err)
		}
		results = append(results, FearGreedData{
			Value:          value,
			Classification: entry.Classification,
			Timestamp:      entry.Timestamp,
		})
	}

	return results, nil
}
