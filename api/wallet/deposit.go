package wallet

import (
	"context"
	"fmt"

	"github.com/kyungw00k/upbit/types"
)

// ListDeposits returns a list of deposits.
// API: GET /deposits
// See https://docs.upbit.com/reference/%EC%9E%85%EA%B8%88-%EB%A6%AC%EC%8A%A4%ED%8A%B8-%EC%A1%B0%ED%9A%8C
func (c *WalletClient) ListDeposits(ctx context.Context, currency string, state string, limit int, page int) ([]types.Deposit, error) {
	var deposits []types.Deposit
	query := map[string]string{}

	if currency != "" {
		query["currency"] = currency
	}
	if state != "" {
		query["state"] = state
	}
	if limit > 0 {
		query["limit"] = fmt.Sprintf("%d", limit)
	}
	if page > 0 {
		query["page"] = fmt.Sprintf("%d", page)
	}

	err := c.client.GET(ctx, "/deposits", query, &deposits)
	if err != nil {
		return nil, err
	}
	return deposits, nil
}

// GetDeposit returns a single deposit by UUID.
// API: GET /deposit
// See https://docs.upbit.com/reference/%EC%9E%85%EA%B8%88-%EB%A6%AC%EC%8A%A4%ED%8A%B8-%EC%A1%B0%ED%9A%8C
func (c *WalletClient) GetDeposit(ctx context.Context, uuid string) (*types.Deposit, error) {
	var deposit types.Deposit
	query := map[string]string{}

	if uuid != "" {
		query["uuid"] = uuid
	}

	err := c.client.GET(ctx, "/deposit", query, &deposit)
	if err != nil {
		return nil, err
	}
	return &deposit, nil
}

// ListDepositAddresses returns a list of deposit addresses.
// API: GET /deposits/coin_addresses
// See https://docs.upbit.com/reference/%EC%9E%85%EA%B8%88-%EC%A3%BC%EC%86%8C-%EC%83%9D%EC%84%B1-%EC%9A%94%EC%B2%AD
func (c *WalletClient) ListDepositAddresses(ctx context.Context) ([]types.DepositAddress, error) {
	var addresses []types.DepositAddress
	err := c.client.GET(ctx, "/deposits/coin_addresses", nil, &addresses)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// GetDepositAddress returns a deposit address for a given currency and network type.
// API: GET /deposits/coin_address
// See https://docs.upbit.com/reference/%EC%9E%85%EA%B8%88-%EC%A3%BC%EC%86%8C-%EC%83%9D%EC%84%B1-%EC%9A%94%EC%B2%AD
func (c *WalletClient) GetDepositAddress(ctx context.Context, currency string, netType string) (*types.DepositAddress, error) {
	var address types.DepositAddress
	query := map[string]string{
		"currency": currency,
		"net_type": netType,
	}

	err := c.client.GET(ctx, "/deposits/coin_address", query, &address)
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// CreateDepositAddress requests generation of a deposit address for a given currency and network type.
// API: POST /deposits/generate_coin_address
// See https://docs.upbit.com/reference/%EC%9E%85%EA%B8%88-%EC%A3%BC%EC%86%8C-%EC%83%9D%EC%84%B1-%EC%9A%94%EC%B2%AD
func (c *WalletClient) CreateDepositAddress(ctx context.Context, currency string, netType string) (*types.CreateDepositAddressResult, error) {
	var result types.CreateDepositAddressResult
	body := map[string]string{
		"currency": currency,
		"net_type": netType,
	}

	err := c.client.POST(ctx, "/deposits/generate_coin_address", body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DepositKRW requests a KRW deposit.
// API: POST /deposits/krw
func (c *WalletClient) DepositKRW(ctx context.Context, amount string, twoFactorType string) (*types.Deposit, error) {
	var deposit types.Deposit
	body := map[string]string{
		"amount":          amount,
		"two_factor_type": twoFactorType,
	}

	err := c.client.POST(ctx, "/deposits/krw", body, &deposit)
	if err != nil {
		return nil, err
	}
	return &deposit, nil
}

// GetAvailableDeposit returns deposit availability information for a digital asset.
// API: GET /deposits/chance/coin
func (c *WalletClient) GetAvailableDeposit(ctx context.Context, currency string, netType string) (*types.AvailableDeposit, error) {
	var info types.AvailableDeposit
	query := map[string]string{
		"currency": currency,
		"net_type": netType,
	}

	err := c.client.GET(ctx, "/deposits/chance/coin", query, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
