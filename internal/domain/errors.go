// ABOUTME: errors.go defines the error hierarchy and error handling for the domain layer.
// It provides typed errors with context and support for error wrapping.
package domain

import "fmt"

// ErrorType defines the category of an error.
type ErrorType string

// Error type constants.
const (
	ErrParse      ErrorType = "PARSE"
	ErrRepository ErrorType = "REPOSITORY"
	ErrValidation ErrorType = "VALIDATION"
	ErrSearch     ErrorType = "SEARCH"
	ErrNotFound   ErrorType = "NOT_FOUND"
	ErrInternal   ErrorType = "INTERNAL"
)

// AppError is a domain-specific error that wraps errors with context.
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

// NewAppError creates a new AppError with the given type, message, and wrapped error.
func NewAppError(errType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap implements errors.Unwrap for error chain traversal.
func (e *AppError) Unwrap() error {
	return e.Err
}

// IsErrorType checks if an error is of a specific AppError type.
func IsErrorType(err error, errType ErrorType) bool {
	appErr, ok := err.(*AppError)
	if !ok {
		return false
	}
	return appErr.Type == errType
}
