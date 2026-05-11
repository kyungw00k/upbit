package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/kyungw00k/upbit/api"
	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/trader/agents"
	"github.com/kyungw00k/upbit/trader/indicator"
	tradertypes "github.com/kyungw00k/upbit/trader/types"
	"github.com/kyungw00k/upbit/types"
)

func formatKRW(v float64) string {
	if math.Abs(v) >= 1_000_000 {
		return fmt.Sprintf("%.1f만원", v/10_000)
	}
	return fmt.Sprintf("%.0f원", v)
}

type trade struct {
	entryIdx   int
	entryPrice float64
	entryTime  string
	exitIdx    int
	exitPrice  float64
	exitTime   string
	pnl        float64
	reason     string
}

func main() {
	market := envOr("MARKET", "KRW-BTC")
	count, _ := strconv.Atoi(envOr("COUNT", "0")) // 0 = unlimited
	interval := envOr("INTERVAL", "60")            // 60 = 1h candles
	from := envOr("FROM", "")                      // e.g. "2025-01-01T00:00:00"
	to := envOr("TO", "")                          // e.g. "2025-12-31T23:59:59"
	initBalance := 330_000.0

	fmt.Printf("=== RSI-2 Backtest: %s ===\n", market)
	fmt.Printf("Initial: %s | Interval: %sm | From: %s | To: %s\n\n",
		formatKRW(initBalance), interval, from, to)

	candles, err := fetchCandles(market, interval, from, to, count)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded %d candles (%s ~ %s)\n\n",
		len(candles), candles[0].Time, candles[len(candles)-1].Time)

	result := runBacktest(candles, market, initBalance)
	printResult(result, market, initBalance)
	printMonthlyReport(result, initBalance)
}

func runBacktest(candles []indicator.CandleData, market string, initBalance float64) []trade {
	agent := agents.NewRSI2Agent()
	balance := initBalance
	feeRate := 0.0005
	positionSize := 0.30 // 10% per trade

	var position *struct {
		entryPrice float64
		volume     float64
		cost       float64
		entryIdx   int
		entryTime  string
		stopPrice  float64 // trailing stop
	}
	var trades []trade
	peak := initBalance
	maxDD := 0.0

	// Need at least 52 candles for SMA50 + RSI(2)
	startIdx := 52

	for i := startIdx; i < len(candles); i++ {
		price := candles[i].Close

		// If position is open, check exit conditions
		if position != nil {
			// Update trailing stop (2 * ATR)
			window := candles[:i+1]
			atrValues := indicator.ATR(toCandleSlice(window), 14)
			if len(atrValues) > 0 && atrValues[len(atrValues)-1] > 0 {
				trailing := price - 2.0*atrValues[len(atrValues)-1]
				if trailing > position.stopPrice {
					position.stopPrice = trailing
				}
			}

			// Calculate RSI(2) for exit check
			closes := make([]float64, i+1)
			for j := 0; j <= i; j++ {
				closes[j] = candles[j].Close
			}
			rsi2Values := indicator.RSI(closes, 2)
			currentRSI2 := rsi2Values[len(rsi2Values)-1]

			exitReason := ""
			if currentRSI2 > 70 {
				exitReason = fmt.Sprintf("RSI2=%.1f overbought", currentRSI2)
			} else if price <= position.stopPrice {
				exitReason = fmt.Sprintf("trailing stop hit (stop=%.0f)", position.stopPrice)
			}

			if exitReason != "" {
				exitValue := position.volume * price
				fee := exitValue * feeRate
				pnl := exitValue - position.cost - fee
				balance += exitValue - fee

				trades = append(trades, trade{
					entryIdx:   position.entryIdx,
					entryPrice: position.entryPrice,
					entryTime:  position.entryTime,
					exitIdx:    i,
					exitPrice:  price,
					exitTime:   candles[i].Time,
					pnl:        pnl,
					reason:     exitReason,
				})

				fmt.Printf("  SELL @ %.0f (%s) PnL: %+.0f KRW | %s\n",
					price, candles[i].Time[:16], pnl, exitReason)

				position = nil
			}
		}

		// If no position, check entry conditions
		if position == nil && i+1 < len(candles) {
			window := candles[:i+1]
			signal, rsi2, reasoning := agent.EvaluateCandles(window)

			if signal == tradertypes.SignalBuy || signal == tradertypes.SignalStrongBuy {
				size := (balance * positionSize) / price
				cost := size * price
				entryFee := cost * feeRate

				if cost+entryFee > balance {
					continue // insufficient funds
				}

				balance -= cost + entryFee

				// Set initial stop loss at 2 * ATR below entry
				atrValues := indicator.ATR(toCandleSlice(window), 14)
				stopPrice := price * 0.94
				if len(atrValues) > 0 && atrValues[len(atrValues)-1] > 0 {
					stopPrice = price - 2.0*atrValues[len(atrValues)-1]
				}

				position = &struct {
					entryPrice float64
					volume     float64
					cost       float64
					entryIdx   int
					entryTime  string
					stopPrice  float64
				}{
					entryPrice: price,
					volume:     size,
					cost:       cost,
					entryIdx:   i,
					entryTime:  candles[i].Time,
					stopPrice:  stopPrice,
				}

				fmt.Printf("  BUY  @ %.0f (%s) RSI2=%.1f | %s\n",
					price, candles[i].Time[:16], rsi2, reasoning)
			}
		}

		// Track drawdown
		equity := balance
		if position != nil {
			equity += position.volume * price
		}
		if equity > peak {
			peak = equity
		}
		dd := (peak - equity) / peak
		if dd > maxDD {
			maxDD = dd
		}
	}

	// Close remaining position at last price
	if position != nil {
		lastPrice := candles[len(candles)-1].Close
		exitValue := position.volume * lastPrice
		fee := exitValue * feeRate
		pnl := exitValue - position.cost - fee
		balance += exitValue - fee
		trades = append(trades, trade{
			entryIdx:   position.entryIdx,
			entryPrice: position.entryPrice,
			entryTime:  position.entryTime,
			exitIdx:    len(candles) - 1,
			exitPrice:  lastPrice,
			exitTime:   candles[len(candles)-1].Time,
			pnl:        pnl,
			reason:     "end of data",
		})
		fmt.Printf("  SELL @ %.0f (%s) PnL: %+.0f KRW | end of data\n",
			lastPrice, candles[len(candles)-1].Time[:16], pnl)
	}

	// Append maxDD as a "metadata" trade (hack for reporting)
	// We'll calculate metrics inline instead

	fmt.Printf("\n--- Summary ---\n")
	totalReturn := (balance - initBalance) / initBalance
	wins := 0
	totalPnL := 0.0
	for _, t := range trades {
		totalPnL += t.pnl
		if t.pnl > 0 {
			wins++
		}
	}
	winRate := 0.0
	if len(trades) > 0 {
		winRate = float64(wins) / float64(len(trades))
	}

	fmt.Printf("Total Return: %.2f%% (%+.0f KRW)\n", totalReturn*100, balance-initBalance)
	fmt.Printf("Total Trades: %d\n", len(trades))
	fmt.Printf("Win Rate:     %.0f%% (%d/%d)\n", winRate*100, wins, len(trades))
	fmt.Printf("Max Drawdown: %.2f%%\n", maxDD*100)

	return trades
}

