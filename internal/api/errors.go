package api

import (
	"encoding/json"
	"fmt"
)

// APIError 구조화된 Upbit API 에러
type APIError struct {
	Code       string                 `json:"name"`
	Message    string                 `json:"message"`
	Detail     map[string]interface{} `json:"-"`
	StatusCode int                    `json:"-"`
}

// apiErrorWrapper Upbit API 에러 JSON 래퍼
// 응답 형식: {"error": {"name": "...", "message": "..."}}
type apiErrorWrapper struct {
	Error apiErrorBody `json:"error"`
}

type apiErrorBody struct {
	Name    interface{} `json:"name"`
	Message string      `json:"message"`
}

// Error error 인터페이스 구현
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("[%s] %s", e.Code, e.Message)
	}
	return e.Message
}

// HTTPStatus RetryableError 인터페이스 구현
func (e *APIError) HTTPStatus() int {
	return e.StatusCode
}

// ExitCode 에러 종류별 프로세스 종료 코드 반환
// 0: 성공, 1: 일반, 2: 인증, 3: rate limit, 4: 입력검증
func (e *APIError) ExitCode() int {
	switch e.StatusCode {
	case 401, 403:
		return 2 // 인증 오류
	case 429:
		return 3 // Rate limit
	case 400, 422:
		return 4 // 입력검증 오류
	default:
		if e.StatusCode >= 400 {
			return 1 // 일반 오류
		}
		return 0
	}
}

// IsAuthError 인증 오류 여부
func (e *APIError) IsAuthError() bool {
	return e.ExitCode() == 2
}

// IsRateLimitError Rate limit 오류 여부
func (e *APIError) IsRateLimitError() bool {
	return e.StatusCode == 429
}

// IsValidationError 입력검증 오류 여부
func (e *APIError) IsValidationError() bool {
	return e.ExitCode() == 4
}

// ParseAPIError HTTP 응답 본문에서 APIError 파싱
// {"error": {"name": "...", "message": "..."}} 형식 처리
func ParseAPIError(body []byte, statusCode int) *APIError {
	var wrapper apiErrorWrapper
	if err := json.Unmarshal(body, &wrapper); err == nil && wrapper.Error.Message != "" {
		apiErr := &APIError{
			Message:    wrapper.Error.Message,
			StatusCode: statusCode,
		}

		// name 필드는 문자열 또는 숫자일 수 있음
		switch v := wrapper.Error.Name.(type) {
		case string:
			apiErr.Code = v
		case float64:
			apiErr.Code = fmt.Sprintf("%d", int(v))
		default:
			apiErr.Code = fmt.Sprintf("%v", v)
		}

		return apiErr
	}

	// 파싱 실패 시 기본 에러 반환
	return &APIError{
		Code:       fmt.Sprintf("%d", statusCode),
		Message:    string(body),
		StatusCode: statusCode,
	}
}

// IsAPIError 에러가 APIError인지 확인
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

// AsAPIError 에러를 APIError로 변환 시도
func AsAPIError(err error) (*APIError, bool) {
	apiErr, ok := err.(*APIError)
	return apiErr, ok
}
