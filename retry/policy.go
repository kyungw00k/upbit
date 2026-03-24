// Package retry provides configurable retry policies with exponential backoff.
package retry

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"time"
)

// RetryableError is an interface for determining whether an error is retryable.
type RetryableError interface {
	error
	HTTPStatus() int
}

// Policy defines a retry policy.
type Policy struct {
	MaxAttempts int
	InitialWait time.Duration
	MaxWait     time.Duration
	Multiplier  float64
}

// Default returns the default retry policy.
// Initial wait 1s, max 3 attempts, max wait 10s, exponential backoff.
func Default() *Policy {
	return &Policy{
		MaxAttempts: 3,
		InitialWait: 1 * time.Second,
		MaxWait:     10 * time.Second,
		Multiplier:  2.0,
	}
}

// ShouldRetry determines whether to retry based on the error and HTTP status code.
func (p *Policy) ShouldRetry(err error, statusCode int) bool {
	// 418: IP banned — fail immediately
	if statusCode == 418 {
		return false
	}

	// 429: Rate limit exceeded — retry
	if statusCode == 429 {
		return true
	}

	// 5xx: Server error — retry
	if statusCode >= 500 && statusCode <= 504 {
		return true
	}

	// Network error — retry
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) {
			return true
		}
	}

	return false
}

// Wait waits before a retry using exponential backoff with jitter.
func (p *Policy) Wait(ctx context.Context, attempt int) error {
	wait := p.InitialWait
	for i := 1; i < attempt; i++ {
		wait = time.Duration(float64(wait) * p.Multiplier)
	}
	if wait > p.MaxWait {
		wait = p.MaxWait
	}

	// Add 10% jitter to prevent thundering herd
	jitter := time.Duration(rand.Float64() * float64(wait) * 0.1)
	wait += jitter

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(wait):
		return nil
	}
}

// Retry executes fn with the default retry policy.
// If fn returns an error implementing RetryableError, the retry decision is based on it.
func Retry(ctx context.Context, fn func() error) error {
	return RetryWithPolicy(ctx, Default(), fn)
}

// RetryWithPolicy executes fn with a custom retry policy.
func RetryWithPolicy(ctx context.Context, policy *Policy, fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Extract HTTP status code (check if error implements RetryableError)
		statusCode := 0
		var retryable RetryableError
		if errors.As(err, &retryable) {
			statusCode = retryable.HTTPStatus()
		}

		if !policy.ShouldRetry(err, statusCode) {
			return err
		}

		// Last attempt — return without waiting
		if attempt == policy.MaxAttempts {
			break
		}

		if waitErr := policy.Wait(ctx, attempt); waitErr != nil {
			return waitErr
		}
	}

	return lastErr
}
