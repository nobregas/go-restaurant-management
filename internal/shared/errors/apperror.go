package errors

import (
	"fmt"
)

type AppError struct {
	Type    ErrorType              `json:"type"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Cause   error                  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (caused by: %v)", e.Type, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func (e *AppError) HTTPStatusCode() int {
	switch e.Type {
	case NOT_FOUND:
		return 404
	case BAD_REQUEST:
		return 400
	case UNAUTHORIZED:
		return 401
	case FORBIDDEN:
		return 403
	case CONFLICT:
		return 409
	case INTERNAL:
		return 500
	default:
		return 500
	}
}
