package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.API.GLMModel != "glm-4.7" {
		t.Errorf("expected default model glm-4.7, got %s", cfg.API.GLMModel)
	}
	if cfg.API.GLMBaseURL != "https://api.z.ai/v1" {
		t.Errorf("expected default base URL, got %s", cfg.API.GLMBaseURL)
	}
	if cfg.Daemon.DecisionInterval != "1h" {
		t.Errorf("expected default interval 1h, got %s", cfg.Daemon.DecisionInterval)
	}
	if cfg.Trading.Mode != "paper" {
		t.Errorf("expected default mode paper, got %s", cfg.Trading.Mode)
	}
	if cfg.Risk.MaxPositionPct != 0.30 {
		t.Errorf("expected max position 0.30, got %f", cfg.Risk.MaxPositionPct)
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
			got, err := dc.GetDecisionInterval()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestIndicatorInterval(t *testing.T) {
	dc := DaemonConfig{}
	got, err := dc.GetIndicatorInterval()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 15*time.Minute {
		t.Errorf("expected 15m, got %v", got)
	}
}

func TestSentimentInterval(t *testing.T) {
	dc := DaemonConfig{}
	got, err := dc.GetSentimentInterval()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 1*time.Hour {
		t.Errorf("expected 1h, got %v", got)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := []byte(`
api:
  glm_model: glm-4.7-flash
daemon:
  decision_interval: "1h"
trading:
  markets:
    - KRW-BTC
    - KRW-ETH
    - KRW-SOL
  mode: real
risk:
  max_daily_loss_pct: 0.05
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

	if cfg.API.GLMModel != "glm-4.7-flash" {
		t.Errorf("expected glm-4.7-flash, got %s", cfg.API.GLMModel)
	}
	dur, _ := cfg.Daemon.GetDecisionInterval()
	if dur != 1*time.Hour {
		t.Errorf("expected 1h, got %v", dur)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected debug, got %s", cfg.Log.Level)
	}
	if cfg.Trading.Mode != "real" {
		t.Errorf("expected real, got %s", cfg.Trading.Mode)
	}
	if len(cfg.Trading.Markets) != 3 {
		t.Errorf("expected 3 markets, got %d", len(cfg.Trading.Markets))
	}
	if cfg.Risk.MaxDailyLossPct != 0.05 {
		t.Errorf("expected 0.05, got %f", cfg.Risk.MaxDailyLossPct)
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
	if cfg.API.GLMModel != "glm-4.7" {
		t.Errorf("expected default model, got %s", cfg.API.GLMModel)
	}
}

func TestEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := []byte(`
api:
  glm_model: from-file
log:
  level: "info"
`)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("GLM_API_KEY", "test-glm-key")
	t.Setenv("GLM_MODEL", "from-env")
	t.Setenv("TRADER_LOG_LEVEL", "debug")
	t.Setenv("TRADER_MODE", "backtest")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.API.GLMModel != "from-env" {
		t.Errorf("expected env override from-env, got %s", cfg.API.GLMModel)
	}
	if cfg.API.GLMAPIKey != "test-glm-key" {
		t.Errorf("expected env override for GLM API key")
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected env override debug, got %s", cfg.Log.Level)
	}
	if cfg.Trading.Mode != "backtest" {
		t.Errorf("expected env override backtest, got %s", cfg.Trading.Mode)
	}
}

func TestEnvOverridesForUpbitKeys(t *testing.T) {
	cfg := DefaultConfig()

	t.Setenv("UPBIT_ACCESS_KEY", "access-test")
	t.Setenv("UPBIT_SECRET_KEY", "secret-test")

	applyEnvOverrides(cfg)

	if cfg.API.UpbitAccessKey != "access-test" {
		t.Errorf("expected upbit access key from env")
	}
	if cfg.API.UpbitSecretKey != "secret-test" {
		t.Errorf("expected upbit secret key from env")
	}
}
