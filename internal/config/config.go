package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config 애플리케이션 설정
type Config struct {
	AccessKey     string `yaml:"access_key"`
	SecretKey     string `yaml:"secret_key"`
	DefaultOutput string `yaml:"default_output"`
}

// configPath 설정 파일 경로 결정
// 우선순위: $XDG_CONFIG_HOME/upbit/config.yaml → ~/.config/upbit/config.yaml → ~/.upbit/config.yaml
func configPath() (string, error) {
	// XDG_CONFIG_HOME 확인
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "upbit", "config.yaml"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home 디렉토리 확인 실패: %w", err)
	}

	// ~/.config/upbit/config.yaml (XDG 기본값)
	xdgDefault := filepath.Join(home, ".config", "upbit", "config.yaml")
	if _, err := os.Stat(xdgDefault); err == nil {
		return xdgDefault, nil
	}

	// ~/.upbit/config.yaml (폴백)
	return filepath.Join(home, ".upbit", "config.yaml"), nil
}

// Load 설정 로드
// 우선순위: 플래그 > 환경변수 > 설정 파일
func Load() (*Config, error) {
	cfg := &Config{}

	// 설정 파일 로드
	path, err := configPath()
	if err != nil {
		return cfg, nil // 설정 파일 없어도 계속 진행
	}

	if err := loadFile(cfg, path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("설정 파일 로드 실패 (%s): %w", path, err)
	}

	// 환경변수 오버라이드
	if v := os.Getenv("UPBIT_ACCESS_KEY"); v != "" {
		cfg.AccessKey = v
	}
	if v := os.Getenv("UPBIT_SECRET_KEY"); v != "" {
		cfg.SecretKey = v
	}

	return cfg, nil
}

// loadFile YAML 파일에서 설정 로드 (key: value 형태의 간단한 파서)
// 주의: access_key, secret_key는 환경변수에서만 읽으므로 파일에서 무시
func loadFile(cfg *Config, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])

		// 따옴표 제거
		value = strings.Trim(value, `"'`)

		switch key {
		case "default_output":
			cfg.DefaultOutput = value
		// access_key, secret_key는 환경변수에서만 읽음 (파일 무시)
		}
	}

	return scanner.Err()
}

// Save 설정 파일에 저장 (파일 권한 0600)
// 주의: access_key, secret_key는 환경변수 전용이므로 파일에 저장하지 않음
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return fmt.Errorf("설정 파일 경로 확인 실패: %w", err)
	}

	// 디렉토리 생성
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("디렉토리 생성 실패 (%s): %w", dir, err)
	}

	content := fmt.Sprintf(
		"default_output: %s\n",
		cfg.DefaultOutput,
	)

	// 0600 권한으로 저장
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("설정 파일 저장 실패: %w", err)
	}

	return nil
}

// Path 현재 설정 파일 경로 반환
func Path() (string, error) {
	return configPath()
}

// ApplyFlags 플래그 값을 설정에 적용 (최우선순위)
func (c *Config) ApplyFlags(accessKey, secretKey string) {
	if accessKey != "" {
		c.AccessKey = accessKey
	}
	if secretKey != "" {
		c.SecretKey = secretKey
	}
}
