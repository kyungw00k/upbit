package retry

import (
	"context"
	"fmt"
	"testing"
)

// testRetryableError 테스트용 RetryableError 구현
type testRetryableError struct {
	statusCode int
	message    string
}

func (e *testRetryableError) Error() string {
	return e.message
}

func (e *testRetryableError) HTTPStatus() int {
	return e.statusCode
}

func TestRetry_SuccessOnFirstAttempt(t *testing.T) {
	callCount := 0
	err := Retry(context.Background(), func() error {
		callCount++
		return nil
	})
	if err != nil {
		t.Fatalf("성공 시 에러가 없어야 함: %v", err)
	}
	if callCount != 1 {
		t.Errorf("성공 시 1회만 호출되어야 함, 실제: %d", callCount)
	}
}

func TestRetry_RetryOn429(t *testing.T) {
	policy := &Policy{
		MaxAttempts: 3,
		InitialWait: 0, // 테스트 속도를 위해 대기 없음
		MaxWait:     0,
		Multiplier:  1.0,
	}

	callCount := 0
	err := RetryWithPolicy(context.Background(), policy, func() error {
		callCount++
		if callCount < 3 {
			return &testRetryableError{statusCode: 429, message: "rate limited"}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("3번째 시도에서 성공해야 함: %v", err)
	}
	if callCount != 3 {
		t.Errorf("3회 호출 기대, 실제: %d", callCount)
	}
}

func TestRetry_NoRetryOn418(t *testing.T) {
	policy := &Policy{
		MaxAttempts: 3,
		InitialWait: 0,
		MaxWait:     0,
		Multiplier:  1.0,
	}

	callCount := 0
	err := RetryWithPolicy(context.Background(), policy, func() error {
		callCount++
		return &testRetryableError{statusCode: 418, message: "IP banned"}
	})
	if err == nil {
		t.Fatal("418 에러 시 즉시 반환되어야 함")
	}
	if callCount != 1 {
		t.Errorf("418은 재시도 없이 1회만 호출되어야 함, 실제: %d", callCount)
	}
}

func TestRetry_MaxAttemptsExceeded(t *testing.T) {
	policy := &Policy{
		MaxAttempts: 3,
		InitialWait: 0,
		MaxWait:     0,
		Multiplier:  1.0,
	}

	callCount := 0
	err := RetryWithPolicy(context.Background(), policy, func() error {
		callCount++
		return &testRetryableError{statusCode: 500, message: "server error"}
	})
	if err == nil {
		t.Fatal("최대 횟수 초과 시 에러가 반환되어야 함")
	}
	if callCount != 3 {
		t.Errorf("최대 3회 호출 기대, 실제: %d", callCount)
	}
}

func TestRetry_NonRetryableError(t *testing.T) {
	policy := &Policy{
		MaxAttempts: 3,
		InitialWait: 0,
		MaxWait:     0,
		Multiplier:  1.0,
	}

	callCount := 0
	err := RetryWithPolicy(context.Background(), policy, func() error {
		callCount++
		return fmt.Errorf("일반 에러 (RetryableError 미구현)")
	})
	if err == nil {
		t.Fatal("재시도 불가 에러 시 즉시 반환되어야 함")
	}
	if callCount != 1 {
		t.Errorf("재시도 불가 에러는 1회만 호출, 실제: %d", callCount)
	}
}

func TestRetry_ContextCancelled(t *testing.T) {
	policy := &Policy{
		MaxAttempts: 5,
		InitialWait: 0,
		MaxWait:     0,
		Multiplier:  1.0,
	}

	ctx, cancel := context.WithCancel(context.Background())

	callCount := 0
	err := RetryWithPolicy(ctx, policy, func() error {
		callCount++
		if callCount == 2 {
			cancel() // 2번째 시도 후 취소
		}
		return &testRetryableError{statusCode: 500, message: "server error"}
	})
	if err == nil {
		t.Fatal("context 취소 시 에러가 반환되어야 함")
	}
}

func TestShouldRetry(t *testing.T) {
	p := Default()

	tests := []struct {
		name     string
		status   int
		expected bool
	}{
		{"418 IP 차단", 418, false},
		{"429 Rate limit", 429, true},
		{"500 서버 에러", 500, true},
		{"502 Bad Gateway", 502, true},
		{"503 Service Unavailable", 503, true},
		{"504 Gateway Timeout", 504, true},
		{"400 Bad Request", 400, false},
		{"401 Unauthorized", 401, false},
		{"404 Not Found", 404, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.ShouldRetry(nil, tt.status)
			if got != tt.expected {
				t.Errorf("ShouldRetry(nil, %d) = %v, 기대: %v", tt.status, got, tt.expected)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	p := Default()
	if p.MaxAttempts != 3 {
		t.Errorf("MaxAttempts 기대: 3, 실제: %d", p.MaxAttempts)
	}
	if p.Multiplier != 2.0 {
		t.Errorf("Multiplier 기대: 2.0, 실제: %f", p.Multiplier)
	}
}
