package quotation

import (
	"context"
	"fmt"

	"github.com/kyungw00k/upbit/types"
)

// GetRecentTrades retrieves recent trade history for the specified market.
// API: GET /trades/ticks?market=KRW-BTC&count=20
// See https://docs.upbit.com/reference/%EC%B5%9C%EA%B7%BC-%EC%B2%B4%EA%B2%B0-%EB%82%B4%EC%97%AD
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
