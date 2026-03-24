package exchange

import (
	"context"

	"github.com/kyungw00k/upbit/types"
)

// GetOrderChance 페어별 주문 가능 정보 조회
// API: GET /orders/chance?market=KRW-BTC
func (c *ExchangeClient) GetOrderChance(ctx context.Context, market string) (*types.OrderChance, error) {
	var chance types.OrderChance
	query := map[string]string{
		"market": market,
	}
	err := c.client.GET(ctx, "/orders/chance", query, &chance)
	if err != nil {
		return nil, err
	}
	return &chance, nil
}
