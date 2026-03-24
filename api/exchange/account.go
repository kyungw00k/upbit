package exchange

import (
	"context"

	"github.com/kyungw00k/upbit/types"
)

// GetAccounts returns the account balances for all assets.
// API: GET /accounts
// See https://docs.upbit.com/reference/%EC%A0%84%EC%B2%B4-%EA%B3%84%EC%A2%8C-%EC%A1%B0%ED%9A%8C
func (c *ExchangeClient) GetAccounts(ctx context.Context) ([]types.Account, error) {
	var accounts []types.Account
	err := c.client.GET(ctx, "/accounts", nil, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
