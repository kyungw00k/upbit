package trader

import (
	"testing"

	"github.com/kyungw00k/upbit/trader/config"
	"github.com/kyungw00k/upbit/trader/types"
)

func defaultRiskConfig() config.RiskConfig {
	return config.RiskConfig{
		MaxPositionPct:  0.10,
		MaxDailyLossPct: 0.03,
		StopLossPct:     0.05,
		MaxPositions:    5,
	}
}

func defaultPortfolio() Portfolio {
	return Portfolio{
		TotalValue: 10_000_000,
		Positions:  map[string]Position{},
		DailyPnL:   0,
	}
}

func validBuySignal() types.CoinAnalysis {
	return types.CoinAnalysis{
		Symbol:     "KRW-BTC",
		Signal:     types.SignalBuy,
		EntryPrice: 75_000_000,
		StopLoss:   73_000_000,
		TakeProfit: 80_000_000,
		Position:   0.001,
		Confidence: 0.8,
		Reasoning:  "test signal",
	}
}

func TestRiskGateApproved(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	signal := validBuySignal()
	// Position value: 0.001 * 75,000,000 = 75,000 KRW = 0.75% of 10M portfolio.
	approved, adjustedSize, reason := rg.Evaluate(signal, defaultPortfolio())
	if !approved {
		t.Errorf("expected approval, got rejection: %s", reason)
	}
	if adjustedSize <= 0 {
		t.Errorf("expected positive adjustedSize, got %f", adjustedSize)
	}
	if reason != "approved" {
		t.Errorf("expected reason 'approved', got %q", reason)
	}
}

func TestRiskGateEmptySymbol(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	signal := validBuySignal()
	signal.Symbol = ""
	approved, _, reason := rg.Evaluate(signal, defaultPortfolio())
	if approved {
		t.Error("should reject empty symbol")
	}
	if reason == "" {
		t.Error("expected non-empty rejection reason")
	}
}

func TestRiskGateZeroEntryPrice(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	signal := validBuySignal()
	signal.EntryPrice = 0
	approved, _, reason := rg.Evaluate(signal, defaultPortfolio())
	if approved {
		t.Error("should reject zero entry price")
	}
	if reason == "" {
		t.Error("expected non-empty rejection reason")
	}
}

func TestRiskGateZeroPosition(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	signal := validBuySignal()
	signal.Position = 0
	approved, _, reason := rg.Evaluate(signal, defaultPortfolio())
	if approved {
		t.Error("should reject zero position size")
	}
	if reason == "" {
		t.Error("expected non-empty rejection reason")
	}
}

func TestRiskGatePositionConcentration(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	portfolio := Portfolio{
		TotalValue: 10_000_000,
		Positions: map[string]Position{
			"KRW-BTC": {Market: "KRW-BTC", EntryPrice: 75_000_000, Volume: 0.005},
		},
		DailyPnL: 0,
	}
	// Existing BTC position: 0.005 * 75,000,000 = 375,000 KRW = 3.75%
	// New position: 0.001 * 75,000,000 = 75,000 KRW = 0.75%
	// Combined: 4.5% — should be approved (under 10%).
	smallSignal := validBuySignal()
	smallSignal.Position = 0.001
	approved, _, _ := rg.Evaluate(smallSignal, portfolio)
	if !approved {
		t.Error("should be approved, combined position within 10%")
	}

	// Large new position that would exceed 10% combined.
	// 0.01 * 75,000,000 = 750,000 KRW = 7.5% + existing 3.75% = 11.25%.
	bigSignal := validBuySignal()
	bigSignal.Position = 0.01
	approved, adjustedSize, reason := rg.Evaluate(bigSignal, portfolio)
	if approved {
		t.Error("should be rejected for position concentration")
	}
	if adjustedSize != 0 {
		t.Errorf("expected adjustedSize 0 for rejected trade, got %f", adjustedSize)
	}
	if reason == "" {
		t.Error("expected non-empty rejection reason")
	}
}

func TestRiskGateDailyLossExceeded(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	// Daily loss of 350,000 on 10M portfolio = 3.5%, exceeds 3% limit.
	portfolio := defaultPortfolio()
	portfolio.DailyPnL = -350_000

	signal := validBuySignal()
	signal.Symbol = "KRW-ETH"
	signal.EntryPrice = 3_000_000
	signal.StopLoss = 2_900_000
	signal.Position = 0.1

	approved, _, reason := rg.Evaluate(signal, portfolio)
	if approved {
		t.Error("should be rejected for daily loss limit")
	}
	if reason == "" {
		t.Error("expected non-empty rejection reason")
	}
}

