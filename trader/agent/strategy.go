package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kyungw00k/upbit/trader/types"
	openai "github.com/sashabaranov/go-openai"
)

// StrategyAgent uses an LLM to generate dynamic trading strategies every cycle.
// It produces an AnalysisReport with strategy details encoded in CoinAnalysis.Reasoning.
type StrategyAgent struct {
	id      string
	client  *openai.Client
	model   string
	timeout time.Duration
}

// NewStrategyAgent creates a StrategyAgent with OpenAI-compatible API config.
func NewStrategyAgent(apiKey, baseURL, model string, timeout time.Duration) *StrategyAgent {
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	if timeout == 0 {
		timeout = 15 * time.Second
	}
	return &StrategyAgent{
		id:      "strategy-agent",
		client:  openai.NewClientWithConfig(cfg),
		model:   model,
		timeout: timeout,
	}
}

// ID returns the agent identifier.
func (s *StrategyAgent) ID() string {
	return s.id
}

// Run calls the LLM to generate a trading strategy and converts it to an AnalysisReport.
func (s *StrategyAgent) Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	market := "KRW-BTC"
	if len(input.Markets) > 0 {
		market = input.Markets[0]
	}

	userPrompt := BuildStrategyPrompt(market, input.Indicators, toSentimentInfo(input.Sentiment))

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	req := openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: StrategySystemPrompt()},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{Type: openai.ChatCompletionResponseFormatTypeJSONObject},
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return s.holdReport(fmt.Sprintf("api error: %v", err)), nil
	}

	if len(resp.Choices) == 0 {
		return s.holdReport("empty response"), nil
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	if content == "" {
		return s.holdReport("empty content"), nil
	}

	content = stripMarkdownFence(content)

	var sr types.StrategyResponse
	if err := json.Unmarshal([]byte(content), &sr); err != nil {
		return s.holdReport(fmt.Sprintf("parse error: %v", err)), nil
	}

	return s.toReport(sr, market), nil
}

func (s *StrategyAgent) holdReport(reason string) types.AnalysisReport {
	log.Printf("[strategy-agent] falling back to HOLD: %s", reason)
	return types.AnalysisReport{
		AgentID:   s.id,
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalHold,
		Coins:     nil,
		Summary:   fmt.Sprintf("HOLD (fallback: %s)", reason),
	}
}

func (s *StrategyAgent) toReport(sr types.StrategyResponse, market string) types.AnalysisReport {
	signal := signalFromAction(sr.Action)

	// Serialize full strategy as JSON in the reasoning field
	// so the pipeline can extract execution params from it.
	strategyJSON, _ := json.Marshal(sr)

	coins := []types.CoinAnalysis{
		{
			Symbol:     market,
			Signal:     signal,
			Confidence: clampConfidence(sr.Confidence),
			Reasoning:  string(strategyJSON),
		},
	}

	return types.AnalysisReport{
		AgentID:   s.id,
		Timestamp: time.Now().UTC(),
		Signal:    signal,
		Coins:     coins,
		Summary:   fmt.Sprintf("[%s] %s — %s", sr.MarketRegime, sr.StrategyName, sr.Reasoning),
	}
}

func signalFromAction(action string) types.Signal {
	switch strings.ToUpper(action) {
	case "BUY":
		return types.SignalBuy
	case "SELL":
		return types.SignalSell
	case "HOLD":
		return types.SignalHold
	default:
		return types.SignalHold
	}
}
