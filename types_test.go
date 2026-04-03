package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSignalValid(t *testing.T) {
	tests := []struct {
		signal Signal
		valid  bool
	}{
		{SignalStrongBuy, true},
		{SignalBuy, true},
		{SignalHold, true},
		{SignalSell, true},
		{SignalStrongSell, true},
		{Signal("INVALID"), false},
		{Signal(""), false},
	}
	for _, tt := range tests {
		if got := tt.signal.Valid(); got != tt.valid {
			t.Errorf("Signal(%q).Valid() = %v, want %v", tt.signal, got, tt.valid)
		}
	}
}

func TestCoinAnalysisMarshalRoundTrip(t *testing.T) {
	original := CoinAnalysis{
		Symbol: "KRW-BTC", Signal: SignalBuy, Confidence: 0.85,
		EntryPrice: 75000000.0, StopLoss: 72000000.0, TakeProfit: 82000000.0,
		Position: 500000.0, Reasoning: "RSI oversold bounce",
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var decoded CoinAnalysis
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if decoded.Symbol != original.Symbol || decoded.Signal != original.Signal {
		t.Errorf("round-trip mismatch")
	}
}

func TestAnalysisReportMarshalRoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	original := AnalysisReport{
		AgentID: "technical-analyst", Timestamp: now, Signal: SignalBuy,
		Coins: []CoinAnalysis{
			{Symbol: "KRW-BTC", Signal: SignalStrongBuy, Confidence: 0.9},
		},
		Summary: "BTC bullish",
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var decoded AnalysisReport
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if decoded.AgentID != original.AgentID || len(decoded.Coins) != 1 {
		t.Errorf("round-trip mismatch")
	}
}

func TestPipelineResultMarshalJSON(t *testing.T) {
	result := PipelineResult{RanAt: time.Now().UTC(), Duration: 1500 * time.Millisecond}
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if raw["duration_ms"].(float64) != 1500 {
		t.Errorf("expected duration_ms=1500, got %v", raw["duration_ms"])
	}
}
