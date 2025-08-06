package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	TypeValidation ErrorType = "VALIDATION_ERROR"
	TypeAuth       ErrorType = "AUTHENTICATION_ERROR"
	TypeDatabase   ErrorType = "DATABASE_ERROR"
	TypeInternal   ErrorType = "INTERNAL_ERROR"
	TypeNotFound   ErrorType = "NOT_FOUND"
)

type AppError struct {
	Type        ErrorType   `json:"type"`
	Message     string      `json:"message"`
	Details     interface{} `json:"details,omitempty"`
	OriginalErr error       `json:"-"` // Store original error for logging, exclude from JSON
	StackTrace  string      `json:"-"` // Optional stack trace for debugging
}

type Response struct {
	Success bool        `json:"success"`
	Error   *AppError   `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// GetOriginalError returns the original underlying error for logging purposes
func (e *AppError) GetOriginalError() error {
	return e.OriginalErr
}

// GetFullDetails returns a map with all error details for structured logging
func (e *AppError) GetFullDetails() map[string]interface{} {
	details := map[string]interface{}{
		"type":    e.Type,
		"message": e.Message,
	}

	if e.Details != nil {
		details["details"] = e.Details
	}

	if e.OriginalErr != nil {
		details["original_error"] = e.OriginalErr.Error()
	}

	if e.StackTrace != "" {
		details["stack_trace"] = e.StackTrace
	}

	return details
}

// NewAppError creates a new structured application error
func NewAppError(errorType ErrorType, message string, details interface{}) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Details: details,
	}
}

// NewAppErrorWithOriginal creates a new structured application error with original error preserved
func NewAppErrorWithOriginal(errorType ErrorType, message string, details interface{}, originalErr error) *AppError {
	return &AppError{
		Type:        errorType,
		Message:     message,
		Details:     details,
		OriginalErr: originalErr,
	}
}

// GetStatusCode maps error types to HTTP status codes
func (e *AppError) GetStatusCode() int {
	switch e.Type {
	case TypeValidation:
		return http.StatusBadRequest
	case TypeAuth:
		return http.StatusUnauthorized
	case TypeDatabase, TypeInternal:
		return http.StatusInternalServerError
	case TypeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// Convenience constructors - Enhanced versions
func NewValidationError(details interface{}) *AppError {
	return NewAppError(TypeValidation, "Validation failed", details)
}

func NewValidationErrorWithOriginal(details interface{}, originalErr error) *AppError {
	return NewAppErrorWithOriginal(TypeValidation, "Validation failed", details, originalErr)
}

func NewNotFoundError(resource string) *AppError {
	return NewAppError(TypeNotFound, resource+" not found", nil)
}

func NewDatabaseError(err error) *AppError {
	return NewAppErrorWithOriginal(TypeDatabase, "Database operation failed", err.Error(), err)
}

func NewInternalError(err error) *AppError {
	return NewAppErrorWithOriginal(TypeInternal, "Internal error", err.Error(), err)
}

// Enhanced As function with better error checking
func As(err error, target interface{}) bool {
	if e, ok := err.(*AppError); ok {
		if t, ok := target.(**AppError); ok {
			*t = e
			return true
		}
	}
	return false
}

func ShouldHideDetails(errorType ErrorType) bool {
	switch errorType {
	case TypeInternal, TypeDatabase, TypeAuth:
		return true
	default:
		return false
	}
}

// Helper function to check if error should be logged with full details
func ShouldLogFullDetails(errorType ErrorType) bool {
	switch errorType {
	case TypeInternal, TypeDatabase:
		return true
	default:
		return false
	}
}
