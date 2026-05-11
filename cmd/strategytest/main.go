package main

import (
	"fmt"
	"math"
	"os"
	"text/tabwriter"

	"github.com/kyungw00k/upbit/trader/indicator"
	"github.com/kyungw00k/upbit/trader/storage"
)

type trade struct {
	entryPrice float64
	entryTime  string
	exitPrice  float64
	exitTime   string
	pnl        float64
	reason     string
}

type backtestResult struct {
	name        string
	category    string
	totalReturn float64
	winRate     float64
	mdd         float64
	trades      int
	wins        int
	finalBal    float64
	sharpe      float64
	monthlyRet  map[string]float64
	lossMonths  int
}

func main() {
	dbPath := "data/candles.db"
	market := "KRW-BTC"
	initBalance := 330_000.0
	feeRate := 0.0005
	posSize := 0.30

	db, err := storage.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB open: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Printf("=== Category Strategy Comparison: %s | Initial: %.0f KRW ===\n\n", market, initBalance)

	type strategy struct {
		name     string
		category string
		fn       func([]indicator.CandleData, float64, float64, float64) []trade
	}

	strategies := []strategy{
		// Mean Reversion
		{"RSI-2 Connors", "MeanReversion", strategyRSI2},
		{"RSI+BB Combined", "MeanReversion", strategyRSIBB},
		{"BB %B Reversion", "MeanReversion", strategyBBPercentB},
		{"Stochastic(14,3)", "MeanReversion", strategyStochastic},
		{"Williams %R", "MeanReversion", strategyWilliamsR},
		{"CCI Reversion", "MeanReversion", strategyCCIReversion},
		{"RSI-2 + Volume", "MeanReversion", strategyRSI2Volume},

		// Trend Following
		{"EMA Cross(20,50)", "TrendFollow", strategyEMACross},
		{"MACD+EMA Filter", "TrendFollow", strategyMACDEMA},
		{"Donchian Breakout", "TrendFollow", strategyDonchian},
		{"SMA Ribbon", "TrendFollow", strategySMARibbon},
		{"ATR Breakout", "TrendFollow", strategyATRBreakout},

		// Momentum
		{"RSI(14) Momentum", "Momentum", strategyRSIMomentum},
		{"MACD Histogram", "Momentum", strategyMACDHistogram},
		{"Rate of Change", "Momentum", strategyROC},

		// Volatility
		{"BB Squeeze", "Volatility", strategyBBSqueeze},
		{"ATR Squeeze", "Volatility", strategyATRSqueeze},
		{"Keltner Breakout", "Volatility", strategyKeltner},
			// Hybrid — regime-adaptive strategies
		{"Regime Switch", "Hybrid", strategyRegimeSwitch},
		{"Signal Voting", "Hybrid", strategyVoting},
		{"MR+Vol Conditional", "Hybrid", strategyMRVolConditional},
		{"Best2 Adaptive", "Hybrid", strategyBest2Adaptive},

		// Williams -- Larry Williams volatility breakout
		{"Larry Williams", "Williams", strategyLarryWilliams},
		{"LW+MA10 Filter", "Williams", strategyLWMA10},
		{"RSI2->LW Pipeline", "Williams", strategyRSI2LWPipeline},
	}

	for _, period := range []struct {
		label string
		from  string
		to    string
	}{
		{"2024", "2024-01-01T00:00:00", "2024-12-31T23:59:59"},
		{"2025", "2025-01-01T00:00:00", "2025-12-31T23:59:59"},
		{"2024+2025", "2024-01-01T00:00:00", "2025-12-31T23:59:59"},
	} {
		candles, err := storage.LoadCandles(db, market, "60", period.from, period.to)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Load %s: %v\n", period.label, err)
			continue
		}
		if len(candles) < 100 {
			continue
		}

		fmt.Printf("━━━ %s (%d candles, 1h) ━━━\n\n", period.label, len(candles))

		var results []backtestResult
		for _, s := range strategies {
			trades := s.fn(candles, initBalance, feeRate, posSize)
			r := analyzeTrades(s.name, s.category, trades, initBalance)
			results = append(results, r)
		}

		printByCategory(results)

		best := results[0]
		for _, r := range results[1:] {
			if r.sharpe > best.sharpe {
				best = r
			}
		}
		fmt.Printf("  >>> Best by Sharpe: %s [%s] (Return: %+.2f%%, Sharpe: %.2f, MDD: %.2f%%)\n\n",
			best.name, best.category, best.totalReturn*100, best.sharpe, best.mdd*100)
	}
}

// ── Mean Reversion ─────────────────────────────────────────

func strategyRSI2(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 52 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		rsi2 := indicator.RSI(closes, 2)
		sma50 := indicator.SMA(closes, 50)
		if math.IsNaN(rsi2[i]) || math.IsNaN(sma50[i]) {
			return "hold", "", 0
		}
		price := c[i].Close
		if price > sma50[i] && rsi2[i] < 10 {
			return "buy", fmt.Sprintf("RSI2=%.1f", rsi2[i]), price - 2*lastATR(c, i)
		}
		if rsi2[i] > 70 {
			return "sell", fmt.Sprintf("RSI2=%.1f", rsi2[i]), 0
		}
		return "hold", "", 0
	})
}

