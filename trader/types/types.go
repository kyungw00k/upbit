package types

import (
	"encoding/json"
	"time"
)

// Signal represents a trading signal strength.
type Signal string

const (
	SignalStrongBuy  Signal = "STRONG_BUY"
	SignalBuy        Signal = "BUY"
	SignalHold       Signal = "HOLD"
	SignalSell       Signal = "SELL"
	SignalStrongSell Signal = "STRONG_SELL"
)

// Valid returns true if the signal is one of the recognized values.
func (s Signal) Valid() bool {
	switch s {
	case SignalStrongBuy, SignalBuy, SignalHold, SignalSell, SignalStrongSell:
		return true
	}
	return false
}

// CoinAnalysis holds analysis results for a single coin.
type CoinAnalysis struct {
	Symbol     string  `json:"symbol"`
	Signal     Signal  `json:"signal"`
	Confidence float64 `json:"confidence"`
	EntryPrice float64 `json:"entry_price"`
	StopLoss   float64 `json:"stop_loss"`
	TakeProfit float64 `json:"take_profit"`
	Position   float64 `json:"position"`
	Reasoning  string  `json:"reasoning"`
}

// AnalysisReport is the output of an agent's analysis.
type AnalysisReport struct {
	AgentID   string         `json:"agent_id"`
	Timestamp time.Time      `json:"timestamp"`
	Signal    Signal         `json:"signal"`
	Coins     []CoinAnalysis `json:"coins"`
	Summary   string         `json:"summary"`
}

// PipelineResult records the outcome of a single pipeline run.
type PipelineResult struct {
	RanAt    time.Time     `json:"ran_at"`
	Duration time.Duration `json:"-"`
}

// MarshalJSON implements json.Marshaler for PipelineResult.
func (r PipelineResult) MarshalJSON() ([]byte, error) {
	type Alias struct {
		RanAt      time.Time `json:"ran_at"`
		DurationMs int64     `json:"duration_ms"`
	}
	return json.Marshal(Alias{
		RanAt:      r.RanAt,
		DurationMs: r.Duration.Milliseconds(),
	})
}

// AgentInput is the data passed to an agent for analysis.
type AgentInput struct {
	Markets    []string       `json:"markets"`
	Indicators map[string]any `json:"indicators,omitempty"`
	Sentiment  *SentimentData `json:"sentiment,omitempty"`
}

// SentimentData holds market sentiment information.
type SentimentData struct {
	FearGreedIndex    int    `json:"fear_greed_index"`
	FearGreedCategory string `json:"fear_greed_category"`
	Source            string `json:"source"`
	Timestamp         string `json:"timestamp"`
}

// StrategyResponse is returned by the LLM strategy agent every decision cycle.
type StrategyResponse struct {
	Action          string            `json:"action"`
	StrategyName    string            `json:"strategy_name"`
	EntryConditions map[string]string `json:"entry_conditions"`
	ExitConditions  ExitConditions    `json:"exit_conditions"`
	PositionSize    float64           `json:"position_size"`
	Confidence      float64           `json:"confidence"`
	Reasoning       string            `json:"reasoning"`
	MarketRegime    string            `json:"market_regime"`
}

// ExitConditions describes when to exit a position.
type ExitConditions struct {
	TakeProfit   string `json:"take_profit"`
	StopLoss     string `json:"stop_loss"`
	TrailingStop string `json:"trailing_stop"`
}

// ExecutionParams holds the finalized parameters for trade execution.
type ExecutionParams struct {
	PosSize    float64 // fraction of capital (0.10 - 0.50)
	TrailMult  float64 // ATR multiplier for trailing stop (1.0 - 3.0)
	StopLoss   float64 // absolute price
	TakeProfit float64 // absolute price (0 = none)
	Confidence float64 // 0.0 - 1.0
}

// DefaultExecutionParams returns the proven backtested defaults (30%, 3.0xATR).
func DefaultExecutionParams() ExecutionParams {
	return ExecutionParams{
		PosSize:    0.30,
		TrailMult:  3.0,
		StopLoss:   0,
		TakeProfit: 0,
		Confidence: 0.5,
	}
}

// Clamp returns a copy with fields clamped to valid ranges.
func (e ExecutionParams) Clamp() ExecutionParams {
	e.PosSize = clampF(e.PosSize, 0.10, 0.50)
	e.TrailMult = clampF(e.TrailMult, 1.0, 3.0)
	e.Confidence = clampF(e.Confidence, 0.0, 1.0)
	if e.StopLoss < 0 {
		e.StopLoss = 0
	}
	if e.TakeProfit < 0 {
		e.TakeProfit = 0
	}
	return e
}

func clampF(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
