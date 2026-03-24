// Package wallet provides access to Upbit's Wallet (deposit/withdrawal) API.
// See https://docs.upbit.com/reference/%EC%9E%85%EA%B8%88-%EB%A6%AC%EC%8A%A4%ED%8A%B8-%EC%A1%B0%ED%9A%8C for API documentation.
package wallet

import (
	"github.com/kyungw00k/upbit/api"
)

// WalletClient is a client for the Wallet API (deposit/withdrawal).
type WalletClient struct {
	client *api.Client
}

// NewWalletClient creates a new WalletClient.
func NewWalletClient(client *api.Client) *WalletClient {
	return &WalletClient{client: client}
}
