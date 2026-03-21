package wallet

import (
	"context"

	"github.com/kyungw00k/upbit/internal/types"
)

// GetServiceStatus 입출금 서비스 상태 조회
// API: GET /status/wallet
func (c *WalletClient) GetServiceStatus(ctx context.Context) ([]types.ServiceStatus, error) {
	var statuses []types.ServiceStatus
	err := c.client.GET(ctx, "/status/wallet", nil, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
