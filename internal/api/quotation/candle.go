package quotation

import (
	"context"
	"fmt"
	"sort"

	"github.com/kyungw00k/upbit/internal/types"
)

// intervalToPath interval 문자열을 Upbit API 캔들 경로로 매핑
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

// GetCandles 캔들(OHLCV) 조회
// interval: "1s", "1m", "3m", "5m", "10m", "15m", "30m", "60m", "240m", "1d", "1w", "1M", "1y"
// API: GET /candles/{type}?market=KRW-BTC&count=200
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

// GetCandlesWithTo to 시각 이전 캔들 조회 (API 1회 호출)
// to가 비어있으면 파라미터 생략 (현재 시각 기준)
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

// GetCandlesAll 자동 페이지네이션으로 대량 캔들 조회
// from: 시작 시각 (빈 문자열이면 count만큼만)
// count: 0이면 from부터 현재까지 전부, >0이면 해당 개수만
// 결과는 시간순(asc) 정렬하여 반환
func (c *QuotationClient) GetCandlesAll(ctx context.Context, market, interval string, from string, count int) ([]types.Candle, error) {
	const pageSize = 200
	var all []types.Candle
	to := "" // 첫 호출: 현재 시각 기준

	for {
		// 남은 개수 계산
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

		// API는 최신→오래된 순으로 반환
		// from이 지정된 경우: from 이전 캔들은 필터링
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

		// count 지정 시 도달 확인
		if count > 0 && len(all) >= count {
			all = all[:count]
			break
		}

		// 다음 페이지를 위한 to 설정: 마지막(가장 오래된) 캔들의 UTC 시각
		lastCandle := batch[len(batch)-1]
		nextTo := lastCandle.CandleDateTimeUtc
		if nextTo == to {
			// 무한 루프 방지
			break
		}
		to = nextTo

		// 반환 개수가 요청보다 적으면 더 이상 데이터 없음
		if len(batch) < fetchCount {
			break
		}
	}

	// 시간순(asc) 정렬 — API는 역순으로 반환하므로
	sort.Slice(all, func(i, j int) bool {
		return all[i].CandleDateTimeKst < all[j].CandleDateTimeKst
	})

	return all, nil
}
