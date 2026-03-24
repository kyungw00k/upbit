package exchange

import (
	"context"

	"github.com/kyungw00k/upbit/types"
)

// GetOrderChance returns order availability information for a given market pair.
// API: GET /orders/chance?market=KRW-BTC
// See https://docs.upbit.com/reference/%EC%A3%BC%EB%AC%B8-%EA%B0%80%EB%8A%A5-%EC%A0%95%EB%B3%B4
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
