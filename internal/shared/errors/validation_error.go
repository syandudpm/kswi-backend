package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validationMessages = map[string]string{
	"required":  "This field is required",
	"email":     "Must be a valid email address",
	"min":       "Must be at least %s characters long",
	"max":       "Must not exceed %s characters",
	"oneof":     "Must be one of: %s",
	"numeric":   "Must be a number",
	"alphanum":  "Must contain only letters and numbers",
	"lowercase": "Must be lowercase",
	"uppercase": "Must be uppercase",
}

func HandleValidationError(err error) *AppError {
	var validationErrors []ValidationError

	// Check for empty body
	if err.Error() == "EOF" {
		return NewAppError(
			TypeValidation,
			"Validation failed",
			[]ValidationError{{
				Field:   "body",
				Message: "Request body is empty. Please provide valid JSON data",
			}},
		)
	}

	// Check for JSON syntax errors
	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return NewAppError(
			TypeValidation,
			"Validation failed",
			[]ValidationError{{
				Field:   "body",
				Message: "Invalid JSON format. Please check your request body",
			}},
		)
	}

	// Check if it's a validator.ValidationErrors
	validatorErrs, ok := err.(validator.ValidationErrors)
	if ok {
		// Process each validation error
		for _, e := range validatorErrs {
			field := strings.ToLower(e.Field())
			message := validationMessages[e.Tag()]

			if message == "" {
				message = "Invalid value"
			}

			// Handle parameters for certain validation tags
			switch e.Tag() {
			case "min", "max", "oneof":
				message = strings.Replace(message, "%s", e.Param(), 1)
			}

			validationErrors = append(validationErrors, ValidationError{
				Field:   field,
				Message: message,
			})
		}
	} else {
		// Try to handle json.UnmarshalTypeError
		var unmarshalTypeError *json.UnmarshalTypeError
		if errors.As(err, &unmarshalTypeError) {
			validationErrors = append(validationErrors, ValidationError{
				Field: unmarshalTypeError.Field,
				Message: fmt.Sprintf("Invalid type for field '%s': expected %v",
					unmarshalTypeError.Field,
					unmarshalTypeError.Type),
			})
		} else {
			// Fallback for other errors
			validationErrors = append(validationErrors, ValidationError{
				Field:   "body",
				Message: fmt.Sprintf("Invalid request: %s", err.Error()),
			})
		}
	}

	return NewAppError(
		TypeValidation,
		"Validation failed",
		validationErrors,
	)
}
