package utils

import (
	"encoding/json"
	"fmt"
	"go-restaurant-management/internal/shared/errors"
	"go-restaurant-management/internal/shared/errors/exceptions"
	"net/http"
)

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("parser: missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func ParseAndValidateJson(r *http.Request, payload any) error {
	if err := ParseJson(r, payload); err != nil {
		return exceptions.NewInvalidJSONError(err)
	}

	if err := ValidateStruct(payload); err != nil {
		return err
	}

	return nil
}

func WriteJson(w http.ResponseWriter, status int, v any) {
	if status == http.StatusNoContent || status == http.StatusNotModified {
		w.WriteHeader(status)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to encode response"}`))
	}
}

func WriteError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		WriteJson(w, appErr.HTTPStatusCode(), appErr)
		return
	}
	genericError := &errors.AppError{
		Type:    errors.INTERNAL,
		Code:    "GENERIC_ERROR",
		Message: "An unexpected error occurred",
		Details: map[string]interface{}{
			"error": err.Error(),
		},
	}
	WriteJson(w, http.StatusInternalServerError, genericError)
}

func WriteErrorWithStatus(w http.ResponseWriter, status int, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		WriteJson(w, status, appErr)
		return
	}

	genericError := &errors.AppError{
		Type:    errors.INTERNAL,
		Code:    "GENERIC_ERROR",
		Message: "An error occurred",
		Details: map[string]interface{}{
			"error": err.Error(),
		},
	}
	WriteJson(w, status, genericError)
}
