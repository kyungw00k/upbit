# Plan 1: Foundation — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Project scaffolding, core types, agent interface, Technical Analyst agent, pipeline orchestrator, and CLI skeleton — producing an analysis-only daemon.

**Architecture:** Single Claude Haiku agent analyzes market data collected from Upbit API. Pipeline orchestrator manages the data collection → analysis → output flow. CLI provides start/stop/status commands.

**Tech Stack:** Go 1.24, Anthropic Go SDK, github.com/kyungw00k/upbit, Cobra CLI, Zap logging

**Design Spec:** `/Users/humphrey.park/Sandbox/upbit/docs/superpowers/specs/2026-04-03-autonomous-trading-daemon-design.md`

---

## File Map

| File | Responsibility |
|------|---------------|
| `cmd/trader/main.go` | CLI entry point, cobra commands |
| `trader/agent.go` | Agent interface + base types |
| `trader/pipeline.go` | Pipeline orchestrator |
| `trader/daemon.go` | Daemon main loop, scheduling |
| `trader/agents/technical.go` | Technical Analyst agent |
| `trader/upbit/client.go` | Upbit API wrapper (data collection) |
| `trader/config/config.go` | Configuration loading |
| `trader/config/config.yaml` | Default config template |
| `trader/types/types.go` | Trading-specific types |
| `configs/config.yaml` | User config file |
| `go.mod` | Module definition |
| `.gitignore` | Git ignore rules |

---

## Task 1: Project Scaffolding

- [ ] **1.1** Initialize Go module and create directory structure

**Working directory:** `/Users/humphrey.park/Sandbox/upbit-trader/`

```bash
mkdir -p /Users/humphrey.park/Sandbox/upbit-trader
cd /Users/humphrey.park/Sandbox/upbit-trader
git init
go mod init github.com/kyungw00k/upbit-trader

# Create directory structure
mkdir -p cmd/trader
mkdir -p trader/agents
mkdir -p trader/upbit
mkdir -p trader/config
mkdir -p trader/types
mkdir -p configs
mkdir -p internal/testutil

# Create placeholder files so directories are tracked by git
touch cmd/trader/main.go
touch trader/agent.go
touch trader/pipeline.go
touch trader/daemon.go
touch trader/agents/technical.go
touch trader/upbit/client.go
touch trader/config/config.go
touch trader/config/config.yaml
touch trader/types/types.go
touch configs/config.yaml
touch internal/testutil/mock.go
```

- [ ] **1.2** Install dependencies

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader

go get github.com/anthropics/anthropic-sdk-go
go get github.com/kyungw00k/upbit
go get github.com/spf13/cobra@latest
go get go.uber.org/zap
go get gopkg.in/yaml.v3

go mod tidy
```

- [ ] **1.3** Create `.gitignore`

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/.gitignore`

```gitignore
# Binaries
/bin/
/trader

# Data
/data/
*.db

# Config (user-specific secrets)
configs/config.yaml

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store

# PID files
*.pid

# Test
coverage.out
```

- [ ] **1.4** Verify build

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go build ./...
# Expected: no errors
go vet ./...
# Expected: no errors
```

- [ ] **1.5** Commit

```bash
git add -A
git commit -m "chore: project scaffolding with Go module and directory structure"
```

---

## Task 2: Configuration System

- [ ] **2.1** Write configuration types and loader

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/config/config.go`

```go
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the root configuration for the trading daemon.
type Config struct {
	API    APIConfig    `yaml:"api"`
	Daemon DaemonConfig `yaml:"daemon"`
	Log    LogConfig    `yaml:"log"`
}

// APIConfig holds API keys and model settings.
type APIConfig struct {
	AnthropicAPIKey string `yaml:"anthropic_api_key"`
	UpbitAccessKey  string `yaml:"upbit_access_key"`
	UpbitSecretKey  string `yaml:"upbit_secret_key"`
	ModelAnalyst    string `yaml:"model_analyst"`
}

// DaemonConfig controls daemon scheduling behavior.
type DaemonConfig struct {
	DecisionInterval string `yaml:"decision_interval"`
}

// LogConfig controls logging behavior.
type LogConfig struct {
	Level string `yaml:"level"`
}

// DecisionInterval returns the parsed decision interval duration.
func (c *DaemonConfig) DecisionInterval() (time.Duration, error) {
	if c.DecisionInterval == "" {
		return 4 * time.Hour, nil
	}
	return time.ParseDuration(c.DecisionInterval)
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			ModelAnalyst: "claude-3-5-haiku-20241022",
		},
		Daemon: DaemonConfig{
			DecisionInterval: "4h",
		},
		Log: LogConfig{
			Level: "info",
		},
	}
}

// Load reads a YAML config file and applies environment variable overrides.
// Environment variables take precedence over file values.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", path, err)
	}

	applyEnvOverrides(cfg)

	return cfg, nil
}

// applyEnvOverrides sets config values from environment variables if present.
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("ANTHROPIC_API_KEY"); v != "" {
		cfg.API.AnthropicAPIKey = v
	}
	if v := os.Getenv("UPBIT_ACCESS_KEY"); v != "" {
		cfg.API.UpbitAccessKey = v
	}
	if v := os.Getenv("UPBIT_SECRET_KEY"); v != "" {
		cfg.API.UpbitSecretKey = v
	}
	if v := os.Getenv("TRADER_MODEL_ANALYST"); v != "" {
		cfg.API.ModelAnalyst = v
	}
	if v := os.Getenv("TRADER_LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
	if v := os.Getenv("TRADER_DECISION_INTERVAL"); v != "" {
		cfg.Daemon.DecisionInterval = v
	}
}
```

- [ ] **2.2** Write default config template

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/config/config.yaml`

```yaml
# Default configuration template for upbit-trader
# Copy to configs/config.yaml and customize.
# Environment variables override file values (see config.go for list).

api:
  # Set via ANTHROPIC_API_KEY env var
  anthropic_api_key: ""
  # Set via UPBIT_ACCESS_KEY env var
  upbit_access_key: ""
  # Set via UPBIT_SECRET_KEY env var
  upbit_secret_key: ""
  # Claude model for analysis agents (default: claude-3-5-haiku-20241022)
  model_analyst: claude-3-5-haiku-20241022

daemon:
  # Decision interval: 4h, 1h, 30m, etc. (default: 4h)
  decision_interval: "4h"

log:
  # Log level: debug, info, warn, error (default: info)
  level: "info"
```

- [ ] **2.3** Write user config file

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/configs/config.yaml`

```yaml
api:
  anthropic_api_key: ""
  upbit_access_key: ""
  upbit_secret_key: ""
  model_analyst: claude-3-5-haiku-20241022

daemon:
  decision_interval: "4h"

log:
  level: "info"
```

- [ ] **2.4** Write tests

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/config/config_test.go`

```go
package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.API.ModelAnalyst != "claude-3-5-haiku-20241022" {
		t.Errorf("expected default model analyst, got %s", cfg.API.ModelAnalyst)
	}
	if cfg.Daemon.DecisionInterval != "4h" {
		t.Errorf("expected default interval 4h, got %s", cfg.Daemon.DecisionInterval)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("expected default log level info, got %s", cfg.Log.Level)
	}
}

func TestDecisionInterval(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Duration
	}{
		{"default empty", "", 4 * time.Hour},
		{"1 hour", "1h", 1 * time.Hour},
		{"30 minutes", "30m", 30 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := DaemonConfig{DecisionInterval: tt.input}
			got, err := dc.DecisionInterval()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := []byte(`
api:
  model_analyst: claude-sonnet-4-20250514
daemon:
  decision_interval: "1h"
log:
  level: "debug"
`)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.API.ModelAnalyst != "claude-sonnet-4-20250514" {
		t.Errorf("expected claude-sonnet-4-20250514, got %s", cfg.API.ModelAnalyst)
	}
	dur, _ := cfg.Daemon.DecisionInterval()
	if dur != 1*time.Hour {
		t.Errorf("expected 1h, got %v", dur)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected debug, got %s", cfg.Log.Level)
	}
}

func TestLoadNonexistentReturnsDefaults(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("Load should not error on missing file, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.API.ModelAnalyst != "claude-3-5-haiku-20241022" {
		t.Errorf("expected default model, got %s", cfg.API.ModelAnalyst)
	}
}

func TestEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := []byte(`
api:
  model_analyst: from-file
log:
  level: "info"
`)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("TRADER_MODEL_ANALYST", "from-env")
	t.Setenv("TRADER_LOG_LEVEL", "debug")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.API.ModelAnalyst != "from-env" {
		t.Errorf("expected env override 'from-env', got %s", cfg.API.ModelAnalyst)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected env override 'debug', got %s", cfg.Log.Level)
	}
}

func TestEnvOverridesForAPIKeys(t *testing.T) {
	cfg := DefaultConfig()

	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-test")
	t.Setenv("UPBIT_ACCESS_KEY", "access-test")
	t.Setenv("UPBIT_SECRET_KEY", "secret-test")

	applyEnvOverrides(cfg)

	if cfg.API.AnthropicAPIKey != "sk-ant-test" {
		t.Errorf("expected anthropic key from env")
	}
	if cfg.API.UpbitAccessKey != "access-test" {
		t.Errorf("expected upbit access key from env")
	}
	if cfg.API.UpbitSecretKey != "secret-test" {
		t.Errorf("expected upbit secret key from env")
	}
}
```

- [ ] **2.5** Run tests

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go test ./trader/config/ -v
# Expected: PASS — all 5 tests pass
```

