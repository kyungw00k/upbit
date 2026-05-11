package agents

import (
	"context"
	"fmt"
	"time"

	"github.com/kyungw00k/upbit/trader/indicator"
	"github.com/kyungw00k/upbit/trader/types"
)

// RSI2Agent implements the Connors RSI-2 mean reversion strategy.
// Buys when RSI(2) is oversold during uptrends, sells when RSI(2) recovers.
// Long-only — suitable for Upbit spot trading.
//
// Multi-timeframe logic:
//   - 4h: trend filter (close > SMA50 → uptrend)
//   - 1h: primary signal (RSI(2) < oversoldThreshold → BUY)
//   - exit: RSI(2) > overboughtThreshold or trailing stop
type RSI2Agent struct {
	OversoldThreshold  float64 // RSI(2) buy threshold (default: 10)
	OverboughtThreshold float64 // RSI(2) sell threshold (default: 70)
	TrendSMAPeriod     int     // SMA period for trend filter (default: 50)
}

// NewRSI2Agent creates an RSI-2 agent with default parameters.
func NewRSI2Agent() *RSI2Agent {
	return &RSI2Agent{
		OversoldThreshold:   10,
		OverboughtThreshold: 70,
		TrendSMAPeriod:      50,
	}
}

// ID returns the agent identifier.
func (a *RSI2Agent) ID() string { return "rsi2-agent" }

// Run evaluates the RSI-2 strategy on the provided indicator data.
func (a *RSI2Agent) Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	market := "KRW-BTC"
	if len(input.Markets) > 0 {
		market = input.Markets[0]
	}

	// Extract indicator values
	ind := input.Indicators
	if ind == nil {
		return holdReport(market, "no indicators"), nil
	}

	rsi2, hasRSI2 := extractFloat(ind, "rsi2")
	closePrice, hasClose := extractFloat(ind, "close")
	sma50, hasSMA := extractFloat(ind, "sma50")

	if !hasRSI2 || !hasClose {
		return holdReport(market, "missing rsi2 or close"), nil
	}

	// Trend filter: only buy in uptrend (close > SMA50)
	uptrend := true
	if hasSMA && closePrice < sma50 {
		uptrend = false
	}

	var signal types.Signal
	var reasoning string
	confidence := 0.5

	switch {
	case uptrend && rsi2 < a.OversoldThreshold:
		signal = types.SignalBuy
		if rsi2 < 5 {
			signal = types.SignalStrongBuy
			confidence = 0.9
		} else {
			confidence = 0.7
		}
		reasoning = fmt.Sprintf("RSI2=%.1f oversold in uptrend (close=%.0f > SMA50=%.0f)", rsi2, closePrice, sma50)

	case !uptrend && rsi2 < a.OversoldThreshold:
		signal = types.SignalHold
		reasoning = fmt.Sprintf("RSI2=%.1f oversold but DOWNTREND (close=%.0f < SMA50=%.0f) — skip", rsi2, closePrice, sma50)
		confidence = 0.3

	case rsi2 > a.OverboughtThreshold:
		signal = types.SignalSell
		if rsi2 > 90 {
			signal = types.SignalStrongSell
			confidence = 0.9
		} else {
			confidence = 0.7
		}
		reasoning = fmt.Sprintf("RSI2=%.1f overbought — exit signal", rsi2)

	default:
		signal = types.SignalHold
		reasoning = fmt.Sprintf("RSI2=%.1f neutral (oversold<%.0f, overbought>%.0f)", rsi2, a.OversoldThreshold, a.OverboughtThreshold)
		confidence = 0.3
	}

	return types.AnalysisReport{
		AgentID:   a.ID(),
		Timestamp: time.Now().UTC(),
		Signal:    signal,
		Coins: []types.CoinAnalysis{
			{
				Symbol:     market,
				Signal:     signal,
				Confidence: confidence,
				Reasoning:  reasoning,
			},
		},
		Summary: reasoning,
	}, nil
}

// EvaluateCandles runs the full RSI-2 evaluation on raw candle data.
// Returns BUY/SELL/HOLD with reasoning. Used for backtesting.
func (a *RSI2Agent) EvaluateCandles(candles []indicator.CandleData) (signal types.Signal, rsi2 float64, reasoning string) {
	if len(candles) < a.TrendSMAPeriod+2 {
		return types.SignalHold, 0, "insufficient data"
	}

	closes := make([]float64, len(candles))
	for i, c := range candles {
		closes[i] = c.Close
	}

	// Calculate RSI(2)
	rsi2Values := indicator.RSI(closes, 2)
	currentRSI2 := rsi2Values[len(rsi2Values)-1]

	// Calculate SMA50 for trend
	smaValues := indicator.SMA(closes, a.TrendSMAPeriod)
	currentSMA50 := smaValues[len(smaValues)-1]
	currentClose := closes[len(closes)-1]

	uptrend := currentClose > currentSMA50

	switch {
	case uptrend && currentRSI2 < a.OversoldThreshold:
		if currentRSI2 < 5 {
			return types.SignalStrongBuy, currentRSI2, fmt.Sprintf("RSI2=%.1f deep oversold in uptrend", currentRSI2)
		}
		return types.SignalBuy, currentRSI2, fmt.Sprintf("RSI2=%.1f oversold in uptrend", currentRSI2)

	case !uptrend && currentRSI2 < a.OversoldThreshold:
		return types.SignalHold, currentRSI2, fmt.Sprintf("RSI2=%.1f oversold but downtrend — skip", currentRSI2)

	case currentRSI2 > a.OverboughtThreshold:
		return types.SignalSell, currentRSI2, fmt.Sprintf("RSI2=%.1f overbought", currentRSI2)

	default:
		return types.SignalHold, currentRSI2, fmt.Sprintf("RSI2=%.1f neutral", currentRSI2)
	}
}

func holdReport(market, reason string) types.AnalysisReport {
	return types.AnalysisReport{
		AgentID:   "rsi2-agent",
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalHold,
		Coins: []types.CoinAnalysis{
			{Symbol: market, Signal: types.SignalHold, Reasoning: reason},
		},
		Summary: reason,
	}
}
