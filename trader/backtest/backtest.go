package backtest

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/kyungw00k/upbit/trader"
	"github.com/kyungw00k/upbit/trader/config"
	"github.com/kyungw00k/upbit/trader/indicator"
	"github.com/kyungw00k/upbit/trader/types"
)

// BacktestResult holds the outcome of a backtest run.
type BacktestResult struct {
	TotalReturn float64
	SharpeRatio float64
	MaxDrawdown float64
	WinRate     float64
	TotalTrades int
	Trades      []BacktestTrade
}

// BacktestTrade records a single completed trade in the backtest.
type BacktestTrade struct {
	Market     string
	Side       string
	EntryPrice float64
	ExitPrice  float64
	PnL        float64
	EntryTime  string
	ExitTime   string
}

// AnnualizationFactor is the number of periods per year for Sharpe ratio
// calculation assuming daily candles.
const AnnualizationFactor = 252

// RiskFreeRate is the annual risk-free rate used for Sharpe ratio.
const RiskFreeRate = 0.02

// Engine runs backtests against historical candle data using an agent and risk gate.
type Engine struct {
	agent    trader.Agent
	riskGate *trader.RiskGate
}

// NewEngine creates a new backtest Engine.
func NewEngine(agent trader.Agent, riskGate *trader.RiskGate) *Engine {
	return &Engine{
		agent:    agent,
		riskGate: riskGate,
	}
}

// Run executes the backtest over the provided candle data.
// candles is a map of market -> sorted (ascending by time) candle data.
// initialBalance is the starting portfolio value in KRW.
func (e *Engine) Run(ctx context.Context, candles map[string][]indicator.CandleData, initialBalance float64) (*BacktestResult, error) {
	if initialBalance <= 0 {
		return nil, fmt.Errorf("backtest: initial balance must be positive")
	}
	if len(candles) == 0 {
		return nil, fmt.Errorf("backtest: no candle data provided")
	}

	// Use the first market for simplicity.
	var market string
	var data []indicator.CandleData
	for m, d := range candles {
		market = m
		data = d
		break
	}

	if len(data) < 50 {
		return nil, fmt.Errorf("backtest: insufficient candle data (%d candles, need at least 50)", len(data))
	}

	balance := initialBalance
	var position *openPosition
	var trades []BacktestTrade
	var equityCurve []float64
	peak := initialBalance
	maxDrawdown := 0.0

	for i := 50; i < len(data); i++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		currentPrice := data[i].Close
		currentTime := data[i].Time

		// Check exit conditions for open position.
		if position != nil {
			exitReason := ""
			// Stop loss hit.
			if currentPrice <= position.stopLoss {
				exitReason = "stop_loss"
			}
			// Take profit hit.
			if currentPrice >= position.takeProfit {
				exitReason = "take_profit"
			}

			if exitReason != "" {
				exitValue := position.volume * currentPrice
				fee := exitValue * trader.DefaultFeeRate
				pnl := exitValue - position.cost - fee
				balance += exitValue - fee

				trades = append(trades, BacktestTrade{
					Market:     market,
					Side:       "long",
					EntryPrice: position.entryPrice,
					ExitPrice:  currentPrice,
					PnL:        pnl,
					EntryTime:  position.entryTime,
					ExitTime:   currentTime,
				})
				position = nil
			}
		}

		// If no open position, evaluate for new entry.
		if position == nil {
			// Calculate indicators from available data up to current candle.
			window := data[:i+1]
			result := indicator.CalculateAll(window)

			indicators := map[string]any{
				"rsi":           result.RSI,
				"close":         currentPrice,
				"volume_ratio":  result.VolRatio,
				"macd":          map[string]any{"histogram": result.MACDHist},
				"bollinger_bands": map[string]any{
					"upper":  result.BBUpper,
					"middle": result.BBMiddle,
					"lower":  result.BBLower,
				},
				"ema": map[string]any{
					"20": result.EMA20,
					"50": result.EMA50,
				},
			}

			input := types.AgentInput{
				Markets:    []string{market},
				Indicators: indicators,
			}

			report, err := e.agent.Run(ctx, input)
			if err != nil {
				log.Printf("[backtest] agent error at %s: %v", currentTime, err)
				continue
			}

			// Log signals for debugging
			if i%25 == 0 || report.Signal != types.SignalHold {
				log.Printf("[backtest] i=%d %s price=%.0f signal=%s rsi=%.1f",
					i, currentTime[:10], currentPrice, report.Signal, result.RSI)
			}

			// Find the coin analysis for our market.
			var coinAnalysis *types.CoinAnalysis
			for _, c := range report.Coins {
				if c.Symbol == market {
					coinAnalysis = &types.CoinAnalysis{
						Symbol:     c.Symbol,
						Signal:     c.Signal,
						Confidence: c.Confidence,
						EntryPrice: currentPrice,
						StopLoss:   0,
						TakeProfit: 0,
						Position:   0,
						Reasoning:  c.Reasoning,
					}
					break
				}
			}

			if coinAnalysis == nil {
				continue
			}

			// Only consider BUY signals for long-only backtest.
			if coinAnalysis.Signal == types.SignalBuy || coinAnalysis.Signal == types.SignalStrongBuy {
				// Use ATR-based stop loss/take profit for dynamic risk management.
				if coinAnalysis.StopLoss == 0 {
					if result.ATR > 0 {
						coinAnalysis.StopLoss = currentPrice - result.ATR*1.5
					} else {
						coinAnalysis.StopLoss = currentPrice * 0.94
					}
				}
				if coinAnalysis.TakeProfit == 0 {
					if result.ATR > 0 {
						coinAnalysis.TakeProfit = currentPrice + result.ATR*3.0
					} else {
						coinAnalysis.TakeProfit = currentPrice * 1.15
					}
				}

				// Size position: use 20% of balance for meaningful exposure.
				positionSize := (balance * 0.10) / currentPrice
				coinAnalysis.Position = positionSize

				portfolio := trader.Portfolio{
					TotalValue: balance,
					Positions:  map[string]trader.Position{},
					DailyPnL:   0,
				}

				approved, adjustedSize, reason := e.riskGate.Evaluate(*coinAnalysis, portfolio)
					log.Printf("[backtest] risk gate: approved=%v size=%.4f reason=%s entry=%.0f sl=%.0f tp=%.0f pos=%.4f",
						approved, adjustedSize, reason, coinAnalysis.EntryPrice, coinAnalysis.StopLoss,
						coinAnalysis.TakeProfit, coinAnalysis.Position)
				if approved && adjustedSize > 0 {
					cost := adjustedSize * currentPrice
					entryFee := cost * trader.DefaultFeeRate
					balance -= cost + entryFee
					position = &openPosition{
						entryPrice: currentPrice,
						volume:     adjustedSize,
						cost:       cost,
						stopLoss:   coinAnalysis.StopLoss,
						takeProfit: coinAnalysis.TakeProfit,
						entryTime:  currentTime,
					}
				}
			}
		}

		// Track equity curve for Sharpe and drawdown.
		currentEquity := balance
		if position != nil {
			currentEquity += position.volume * currentPrice
		}
		equityCurve = append(equityCurve, currentEquity)

		if currentEquity > peak {
			peak = currentEquity
		}
		drawdown := (peak - currentEquity) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	// Close any remaining position at the last price.
	if position != nil {
		lastPrice := data[len(data)-1].Close
		lastTime := data[len(data)-1].Time
		exitValue := position.volume * lastPrice
		fee := exitValue * trader.DefaultFeeRate
		pnl := exitValue - position.cost - fee
		balance += exitValue - fee
		trades = append(trades, BacktestTrade{
			Market:     market,
			Side:       "long",
			EntryPrice: position.entryPrice,
			ExitPrice:  lastPrice,
			PnL:        pnl,
			EntryTime:  position.entryTime,
			ExitTime:   lastTime,
		})
	}

	// Calculate metrics.
	finalEquity := balance
	totalReturn := (finalEquity - initialBalance) / initialBalance

	sharpeRatio := calculateSharpeRatio(equityCurve)
	winRate := calculateWinRate(trades)

	return &BacktestResult{
		TotalReturn: totalReturn,
		SharpeRatio: sharpeRatio,
		MaxDrawdown: maxDrawdown,
		WinRate:     winRate,
		TotalTrades: len(trades),
		Trades:      trades,
	}, nil
}

