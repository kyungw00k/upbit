package wallet

import (
	"context"

	"github.com/kyungw00k/upbit/internal/types"
)

// TravelRuleVerifyByTxIDRequest TxID 기반 트래블룰 검증 요청 파라미터
type TravelRuleVerifyByTxIDRequest struct {
	TxID     string `json:"txid"`
	VaspUUID string `json:"vasp_uuid"`
	Currency string `json:"currency"`
	NetType  string `json:"net_type"`
}

// TravelRuleVerifyByUUIDRequest UUID 기반 트래블룰 검증 요청 파라미터
type TravelRuleVerifyByUUIDRequest struct {
	DepositUUID string `json:"deposit_uuid"`
	VaspUUID    string `json:"vasp_uuid"`
}

// GetTravelRuleVASPs 트래블룰 지원 거래소 목록 조회
// API: GET /travel_rule/vasps
func (c *WalletClient) GetTravelRuleVASPs(ctx context.Context) ([]types.VASP, error) {
	var vasps []types.VASP
	err := c.client.GET(ctx, "/travel_rule/vasps", nil, &vasps)
	if err != nil {
		return nil, err
	}
	return vasps, nil
}

// VerifyTravelRuleByTxID TxID 기반 트래블룰 검증 요청
// API: POST /travel_rule/deposit/txid
func (c *WalletClient) VerifyTravelRuleByTxID(ctx context.Context, req *TravelRuleVerifyByTxIDRequest) (*types.TravelRuleVerification, error) {
	var result types.TravelRuleVerification
	err := c.client.POST(ctx, "/travel_rule/deposit/txid", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyTravelRuleByUUID UUID 기반 트래블룰 검증 요청
// API: POST /travel_rule/deposit/uuid
func (c *WalletClient) VerifyTravelRuleByUUID(ctx context.Context, req *TravelRuleVerifyByUUIDRequest) (*types.TravelRuleVerification, error) {
	var result types.TravelRuleVerification
	err := c.client.POST(ctx, "/travel_rule/deposit/uuid", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
