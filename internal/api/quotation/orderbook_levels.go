package quotation

import (
	"context"
	"strings"
)

// OrderbookLevel 호가 모아보기 단위 정보
type OrderbookLevel struct {
	Market          string    `json:"market"`
	SupportedLevels []float64 `json:"supported_levels"`
}

// GetOrderbookLevels 호가 모아보기 단위 조회 (인증 불필요)
// API: GET /orderbook/supported_levels?markets=KRW-BTC,KRW-ETH
func (c *QuotationClient) GetOrderbookLevels(ctx context.Context, markets []string) ([]OrderbookLevel, error) {
	var levels []OrderbookLevel
	query := map[string]string{
		"markets": strings.Join(markets, ","),
	}
	err := c.client.GET(ctx, "/orderbook/supported_levels", query, &levels)
	if err != nil {
		return nil, err
	}
	return levels, nil
}
