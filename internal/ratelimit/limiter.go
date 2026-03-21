package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Group Rate limit 그룹 타입
type Group string

const (
	GroupMarket         Group = "market"
	GroupCandle         Group = "candle"
	GroupTicker         Group = "ticker"
	GroupOrderbook      Group = "orderbook"
	GroupTrade          Group = "trade"
	GroupDefault        Group = "default"
	GroupOrder          Group = "order"
	GroupOrderTest      Group = "order-test"
	GroupOrderCancelAll Group = "order-cancel-all"
)

// 그룹별 초당 요청 수
const (
	rateMarket         = 10.0
	rateCandle         = 10.0
	rateTicker         = 10.0
	rateOrderbook      = 10.0
	rateTrade          = 10.0
	rateDefault        = 30.0
	rateOrder          = 8.0
	rateOrderTest      = 8.0
	rateOrderCancelAll = 0.5
)

// bucket Token bucket 구현체
type bucket struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	rate     float64 // 초당 토큰 보충량
	lastTime time.Time
}

func newBucket(rate float64) *bucket {
	return &bucket{
		tokens:   rate, // 초기에 가득 채움
		capacity: rate,
		rate:     rate,
		lastTime: time.Now(),
	}
}

// wait 토큰 소비 또는 대기
func (b *bucket) wait(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastTime).Seconds()
	b.tokens = min(b.capacity, b.tokens+elapsed*b.rate)
	b.lastTime = now

	if b.tokens >= 1 {
		b.tokens--
		return nil
	}

	// 대기 시간 계산
	waitDuration := time.Duration((1-b.tokens)/b.rate*1000) * time.Millisecond

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(waitDuration):
		b.tokens = 0
		b.lastTime = time.Now()
		return nil
	}
}

// setRate Rate 동적 변경
func (b *bucket) setRate(rate float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if rate > 0 {
		b.rate = rate
		b.capacity = rate
	}
}

// Limiter Rate limiter (그룹별 Token bucket)
type Limiter struct {
	buckets map[Group]*bucket
}

// NewLimiter 새 Rate limiter 생성
func NewLimiter() *Limiter {
	return &Limiter{
		buckets: map[Group]*bucket{
			GroupMarket:         newBucket(rateMarket),
			GroupCandle:         newBucket(rateCandle),
			GroupTicker:         newBucket(rateTicker),
			GroupOrderbook:      newBucket(rateOrderbook),
			GroupTrade:          newBucket(rateTrade),
			GroupDefault:        newBucket(rateDefault),
			GroupOrder:          newBucket(rateOrder),
			GroupOrderTest:      newBucket(rateOrderTest),
			GroupOrderCancelAll: newBucket(rateOrderCancelAll),
		},
	}
}

// Wait 해당 그룹의 Rate limit 대기
func (l *Limiter) Wait(ctx context.Context, group Group) error {
	b, ok := l.buckets[group]
	if !ok {
		b = l.buckets[GroupDefault]
	}
	return b.wait(ctx)
}

// UpdateFromHeader Remaining-Req 헤더를 파싱하여 Rate limiter 동적 조절
// 헤더 형식: "group=market; min=599; sec=9"
func (l *Limiter) UpdateFromHeader(header string, group Group) {
	if header == "" {
		return
	}

	remaining, perSec := parseRemainingReq(header)

	b, ok := l.buckets[group]
	if !ok {
		return
	}

	// remaining이 임계값 이하이면 속도 제한
	if remaining >= 0 && remaining <= 2 {
		b.setRate(1.0) // 긴급 제한: 초당 1회
		return
	}

	// 서버가 알려준 남은 초당 허용량으로 동적 조절
	if perSec > 0 {
		b.setRate(float64(perSec))
	}
}

// parseRemainingReq Remaining-Req 헤더 파싱
// 헤더 형식: "group=market; min=599; sec=9"
// 반환: (remaining, perSec) — sec 값을 remaining 및 perSec 양쪽에 사용
func parseRemainingReq(header string) (int, int) {
	remaining := -1
	perSec := -1

	parts := strings.Split(header, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		switch key {
		case "sec":
			if n, err := strconv.Atoi(val); err == nil {
				remaining = n
				perSec = n
			}
		case "min":
			// 분 단위는 참고용으로만 사용
			_ = val
		}
	}

	return remaining, perSec
}

// GroupFromPath HTTP 경로에서 Rate limit 그룹 자동 판별
// 실제 Upbit API 엔드포인트 경로 기준
func GroupFromPath(path string) Group {
	switch {
	case strings.HasPrefix(path, "/trading_pairs"):
		return GroupMarket
	case strings.HasPrefix(path, "/tickers") || strings.HasPrefix(path, "/ticker"):
		return GroupTicker
	case strings.HasPrefix(path, "/candles/"):
		return GroupCandle
	case strings.HasPrefix(path, "/orderbooks"):
		return GroupOrderbook
	case strings.HasPrefix(path, "/trades/"):
		return GroupTrade
	case strings.HasPrefix(path, "/orders/test"):
		return GroupOrderTest
	case strings.HasPrefix(path, "/orders/batch"):
		return GroupOrderCancelAll
	case strings.HasPrefix(path, "/orders"):
		return GroupOrder
	default:
		return GroupDefault
	}
}

// String 그룹명 문자열 반환
func (g Group) String() string {
	return string(g)
}

// Validate 그룹 유효성 검사
func (g Group) Validate() error {
	switch g {
	case GroupMarket, GroupCandle, GroupTicker, GroupOrderbook, GroupTrade,
		GroupDefault, GroupOrder, GroupOrderTest, GroupOrderCancelAll:
		return nil
	default:
		return fmt.Errorf("알 수 없는 rate limit 그룹: %s", g)
	}
}
