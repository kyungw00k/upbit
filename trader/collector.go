package trader

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/kyungw00k/upbit/api"
	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/trader/indicator"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	upbitTypes "github.com/kyungw00k/upbit/types"
)

// CollectAndCalculate fetches candle data for each market, computes technical
// indicators using CalculateAll, and upserts the results into the PocketBase
// "indicators" collection.
func CollectAndCalculate(app core.App, apiClient *api.Client, markets []string, interval string) error {
	qc := quotation.NewQuotationClient(apiClient)

	for _, market := range markets {
		if err := collectOne(app, qc, market, interval); err != nil {
			log.Printf("[collector] %s/%s 수집 실패: %v", market, interval, err)
			continue
		}
	}

	return nil
}

func collectOne(app core.App, qc *quotation.QuotationClient, market, interval string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Fetch 500 candles for EMA200 stability.
	candles, err := qc.GetCandlesAll(ctx, market, interval, "", 500)
	if err != nil {
		return fmt.Errorf("캔들 조회 실패: %w", err)
	}

	if len(candles) < 50 {
		return fmt.Errorf("캔들 데이터 부족: %d개 (최소 50개 필요)", len(candles))
	}

	// Convert to indicator CandleData.
	ic := make([]indicator.CandleData, len(candles))
	for i, c := range candles {
		ic[i] = indicator.CandleData{
			Open:   c.OpeningPrice,
			High:   c.HighPrice,
			Low:    c.LowPrice,
			Close:  c.TradePrice,
			Volume: c.CandleAccTradeVolume,
			Time:   c.CandleDateTimeKst,
		}
	}

	// Calculate all indicators.
	result := indicator.CalculateAll(ic)

	// Use the last candle's timestamp for the record.
	lastCandle := candles[len(candles)-1]

	// Upsert into indicators collection.
	return upsertIndicator(app, market, interval, lastCandle, result)
}

func upsertIndicator(app core.App, market, interval string, candle upbitTypes.Candle, r indicator.Result) error {
	col, err := app.FindCollectionByNameOrId("indicators")
	if err != nil {
		return fmt.Errorf("indicators 컬렉션 조회 실패: %w", err)
	}

	// Try to find existing record by market + interval + candle_time.
	filter := "market={:market} && interval={:interval} && candle_time={:candleTime}"
	existing, findErr := app.FindFirstRecordByFilter(col, filter, dbx.Params{
		"market":    market,
		"interval":  interval,
		"candleTime": candle.CandleDateTimeKst,
	})

	var record *core.Record
	if findErr == nil && existing != nil {
		record = existing
	} else {
		record = core.NewRecord(col)
		record.Set("market", market)
		record.Set("interval", interval)
		record.Set("candle_time", candle.CandleDateTimeKst)
	}

	record.Set("open", candle.OpeningPrice)
	record.Set("high", candle.HighPrice)
	record.Set("low", candle.LowPrice)
	record.Set("close", candle.TradePrice)
	record.Set("volume", candle.CandleAccTradeVolume)

	if !math.IsNaN(r.RSI) {
		record.Set("rsi_14", r.RSI)
	}
	if !math.IsNaN(r.MACDLine) {
		record.Set("macd_line", r.MACDLine)
	}
	if !math.IsNaN(r.MACDSignal) {
		record.Set("macd_signal", r.MACDSignal)
	}
	if !math.IsNaN(r.MACDHist) {
		record.Set("macd_hist", r.MACDHist)
	}
	if !math.IsNaN(r.BBUpper) {
		record.Set("bb_upper", r.BBUpper)
	}
	if !math.IsNaN(r.BBMiddle) {
		record.Set("bb_middle", r.BBMiddle)
	}
	if !math.IsNaN(r.BBLower) {
		record.Set("bb_lower", r.BBLower)
	}
	if !math.IsNaN(r.EMA20) {
		record.Set("ema_20", r.EMA20)
	}
	if !math.IsNaN(r.EMA50) {
		record.Set("ema_50", r.EMA50)
	}
	if !math.IsNaN(r.ATR) {
		record.Set("atr_14", r.ATR)
	}
	if r.VolRatio > 0 {
		record.Set("vol_ratio", r.VolRatio)
	}
	record.Set("trend", r.Trend)

	if err := app.Save(record); err != nil {
		return fmt.Errorf("indicators 레코드 저장 실패: %w", err)
	}

	return nil
}
