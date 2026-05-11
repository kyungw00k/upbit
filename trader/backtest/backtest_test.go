package backtest

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/kyungw00k/upbit/trader"
	"github.com/kyungw00k/upbit/trader/indicator"
	"github.com/kyungw00k/upbit/trader/types"
)

// mockBuyAgent always returns BUY.
type mockBuyAgent struct{}

func (m *mockBuyAgent) ID() string { return "mock-buy-agent" }
func (m *mockBuyAgent) Run(_ context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	market := "KRW-BTC"
	if len(input.Markets) > 0 {
		market = input.Markets[0]
	}
	return types.AnalysisReport{
		AgentID:   m.ID(),
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalBuy,
		Coins: []types.CoinAnalysis{
			{
				Symbol:     market,
				Signal:     types.SignalBuy,
				Confidence: 0.8,
				Reasoning:  "mock buy",
			},
		},
		Summary: "BUY",
	}, nil
}

// mockSellAgent always returns SELL.
type mockSellAgent struct{}

func (m *mockSellAgent) ID() string { return "mock-sell-agent" }
func (m *mockSellAgent) Run(_ context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	market := "KRW-BTC"
	if len(input.Markets) > 0 {
		market = input.Markets[0]
	}
	return types.AnalysisReport{
		AgentID:   m.ID(),
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalSell,
		Coins: []types.CoinAnalysis{
			{
				Symbol:     market,
				Signal:     types.SignalSell,
				Confidence: 0.8,
				Reasoning:  "mock sell",
			},
		},
		Summary: "SELL",
	}, nil
}

// mockHoldAgent always returns HOLD.
type mockHoldAgent struct{}

func (m *mockHoldAgent) ID() string { return "mock-hold-agent" }
func (m *mockHoldAgent) Run(_ context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	market := "KRW-BTC"
	if len(input.Markets) > 0 {
		market = input.Markets[0]
	}
	return types.AnalysisReport{
		AgentID:   m.ID(),
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalHold,
		Coins: []types.CoinAnalysis{
			{
				Symbol:     market,
				Signal:     types.SignalHold,
				Confidence: 0.3,
				Reasoning:  "mock hold",
			},
		},
		Summary: "HOLD",
	}, nil
}

// generateUptrendCandles creates candle data in a steady uptrend.
func generateUptrendCandles(n int) []indicator.CandleData {
	candles := make([]indicator.CandleData, n)
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < n; i++ {
		price := 1_000_000 + float64(i)*5_000
		candles[i] = indicator.CandleData{
			Open:   price - 1000,
			High:   price + 2000,
			Low:    price - 2000,
			Close:  price,
			Volume: 100 + float64(i),
			Time:   baseTime.Add(time.Duration(i) * time.Hour).Format("2006-01-02T15:04:05"),
		}
	}
	return candles
}

// generateDowntrendCandles creates candle data in a steady downtrend.
func generateDowntrendCandles(n int) []indicator.CandleData {
	candles := make([]indicator.CandleData, n)
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < n; i++ {
		price := 5_000_000 - float64(i)*5_000
		if price < 100_000 {
			price = 100_000
		}
		candles[i] = indicator.CandleData{
			Open:   price + 1000,
			High:   price + 2000,
			Low:    price - 2000,
			Close:  price,
			Volume: 100 + float64(i),
			Time:   baseTime.Add(time.Duration(i) * time.Hour).Format("2006-01-02T15:04:05"),
		}
	}
	return candles
}

