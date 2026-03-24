package exchange

import (
	"context"

	"github.com/kyungw00k/upbit/types"
)

// GetAPIKeys API 키 목록 조회
// API: GET /api_keys
func (c *ExchangeClient) GetAPIKeys(ctx context.Context) ([]types.APIKey, error) {
	var keys []types.APIKey
	err := c.client.GET(ctx, "/api_keys", nil, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
