package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/kyungw00k/upbit/trader/types"
	openai "github.com/sashabaranov/go-openai"
)

// GLMAgent implements Agent using a GLM (OpenAI-compatible) API.
type GLMAgent struct {
	id     string
	client *openai.Client
	model  string
}

// NewGLMAgent creates a GLMAgent with the given API configuration.
func NewGLMAgent(apiKey, baseURL, model string) *GLMAgent {
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	return &GLMAgent{
		id:     "glm-analyst",
		client: openai.NewClientWithConfig(cfg),
		model:  model,
	}
}

// ID returns the agent identifier.
func (g *GLMAgent) ID() string {
	return g.id
}

// Run sends the market data to GLM and returns an AnalysisReport.
func (g *GLMAgent) Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	market := "KRW-BTC"
	if len(input.Markets) > 0 {
		market = input.Markets[0]
	}

	userPrompt := BuildAnalysisPrompt(market, input.Indicators, toSentimentInfo(input.Sentiment))

	req := openai.ChatCompletionRequest{
		Model: g.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: SystemPrompt()},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{Type: openai.ChatCompletionResponseFormatTypeJSONObject},
	}

	resp, err := g.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return types.AnalysisReport{}, fmt.Errorf("glm api call: %w", err)
	}

	if len(resp.Choices) == 0 {
		return types.AnalysisReport{}, fmt.Errorf("glm returned no choices")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	if content == "" {
		return types.AnalysisReport{}, fmt.Errorf("glm returned empty response")
	}

	// Strip markdown code fences if present (GLM sometimes wraps JSON in ```json ... ```).
	content = stripMarkdownFence(content)

	var result glmResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return types.AnalysisReport{}, fmt.Errorf("parsing glm response: %w\nraw: %s", err, content)
	}

	// Validate signal.
	signal := types.Signal(strings.ToUpper(result.Signal))
	if !signal.Valid() {
		signal = types.SignalHold
	}

	// Build coin analyses.
	coins := make([]types.CoinAnalysis, 0, len(result.Coins))
	for _, c := range result.Coins {
		cs := types.Signal(strings.ToUpper(c.Signal))
		if !cs.Valid() {
			cs = types.SignalHold
		}
		conf := clampConfidence(c.Confidence)
		coins = append(coins, types.CoinAnalysis{
			Symbol:     c.Symbol,
			Signal:     cs,
			Confidence: conf,
			Reasoning:  c.Reasoning,
		})
	}

	return types.AnalysisReport{
		AgentID:   g.id,
		Timestamp: time.Now().UTC(),
		Signal:    signal,
		Coins:     coins,
		Summary:   result.Summary,
	}, nil
}

// glmResponse mirrors the expected JSON structure from GLM.
type glmResponse struct {
	Signal     string        `json:"signal"`
	Confidence float64       `json:"confidence"`
	Reasoning  string        `json:"reasoning"`
	Coins      []glmCoin     `json:"coins"`
	Summary    string        `json:"summary"`
}

// glmCoin is a per-coin signal in the GLM response.
type glmCoin struct {
	Symbol     string  `json:"symbol"`
	Signal     string  `json:"signal"`
	Confidence float64 `json:"confidence"`
	Reasoning  string  `json:"reasoning"`
}

// stripMarkdownFence removes ```json ... ``` wrappers.
func stripMarkdownFence(s string) string {
	re := regexp.MustCompile("(?s)^```(?:json)?\\s*\\n?(.*?)\\n?```$")
	matches := re.FindStringSubmatch(s)
	if len(matches) == 2 {
		return strings.TrimSpace(matches[1])
	}
	return s
}

// clampConfidence ensures a value is between 0 and 1.
func clampConfidence(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// toSentimentInfo converts *types.SentimentData to SentimentInfo for prompt building.
func toSentimentInfo(sd *types.SentimentData) *SentimentInfo {
	if sd == nil {
		return nil
	}
	return &SentimentInfo{
		Index:    sd.FearGreedIndex,
		Category: sd.FearGreedCategory,
	}
}
