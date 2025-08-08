package utils

import (
	"fmt"
	"go-restaurant-management/internal/shared/errors/exceptions"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func init() {
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateStruct(s interface{}) error {
	if err := Validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return FormatValidationError(validationErrors)
		}
		return exceptions.NewValidationError("validation", err.Error())
	}
	return nil
}

func FormatValidationError(errs validator.ValidationErrors) error {
	var errorMessages []map[string]interface{}

	for _, err := range errs {
		field := err.Field()
		message := getValidationMessage(err)

		errorMessages = append(errorMessages, map[string]interface{}{
			"field":   field,
			"message": message,
			"tag":     err.Tag(),
			"value":   err.Value(),
		})
	}

	if len(errorMessages) == 1 {
		return exceptions.NewValidationError(
			errorMessages[0]["field"].(string),
			errorMessages[0]["message"].(string),
		)
	}

	return exceptions.NewMultipleValidationErrors(errorMessages)
}

func getValidationMessage(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("The field %s is required", field)
	case "email":
		return fmt.Sprintf("The field %s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("The field %s must have at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("The field %s must have at most %s characters", field, err.Param())
	case "len":
		return fmt.Sprintf("The field %s must have exactly %s characters", field, err.Param())
	case "numeric":
		return fmt.Sprintf("The field %s must contain only numbers", field)
	case "eqfield":
		return fmt.Sprintf("The field %s must be equal to %s", field, err.Param())
	default:
		return fmt.Sprintf("The field %s is invalid: %s", field, err.Tag())
	}
}
