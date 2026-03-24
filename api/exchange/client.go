// Package exchange provides access to Upbit's Exchange (trading) API.
// See https://docs.upbit.com/reference/%EC%9E%90%EC%82%B0 for API documentation.
package exchange

import (
	"github.com/kyungw00k/upbit/api"
)

// ExchangeClient is the client for the Exchange (trading) API.
type ExchangeClient struct {
	client *api.Client
}

// NewExchangeClient creates a new ExchangeClient.
func NewExchangeClient(client *api.Client) *ExchangeClient {
	return &ExchangeClient{client: client}
}
