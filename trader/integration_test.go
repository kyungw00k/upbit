package trader

import (
	"context"
	"testing"
	"time"

	"github.com/kyungw00k/upbit/trader/agent"
	"github.com/kyungw00k/upbit/trader/agents"
	"github.com/kyungw00k/upbit/trader/config"
	"github.com/kyungw00k/upbit/trader/indicator"
	"github.com/kyungw00k/upbit/trader/types"
)

// TestIntegration_IndicatorToAgentToRiskToExecution verifies the full pipeline:
// indicator calculation → technical agent → risk gate → paper executor.
func TestIntegration_IndicatorToAgentToRiskToExecution(t *testing.T) {
	candles := generateUptrendCandles(50)

	// Step 1: Calculate indicators
	result := indicator.CalculateAll(candles)
	if result.RSI < 0 || result.RSI > 100 {
		t.Fatalf("RSI out of range: %f", result.RSI)
	}
	t.Logf("Indicators: RSI=%.1f MACD_Hist=%.4f Trend=%s BB=[%.0f,%.0f,%.0f]",
		result.RSI, result.MACDHist, result.Trend, result.BBLower, result.BBMiddle, result.BBUpper)

	// Step 2: Run technical agent
	techAgent := agents.NewTechnicalAgent()
	input := types.AgentInput{
		Markets: []string{"KRW-BTC"},
		Indicators: map[string]any{
			"KRW-BTC": map[string]float64{
				"rsi":            result.RSI,
				"macd_histogram": result.MACDHist,
				"bb_upper":       result.BBUpper,
				"bb_middle":      result.BBMiddle,
				"bb_lower":       result.BBLower,
				"ema20":          result.EMA20,
				"ema50":          result.EMA50,
				"atr":            result.ATR,
				"vol_ratio":      result.VolRatio,
				"close":          candles[len(candles)-1].Close,
			},
		},
	}

	report, err := techAgent.Run(context.Background(), input)
	if err != nil {
		t.Fatalf("TechnicalAgent.Run failed: %v", err)
	}
	t.Logf("TechnicalAgent: signal=%s coins=%d summary=%s",
		report.Signal, len(report.Coins), report.Summary)

	// Step 3: Risk Gate
	riskGate := NewRiskGate(config.DefaultConfig().Risk)
	portfolio := Portfolio{
		TotalValue: 10_000_000,
		Positions:  map[string]Position{},
		DailyPnL:   0,
	}

	for _, coin := range report.Coins {
		approved, adjustedSize, reason := riskGate.Evaluate(coin, portfolio)
		t.Logf("RiskGate: symbol=%s approved=%v size=%.0f reason=%s",
			coin.Symbol, approved, adjustedSize, reason)

		if approved && coin.EntryPrice > 0 && adjustedSize > 0 {
			executor := &PaperExecutor{}
			tr, err := executor.Execute(context.Background(), TradeRequest{
				Market: coin.Symbol,
				Side:   "bid",
				Price:  coin.EntryPrice,
				Volume: adjustedSize / coin.EntryPrice,
				Mode:   "paper",
			})
			if err != nil {
				t.Errorf("PaperExecutor failed: %v", err)
			} else {
				t.Logf("Trade: uuid=%s filled=%.0f fee=%.2f",
					tr.OrderUUID, tr.FilledPrice, tr.Fee)
			}
		}
	}
}

// TestIntegration_PipelineWithMultipleAgents verifies signal merging.
func TestIntegration_PipelineWithMultipleAgents(t *testing.T) {
	cfg := config.DefaultConfig()
	riskGate := NewRiskGate(cfg.Risk)
	executor := &PaperExecutor{}

	techAgent := agents.NewTechnicalAgent()
	mockGLM := NewMockAgent("glm-analyst")
	mockGLM.Report = types.AnalysisReport{
		AgentID:   "glm-analyst",
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalBuy,
		Coins: []types.CoinAnalysis{
			{
				Symbol:     "KRW-BTC",
				Signal:     types.SignalBuy,
				Confidence: 0.75,
				EntryPrice: 75_000_000,
				StopLoss:   72_000_000,
				TakeProfit: 82_000_000,
				Position:   500_000,
				Reasoning:  "Bullish momentum",
			},
		},
		Summary: "BTC showing bullish signals",
	}

	pipeline := NewPipeline([]Agent{techAgent, mockGLM}, riskGate, executor, cfg)

	marketData := map[string]any{
		"KRW-BTC": map[string]float64{
			"rsi": 35, "macd_histogram": 50, "bb_upper": 78_000_000,
			"bb_middle": 75_000_000, "bb_lower": 72_000_000,
			"ema20": 74_500_000, "ema50": 73_000_000,
			"atr": 800_000, "vol_ratio": 1.8, "close": 75_000_000,
		},
	}

	result, err := pipeline.Run(context.Background(), marketData, nil)
	if err != nil {
		t.Fatalf("Pipeline.Run failed: %v", err)
	}
	t.Logf("Pipeline: duration=%v", result.Duration)
}

// TestIntegration_DaemonStartStop verifies daemon lifecycle.
func TestIntegration_DaemonStartStop(t *testing.T) {
	cfg := config.DefaultConfig()
	riskGate := NewRiskGate(cfg.Risk)
	executor := &PaperExecutor{}

	mockAgent := NewMockAgent("test")
	mockAgent.Report = types.AnalysisReport{
		AgentID:   "test",
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalHold,
		Summary:   "Integration test",
	}

	pipeline := NewPipeline([]Agent{mockAgent}, riskGate, executor, cfg)
	daemon := NewDaemon(pipeline, 1*time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- daemon.Start(ctx) }()

	time.Sleep(100 * time.Millisecond)
	daemon.Stop()
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("daemon error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("daemon did not stop")
	}
	t.Log("Daemon lifecycle OK")
}

// TestIntegration_PromptBuilderWithRealIndicators verifies GLM prompt
// generation with actual indicator data.
func TestIntegration_PromptBuilderWithRealIndicators(t *testing.T) {
	candles := generateUptrendCandles(50)
	result := indicator.CalculateAll(candles)

	indicators := map[string]any{
		"rsi":            result.RSI,
		"macd_line":      result.MACDLine,
		"macd_signal":    result.MACDSignal,
		"macd_histogram": result.MACDHist,
		"bb_upper":       result.BBUpper,
		"bb_middle":      result.BBMiddle,
		"bb_lower":       result.BBLower,
		"ema20":          result.EMA20,
		"ema50":          result.EMA50,
		"atr":            result.ATR,
		"vol_ratio":      result.VolRatio,
		"trend":          result.Trend,
	}

	sentiment := &agent.SentimentInfo{
		Index:    45,
		Category: "Fear",
	}

	prompt := agent.BuildAnalysisPrompt("KRW-BTC", indicators, sentiment)
	if len(prompt) < 100 {
		t.Errorf("Prompt too short: %d chars", len(prompt))
	}
	t.Logf("Prompt: %d chars — %s", len(prompt), truncate(prompt, 150))
}

// --- Helpers ---

func generateUptrendCandles(n int) []indicator.CandleData {
	candles := make([]indicator.CandleData, n)
	basePrice := 70_000_000.0
	for i := 0; i < n; i++ {
		open := basePrice + float64(i)*100_000
		close_ := open + 150_000
		high := close_ + 50_000
		low := open - 50_000
		volume := 1000.0 + float64(i)*10
		candles[i] = indicator.CandleData{
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close_,
			Volume: volume,
			Time:   time.Now().Add(-time.Duration(n-i) * time.Hour).Format(time.RFC3339),
		}
	}
	return candles
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
