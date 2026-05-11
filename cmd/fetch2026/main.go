package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kyungw00k/upbit/api"
	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/trader/indicator"
	"github.com/kyungw00k/upbit/trader/storage"
	"github.com/kyungw00k/upbit/types"
)

func main() {
	market := "KRW-BTC"
	dbPath := "data/candles.db"

	db, err := storage.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB open: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	apiClient := api.NewClient("", "")
	qClient := quotation.NewQuotationClient(apiClient)

	for _, iv := range []string{"30m", "60m", "240m"} {
		fmt.Printf("[%s] downloading 2026-01-01 ~ now ...\n", iv)
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		raw, err := qClient.GetCandlesAll(ctx, market, iv, "2026-01-01T00:00:00", 0)
		cancel()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] failed: %v\n", iv, err)
			continue
		}

		filtered := make([]types.Candle, 0)
		for _, r := range raw {
			if r.CandleDateTimeKst >= "2026-01-01T00:00:00" {
				filtered = append(filtered, r)
			}
		}

		candles := make([]indicator.CandleData, len(filtered))
		for i, r := range filtered {
			candles[i] = indicator.CandleData{
				Open:   r.OpeningPrice,
				High:   r.HighPrice,
				Low:    r.LowPrice,
				Close:  r.TradePrice,
				Volume: r.CandleAccTradeVolume,
				Time:   r.CandleDateTimeKst,
			}
		}

		inserted, err := storage.InsertCandles(db, market, iv[:len(iv)-1], candles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] insert failed: %v\n", iv, err)
			continue
		}
		fmt.Printf("[%s] stored %d candles (2026)\n", iv, inserted)
	}

	fmt.Println("\n=== Updated Summary ===")
	for _, iv := range []string{"30", "60", "240"} {
		count, _ := storage.CandleCount(db, market, iv)
		min, max, _ := storage.TimeRange(db, market, iv)
		fmt.Printf("%-6s %6d candles  %s ~ %s\n", iv+"m", count, min, max)
	}
}