func strategyRSIBB(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 20 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		rsi := indicator.RSI(closes, 14)
		u, _, l := bbAt(closes, i)
		if math.IsNaN(rsi[i]) || math.IsNaN(u) {
			return "hold", "", 0
		}
		price := c[i].Close
		bbPos := bbPosition(price, u, l)
		if rsi[i] < 35 && bbPos < 0.3 {
			return "buy", fmt.Sprintf("RSI=%.0f BB=%.0f%%", rsi[i], bbPos*100), price - 1.5*lastATR(c, i)
		}
		if rsi[i] > 65 || price > u {
			return "sell", "exit", 0
		}
		return "hold", "", 0
	})
}

func strategyBBPercentB(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 20 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		u, _, l := bbAt(closes, i)
		if math.IsNaN(u) {
			return "hold", "", 0
		}
		price := c[i].Close
		pctB := (price - l) / (u - l)
		if pctB < 0.05 && i > 0 {
			prevPctB := (c[i-1].Close - l) / (u - l)
			if prevPctB < 0 {
				return "buy", fmt.Sprintf("%%B recovery %.2f→%.2f", prevPctB, pctB), price - 2*lastATR(c, i)
			}
		}
		if pctB > 0.8 {
			return "sell", fmt.Sprintf("%%B=%.2f", pctB), 0
		}
		return "hold", "", 0
	})
}

func strategyStochastic(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 17 {
			return "hold", "", 0
		}
		k, d := stochastic(c, i, 14, 3)
		if math.IsNaN(k) || math.IsNaN(d) {
			return "hold", "", 0
		}
		if k < 20 && d < 20 {
			prevK, prevD := stochastic(c, i-1, 14, 3)
			if !math.IsNaN(prevK) && !math.IsNaN(prevD) && prevK <= prevD && k > d {
				return "buy", fmt.Sprintf("Stoch cross %.0f/%.0f", k, d), c[i].Close - 2*lastATR(c, i)
			}
		}
		if k > 80 {
			return "sell", fmt.Sprintf("Stoch %.0f", k), 0
		}
		return "hold", "", 0
	})
}

func strategyWilliamsR(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 14 {
			return "hold", "", 0
		}
		wr := williamsR(c, i, 14)
		if math.IsNaN(wr) {
			return "hold", "", 0
		}
		price := c[i].Close
		if wr < -80 {
			prevWR := williamsR(c, i-1, 14)
			if !math.IsNaN(prevWR) && wr > prevWR {
				return "buy", fmt.Sprintf("%%R=%.0f bounce", wr), price - 2*lastATR(c, i)
			}
		}
		if wr > -20 {
			return "sell", fmt.Sprintf("%%R=%.0f", wr), 0
		}
		return "hold", "", 0
	})
}

func strategyCCIReversion(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 20 {
			return "hold", "", 0
		}
		cci := calcCCI(c, i, 20)
		if math.IsNaN(cci) {
			return "hold", "", 0
		}
		price := c[i].Close
		if cci < -100 {
			return "buy", fmt.Sprintf("CCI=%.0f", cci), price - 2*lastATR(c, i)
		}
		if cci > 100 {
			return "sell", fmt.Sprintf("CCI=%.0f", cci), 0
		}
		return "hold", "", 0
	})
}

func strategyRSI2Volume(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 52 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		rsi2 := indicator.RSI(closes, 2)
		sma50 := indicator.SMA(closes, 50)
		if math.IsNaN(rsi2[i]) || math.IsNaN(sma50[i]) {
			return "hold", "", 0
		}
		price := c[i].Close
		vr := volumeRatio(c, i, 20)
		if price > sma50[i] && rsi2[i] < 10 && vr > 0.8 {
			return "buy", fmt.Sprintf("RSI2=%.1f vol=%.1fx", rsi2[i], vr), price - 2*lastATR(c, i)
		}
		if rsi2[i] > 70 {
			return "sell", fmt.Sprintf("RSI2=%.1f", rsi2[i]), 0
		}
		return "hold", "", 0
	})
}

// ── Trend Following ────────────────────────────────────────

func strategyEMACross(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 51 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		ema20 := indicator.EMA(closes, 20)
		ema50 := indicator.EMA(closes, 50)
		if math.IsNaN(ema20[i]) || math.IsNaN(ema50[i]) || math.IsNaN(ema20[i-1]) || math.IsNaN(ema50[i-1]) {
			return "hold", "", 0
		}
		prevDiff := ema20[i-1] - ema50[i-1]
		currDiff := ema20[i] - ema50[i]
		if prevDiff <= 0 && currDiff > 0 {
			return "buy", "EMA golden cross", c[i].Close - 2*lastATR(c, i)
		}
		if prevDiff >= 0 && currDiff < 0 {
			return "sell", "EMA death cross", 0
		}
		return "hold", "", 0
	})
}

