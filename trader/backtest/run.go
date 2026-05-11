package backtest

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/kyungw00k/upbit/trader"
	"github.com/kyungw00k/upbit/trader/agents"
	"github.com/kyungw00k/upbit/trader/indicator"
)

// RunBacktest executes a backtest with the given candle data and prints results.
// This is the main entry point for CLI-based backtesting.
func RunBacktest(ctx context.Context, market string, candles []indicator.CandleData, initialBalance float64) (*BacktestResult, error) {
	techAgent := agents.NewTechnicalAgent()
	riskGate := trader.NewRiskGate(DefaultRiskConfig())

	engine := NewEngine(techAgent, riskGate)

	result, err := engine.Run(ctx, map[string][]indicator.CandleData{
		market: candles,
	}, initialBalance)
	if err != nil {
		return nil, fmt.Errorf("backtest failed: %w", err)
	}

	return result, nil
}

// PrintResult formats and prints backtest results to stdout.
func PrintResult(market string, result *BacktestResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "\n=== Backtest Results: %s ===\n\n", market)
	fmt.Fprintf(w, "Total Return:\t%.2f%%\n", result.TotalReturn*100)
	fmt.Fprintf(w, "Sharpe Ratio:\t%.4f\n", result.SharpeRatio)
	fmt.Fprintf(w, "Max Drawdown:\t%.2f%%\n", result.MaxDrawdown*100)
	fmt.Fprintf(w, "Win Rate:\t%.2f%%\n", result.WinRate*100)
	fmt.Fprintf(w, "Total Trades:\t%d\n", result.TotalTrades)
	w.Flush()

	if len(result.Trades) > 0 {
		fmt.Fprintf(w, "\n--- Trade Log ---\n")
		fmt.Fprintf(w, "#\tSide\tEntry\tExit\tPnL\tEntry Time\tExit Time\n")
		for i, t := range result.Trades {
			fmt.Fprintf(w, "%d\t%s\t%.0f\t%.0f\t%.0f\t%s\t%s\n",
				i+1, t.Side, t.EntryPrice, t.ExitPrice, t.PnL, t.EntryTime, t.ExitTime)
		}
		w.Flush()
	}
	fmt.Println()
}

// RunBacktestFromCandles converts raw candle data and runs the backtest.
func RunBacktestFromCandles(ctx context.Context, market string, rawCandles []struct {
	Open, High, Low, Close, Volume float64
	Time                           string
}, initialBalance float64) (*BacktestResult, error) {
	if len(rawCandles) < 50 {
		return nil, fmt.Errorf("need at least 50 candles, got %d", len(rawCandles))
	}

	candles := make([]indicator.CandleData, len(rawCandles))
	for i, rc := range rawCandles {
		candles[i] = indicator.CandleData{
			Open:   rc.Open,
			High:   rc.High,
			Low:    rc.Low,
			Close:  rc.Close,
			Volume: rc.Volume,
			Time:   rc.Time,
		}
	}

	return RunBacktest(ctx, market, candles, initialBalance)
}

// init ensures the log package doesn't add timestamps when running backtest.
func init() {
	log.SetFlags(0)
}
