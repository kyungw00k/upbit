package quotation

import (
	"context"
	"strings"
)

// TickSizeInfo 호가 정책 정보
type TickSizeInfo struct {
	Market          string   `json:"market"`
	QuoteCurrency   string   `json:"quote_currency"`
	TickSize        string   `json:"tick_size"`
	SupportedLevels []string `json:"supported_levels"`
}

// GetTickSizes 호가 정책 조회 (인증 불필요)
// API: GET /orderbook/instruments?markets=KRW-BTC,KRW-ETH
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