func strategyMACDEMA(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 50 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		_, _, hist := indicator.MACD(closes, 12, 26, 9)
		if math.IsNaN(hist[i]) || math.IsNaN(hist[i-1]) {
			return "hold", "", 0
		}
		ema20 := indicator.EMA(closes, 20)
		ema50 := indicator.EMA(closes, 50)
		if math.IsNaN(ema20[i]) || math.IsNaN(ema50[i]) {
			return "hold", "", 0
		}
		if hist[i-1] <= 0 && hist[i] > 0 && ema20[i] > ema50[i] {
			return "buy", "MACD cross in uptrend", c[i].Close - 2*lastATR(c, i)
		}
		if hist[i-1] >= 0 && hist[i] < 0 {
			return "sell", "MACD bearish cross", 0
		}
		return "hold", "", 0
	})
}

func strategyDonchian(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 20 {
			return "hold", "", 0
		}
		high, low := donchian(c, i, 20)
		price := c[i].Close
		if price > high {
			return "buy", fmt.Sprintf("Donchian %.0f", high), price - 2*lastATR(c, i)
		}
		if price < low {
			return "sell", fmt.Sprintf("Donchian %.0f", low), 0
		}
		return "hold", "", 0
	})
}

func strategySMARibbon(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 50 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		sma10 := indicator.SMA(closes, 10)
		sma30 := indicator.SMA(closes, 30)
		sma50 := indicator.SMA(closes, 50)
		if math.IsNaN(sma10[i]) || math.IsNaN(sma30[i]) || math.IsNaN(sma50[i]) {
			return "hold", "", 0
		}
		bullish := sma10[i] > sma30[i] && sma30[i] > sma50[i]
		prevBull := !math.IsNaN(sma10[i-1]) && !math.IsNaN(sma30[i-1]) && !math.IsNaN(sma50[i-1]) &&
			sma10[i-1] > sma30[i-1] && sma30[i-1] > sma50[i-1]
		if !prevBull && bullish {
			return "buy", "SMA ribbon bullish", c[i].Close - 2*lastATR(c, i)
		}
		if prevBull && !bullish {
			return "sell", "SMA ribbon broken", 0
		}
		return "hold", "", 0
	})
}

func strategyATRBreakout(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 15 {
			return "hold", "", 0
		}
		atr := lastATR(c, i)
		if atr <= 0 {
			return "hold", "", 0
		}
		prevClose := c[i-1].Close
		price := c[i].Close
		if price > prevClose+1.5*atr {
			return "buy", fmt.Sprintf("ATR breakout +%.1fσ", (price-prevClose)/atr), price - 2*atr
		}
		if price < prevClose-1.5*atr {
			return "sell", fmt.Sprintf("ATR breakdown"), 0
		}
		return "hold", "", 0
	})
}

// ── Momentum ───────────────────────────────────────────────

func strategyRSIMomentum(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 15 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		rsi := indicator.RSI(closes, 14)
		if math.IsNaN(rsi[i]) || math.IsNaN(rsi[i-1]) {
			return "hold", "", 0
		}
		if rsi[i] > 30 && rsi[i-1] <= 30 {
			return "buy", fmt.Sprintf("RSI cross %.0f→%.0f", rsi[i-1], rsi[i]), c[i].Close - 2*lastATR(c, i)
		}
		if rsi[i] > 70 {
			return "sell", fmt.Sprintf("RSI=%.0f", rsi[i]), 0
		}
		return "hold", "", 0
	})
}

func strategyMACDHistogram(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 35 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		_, _, hist := indicator.MACD(closes, 12, 26, 9)
		if math.IsNaN(hist[i]) || math.IsNaN(hist[i-1]) || math.IsNaN(hist[i-2]) {
			return "hold", "", 0
		}
		if hist[i] > hist[i-1] && hist[i-1] > hist[i-2] && hist[i-2] < 0 {
			return "buy", "MACD hist momentum", c[i].Close - 2*lastATR(c, i)
		}
		if hist[i] < 0 && hist[i-1] >= 0 {
			return "sell", "MACD hist negative", 0
		}
		return "hold", "", 0
	})
}

func strategyROC(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 13 {
			return "hold", "", 0
		}
		roc := (c[i].Close - c[i-12].Close) / c[i-12].Close * 100
		prevROC := (c[i-1].Close - c[i-13].Close) / c[i-13].Close * 100
		if roc > -2 && prevROC <= -2 {
			return "buy", fmt.Sprintf("ROC %.1f→%.1f%%", prevROC, roc), c[i].Close - 2*lastATR(c, i)
		}
		if roc > 5 {
			return "sell", fmt.Sprintf("ROC=%.1f%%", roc), 0
		}
		return "hold", "", 0
	})
}

// ── Volatility ─────────────────────────────────────────────

