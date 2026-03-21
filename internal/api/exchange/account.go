package exchange

import (
	"context"

	"github.com/kyungw00k/upbit/internal/types"
)

// GetAccounts 계정 잔고 조회
// API: GET /accounts
func (c *ExchangeClient) GetAccounts(ctx context.Context) ([]types.Account, error) {
	var accounts []types.Account
	err := c.client.GET(ctx, "/accounts", nil, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
