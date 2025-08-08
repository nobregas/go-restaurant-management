package middleware

import (
	"fmt"
	"go-restaurant-management/internal/shared/errors"
	"go-restaurant-management/internal/shared/utils"
	"log"
	"net/http"
	"runtime/debug"
)

func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
				handlePanic(w, err)
			}
		}()
		next(w, r)
	}
}

func handlePanic(w http.ResponseWriter, err interface{}) {
	switch e := err.(type) {
	case *errors.AppError:
		utils.WriteError(w, e)

	case error:
		internalErr := &errors.AppError{
			Type:    errors.INTERNAL,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error occurred",
			Details: map[string]interface{}{
				"error": e.Error(),
			},
			Cause: e,
		}
		utils.WriteError(w, internalErr)

	default:
		genericErr := &errors.AppError{
			Type:    errors.INTERNAL,
			Code:    "UNEXPECTED_ERROR",
			Message: "An unexpected error occurred",
			Details: map[string]interface{}{
				"panic_value": fmt.Sprintf("%v", err),
			},
		}
		utils.WriteError(w, genericErr)
	}
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func ErrorHandlerFunc(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			utils.WriteError(w, err)
		}
	}
}
