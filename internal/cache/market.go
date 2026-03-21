package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// MarketCache 마켓 목록 파일 캐시
// JSON 파일로 저장, TTL 기반 만료 (기본 1시간)
type MarketCache struct {
	dir string
	ttl time.Duration
}

type marketCacheData struct {
	UpdatedAt time.Time `json:"updated_at"`
	Markets   []string  `json:"markets"`
}

// NewMarketCache 마켓 캐시 생성
func NewMarketCache(ttl time.Duration) (*MarketCache, error) {
	dir, err := CandleCacheDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &MarketCache{dir: dir, ttl: ttl}, nil
}

// Get 캐시된 마켓 목록 반환. 만료 시 nil 반환.
func (c *MarketCache) Get() map[string]bool {
	path := filepath.Join(c.dir, "markets.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var cached marketCacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil
	}

	if time.Since(cached.UpdatedAt) > c.ttl {
		return nil // 만료
	}

	result := make(map[string]bool, len(cached.Markets))
	for _, m := range cached.Markets {
		result[m] = true
	}
	return result
}

// Set 마켓 목록을 캐시에 저장
func (c *MarketCache) Set(markets []string) error {
	cached := marketCacheData{
		UpdatedAt: time.Now(),
		Markets:   markets,
	}
	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}
	path := filepath.Join(c.dir, "markets.json")
	return os.WriteFile(path, data, 0600)
}
