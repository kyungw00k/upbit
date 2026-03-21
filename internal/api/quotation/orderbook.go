package quotation

import (
	"context"
	"strings"

	"github.com/kyungw00k/upbit/internal/types"
)

// GetOrderbooks 호가 조회
// API: GET /orderbook?markets=KRW-BTC,KRW-ETH
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
