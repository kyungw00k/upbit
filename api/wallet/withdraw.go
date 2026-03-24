package wallet

import (
	"context"
	"fmt"

	"github.com/kyungw00k/upbit/types"
)

// ListWithdrawals 출금 목록 조회
// API: GET /withdraws
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

// GetWithdrawal 개별 출금 조회
// API: GET /withdraw
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

// WithdrawCoinRequest 디지털 자산 출금 요청 파라미터
type WithdrawCoinRequest struct {
	Currency         string `json:"currency"`
	NetType          string `json:"net_type"`
	Amount           string `json:"amount"`
	Address          string `json:"address"`
	SecondaryAddress string `json:"secondary_address,omitempty"`
	TransactionType  string `json:"transaction_type,omitempty"`
}

// WithdrawCoin 디지털 자산 출금 요청
// API: POST /withdraws/coin
func (c *WalletClient) WithdrawCoin(ctx context.Context, req *WithdrawCoinRequest) (*types.Withdrawal, error) {
	var withdrawal types.Withdrawal
	err := c.client.POST(ctx, "/withdraws/coin", req, &withdrawal)
	if err != nil {
		return nil, err
	}
	return &withdrawal, nil
}

// WithdrawKRW 원화 출금 요청
// API: POST /withdraws/krw
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

// CancelWithdrawal 디지털 자산 출금 취소 요청
// API: DELETE /withdraws/coin?uuid=xxx
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

// GetAvailableWithdrawal 출금 가능 정보 조회
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

// ListWithdrawalAddresses 출금 허용 주소 목록 조회
// API: GET /withdraws/coin_addresses
func (c *WalletClient) ListWithdrawalAddresses(ctx context.Context) ([]types.WithdrawalAddress, error) {
	var addresses []types.WithdrawalAddress
	err := c.client.GET(ctx, "/withdraws/coin_addresses", nil, &addresses)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}
