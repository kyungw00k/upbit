package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kyungw00k/upbit/trader/backtest"
	"github.com/kyungw00k/upbit/trader/indicator"
)

// upbitCandle matches the Upbit API candle response.
type upbitCandle struct {
	Market               string  `json:"market"`
	CandleDateTimeKst    string  `json:"candle_date_time_kst"`
	OpeningPrice         float64 `json:"opening_price"`
	HighPrice            float64 `json:"high_price"`
	LowPrice             float64 `json:"low_price"`
	TradePrice           float64 `json:"trade_price"`
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
}

const (
	defaultMarket       = "KRW-BTC"
	defaultDays         = 200
	defaultInitBalance  = 10_000_000 // 1000만원
	upbitCandleURL      = "https://api.upbit.com/v1/candles/days"
)

func main() {
	market := envOr("MARKET", defaultMarket)
	days, _ := strconv.Atoi(envOr("DAYS", strconv.Itoa(defaultDays)))
	initBalance := defaultInitBalance

	fmt.Printf("Fetching %d daily candles for %s...\n", days, market)

	candles, err := fetchCandles(market, days)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Got %d candles (%s ~ %s)\n\n",
		len(candles), candles[0].Time, candles[len(candles)-1].Time)

	result, err := backtest.RunBacktest(context.Background(), market, candles, float64(initBalance))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Backtest error: %v\n", err)
		os.Exit(1)
	}

	backtest.PrintResult(market, result)
}

func fetchCandles(market string, count int) ([]indicator.CandleData, error) {
	url := fmt.Sprintf("%s?market=%s&count=%d", upbitCandleURL, market, count)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var raw []upbitCandle
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}

	// API returns newest first, reverse to chronological order
	result := make([]indicator.CandleData, len(raw))
	for i, j := 0, len(raw)-1; j >= 0; i, j = i+1, j-1 {
		r := raw[j]
		result[i] = indicator.CandleData{
			Open:   r.OpeningPrice,
			High:   r.HighPrice,
			Low:    r.LowPrice,
			Close:  r.TradePrice,
			Volume: r.CandleAccTradeVolume,
			Time:   r.CandleDateTimeKst,
		}
	}

	return result, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
