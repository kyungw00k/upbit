package quotation

import (
	"context"
	"fmt"

	"github.com/kyungw00k/upbit/types"
)

// GetRecentTrades 최근 체결 내역 조회
// API: GET /trades/ticks?market=KRW-BTC&count=20
func (c *QuotationClient) GetRecentTrades(ctx context.Context, market string, count int) ([]types.Trade, error) {
	var trades []types.Trade
	query := map[string]string{
		"market": market,
		"count":  fmt.Sprintf("%d", count),
	}
	err := c.client.GET(ctx, "/trades/ticks", query, &trades)
	if err != nil {
		return nil, err
	}
	return trades, nil
}