func strategyBBSqueeze(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 40 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		u, _, l := bbAt(closes, i)
		if math.IsNaN(u) {
			return "hold", "", 0
		}
		bw := u - l
		minBW := bw
		for j := i - 20; j < i; j++ {
			uj, _, lj := bbAt(closes, j)
			if !math.IsNaN(uj) {
				w := uj - lj
				if w < minBW {
					minBW = w
				}
			}
		}
		price := c[i].Close
		vr := volumeRatio(c, i, 20)
		if bw <= minBW*1.05 && price > u && vr > 1.3 {
			return "buy", fmt.Sprintf("BB squeeze vol=%.1fx", vr), price - 2*lastATR(c, i)
		}
		return "hold", "", 0
	})
}

func strategyATRSqueeze(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 28 {
			return "hold", "", 0
		}
		atr := lastATR(c, i)
		prevATR := lastATR(c, i-14)
		if atr <= 0 || prevATR <= 0 {
			return "hold", "", 0
		}
		ratio := atr / prevATR
		price := c[i].Close
		if ratio < 0.7 {
			closes := closesUpTo(c, i)
			sma := indicator.SMA(closes, 20)
			if !math.IsNaN(sma[i]) && price > sma[i]*1.005 {
				return "buy", fmt.Sprintf("ATR squeeze %.2fx", ratio), price - 2*atr
			}
		}
		if ratio > 2.0 {
			return "sell", "ATR expansion", 0
		}
		return "hold", "", 0
	})
}

func strategyKeltner(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 20 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		ema20 := indicator.EMA(closes, 20)
		atr := lastATR(c, i)
		if math.IsNaN(ema20[i]) || atr <= 0 {
			return "hold", "", 0
		}
		upper := ema20[i] + 2*atr
		lower := ema20[i] - 2*atr
		price := c[i].Close
		if price > upper {
			return "buy", "Keltner upper", price - 2*atr
		}
		if price < lower {
			return "sell", "Keltner lower", 0
		}
		return "hold", "", 0
	})
}

// ── Engine ─────────────────────────────────────────────────

func runStrategy(candles []indicator.CandleData, balance, feeRate, posSize float64, evaluate func([]indicator.CandleData, int) (string, string, float64)) []trade {
	type position struct {
		price, volume, cost, stop float64
		entryTime                 string
	}
	var pos *position
	var trades []trade
	peak := balance

	for i := 0; i < len(candles); i++ {
		price := candles[i].Close

		if pos != nil {
			exit := false
			reason := ""
			if pos.stop > 0 && price <= pos.stop {
				exit = true
				reason = fmt.Sprintf("trailing stop (%.0f)", pos.stop)
			}
			if !exit {
				sig, r, _ := evaluate(candles, i)
				if sig == "sell" {
					exit = true
					reason = r
				}
			}
			if !exit {
				atr := lastATR(candles, i)
				if atr > 0 {
					if trail := price - 2*atr; trail > pos.stop {
						pos.stop = trail
					}
				}
			}
			if exit {
				exitVal := pos.volume * price
				fee := exitVal * feeRate
				pnl := exitVal - pos.cost - fee
				balance += exitVal - fee
				trades = append(trades, trade{pos.price, pos.entryTime, price, candles[i].Time, pnl, reason})
				pos = nil
			}
		}

		if pos == nil && i+1 < len(candles) {
			sig, _, stop := evaluate(candles, i)
			if sig == "buy" {
				size := (balance * posSize) / price
				cost := size * price
				entryFee := cost * feeRate
				if cost+entryFee > balance {
					continue
				}
				balance -= cost + entryFee
				pos = &position{price, size, cost, stop, candles[i].Time}
			}
		}

		equity := balance
		if pos != nil {
			equity += pos.volume * price
		}
		if equity > peak {
			peak = equity
		}
	}

	if pos != nil {
		last := candles[len(candles)-1].Close
		exitVal := pos.volume * last
		fee := exitVal * feeRate
		pnl := exitVal - pos.cost - fee
		balance += exitVal - fee
		trades = append(trades, trade{pos.price, pos.entryTime, last, candles[len(candles)-1].Time, pnl, "end"})
	}

	return trades
}

func analyzeTrades(name, category string, trades []trade, initBalance float64) backtestResult {
	balance := initBalance
	wins := 0
	peak := initBalance
	maxDD := 0.0
	monthlyRet := make(map[string]float64)

	for _, t := range trades {
		balance += t.pnl
		if t.pnl > 0 {
			wins++
		}
		if balance > peak {
			peak = balance
		}
		if dd := (peak - balance) / peak; dd > maxDD {
			maxDD = dd
		}
		monthlyRet[t.exitTime[:7]] += t.pnl
	}

	lossMonths := 0
	for _, pnl := range monthlyRet {
		if pnl < 0 {
			lossMonths++
		}
	}

	winRate := 0.0
	if len(trades) > 0 {
		winRate = float64(wins) / float64(len(trades))
	}

	sharpe := 0.0
	if maxDD > 0 {
		sharpe = ((balance-initBalance)/initBalance - 0.02) / maxDD
	}

	return backtestResult{name, category, (balance - initBalance) / initBalance, winRate, maxDD,
		len(trades), wins, balance, sharpe, monthlyRet, lossMonths}
}

