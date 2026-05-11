package main

import (
	"fmt"
	"math"
	"os"
	"text/tabwriter"

	"github.com/kyungw00k/upbit/trader/indicator"
	"github.com/kyungw00k/upbit/trader/storage"
)

// Portfolio v2: RSI-2 (MeanReversion) + Keltner Breakout (Volatility)
// Truly different signal sources → genuine diversification
// Weighted by Sharpe ratio instead of equal split

type subTrade struct {
	strategy   string
	entryPrice float64
	entryTime  string
	exitPrice  float64
	exitTime   string
	pnl        float64
	reason     string
}

type subPosition struct {
	strategy  string
	price     float64
	volume    float64
	cost      float64
	entryTime string
	stop      float64
}

type signalFunc func(c []indicator.CandleData, i int) (signal, reason string, stop float64)

func main() {
	dbPath := "data/candles.db"
	market := "KRW-BTC"
	initBalance := 330_000.0

	db, err := storage.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB open: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Printf("=== Portfolio v2: RSI-2 + Keltner | %s | Initial: %.0f KRW ===\n\n", market, initBalance)

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
		if err != nil || len(candles) < 100 {
			continue
		}
		fmt.Printf("━━━ %s (%d candles) ━━━\n\n", period.label, len(candles))

		// Phase 1: Run each strategy standalone to get Sharpe ratios
		sharpes := runStandalone(candles, initBalance)

		// Phase 2: Weighted portfolio
		runWeightedPortfolio(candles, initBalance, sharpes)
	}
}

// Run each strategy standalone to determine weights
func runStandalone(candles []indicator.CandleData, initBalance float64) map[string]float64 {
	strategies := getStrategies()
	sharpes := make(map[string]float64)

	fmt.Println("--- Standalone Performance ---")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Strategy\tReturn\tMDD\tTrades\tWinRate\tSharpe\n")

	for _, s := range strategies {
		bal, _, mdd, trades := runSingleStrategy(candles, initBalance, s.sig)
		ret := (bal - initBalance) / initBalance
		wins := 0
		for _, t := range trades {
			if t.pnl > 0 {
				wins++
			}
		}
		wr := 0.0
		if len(trades) > 0 {
			wr = float64(wins) / float64(len(trades)) * 100
		}
		sharpe := 0.0
		if mdd > 0 {
			sharpe = (ret - 0.02) / mdd
		}
		sharpes[s.name] = sharpe
		fmt.Fprintf(w, "%s\t%+.2f%%\t%.2f%%\t%d\t%.0f%%\t%.2f\n",
			s.name, ret*100, mdd*100, len(trades), wr, sharpe)
	}
	w.Flush()
	fmt.Println()

	return sharpes
}

func getStrategies() []struct {
	name string
	sig  signalFunc
} {
	return []struct {
		name string
		sig  signalFunc
	}{
		{"RSI-2", sigRSI2},
		{"Keltner", sigKeltner},
	}
}

