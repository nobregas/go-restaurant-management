package exceptions

import (
	"fmt"
	"go-restaurant-management/internal/shared/errors"
)

func NewEntityNotFound(entity string, id interface{}) *errors.AppError {
	return &errors.AppError{
		Type:    errors.NOT_FOUND,
		Code:    "ENTITY_NOT_FOUND",
		Message: fmt.Sprintf("%s not found", entity),
		Details: map[string]interface{}{
			"entity": entity,
			"id":     id,
		},
	}
}

func NewValidationError(field string, reason string) *errors.AppError {
	return &errors.AppError{
		Type:    errors.BAD_REQUEST,
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Details: map[string]interface{}{
			"field":  field,
			"reason": reason,
		},
	}
}

func NewInvalidJSONError(cause error) *errors.AppError {
	return &errors.AppError{
		Type:    errors.BAD_REQUEST,
		Code:    "INVALID_JSON",
		Message: "Invalid JSON format",
		Details: map[string]interface{}{
			"reason": "The request body contains invalid JSON",
		},
		Cause: cause,
	}
}

func NewConflictError(field string, reason string) *errors.AppError {
	return &errors.AppError{
		Type:    errors.CONFLICT,
		Code:    "CONFLICT_ERROR",
		Message: "Resource conflict",
		Details: map[string]interface{}{
			"field":  field,
			"reason": reason,
		},
	}
}

func NewInternalServerError(reason string) *errors.AppError {
	return &errors.AppError{
		Type:    errors.INTERNAL,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal Server Error",
		Details: map[string]interface{}{
			"reason": reason,
		},
	}
}

func NewMethodNotAllowedError(method, path string) *errors.AppError {
	return &errors.AppError{
		Type:    errors.BAD_REQUEST,
		Code:    "METHOD_NOT_ALLOWED",
		Message: "Method not allowed",
		Details: map[string]interface{}{
			"method": method,
			"path":   path,
			"reason": fmt.Sprintf("Method %s is not allowed for path %s", method, path),
		},
	}
}

func NewRouteNotFoundError(path string) *errors.AppError {
	return &errors.AppError{
		Type:    errors.NOT_FOUND,
		Code:    "ROUTE_NOT_FOUND",
		Message: "Route not found",
		Details: map[string]interface{}{
			"path":   path,
			"reason": fmt.Sprintf("Route %s not found", path),
		},
	}
}

func NewUnauthorizedError(reason string) *errors.AppError {
	return &errors.AppError{
		Type:    errors.UNAUTHORIZED,
		Code:    "UNAUTHORIZED",
		Message: "Authentication failed",
		Details: map[string]interface{}{
			"reason": reason,
		},
	}
}

func NewMultipleValidationErrors(errors_ []map[string]interface{}) *errors.AppError {
	return &errors.AppError{
		Type:    errors.BAD_REQUEST,
		Code:    "MULTIPLE_VALIDATION_ERRORS",
		Message: "Multiple validation errors occurred",
		Details: map[string]interface{}{
			"errors": errors_,
		},
	}
}
