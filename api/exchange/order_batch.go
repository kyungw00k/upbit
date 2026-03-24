package exchange

import (
	"context"
	"strings"

	"github.com/kyungw00k/upbit/types"
)

// CancelAndNewOrderRequest holds the parameters for cancelling an order and placing a replacement.
type CancelAndNewOrderRequest struct {
	PrevOrderUUID       string `json:"prev_order_uuid,omitempty"`       // UUID of the order to cancel
	PrevOrderIdentifier string `json:"prev_order_identifier,omitempty"` // identifier of the order to cancel
	NewOrdType          string `json:"new_ord_type"`                    // limit, price, market, best
	NewVolume           string `json:"new_volume,omitempty"`            // new order volume (or "remain_only")
	NewPrice            string `json:"new_price,omitempty"`             // new order price or total amount
	NewIdentifier       string `json:"new_identifier,omitempty"`        // new order identifier
	NewTimeInForce      string `json:"new_time_in_force,omitempty"`    // ioc, fok, post_only
	NewSMPType          string `json:"new_smp_type,omitempty"`         // cancel_maker, cancel_taker, reduce
}

// CancelOrder cancels a single order by UUID.
// API: DELETE /order?uuid=xxx
// See https://docs.upbit.com/reference/%EC%A3%BC%EB%AC%B8-%EC%B7%A8%EC%86%8C
func (c *ExchangeClient) CancelOrder(ctx context.Context, uuid string) (*types.Order, error) {
	var order types.Order
	query := map[string]string{
		"uuid": uuid,
	}
	err := c.client.DELETE(ctx, "/order", query, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CancelOrdersByIDs cancels multiple orders by UUID list (max 20).
// API: DELETE /orders/uuids?uuids[]=xxx&uuids[]=yyy
// Only query parameter format is supported (request body not allowed).
// See https://docs.upbit.com/reference/%EC%A3%BC%EB%AC%B8-%EC%B7%A8%EC%86%8C
func (c *ExchangeClient) CancelOrdersByIDs(ctx context.Context, uuids []string) (*types.BatchCancelResult, error) {
	var result types.BatchCancelResult

	// build array query string: uuids[]=uuid1&uuids[]=uuid2
	parts := make([]string, len(uuids))
	for i, id := range uuids {
		parts[i] = "uuids[]=" + id
	}
	rawQuery := strings.Join(parts, "&")

	err := c.client.DELETEWithRawQuery(ctx, "/orders/uuids", rawQuery, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchCancelOrders cancels open orders in bulk (max 300, WAIT status only).
// API: DELETE /orders/open
// cancel_side: all, bid, ask (default: all)
// Scope can be limited by pairs or quote_currencies.
// See https://docs.upbit.com/reference/%EC%A3%BC%EB%AC%B8-%EC%B7%A8%EC%86%8C
func (c *ExchangeClient) BatchCancelOrders(ctx context.Context, cancelSide string, pairs string, quoteCurrencies string) (*types.BatchCancelResult, error) {
	var result types.BatchCancelResult
	query := map[string]string{}

	if cancelSide != "" {
		query["cancel_side"] = cancelSide
	}
	if pairs != "" {
		query["pairs"] = pairs
	}
	if quoteCurrencies != "" {
		query["quote_currencies"] = quoteCurrencies
	}

	err := c.client.DELETE(ctx, "/orders/open", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelAndNewOrder cancels an existing order and places a replacement order.
// API: POST /orders/cancel_and_new
// See https://docs.upbit.com/reference/%EC%A3%BC%EB%AC%B8-%EC%A0%95%EC%A0%95
func (c *ExchangeClient) CancelAndNewOrder(ctx context.Context, req *CancelAndNewOrderRequest) (*types.CancelAndNewOrderResult, error) {
	var result types.CancelAndNewOrderResult
	err := c.client.POST(ctx, "/orders/cancel_and_new", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