func printByCategory(results []backtestResult) {
	for _, cat := range []string{"MeanReversion", "TrendFollow", "Momentum", "Volatility", "Hybrid", "Williams"} {
		fmt.Printf("  [%s]\n", cat)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "    Strategy\tReturn\tWinRate\tMDD\tTrades\tSharpe\tLossMo\n")
		for _, r := range results {
			if r.category != cat {
				continue
			}
			fmt.Fprintf(w, "    %s\t%+.2f%%\t%.0f%%\t%.2f%%\t%d\t%.2f\t%d\n",
				r.name, r.totalReturn*100, r.winRate*100, r.mdd*100, r.trades, r.sharpe, r.lossMonths)
		}
		w.Flush()
		fmt.Println()
	}
}

// ── Indicator Helpers ──────────────────────────────────────

func closesUpTo(c []indicator.CandleData, i int) []float64 {
	r := make([]float64, i+1)
	for j := 0; j <= i; j++ {
		r[j] = c[j].Close
	}
	return r
}

func lastNonNaN(vals []float64) float64 {
	for i := len(vals) - 1; i >= 0; i-- {
		if !math.IsNaN(vals[i]) {
			return vals[i]
		}
	}
	return 0
}

func lastATR(c []indicator.CandleData, i int) float64 {
	if i < 14 {
		return 0
	}
	return lastNonNaN(indicator.ATR(c[:i+1], 14))
}

func bbAt(closes []float64, i int) (upper, middle, lower float64) {
	u, m, l := indicator.BollingerBands(closes, 20, 2.0)
	return u[i], m[i], l[i]
}

func bbPosition(price, upper, lower float64) float64 {
	if bw := upper - lower; bw != 0 {
		return (price - lower) / bw
	}
	return 0.5
}

func volumeRatio(c []indicator.CandleData, i, period int) float64 {
	if i < period {
		return 1
	}
	var sum float64
	for j := i - period; j < i; j++ {
		sum += c[j].Volume
	}
	if avg := sum / float64(period); avg > 0 {
		return c[i].Volume / avg
	}
	return 1
}

func stochastic(c []indicator.CandleData, i, kPeriod, dPeriod int) (k, d float64) {
	if i < kPeriod {
		return math.NaN(), math.NaN()
	}
	var lowN, highN float64 = c[i].Low, c[i].High
	for j := i - kPeriod + 1; j <= i; j++ {
		if c[j].Low < lowN {
			lowN = c[j].Low
		}
		if c[j].High > highN {
			highN = c[j].High
		}
	}
	if highN == lowN {
		return 50, 50
	}
	k = (c[i].Close - lowN) / (highN - lowN) * 100
	if i < kPeriod+dPeriod {
		return k, k
	}
	var sum float64
	for j := 0; j < dPeriod; j++ {
		var lN, hN float64 = c[i-j].Low, c[i-j].High
		for m := i - j - kPeriod + 1; m <= i-j; m++ {
			if c[m].Low < lN {
				lN = c[m].Low
			}
			if c[m].High > hN {
				hN = c[m].High
			}
		}
		if hN != lN {
			sum += (c[i-j].Close - lN) / (hN - lN) * 100
		}
	}
	return k, sum / float64(dPeriod)
}

func williamsR(c []indicator.CandleData, i, period int) float64 {
	if i < period {
		return math.NaN()
	}
	var hN, lN float64 = c[i].High, c[i].Low
	for j := i - period + 1; j <= i; j++ {
		if c[j].High > hN {
			hN = c[j].High
		}
		if c[j].Low < lN {
			lN = c[j].Low
		}
	}
	if hN == lN {
		return -50
	}
	return (hN - c[i].Close) / (hN - lN) * -100
}

func calcCCI(c []indicator.CandleData, i, period int) float64 {
	if i < period {
		return math.NaN()
	}
	tps := make([]float64, i+1)
	for j := 0; j <= i; j++ {
		tps[j] = (c[j].High + c[j].Low + c[j].Close) / 3
	}
	smaTP := indicator.SMA(tps, period)
	if math.IsNaN(smaTP[i]) {
		return math.NaN()
	}
	var md float64
	for j := i - period + 1; j <= i; j++ {
		md += math.Abs(tps[j] - smaTP[i])
	}
	md /= float64(period)
	if md == 0 {
		return 0
	}
	return (tps[i] - smaTP[i]) / (0.015 * md)
}

func donchian(c []indicator.CandleData, i, period int) (high, low float64) {
	high = c[i-1].High
	low = c[i-1].Low
	for j := i - period; j < i; j++ {
		if j < 0 {
			continue
		}
		if c[j].High > high {
			high = c[j].High
		}
		if c[j].Low < low {
			low = c[j].Low
		}
	}
	return
}