- [ ] **2.6** Commit

```bash
git add trader/config/config.go trader/config/config.yaml trader/config/config_test.go configs/config.yaml
git commit -m "feat(config): YAML configuration with env var overrides"
```

---

## Task 3: Core Types

- [ ] **3.1** Write trading types

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/types/types.go`

```go
package types

import (
	"encoding/json"
	"time"

	"github.com/kyungw00k/upbit/types"
)

// Signal represents a trading signal strength.
type Signal string

const (
	SignalStrongBuy Signal = "STRONG_BUY"
	SignalBuy       Signal = "BUY"
	SignalHold      Signal = "HOLD"
	SignalSell      Signal = "SELL"
	SignalStrongSell Signal = "STRONG_SELL"
)

// Valid returns true if the signal is a recognized value.
func (s Signal) Valid() bool {
	switch s {
	case SignalStrongBuy, SignalBuy, SignalHold, SignalSell, SignalStrongSell:
		return true
	}
	return false
}

// CoinAnalysis holds per-coin analysis results from an agent.
type CoinAnalysis struct {
	Symbol     string  `json:"symbol"`
	Signal     Signal  `json:"signal"`
	Confidence float64 `json:"confidence"`
	EntryPrice float64 `json:"entry_price"`
	StopLoss   float64 `json:"stop_loss"`
	TakeProfit float64 `json:"take_profit"`
	Position   float64 `json:"position"`
	Reasoning  string  `json:"reasoning"`
}

// AnalysisReport is the output of an agent's analysis run.
type AnalysisReport struct {
	AgentID   string         `json:"agent_id"`
	Timestamp time.Time      `json:"timestamp"`
	Signal    Signal         `json:"signal"`
	Coins     []CoinAnalysis `json:"coins"`
	Summary   string         `json:"summary"`
}

// MarketSnapshot collects all market data needed for analysis.
type MarketSnapshot struct {
	Markets    []types.Market    `json:"markets"`
	Tickers    []types.Ticker    `json:"tickers"`
	Orderbooks []types.Orderbook `json:"orderbooks,omitempty"`
	Candles    map[string][]types.Candle `json:"candles,omitempty"`
	CapturedAt time.Time        `json:"captured_at"`
}

// TechnicalDigest holds per-coin technical indicator summaries.
type TechnicalDigest struct {
	Coins map[string]CoinTechnical `json:"coins"`
}

// CoinTechnical holds technical indicators for a single coin.
type CoinTechnical struct {
	Symbol         string  `json:"symbol"`
	RSI            float64 `json:"rsi,omitempty"`
	MACD           float64 `json:"macd,omitempty"`
	MACDSignal     float64 `json:"macd_signal,omitempty"`
	MACDHistogram  float64 `json:"macd_histogram,omitempty"`
	SMA20          float64 `json:"sma_20,omitempty"`
	SMA50          float64 `json:"sma_50,omitempty"`
	EMA12          float64 `json:"ema_12,omitempty"`
	EMA26          float64 `json:"ema_26,omitempty"`
	BollingerUpper float64 `json:"bollinger_upper,omitempty"`
	BollingerLower float64 `json:"bollinger_lower,omitempty"`
	VolumeRatio    float64 `json:"volume_ratio,omitempty"`
}

// AgentInput is the input passed to an agent's Run method.
type AgentInput struct {
	MarketData *MarketSnapshot  `json:"market_data"`
	TechData   *TechnicalDigest `json:"tech_data,omitempty"`
	Reports    []AnalysisReport `json:"reports,omitempty"`
}

// PipelineResult holds the combined results of a pipeline run.
type PipelineResult struct {
	Reports    []AnalysisReport `json:"reports"`
	RanAt      time.Time        `json:"ran_at"`
	Duration   time.Duration    `json:"duration_ms"`
	Error      string           `json:"error,omitempty"`
}

// MarshalJSON implements custom JSON for PipelineResult with duration in ms.
func (r *PipelineResult) MarshalJSON() ([]byte, error) {
	type Alias PipelineResult
	return json.Marshal(&struct {
		DurationMs int64 `json:"duration_ms"`
		*Alias
	}{
		DurationMs: r.Duration.Milliseconds(),
		Alias:      (*Alias)(r),
	})
}
```

- [ ] **3.2** Write tests

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/types/types_test.go`

```go
package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSignalValid(t *testing.T) {
	tests := []struct {
		signal Signal
		valid  bool
	}{
		{SignalStrongBuy, true},
		{SignalBuy, true},
		{SignalHold, true},
		{SignalSell, true},
		{SignalStrongSell, true},
		{Signal("INVALID"), false},
		{Signal(""), false},
	}

	for _, tt := range tests {
		if got := tt.signal.Valid(); got != tt.valid {
			t.Errorf("Signal(%q).Valid() = %v, want %v", tt.signal, got, tt.valid)
		}
	}
}

func TestCoinAnalysisMarshalRoundTrip(t *testing.T) {
	original := CoinAnalysis{
		Symbol:     "KRW-BTC",
		Signal:     SignalBuy,
		Confidence: 0.85,
		EntryPrice: 75000000.0,
		StopLoss:   72000000.0,
		TakeProfit: 82000000.0,
		Position:   500000.0,
		Reasoning:  "RSI oversold bounce with MACD crossover",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded CoinAnalysis
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Symbol != original.Symbol {
		t.Errorf("Symbol mismatch: got %s, want %s", decoded.Symbol, original.Symbol)
	}
	if decoded.Signal != original.Signal {
		t.Errorf("Signal mismatch: got %s, want %s", decoded.Signal, original.Signal)
	}
	if decoded.Confidence != original.Confidence {
		t.Errorf("Confidence mismatch: got %f, want %f", decoded.Confidence, original.Confidence)
	}
	if decoded.EntryPrice != original.EntryPrice {
		t.Errorf("EntryPrice mismatch: got %f, want %f", decoded.EntryPrice, original.EntryPrice)
	}
}

func TestAnalysisReportMarshalRoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	original := AnalysisReport{
		AgentID:   "technical-analyst",
		Timestamp: now,
		Signal:    SignalBuy,
		Coins: []CoinAnalysis{
			{
				Symbol:     "KRW-BTC",
				Signal:     SignalStrongBuy,
				Confidence: 0.9,
				EntryPrice: 75000000,
				Reasoning:  "Strong breakout above resistance",
			},
			{
				Symbol:     "KRW-ETH",
				Signal:     SignalHold,
				Confidence: 0.5,
				Reasoning:  "Consolidation phase",
			},
		},
		Summary: "BTC showing bullish momentum, ETH consolidating",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded AnalysisReport
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.AgentID != original.AgentID {
		t.Errorf("AgentID mismatch")
	}
	if decoded.Timestamp != original.Timestamp {
		t.Errorf("Timestamp mismatch")
	}
	if len(decoded.Coins) != 2 {
		t.Fatalf("Expected 2 coins, got %d", len(decoded.Coins))
	}
	if decoded.Coins[0].Symbol != "KRW-BTC" {
		t.Errorf("First coin symbol mismatch")
	}
	if decoded.Coins[1].Signal != SignalHold {
		t.Errorf("Second coin signal mismatch")
	}
}

func TestPipelineResultMarshalJSON(t *testing.T) {
	result := PipelineResult{
		RanAt:    time.Now().UTC(),
		Duration: 1500 * time.Millisecond,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if raw["duration_ms"].(float64) != 1500 {
		t.Errorf("expected duration_ms=1500, got %v", raw["duration_ms"])
	}
}

func TestAgentInputMarshalRoundTrip(t *testing.T) {
	original := AgentInput{
		MarketData: &MarketSnapshot{
			CapturedAt: time.Now().UTC(),
		},
		Reports: []AnalysisReport{
			{
				AgentID: "test",
				Signal:  SignalHold,
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded AgentInput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(decoded.Reports) != 1 {
		t.Errorf("expected 1 report, got %d", len(decoded.Reports))
	}
}

func TestTechnicalDigest(t *testing.T) {
	digest := TechnicalDigest{
		Coins: map[string]CoinTechnical{
			"KRW-BTC": {
				Symbol:        "KRW-BTC",
				RSI:           65.5,
				MACD:          120.3,
				MACDSignal:    98.7,
				MACDHistogram: 21.6,
				SMA20:         74000000,
				BollingerUpper: 78000000,
				BollingerLower: 70000000,
			},
		},
	}

	data, err := json.Marshal(digest)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded TechnicalDigest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	btc, ok := decoded.Coins["KRW-BTC"]
	if !ok {
		t.Fatal("KRW-BTC not found in decoded digest")
	}
	if btc.RSI != 65.5 {
		t.Errorf("RSI mismatch: got %f, want 65.5", btc.RSI)
	}
}
```

