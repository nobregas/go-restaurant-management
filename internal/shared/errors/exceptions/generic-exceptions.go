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
		Code:    "BAD_REQUEST",
		Message: "BAD Request error",
		Details: map[string]interface{}{
			"field":  field,
			"reason": reason,
		},
	}
}

func NewConflictError(field string, reason string) *errors.AppError {
	return &errors.AppError{
		Type:    errors.CONFLICT,
		Code:    "Conflict",
		Message: "Conflict Request error",
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
