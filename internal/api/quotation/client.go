package quotation

import (
	"github.com/kyungw00k/upbit/internal/api"
)

// QuotationClient Quotation API (시세 조회) 클라이언트
type QuotationClient struct {
	client *api.Client
}

// NewQuotationClient QuotationClient 생성
func NewQuotationClient(client *api.Client) *QuotationClient {
	return &QuotationClient{client: client}
}