- [ ] **3.3** Run tests

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go test ./trader/types/ -v
# Expected: PASS — all 6 tests pass
```

- [ ] **3.4** Commit

```bash
git add trader/types/
git commit -m "feat(types): core trading types with JSON serialization"
```

---

## Task 4: Agent Interface
- [ ] **4.1** Write agent interface and base types

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/agent.go`

```go
package trader

import (
	"context"

	"github.com/kyungw00k/upbit-trader/trader/types"
)

// Agent is the interface that all trading agents must implement.
type Agent interface {
	// ID returns the unique identifier for this agent.
	ID() string

	// Run executes the agent with the given input and returns an analysis report.
	Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error)
}

// MockAgent is a test agent that returns a predefined report.
type MockAgent struct {
	id      string
	Report  types.AnalysisReport
	Err     error
	RunFunc func(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error)
}

// NewMockAgent creates a new mock agent for testing.
func NewMockAgent(id string) *MockAgent {
	return &MockAgent{id: id}
}

// ID returns the mock agent's identifier.
func (m *MockAgent) ID() string {
	return m.id
}

// Run returns the predefined report or calls the custom run function.
func (m *MockAgent) Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	if m.RunFunc != nil {
		return m.RunFunc(ctx, input)
	}
	if m.Err != nil {
		return types.AnalysisReport{}, m.Err
	}
	return m.Report, nil
}
```

- [ ] **4.2** Write tests

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/agent_test.go`

```go
package trader

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyungw00k/upbit-trader/trader/types"
)

func TestMockAgentID(t *testing.T) {
	agent := NewMockAgent("test-agent")
	if agent.ID() != "test-agent" {
		t.Errorf("expected ID 'test-agent', got %s", agent.ID())
	}
}

