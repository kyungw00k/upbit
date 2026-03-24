package exchange

import (
	"github.com/kyungw00k/upbit/api"
)

// ExchangeClient Exchange API (거래) 클라이언트
type ExchangeClient struct {
	client *api.Client
}

// NewExchangeClient ExchangeClient 생성
func NewExchangeClient(client *api.Client) *ExchangeClient {
	return &ExchangeClient{client: client}
}
