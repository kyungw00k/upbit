package agent

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kyungw00k/upbit/trader/types"
)

func TestStripMarkdownFence(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain json", `{"signal":"HOLD"}`, `{"signal":"HOLD"}`},
		{"json fence", "```json\n{\"signal\":\"HOLD\"}\n```", `{"signal":"HOLD"}`},
		{"bare fence", "```\n{\"signal\":\"HOLD\"}\n```", `{"signal":"HOLD"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripMarkdownFence(tt.input)
			if got != tt.want {
				t.Errorf("stripMarkdownFence() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGLMAgent_NormalResponse(t *testing.T) {
	respBody := glmResponse{
		Signal:   "STRONG_BUY",
		Summary:  "Bullish momentum detected",
		Coins:    []glmCoin{{Symbol: "KRW-BTC", Signal: "STRONG_BUY", Confidence: 0.85, Reasoning: "Oversold RSI bounce"}},
	}
	respJSON, _ := json.Marshal(respBody)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		json.Unmarshal(body, &req)

		// Verify JSON mode is requested.
		if rf, ok := req["response_format"].(map[string]any); ok {
			if rf["type"] != "json_object" {
				t.Errorf("expected response_format type json_object, got %v", rf["type"])
			}
		}

		resp := map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]any{
						"content": string(respJSON),
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	agent := NewGLMAgent("test-key", server.URL, "glm-4.7")
	report, err := agent.Run(t.Context(), types.AgentInput{
		Markets:    []string{"KRW-BTC"},
		Indicators: map[string]any{"rsi": 28.5},
	})
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if report.Signal != types.SignalStrongBuy {
		t.Errorf("Signal = %q, want STRONG_BUY", report.Signal)
	}
	if report.AgentID != "glm-analyst" {
		t.Errorf("AgentID = %q, want glm-analyst", report.AgentID)
	}
	if len(report.Coins) != 1 || report.Coins[0].Symbol != "KRW-BTC" {
		t.Errorf("Coins = %v, want 1 coin KRW-BTC", report.Coins)
	}
	if report.Coins[0].Confidence != 0.85 {
		t.Errorf("Confidence = %v, want 0.85", report.Coins[0].Confidence)
	}
}

func TestGLMAgent_MarkdownFenceResponse(t *testing.T) {
	inner := `{"signal":"HOLD","confidence":0.5,"reasoning":"mixed signals","coins":[{"symbol":"KRW-ETH","signal":"HOLD","confidence":0.5,"reasoning":"neutral"}],"summary":"Neutral market"}`
	wrapped := "```json\n" + inner + "\n```"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"content": wrapped}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	agent := NewGLMAgent("test-key", server.URL, "glm-4.7")
	report, err := agent.Run(t.Context(), types.AgentInput{Markets: []string{"KRW-ETH"}})
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if report.Signal != types.SignalHold {
		t.Errorf("Signal = %q, want HOLD", report.Signal)
	}
}

func TestGLMAgent_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{"choices": []map[string]any{}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	agent := NewGLMAgent("test-key", server.URL, "glm-4.7")
	_, err := agent.Run(t.Context(), types.AgentInput{Markets: []string{"KRW-BTC"}})
	if err == nil {
		t.Fatal("expected error for empty response")
	}
	if !strings.Contains(err.Error(), "no choices") {
		t.Errorf("error = %q, want 'no choices'", err.Error())
	}
}

func TestGLMAgent_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"content": "this is not json"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	agent := NewGLMAgent("test-key", server.URL, "glm-4.7")
	_, err := agent.Run(t.Context(), types.AgentInput{Markets: []string{"KRW-BTC"}})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "parsing glm response") {
		t.Errorf("error = %q, want 'parsing glm response'", err.Error())
	}
}

func TestBuildAnalysisPrompt(t *testing.T) {
	indicators := map[string]any{
		"rsi":           45.3,
		"trend":         "neutral",
		"volume_ratio":  1.2,
		"macd":          map[string]any{"line": 100.0, "signal": 95.0, "histogram": 5.0},
		"bollinger_bands": map[string]any{"upper": 80000.0, "middle": 75000.0, "lower": 70000.0},
		"ema":           map[string]any{"20": 74500.0, "50": 74000.0},
		"atr":           1500.0,
	}
	sentiment := &SentimentInfo{Index: 55, Category: "Greed"}

	prompt := BuildAnalysisPrompt("KRW-BTC", indicators, sentiment)

	checks := []string{
		"KRW-BTC",
		"45.3",
		"neutral",
		"1.2",
		"100",   // MACD line
		"80000", // BB upper
		"74500", // EMA 20
		"1500",  // ATR
		"55",    // Fear & Greed
		"Greed",
	}
	for _, c := range checks {
		if !strings.Contains(prompt, c) {
			t.Errorf("prompt missing %q\n%s", c, prompt)
		}
	}
}

func TestClampConfidence(t *testing.T) {
	tests := []struct {
		in, want float64
	}{
		{-0.5, 0},
		{0.0, 0},
		{0.5, 0.5},
		{1.0, 1.0},
		{1.5, 1.0},
	}
	for _, tt := range tests {
		got := clampConfidence(tt.in)
		if got != tt.want {
			t.Errorf("clampConfidence(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