func TestMockAgentReturnsReport(t *testing.T) {
	expected := types.AnalysisReport{
		AgentID:   "test",
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalBuy,
		Summary:   "Test summary",
		Coins: []types.CoinAnalysis{
			{Symbol: "KRW-BTC", Signal: types.SignalStrongBuy, Confidence: 0.9},
		},
	}

	agent := NewMockAgent("test")
	agent.Report = expected

	report, err := agent.Run(context.Background(), types.AgentInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.AgentID != expected.AgentID {
		t.Errorf("AgentID mismatch")
	}
	if report.Signal != expected.Signal {
		t.Errorf("Signal mismatch")
	}
	if len(report.Coins) != 1 {
		t.Errorf("expected 1 coin, got %d", len(report.Coins))
	}
}

func TestMockAgentReturnsError(t *testing.T) {
	agent := NewMockAgent("failing")
	agent.Err = errors.New("agent failure")

	_, err := agent.Run(context.Background(), types.AgentInput{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "agent failure" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMockAgentCustomRunFunc(t *testing.T) {
	called := false
	agent := NewMockAgent("custom")
	agent.RunFunc = func(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
		called = true
		return types.AnalysisReport{AgentID: "custom", Signal: types.SignalHold}, nil
	}

	report, err := agent.Run(context.Background(), types.AgentInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("RunFunc was not called")
	}
	if report.Signal != types.SignalHold {
		t.Errorf("expected HOLD signal, got %s", report.Signal)
	}
}

// TestAgentInterfaceSatisfaction verifies MockAgent satisfies the Agent interface.
func TestAgentInterfaceSatisfaction(t *testing.T) {
	var _ Agent = (*MockAgent)(nil)
}
```

- [ ] **4.3** Run tests

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go test ./trader/ -v -run "TestMock|TestAgent"
# Expected: PASS — all 5 tests pass
```

- [ ] **4.4** Commit

```bash
git add trader/agent.go trader/agent_test.go
git commit -m "feat(agent): agent interface with mock agent for testing"
```

---

## Task 5: Upbit Client Wrapper

- [ ] **5.1** Write the Upbit data collector interface and implementation

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/upbit/client.go`

```go
package upbit

import (
	"context"
	"fmt"
	"sort"

	"github.com/kyungw00k/upbit/api"
	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit-trader/trader/types"
	uptypes "github.com/kyungw00k/upbit/types"
)

// Collector defines the interface for collecting market data from Upbit.
type Collector interface {
	// GetMarketSnapshot collects current market data for all KRW-paired markets.
	GetMarketSnapshot(ctx context.Context) (*types.MarketSnapshot, error)

	// GetCandles retrieves candle data for a specific market.
	GetCandles(ctx context.Context, market string, interval string, count int) ([]uptypes.Candle, error)
}

// Client wraps the Upbit API client for data collection.
type Client struct {
	quoteClient *quotation.QuotationClient
	topN        int // number of top coins by volume to include in candles
}

// NewClient creates a new Upbit data collector client.
func NewClient(accessKey, secretKey string) *Client {
	apiClient := api.NewClient(accessKey, secretKey)
	return &Client{
		quoteClient: quotation.NewQuotationClient(apiClient),
		topN:        20,
	}
}

// WithTopN sets the number of top coins by volume to fetch candles for.
func (c *Client) WithTopN(n int) *Client {
	c.topN = n
	return c
}

// GetMarketSnapshot collects tickers for all KRW markets, sorted by volume.
func (c *Client) GetMarketSnapshot(ctx context.Context) (*types.MarketSnapshot, error) {
	markets, err := c.quoteClient.GetMarkets(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting markets: %w", err)
	}

	// Filter to KRW-paired markets only
	var krwMarkets []uptypes.Market
	for _, m := range markets {
		if len(m.Market) > 4 && m.Market[:4] == "KRW-" {
			krwMarkets = append(krwMarkets, m)
		}
	}

	tickers, err := c.quoteClient.GetAllTickers(ctx, []string{"KRW"})
	if err != nil {
		return nil, fmt.Errorf("getting tickers: %w", err)
	}

	// Sort tickers by 24h trade price descending to find top coins
	sorted := make([]uptypes.Ticker, len(tickers))
	copy(sorted, tickers)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].AccTradePrice24h > sorted[j].AccTradePrice24h
	})

	// Fetch candles for top N coins
	candles := make(map[string][]uptypes.Candle)
	topMarkets := min(c.topN, len(sorted))
	for i := 0; i < topMarkets; i++ {
		candleData, err := c.quoteClient.GetCandles(ctx, sorted[i].Market, "days", 30)
		if err != nil {
			// Log warning but don't fail the entire snapshot
			continue
		}
		candles[sorted[i].Market] = candleData
	}

	return &types.MarketSnapshot{
		Markets:    krwMarkets,
		Tickers:    tickers,
		Candles:    candles,
		CapturedAt: uptypes.TimestampToTime(tickers[0].TradeTimestamp),
	}, nil
}

// GetCandles retrieves candle data for a specific market.
func (c *Client) GetCandles(ctx context.Context, market string, interval string, count int) ([]uptypes.Candle, error) {
	return c.quoteClient.GetCandles(ctx, market, interval, count)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

**NOTE:** The above uses `uptypes.TimestampToTime` which may not exist. If the upstream `types` package does not have this helper, use `time.Unix(timestamp/1000, 0)` instead. Check by running:

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
grep -r "TimestampToTime" $(go env GOMODCACHE)/github.com/kyungw00k/upbit@*/types/
```

If not found, replace the `CapturedAt` line in `GetMarketSnapshot` with:

```go
CapturedAt: time.Now().UTC(),
```

- [ ] **5.2** Write mock collector for testing

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/internal/testutil/mock_collector.go`

```go
package testutil

import (
	"context"
	"time"

	"github.com/kyungw00k/upbit-trader/trader/types"
	uptypes "github.com/kyungw00k/upbit/types"
)

// MockCollector is a mock implementation of the upbit.Collector interface.
type MockCollector struct {
	Snapshot *types.MarketSnapshot
	Err      error
}

// NewMockCollector creates a mock collector with sample data.
func NewMockCollector() *MockCollector {
	now := time.Now().UTC()
	return &MockCollector{
		Snapshot: &types.MarketSnapshot{
			Markets: []uptypes.Market{
				{Market: "KRW-BTC", KoreanName: "비트코인", EnglishName: "Bitcoin"},
				{Market: "KRW-ETH", KoreanName: "이더리움", EnglishName: "Ethereum"},
				{Market: "KRW-SOL", KoreanName: "솔라나", EnglishName: "Solana"},
			},
			Tickers: []uptypes.Ticker{
				{Market: "KRW-BTC", TradePrice: 75000000, ChangeRate: 0.025, AccTradePrice24h: 1500000000000, OpeningPrice: 73000000, HighPrice: 76000000, LowPrice: 72000000, TradeTimestamp: now.UnixMilli()},
				{Market: "KRW-ETH", TradePrice: 5200000, ChangeRate: -0.015, AccTradePrice24h: 800000000000, OpeningPrice: 5300000, HighPrice: 5350000, LowPrice: 5100000, TradeTimestamp: now.UnixMilli()},
				{Market: "KRW-SOL", TradePrice: 180000, ChangeRate: 0.05, AccTradePrice24h: 500000000000, OpeningPrice: 170000, HighPrice: 185000, LowPrice: 168000, TradeTimestamp: now.UnixMilli()},
			},
			CapturedAt: now,
		},
	}
}

// GetMarketSnapshot returns the mock snapshot.
func (m *MockCollector) GetMarketSnapshot(ctx context.Context) (*types.MarketSnapshot, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Snapshot, nil
}

// GetCandles returns empty candle data for the mock.
func (m *MockCollector) GetCandles(ctx context.Context, market string, interval string, count int) ([]uptypes.Candle, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return []uptypes.Candle{}, nil
}
```

- [ ] **5.3** Write tests

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/internal/testutil/mock_collector_test.go`

```go
package testutil

import (
	"context"
	"errors"
	"testing"
)

func TestMockCollectorReturnsSnapshot(t *testing.T) {
	collector := NewMockCollector()

	snapshot, err := collector.GetMarketSnapshot(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(snapshot.Markets) != 3 {
		t.Errorf("expected 3 markets, got %d", len(snapshot.Markets))
	}
	if len(snapshot.Tickers) != 3 {
		t.Errorf("expected 3 tickers, got %d", len(snapshot.Tickers))
	}
	if snapshot.Tickers[0].Market != "KRW-BTC" {
		t.Errorf("expected first ticker KRW-BTC, got %s", snapshot.Tickers[0].Market)
	}
	if snapshot.Tickers[0].TradePrice != 75000000 {
		t.Errorf("expected BTC price 75000000, got %f", snapshot.Tickers[0].TradePrice)
	}
}

func TestMockCollectorReturnsError(t *testing.T) {
	collector := NewMockCollector()
	collector.Err = errors.New("api error")

	_, err := collector.GetMarketSnapshot(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMockCollectorGetCandles(t *testing.T) {
	collector := NewMockCollector()

	candles, err := collector.GetCandles(context.Background(), "KRW-BTC", "days", 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(candles) != 0 {
		t.Errorf("expected empty candles, got %d", len(candles))
	}
}
```

- [ ] **5.4** Run tests

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go test ./internal/testutil/ -v
# Expected: PASS — all 3 tests pass
```

- [ ] **5.5** Commit

```bash
git add trader/upbit/ internal/testutil/mock_collector.go internal/testutil/mock_collector_test.go
git commit -m "feat(upbit): data collector wrapper with mock for testing"
```

---

## Task 6: Technical Analyst Agent

- [ ] **6.1** Write the Technical Analyst agent

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/agents/technical.go`

```go
package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/kyungw00k/upbit-trader/trader/types"
	"go.uber.org/zap"
)

const (
	// TechnicalAnalystID is the unique identifier for the Technical Analyst agent.
	TechnicalAnalystID = "technical-analyst"

	// technicalAnalystSystemPrompt is the system prompt for Claude.
	technicalAnalystSystemPrompt = `You are an expert cryptocurrency technical analyst. Analyze the provided market data and return a structured JSON response.

Your analysis should cover:
1. Overall market sentiment based on price action and volume
2. Per-coin technical analysis with specific entry/stop/target levels
3. Confidence scores based on signal strength

You MUST respond with valid JSON only (no markdown, no explanation outside JSON) in this exact format:
{
  "signal": "HOLD",
  "summary": "Brief overall market assessment (1-2 sentences)",
  "coins": [
    {
      "symbol": "KRW-BTC",
      "signal": "STRONG_BUY",
      "confidence": 0.85,
      "entry_price": 75000000,
      "stop_loss": 72000000,
      "take_profit": 82000000,
      "position": 500000,
      "reasoning": "RSI at 35 showing oversold, MACD bullish crossover, volume increasing"
    }
  ]
}

Signal values: STRONG_BUY, BUY, HOLD, SELL, STRONG_SELL
Confidence: 0.0 to 1.0
Position: suggested position size in KRW (0 means no position)
Entry/stop/target prices should be in the coin's quote currency (KRW)

Only include coins with a BUY or STRONG_BUY signal in the coins array. Exclude HOLD/SELL coins from the array but mention them in the summary if notable.`
)

// TechnicalAnalyst is the Technical Analyst agent that uses Claude to analyze market data.
type TechnicalAnalyst struct {
	client *anthropic.Client
	model  string
	logger *zap.Logger
}

// NewTechnicalAnalyst creates a new Technical Analyst agent.
func NewTechnicalAnalyst(apiKey string, model string, logger *zap.Logger) *TechnicalAnalyst {
	var opts []anthropic.ClientOption
	if apiKey != "" {
		opts = append(opts, anthropic.WithAPIKey(apiKey))
	}

	client := anthropic.NewClient(opts...)

	if model == "" {
		model = string(anthropic.ModelClaude3_5Haiku20241022)
	}

	return &TechnicalAnalyst{
		client: client,
		model:  model,
		logger: logger.Named("technical-analyst"),
	}
}

// ID returns the agent's unique identifier.
func (ta *TechnicalAnalyst) ID() string {
	return TechnicalAnalystID
}

// Run executes the Technical Analyst against the provided market data.
func (ta *TechnicalAnalyst) Run(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
	if input.MarketData == nil {
		return types.AnalysisReport{}, fmt.Errorf("market data is required")
	}

	ta.logger.Info("starting technical analysis",
		zap.Int("tickers", len(input.MarketData.Tickers)),
		zap.Int("markets", len(input.MarketData.Markets)),
	)

	// Serialize market data for the Claude prompt
	marketDataJSON, err := json.MarshalIndent(input.MarketData, "", "  ")
	if err != nil {
		return types.AnalysisReport{}, fmt.Errorf("serializing market data: %w", err)
	}

	ta.logger.Debug("sending market data to Claude",
		zap.Int("json_size", len(marketDataJSON)),
	)

	message, err := ta.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(ta.model),
		MaxTokens: 4096,
		System: []anthropic.TextBlockParam{
			{Text: technicalAnalystSystemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(string(marketDataJSON))),
		},
	})
	if err != nil {
		return types.AnalysisReport{}, fmt.Errorf("Claude API call failed: %w", err)
	}

	if len(message.Content) == 0 {
		return types.AnalysisReport{}, fmt.Errorf("empty response from Claude")
	}

	// Extract text content
	var responseText string
	for _, block := range message.Content {
		if block.Type == anthropic.ContentBlockTypeText {
			responseText = block.Text
			break
		}
	}

	if responseText == "" {
		return types.AnalysisReport{}, fmt.Errorf("no text content in Claude response")
	}

	ta.logger.Debug("received Claude response",
		zap.Int("response_size", len(responseText)),
	)

	// Parse the JSON response
	report, err := ta.parseResponse(responseText)
	if err != nil {
		return types.AnalysisReport{}, fmt.Errorf("parsing Claude response: %w", err)
	}

	ta.logger.Info("technical analysis complete",
		zap.String("signal", string(report.Signal)),
		zap.Int("coins_analyzed", len(report.Coins)),
		zap.String("summary", report.Summary),
	)

	return report, nil
}

// parseResponse parses the Claude JSON response into an AnalysisReport.
func (ta *TechnicalAnalyst) parseResponse(responseText string) (types.AnalysisReport, error) {
	var raw struct {
		Signal  string              `json:"signal"`
		Summary string              `json:"summary"`
		Coins   []types.CoinAnalysis `json:"coins"`
	}

	if err := json.Unmarshal([]byte(responseText), &raw); err != nil {
		ta.logger.Error("failed to parse Claude response as JSON",
			zap.String("response", responseText),
			zap.Error(err),
		)
		return types.AnalysisReport{}, fmt.Errorf("JSON parse error: %w", err)
	}

	signal := types.Signal(raw.Signal)
	if !signal.Valid() {
		return types.AnalysisReport{}, fmt.Errorf("invalid signal: %s", raw.Signal)
	}

	// Validate each coin analysis
	for i, coin := range raw.Coins {
		if !coin.Signal.Valid() {
			return types.AnalysisReport{}, fmt.Errorf("coin[%d] (%s): invalid signal: %s", i, coin.Symbol, coin.Signal)
		}
		if coin.Confidence < 0 || coin.Confidence > 1 {
			return types.AnalysisReport{}, fmt.Errorf("coin[%d] (%s): confidence must be 0-1, got %f", i, coin.Symbol, coin.Confidence)
		}
	}

	return types.AnalysisReport{
		AgentID:   TechnicalAnalystID,
		Timestamp: time.Now().UTC(),
		Signal:    signal,
		Coins:     raw.Coins,
		Summary:   raw.Summary,
	}, nil
}
```

- [ ] **6.2** Write tests

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/agents/technical_test.go`

```go
package agents

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kyungw00k/upbit-trader/trader/types"
	"go.uber.org/zap"
)

func TestTechnicalAnalystID(t *testing.T) {
	ta := NewTechnicalAnalyst("", "claude-3-5-haiku-20241022", zap.NewNop())
	if ta.ID() != TechnicalAnalystID {
		t.Errorf("expected ID %s, got %s", TechnicalAnalystID, ta.ID())
	}
}

func TestTechnicalAnalystParseResponse(t *testing.T) {
	ta := NewTechnicalAnalyst("", "claude-3-5-haiku-20241022", zap.NewNop())

	validResponse := `{
		"signal": "BUY",
		"summary": "Market showing bullish momentum with BTC leading",
		"coins": [
			{
				"symbol": "KRW-BTC",
				"signal": "STRONG_BUY",
				"confidence": 0.85,
				"entry_price": 75000000,
				"stop_loss": 72000000,
				"take_profit": 82000000,
				"position": 500000,
				"reasoning": "RSI oversold bounce"
			}
		]
	}`

	report, err := ta.parseResponse(validResponse)
	if err != nil {
		t.Fatalf("parseResponse failed: %v", err)
	}

	if report.AgentID != TechnicalAnalystID {
		t.Errorf("expected agent ID %s, got %s", TechnicalAnalystID, report.AgentID)
	}
	if report.Signal != types.SignalBuy {
		t.Errorf("expected BUY signal, got %s", report.Signal)
	}
	if report.Summary != "Market showing bullish momentum with BTC leading" {
		t.Errorf("unexpected summary: %s", report.Summary)
	}
	if len(report.Coins) != 1 {
		t.Fatalf("expected 1 coin, got %d", len(report.Coins))
	}
	if report.Coins[0].Symbol != "KRW-BTC" {
		t.Errorf("expected KRW-BTC, got %s", report.Coins[0].Symbol)
	}
	if report.Coins[0].Confidence != 0.85 {
		t.Errorf("expected confidence 0.85, got %f", report.Coins[0].Confidence)
	}
	if report.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestTechnicalAnalystParseResponseInvalidSignal(t *testing.T) {
	ta := NewTechnicalAnalyst("", "claude-3-5-haiku-20241022", zap.NewNop())

	invalidResponse := `{"signal": "INVALID", "summary": "test", "coins": []}`

	_, err := ta.parseResponse(invalidResponse)
	if err == nil {
		t.Fatal("expected error for invalid signal")
	}
}

func TestTechnicalAnalystParseResponseInvalidCoinSignal(t *testing.T) {
	ta := NewTechnicalAnalyst("", "claude-3-5-haiku-20241022", zap.NewNop())

	invalidCoinResponse := `{
		"signal": "BUY",
		"summary": "test",
		"coins": [{"symbol": "KRW-BTC", "signal": "MAYBE", "confidence": 0.5}]
	}`

	_, err := ta.parseResponse(invalidCoinResponse)
	if err == nil {
		t.Fatal("expected error for invalid coin signal")
	}
}

func TestTechnicalAnalystParseResponseInvalidConfidence(t *testing.T) {
	ta := NewTechnicalAnalyst("", "claude-3-5-haiku-20241022", zap.NewNop())

	invalidConfidence := `{
		"signal": "HOLD",
		"summary": "test",
		"coins": [{"symbol": "KRW-BTC", "signal": "BUY", "confidence": 1.5}]
	}`

	_, err := ta.parseResponse(invalidConfidence)
	if err == nil {
		t.Fatal("expected error for invalid confidence")
	}
}

func TestTechnicalAnalystParseResponseInvalidJSON(t *testing.T) {
	ta := NewTechnicalAnalyst("", "claude-3-5-haiku-20241022", zap.NewNop())

	_, err := ta.parseResponse("not json at all")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestTechnicalAnalystRunWithMockServer(t *testing.T) {
	// Create a mock Anthropic API server
	claudeResponse := types.AnalysisReport{
		Signal:  types.SignalBuy,
		Summary: "BTC bullish, ETH neutral",
		Coins: []types.CoinAnalysis{
			{Symbol: "KRW-BTC", Signal: types.SignalStrongBuy, Confidence: 0.9, EntryPrice: 75000000, Reasoning: "Breakout"},
		},
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"id":    "msg_test",
		"type":  "message",
		"role":  "assistant",
		"model": "claude-3-5-haiku-20241022",
		"content": []map[string]interface{}{
			{"type": "text", "text": func() string {
				b, _ := json.Marshal(map[string]interface{}{
					"signal":  "BUY",
					"summary": "BTC bullish, ETH neutral",
					"coins":   claudeResponse.Coins,
				})
				return string(b)
			}()},
		},
		"stop_reason": "end_turn",
		"usage":       map[string]int{"input_tokens": 100, "output_tokens": 50},
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
	}))
	defer server.Close()

	ta := NewTechnicalAnalystWithURL(server.URL, zap.NewNop())

	input := types.AgentInput{
		MarketData: &types.MarketSnapshot{
			Tickers: []interface{}{
				map[string]interface{}{"market": "KRW-BTC", "trade_price": 75000000.0},
			},
			CapturedAt: time.Now().UTC(),
		},
	}

	report, err := ta.Run(context.Background(), input)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if report.Signal != types.SignalBuy {
		t.Errorf("expected BUY, got %s", report.Signal)
	}
	if len(report.Coins) != 1 {
		t.Fatalf("expected 1 coin, got %d", len(report.Coins))
	}
}

func TestTechnicalAnalystRunNilMarketData(t *testing.T) {
	ta := NewTechnicalAnalyst("", "claude-3-5-haiku-20241022", zap.NewNop())

	_, err := ta.Run(context.Background(), types.AgentInput{})
	if err == nil {
		t.Fatal("expected error for nil market data")
	}
}
```

**NOTE:** The `TestTechnicalAnalystRunWithMockServer` test uses `NewTechnicalAnalystWithURL` which needs to be added to `technical.go`. Also the test references `interface{}` types which won't match `types.Ticker`. Add this helper to `technical.go`:

```go
// NewTechnicalAnalystWithURL creates a Technical Analyst pointed at a custom URL (for testing).
func NewTechnicalAnalystWithURL(baseURL string, logger *zap.Logger) *TechnicalAnalyst {
	client := anthropic.NewClient(
		anthropic.WithBaseURL(baseURL),
		anthropic.WithAPIKey("test-key"),
	)
	return &TechnicalAnalyst{
		client: client,
		model:  string(anthropic.ModelClaude3_5Haiku20241022),
		logger: logger.Named("technical-analyst-test"),
	}
}
```

Also, update the `TestTechnicalAnalystRunWithMockServer` test to use proper types:

```go
// Replace the MarketData in the test with proper upbit types:
import uptypes "github.com/kyungw00k/upbit/types"

input := types.AgentInput{
	MarketData: &types.MarketSnapshot{
		Tickers: []uptypes.Ticker{
			{Market: "KRW-BTC", TradePrice: 75000000, ChangeRate: 0.025, AccTradePrice24h: 1.5e12},
		},
		CapturedAt: time.Now().UTC(),
	},
}
```

- [ ] **6.3** Run tests

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go test ./trader/agents/ -v
# Expected: PASS — all 7 tests pass
```

- [ ] **6.4** Commit

```bash
git add trader/agents/
git commit -m "feat(agents): Technical Analyst agent with Claude Haiku integration"
```

---

## Task 7: Pipeline Orchestrator

- [ ] **7.1** Write the pipeline orchestrator

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/pipeline.go`

```go
package trader

import (
	"context"
	"fmt"
	"time"

	"github.com/kyungw00k/upbit-trader/trader/types"
	"go.uber.org/zap"
)

// Collector defines the interface for collecting market data.
type Collector interface {
	GetMarketSnapshot(ctx context.Context) (*types.MarketSnapshot, error)
}

// Pipeline orchestrates the data collection → analysis → output flow.
type Pipeline struct {
	collector Collector
	agents    []Agent
	logger    *zap.Logger
}

// NewPipeline creates a new pipeline with the given collector and agents.
func NewPipeline(collector Collector, agents []Agent, logger *zap.Logger) *Pipeline {
	return &Pipeline{
		collector: collector,
		agents:    agents,
		logger:    logger.Named("pipeline"),
	}
}

// Run executes the full pipeline once and returns the combined results.
func (p *Pipeline) Run(ctx context.Context) (*types.PipelineResult, error) {
	start := time.Now()
	p.logger.Info("pipeline started")

	result := &types.PipelineResult{
		RanAt: start.UTC(),
	}

	// Phase 1: Collect market data
	p.logger.Info("collecting market data")
	snapshot, err := p.collector.GetMarketSnapshot(ctx)
	if err != nil {
		result.Error = fmt.Sprintf("data collection failed: %v", err)
		p.logger.Error("data collection failed", zap.Error(err))
		return result, fmt.Errorf("data collection: %w", err)
	}

	p.logger.Info("market data collected",
		zap.Int("tickers", len(snapshot.Tickers)),
		zap.Int("markets", len(snapshot.Markets)),
	)

	// Phase 2: Run each agent with the collected data
	input := types.AgentInput{
		MarketData: snapshot,
	}

	for _, agent := range p.agents {
		p.logger.Info("running agent", zap.String("agent_id", agent.ID()))

		report, err := agent.Run(ctx, input)
		if err != nil {
			p.logger.Error("agent failed",
				zap.String("agent_id", agent.ID()),
				zap.Error(err),
			)
			// Continue with other agents even if one fails
			result.Reports = append(result.Reports, types.AnalysisReport{
				AgentID:   agent.ID(),
				Timestamp: time.Now().UTC(),
				Error:     err.Error(), // Store error info in report
			})
			continue
		}

		p.logger.Info("agent completed",
			zap.String("agent_id", agent.ID()),
			zap.String("signal", string(report.Signal)),
			zap.Int("coins", len(report.Coins)),
		)

		// Pass this report to subsequent agents
		input.Reports = append(input.Reports, report)
		result.Reports = append(result.Reports, report)
	}

	result.Duration = time.Since(start)
	p.logger.Info("pipeline completed",
		zap.Duration("duration", result.Duration),
		zap.Int("reports", len(result.Reports)),
	)

	return result, nil
}
```

**NOTE:** The `AnalysisReport` struct needs an `Error` field for storing agent failure info. Add to `trader/types/types.go`:

```go
// Add to AnalysisReport struct:
Error string `json:"error,omitempty"`
```

- [ ] **7.2** Write tests

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/pipeline_test.go`

```go
package trader

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyungw00k/upbit-trader/trader/types"
	testutil "github.com/kyungw00k/upbit-trader/internal/testutil"
	"go.uber.org/zap"
)

func TestPipelineRunSingleAgent(t *testing.T) {
	logger := zap.NewNop()
	collector := testutil.NewMockCollector()

	agent := NewMockAgent("test-agent")
	agent.Report = types.AnalysisReport{
		AgentID:   "test-agent",
		Timestamp: time.Now().UTC(),
		Signal:    types.SignalBuy,
		Summary:   "Test analysis",
		Coins: []types.CoinAnalysis{
			{Symbol: "KRW-BTC", Signal: types.SignalStrongBuy, Confidence: 0.9},
		},
	}

	pipeline := NewPipeline(collector, []Agent{agent}, logger)
	result, err := pipeline.Run(context.Background())

	if err != nil {
		t.Fatalf("pipeline Run failed: %v", err)
	}
	if len(result.Reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(result.Reports))
	}
	if result.Reports[0].Signal != types.SignalBuy {
		t.Errorf("expected BUY signal, got %s", result.Reports[0].Signal)
	}
	if result.Duration <= 0 {
		t.Error("expected positive duration")
	}
	if !result.RanAt.Before(time.Now().UTC()) {
		t.Error("expected RanAt to be in the past")
	}
}

func TestPipelineRunMultipleAgents(t *testing.T) {
	logger := zap.NewNop()
	collector := testutil.NewMockCollector()

	agent1 := NewMockAgent("agent-1")
	agent1.Report = types.AnalysisReport{
		AgentID: "agent-1",
		Signal:  types.SignalBuy,
		Summary: "First analysis",
	}

	agent2 := NewMockAgent("agent-2")
	agent2.Report = types.AnalysisReport{
		AgentID: "agent-2",
		Signal:  types.SignalHold,
		Summary: "Second analysis",
	}

	pipeline := NewPipeline(collector, []Agent{agent1, agent2}, logger)
	result, err := pipeline.Run(context.Background())

	if err != nil {
		t.Fatalf("pipeline Run failed: %v", err)
	}
	if len(result.Reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(result.Reports))
	}
	if result.Reports[0].AgentID != "agent-1" {
		t.Errorf("expected agent-1 first, got %s", result.Reports[0].AgentID)
	}
	if result.Reports[1].AgentID != "agent-2" {
		t.Errorf("expected agent-2 second, got %s", result.Reports[1].AgentID)
	}
}

