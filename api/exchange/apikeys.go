package exchange

import (
	"context"

	"github.com/kyungw00k/upbit/types"
)

// GetAPIKeys returns the list of API keys for the authenticated account.
// API: GET /api_keys
// See https://docs.upbit.com/reference/open-api-%ED%82%A4-%EB%A6%AC%EC%8A%A4%ED%8A%B8-%EC%A1%B0%ED%9A%8C
func (c *ExchangeClient) GetAPIKeys(ctx context.Context) ([]types.APIKey, error) {
	var keys []types.APIKey
	err := c.client.GET(ctx, "/api_keys", nil, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
