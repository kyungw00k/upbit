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
	Volume      string `json:"volume,omitempty"`         // order volume
	Price       string `json:"price,omitempty"`          // order price or total amount
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

// GetOrdersByUUIDs retrieves multiple orders by UUID list (max 100).
// API: GET /orders/uuids?uuids[]=xxx&uuids[]=yyy
// See https://docs.upbit.com/reference/id%EB%A1%9C-%EC%A3%BC%EB%AC%B8%EC%A1%B0%ED%9A%8C
func (c *ExchangeClient) GetOrdersByUUIDs(ctx context.Context, uuids []string) ([]types.Order, error) {
	var orders []types.Order

	// build array query string: uuids[]=uuid1&uuids[]=uuid2
	parts := make([]string, len(uuids))
	for i, id := range uuids {
		parts[i] = "uuids[]=" + id
	}
	rawQuery := strings.Join(parts, "&")

	err := c.client.GETWithRawQuery(ctx, "/orders/uuids", rawQuery, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
