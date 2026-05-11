package trader

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/kyungw00k/upbit/trader/config"
	"github.com/kyungw00k/upbit/trader/types"
)

// Agent weights for signal merging.
const (
	WeightKeltnerAgent  = 0.4
	WeightStrategyAgent = 0.6
)

// Signal numeric values for weighted merging.
const (
	SignalValueStrongBuy  = 2.0
	SignalValueBuy        = 1.0
	SignalValueHold       = 0.0
	SignalValueSell       = -1.0
	SignalValueStrongSell = -2.0
)

// Pipeline orchestrates the full decision flow:
// agents -> signal merge -> risk gate -> execution.
type Pipeline struct {
	agents   []Agent
	riskGate *RiskGate
	executor Executor
	config   *config.Config
}

// NewPipeline creates a new Pipeline with the given components.
func NewPipeline(agents []Agent, riskGate *RiskGate, executor Executor, cfg *config.Config) *Pipeline {
	return &Pipeline{
		agents:   agents,
		riskGate: riskGate,
		executor: executor,
		config:   cfg,
	}
}

// Run executes the full pipeline: collect agent reports, merge signals,
// evaluate risk, and execute trades if approved.
func (p *Pipeline) Run(ctx context.Context, marketData map[string]any, sentiment *types.SentimentData) (*types.PipelineResult, error) {
	start := time.Now().UTC()

	if len(p.agents) == 0 {
		return &types.PipelineResult{
			RanAt:    start,
			Duration: time.Since(start),
		}, nil
	}

	// Step 1: Run each agent and collect AnalysisReports.
	input := types.AgentInput{
		Markets:    p.config.Trading.Markets,
		Indicators: marketData,
		Sentiment:  sentiment,
	}

	reports := make([]types.AnalysisReport, 0, len(p.agents))
	for _, agent := range p.agents {
		report, err := agent.Run(ctx, input)
		if err != nil {
			log.Printf("[pipeline] agent %s failed: %v", agent.ID(), err)
			continue
		}
		reports = append(reports, report)
	}

	if len(reports) == 0 {
		return &types.PipelineResult{
			RanAt:    start,
			Duration: time.Since(start),
		}, fmt.Errorf("all agents failed")
	}

	// Step 2: Merge signals using weighted voting.
	merged := p.mergeSignals(reports)

	// Step 3: If merged signal is HOLD, skip execution.
	if merged.Signal == types.SignalHold {
		log.Printf("[pipeline] merged signal is HOLD, skipping execution")
		return &types.PipelineResult{
			RanAt:    start,
			Duration: time.Since(start),
		}, nil
	}

	// Step 3.5: Extract execution params from strategy agent's response.
	execParams := p.extractExecutionParams(reports)
	log.Printf("[pipeline] execution params: posSize=%.0f%% trailMult=%.1fx stop=%.0f tp=%.0f conf=%.2f",
		execParams.PosSize*100, execParams.TrailMult, execParams.StopLoss, execParams.TakeProfit, execParams.Confidence)

	// Step 4: Evaluate risk for each coin in the merged report.
	for _, coin := range merged.Coins {
		if coin.Signal == types.SignalHold {
			continue
		}

		// Apply execution params from strategy agent.
		if coin.EntryPrice > 0 && execParams.StopLoss > 0 {
			coin.StopLoss = execParams.StopLoss
		}
		if coin.EntryPrice > 0 && execParams.TakeProfit > 0 {
			coin.TakeProfit = execParams.TakeProfit
		}
		if execParams.PosSize > 0 {
			coin.Position = execParams.PosSize
		}

		portfolio := Portfolio{
			TotalValue: p.config.Trading.InitialBalance,
			Positions:  map[string]Position{},
			DailyPnL:   0,
		}

		approved, adjustedSize, reason := p.riskGate.Evaluate(coin, portfolio)
		if !approved {
			log.Printf("[pipeline] risk gate rejected %s: %s", coin.Symbol, reason)
			continue
		}

		// Step 5: Execute the approved trade.
		side := "bid"
		if coin.Signal == types.SignalSell || coin.Signal == types.SignalStrongSell {
			side = "ask"
		}

		req := TradeRequest{
			Market: coin.Symbol,
			Side:   side,
			Price:  coin.EntryPrice,
			Volume: adjustedSize,
			Mode:   p.config.Trading.Mode,
		}

		result, err := p.executor.Execute(ctx, req)
		if err != nil {
			log.Printf("[pipeline] execution failed for %s: %v", coin.Symbol, err)
			continue
		}
		log.Printf("[pipeline] executed %s %s: uuid=%s filled=%.2f fee=%.2f",
			side, coin.Symbol, result.OrderUUID, result.FilledPrice, result.Fee)
	}

	return &types.PipelineResult{
		RanAt:    start,
		Duration: time.Since(start),
	}, nil
}

// extractExecutionParams parses StrategyResponse JSON from the strategy agent report
// and converts it to ExecutionParams. Falls back to defaults on failure.
func (p *Pipeline) extractExecutionParams(reports []types.AnalysisReport) types.ExecutionParams {
	for _, report := range reports {
		if report.AgentID != "strategy-agent" || len(report.Coins) == 0 {
			continue
		}

		var sr types.StrategyResponse
		if err := json.Unmarshal([]byte(report.Coins[0].Reasoning), &sr); err != nil {
			log.Printf("[pipeline] failed to parse strategy response: %v", err)
			continue
		}

		params := types.DefaultExecutionParams()
		if sr.PositionSize > 0 {
			params.PosSize = sr.PositionSize
		}
		if trailMult, err := parseATRMult(sr.ExitConditions.TrailingStop); err == nil {
			params.TrailMult = trailMult
		}
		params.Confidence = sr.Confidence
		return params.Clamp()
	}

	return types.DefaultExecutionParams()
}