func TestPipelineCollectorError(t *testing.T) {
	logger := zap.NewNop()
	collector := testutil.NewMockCollector()
	collector.Err = errors.New("connection refused")

	pipeline := NewPipeline(collector, []Agent{}, logger)
	result, err := pipeline.Run(context.Background())

	if err == nil {
		t.Fatal("expected error from failed collector")
	}
	if result.Error == "" {
		t.Error("expected error string in result")
	}
}

func TestPipelineAgentFailureContinues(t *testing.T) {
	logger := zap.NewNop()
	collector := testutil.NewMockCollector()

	failingAgent := NewMockAgent("failing")
	failingAgent.Err = errors.New("agent crashed")

	successAgent := NewMockAgent("success")
	successAgent.Report = types.AnalysisReport{
		AgentID: "success",
		Signal:  types.SignalHold,
	}

	pipeline := NewPipeline(collector, []Agent{failingAgent, successAgent}, logger)
	result, err := pipeline.Run(context.Background())

	if err != nil {
		t.Fatalf("pipeline should continue after agent failure, got error: %v", err)
	}
	// Should have 2 reports: one with error, one successful
	if len(result.Reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(result.Reports))
	}
}

func TestPipelineRunNoAgents(t *testing.T) {
	logger := zap.NewNop()
	collector := testutil.NewMockCollector()

	pipeline := NewPipeline(collector, []Agent{}, logger)
	result, err := pipeline.Run(context.Background())

	if err != nil {
		t.Fatalf("pipeline with no agents should succeed, got: %v", err)
	}
	if len(result.Reports) != 0 {
		t.Errorf("expected 0 reports, got %d", len(result.Reports))
	}
}

