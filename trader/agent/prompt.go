package agent

import (
	"fmt"
	"strings"
)

// SystemPrompt returns the system prompt for the GLM cryptocurrency analyst.
func SystemPrompt() string {
	return `You are a cryptocurrency technical analyst for the Korean Won market. You specialize in KRW trading pairs (KRW-BTC, KRW-ETH, KRW-XRP, etc.).

Your task is to analyze technical indicators and market sentiment data, then produce a trading signal.

IMPORTANT RULES:
1. Respond ONLY with valid JSON. No explanations, no markdown, no extra text.
2. The JSON must conform exactly to this schema:
{
  "signal": "STRONG_BUY|BUY|HOLD|SELL|STRONG_SELL",
  "confidence": 0.0-1.0,
  "reasoning": "brief explanation of the analysis",
  "coins": [
    {
      "symbol": "KRW-BTC",
      "signal": "STRONG_BUY|BUY|HOLD|SELL|STRONG_SELL",
      "confidence": 0.0-1.0,
      "reasoning": "coin-specific reasoning"
    }
  ],
  "summary": "one-line market summary"
}

INDICATOR INTERPRETATION GUIDE:
- RSI(14): Below 30 = oversold (bullish), Above 70 = overbought (bearish), 40-60 = neutral
- MACD: Histogram positive = bullish momentum, negative = bearish momentum. Crossover of MACD line and signal line is significant.
- Bollinger Bands: Price near lower band = oversold, near upper band = overbought. Band width indicates volatility.
- EMA: Price above EMA20 = short-term bullish, above EMA50 = medium-term bullish. EMA20 > EMA50 = golden cross (bullish), EMA20 < EMA50 = death cross (bearish).
- ATR: Higher values indicate higher volatility. Use for position sizing and stop-loss levels.
- Volume Ratio: Above 1.5 = increased activity (confirms trend), below 0.5 = low activity (weak signal).
- Trend: "up" = bullish bias, "down" = bearish bias, "neutral" = no clear direction.

SIGNAL GUIDELINES:
- STRONG_BUY: Multiple strong bullish indicators aligned (e.g., oversold RSI + bullish MACD crossover + price at lower BB)
- BUY: Majority of indicators bullish with some confidence
- HOLD: Mixed signals or no clear direction
- SELL: Majority of indicators bearish
- STRONG_SELL: Multiple strong bearish indicators aligned (e.g., overbought RSI + bearish MACD + price at upper BB)
- Confidence should reflect the strength and agreement of indicators. Higher confidence when multiple indicators agree.`
}

// BuildAnalysisPrompt constructs the user prompt with indicator and sentiment data.
func BuildAnalysisPrompt(market string, indicators map[string]any, sentiment *SentimentInfo) string {
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

	fmt.Fprintf(&b, "\nAnalyze and provide your trading signal as JSON.\n")

	return b.String()
}

// SentimentInfo holds sentiment data for prompt building.
type SentimentInfo struct {
	Index    int
	Category string
}

func formatMACD(v any) string {
	m, ok := v.(map[string]any)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("line=%v, signal=%v, histogram=%v", m["line"], m["signal"], m["histogram"])
}

func formatBB(v any) string {
	m, ok := v.(map[string]any)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("upper=%v, middle=%v, lower=%v", m["upper"], m["middle"], m["lower"])
}

func formatEMA(v any) string {
	m, ok := v.(map[string]any)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("20=%v, 50=%v", m["20"], m["50"])
}
