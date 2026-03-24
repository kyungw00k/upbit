package wallet

import (
	"context"

	"github.com/kyungw00k/upbit/types"
)

// TravelRuleVerifyByTxIDRequest holds parameters for a Travel Rule verification request by transaction ID.
type TravelRuleVerifyByTxIDRequest struct {
	TxID     string `json:"txid"`
	VaspUUID string `json:"vasp_uuid"`
	Currency string `json:"currency"`
	NetType  string `json:"net_type"`
}

// TravelRuleVerifyByUUIDRequest holds parameters for a Travel Rule verification request by deposit UUID.
type TravelRuleVerifyByUUIDRequest struct {
	DepositUUID string `json:"deposit_uuid"`
	VaspUUID    string `json:"vasp_uuid"`
}

// GetTravelRuleVASPs returns a list of VASPs that support Travel Rule.
// API: GET /travel_rule/vasps
// See https://docs.upbit.com/reference/%ED%8A%B8%EB%9E%98%EB%B8%94%EB%A3%B0
func (c *WalletClient) GetTravelRuleVASPs(ctx context.Context) ([]types.VASP, error) {
	var vasps []types.VASP
	err := c.client.GET(ctx, "/travel_rule/vasps", nil, &vasps)
	if err != nil {
		return nil, err
	}
	return vasps, nil
}

// VerifyTravelRuleByTxID submits a Travel Rule verification request using a transaction ID.
// API: POST /travel_rule/deposit/txid
// See https://docs.upbit.com/reference/%ED%8A%B8%EB%9E%98%EB%B8%94%EB%A3%B0
func (c *WalletClient) VerifyTravelRuleByTxID(ctx context.Context, req *TravelRuleVerifyByTxIDRequest) (*types.TravelRuleVerification, error) {
	var result types.TravelRuleVerification
	err := c.client.POST(ctx, "/travel_rule/deposit/txid", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyTravelRuleByUUID submits a Travel Rule verification request using a deposit UUID.
// API: POST /travel_rule/deposit/uuid
// See https://docs.upbit.com/reference/%ED%8A%B8%EB%9E%98%EB%B8%94%EB%A3%B0
func (c *WalletClient) VerifyTravelRuleByUUID(ctx context.Context, req *TravelRuleVerifyByUUIDRequest) (*types.TravelRuleVerification, error) {
	var result types.TravelRuleVerification
	err := c.client.POST(ctx, "/travel_rule/deposit/uuid", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
