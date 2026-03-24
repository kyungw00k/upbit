package quotation

import (
	"context"
	"fmt"
	"sort"

	"github.com/kyungw00k/upbit/types"
)

// intervalToPath maps an interval string to the corresponding Upbit API candle path.
// "1s"  → /candles/seconds
// "1m"  → /candles/minutes/1
// "3m"  → /candles/minutes/3
// "5m"  → /candles/minutes/5
// "10m" → /candles/minutes/10
// "15m" → /candles/minutes/15
// "30m" → /candles/minutes/30
// "60m" → /candles/minutes/60
// "240m"→ /candles/minutes/240
// "1d"  → /candles/days
// "1w"  → /candles/weeks
// "1M"  → /candles/months
// "1y"  → /candles/years
func intervalToPath(interval string) (string, error) {
	switch interval {
	case "1s":
		return "/candles/seconds", nil
	case "1m":
		return "/candles/minutes/1", nil
	case "3m":
		return "/candles/minutes/3", nil
	case "5m":
		return "/candles/minutes/5", nil
	case "10m":
		return "/candles/minutes/10", nil
	case "15m":
		return "/candles/minutes/15", nil
	case "30m":
		return "/candles/minutes/30", nil
	case "60m":
		return "/candles/minutes/60", nil
	case "240m":
		return "/candles/minutes/240", nil
	case "1d":
		return "/candles/days", nil
	case "1w":
		return "/candles/weeks", nil
	case "1M":
		return "/candles/months", nil
	case "1y":
		return "/candles/years", nil
	default:
		return "", fmt.Errorf("지원하지 않는 캔들 간격: %s (지원: 1s, 1m, 3m, 5m, 10m, 15m, 30m, 60m, 240m, 1d, 1w, 1M, 1y)", interval)
	}
}

// GetCandles retrieves OHLCV candles for the specified market and interval.
// interval: "1s", "1m", "3m", "5m", "10m", "15m", "30m", "60m", "240m", "1d", "1w", "1M", "1y"
// API: GET /candles/{type}?market=KRW-BTC&count=200
// See https://docs.upbit.com/reference/%EB%B6%84minute-%EC%BA%94%EB%93%A4-1
func (c *QuotationClient) GetCandles(ctx context.Context, market string, interval string, count int) ([]types.Candle, error) {
	path, err := intervalToPath(interval)
	if err != nil {
		return nil, err
	}

	var candles []types.Candle
	query := map[string]string{
		"market": market,
		"count":  fmt.Sprintf("%d", count),
	}
	err = c.client.GET(ctx, path, query, &candles)
	if err != nil {
		return nil, err
	}
	return candles, nil
}

// GetCandlesWithTo retrieves candles before the specified time (single API call).
// If to is empty, the parameter is omitted and the current time is used as the reference.
func (c *QuotationClient) GetCandlesWithTo(ctx context.Context, market, interval string, count int, to string) ([]types.Candle, error) {
	path, err := intervalToPath(interval)
	if err != nil {
		return nil, err
	}

	var candles []types.Candle
	query := map[string]string{
		"market": market,
		"count":  fmt.Sprintf("%d", count),
	}
	if to != "" {
		query["to"] = to
	}
	err = c.client.GET(ctx, path, query, &candles)
	if err != nil {
		return nil, err
	}
	return candles, nil
}

// GetCandlesAll retrieves a large number of candles using automatic pagination.
// from: start time (if empty, only count candles are fetched)
// count: if 0, fetches all candles from 'from' to now; if >0, fetches that many candles
// Results are returned sorted in ascending chronological order.
func (c *QuotationClient) GetCandlesAll(ctx context.Context, market, interval string, from string, count int) ([]types.Candle, error) {
	const pageSize = 200
	var all []types.Candle
	to := "" // first call: use current time as reference

	for {
		// calculate remaining count
		fetchCount := pageSize
		if count > 0 {
			remaining := count - len(all)
			if remaining <= 0 {
				break
			}
			if remaining < pageSize {
				fetchCount = remaining
			}
		}

		batch, err := c.GetCandlesWithTo(ctx, market, interval, fetchCount, to)
		if err != nil {
			return nil, fmt.Errorf("캔들 조회 실패: %w", err)
		}

		if len(batch) == 0 {
			break
		}

		// API returns candles in descending order (newest first)
		// if from is specified, filter out candles older than from
		if from != "" {
			filtered := make([]types.Candle, 0, len(batch))
			reachedFrom := false
			for _, candle := range batch {
				if candle.CandleDateTimeKst < from {
					reachedFrom = true
					continue
				}
				filtered = append(filtered, candle)
			}
			batch = filtered
			if reachedFrom {
				all = append(all, batch...)
				break
			}
		}

		all = append(all, batch...)

		// check if requested count has been reached
		if count > 0 && len(all) >= count {
			all = all[:count]
			break
		}

		// set to for the next page: UTC time of the last (oldest) candle in the batch
		lastCandle := batch[len(batch)-1]
		nextTo := lastCandle.CandleDateTimeUtc
		if nextTo == to {
			// prevent infinite loop
			break
		}
		to = nextTo

		// if fewer candles than requested were returned, no more data is available
		if len(batch) < fetchCount {
			break
		}
	}

	// sort in ascending chronological order — API returns in descending order
	sort.Slice(all, func(i, j int) bool {
		return all[i].CandleDateTimeKst < all[j].CandleDateTimeKst
	})

	return all, nil
}
