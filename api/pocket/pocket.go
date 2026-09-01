package pocket

import (
	"context"
	"strconv"
	"strings"

	"github.com/kyungw00k/upbit/types"
)

// ListPockets returns the caller's pockets (main + sub).
// API: GET /pockets
// See https://docs.upbit.com/reference/list-pockets
func (c *PocketClient) ListPockets(ctx context.Context) ([]types.Pocket, error) {
	var pockets []types.Pocket
	if err := c.client.GET(ctx, "/pockets", nil, &pockets); err != nil {
		return nil, err
	}
	return pockets, nil
}

// SubPocketBalance returns the balance of a sub pocket.
// API: GET /pockets/assets?uuid=xxx
// See https://docs.upbit.com/reference/get-sub-pocket-balance
func (c *PocketClient) SubPocketBalance(ctx context.Context, uuid string) ([]types.Account, error) {
	var accounts []types.Account
	query := map[string]string{"uuid": uuid}
	if err := c.client.GET(ctx, "/pockets/assets", query, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

// ListPocketAPIKeys returns API keys grouped by pocket.
// API: GET /pockets/api_keys?uuids[]=xxx&include_expired=true
// See https://docs.upbit.com/reference/list-pocket-api-keys
func (c *PocketClient) ListPocketAPIKeys(ctx context.Context, uuids []string, includeExpired bool) ([]types.PocketAPIKeys, error) {
	var keys []types.PocketAPIKeys
	var parts []string
	for _, id := range uuids {
		if id != "" {
			parts = append(parts, "uuids[]="+id)
		}
	}
	if includeExpired {
		parts = append(parts, "include_expired=true")
	}
	rawQuery := strings.Join(parts, "&")
	if err := c.client.GETWithRawQuery(ctx, "/pockets/api_keys", rawQuery, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

// ListTransfersOptions filters pocket transfer history.
// Direction: "in" or "out". States: e.g. WAIT, OK. OrderBy: asc, desc (default desc).
type ListTransfersOptions struct {
	Direction   string
	States      []string
	UUIDs       []string
	Identifiers []string
	StartTime   string
	EndTime     string
	Currency    string
	Limit       int
	OrderBy     string
}

// ListTransfers returns transfers between the caller's pockets.
// API: GET /pockets/transfers
// See https://docs.upbit.com/reference/list-transfers
func (c *PocketClient) ListTransfers(ctx context.Context, opts ListTransfersOptions) ([]types.PocketTransfer, error) {
	var transfers []types.PocketTransfer
	rawQuery := buildTransferQuery(opts, "", "")
	if err := c.client.GETWithRawQuery(ctx, "/pockets/transfers", rawQuery, &transfers); err != nil {
		return nil, err
	}
	return transfers, nil
}

// TransferRequest moves assets from the API key's pocket to another pocket.
// API: POST /pockets/transfers
// See https://docs.upbit.com/reference/transfer
type TransferRequest struct {
	To         string `json:"to"`
	Currency   string `json:"currency"`
	Amount     string `json:"amount"`
	Identifier string `json:"identifier,omitempty"`
}

// Transfer moves assets between pockets (source: the API key's pocket).
func (c *PocketClient) Transfer(ctx context.Context, req *TransferRequest) (*types.PocketTransfer, error) {
	var transfer types.PocketTransfer
	if err := c.client.POST(ctx, "/pockets/transfers", req, &transfer); err != nil {
		return nil, err
	}
	return &transfer, nil
}

// ListUniversalTransfersOptions filters main-pocket transfer history.
type ListUniversalTransfersOptions struct {
	ListTransfersOptions
	From string // sender pocket UUID
	To   string // receiver pocket UUID
}

// ListUniversalTransfers returns main-pocket asset transfer history.
// API: GET /pockets/universal_transfers
// See https://docs.upbit.com/reference/list-universal-transfers
func (c *PocketClient) ListUniversalTransfers(ctx context.Context, opts ListUniversalTransfersOptions) ([]types.PocketTransfer, error) {
	var transfers []types.PocketTransfer
	rawQuery := buildTransferQuery(opts.ListTransfersOptions, opts.From, opts.To)
	if err := c.client.GETWithRawQuery(ctx, "/pockets/universal_transfers", rawQuery, &transfers); err != nil {
		return nil, err
	}
	return transfers, nil
}

// UniversalTransferRequest moves main-pocket assets to/from a sub pocket.
// API: POST /pockets/universal_transfers
// See https://docs.upbit.com/reference/universal-transfer
type UniversalTransferRequest struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Currency   string `json:"currency"`
	Amount     string `json:"amount"`
	Identifier string `json:"identifier,omitempty"`
}

// UniversalTransfer moves main-pocket assets between pockets explicitly.
func (c *PocketClient) UniversalTransfer(ctx context.Context, req *UniversalTransferRequest) (*types.PocketTransfer, error) {
	var transfer types.PocketTransfer
	if err := c.client.POST(ctx, "/pockets/universal_transfers", req, &transfer); err != nil {
		return nil, err
	}
	return &transfer, nil
}

// buildTransferQuery builds the query string for transfer list endpoints.
// Array params keep the raw uuids[]=a form required by Upbit.
func buildTransferQuery(opts ListTransfersOptions, from, to string) string {
	var parts []string
	add := func(k, v string) {
		if v != "" {
			parts = append(parts, k+"="+v)
		}
	}
	addAll := func(k string, vs []string) {
		for _, v := range vs {
			if v != "" {
				parts = append(parts, k+"[]="+v)
			}
		}
	}
	add("direction", opts.Direction)
	addAll("states", opts.States)
	addAll("uuids", opts.UUIDs)
	addAll("identifiers", opts.Identifiers)
	add("start_time", opts.StartTime)
	add("end_time", opts.EndTime)
	add("currency", opts.Currency)
	if opts.Limit > 0 {
		parts = append(parts, "limit="+strconv.Itoa(opts.Limit))
	}
	add("order_by", opts.OrderBy)
	add("from", from)
	add("to", to)
	return strings.Join(parts, "&")
}
