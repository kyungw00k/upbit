package wallet

import (
	"context"

	"github.com/kyungw00k/upbit/types"
)

// GetServiceStatus returns the deposit/withdrawal service status for all currencies.
// API: GET /status/wallet
// See https://docs.upbit.com/reference/%EC%9E%85%EC%B6%9C%EA%B8%88-%EC%83%81%ED%83%9C
func (c *WalletClient) GetServiceStatus(ctx context.Context) ([]types.ServiceStatus, error) {
	var statuses []types.ServiceStatus
	err := c.client.GET(ctx, "/status/wallet", nil, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
