package agents

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/kyungw00k/upbit/trader/types"
)

// TechnicalAgent is a rule-based agent that generates trading signals from
// technical indicator thresholds without any LLM dependency.
type TechnicalAgent struct{}

// NewTechnicalAgent creates a new TechnicalAgent.
func NewTechnicalAgent() *TechnicalAgent {
	return &TechnicalAgent{}
}

// ID returns the agent identifier.
func (t *TechnicalAgent) ID() string {
	return "technical-agent"
}

// Run evaluates the indicators in AgentInput and returns an AnalysisReport.
func (t *TechnicalAgent) Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	score := 0.0

	// RSI scoring — weight 30%.
	if v, ok := extractFloat(input.Indicators, "rsi"); ok {
		score += scoreRSI(v)
	}

	// MACD histogram scoring — weight 25%.
	if v, ok := extractMACDHist(input.Indicators); ok {
		score += scoreMACD(v)
	}

	// Bollinger Bands scoring — weight 20%.
	if s, ok := scoreBollinger(input.Indicators); ok {
		score += s
	}

	// EMA scoring — weight 15%.
	if s, ok := scoreEMA(input.Indicators); ok {
		score += s
	}

	// Volume scoring — weight 10%.
	if v, ok := extractFloat(input.Indicators, "volume_ratio"); ok {
		score += scoreVolume(v)
	}

	signal := scoreToSignal(score)
	confidence := math.Abs(score) / 100.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	market := "KRW-BTC"
	if len(input.Markets) > 0 {
		market = input.Markets[0]
	}

	reasoning := fmt.Sprintf("Rule-based score: %.1f → %s", score, signal)

	coins := []types.CoinAnalysis{
		{
			Symbol:     market,
			Signal:     signal,
			Confidence: confidence,
			Reasoning:  reasoning,
		},
	}

	return types.AnalysisReport{
		AgentID:   t.ID(),
		Timestamp: time.Now().UTC(),
		Signal:    signal,
		Coins:     coins,
		Summary:   reasoning,
	}, nil
}

// scoreRSI returns a score from -30 to +30 based on RSI.
// RSI < 30 → +30 (oversold), RSI 40 → +15, RSI 50 → 0, RSI 60 → -15, RSI > 70 → -30.
func scoreRSI(rsi float64) float64 {
	if rsi < 30 {
		return 30
	}
	if rsi > 70 {
		return -30
	}
	return 30 - (rsi-30)*(60.0/40.0)
}

// scoreMACD returns a continuous score based on histogram magnitude.
func scoreMACD(hist float64) float64 {
	if hist > 0 {
		score := math.Min(hist*0.001, 25)
		if score < 5 {
			score = 5 // minimum bullish bias
		}
		return score
	}
	if hist < 0 {
		score := math.Max(hist*0.001, -25)
		if score > -5 {
			score = -5
		}
		return score
	}
	return 0
}

// scoreBollinger returns -20 to +20 based on close position relative to bands.
func scoreBollinger(indicators map[string]any) (float64, bool) {
	bb, ok := indicators["bollinger_bands"]
	if !ok {
		return 0, false
	}
	m, ok := bb.(map[string]any)
	if !ok {
		return 0, false
	}
	upper, ok1 := toFloat(m["upper"])
	lower, ok2 := toFloat(m["lower"])
	if !ok1 || !ok2 {
		return 0, false
	}
	closePrice, ok := extractFloat(indicators, "close")
	if !ok {
		return 0, false
	}

	bandWidth := upper - lower
	if bandWidth == 0 {
		return 0, true
	}
	position := (closePrice - lower) / bandWidth // 0 at lower, 1 at upper.

	// Near lower band → +20 (oversold), near upper → -20 (overbought).
	return 20 - position*40, true
}

// scoreEMA returns a trend-filtered score. Penalizes buying in downtrends heavily.
func scoreEMA(indicators map[string]any) (float64, bool) {
	ema, ok := indicators["ema"]
	if !ok {
		return 0, false
	}
	m, ok := ema.(map[string]any)
	if !ok {
		return 0, false
	}
	closePrice, ok := extractFloat(indicators, "close")
	if !ok {
		return 0, false
	}

	var score float64
	var e20, e50 float64
	var hasE20, hasE50 bool

	if v, ok := toFloat(m["20"]); ok {
		e20 = v
		hasE20 = true
	}
	if v, ok := toFloat(m["50"]); ok {
		e50 = v
		hasE50 = true
	}

	// Trend filter: if close < both EMAs, strong penalty (don't buy in downtrend)
	if hasE20 && hasE50 && closePrice < e20 && closePrice < e50 {
		return -30, true
	}

	if hasE20 {
		if closePrice > e20 {
			score += 15
		} else {
			score -= 15
		}
	}
	if hasE50 {
		if closePrice > e50 {
			score += 15
		} else {
			score -= 15
		}
	}
	return score, true
}

// scoreVolume returns +10 if volume ratio > 1.5 (confirms trend).
func scoreVolume(ratio float64) float64 {
	if ratio > 1.5 {
		return 10
	}
	if ratio < 0.5 {
		return -10
	}
	return 0
}

// scoreToSignal maps a numeric score to a Signal.
func scoreToSignal(score float64) types.Signal {
	switch {
	case score > 50:
		return types.SignalStrongBuy
	case score > 10:
		return types.SignalBuy
	case score < -50:
		return types.SignalStrongSell
	case score < -10:
		return types.SignalSell
	default:
		return types.SignalHold
	}
}

// extractFloat extracts a float64 from the indicators map.
func extractFloat(m map[string]any, key string) (float64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	return toFloat(v)
}

// extractMACDHist extracts the MACD histogram value.
func extractMACDHist(m map[string]any) (float64, bool) {
	v, ok := m["macd"]
	if !ok {
		return 0, false
	}
	mm, ok := v.(map[string]any)
	if !ok {
		return 0, false
	}
	return toFloat(mm["histogram"])
}

// toFloat converts a value to float64, handling json.Number, float64, int, etc.
func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	default:
		return 0, false
	}
}
