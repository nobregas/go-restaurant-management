package middleware

import (
	"fmt"
	"go-restaurant-management/internal/shared/errors"
	"go-restaurant-management/internal/shared/utils"
	"net/http"
)

func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				handleError(w, err)
			}
		}()
		next(w, r)
	}
}

func handleError(w http.ResponseWriter, err interface{}) {
	switch e := err.(type) {
	case *errors.AppError:
		writeAppError(w, e)

	case error:
		utils.WriteError(w, http.StatusInternalServerError, e)

	default:
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("unexpected error: %v", err))
	}
}

func writeAppError(w http.ResponseWriter, e *errors.AppError) {
	status := http.StatusInternalServerError
	switch e.Type {
	case errors.BAD_REQUEST:
		status = http.StatusBadRequest
	case errors.UNAUTHORIZED:
		status = http.StatusUnauthorized
	case errors.FORBIDDEN:
		status = http.StatusForbidden
	case errors.NOT_FOUND:
		status = http.StatusNotFound
	case errors.CONFLICT:
		status = http.StatusConflict
	case errors.INTERNAL:
		status = http.StatusInternalServerError
	}

	utils.WriteJson(w, status, map[string]interface{}{
		"error":   e.Code,
		"message": e.Message,
		"details": e.Details,
	})
}
