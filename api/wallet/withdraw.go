package wallet

import (
	"context"
	"fmt"

	"github.com/kyungw00k/upbit/types"
)

// ListWithdrawals returns a list of withdrawals.
// API: GET /withdraws
// See https://docs.upbit.com/reference/%EC%B6%9C%EA%B8%88-%EB%A6%AC%EC%8A%A4%ED%8A%B8-%EC%A1%B0%ED%9A%8C
func (c *WalletClient) ListWithdrawals(ctx context.Context, currency string, state string, limit int, page int) ([]types.Withdrawal, error) {
	var withdrawals []types.Withdrawal
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

	err := c.client.GET(ctx, "/withdraws", query, &withdrawals)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}

// GetWithdrawal returns a single withdrawal by UUID.
// API: GET /withdraw
// See https://docs.upbit.com/reference/%EC%B6%9C%EA%B8%88-%EB%A6%AC%EC%8A%A4%ED%8A%B8-%EC%A1%B0%ED%9A%8C
func (c *WalletClient) GetWithdrawal(ctx context.Context, uuid string) (*types.Withdrawal, error) {
	var withdrawal types.Withdrawal
	query := map[string]string{}

	if uuid != "" {
		query["uuid"] = uuid
	}

	err := c.client.GET(ctx, "/withdraw", query, &withdrawal)
	if err != nil {
		return nil, err
	}
	return &withdrawal, nil
}

// WithdrawCoinRequest holds parameters for a digital asset withdrawal request.
type WithdrawCoinRequest struct {
	Currency         string `json:"currency"`
	NetType          string `json:"net_type"`
	Amount           string `json:"amount"`
	Address          string `json:"address"`
	SecondaryAddress string `json:"secondary_address,omitempty"`
	TransactionType  string `json:"transaction_type,omitempty"`
}

// WithdrawCoin submits a digital asset withdrawal request.
// API: POST /withdraws/coin
// See https://docs.upbit.com/reference/%EC%B6%9C%EA%B8%88%ED%95%98%EA%B8%B0
func (c *WalletClient) WithdrawCoin(ctx context.Context, req *WithdrawCoinRequest) (*types.Withdrawal, error) {
	var withdrawal types.Withdrawal
	err := c.client.POST(ctx, "/withdraws/coin", req, &withdrawal)
	if err != nil {
		return nil, err
	}
	return &withdrawal, nil
}

// WithdrawKRW submits a KRW withdrawal request.
// API: POST /withdraws/krw
// See https://docs.upbit.com/reference/%EC%B6%9C%EA%B8%88%ED%95%98%EA%B8%B0
func (c *WalletClient) WithdrawKRW(ctx context.Context, amount string, twoFactorType string) (*types.Withdrawal, error) {
	var withdrawal types.Withdrawal
	body := map[string]string{
		"amount":          amount,
		"two_factor_type": twoFactorType,
	}

	err := c.client.POST(ctx, "/withdraws/krw", body, &withdrawal)
	if err != nil {
		return nil, err
	}
	return &withdrawal, nil
}

// CancelWithdrawal cancels a digital asset withdrawal by UUID.
// API: DELETE /withdraws/coin?uuid=xxx
// See https://docs.upbit.com/reference/%EC%B6%9C%EA%B8%88-%EC%B7%A8%EC%86%8C
func (c *WalletClient) CancelWithdrawal(ctx context.Context, uuid string) (*types.Withdrawal, error) {
	var withdrawal types.Withdrawal
	query := map[string]string{
		"uuid": uuid,
	}

	err := c.client.DELETE(ctx, "/withdraws/coin", query, &withdrawal)
	if err != nil {
		return nil, err
	}
	return &withdrawal, nil
}

// GetAvailableWithdrawal returns withdrawal availability information for a given currency.
// API: GET /withdraws/chance
func (c *WalletClient) GetAvailableWithdrawal(ctx context.Context, currency string, netType string) (*types.WithdrawalChance, error) {
	var chance types.WithdrawalChance
	query := map[string]string{
		"currency": currency,
	}
	if netType != "" {
		query["net_type"] = netType
	}

	err := c.client.GET(ctx, "/withdraws/chance", query, &chance)
	if err != nil {
		return nil, err
	}
	return &chance, nil
}

// ListWithdrawalAddresses returns the list of allowed withdrawal addresses.
// API: GET /withdraws/coin_addresses
func (c *WalletClient) ListWithdrawalAddresses(ctx context.Context) ([]types.WithdrawalAddress, error) {
	var addresses []types.WithdrawalAddress
	err := c.client.GET(ctx, "/withdraws/coin_addresses", nil, &addresses)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}
