package agents

import (
	"testing"

	"github.com/kyungw00k/upbit/trader/types"
)

func TestTechnicalAgent_OverboughtSELL(t *testing.T) {
	agent := NewTechnicalAgent()
	input := types.AgentInput{
		Markets: []string{"KRW-BTC"},
		Indicators: map[string]any{
			"rsi":  80.0,
			"macd": map[string]any{"line": 100.0, "signal": 110.0, "histogram": -10.0},
			"bollinger_bands": map[string]any{"upper": 80000.0, "middle": 75000.0, "lower": 70000.0},
			"ema":   map[string]any{"20": 74000.0, "50": 73000.0},
			"close": 79500.0,
		},
	}

	report, err := agent.Run(t.Context(), input)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if report.Signal != types.SignalSell && report.Signal != types.SignalStrongSell {
		t.Errorf("Signal = %q, want SELL or STRONG_SELL", report.Signal)
	}
	if report.AgentID != "technical-agent" {
		t.Errorf("AgentID = %q, want technical-agent", report.AgentID)
	}
}

func TestTechnicalAgent_OversoldBUY(t *testing.T) {
	agent := NewTechnicalAgent()
	input := types.AgentInput{
		Markets: []string{"KRW-BTC"},
		Indicators: map[string]any{
			"rsi":  20.0,
			"macd": map[string]any{"line": 110.0, "signal": 100.0, "histogram": 10.0},
			"bollinger_bands": map[string]any{"upper": 80000.0, "middle": 75000.0, "lower": 70000.0},
			"ema":   map[string]any{"20": 76000.0, "50": 77000.0},
			"close": 70500.0,
		},
	}

	report, err := agent.Run(t.Context(), input)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if report.Signal != types.SignalBuy && report.Signal != types.SignalStrongBuy {
		t.Errorf("Signal = %q, want BUY or STRONG_BUY", report.Signal)
	}
}

func TestTechnicalAgent_NeutralHOLD(t *testing.T) {
	agent := NewTechnicalAgent()
	input := types.AgentInput{
		Markets: []string{"KRW-BTC"},
		Indicators: map[string]any{
			"rsi":  50.0,
			"macd": map[string]any{"line": 100.0, "signal": 100.0, "histogram": 0.0},
			"bollinger_bands": map[string]any{"upper": 80000.0, "middle": 75000.0, "lower": 70000.0},
			"ema":   map[string]any{"20": 74999.0, "50": 75001.0},
			"close": 75000.0,
		},
	}

	report, err := agent.Run(t.Context(), input)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if report.Signal != types.SignalHold {
		t.Errorf("Signal = %q, want HOLD", report.Signal)
	}
}

func TestTechnicalAgent_EmptyIndicators(t *testing.T) {
	agent := NewTechnicalAgent()
	input := types.AgentInput{
		Markets:    []string{"KRW-BTC"},
		Indicators: map[string]any{},
	}

	report, err := agent.Run(t.Context(), input)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if report.Signal != types.SignalHold {
		t.Errorf("Signal = %q, want HOLD for empty indicators", report.Signal)
	}
	if report.Coins[0].Confidence != 0 {
		t.Errorf("Confidence = %v, want 0 for empty indicators", report.Coins[0].Confidence)
	}
}

func TestScoreRSI(t *testing.T) {
	tests := []struct {
		rsi  float64
		want float64
	}{
		{20, 30},    // Deep oversold → max bullish
		{30, 30},    // Oversold threshold
		{50, 0},     // Neutral
		{70, -30},   // Overbought threshold
		{80, -30},   // Deep overbought → max bearish
	}
	for _, tt := range tests {
		got := scoreRSI(tt.rsi)
		if got != tt.want {
			t.Errorf("scoreRSI(%v) = %v, want %v", tt.rsi, got, tt.want)
		}
	}
}

func TestScoreMACD(t *testing.T) {
	tests := []struct {
		hist float64
		want float64
	}{
		{5.0, 25},
		{-5.0, -25},
		{0, 0},
	}
	for _, tt := range tests {
		got := scoreMACD(tt.hist)
		if got != tt.want {
			t.Errorf("scoreMACD(%v) = %v, want %v", tt.hist, got, tt.want)
		}
	}
}

func TestScoreToSignal(t *testing.T) {
	tests := []struct {
		score float64
		want  types.Signal
	}{
		{80, types.SignalStrongBuy},
		{40, types.SignalBuy},
		{5, types.SignalHold},
		{-30, types.SignalSell},
		{-70, types.SignalStrongSell},
	}
	for _, tt := range tests {
		got := scoreToSignal(tt.score)
		if got != tt.want {
			t.Errorf("scoreToSignal(%v) = %v, want %v", tt.score, got, tt.want)
		}
	}
}
