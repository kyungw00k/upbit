package quotation

import (
	"context"
	"strings"

	"github.com/kyungw00k/upbit/types"
)

// GetOrderbooks retrieves the order book for the specified markets.
// API: GET /orderbook?markets=KRW-BTC,KRW-ETH
// See https://docs.upbit.com/reference/%ED%98%B8%EA%B0%80-%EC%A0%95%EB%B3%B4-%EC%A1%B0%ED%9A%8C
func (c *QuotationClient) GetOrderbooks(ctx context.Context, markets []string) ([]types.Orderbook, error) {
	var orderbooks []types.Orderbook
	query := map[string]string{
		"markets": strings.Join(markets, ","),
	}
	err := c.client.GET(ctx, "/orderbook", query, &orderbooks)
	if err != nil {
		return nil, err
	}
	return orderbooks, nil
}