func TestBacktestUptrendPositiveReturn(t *testing.T) {
	rg := trader.NewRiskGate(DefaultRiskConfig())
	agent := &mockBuyAgent{}
	engine := NewEngine(agent, rg)

	candles := map[string][]indicator.CandleData{
		"KRW-BTC": generateUptrendCandles(200),
	}

	result, err := engine.Run(context.Background(), candles, 10_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.TotalTrades == 0 {
		t.Error("expected at least one trade in uptrend with buy agent")
	}
	// In an uptrend with a buy agent, we expect positive total return.
	if result.TotalReturn <= 0 {
		t.Errorf("expected positive TotalReturn in uptrend, got %.4f", result.TotalReturn)
	}
}

func TestBacktestDowntrendWithBuyAgent(t *testing.T) {
	rg := trader.NewRiskGate(DefaultRiskConfig())
	agent := &mockBuyAgent{}
	engine := NewEngine(agent, rg)

	candles := map[string][]indicator.CandleData{
		"KRW-BTC": generateDowntrendCandles(200),
	}

	result, err := engine.Run(context.Background(), candles, 10_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// In a downtrend buying should lead to losses or at least very few trades.
	if result.TotalTrades > 0 && result.TotalReturn > 0 {
		t.Logf("downtrend with buy agent: TotalReturn=%.4f, trades=%d", result.TotalReturn, result.TotalTrades)
	}
}

func TestBacktestHoldAgentNoTrades(t *testing.T) {
	rg := trader.NewRiskGate(DefaultRiskConfig())
	agent := &mockHoldAgent{}
	engine := NewEngine(agent, rg)

	candles := map[string][]indicator.CandleData{
		"KRW-BTC": generateUptrendCandles(100),
	}

	result, err := engine.Run(context.Background(), candles, 10_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalTrades != 0 {
		t.Errorf("expected 0 trades with hold agent, got %d", result.TotalTrades)
	}
}

func TestBacktestMaxDrawdownNonNegative(t *testing.T) {
	rg := trader.NewRiskGate(DefaultRiskConfig())
	agent := &mockBuyAgent{}
	engine := NewEngine(agent, rg)

	candles := map[string][]indicator.CandleData{
		"KRW-BTC": generateUptrendCandles(200),
	}

	result, err := engine.Run(context.Background(), candles, 10_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MaxDrawdown < 0 {
		t.Errorf("MaxDrawdown should be >= 0, got %.4f", result.MaxDrawdown)
	}
	if result.MaxDrawdown > 1 {
		t.Errorf("MaxDrawdown should be <= 1, got %.4f", result.MaxDrawdown)
	}
}

func TestBacktestWinRateInRange(t *testing.T) {
	rg := trader.NewRiskGate(DefaultRiskConfig())
	agent := &mockBuyAgent{}
	engine := NewEngine(agent, rg)

	candles := map[string][]indicator.CandleData{
		"KRW-BTC": generateUptrendCandles(200),
	}

	result, err := engine.Run(context.Background(), candles, 10_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.WinRate < 0 || result.WinRate > 1 {
		t.Errorf("WinRate should be between 0 and 1, got %.4f", result.WinRate)
	}
}

func TestBacktestInsufficientCandles(t *testing.T) {
	rg := trader.NewRiskGate(DefaultRiskConfig())
	agent := &mockBuyAgent{}
	engine := NewEngine(agent, rg)

	candles := map[string][]indicator.CandleData{
		"KRW-BTC": generateUptrendCandles(10),
	}

	_, err := engine.Run(context.Background(), candles, 10_000_000)
	if err == nil {
		t.Fatal("expected error for insufficient candle data")
	}
}

func TestBacktestEmptyCandles(t *testing.T) {
	rg := trader.NewRiskGate(DefaultRiskConfig())
	agent := &mockBuyAgent{}
	engine := NewEngine(agent, rg)

	_, err := engine.Run(context.Background(), map[string][]indicator.CandleData{}, 10_000_000)
	if err == nil {
		t.Fatal("expected error for empty candle data")
	}
}

func TestBacktestZeroBalance(t *testing.T) {
	rg := trader.NewRiskGate(DefaultRiskConfig())
	agent := &mockBuyAgent{}
	engine := NewEngine(agent, rg)

	candles := map[string][]indicator.CandleData{
		"KRW-BTC": generateUptrendCandles(100),
	}

	_, err := engine.Run(context.Background(), candles, 0)
	if err == nil {
		t.Fatal("expected error for zero initial balance")
	}
}

func TestCalculateSharpeRatioEmpty(t *testing.T) {
	result := calculateSharpeRatio([]float64{})
	if result != 0 {
		t.Errorf("expected 0 for empty equity curve, got %f", result)
	}
}

func TestCalculateSharpeRatioSingleValue(t *testing.T) {
	result := calculateSharpeRatio([]float64{100})
	if result != 0 {
		t.Errorf("expected 0 for single value equity curve, got %f", result)
	}
}

func TestCalculateWinRateEmpty(t *testing.T) {
	result := calculateWinRate(nil)
	if result != 0 {
		t.Errorf("expected 0 for empty trades, got %f", result)
	}
}

func TestBacktestSharpeRatioFinite(t *testing.T) {
	rg := trader.NewRiskGate(DefaultRiskConfig())
	agent := &mockBuyAgent{}
	engine := NewEngine(agent, rg)

	candles := map[string][]indicator.CandleData{
		"KRW-BTC": generateUptrendCandles(200),
	}

	result, err := engine.Run(context.Background(), candles, 10_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.IsInf(result.SharpeRatio, 0) || math.IsNaN(result.SharpeRatio) {
		t.Errorf("SharpeRatio should be finite, got %f", result.SharpeRatio)
	}
}