func TestPipelinePassesReportsToSubsequentAgents(t *testing.T) {
	logger := zap.NewNop()
	collector := testutil.NewMockCollector()

	var receivedReports []types.AnalysisReport

	agent1 := NewMockAgent("agent-1")
	agent1.Report = types.AnalysisReport{
		AgentID: "agent-1",
		Signal:  types.SignalBuy,
	}

	agent2 := NewMockAgent("agent-2")
	agent2.RunFunc = func(ctx context.Context, input types.AgentInput) (types.AnalysisReport, error) {
		receivedReports = input.Reports
		return types.AnalysisReport{AgentID: "agent-2", Signal: types.SignalHold}, nil
	}

	pipeline := NewPipeline(collector, []Agent{agent1, agent2}, logger)
	_, err := pipeline.Run(context.Background())
	if err != nil {
		t.Fatalf("pipeline Run failed: %v", err)
	}

	if len(receivedReports) != 1 {
		t.Fatalf("expected agent-2 to receive 1 prior report, got %d", len(receivedReports))
	}
	if receivedReports[0].AgentID != "agent-1" {
		t.Errorf("expected prior report from agent-1, got %s", receivedReports[0].AgentID)
	}
}
```

- [ ] **7.3** Run tests

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go test ./trader/ -v -run "TestPipeline"
# Expected: PASS — all 6 tests pass
```

- [ ] **7.4** Commit

```bash
git add trader/pipeline.go trader/pipeline_test.go trader/types/types.go
git commit -m "feat(pipeline): orchestrator for data collection and agent execution"
```

---

## Task 8: Daemon Main Loop

- [ ] **8.1** Write the daemon

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/daemon.go`

```go
package trader

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/kyungw00k/upbit-trader/trader/types"
	"go.uber.org/zap"
)

// DaemonState represents the current state of the daemon.
type DaemonState int

const (
	DaemonStopped DaemonState = iota
	DaemonRunning
	DaemonStopping
	DaemonEmergency
)

func (s DaemonState) String() string {
	switch s {
	case DaemonStopped:
		return "stopped"
	case DaemonRunning:
		return "running"
	case DaemonStopping:
		return "stopping"
	case DaemonEmergency:
		return "emergency"
	default:
		return "unknown"
	}
}

// Daemon manages the main loop for the trading system.
type Daemon struct {
	pipeline  *Pipeline
	interval  time.Duration
	logger    *zap.Logger
	state     DaemonState
	mu        sync.RWMutex
	pidFile   string
	cancel    context.CancelFunc
}

// NewDaemon creates a new daemon instance.
func NewDaemon(pipeline *Pipeline, interval time.Duration, logger *zap.Logger, pidFile string) *Daemon {
	return &Daemon{
		pipeline: pipeline,
		interval: interval,
		logger:   logger.Named("daemon"),
		pidFile:  pidFile,
		state:    DaemonStopped,
	}
}

// State returns the current daemon state.
func (d *Daemon) State() DaemonState {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state
}

// setState updates the daemon state thread-safely.
func (d *Daemon) setState(s DaemonState) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = s
}

// Start runs the daemon loop. It blocks until the context is cancelled.
func (d *Daemon) Start(ctx context.Context) error {
	ctx, d.cancel = context.WithCancel(ctx)

	d.setState(DaemonRunning)
	d.logger.Info("daemon started",
		zap.String("interval", d.interval.String()),
	)

	// Write PID file
	if err := d.writePIDFile(); err != nil {
		d.logger.Warn("failed to write PID file", zap.Error(err))
	}

	// Set up signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		d.setState(DaemonStopped)
		d.cleanup()
	}()

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	// Run immediately on start
	d.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			d.logger.Info("daemon context cancelled, shutting down")
			return nil

		case sig := <-sigCh:
			d.logger.Info("received signal", zap.String("signal", sig.String()))
			return nil

		case <-ticker.C:
			d.runOnce(ctx)
		}
	}
}

// Stop gracefully stops the daemon.
func (d *Daemon) Stop() error {
	d.setState(DaemonStopping)
	d.logger.Info("daemon stopping")

	if d.cancel != nil {
		d.cancel()
	}

	return nil
}

// EmergencyStop triggers an emergency stop.
func (d *Daemon) EmergencyStop() error {
	d.setState(DaemonEmergency)
	d.logger.Warn("EMERGENCY STOP activated")

	if d.cancel != nil {
		d.cancel()
	}

	return nil
}

