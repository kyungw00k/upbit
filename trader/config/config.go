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

// GetDecisionInterval returns the parsed decision interval duration.
func (c *DaemonConfig) GetDecisionInterval() (time.Duration, error) {
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