// openPosition tracks a currently open position in the backtest.
type openPosition struct {
	entryPrice float64
	volume     float64
	cost       float64
	stopLoss   float64
	takeProfit float64
	entryTime  string
}

// calculateSharpeRatio computes the annualized Sharpe ratio from an equity curve.
func calculateSharpeRatio(equityCurve []float64) float64 {
	if len(equityCurve) < 2 {
		return 0
	}

	// Calculate daily returns.
	returns := make([]float64, len(equityCurve)-1)
	for i := 1; i < len(equityCurve); i++ {
		if equityCurve[i-1] > 0 {
			returns[i-1] = (equityCurve[i] - equityCurve[i-1]) / equityCurve[i-1]
		}
	}

	if len(returns) == 0 {
		return 0
	}

	// Mean return.
	var sum float64
	for _, r := range returns {
		sum += r
	}
	meanReturn := sum / float64(len(returns))

	// Standard deviation.
	var variance float64
	for _, r := range returns {
		diff := r - meanReturn
		variance += diff * diff
	}
	variance /= float64(len(returns))
	stdDev := math.Sqrt(variance)

	if stdDev == 0 {
		return 0
	}

	// Annualize.
	dailyRiskFree := RiskFreeRate / float64(AnnualizationFactor)
	sharpe := (meanReturn - dailyRiskFree) / stdDev * math.Sqrt(float64(AnnualizationFactor))

	return sharpe
}

// calculateWinRate returns the fraction of profitable trades.
func calculateWinRate(trades []BacktestTrade) float64 {
	if len(trades) == 0 {
		return 0
	}
	wins := 0
	for _, t := range trades {
		if t.PnL > 0 {
			wins++
		}
	}
	return float64(wins) / float64(len(trades))
}

// DefaultRiskConfig returns a default risk configuration suitable for backtesting.
func DefaultRiskConfig() config.RiskConfig {
	return config.RiskConfig{
		MaxPositionPct:  trader.DefaultMaxPositionPct,
		MaxDailyLossPct: trader.DefaultMaxDailyLossPct,
		StopLossPct:     trader.DefaultStopLossPct,
		MaxPositions:    trader.DefaultMaxPositions,
	}
}
