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
	dur, _ := cfg.Daemon.GetDecisionInterval()
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
		t.Errorf("expected env override from-env, got %s", cfg.API.ModelAnalyst)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected env override debug, got %s", cfg.Log.Level)
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