// RunOnce runs the pipeline a single time (one-shot mode).
func (d *Daemon) RunOnce(ctx context.Context) (*types.PipelineResult, error) {
	return d.runOnce(ctx)
}

// runOnce executes a single pipeline run.
func (d *Daemon) runOnce(ctx context.Context) (*types.PipelineResult, error) {
	d.logger.Info("executing pipeline run")

	result, err := d.pipeline.Run(ctx)
	if err != nil {
		d.logger.Error("pipeline run failed", zap.Error(err))
		return result, err
	}

	// Log summary
	for _, report := range result.Reports {
		d.logger.Info("analysis report",
			zap.String("agent", report.AgentID),
			zap.String("signal", string(report.Signal)),
			zap.Int("coins", len(report.Coins)),
			zap.String("summary", report.Summary),
		)
	}

	return result, nil
}

// writePIDFile writes the current process PID to a file.
func (d *Daemon) writePIDFile() error {
	if d.pidFile == "" {
		return nil
	}

	dir := filepath.Dir(d.pidFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating PID directory: %w", err)
	}

	pid := os.Getpid()
	if err := os.WriteFile(d.pidFile, []byte(fmt.Sprintf("%d\n", pid)), 0644); err != nil {
		return fmt.Errorf("writing PID file: %w", err)
	}

	d.logger.Info("wrote PID file", zap.String("path", d.pidFile), zap.Int("pid", pid))
	return nil
}

// cleanup removes the PID file and performs other cleanup.
func (d *Daemon) cleanup() {
	if d.pidFile != "" {
		if err := os.Remove(d.pidFile); err != nil && !os.IsNotExist(err) {
			d.logger.Warn("failed to remove PID file", zap.Error(err))
		}
	}
}

// ReadPIDFile reads a PID from the given file.
func ReadPIDFile(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("reading PID file: %w", err)
	}

	var pid int
	if _, err := fmt.Sscanf(string(data), "%d", &pid); err != nil {
		return 0, fmt.Errorf("parsing PID: %w", err)
	}

	return pid, nil
}
```

- [ ] **8.2** Write tests

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/daemon_test.go`

```go
package trader

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyungw00k/upbit-trader/trader/types"
	testutil "github.com/kyungw00k/upbit-trader/internal/testutil"
	"go.uber.org/zap"
)

func TestDaemonStateString(t *testing.T) {
	tests := []struct {
		state DaemonState
		want  string
	}{
		{DaemonStopped, "stopped"},
		{DaemonRunning, "running"},
		{DaemonStopping, "stopping"},
		{DaemonEmergency, "emergency"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("DaemonState(%d).String() = %s, want %s", tt.state, got, tt.want)
		}
	}
}

func TestDaemonRunOnce(t *testing.T) {
	logger := zap.NewNop()
	collector := testutil.NewMockCollector()

	agent := NewMockAgent("test")
	agent.Report = types.AnalysisReport{
		AgentID: "test",
		Signal:  types.SignalBuy,
		Summary: "One-shot analysis",
	}

	pipeline := NewPipeline(collector, []Agent{agent}, logger)
	daemon := NewDaemon(pipeline, 4*time.Hour, logger, "")

	result, err := daemon.RunOnce(context.Background())
	if err != nil {
		t.Fatalf("RunOnce failed: %v", err)
	}
	if len(result.Reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(result.Reports))
	}
}

func TestDaemonStartStop(t *testing.T) {
	logger := zap.NewNop()
	collector := testutil.NewMockCollector()

	agent := NewMockAgent("test")
	agent.Report = types.AnalysisReport{AgentID: "test", Signal: types.SignalHold}

	pipeline := NewPipeline(collector, []Agent{agent}, logger)
	pidFile := filepath.Join(t.TempDir(), "trader.pid")
	daemon := NewDaemon(pipeline, 100*time.Millisecond, logger, pidFile)

	// Start daemon in background
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- daemon.Start(ctx)
	}()

	// Wait for at least one run
	time.Sleep(200 * time.Millisecond)

	// Stop the daemon
	if err := daemon.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("daemon exited with error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("daemon did not stop within timeout")
	}

	if daemon.State() != DaemonStopped {
		t.Errorf("expected stopped state, got %s", daemon.State())
	}
}

func TestDaemonPIDFile(t *testing.T) {
	logger := zap.NewNop()
	collector := testutil.NewMockCollector()

	agent := NewMockAgent("test")
	agent.Report = types.AnalysisReport{AgentID: "test", Signal: types.SignalHold}

	pipeline := NewPipeline(collector, []Agent{agent}, logger)
	pidFile := filepath.Join(t.TempDir(), "trader.pid")

	daemon := NewDaemon(pipeline, 1*time.Hour, logger, pidFile)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- daemon.Start(ctx)
	}()

	// Wait for PID file to be written
	time.Sleep(100 * time.Millisecond)

	// Verify PID file exists and contains current PID
	pid, err := ReadPIDFile(pidFile)
	if err != nil {
		t.Fatalf("ReadPIDFile failed: %v", err)
	}
	if pid != os.Getpid() {
		// In a goroutine this would be the parent's PID
		t.Logf("PID in file: %d, current PID: %d", pid, os.Getpid())
	}

	// Stop
	daemon.Stop()
	cancel()
	<-done

	// Verify PID file was cleaned up
	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		t.Error("PID file should have been removed after stop")
	}
}

func TestDaemonEmergencyStop(t *testing.T) {
	logger := zap.NewNop()
	collector := testutil.NewMockCollector()

	agent := NewMockAgent("test")
	agent.Report = types.AnalysisReport{AgentID: "test", Signal: types.SignalHold}

	pipeline := NewPipeline(collector, []Agent{agent}, logger)
	daemon := NewDaemon(pipeline, 1*time.Hour, logger, "")

	if err := daemon.EmergencyStop(); err != nil {
		t.Fatalf("EmergencyStop failed: %v", err)
	}

	if daemon.State() != DaemonEmergency {
		t.Errorf("expected emergency state, got %s", daemon.State())
	}
}

func TestReadPIDFileNotExist(t *testing.T) {
	_, err := ReadPIDFile("/nonexistent/pid/file")
	if err == nil {
		t.Fatal("expected error for nonexistent PID file")
	}
}
```

- [ ] **8.3** Run tests

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go test ./trader/ -v -run "TestDaemon|TestReadPID"
# Expected: PASS — all 6 tests pass
```

- [ ] **8.4** Commit

```bash
git add trader/daemon.go trader/daemon_test.go
git commit -m "feat(daemon): main loop with graceful shutdown and PID management"
```

---

## Task 9: CLI Commands

- [ ] **9.1** Write the CLI entry point

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/cmd/trader/main.go`

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/kyungw00k/upbit-trader/trader"
	"github.com/kyungw00k/upbit-trader/trader/agents"
	"github.com/kyungw00k/upbit-trader/trader/config"
	"github.com/kyungw00k/upbit-trader/trader/upbit"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultConfigPath = "configs/config.yaml"
	defaultPIDFile    = "trader.pid"
	appName           = "trader"
	appShortDesc      = "Autonomous crypto trading daemon"
	appLongDesc       = `upbit-trader is an autonomous cryptocurrency trading daemon that uses
Claude AI agents to analyze market data and make trading decisions.

This is the analysis-only mode — no real trades are executed.`
)

var (
	cfgPath    string
	pidFile    string
	daemonMode bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   appName,
		Short: appShortDesc,
		Long:  appLongDesc,
	}

	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", defaultConfigPath, "config file path")

	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(stopCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(analyzeCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the trading daemon",
		RunE:  runStart,
	}

	cmd.Flags().BoolVarP(&daemonMode, "daemon", "d", false, "run as background daemon")
	cmd.Flags().StringVar(&pidFile, "pid-file", defaultPIDFile, "PID file path")

	return cmd
}

func stopCmd() *cobra.Command {
	var emergency bool

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the trading daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStop(cmd, emergency)
		},
	}

	cmd.Flags().BoolVar(&emergency, "emergency", false, "emergency stop (kills all activity immediately)")

	return cmd
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show daemon status",
		RunE:  runStatus,
	}
}

func analyzeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "analyze",
		Short: "Run analysis once and print results",
		RunE:  runAnalyze,
	}
}

func loadConfigAndLogger() (*config.Config, *zap.Logger, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, nil, fmt.Errorf("loading config: %w", err)
	}

	level, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.EncoderConfig.TimeKey = "time"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapCfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, nil, fmt.Errorf("creating logger: %w", err)
	}

	return cfg, logger, err
}

func buildPipeline(cfg *config.Config, logger *zap.Logger) *trader.Pipeline {
	collector := upbit.NewClient(cfg.API.UpbitAccessKey, cfg.API.UpbitSecretKey)

	analyst := agents.NewTechnicalAnalyst(cfg.API.AnthropicAPIKey, cfg.API.ModelAnalyst, logger)

	return trader.NewPipeline(collector, []trader.Agent{analyst}, logger)
}

