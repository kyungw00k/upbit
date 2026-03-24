package quotation

import (
	"context"
	"strings"
)

// TickSizeInfo holds the tick size policy information for a market.
type TickSizeInfo struct {
	Market          string   `json:"market"`
	QuoteCurrency   string   `json:"quote_currency"`
	TickSize        string   `json:"tick_size"`
	SupportedLevels []string `json:"supported_levels"`
}

// GetTickSizes retrieves the tick size policy for the specified markets (no authentication required).
// API: GET /orderbook/instruments?markets=KRW-BTC,KRW-ETH
// See https://docs.upbit.com/reference/%ED%98%B8%EA%B0%80-%EB%8B%A8%EC%9C%84-%EC%A1%B0%ED%9A%8C
func (c *QuotationClient) GetTickSizes(ctx context.Context, markets []string) ([]TickSizeInfo, error) {
	var tickSizes []TickSizeInfo
	query := map[string]string{
		"markets": strings.Join(markets, ","),
	}
	err := c.client.GET(ctx, "/orderbook/instruments", query, &tickSizes)
	if err != nil {
		return nil, err
	}
	return tickSizes, nil
}
