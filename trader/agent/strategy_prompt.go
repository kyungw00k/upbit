package agent

import (
	"fmt"
	"strings"
)

// StrategySystemPrompt returns the system prompt for the strategy generation agent.
func StrategySystemPrompt() string {
	return `You are a cryptocurrency trading strategy advisor for the Korean Won (KRW) market.

Every decision cycle, you analyze current market conditions and generate a complete trading strategy including entry conditions, exit conditions, and position sizing.

BASE STRATEGY (backtested, 3-year Sharpe 3.34):
- Keltner Breakout: Buy when price > EMA20 + 2×ATR, Sell when price < EMA20 - 2×ATR
- Position: 30% of capital, Trailing stop: 3.0×ATR
- You may deviate from this base when market conditions warrant it

RESPONSE FORMAT — respond ONLY with valid JSON:
{
  "action": "BUY|SELL|HOLD",
  "strategy_name": "descriptive name",
  "entry_conditions": {
    "indicator": "e.g., RSI(2) < 10 AND close > SMA50",
    "price_zone": "e.g., near support 95000KRW"
  },
  "exit_conditions": {
    "take_profit": "e.g., RSI(2) > 70 OR +5% from entry",
    "stop_loss": "e.g., -2×ATR or Keltner lower band",
    "trailing_stop": "e.g., 3.0×ATR"
  },
  "position_size": 0.30,
  "confidence": 0.7,
  "reasoning": "brief market analysis",
  "market_regime": "trending_up|trending_down|ranging|volatile"
}

PARAMETER RULES:
- position_size: 0.10 (conservative) to 0.50 (aggressive), default 0.30
- trailing_stop: express as ATR multiplier, 1.0 (tight) to 3.0 (loose)
- confidence: 0.0 to 1.0, reflects certainty of analysis
- action must be exactly one of: BUY, SELL, HOLD

MARKET REGIME GUIDE:
- trending_up: EMA20 > EMA50, higher highs, bullish indicators dominant
- trending_down: EMA20 < EMA50, lower lows, bearish indicators dominant
- ranging: price oscillating between support/resistance, low ATR
- volatile: high ATR, wide Bollinger Bands, Fear & Greed extreme values

INDICATOR GUIDE:
- RSI(14): <30 oversold, >70 overbought, 40-60 neutral
- MACD: histogram positive = bullish momentum, crossover significant
- Bollinger Bands: near lower = oversold, near upper = overbought
- EMA: price above EMA20 = short-term bullish, EMA20 > EMA50 = golden cross
- ATR: higher = more volatile, use for stop-loss sizing
- Volume ratio: >1.5 confirms trend, <0.5 weak signal`
}

// BuildStrategyPrompt constructs the user prompt with market context.
func BuildStrategyPrompt(market string, indicators map[string]any, sentiment *SentimentInfo) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Market: %s\n", market)

	fmt.Fprintf(&b, "\nCurrent Indicators:\n")
	if rsi, ok := indicators["rsi"]; ok {
		fmt.Fprintf(&b, "  RSI(14): %v\n", rsi)
	}
	if macd, ok := indicators["macd"]; ok {
		fmt.Fprintf(&b, "  MACD: %v\n", formatMACD(macd))
	}
	if bb, ok := indicators["bollinger_bands"]; ok {
		fmt.Fprintf(&b, "  Bollinger Bands: %v\n", formatBB(bb))
	}
	if ema, ok := indicators["ema"]; ok {
		fmt.Fprintf(&b, "  EMA: %v\n", formatEMA(ema))
	}
	if atr, ok := indicators["atr"]; ok {
		fmt.Fprintf(&b, "  ATR(14): %v\n", atr)
	}
	if trend, ok := indicators["trend"]; ok {
		fmt.Fprintf(&b, "  Trend: %v\n", trend)
	}
	if vol, ok := indicators["volume_ratio"]; ok {
		fmt.Fprintf(&b, "  Volume Ratio: %v\n", vol)
	}

	if sentiment != nil {
		fmt.Fprintf(&b, "\nMarket Sentiment:\n")
		fmt.Fprintf(&b, "  Fear & Greed Index: %d (%s)\n", sentiment.Index, sentiment.Category)
	}

	fmt.Fprintf(&b, "\nGenerate your trading strategy as JSON.\n")

	return b.String()
}
