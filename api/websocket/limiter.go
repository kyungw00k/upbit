package websocket

import (
	"context"
	"sync"
	"time"
)

// Upbit WebSocket rate limits (https://docs.upbit.com/reference/rate-limits):
//   - websocket-connect: 5 connections/sec (per IP without auth, per pocket with auth)
//   - websocket-message: 5 messages/sec and 100 messages/min per connection
const (
	maxConnectsPerSec = 5.0
	maxMessagesPerSec = 5.0
	maxMessagesPerMin = 100.0
)

// connectLimiter throttles dial attempts process-wide: the connect budget
// is shared per IP (unauthenticated) or per pocket (authenticated), so all
// clients in this process draw from the same allowance.
var connectLimiter = newMsgLimiter(maxConnectsPerSec, 0)

// msgLimiter is a dual-window token bucket (per-second and per-minute).
// A zero rate disables that window.
type msgLimiter struct {
	mu         sync.Mutex
	perSec     float64
	perMin     float64
	secTokens  float64
	minTokens  float64
	lastRefill time.Time
}

func newMsgLimiter(perSec, perMin float64) *msgLimiter {
	return &msgLimiter{
		perSec:     perSec,
		perMin:     perMin,
		secTokens:  perSec,
		minTokens:  perMin,
		lastRefill: time.Now(),
	}
}

// refillLocked refills both windows based on elapsed time.
// Caller must hold l.mu.
func (l *msgLimiter) refillLocked() {
	now := time.Now()
	elapsed := now.Sub(l.lastRefill).Seconds()
	if elapsed <= 0 {
		return
	}
	if l.perSec > 0 {
		l.secTokens = min(l.perSec, l.secTokens+elapsed*l.perSec)
	}
	if l.perMin > 0 {
		l.minTokens = min(l.perMin, l.minTokens+elapsed*l.perMin/60)
	}
	l.lastRefill = now
}

// allow reports whether one message may be sent right now.
func (l *msgLimiter) allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refillLocked()
	if l.perSec > 0 && l.secTokens < 1 {
		return false
	}
	if l.perMin > 0 && l.minTokens < 1 {
		return false
	}
	l.secTokens--
	l.minTokens--
	return true
}

// wait blocks until one message may be sent or ctx is done.
func (l *msgLimiter) wait(ctx context.Context) error {
	for {
		if l.allow() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(20 * time.Millisecond):
		}
	}
}

// reset restores full budgets (new connection, fresh allowance).
func (l *msgLimiter) reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.secTokens = l.perSec
	l.minTokens = l.perMin
	l.lastRefill = time.Now()
}
