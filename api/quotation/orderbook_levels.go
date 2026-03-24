package quotation

import (
	"context"
	"strings"
)

// OrderbookLevel holds the supported aggregation levels for an order book market.
type OrderbookLevel struct {
	Market          string    `json:"market"`
	SupportedLevels []float64 `json:"supported_levels"`
}

// GetOrderbookLevels retrieves the supported order book aggregation levels for the specified markets (no authentication required).
// API: GET /orderbook/supported_levels?markets=KRW-BTC,KRW-ETH
// See https://docs.upbit.com/reference/%ED%98%B8%EA%B0%80-%EB%AA%A8%EC%95%84%EB%B3%B4%EA%B8%B0-%EB%8B%A8%EC%9C%84-%EC%A1%B0%ED%9A%8C
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
