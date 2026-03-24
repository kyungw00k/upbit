package quotation

import (
	"context"

	"github.com/kyungw00k/upbit/types"
)

// GetMarkets retrieves the list of all available markets (trading pairs).
// API: GET /market/all?is_details=true
// See https://docs.upbit.com/reference/%EB%A7%88%EC%BC%93-%EC%BD%94%EB%93%9C-%EC%A1%B0%ED%9A%8C
func (c *QuotationClient) GetMarkets(ctx context.Context) ([]types.Market, error) {
	var markets []types.Market
	query := map[string]string{
		"is_details": "true",
	}
	err := c.client.GET(ctx, "/market/all", query, &markets)
	if err != nil {
		return nil, err
	}
	return markets, nil
}
