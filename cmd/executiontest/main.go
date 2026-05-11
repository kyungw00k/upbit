package main

import (
	"fmt"
	"math"
	"os"
	"text/tabwriter"

	"github.com/kyungw00k/upbit/trader/indicator"
	"github.com/kyungw00k/upbit/trader/storage"
)

type optTrade struct {
	pnl float64
}

func main() {
	dbPath := "data/candles.db"
	market := "KRW-BTC"
	initBalance := 330_000.0
	feeRate := 0.0005
	posSize := 0.30
	trailMult := 3.0

	db, err := storage.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB open: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Printf("=== 2026 Out-of-Sample Test ===\n")
	fmt.Printf("Strategy: Keltner Breakout | Pos: %.0f%% | Trail: %.1f×ATR\n\n", posSize*100, trailMult)

	periods := []struct {
		label string
		from  string
		to    string
	}{
		{"2024", "2024-01-01T00:00:00", "2024-12-31T23:59:59"},
		{"2025", "2025-01-01T00:00:00", "2025-12-31T23:59:59"},
		{"2026 (YTD)", "2026-01-01T00:00:00", "2026-12-31T23:59:59"},
		{"2024-2026", "2024-01-01T00:00:00", "2026-12-31T23:59:59"},
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Period\tReturn\tMDD\tSharpe\tTrades\tWinRate\tPnL(KRW)\n")

	for _, p := range periods {
		candles, err := storage.LoadCandles(db, market, "60", p.from, p.to)
		if err != nil || len(candles) < 100 {
			fmt.Printf("[%s] insufficient data (%d candles)\n", p.label, len(candles))
			continue
		}

		balance, mdd, trades := runBacktest(candles, initBalance, feeRate, posSize, trailMult)
		ret := (balance - initBalance) / initBalance
		sharpe := 0.0
		if mdd > 0 {
			sharpe = (ret - 0.02) / mdd
		}
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
		fmt.Fprintf(w, "%s\t%+.2f%%\t%.2f%%\t%.2f\t%d\t%.0f%%\t%+.0f\n",
			p.label, ret*100, mdd*100, sharpe, len(trades), wr, balance-initBalance)
	}
	w.Flush()

	// Detailed 2026 monthly breakdown
	fmt.Printf("\n━━━ 2026 Monthly Detail ━━━\n")
	candles2026, _ := storage.LoadCandles(db, market, "60", "2026-01-01T00:00:00", "2026-12-31T23:59:59")
	if len(candles2026) >= 100 {
		runMonthly(candles2026, initBalance, feeRate, posSize, trailMult)
	}
}

func runMonthly(candles []indicator.CandleData, initBalance, feeRate, posSize, trailMult float64) {
	type position struct {
		price, volume, cost, stop float64
		entryTime                 string
	}
	var pos *position
	balance := initBalance
	peak := initBalance
	maxDD := 0.0
	monthlyPnL := make(map[string]float64)

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
				sig, r, _ := sigKeltner(candles, i)
				if sig == "sell" {
					exit = true
					reason = r
				}
			}
			if !exit {
				if atr > 0 {
					if trail := price - trailMult*atr; trail > pos.stop {
						pos.stop = trail
					}
				}
			}
			if exit {
				exitVal := pos.volume * price
				fee := exitVal * feeRate
				pnl := exitVal - pos.cost - fee
				balance += exitVal - fee
				month := candles[i].Time[:7]
				monthlyPnL[month] += pnl
				fmt.Printf("  SELL @ %12.0f (%s) PnL: %+8.0f | %s\n",
					price, candles[i].Time[:16], pnl, reason)
				pos = nil
			}
		}

		if pos == nil && i+1 < len(candles) {
			sig, _, stop := sigKeltner(candles, i)
			if sig == "buy" {
				size := (balance * posSize) / price
				cost := size * price
				entryFee := cost * feeRate
				if cost+entryFee > balance {
					continue
				}
				balance -= cost + entryFee
				atr := lastATR(candles, i)
				adjStop := stop
				if atr > 0 {
					adjStop = price - trailMult*atr
				}
				pos = &position{price, size, cost, adjStop, candles[i].Time}
				fmt.Printf("  BUY  @ %12.0f (%s) stop=%.0f\n",
					price, candles[i].Time[:16], adjStop)
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
		monthlyPnL[candles[len(candles)-1].Time[:7]] += pnl
		fmt.Printf("  SELL @ %12.0f (%s) PnL: %+8.0f | end of data\n",
			lastPrice, candles[len(candles)-1].Time[:16], pnl)
	}

	fmt.Printf("\n--- 2026 Monthly Returns ---\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Month\tPnL\tCumul%%\n")
	months := sortedKeys(monthlyPnL)
	running := initBalance
	for _, m := range months {
		running += monthlyPnL[m]
		cumul := (running - initBalance) / initBalance * 100
		fmt.Fprintf(w, "%s\t%+8.0f\t%+.2f%%\n", m, monthlyPnL[m], cumul)
	}
	w.Flush()

	fmt.Printf("\nFinal Balance: %.0f KRW (%+.2f%%)\n", balance, (balance-initBalance)/initBalance*100)
	fmt.Printf("Max Drawdown:  %.2f%%\n", maxDD*100)
}

func runBacktest(candles []indicator.CandleData, initBalance, feeRate, posSize, trailMult float64) (float64, float64, []optTrade) {
	type position struct {
		price, volume, cost, stop float64
		entryTime                 string
	}
	var pos *position
	balance := initBalance
	var trades []optTrade
	peak := initBalance
	maxDD := 0.0

	for i := 0; i < len(candles); i++ {
		price := candles[i].Close

		if pos != nil {
			exit := false
			atr := lastATR(candles, i)
			if pos.stop > 0 && price <= pos.stop {
				exit = true
			}
			if !exit {
				sig, _, _ := sigKeltner(candles, i)
				if sig == "sell" {
					exit = true
				}
			}
			if !exit {
				if atr > 0 {
					if trail := price - trailMult*atr; trail > pos.stop {
						pos.stop = trail
					}
				}
			}
			if exit {
				exitVal := pos.volume * price
				fee := exitVal * feeRate
				pnl := exitVal - pos.cost - fee
				balance += exitVal - fee
				trades = append(trades, optTrade{pnl})
				pos = nil
			}
		}

		if pos == nil && i+1 < len(candles) {
			sig, _, stop := sigKeltner(candles, i)
			if sig == "buy" {
				size := (balance * posSize) / price
				cost := size * price
				entryFee := cost * feeRate
				if cost+entryFee > balance {
					continue
				}
				balance -= cost + entryFee
				atr := lastATR(candles, i)
				adjStop := stop
				if atr > 0 {
					adjStop = price - trailMult*atr
				}
				pos = &position{price, size, cost, adjStop, candles[i].Time}
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

	if pos != nil {
		lastPrice := candles[len(candles)-1].Close
		exitVal := pos.volume * lastPrice
		fee := exitVal * feeRate
		pnl := exitVal - pos.cost - fee
		balance += exitVal - fee
		trades = append(trades, optTrade{pnl})
	}

	return balance, maxDD, trades
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
