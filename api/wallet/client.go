package wallet

import (
	"github.com/kyungw00k/upbit/api"
)

// WalletClient Wallet API (입출금) 클라이언트
type WalletClient struct {
	client *api.Client
}

// NewWalletClient WalletClient 생성
func NewWalletClient(client *api.Client) *WalletClient {
	return &WalletClient{client: client}
}
