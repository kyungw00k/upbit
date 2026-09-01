// Package pocket provides access to Upbit's Pocket API (main/sub account partitions).
// See https://docs.upbit.com/reference/pocket-overview for API documentation.
package pocket

import (
	"github.com/kyungw00k/upbit/api"
)

// PocketClient is a client for the Pocket API.
type PocketClient struct {
	client *api.Client
}

// NewPocketClient creates a new PocketClient.
func NewPocketClient(client *api.Client) *PocketClient {
	return &PocketClient{client: client}
}