// ── Hybrid Strategies ──────────────────────────────────────

// detectRegime classifies current market state.
// Returns: "trend", "reversion", "volatile", "quiet"
func detectRegime(c []indicator.CandleData, i int) string {
	if i < 50 {
		return "quiet"
	}
	closes := closesUpTo(c, i)

	// ADX-like: trend strength via EMA separation
	ema20 := indicator.EMA(closes, 20)
	ema50 := indicator.EMA(closes, 50)
	if math.IsNaN(ema20[i]) || math.IsNaN(ema50[i]) {
		return "quiet"
	}
	trendStrength := math.Abs(ema20[i]-ema50[i]) / ema50[i] * 100

	// Volatility: BB bandwidth relative to price
	u, _, l := bbAt(closes, i)
	if math.IsNaN(u) {
		return "quiet"
	}
	bbWidth := (u - l) / ((u + l) / 2) * 100

	// RSI for mean reversion context
	rsi := indicator.RSI(closes, 14)
	rsiVal := rsi[i]

	switch {
	case bbWidth > 6 && trendStrength > 1.5:
		return "trend" // strong directional move with volatility
	case bbWidth < 3:
		return "reversion" // low volatility → mean reversion works
	case trendStrength > 1.0:
		return "trend"
	case bbWidth > 5:
		return "volatile"
	default:
		if !math.IsNaN(rsiVal) && (rsiVal < 30 || rsiVal > 70) {
			return "reversion"
		}
		return "quiet"
	}
}

// strategyRegimeSwitch selects strategy based on detected market regime.
func strategyRegimeSwitch(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		regime := detectRegime(c, i)
		price := c[i].Close
		atr := lastATR(c, i)

		switch regime {
		case "reversion":
			// RSI-2 oversold in uptrend
			if i < 52 {
				return "hold", "", 0
			}
			closes := closesUpTo(c, i)
			rsi2 := indicator.RSI(closes, 2)
			sma50 := indicator.SMA(closes, 50)
			if math.IsNaN(rsi2[i]) || math.IsNaN(sma50[i]) {
				return "hold", "", 0
			}
			if price > sma50[i] && rsi2[i] < 10 {
				return "buy", fmt.Sprintf("[reversion] RSI2=%.1f", rsi2[i]), price - 2*atr
			}
			if rsi2[i] > 70 {
				return "sell", "[reversion] RSI2 overbought", 0
			}

		case "volatile":
			// Keltner breakout
			if i < 20 {
				return "hold", "", 0
			}
			closes := closesUpTo(c, i)
			ema20 := indicator.EMA(closes, 20)
			if math.IsNaN(ema20[i]) || atr <= 0 {
				return "hold", "", 0
			}
			upper := ema20[i] + 2*atr
			lower := ema20[i] - 2*atr
			if price > upper {
				return "buy", "[volatile] Keltner breakout", price - 2*atr
			}
			if price < lower {
				return "sell", "[volatile] Keltner breakdown", 0
			}

		case "trend":
			// SMA Ribbon alignment
			if i < 50 {
				return "hold", "", 0
			}
			closes := closesUpTo(c, i)
			sma10 := indicator.SMA(closes, 10)
			sma30 := indicator.SMA(closes, 30)
			sma50 := indicator.SMA(closes, 50)
			if math.IsNaN(sma10[i]) || math.IsNaN(sma30[i]) || math.IsNaN(sma50[i]) {
				return "hold", "", 0
			}
			bull := sma10[i] > sma30[i] && sma30[i] > sma50[i]
			prevBull := !math.IsNaN(sma10[i-1]) && !math.IsNaN(sma30[i-1]) && !math.IsNaN(sma50[i-1]) &&
				sma10[i-1] > sma30[i-1] && sma30[i-1] > sma50[i-1]
			if !prevBull && bull {
				return "buy", "[trend] SMA ribbon bullish", price - 2*atr
			}
			if prevBull && !bull {
				return "sell", "[trend] SMA ribbon broken", 0
			}
		}
		return "hold", "", 0
	})
}

