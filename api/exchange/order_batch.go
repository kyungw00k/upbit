package exchange

import (
	"context"
	"strings"

	"github.com/kyungw00k/upbit/types"
)

// CancelAndNewOrderRequest 취소 후 재주문 요청 파라미터
type CancelAndNewOrderRequest struct {
	PrevOrderUUID       string `json:"prev_order_uuid,omitempty"`       // 취소할 주문 UUID
	PrevOrderIdentifier string `json:"prev_order_identifier,omitempty"` // 취소할 주문 Identifier
	NewOrdType          string `json:"new_ord_type"`                    // limit, price, market, best
	NewVolume           string `json:"new_volume,omitempty"`            // 신규 주문 수량 (또는 "remain_only")
	NewPrice            string `json:"new_price,omitempty"`             // 신규 주문 단가 또는 총액
	NewIdentifier       string `json:"new_identifier,omitempty"`       // 신규 주문 식별자
	NewTimeInForce      string `json:"new_time_in_force,omitempty"`    // ioc, fok, post_only
	NewSMPType          string `json:"new_smp_type,omitempty"`         // cancel_maker, cancel_taker, reduce
}

// CancelOrder 개별 주문 취소
// API: DELETE /order?uuid=xxx
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

// CancelOrdersByIDs UUID 목록으로 주문 취소 (최대 20개)
// API: DELETE /orders/uuids?uuids[]=xxx&uuids[]=yyy
// 쿼리 파라미터 형식만 지원 (body 사용 불가)
func (c *ExchangeClient) CancelOrdersByIDs(ctx context.Context, uuids []string) (*types.BatchCancelResult, error) {
	var result types.BatchCancelResult

	// 배열 쿼리 문자열 생성: uuids[]=uuid1&uuids[]=uuid2
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

// BatchCancelOrders 주문 일괄 취소 (최대 300개, WAIT 상태만)
// API: DELETE /orders/open
// cancel_side: all, bid, ask (기본값 all)
// pairs 또는 quote_currencies로 대상 한정 가능
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

// CancelAndNewOrder 취소 후 재주문
// API: POST /orders/cancel_and_new
func (c *ExchangeClient) CancelAndNewOrder(ctx context.Context, req *CancelAndNewOrderRequest) (*types.CancelAndNewOrderResult, error) {
	var result types.CancelAndNewOrderResult
	err := c.client.POST(ctx, "/orders/cancel_and_new", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