// parseATRMult extracts the ATR multiplier from a string like "3.0×ATR" or "3.0xATR".
func parseATRMult(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty")
	}
	// Try direct float parse first.
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return v, nil
	}
	// Extract number before "×" or "x" or "ATR".
	for i, c := range s {
		if c == '×' || c == 'x' || c == 'X' {
			v, err := strconv.ParseFloat(s[:i], 64)
			if err == nil {
				return v, nil
			}
		}
	}
	return 0, fmt.Errorf("cannot parse: %s", s)
}

// mergeSignals combines multiple agent reports into a single merged report
// using weighted voting.
func (p *Pipeline) mergeSignals(reports []types.AnalysisReport) types.AnalysisReport {
	if len(reports) == 0 {
		return types.AnalysisReport{Signal: types.SignalHold}
	}
	if len(reports) == 1 {
		return reports[0]
	}

	// Build weights map by agent ID.
	weights := p.agentWeights()

	// Weighted average of overall signal and confidence.
	var totalWeight, weightedSignal, weightedConfidence float64
	var summary string

	for _, report := range reports {
		w, ok := weights[report.AgentID]
		if !ok {
			w = 0.5 // default weight for unknown agents
		}
		weightedSignal += signalToValue(report.Signal) * w
		weightedConfidence += report.Coins[0].Confidence * w
		totalWeight += w
		if summary == "" {
			summary = report.Summary
		}
	}

	if totalWeight == 0 {
		return types.AnalysisReport{Signal: types.SignalHold}
	}

	avgSignal := weightedSignal / totalWeight
	avgConfidence := weightedConfidence / totalWeight
	if avgConfidence > 1.0 {
		avgConfidence = 1.0
	}

	mergedSignal := valueToSignal(avgSignal)

	// Merge coins from all reports using the same weighted approach.
	mergedCoins := p.mergeCoins(reports, weights)

	return types.AnalysisReport{
		AgentID:   "pipeline-merged",
		Timestamp: time.Now().UTC(),
		Signal:    mergedSignal,
		Coins:     mergedCoins,
		Summary:   summary,
	}
}

// mergeCoins merges per-coin analyses from multiple agents.
func (p *Pipeline) mergeCoins(reports []types.AnalysisReport, weights map[string]float64) []types.CoinAnalysis {
	type coinAccumulator struct {
		weightedSignal     float64
		weightedConfidence float64
		totalWeight        float64
		entryPrice         float64
		stopLoss           float64
		takeProfit         float64
		position           float64
		reasoning          string
	}

	coinMap := make(map[string]*coinAccumulator)

	for _, report := range reports {
		w, ok := weights[report.AgentID]
		if !ok {
			w = 0.5
		}
		for _, coin := range report.Coins {
			acc, exists := coinMap[coin.Symbol]
			if !exists {
				acc = &coinAccumulator{}
				coinMap[coin.Symbol] = acc
			}
			acc.weightedSignal += signalToValue(coin.Signal) * w
			acc.weightedConfidence += coin.Confidence * w
			acc.totalWeight += w
			if coin.EntryPrice > 0 {
				acc.entryPrice = coin.EntryPrice
			}
			if coin.StopLoss > 0 {
				acc.stopLoss = coin.StopLoss
			}
			if coin.TakeProfit > 0 {
				acc.takeProfit = coin.TakeProfit
			}
			if coin.Position > 0 {
				acc.position = coin.Position
			}
			if coin.Reasoning != "" {
				acc.reasoning = coin.Reasoning
			}
		}
	}

	coins := make([]types.CoinAnalysis, 0, len(coinMap))
	for symbol, acc := range coinMap {
		if acc.totalWeight == 0 {
			continue
		}
		avgSignal := acc.weightedSignal / acc.totalWeight
		avgConfidence := acc.weightedConfidence / acc.totalWeight
		if avgConfidence > 1.0 {
			avgConfidence = 1.0
		}

		coins = append(coins, types.CoinAnalysis{
			Symbol:     symbol,
			Signal:     valueToSignal(avgSignal),
			Confidence: avgConfidence,
			EntryPrice: acc.entryPrice,
			StopLoss:   acc.stopLoss,
			TakeProfit: acc.takeProfit,
			Position:   acc.position,
			Reasoning:  acc.reasoning,
		})
	}

	return coins
}

// agentWeights returns a map of agent ID to its weight for signal merging.
func (p *Pipeline) agentWeights() map[string]float64 {
	return map[string]float64{
		"keltner-agent":  WeightKeltnerAgent,
		"strategy-agent": WeightStrategyAgent,
	}
}

// signalToValue converts a Signal to its numeric representation.
func signalToValue(s types.Signal) float64 {
	switch s {
	case types.SignalStrongBuy:
		return SignalValueStrongBuy
	case types.SignalBuy:
		return SignalValueBuy
	case types.SignalHold:
		return SignalValueHold
	case types.SignalSell:
		return SignalValueSell
	case types.SignalStrongSell:
		return SignalValueStrongSell
	default:
		return SignalValueHold
	}
}

// valueToSignal converts a numeric value back to a Signal.
func valueToSignal(v float64) types.Signal {
	switch {
	case v >= 1.5:
		return types.SignalStrongBuy
	case v >= 0.5:
		return types.SignalBuy
	case v <= -1.5:
		return types.SignalStrongSell
	case v <= -0.5:
		return types.SignalSell
	default:
		return types.SignalHold
	}
}
