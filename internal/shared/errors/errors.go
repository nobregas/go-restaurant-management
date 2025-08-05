package errors

type ErrorType string

const (
	NOT_FOUND    ErrorType = "NOT_FOUND"
	BAD_REQUEST  ErrorType = "BAD_REQUEST"
	UNEXPECTED   ErrorType = "UNEXPECTED"
	UNAUTHORIZED ErrorType = "UNAUTHORIZED"
	FORBIDDEN    ErrorType = "FORBIDDEN"
	CONFLICT     ErrorType = "CONFLICT"
	INTERNAL     ErrorType = "INTERNAL"
)
