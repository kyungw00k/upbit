package wallet

import (
	"context"
	"fmt"

	"github.com/kyungw00k/upbit/internal/types"
)

// ListDeposits 입금 목록 조회
// API: GET /deposits
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

// GetDeposit 개별 입금 조회
// API: GET /deposit
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

// ListDepositAddresses 입금 주소 목록 조회
// API: GET /deposits/coin_addresses
func (c *WalletClient) ListDepositAddresses(ctx context.Context) ([]types.DepositAddress, error) {
	var addresses []types.DepositAddress
	err := c.client.GET(ctx, "/deposits/coin_addresses", nil, &addresses)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// GetDepositAddress 개별 입금 주소 조회
// API: GET /deposits/coin_address
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

// CreateDepositAddress 입금 주소 생성 요청
// API: POST /deposits/generate_coin_address
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

// DepositKRW 원화 입금
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

// GetAvailableDeposit 디지털 자산 입금 가능 정보 조회
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
