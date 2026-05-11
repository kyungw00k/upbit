package agents

import (
	"context"
	"fmt"
	"time"

	"github.com/kyungw00k/upbit/trader/types"
)

// KeltnerAgent is a rule-based agent using the Keltner Breakout strategy.
// Proven parameters from 3-year backtest: EMA20 + 2×ATR channel, 3.0×ATR trailing stop.
type KeltnerAgent struct {
	KeltnerMult float64 // Channel multiplier (default 2.0)
	EMAPeriod   int     // EMA period (default 20)
	ATRPeriod   int     // ATR period (default 14)
	PosSize     float64 // Position size fraction (default 0.30)
	TrailMult   float64 // Trailing stop ATR multiplier (default 3.0)
}

// NewKeltnerAgent creates a KeltnerAgent with proven default parameters.
func NewKeltnerAgent() *KeltnerAgent {
	return &KeltnerAgent{
		KeltnerMult: 2.0,
		EMAPeriod:   20,
		ATRPeriod:   14,
		PosSize:     0.30,
		TrailMult:   3.0,
	}
}

// ID returns the agent identifier.
func (k *KeltnerAgent) ID() string {
	return "keltner-agent"
}

// Run evaluates the Keltner channel breakout from indicator data.
func (k *KeltnerAgent) Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	market := "KRW-BTC"
	if len(input.Markets) > 0 {
		market = input.Markets[0]
	}

	emaVal, hasEMA := extractFloat(input.Indicators, "ema")
	atrVal, hasATR := extractFloat(input.Indicators, "atr")
	closeVal, hasClose := extractFloat(input.Indicators, "close")

	if !hasEMA || !hasATR || !hasClose {
		return k.holdReport(market, "insufficient indicator data"), nil
	}

	// Keltner channel
	upper := emaVal + k.KeltnerMult*atrVal
	lower := emaVal - k.KeltnerMult*atrVal

	var signal types.Signal
	var reasoning string
	var stopLoss float64

	if closeVal > upper {
		signal = types.SignalBuy
		stopLoss = closeVal - k.TrailMult*atrVal
		reasoning = fmt.Sprintf("Keltner breakout: price %.0f > upper %.0f (EMA20=%.0f, ATR=%.0f, stop=%.0f)",
			closeVal, upper, emaVal, atrVal, stopLoss)
	} else if closeVal < lower {
		signal = types.SignalSell
		reasoning = fmt.Sprintf("Keltner breakdown: price %.0f < lower %.0f (EMA20=%.0f, ATR=%.0f)",
			closeVal, lower, emaVal, atrVal)
	} else {
		signal = types.SignalHold
		reasoning = fmt.Sprintf("Within Keltner channel: %.0f < price %.0f < %.0f",
			lower, closeVal, upper)
	}

	confidence := 0.5
	if signal != types.SignalHold {
		// Higher confidence on stronger breakouts
		breakoutDist := 0.0
		if closeVal > upper && upper > 0 {
			breakoutDist = (closeVal - upper) / upper * 100
		} else if closeVal < lower && lower > 0 {
			breakoutDist = (lower - closeVal) / lower * 100
		}
		confidence = 0.5 + min(breakoutDist*0.1, 0.5)
	}

	coins := []types.CoinAnalysis{
		{
			Symbol:     market,
			Signal:     signal,
			Confidence: confidence,
			EntryPrice: closeVal,
			StopLoss:   stopLoss,
			Position:   k.PosSize,
			Reasoning:  reasoning,
		},
	}

	return types.AnalysisReport{
		AgentID:   k.ID(),
		Timestamp: time.Now().UTC(),
		Signal:    signal,
		Coins:     coins,
		Summary:   reasoning,
	}, nil
}

func (k *KeltnerAgent) holdReport(market, reason string) types.AnalysisReport {
	return types.AnalysisReport{
		AgentID:   k.ID(),
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalHold,
		Coins: []types.CoinAnalysis{
			{Symbol: market, Signal: types.SignalHold, Reasoning: reason},
		},
		Summary: fmt.Sprintf("HOLD: %s", reason),
	}
}
