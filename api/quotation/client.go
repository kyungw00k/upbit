// Package quotation provides access to Upbit's Quotation (market data) API.
// See https://docs.upbit.com/reference/%EC%8B%9C%EC%84%B8-%EC%A2%85%EB%AA%A9-%EC%A1%B0%ED%9A%8C for API documentation.
package quotation

import (
	"github.com/kyungw00k/upbit/api"
)

// QuotationClient is a client for the Quotation API (market data).
type QuotationClient struct {
	client *api.Client
}

// NewQuotationClient creates a new QuotationClient.
func NewQuotationClient(client *api.Client) *QuotationClient {
	return &QuotationClient{client: client}
}
