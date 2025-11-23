package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAppError(t *testing.T) {
	tests := []struct {
		name    string
		errType ErrorType
		message string
		err     error
		want    *AppError
	}{
		{
			name:    "parse error without wrapped error",
			errType: ErrParse,
			message: "invalid JSON",
			err:     nil,
			want: &AppError{
				Type:    ErrParse,
				Message: "invalid JSON",
				Err:     nil,
			},
		},
		{
			name:    "repository error with wrapped error",
			errType: ErrRepository,
			message: "failed to read file",
			err:     errors.New("permission denied"),
			want: &AppError{
				Type:    ErrRepository,
				Message: "failed to read file",
				Err:     errors.New("permission denied"),
			},
		},
		{
			name:    "validation error",
			errType: ErrValidation,
			message: "conversation ID cannot be empty",
			err:     nil,
			want: &AppError{
				Type:    ErrValidation,
				Message: "conversation ID cannot be empty",
				Err:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAppError(tt.errType, tt.message, tt.err)
			assert.Equal(t, tt.want.Type, got.Type)
			assert.Equal(t, tt.want.Message, got.Message)
			assert.Equal(t, tt.want.Err, got.Err)
		})
	}
}

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name    string
		appErr  *AppError
		wantMsg string
	}{
		{
			name:    "error without wrapped error",
			appErr:  NewAppError(ErrParse, "invalid JSON", nil),
			wantMsg: "[PARSE] invalid JSON",
		},
		{
			name:    "error with wrapped error",
			appErr:  NewAppError(ErrRepository, "failed to read", errors.New("file not found")),
			wantMsg: "[REPOSITORY] failed to read: file not found",
		},
		{
			name:    "validation error",
			appErr:  NewAppError(ErrValidation, "empty title", nil),
			wantMsg: "[VALIDATION] empty title",
		},
		{
			name:    "unknown error type",
			appErr:  NewAppError(ErrorType("UNKNOWN"), "something went wrong", nil),
			wantMsg: "[UNKNOWN] something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.appErr.Error()
			assert.Equal(t, tt.wantMsg, got)
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	innerErr := errors.New("original error")
	appErr := NewAppError(ErrRepository, "wrap failed", innerErr)

	unwrapped := errors.Unwrap(appErr)
	require.NotNil(t, unwrapped)
	assert.Equal(t, innerErr, unwrapped)
}

func TestIsErrorType(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		errType  ErrorType
		wantTrue bool
	}{
		{
			name:     "matching error type",
			err:      NewAppError(ErrParse, "msg", nil),
			errType:  ErrParse,
			wantTrue: true,
		},
		{
			name:     "non-matching error type",
			err:      NewAppError(ErrParse, "msg", nil),
			errType:  ErrRepository,
			wantTrue: false,
		},
		{
			name:     "not an AppError",
			err:      errors.New("generic error"),
			errType:  ErrParse,
			wantTrue: false,
		},
		{
			name:     "nil error",
			err:      nil,
			errType:  ErrParse,
			wantTrue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsErrorType(tt.err, tt.errType)
			assert.Equal(t, tt.wantTrue, got)
		})
	}
}