func runSingleStrategy(candles []indicator.CandleData, initBalance float64, sig signalFunc) (float64, []subTrade, float64, []subTrade) {
	feeRate := 0.0005
	posSize := 0.30
	balance := initBalance

	type position struct {
		price, volume, cost, stop float64
		entryTime                 string
	}
	var pos *position
	var trades []subTrade
	peak := initBalance
	maxDD := 0.0

	for i := 0; i < len(candles); i++ {
		price := candles[i].Close

		if pos != nil {
			exit := false
			reason := ""

			atr := lastATR(candles, i)
			if pos.stop > 0 && price <= pos.stop {
				exit = true
				reason = fmt.Sprintf("trailing stop (%.0f)", pos.stop)
			}

			if !exit {
				sig2, r, _ := sig(candles, i)
				if sig2 == "sell" {
					exit = true
					reason = r
				}
			}

			if !exit {
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
				trades = append(trades, subTrade{"", pos.price, pos.entryTime, price,
					candles[i].Time, pnl, reason})
				pos = nil
			}
		}

		if pos == nil && i+1 < len(candles) {
			sig2, _, stop := sig(candles, i)
			if sig2 == "buy" {
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
		if dd := (peak - equity) / peak; dd > maxDD {
			maxDD = dd
		}
	}

	// Close remaining
	if pos != nil {
		lastPrice := candles[len(candles)-1].Close
		exitVal := pos.volume * lastPrice
		fee := exitVal * feeRate
		pnl := exitVal - pos.cost - fee
		balance += exitVal - fee
		trades = append(trades, subTrade{"", pos.price, pos.entryTime, lastPrice,
			candles[len(candles)-1].Time, pnl, "end of data"})
	}

	return balance, trades, maxDD, trades
}

func runWeightedPortfolio(candles []indicator.CandleData, initBalance float64, sharpes map[string]float64) {
	feeRate := 0.0005

	strategies := getStrategies()

	// Calculate weights from Sharpe (min 0.2 each)
	totalSharpe := 0.0
	for _, s := range strategies {
		w := math.Max(sharpes[s.name], 0.1)
		totalSharpe += w
	}
	weights := make(map[string]float64)
	for _, s := range strategies {
		w := math.Max(sharpes[s.name], 0.1)
		weights[s.name] = w / totalSharpe
	}

	fmt.Printf("--- Weighted Portfolio (Sharpe-based) ---\n")
	for _, s := range strategies {
		fmt.Printf("  %s: weight=%.0f%% (sharpe=%.2f)\n", s.name, weights[s.name]*100, sharpes[s.name])
	}
	fmt.Println()

	positions := make(map[string]*subPosition)
	var allTrades []subTrade
	balances := make(map[string]float64)
	for _, s := range strategies {
		balances[s.name] = initBalance * weights[s.name]
	}

	peak := initBalance
	maxDD := 0.0
	monthlyPnL := make(map[string]float64)
	strategyMonthly := make(map[string]map[string]float64)
	for _, s := range strategies {
		strategyMonthly[s.name] = make(map[string]float64)
	}

	// Track correlation: how often strategies trade at the same time
	type entry struct{ strat, time string }
	var entries []entry

	for i := 0; i < len(candles); i++ {
		price := candles[i].Close

		for _, s := range strategies {
			pos := positions[s.name]
			bal := balances[s.name]

			if pos != nil {
				exit := false
				reason := ""

				atr := lastATR(candles, i)
				if pos.stop > 0 && price <= pos.stop {
					exit = true
					reason = fmt.Sprintf("trailing stop (%.0f)", pos.stop)
				}

				if !exit {
					sig, r, _ := s.sig(candles, i)
					if sig == "sell" {
						exit = true
						reason = r
					}
				}

				if !exit {
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
					bal += exitVal - fee
					allTrades = append(allTrades, subTrade{s.name, pos.price, pos.entryTime, price, candles[i].Time, pnl, reason})
					month := candles[i].Time[:7]
					monthlyPnL[month] += pnl
					strategyMonthly[s.name][month] += pnl

					fmt.Printf("  [%-8s] SELL @ %10.0f (%s) PnL: %+8.0f | %s\n",
						s.name, price, candles[i].Time[:16], pnl, reason)

					delete(positions, s.name)
					balances[s.name] = bal
				}
			}

			if positions[s.name] == nil && i+1 < len(candles) {
				sig, _, stop := s.sig(candles, i)
				if sig == "buy" {
					size := (bal * 0.90) / price
					cost := size * price
					entryFee := cost * feeRate
					if cost+entryFee > bal {
						continue
					}
					bal -= cost + entryFee
					positions[s.name] = &subPosition{s.name, price, size, cost, candles[i].Time, stop}
					balances[s.name] = bal

					entries = append(entries, entry{s.name, candles[i].Time[:16]})
					fmt.Printf("  [%-8s] BUY  @ %10.0f (%s) stop=%.0f\n",
						s.name, price, candles[i].Time[:16], stop)
				}
			}
		}

		equity := 0.0
		for _, s := range strategies {
			equity += balances[s.name]
			if pos, ok := positions[s.name]; ok {
				equity += pos.volume * price
			}
		}
		if equity > peak {
			peak = equity
		}
		if dd := (peak - equity) / peak; dd > maxDD {
			maxDD = dd
		}
	}

	// Close remaining
	lastPrice := candles[len(candles)-1].Close
	for _, s := range strategies {
		if pos, ok := positions[s.name]; ok {
			exitVal := pos.volume * lastPrice
			fee := exitVal * feeRate
			pnl := exitVal - pos.cost - fee
			balances[s.name] += exitVal - fee
			allTrades = append(allTrades, subTrade{s.name, pos.price, pos.entryTime, lastPrice,
				candles[len(candles)-1].Time, pnl, "end of data"})
			monthlyPnL[candles[len(candles)-1].Time[:7]] += pnl
			strategyMonthly[s.name][candles[len(candles)-1].Time[:7]] += pnl
		}
	}

	// Check signal overlap
	overlapCount := 0
	for i, e1 := range entries {
		for _, e2 := range entries[i+1:] {
			if e1.strat != e2.strat && e1.time == e2.time {
				overlapCount++
			}
		}
	}
	fmt.Printf("\n  Signal overlap: %d/%d entries at same time\n", overlapCount, len(entries))

	// Results
	fmt.Printf("\n--- Portfolio v2 Summary ---\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Strategy\tAlloc\tBalance\tReturn\tTrades\tWinRate\n")
	totalBalance := 0.0
	for _, s := range strategies {
		bal := balances[s.name]
		totalBalance += bal
		subInit := initBalance * weights[s.name]
		wins, count := 0, 0
		for _, t := range allTrades {
			if t.strategy == s.name {
				count++
				if t.pnl > 0 {
					wins++
				}
			}
		}
		wr := 0.0
		if count > 0 {
			wr = float64(wins) / float64(count) * 100
		}
		fmt.Fprintf(w, "%s\t%.0f%%\t%.0f\t%+.2f%%\t%d\t%.0f%%\n",
			s.name, weights[s.name]*100, bal, (bal-subInit)/subInit*100, count, wr)
	}
	fmt.Fprintf(w, "COMBINED\t---\t%.0f\t%+.2f%%\t%d\t---\n",
		totalBalance, (totalBalance-initBalance)/initBalance*100, len(allTrades))
	w.Flush()

	totalReturn := (totalBalance - initBalance) / initBalance
	sharpe := 0.0
	if maxDD > 0 {
		sharpe = (totalReturn - 0.02) / maxDD
	}
	fmt.Printf("\nTotal Return: %+.2f%% (%+.0f KRW)\n", totalReturn*100, totalBalance-initBalance)
	fmt.Printf("Max Drawdown: %.2f%%\n", maxDD*100)
	fmt.Printf("Sharpe:       %.2f\n", sharpe)
	fmt.Printf("Total Trades: %d\n", len(allTrades))
	wins := 0
	for _, t := range allTrades {
		if t.pnl > 0 {
			wins++
		}
	}
	if len(allTrades) > 0 {
		fmt.Printf("Win Rate:     %.0f%% (%d/%d)\n", float64(wins)/float64(len(allTrades))*100, wins, len(allTrades))
	}

	// Monthly breakdown
	fmt.Printf("\n--- Monthly Returns ---\n")
	months := sortedKeys(monthlyPnL)
	runningBal := initBalance
	w2 := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w2, "Month\tRSI-2\tKeltner\tCombined\tCumul%%\n")
	for _, m := range months {
		r1 := strategyMonthly["RSI-2"][m]
		r2 := strategyMonthly["Keltner"][m]
		combined := r1 + r2
		runningBal += combined
		cumul := (runningBal - initBalance) / initBalance * 100
		fmt.Fprintf(w2, "%s\t%+8.0f\t%+8.0f\t%+8.0f\t%+.2f%%\n",
			m, r1, r2, combined, cumul)
	}
	w2.Flush()
	fmt.Println()
}

func sortedKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}

// ── Strategy Signals ───────────────────────────────────────

func sigRSI2(c []indicator.CandleData, i int) (string, string, float64) {
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
		return "buy", fmt.Sprintf("RSI2=%.1f oversold", rsi2[i]), price - 2*lastATR(c, i)
	}
	if rsi2[i] > 70 {
		return "sell", fmt.Sprintf("RSI2=%.1f overbought", rsi2[i]), 0
	}
	return "hold", "", 0
}

func sigKeltner(c []indicator.CandleData, i int) (string, string, float64) {
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
		return "buy", fmt.Sprintf("Keltner upper (%.0f)", upper), price - 2*atr
	}
	if price < lower {
		return "sell", fmt.Sprintf("Keltner lower (%.0f)", lower), 0
	}
	return "hold", "", 0
}

// ── Helpers ────────────────────────────────────────────────

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