func TestRiskGateDailyLossExactlyAtLimit(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	// Daily loss of exactly 300,000 on 10M = 3.0%, should be approved (not exceeded).
	portfolio := defaultPortfolio()
	portfolio.DailyPnL = -300_000

	approved, _, _ := rg.Evaluate(validBuySignal(), portfolio)
	if !approved {
		t.Error("should be approved when daily loss exactly at limit")
	}
}

func TestRiskGateMaxPositionsExceeded(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	// Already have 5 positions (max), adding a 6th should fail.
	positions := map[string]Position{
		"KRW-BTC":  {Market: "KRW-BTC", EntryPrice: 75_000_000, Volume: 0.001},
		"KRW-ETH":  {Market: "KRW-ETH", EntryPrice: 3_000_000, Volume: 0.1},
		"KRW-SOL":  {Market: "KRW-SOL", EntryPrice: 150_000, Volume: 5},
		"KRW-XRP":  {Market: "KRW-XRP", EntryPrice: 500, Volume: 1000},
		"KRW-DOGE": {Market: "KRW-DOGE", EntryPrice: 300, Volume: 2000},
	}
	portfolio := defaultPortfolio()
	portfolio.Positions = positions

	signal := validBuySignal()
	signal.Symbol = "KRW-ADA"
	signal.EntryPrice = 800
	signal.StopLoss = 780
	signal.Position = 100

	approved, _, reason := rg.Evaluate(signal, portfolio)
	if approved {
		t.Error("should be rejected for max positions exceeded")
	}
	if reason == "" {
		t.Error("expected non-empty rejection reason")
	}
}

func TestRiskGateStopLossMissing(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	signal := validBuySignal()
	signal.StopLoss = 0

	approved, adjustedSize, reason := rg.Evaluate(signal, defaultPortfolio())
	if approved {
		t.Error("should be rejected for missing stop loss")
	}
	if adjustedSize != 0 {
		t.Errorf("expected adjustedSize 0, got %f", adjustedSize)
	}
	if reason == "" {
		t.Error("expected non-empty rejection reason")
	}
}

func TestRiskGateStopLossTooFar(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	// Stop loss at 10% below entry, exceeds max 5%.
	signal := validBuySignal()
	signal.StopLoss = 67_500_000

	approved, _, reason := rg.Evaluate(signal, defaultPortfolio())
	if approved {
		t.Error("should be rejected for stop loss too far")
	}
	if reason == "" {
		t.Error("expected non-empty rejection reason")
	}
}

func TestRiskGateStopLossAboveEntry(t *testing.T) {
	rg := NewRiskGate(defaultRiskConfig())
	signal := validBuySignal()
	signal.StopLoss = 80_000_000 // above entry

	approved, _, reason := rg.Evaluate(signal, defaultPortfolio())
	if approved {
		t.Error("should be rejected for stop loss above entry price")
	}
	if reason == "" {
		t.Error("expected non-empty rejection reason")
	}
}

func TestRiskGateExistingPositionUpdateDoesNotCount(t *testing.T) {
	rg := NewRiskGate(config.RiskConfig{
		MaxPositionPct:  0.10,
		MaxDailyLossPct: 0.03,
		StopLossPct:     0.05,
		MaxPositions:    2,
	})
	positions := map[string]Position{
		"KRW-BTC": {Market: "KRW-BTC", EntryPrice: 75_000_000, Volume: 0.001},
		"KRW-ETH": {Market: "KRW-ETH", EntryPrice: 3_000_000, Volume: 0.1},
	}
	portfolio := defaultPortfolio()
	portfolio.Positions = positions

	// Updating existing BTC position should not trigger max positions limit.
	approved, _, _ := rg.Evaluate(validBuySignal(), portfolio)
	if !approved {
		t.Error("updating existing position should not trigger max positions limit")
	}
}

func TestRiskGateDefaultConfig(t *testing.T) {
	// Empty config should apply defaults.
	rg := NewRiskGate(config.RiskConfig{})
	approved, _, reason := rg.Evaluate(validBuySignal(), defaultPortfolio())
	if !approved {
		t.Errorf("should be approved with default config: %s", reason)
	}
}
