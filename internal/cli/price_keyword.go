package cli

import (
	"context"
	"fmt"

	"github.com/kyungw00k/upbit/api"
	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/internal/i18n"
)

// priceKeywords maps keyword strings to Ticker field extractors
var priceKeywords = map[string]func(t *tickerPrices) float64{
	"now":  func(t *tickerPrices) float64 { return t.tradePrice },
	"open": func(t *tickerPrices) float64 { return t.openingPrice },
	"low":  func(t *tickerPrices) float64 { return t.lowPrice },
	"high": func(t *tickerPrices) float64 { return t.highPrice },
}

type tickerPrices struct {
	tradePrice   float64
	openingPrice float64
	lowPrice     float64
	highPrice    float64
}

// isPriceKeyword checks if the given string is a price keyword
func isPriceKeyword(s string) bool {
	_, ok := priceKeywords[s]
	return ok
}

// resolvePriceKeyword resolves a price keyword (now/open/low/high) to actual price string
func resolvePriceKeyword(ctx context.Context, client *api.Client, market, keyword string) (string, error) {
	extractor, ok := priceKeywords[keyword]
	if !ok {
		return keyword, nil
	}

	qc := quotation.NewQuotationClient(client)
	tickers, err := qc.GetTickers(ctx, []string{market})
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T(i18n.ErrTickerFetch), err)
	}
	if len(tickers) == 0 {
		return "", fmt.Errorf("%s", i18n.Tf(i18n.ErrTickerEmpty, market))
	}

	t := &tickerPrices{
		tradePrice:   tickers[0].TradePrice,
		openingPrice: tickers[0].OpeningPrice,
		lowPrice:     tickers[0].LowPrice,
		highPrice:    tickers[0].HighPrice,
	}

	price := extractor(t)
	resolved := smartPrice(price)
	return resolved, nil
}
