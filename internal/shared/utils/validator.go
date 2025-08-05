package utils

import (
	"fmt"
	"go-restaurant-management/internal/shared/errors/exceptions"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func FormatValidationError(errs validator.ValidationErrors) error {
	var errorList []string

	for _, err := range errs {
		// if json tag is set, use it
		field := strings.ToLower(err.Field())
		if jsonTag := getJsonTag(err); jsonTag != "" {
			field = jsonTag
		}

		// create friendly message
		switch err.Tag() {
		case "required":
			errorList = append(errorList, fmt.Sprintf("The field %s is required", field))
		case "email":
			errorList = append(errorList, fmt.Sprintf("The field %s must be a valid email", field))
		case "min":
			errorList = append(errorList, fmt.Sprintf("The field %s must have at least %s characters", field, err.Param()))
		case "max":
			errorList = append(errorList, fmt.Sprintf("The field %s must have at most %s characters", field, err.Param()))
		case "eqfield":
			errorList = append(errorList, fmt.Sprintf("The field %s must be equal to %s", field, err.Param()))
		case "cpf":
			errorList = append(errorList, fmt.Sprintf("The field %s must be a valid CPF", field))
		default:
			errorList = append(errorList, fmt.Sprintf("Erro no campo %s: %s", field, err.Tag()))
		}
	}

	return exceptions.NewValidationError("", strings.Join(errorList, ". "))
}

func getJsonTag(err validator.FieldError) string {
	if field, ok := err.Type().FieldByName(err.Field()); ok {
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			// remove json tag
			return strings.Split(jsonTag, ",")[0]
		}
	}
	return ""
}
