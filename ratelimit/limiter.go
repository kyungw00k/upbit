// Package ratelimit implements request rate limiting for the Upbit API.
// See https://docs.upbit.com/reference/%EC%9A%94%EC%B2%AD-%EC%88%98-%EC%A0%9C%ED%95%9C for rate limit documentation.
package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Group is the rate limit group type.
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

// Requests per second per group.
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

// bucket is a token bucket implementation.
type bucket struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	rate     float64 // tokens replenished per second
	lastTime time.Time
}

func newBucket(rate float64) *bucket {
	return &bucket{
		tokens:   rate, // start full
		capacity: rate,
		rate:     rate,
		lastTime: time.Now(),
	}
}

// wait consumes a token or blocks until one is available.
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

	// Calculate how long to wait for the next token.
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

// setRate dynamically updates the bucket rate.
func (b *bucket) setRate(rate float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if rate > 0 {
		b.rate = rate
		b.capacity = rate
	}
}

// Limiter is a per-group token bucket rate limiter.
type Limiter struct {
	buckets map[Group]*bucket
}

// NewLimiter creates a new Limiter with default per-group rates.
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

// Wait blocks until a request token is available for the given group.
func (l *Limiter) Wait(ctx context.Context, group Group) error {
	b, ok := l.buckets[group]
	if !ok {
		b = l.buckets[GroupDefault]
	}
	return b.wait(ctx)
}

// UpdateFromHeader parses a Remaining-Req header and dynamically adjusts the rate limiter.
// Header format: "group=market; min=599; sec=9"
func (l *Limiter) UpdateFromHeader(header string, group Group) {
	if header == "" {
		return
	}

	remaining, perSec := parseRemainingReq(header)

	b, ok := l.buckets[group]
	if !ok {
		return
	}

	// Throttle aggressively when remaining requests are critically low.
	if remaining >= 0 && remaining <= 2 {
		b.setRate(1.0) // emergency limit: 1 request per second
		return
	}

	// Adjust rate to match the server-reported per-second allowance.
	if perSec > 0 {
		b.setRate(float64(perSec))
	}
}

// parseRemainingReq parses a Remaining-Req header.
// Header format: "group=market; min=599; sec=9"
// Returns (remaining, perSec) — both derived from the sec field.
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
			// minute value is used for reference only
			_ = val
		}
	}

	return remaining, perSec
}

// GroupFromPath infers the rate limit group from an HTTP request path.
// Based on actual Upbit API endpoint paths.
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

// String returns the string representation of the group name.
func (g Group) String() string {
	return string(g)
}

// Validate returns an error if the group is not a recognized rate limit group.
func (g Group) Validate() error {
	switch g {
	case GroupMarket, GroupCandle, GroupTicker, GroupOrderbook, GroupTrade,
		GroupDefault, GroupOrder, GroupOrderTest, GroupOrderCancelAll:
		return nil
	default:
		return fmt.Errorf("unknown rate limit group: %s", g)
	}
}
