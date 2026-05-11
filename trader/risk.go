package trader

import (
	"fmt"

	"github.com/kyungw00k/upbit/trader/config"
	"github.com/kyungw00k/upbit/trader/types"
)

// Default risk configuration values used when config fields are zero.
const (
	DefaultMaxPositionPct  = 0.10
	DefaultMaxDailyLossPct = 0.03
	DefaultStopLossPct     = 0.10
	DefaultMaxPositions    = 5

	DefaultFeeRate = 0.0005 // 0.05% trading fee
)

// Position holds details about an open position for a single market.
type Position struct {
	Market     string
	EntryPrice float64
	Volume     float64
}

// Portfolio holds the current state of the trading portfolio.
type Portfolio struct {
	TotalValue float64
	Positions  map[string]Position
	DailyPnL   float64
}

// RiskGate evaluates trading signals against risk management rules.
// It is pure Go with no LLM dependency.
type RiskGate struct {
	config config.RiskConfig
}

// NewRiskGate creates a RiskGate with the given configuration.
// Applies default values for any zero config fields.
func NewRiskGate(cfg config.RiskConfig) *RiskGate {
	if cfg.MaxPositionPct <= 0 {
		cfg.MaxPositionPct = DefaultMaxPositionPct
	}
	if cfg.MaxDailyLossPct <= 0 {
		cfg.MaxDailyLossPct = DefaultMaxDailyLossPct
	}
	if cfg.StopLossPct <= 0 {
		cfg.StopLossPct = DefaultStopLossPct
	}
	if cfg.MaxPositions <= 0 {
		cfg.MaxPositions = DefaultMaxPositions
	}
	return &RiskGate{config: cfg}
}

// Evaluate checks whether a proposed trade passes all risk rules.
// Returns (approved, adjustedSize, reason).
// If approved is false, adjustedSize will be 0 and reason explains why.
func (r *RiskGate) Evaluate(signal types.CoinAnalysis, portfolio Portfolio) (bool, float64, string) {
	// Validate required signal fields.
	if signal.Symbol == "" {
		return false, 0, "rejected: symbol is empty"
	}
	if signal.EntryPrice <= 0 {
		return false, 0, "rejected: entry_price must be positive"
	}
	if signal.Position <= 0 {
		return false, 0, "rejected: position size must be positive"
	}

	// Rule 1: Stop loss must be set within StopLossPct of entry price.
	if signal.StopLoss <= 0 {
		return false, 0, "rejected: stop_loss is missing"
	}
	if signal.StopLoss >= signal.EntryPrice {
		return false, 0, "rejected: stop_loss must be below entry_price for long positions"
	}
	stopLossPct := (signal.EntryPrice - signal.StopLoss) / signal.EntryPrice
	if stopLossPct > r.config.StopLossPct {
		return false, 0, fmt.Sprintf("rejected: stop_loss_pct %.4f exceeds max %.4f", stopLossPct, r.config.StopLossPct)
	}

	// Rule 2: Single coin position must not exceed MaxPositionPct of total portfolio.
	positionValue := signal.Position * signal.EntryPrice
	if portfolio.TotalValue > 0 {
		positionPct := positionValue / portfolio.TotalValue
		// Include existing position for the same market.
		if existing, ok := portfolio.Positions[signal.Symbol]; ok {
			existingValue := existing.Volume * existing.EntryPrice
			positionPct = (positionValue + existingValue) / portfolio.TotalValue
		}
		if positionPct > r.config.MaxPositionPct {
			return false, 0, fmt.Sprintf("rejected: position_pct %.4f exceeds max %.4f", positionPct, r.config.MaxPositionPct)
		}
	}

	// Rule 3: Daily loss must not exceed MaxDailyLossPct of total portfolio.
	if portfolio.TotalValue > 0 && portfolio.DailyPnL < 0 {
		dailyLossPct := -portfolio.DailyPnL / portfolio.TotalValue
		if dailyLossPct > r.config.MaxDailyLossPct {
			return false, 0, fmt.Sprintf("rejected: daily_loss_pct %.4f exceeds max %.4f", dailyLossPct, r.config.MaxDailyLossPct)
		}
	}

	// Rule 4: Total number of concurrent positions must not exceed MaxPositions.
	// Count existing positions excluding the current market (which would be an update).
	existingCount := len(portfolio.Positions)
	if _, exists := portfolio.Positions[signal.Symbol]; !exists {
		existingCount++
	}
	if existingCount > r.config.MaxPositions {
		return false, 0, fmt.Sprintf("rejected: position count %d exceeds max %d", existingCount, r.config.MaxPositions)
	}

	// All checks passed.
	return true, signal.Position, "approved"
}