// strategyVoting uses 3 strategies voting — trade when at least 2 agree.
func strategyVoting(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		price := c[i].Close
		atr := lastATR(c, i)

		buyVotes := 0
		sellVotes := 0
		var reasons []string

		// Vote 1: RSI-2
		if i >= 52 {
			closes := closesUpTo(c, i)
			rsi2 := indicator.RSI(closes, 2)
			sma50 := indicator.SMA(closes, 50)
			if !math.IsNaN(rsi2[i]) && !math.IsNaN(sma50[i]) {
				if price > sma50[i] && rsi2[i] < 10 {
					buyVotes++
					reasons = append(reasons, "RSI2")
				} else if rsi2[i] > 70 {
					sellVotes++
				}
			}
		}

		// Vote 2: RSI(14)+BB
		if i >= 20 {
			closes := closesUpTo(c, i)
			rsi := indicator.RSI(closes, 14)
			u, _, l := bbAt(closes, i)
			if !math.IsNaN(rsi[i]) && !math.IsNaN(u) {
				bbPos := bbPosition(price, u, l)
				if rsi[i] < 35 && bbPos < 0.3 {
					buyVotes++
					reasons = append(reasons, "RSI+BB")
				} else if rsi[i] > 65 || price > u {
					sellVotes++
				}
			}
		}

		// Vote 3: Keltner
		if i >= 20 {
			closes := closesUpTo(c, i)
			ema20 := indicator.EMA(closes, 20)
			if !math.IsNaN(ema20[i]) && atr > 0 {
				if price > ema20[i]+2*atr {
					buyVotes++
					reasons = append(reasons, "Keltner")
				} else if price < ema20[i]-2*atr {
					sellVotes++
				}
			}
		}

		if buyVotes >= 2 {
			return "buy", fmt.Sprintf("[vote %d] %v", buyVotes, reasons), price - 2*atr
		}
		if sellVotes >= 2 {
			return "sell", fmt.Sprintf("[vote %d]", sellVotes), 0
		}
		return "hold", "", 0
	})
}

// strategyMRVolConditional uses MeanReversion primary, Volatility as confirmation.
func strategyMRVolConditional(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 52 {
			return "hold", "", 0
		}
		price := c[i].Close
		atr := lastATR(c, i)
		closes := closesUpTo(c, i)

		// Primary: RSI-2 mean reversion signal
		rsi2 := indicator.RSI(closes, 2)
		sma50 := indicator.SMA(closes, 50)
		if math.IsNaN(rsi2[i]) || math.IsNaN(sma50[i]) {
			return "hold", "", 0
		}

		rsi2Buy := price > sma50[i] && rsi2[i] < 10
		rsi2Sell := rsi2[i] > 70

		// Confirmation: BB bandwidth (volatility state)
		u, _, l := bbAt(closes, i)
		if math.IsNaN(u) {
			return "hold", "", 0
		}
		bbWidth := (u - l) / ((u + l) / 2) * 100

		// Low vol (bbWidth < 4) → full confidence in mean reversion
		// High vol (bbWidth > 6) → also accept, but use wider stop
		// Medium → only trade if RSI2 < 5 (stronger signal)

		if rsi2Buy {
			if bbWidth < 4 {
				return "buy", fmt.Sprintf("[MR+Vol] RSI2=%.1f lowVol=%.1f%%", rsi2[i], bbWidth), price - 2*atr
			}
			if bbWidth > 6 {
				return "buy", fmt.Sprintf("[MR+Vol] RSI2=%.1f highVol=%.1f%%", rsi2[i], bbWidth), price - 3*atr
			}
			if rsi2[i] < 5 {
				return "buy", fmt.Sprintf("[MR+Vol] RSI2=%.1f deep oversold", rsi2[i]), price - 2*atr
			}
		}

		if rsi2Sell {
			return "sell", fmt.Sprintf("[MR+Vol] RSI2=%.1f overbought", rsi2[i]), 0
		}
		return "hold", "", 0
	})
}

// strategyBest2Adaptive combines RSI-2 (mean reversion) and Keltner (volatility).
// Uses RSI-2 in low-volatility regimes, Keltner in high-volatility regimes,
// and requires both to agree when volatility is moderate.
func strategyBest2Adaptive(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 52 {
			return "hold", "", 0
		}
		price := c[i].Close
		atr := lastATR(c, i)
		closes := closesUpTo(c, i)

		// Regime detection via BB width
		u, _, l := bbAt(closes, i)
		if math.IsNaN(u) {
			return "hold", "", 0
		}
		bbWidth := (u - l) / ((u + l) / 2) * 100

		// RSI-2 signal
		rsi2 := indicator.RSI(closes, 2)
		sma50 := indicator.SMA(closes, 50)
		rsi2Buy := !math.IsNaN(rsi2[i]) && !math.IsNaN(sma50[i]) && price > sma50[i] && rsi2[i] < 10
		rsi2Sell := !math.IsNaN(rsi2[i]) && rsi2[i] > 70

		// Keltner signal
		ema20 := indicator.EMA(closes, 20)
		keltBuy := !math.IsNaN(ema20[i]) && atr > 0 && price > ema20[i]+2*atr
		keltSell := !math.IsNaN(ema20[i]) && atr > 0 && price < ema20[i]-2*atr

		// Low volatility → trust mean reversion
		if bbWidth < 4 {
			if rsi2Buy {
				return "buy", fmt.Sprintf("[lowVol→MR] RSI2=%.1f", rsi2[i]), price - 2*atr
			}
			if rsi2Sell {
				return "sell", "[lowVol→MR] RSI2 overbought", 0
			}
		}

		// High volatility → trust volatility strategy
		if bbWidth > 5 {
			if keltBuy {
				return "buy", fmt.Sprintf("[highVol→Vol] BBW=%.1f%%", bbWidth), price - 2*atr
			}
			if keltSell {
				return "sell", "[highVol→Vol] Keltner breakdown", 0
			}
		}

		// Moderate → both must agree
		if rsi2Buy && keltBuy {
			return "buy", "[both agree] MR+Vol bullish", price - 2*atr
		}
		if rsi2Sell && keltSell {
			return "sell", "[both agree] MR+Vol bearish", 0
		}
		return "hold", "", 0
	})
}

