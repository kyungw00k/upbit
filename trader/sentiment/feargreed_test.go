package sentiment

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchLatestFearGreedIndex_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify limit parameter is present
		if r.URL.Query().Get("limit") != "1" {
			t.Errorf("expected limit=1, got %s", r.URL.Query().Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"value": "45",
					"value_classification": "Fear",
					"timestamp": "1715299200"
				}
			]
		}`))
	}))
	defer server.Close()

	result, err := fetch(context.Background(), server.URL+"?limit=1&format=json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].Value != 45 {
		t.Errorf("expected value 45, got %d", result[0].Value)
	}

	if result[0].Classification != "Fear" {
		t.Errorf("expected classification Fear, got %s", result[0].Classification)
	}

	if result[0].Timestamp != "1715299200" {
		t.Errorf("expected timestamp 1715299200, got %s", result[0].Timestamp)
	}
}

func TestFetchFearGreedIndex_MultipleEntries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "3" {
			t.Errorf("expected limit=3, got %s", r.URL.Query().Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"value": "75",
					"value_classification": "Greed",
					"timestamp": "1715299200"
				},
				{
					"value": "50",
					"value_classification": "Neutral",
					"timestamp": "1715212800"
				},
				{
					"value": "25",
					"value_classification": "Extreme Fear",
					"timestamp": "1715126400"
				}
			]
		}`))
	}))
	defer server.Close()

	results, err := fetch(context.Background(), server.URL+"?limit=3&format=json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	if results[0].Value != 75 || results[0].Classification != "Greed" {
		t.Errorf("first entry: expected 75/Greed, got %d/%s", results[0].Value, results[0].Classification)
	}

	if results[1].Value != 50 || results[1].Classification != "Neutral" {
		t.Errorf("second entry: expected 50/Neutral, got %d/%s", results[1].Value, results[1].Classification)
	}

	if results[2].Value != 25 || results[2].Classification != "Extreme Fear" {
		t.Errorf("third entry: expected 25/Extreme Fear, got %d/%s", results[2].Value, results[2].Classification)
	}
}

func TestFetchFearGreedIndex_EmptyData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": []}`))
	}))
	defer server.Close()

	_, err := fetch(context.Background(), server.URL+"?limit=1&format=json")
	if err == nil {
		t.Fatal("expected error for empty data, got nil")
	}
}

func TestFetchFearGreedIndex_InvalidValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"value": "not_a_number",
					"value_classification": "Fear",
					"timestamp": "1715299200"
				}
			]
		}`))
	}))
	defer server.Close()

	_, err := fetch(context.Background(), server.URL+"?limit=1&format=json")
	if err == nil {
		t.Fatal("expected error for invalid value, got nil")
	}
}

func TestFetchFearGreedIndex_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := fetch(context.Background(), server.URL+"?limit=1&format=json")
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
}

func TestFetchFearGreedIndex_LimitParameter(t *testing.T) {
	var capturedLimit string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedLimit = r.URL.Query().Get("limit")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"value": "45",
					"value_classification": "Fear",
					"timestamp": "1715299200"
				}
			]
		}`))
	}))
	defer server.Close()

	_, err := fetch(context.Background(), server.URL+"?limit=10&format=json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedLimit != "10" {
		t.Errorf("expected limit parameter '10', got '%s'", capturedLimit)
	}
}

func TestFetchFearGreedIndex_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := fetch(ctx, server.URL+"?limit=1&format=json")
	if err == nil {
		t.Fatal("expected error due to context cancellation, got nil")
	}
}

func TestFetchLatestFearGreedIndex_NoData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": []}`))
	}))
	defer server.Close()

	_, err := fetch(context.Background(), server.URL+"?limit=1&format=json")
	if err == nil {
		t.Fatal("expected error for no data, got nil")
	}
}
