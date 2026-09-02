package exchange

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyungw00k/upbit/types"
)

// OrderRequest holds the parameters for placing an order.
type OrderRequest struct {
	Market      string `json:"market"`
	Side        string `json:"side"`                    // bid, ask
	OrdType     string `json:"ord_type"`                // limit, price, market, best
	Volume      string `json:"volume,omitempty"`        // order volume
	Price       string `json:"price,omitempty"`         // order price or total amount
	WatchPrice  string `json:"watch_price,omitempty"`   // trigger price for reserved orders
	TimeInForce string `json:"time_in_force,omitempty"` // ioc, fok, post_only
	SMPType     string `json:"smp_type,omitempty"`      // cancel_maker, cancel_taker, reduce
	Identifier  string `json:"identifier,omitempty"`    // client-assigned order identifier
}

// CreateOrder places a new order.
// API: POST /orders
// See https://docs.upbit.com/reference/%EC%A3%BC%EB%AC%B8%ED%95%98%EA%B8%B0
func (c *ExchangeClient) CreateOrder(ctx context.Context, req *OrderRequest) (*types.Order, error) {
	var order types.Order
	err := c.client.POST(ctx, "/orders", req, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// TestOrder places a test order (not executed against the market).
// API: POST /orders/test
// See https://docs.upbit.com/reference/%EC%A3%BC%EB%AC%B8%ED%95%98%EA%B8%B0
func (c *ExchangeClient) TestOrder(ctx context.Context, req *OrderRequest) (*types.Order, error) {
	var order types.Order
	err := c.client.POST(ctx, "/orders/test", req, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrder retrieves a single order by UUID.
// API: GET /order?uuid=xxx
// See https://docs.upbit.com/reference/%EA%B0%9C%EB%B3%84-%EC%A3%BC%EB%AC%B8-%EC%A1%B0%ED%9A%8C
func (c *ExchangeClient) GetOrder(ctx context.Context, uuid string) (*types.Order, error) {
	var order types.Order
	query := map[string]string{
		"uuid": uuid,
	}
	err := c.client.GET(ctx, "/order", query, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrderByIdentifier retrieves a single order by client-assigned identifier.
// API: GET /order?identifier=xxx
// See https://docs.upbit.com/reference/%EA%B0%9C%EB%B3%84-%EC%A3%BC%EB%AC%B8-%EC%A1%B0%ED%9A%8C
func (c *ExchangeClient) GetOrderByIdentifier(ctx context.Context, identifier string) (*types.Order, error) {
	var order types.Order
	query := map[string]string{
		"identifier": identifier,
	}
	err := c.client.GET(ctx, "/order", query, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ListOpenOrders returns the list of open (waiting) orders.
// API: GET /orders/open
// See https://docs.upbit.com/reference/%EB%8C%80%EA%B8%B0-%EC%A3%BC%EB%AC%B8-%EC%A1%B0%ED%9A%8C
func (c *ExchangeClient) ListOpenOrders(ctx context.Context, market string, limit int, page int) ([]types.Order, error) {
	var orders []types.Order
	query := map[string]string{}

	if market != "" {
		query["market"] = market
	}
	if limit > 0 {
		query["limit"] = fmt.Sprintf("%d", limit)
	}
	if page > 0 {
		query["page"] = fmt.Sprintf("%d", page)
	}

	err := c.client.GET(ctx, "/orders/open", query, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// ListClosedOrders returns the list of closed (completed or cancelled) orders.
// API: GET /orders/closed
// See https://docs.upbit.com/reference/%EC%A2%85%EB%A3%8C-%EC%A3%BC%EB%AC%B8-%EC%A1%B0%ED%9A%8C
func (c *ExchangeClient) ListClosedOrders(ctx context.Context, market string, limit int, page int) ([]types.Order, error) {
	var orders []types.Order
	query := map[string]string{}

	if market != "" {
		query["market"] = market
	}
	if limit > 0 {
		query["limit"] = fmt.Sprintf("%d", limit)
	}
	if page > 0 {
		query["page"] = fmt.Sprintf("%d", page)
	}

	err := c.client.GET(ctx, "/orders/closed", query, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrdersByIDsOptions filters the batch order query (GET /orders/uuids).
// At most 100 UUIDs or 100 identifiers per request.
type GetOrdersByIDsOptions struct {
	UUIDs       []string
	Identifiers []string
	Market      string
	OrderBy     string // asc, desc (default desc)
}

// GetOrdersByIDs retrieves multiple orders by UUID or identifier list.
// API: GET /orders/uuids?uuids[]=xxx&identifiers[]=yyy&market=KRW-BTC&order_by=desc
// See https://docs.upbit.com/reference/list-orders-by-ids
func (c *ExchangeClient) GetOrdersByIDs(ctx context.Context, opts GetOrdersByIDsOptions) ([]types.Order, error) {
	var orders []types.Order

	var parts []string
	for _, id := range opts.UUIDs {
		if id != "" {
			parts = append(parts, "uuids[]="+id)
		}
	}
	for _, id := range opts.Identifiers {
		if id != "" {
			parts = append(parts, "identifiers[]="+id)
		}
	}
	if opts.Market != "" {
		parts = append(parts, "market="+opts.Market)
	}
	if opts.OrderBy != "" {
		parts = append(parts, "order_by="+opts.OrderBy)
	}
	rawQuery := strings.Join(parts, "&")

	err := c.client.GETWithRawQuery(ctx, "/orders/uuids", rawQuery, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrdersByUUIDs retrieves multiple orders by UUID list (max 100).
// API: GET /orders/uuids?uuids[]=xxx&uuids[]=yyy
func (c *ExchangeClient) GetOrdersByUUIDs(ctx context.Context, uuids []string) ([]types.Order, error) {
	return c.GetOrdersByIDs(ctx, GetOrdersByIDsOptions{UUIDs: uuids})
}
