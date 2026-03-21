package quotation

import (
	"context"

	"github.com/kyungw00k/upbit/internal/types"
)

// GetMarkets 마켓(거래쌍) 목록 조회
// API: GET /market/all?is_details=true
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