func printResult(trades []trade, market string, initBalance float64) {
	if len(trades) == 0 {
		return
	}

	fmt.Printf("\n--- Trade Log ---\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "#\tEntry\tExit\tPnL(KRW)\tReason\n")
	for i, t := range trades {
		fmt.Fprintf(w, "%d\t%.0f\t%.0f\t%+.0f\t%s\n",
			i+1, t.entryPrice, t.exitPrice, t.pnl, t.reason)
	}
	w.Flush()
}

func toCandleSlice(candles []indicator.CandleData) []indicator.CandleData {
	return candles
}

func fetchCandles(market, unit, from, to string, total int) ([]indicator.CandleData, error) {
	apiClient := api.NewClient("", "")
	qClient := quotation.NewQuotationClient(apiClient)

	interval := unit + "m"

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	raw, err := qClient.GetCandlesAll(ctx, market, interval, from, total)
	if err != nil {
		return nil, fmt.Errorf("GetCandlesAll failed: %w", err)
	}

	// Filter by `to` if specified
	if to != "" {
		filtered := make([]types.Candle, 0, len(raw))
		for _, r := range raw {
			if r.CandleDateTimeKst <= to {
				filtered = append(filtered, r)
			}
		}
		raw = filtered
	}

	result := make([]indicator.CandleData, len(raw))
	for i, r := range raw {
		result[i] = indicator.CandleData{
			Open:   r.OpeningPrice,
			High:   r.HighPrice,
			Low:    r.LowPrice,
			Close:  r.TradePrice,
			Volume: r.CandleAccTradeVolume,
			Time:   r.CandleDateTimeKst,
		}
	}

	return result, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type monthlyStats struct {
	month     string
	trades    int
	pnl       float64
	winCount  int
	startEquity float64
	endEquity  float64
}

func printMonthlyReport(trades []trade, initBalance float64) {
	if len(trades) == 0 {
		fmt.Println("\n거래 내역이 없습니다.")
		return
	}

	// Group trades by exit month
	monthMap := make(map[string]*monthlyStats)
	var monthOrder []string

	runningEquity := initBalance

	for _, t := range trades {
		// Parse exit time to get month (format: "2025-01-15T10:00:00")
		month := t.exitTime[:7] // "2025-01"
		if _, exists := monthMap[month]; !exists {
			monthMap[month] = &monthlyStats{
				month:       month,
				startEquity: runningEquity,
			}
			monthOrder = append(monthOrder, month)
		}
		ms := monthMap[month]
		ms.trades++
		ms.pnl += t.pnl
		if t.pnl > 0 {
			ms.winCount++
		}
	}

	fmt.Printf("\n=== 월별 수익률 리포트 (초기 자본: %s) ===\n", formatKRW(initBalance))
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "월\t거래수\t승률\t월PnL\t월수익률\t누적수익률\t잔고\n")

	cumulativeReturn := 0.0
	runningEquity = initBalance

	for _, month := range monthOrder {
		ms := monthMap[month]
		winRate := 0.0
		if ms.trades > 0 {
			winRate = float64(ms.winCount) / float64(ms.trades) * 100
		}
		monthlyReturn := ms.pnl / runningEquity * 100
		cumulativeReturn += monthlyReturn
		runningEquity += ms.pnl

		fmt.Fprintf(w, "%s\t%d\t%.0f%%\t%s\t%+.2f%%\t%+.2f%%\t%s\n",
			month, ms.trades, winRate, formatKRW(ms.pnl),
			monthlyReturn, cumulativeReturn, formatKRW(runningEquity))
	}
	w.Flush()

	totalReturn := (runningEquity - initBalance) / initBalance * 100
	fmt.Printf("\n최종 잔고: %s (총 수익률: %+.2f%%)\n", formatKRW(runningEquity), totalReturn)
}
