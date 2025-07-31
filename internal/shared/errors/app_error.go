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
	Type    ErrorType   `json:"type"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Error   *AppError   `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// NewAppError creates a new structured application error
func NewAppError(errorType ErrorType, message string, details interface{}) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Details: details,
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

// Convenience constructors
func NewValidationError(details interface{}) *AppError {
	return NewAppError(TypeValidation, "Validation failed", details)
}

func NewNotFoundError(resource string) *AppError {
	return NewAppError(TypeNotFound, resource+" not found", nil)
}

func NewDatabaseError(err error) *AppError {
	return NewAppError(TypeDatabase, "Database operation failed", err.Error())
}

func NewInternalError(err error) *AppError {
	return NewAppError(TypeInternal, "Internal error", err.Error())
}

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
