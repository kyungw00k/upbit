package trader

import (
	"context"
	"testing"
	"time"

	"github.com/kyungw00k/upbit/trader/config"
	"github.com/kyungw00k/upbit/trader/types"
)

func testPipelineConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Trading.Mode = "paper"
	cfg.Trading.Markets = []string{"KRW-BTC"}
	cfg.Trading.InitialBalance = 10_000_000
	return cfg
}

func TestPipelineHoldSignalNoExecution(t *testing.T) {
	holdAgent := NewMockAgent("technical-agent")
	holdAgent.Report = types.AnalysisReport{
		AgentID:   "technical-agent",
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalHold,
		Coins: []types.CoinAnalysis{
			{
				Symbol:     "KRW-BTC",
				Signal:     types.SignalHold,
				Confidence: 0.5,
				Reasoning:  "no clear signal",
			},
		},
		Summary: "HOLD",
	}

	rg := NewRiskGate(config.DefaultConfig().Risk)
	executor := NewPaperExecutor()
	pipeline := NewPipeline([]Agent{holdAgent}, rg, executor, testPipelineConfig())

	result, err := pipeline.Run(context.Background(), map[string]any{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestPipelineBuySignalApproved(t *testing.T) {
	buyAgent := NewMockAgent("glm-analyst")
	buyAgent.Report = types.AnalysisReport{
		AgentID:   "glm-analyst",
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalBuy,
		Coins: []types.CoinAnalysis{
			{
				Symbol:     "KRW-BTC",
				Signal:     types.SignalBuy,
				Confidence: 0.8,
				EntryPrice: 75_000_000,
				StopLoss:   73_000_000,
				TakeProfit: 80_000_000,
				Position:   0.001,
				Reasoning:  "bullish momentum",
			},
		},
		Summary: "BUY BTC",
	}

	rg := NewRiskGate(config.DefaultConfig().Risk)
	executor := NewPaperExecutor()
	pipeline := NewPipeline([]Agent{buyAgent}, rg, executor, testPipelineConfig())

	result, err := pipeline.Run(context.Background(), map[string]any{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestPipelineRiskGateRejection(t *testing.T) {
	// Agent produces a BUY signal with missing stop loss.
	// The pipeline should run without error but risk gate rejects the trade.
	buyAgent := NewMockAgent("glm-analyst")
	buyAgent.Report = types.AnalysisReport{
		AgentID:   "glm-analyst",
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalBuy,
		Coins: []types.CoinAnalysis{
			{
				Symbol:     "KRW-BTC",
				Signal:     types.SignalBuy,
				Confidence: 0.9,
				EntryPrice: 75_000_000,
				StopLoss:   0,
				TakeProfit: 80_000_000,
				Position:   0.001,
				Reasoning:  "missing stop loss triggers rejection",
			},
		},
		Summary: "BUY BTC",
	}

	rg := NewRiskGate(config.DefaultConfig().Risk)
	executor := NewPaperExecutor()
	pipeline := NewPipeline([]Agent{buyAgent}, rg, executor, testPipelineConfig())

	result, err := pipeline.Run(context.Background(), map[string]any{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestPipelineEmptyAgents(t *testing.T) {
	rg := NewRiskGate(config.DefaultConfig().Risk)
	executor := NewPaperExecutor()
	pipeline := NewPipeline([]Agent{}, rg, executor, testPipelineConfig())

	result, err := pipeline.Run(context.Background(), map[string]any{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestPipelineAllAgentsFail(t *testing.T) {
	failAgent := NewMockAgent("failing-agent")
	failAgent.Err = context.DeadlineExceeded

	rg := NewRiskGate(config.DefaultConfig().Risk)
	executor := NewPaperExecutor()
	pipeline := NewPipeline([]Agent{failAgent}, rg, executor, testPipelineConfig())

	_, err := pipeline.Run(context.Background(), map[string]any{}, nil)
	if err == nil {
		t.Fatal("expected error when all agents fail")
	}
}

func TestPipelineTwoAgents(t *testing.T) {
	technicalAgent := NewMockAgent("technical-agent")
	technicalAgent.Report = types.AnalysisReport{
		AgentID:   "technical-agent",
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalBuy,
		Coins: []types.CoinAnalysis{
			{
				Symbol:     "KRW-BTC",
				Signal:     types.SignalBuy,
				Confidence: 0.7,
				EntryPrice: 75_000_000,
				StopLoss:   73_000_000,
				TakeProfit: 80_000_000,
				Position:   0.001,
				Reasoning:  "RSI oversold",
			},
		},
		Summary: "Technical BUY",
	}

	glmAgent := NewMockAgent("glm-analyst")
	glmAgent.Report = types.AnalysisReport{
		AgentID:   "glm-analyst",
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalStrongBuy,
		Coins: []types.CoinAnalysis{
			{
				Symbol:     "KRW-BTC",
				Signal:     types.SignalStrongBuy,
				Confidence: 0.9,
				EntryPrice: 75_000_000,
				StopLoss:   73_000_000,
				TakeProfit: 82_000_000,
				Position:   0.002,
				Reasoning:  "GLM bullish",
			},
		},
		Summary: "GLM STRONG_BUY",
	}

	rg := NewRiskGate(config.DefaultConfig().Risk)
	executor := NewPaperExecutor()
	pipeline := NewPipeline([]Agent{technicalAgent, glmAgent}, rg, executor, testPipelineConfig())

	result, err := pipeline.Run(context.Background(), map[string]any{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestSignalToValueRoundTrip(t *testing.T) {
	tests := []struct {
		signal types.Signal
		value  float64
	}{
		{types.SignalStrongBuy, SignalValueStrongBuy},
		{types.SignalBuy, SignalValueBuy},
		{types.SignalHold, SignalValueHold},
		{types.SignalSell, SignalValueSell},
		{types.SignalStrongSell, SignalValueStrongSell},
	}
	for _, tt := range tests {
		v := signalToValue(tt.signal)
		if v != tt.value {
			t.Errorf("signalToValue(%s) = %f, want %f", tt.signal, v, tt.value)
		}
	}
}

func TestValueToSignal(t *testing.T) {
	tests := []struct {
		value  float64
		signal types.Signal
	}{
		{2.0, types.SignalStrongBuy},
		{1.5, types.SignalStrongBuy},
		{1.0, types.SignalBuy},
		{0.5, types.SignalBuy},
		{0.0, types.SignalHold},
		{-0.5, types.SignalSell},
		{-1.0, types.SignalSell},
		{-1.5, types.SignalStrongSell},
		{-2.0, types.SignalStrongSell},
	}
	for _, tt := range tests {
		s := valueToSignal(tt.value)
		if s != tt.signal {
			t.Errorf("valueToSignal(%f) = %s, want %s", tt.value, s, tt.signal)
		}
	}
}
