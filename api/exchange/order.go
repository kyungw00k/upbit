package exchange

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyungw00k/upbit/types"
)

// OrderRequest 주문 생성 요청 파라미터
type OrderRequest struct {
	Market      string `json:"market"`
	Side        string `json:"side"`                    // bid, ask
	OrdType     string `json:"ord_type"`                // limit, price, market, best
	Volume      string `json:"volume,omitempty"`         // 주문 수량
	Price       string `json:"price,omitempty"`          // 주문 단가 또는 총액
	TimeInForce string `json:"time_in_force,omitempty"` // ioc, fok, post_only
	SMPType     string `json:"smp_type,omitempty"`      // cancel_maker, cancel_taker, reduce
	Identifier  string `json:"identifier,omitempty"`    // 클라이언트 지정 주문 식별자
}

// CreateOrder 주문 생성
// API: POST /orders
func (c *ExchangeClient) CreateOrder(ctx context.Context, req *OrderRequest) (*types.Order, error) {
	var order types.Order
	err := c.client.POST(ctx, "/orders", req, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// TestOrder 테스트 주문 (실제 체결 안됨)
// API: POST /orders/test
func (c *ExchangeClient) TestOrder(ctx context.Context, req *OrderRequest) (*types.Order, error) {
	var order types.Order
	err := c.client.POST(ctx, "/orders/test", req, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrder UUID로 개별 주문 조회
// API: GET /order?uuid=xxx
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

// GetOrderByIdentifier Identifier로 개별 주문 조회
// API: GET /order?identifier=xxx
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

// ListOpenOrders 체결 대기 주문 목록 조회
// API: GET /orders/open
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

// ListClosedOrders 종료 주문 목록 조회
// API: GET /orders/closed
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

// GetOrdersByUUIDs UUID 목록으로 복수 주문 조회 (최대 100개)
// API: GET /orders/uuids?uuids[]=xxx&uuids[]=yyy
func (c *ExchangeClient) GetOrdersByUUIDs(ctx context.Context, uuids []string) ([]types.Order, error) {
	var orders []types.Order

	// 배열 쿼리 문자열 생성: uuids[]=uuid1&uuids[]=uuid2
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
