package quotation

import (
	"context"
	"strings"

	"github.com/kyungw00k/upbit/types"
)

// GetTickers retrieves the current price for the specified markets.
// API: GET /ticker?markets=KRW-BTC,KRW-ETH
// See https://docs.upbit.com/reference/%ED%98%84%EC%9E%AC%EA%B0%80-%EC%A0%95%EB%B3%B4
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

// GetAllTickers retrieves the current price for all markets by quote currency.
// API: GET /ticker/all?quote_currencies=KRW,BTC,USDT
// See https://docs.upbit.com/reference/%ED%98%84%EC%9E%AC%EA%B0%80-%EC%A0%95%EB%B3%B4
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