// ── Larry Williams Volatility Breakout ─────────────────────

// strategyLarryWilliams implements the classic Larry Williams volatility breakout:
// Buy when price breaks above previous day's high + K% of (high-low) range.
// K = 0.5 (50% of range). Exit at close or next day.
func strategyLarryWilliams(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 24 {
			return "hold", "", 0
		}

		// Use previous candle's range (1h candle as "day" proxy)
		prevHigh := c[i-1].High
		prevLow := c[i-1].Low
		prevRange := prevHigh - prevLow
		if prevRange <= 0 {
			return "hold", "", 0
		}

		price := c[i].Close
		breakoutLevel := prevHigh + 0.5*prevRange

		if price > breakoutLevel {
			atr := lastATR(c, i)
			return "buy", fmt.Sprintf("LW breakout %.0f (range=%.0f)", breakoutLevel, prevRange), price - 2*atr
		}
		if price < prevLow {
			return "sell", "below prev low", 0
		}
		return "hold", "", 0
	})
}

// strategyLWMA10 adds MA10 uptrend filter to Larry Williams breakout.
// Only buy when price is above MA10 (uptrend confirmation).
func strategyLWMA10(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 24 {
			return "hold", "", 0
		}

		closes := closesUpTo(c, i)
		ma10 := indicator.SMA(closes, 10)
		if math.IsNaN(ma10[i]) {
			return "hold", "", 0
		}

		price := c[i].Close

		// Only trade in uptrend: price > MA10
		if price <= ma10[i] {
			return "hold", "", 0
		}

		prevHigh := c[i-1].High
		prevLow := c[i-1].Low
		prevRange := prevHigh - prevLow
		if prevRange <= 0 {
			return "hold", "", 0
		}

		breakoutLevel := prevHigh + 0.5*prevRange

		if price > breakoutLevel {
			atr := lastATR(c, i)
			return "buy", fmt.Sprintf("LW+MA10 breakout (MA10=%.0f)", ma10[i]), price - 2*atr
		}

		// Sell when price drops below MA10
		if price < ma10[i]*0.99 {
			return "sell", "below MA10", 0
		}
		return "hold", "", 0
	})
}

// strategyRSI2LWPipeline combines RSI-2 entry with Larry Williams volatility confirmation.
// Step 1: RSI-2 detects oversold condition in uptrend → candidate
// Step 2: Larry Williams breakout confirms entry with momentum
// Step 3: Exit on RSI-2 overbought or trailing stop
func strategyRSI2LWPipeline(candles []indicator.CandleData, balance, feeRate, posSize float64) []trade {
	// Track RSI-2 oversold state
	rsi2Oversold := false
	rsi2OversoldIdx := 0

	return runStrategy(candles, balance, feeRate, posSize, func(c []indicator.CandleData, i int) (string, string, float64) {
		if i < 52 {
			return "hold", "", 0
		}
		closes := closesUpTo(c, i)
		price := c[i].Close
		atr := lastATR(c, i)

		// RSI-2 state tracking
		rsi2 := indicator.RSI(closes, 2)
		sma50 := indicator.SMA(closes, 50)
		if math.IsNaN(rsi2[i]) || math.IsNaN(sma50[i]) {
			return "hold", "", 0
		}

		// Detect RSI-2 oversold in uptrend
		if price > sma50[i] && rsi2[i] < 10 {
			rsi2Oversold = true
			rsi2OversoldIdx = i
		}

		// RSI-2 overbought → exit
		if rsi2[i] > 70 {
			rsi2Oversold = false
			return "sell", fmt.Sprintf("RSI2=%.1f overbought", rsi2[i]), 0
		}

		// If RSI-2 was recently oversold (within last 5 candles), look for LW breakout
		if rsi2Oversold && i-rsi2OversoldIdx <= 5 {
			prevHigh := c[i-1].High
			prevLow := c[i-1].Low
			prevRange := prevHigh - prevLow
			if prevRange > 0 {
				breakoutLevel := prevHigh + 0.3*prevRange // Lower threshold since RSI-2 already confirmed
				if price > breakoutLevel {
					rsi2Oversold = false
					return "buy", fmt.Sprintf("RSI2→LW RSI2=%.1f breakout", rsi2[i]), price - 2*atr
				}
			}
		}

		// MA10 trend filter for additional sell
		sma10 := indicator.SMA(closes, 10)
		if !math.IsNaN(sma10[i]) && price < sma10[i]*0.99 {
			rsi2Oversold = false
			return "sell", "below MA10", 0
		}

		return "hold", "", 0
	})
}
