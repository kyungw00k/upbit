package quotation

import (
	"context"
	"strings"

	"github.com/kyungw00k/upbit/types"
)

// GetTickers 현재가 조회
// API: GET /ticker?markets=KRW-BTC,KRW-ETH
func (c *QuotationClient) GetTickers(ctx context.Context, markets []string) ([]types.Ticker, error) {
	var tickers []types.Ticker
	query := map[string]string{
		"markets": strings.Join(markets, ","),
	}
	err := c.client.GET(ctx, "/ticker", query, &tickers)
	if err != nil {
		return nil, err
	}
	return tickers, nil
}

// GetAllTickers 마켓 단위 전체 현재가 조회
// API: GET /ticker/all?quote_currencies=KRW,BTC,USDT
func (c *QuotationClient) GetAllTickers(ctx context.Context, quoteCurrencies []string) ([]types.Ticker, error) {
	var tickers []types.Ticker
	query := map[string]string{
		"quote_currencies": strings.Join(quoteCurrencies, ","),
	}
	err := c.client.GET(ctx, "/ticker/all", query, &tickers)
	if err != nil {
		return nil, err
	}
	return tickers, nil
}
