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
	market := envOr("MARKET", "KRW-BTC")
	dbPath := envOr("DB", "data/candles.db")

	db, err := storage.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB open failed: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	intervals := []struct {
		code     string
		from     string
		to       string
		maxCount int
	}{
		{"30", "2024-01-01T00:00:00", "2025-12-31T23:59:59", 0},
		{"60", "2024-01-01T00:00:00", "2025-12-31T23:59:59", 0},
		{"240", "2024-01-01T00:00:00", "2025-12-31T23:59:59", 0},
		{"1d", "2024-01-01T00:00:00", "2025-12-31T23:59:59", 0},
	}

	apiClient := api.NewClient("", "")
	qClient := quotation.NewQuotationClient(apiClient)

	for _, iv := range intervals {
		label := iv.code + "m"
		if iv.code == "day" {
			label = "day"
		}

		existing, _ := storage.CandleCount(db, market, iv.code)
		fmt.Printf("[%s] existing: %d candles\n", label, existing)

		if existing > 0 {
			min, max, _ := storage.TimeRange(db, market, iv.code)
			fmt.Printf("[%s] range: %s ~ %s (skip)\n", label, min, max)
			continue
		}

		fmt.Printf("[%s] downloading %s ~ %s ...\n", label, iv.from, iv.to)

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		raw, err := qClient.GetCandlesAll(ctx, market, label, iv.from, iv.maxCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] download failed: %v\n", label, err)
			continue
		}

		// Filter by to
		filtered := raw
		if iv.to != "" {
			f := make([]types.Candle, 0, len(raw))
			for _, r := range raw {
				if r.CandleDateTimeKst <= iv.to {
					f = append(f, r)
				}
			}
			filtered = f
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

		inserted, err := storage.InsertCandles(db, market, iv.code, candles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] insert failed: %v\n", label, err)
			continue
		}
		fmt.Printf("[%s] stored %d / %d candles\n\n", label, inserted, len(candles))
	}

	// Summary
	fmt.Println("\n=== Storage Summary ===")
	for _, iv := range intervals {
		label := iv.code + "m"
		if iv.code == "day" {
			label = "day"
		}
		count, _ := storage.CandleCount(db, market, iv.code)
		min, max, _ := storage.TimeRange(db, market, iv.code)
		fmt.Printf("%-6s %6d candles  %s ~ %s\n", label, count, min, max)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
