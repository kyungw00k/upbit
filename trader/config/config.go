package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the root configuration for the trading daemon.
type Config struct {
	API       APIConfig       `yaml:"api"`
	Daemon    DaemonConfig    `yaml:"daemon"`
	Trading   TradingConfig   `yaml:"trading"`
	Risk      RiskConfig      `yaml:"risk"`
	Log       LogConfig       `yaml:"log"`
	Sentiment SentimentConfig `yaml:"sentiment"`
	Strategy  StrategyConfig  `yaml:"strategy"`
}

// APIConfig holds API keys and model settings.
type APIConfig struct {
	GLMAPIKey      string `yaml:"glm_api_key"`
	GLMBaseURL     string `yaml:"glm_base_url"`
	GLMModel       string `yaml:"glm_model"`
	UpbitAccessKey string `yaml:"upbit_access_key"`
	UpbitSecretKey string `yaml:"upbit_secret_key"`
}

// DaemonConfig controls daemon scheduling behavior.
type DaemonConfig struct {
	DecisionInterval string `yaml:"decision_interval"`
	IndicatorInterval string `yaml:"indicator_interval"`
	SentimentInterval string `yaml:"sentiment_interval"`
}

// TradingConfig controls which markets to trade and how.
type TradingConfig struct {
	Markets        []string `yaml:"markets"`
	Mode           string   `yaml:"mode"`
	InitialBalance float64  `yaml:"initial_balance"`
	Timeframes     []string `yaml:"timeframes"`
}

// RiskConfig controls risk management parameters.
type RiskConfig struct {
	MaxPositionPct  float64 `yaml:"max_position_pct"`
	MaxDailyLossPct float64 `yaml:"max_daily_loss_pct"`
	StopLossPct     float64 `yaml:"stop_loss_pct"`
	MaxPositions    int     `yaml:"max_positions"`
}

// LogConfig controls logging behavior.
type LogConfig struct {
	Level string `yaml:"level"`
}

// SentimentConfig controls sentiment data collection.
type SentimentConfig struct {
	FearGreedAPIURL string `yaml:"fear_greed_api_url"`
	FearGreedLimit  int    `yaml:"fear_greed_limit"`
}

// StrategyConfig controls the LLM strategy agent behavior.
type StrategyConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Timeout  string `yaml:"timeout"`
	BaseURL  string `yaml:"base_url"`
	Model    string `yaml:"model"`
}

// GetDecisionInterval returns the parsed decision interval duration.
func (c *DaemonConfig) GetDecisionInterval() (time.Duration, error) {
	if c.DecisionInterval == "" {
		return 4 * time.Hour, nil
	}
	return time.ParseDuration(c.DecisionInterval)
}

// GetIndicatorInterval returns the parsed indicator interval duration.
func (c *DaemonConfig) GetIndicatorInterval() (time.Duration, error) {
	if c.IndicatorInterval == "" {
		return 15 * time.Minute, nil
	}
	return time.ParseDuration(c.IndicatorInterval)
}

// GetSentimentInterval returns the parsed sentiment interval duration.
func (c *DaemonConfig) GetSentimentInterval() (time.Duration, error) {
	if c.SentimentInterval == "" {
		return 1 * time.Hour, nil
	}
	return time.ParseDuration(c.SentimentInterval)
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			GLMBaseURL: "https://api.z.ai/v1",
			GLMModel:   "glm-4.7",
		},
		Daemon: DaemonConfig{
			DecisionInterval:  "1h",
			IndicatorInterval: "15m",
			SentimentInterval: "1h",
		},
		Trading: TradingConfig{
			Markets:        []string{"KRW-BTC", "KRW-ETH"},
			Mode:           "paper",
			InitialBalance: 10000000,
			Timeframes:     []string{"1h", "4h", "1d", "1w"},
		},
		Risk: RiskConfig{
			MaxPositionPct:  0.30,
			MaxDailyLossPct: 0.03,
			StopLossPct:     0.05,
			MaxPositions:    5,
		},
		Log: LogConfig{
			Level: "info",
		},
		Sentiment: SentimentConfig{
			FearGreedAPIURL: "https://api.alternative.me/fng/",
			FearGreedLimit:  30,
		},
		Strategy: StrategyConfig{
			Enabled: true,
			Timeout: "15s",
		},
	}
}

// Load reads a YAML config file and applies environment variable overrides.
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
// Supports both GLM_* and OPENAI_* env vars (OPENAI_* takes precedence).
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("GLM_API_KEY"); v != "" {
		cfg.API.GLMAPIKey = v
	}
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		cfg.API.GLMAPIKey = v
	}
	if v := os.Getenv("GLM_BASE_URL"); v != "" {
		cfg.API.GLMBaseURL = v
	}
	if v := os.Getenv("OPENAI_BASE_URL"); v != "" {
		cfg.API.GLMBaseURL = v
	}
	if v := os.Getenv("GLM_MODEL"); v != "" {
		cfg.API.GLMModel = v
	}
	if v := os.Getenv("OPENAI_MODEL"); v != "" {
		cfg.API.GLMModel = v
	}
	if v := os.Getenv("UPBIT_ACCESS_KEY"); v != "" {
		cfg.API.UpbitAccessKey = v
	}
	if v := os.Getenv("UPBIT_SECRET_KEY"); v != "" {
		cfg.API.UpbitSecretKey = v
	}
	if v := os.Getenv("TRADER_LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
	if v := os.Getenv("TRADER_DECISION_INTERVAL"); v != "" {
		cfg.Daemon.DecisionInterval = v
	}
	if v := os.Getenv("TRADER_MODE"); v != "" {
		cfg.Trading.Mode = v
	}
	if v := os.Getenv("STRATEGY_ENABLED"); v == "true" {
		cfg.Strategy.Enabled = true
	}
	if v := os.Getenv("STRATEGY_ENABLED"); v == "false" {
		cfg.Strategy.Enabled = false
	}
	if v := os.Getenv("STRATEGY_TIMEOUT"); v != "" {
		cfg.Strategy.Timeout = v
	}
}