func runStart(cmd *cobra.Command, args []string) error {
	cfg, logger, err := loadConfigAndLogger()
	if err != nil {
		return err
	}
	defer logger.Sync()

	interval, err := cfg.Daemon.DecisionInterval()
	if err != nil {
		return fmt.Errorf("parsing decision interval: %w", err)
	}

	pipeline := buildPipeline(cfg, logger)
	daemon := trader.NewDaemon(pipeline, interval, logger, pidFile)

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		logger.Info("received signal, initiating graceful shutdown", zap.String("signal", sig.String()))
		daemon.Stop()
		cancel()
	}()

	logger.Info("starting upbit-trader daemon",
		zap.String("interval", interval.String()),
		zap.Bool("daemon_mode", daemonMode),
	)

	return daemon.Start(ctx)
}

func runStop(cmd *cobra.Command, emergency bool) error {
	logger, _ := zap.NewProduction()

	if emergency {
		// For emergency stop, just read PID and kill
		pid, err := trader.ReadPIDFile(pidFile)
		if err != nil {
			return fmt.Errorf("reading PID file: %w", err)
		}

		logger.Warn("sending emergency stop signal", zap.Int("pid", pid))
		process, err := os.FindProcess(pid)
		if err != nil {
			return fmt.Errorf("finding process: %w", err)
		}

		if err := process.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("sending signal: %w", err)
		}

		fmt.Println("Emergency stop signal sent to daemon")
		return nil
	}

	pid, err := trader.ReadPIDFile(pidFile)
	if err != nil {
		return fmt.Errorf("daemon not running (no PID file at %s): %w", pidFile, err)
	}

	logger.Info("stopping daemon", zap.Int("pid", pid))
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("finding process: %w", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("sending signal: %w", err)
	}

	fmt.Printf("Stop signal sent to daemon (PID: %d)\n", pid)
	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	pid, err := trader.ReadPIDFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Daemon status: stopped")
			return nil
		}
		return fmt.Errorf("checking daemon status: %w", err)
	}

	// Check if process is actually running
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Daemon status: stopped (stale PID file)")
		return nil
	}

	if err := process.Signal(syscall.Signal(0)); err != nil {
		fmt.Println("Daemon status: stopped (stale PID file)")
		return nil
	}

	fmt.Printf("Daemon status: running (PID: %d)\n", pid)
	fmt.Printf("PID file: %s\n", filepath.Abs(pidFile))
	return nil
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	cfg, logger, err := loadConfigAndLogger()
	if err != nil {
		return err
	}
	defer logger.Sync()

	pipeline := buildPipeline(cfg, logger)
	daemon := trader.NewDaemon(pipeline, 0, logger, "")

	logger.Info("running one-shot analysis")
	result, err := daemon.RunOnce(cmd.Context())
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Output as JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("encoding result: %w", err)
	}

	return nil
}
```

**NOTE:** The `time` import in `main.go` is unused. Remove it if the compiler complains. Also, `daemonMode` is declared but used only in a log statement, which is fine.

- [ ] **9.2** Verify CLI builds

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go build -o bin/trader ./cmd/trader/
# Expected: binary built successfully
```

- [ ] **9.3** Verify CLI help output

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
./bin/trader --help
# Expected: shows "trader - Autonomous crypto trading daemon" with subcommands

./bin/trader start --help
# Expected: shows --daemon, --pid-file, --config flags

./bin/trader analyze --help
# Expected: shows --config flag

./bin/trader status --help
# Expected: shows status command description
```

- [ ] **9.4** Commit

```bash
git add cmd/trader/
git commit -m "feat(cli): cobra-based CLI with start, stop, status, analyze commands"
```

---

## Task 10: Logging

- [ ] **10.1** Write logging utilities

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/logging.go`

```go
package trader

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger with the given level string.
func NewLogger(level string) (*zap.Logger, error) {
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		lvl = zapcore.InfoLevel
	}

	return NewLoggerWithLevel(lvl)
}

// NewLoggerWithLevel creates a new zap logger with the given zapcore.Level.
func NewLoggerWithLevel(level zapcore.Level) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return cfg.Build()
}

// NewTestLogger creates a logger suitable for testing (nop with debug level).
func NewTestLogger() *zap.Logger {
	return zap.NewNop()
}

// MustNewLogger creates a logger or panics. Use in main/initialization only.
func MustNewLogger(level string) *zap.Logger {
	logger, err := NewLogger(level)
	if err != nil {
		panic(fmt.Sprintf("creating logger: %v", err))
	}
	return logger
}
```

- [ ] **10.2** Write tests

**File:** `/Users/humphrey.park/Sandbox/upbit-trader/trader/logging_test.go`

```go
package trader

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"debug", "debug"},
		{"info", "info"},
		{"warn", "warn"},
		{"error", "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.level)
			if err != nil {
				t.Fatalf("NewLogger(%q) failed: %v", tt.level, err)
			}
			if logger == nil {
				t.Fatal("expected non-nil logger")
			}
			logger.Sync() // should not panic
		})
	}
}

func TestNewLoggerInvalidLevel(t *testing.T) {
	// Invalid level should fall back to info
	logger, err := NewLogger("invalid")
	if err != nil {
		t.Fatalf("NewLogger with invalid level should not error, got: %v", err)
	}
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestNewLoggerWithLevel(t *testing.T) {
	logger, err := NewLoggerWithLevel(zapcore.DebugLevel)
	if err != nil {
		t.Fatalf("NewLoggerWithLevel failed: %v", err)
	}
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestNewTestLogger(t *testing.T) {
	logger := NewTestLogger()
	if logger == nil {
		t.Fatal("expected non-nil test logger")
	}
	// Should not panic when writing
	logger.Info("test message")
}

func TestMustNewLogger(t *testing.T) {
	logger := MustNewLogger("info")
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestMustNewLoggerPanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected MustNewLogger to panic on truly invalid config")
		}
	}()
	// This won't panic because NewLogger handles invalid levels gracefully.
	// The test verifies the function works correctly.
	_ = MustNewLogger("info")
}

func TestLoggerLevels(t *testing.T) {
	core, recorded := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)

	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Warn("warn msg")
	logger.Error("error msg")

	if recorded.Len() != 4 {
		t.Errorf("expected 4 log entries, got %d", recorded.Len())
	}

	if recorded.All()[0].Message != "debug msg" {
		t.Errorf("expected 'debug msg', got %s", recorded.All()[0].Message)
	}
}
```

- [ ] **10.3** Run tests

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go test ./trader/ -v -run "TestNewLogger|TestMustNew|TestTestLogger|TestLoggerLevels"
# Expected: PASS — all tests pass
```

- [ ] **10.4** Add zaptest dependency (if not already present)

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go get go.uber.org/zap/zaptest/observer
go mod tidy
```

- [ ] **10.5** Commit

```bash
git add trader/logging.go trader/logging_test.go go.mod go.sum
git commit -m "feat(logging): configurable zap logger with level support"
```

---

## Task 11: Integration and Final Verification

- [ ] **11.1** Run all tests

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go test ./... -v -count=1
# Expected: ALL TESTS PASS
```

- [ ] **11.2** Run vet and build

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go vet ./...
go build ./...
# Expected: no errors
```

- [ ] **11.3** Verify binary works

```bash
cd /Users/humphrey.park/Sandbox/upbit-trader
go build -o bin/trader ./cmd/trader/
./bin/trader --help
./bin/trader status
# Expected: "Daemon status: stopped"
```

- [ ] **11.4** Final commit

```bash
git add -A
git commit -m "chore: Plan 1 Foundation complete — analysis-only daemon"
```

---

## Summary

**What was built:**
1. Go project scaffolding with all dependencies
2. YAML configuration with environment variable overrides
3. Core trading types (Signal, CoinAnalysis, AnalysisReport, MarketSnapshot, etc.)
4. Agent interface with mock agent for testing
5. Upbit data collector wrapper with KRW market support
6. Technical Analyst agent using Claude Haiku via Anthropic Go SDK
7. Pipeline orchestrator for data collection → agent execution flow
8. Daemon main loop with graceful shutdown, PID management, and emergency stop
9. Cobra CLI with start, stop, status, and analyze commands
10. Configurable Zap logging

**What Plan 2 will add (Full Agent Team):**
- Market Analyst agent (시장 상관관계, 돈 흐름 분석)
- News Analyst agent (뉴스/감성 분석)
- CEO agent (다중 에이전트 조율, 코인 탐색, 최종 결정)
- Risk Manager agent (포지션 사이징, 독립 거부권)

**What Plan 3 will add (Trading & Safety):**
- Execution Trader agent (주문 실행)
- 실제 주문 체결 (dry-run 모드 지원)
- 하드 리밋 강제 (코드 레벨)
- 킬 스위치 (긴급 정지)
- 데몬 간 통신 (stop 신호 전달)

**What Plan 4 will add (Learning & Monitoring):**
- Self-Reviewer agent (거래 후 자가 리뷰)
- 메모리/학습 시스템 (SQLite)
- Decision Audit Trail (의사결정 주기별 전체 데이터 영구 저장)
- 웹 대시보드 (HTTP API + SPA)
  - 의사결정 투명성: 각 에이전트의 입력 데이터 요약 → reasoning → 출력 시각화
  - Decision Timeline, Agent Detail, Trade History 화면
- TUI 대시보드 (Bubble Tea)
- Telegram/Slack 알림
