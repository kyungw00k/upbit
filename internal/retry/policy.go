package retry

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"time"
)

// RetryableError 재시도 가능 여부를 판단하기 위한 인터페이스
type RetryableError interface {
	error
	HTTPStatus() int
}

// Policy 재시도 정책
type Policy struct {
	MaxAttempts int
	InitialWait time.Duration
	MaxWait     time.Duration
	Multiplier  float64
}

// Default 기본 재시도 정책 반환
// 초기 1초, 최대 3회, 최대 대기 10초, 지수 백오프
func Default() *Policy {
	return &Policy{
		MaxAttempts: 3,
		InitialWait: 1 * time.Second,
		MaxWait:     10 * time.Second,
		Multiplier:  2.0,
	}
}

// ShouldRetry 에러 및 HTTP 상태코드에 따라 재시도 여부 결정
func (p *Policy) ShouldRetry(err error, statusCode int) bool {
	// 418: IP 차단 — 즉시 실패
	if statusCode == 418 {
		return false
	}

	// 429: Rate limit — 재시도
	if statusCode == 429 {
		return true
	}

	// 5xx: 서버 오류 — 재시도
	if statusCode >= 500 && statusCode <= 504 {
		return true
	}

	// 네트워크 오류 — 재시도
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) {
			return true
		}
	}

	return false
}

// Wait 재시도 전 대기 (지수 백오프 + jitter)
func (p *Policy) Wait(ctx context.Context, attempt int) error {
	wait := p.InitialWait
	for i := 1; i < attempt; i++ {
		wait = time.Duration(float64(wait) * p.Multiplier)
	}
	if wait > p.MaxWait {
		wait = p.MaxWait
	}

	// 10% jitter 추가 (thundering herd 방지)
	jitter := time.Duration(rand.Float64() * float64(wait) * 0.1)
	wait += jitter

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(wait):
		return nil
	}
}

// Retry 재시도 로직 실행
// fn: 실행할 함수 — RetryableError를 구현하는 에러를 반환하면 재시도 판단
func Retry(ctx context.Context, fn func() error) error {
	return RetryWithPolicy(ctx, Default(), fn)
}

// RetryWithPolicy 커스텀 정책으로 재시도 실행
func RetryWithPolicy(ctx context.Context, policy *Policy, fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// HTTP 상태코드 추출 (RetryableError 구현 여부 확인)
		statusCode := 0
		var retryable RetryableError
		if errors.As(err, &retryable) {
			statusCode = retryable.HTTPStatus()
		}

		if !policy.ShouldRetry(err, statusCode) {
			return err
		}

		// 마지막 시도이면 대기 없이 반환
		if attempt == policy.MaxAttempts {
			break
		}

		if waitErr := policy.Wait(ctx, attempt); waitErr != nil {
			return waitErr
		}
	}

	return lastErr
}
