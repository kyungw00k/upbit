package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

// APIError represents a structured Upbit API error.
type APIError struct {
	Code       string                 `json:"name"`
	Message    string                 `json:"message"`
	Detail     map[string]interface{} `json:"-"`
	StatusCode int                    `json:"-"`
}

// apiErrorWrapper is the JSON wrapper for Upbit API errors.
// Response format: {"error": {"name": "...", "message": "..."}}
type apiErrorWrapper struct {
	Error apiErrorBody `json:"error"`
}

type apiErrorBody struct {
	Name    interface{} `json:"name"`
	Message string      `json:"message"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("[%s] %s", e.Code, e.Message)
	}
	return e.Message
}

// HTTPStatus implements the RetryableError interface.
func (e *APIError) HTTPStatus() int {
	return e.StatusCode
}

// ExitCode returns the process exit code for the error type.
// 0: success, 1: general, 2: auth, 3: rate limit, 4: input validation
func (e *APIError) ExitCode() int {
	switch e.StatusCode {
	case 401, 403:
		return 2 // auth error
	case 429:
		return 3 // rate limit
	case 400, 422:
		return 4 // input validation error
	default:
		if e.StatusCode >= 400 {
			return 1 // general error
		}
		return 0
	}
}

// IsAuthError reports whether the error is an authentication error.
func (e *APIError) IsAuthError() bool {
	return e.ExitCode() == 2
}

// IsRateLimitError reports whether the error is a rate limit error.
func (e *APIError) IsRateLimitError() bool {
	return e.StatusCode == 429
}

// IsValidationError reports whether the error is an input validation error.
func (e *APIError) IsValidationError() bool {
	return e.ExitCode() == 4
}

// ParseAPIError parses an APIError from an HTTP response body.
// Handles the {"error": {"name": "...", "message": "..."}} format.
func ParseAPIError(body []byte, statusCode int) *APIError {
	var wrapper apiErrorWrapper
	if err := json.Unmarshal(body, &wrapper); err == nil && wrapper.Error.Message != "" {
		apiErr := &APIError{
			Message:    wrapper.Error.Message,
			StatusCode: statusCode,
		}

		// The name field may be a string or a number.
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

	// Fallback: return a generic error if parsing fails.
	return &APIError{
		Code:       fmt.Sprintf("%d", statusCode),
		Message:    string(body),
		StatusCode: statusCode,
	}
}

// IsAPIError reports whether err is (or wraps) an *APIError.
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// AsAPIError attempts to extract an *APIError from err, unwrapping if necessary.
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	ok := errors.As(err, &apiErr)
	return apiErr, ok
}
