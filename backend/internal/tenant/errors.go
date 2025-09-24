package tenant

import (
	"errors"
	"fmt"
)

// Custom error types for better error handling
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "validation"
	ErrorTypeNotFound     ErrorType = "not_found"
	ErrorTypeConflict     ErrorType = "conflict"
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	ErrorTypeForbidden    ErrorType = "forbidden"
	ErrorTypeInternal     ErrorType = "internal"
)

// TenantError represents a tenant-specific error with type information
type TenantError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"`
	Code    string    `json:"code,omitempty"`
}

func (e *TenantError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s (field: %s)", e.Type, e.Message, e.Field)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Error constructors
func NewValidationError(message, field string) *TenantError {
	return &TenantError{
		Type:    ErrorTypeValidation,
		Message: message,
		Field:   field,
		Code:    "VALIDATION_ERROR",
	}
}

func NewNotFoundError(message string) *TenantError {
	return &TenantError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Code:    "NOT_FOUND",
	}
}

func NewConflictError(message string) *TenantError {
	return &TenantError{
		Type:    ErrorTypeConflict,
		Message: message,
		Code:    "CONFLICT",
	}
}

func NewUnauthorizedError(message string) *TenantError {
	return &TenantError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Code:    "UNAUTHORIZED",
	}
}

func NewForbiddenError(message string) *TenantError {
	return &TenantError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Code:    "FORBIDDEN",
	}
}

func NewInternalError(message string) *TenantError {
	return &TenantError{
		Type:    ErrorTypeInternal,
		Message: message,
		Code:    "INTERNAL_ERROR",
	}
}

// Helper functions to check error types
func IsValidationError(err error) bool {
	var tenantErr *TenantError
	return errors.As(err, &tenantErr) && tenantErr.Type == ErrorTypeValidation
}

func IsNotFoundError(err error) bool {
	var tenantErr *TenantError
	return errors.As(err, &tenantErr) && tenantErr.Type == ErrorTypeNotFound
}

func IsConflictError(err error) bool {
	var tenantErr *TenantError
	return errors.As(err, &tenantErr) && tenantErr.Type == ErrorTypeConflict
}

func IsUnauthorizedError(err error) bool {
	var tenantErr *TenantError
	return errors.As(err, &tenantErr) && tenantErr.Type == ErrorTypeUnauthorized
}

func IsForbiddenError(err error) bool {
	var tenantErr *TenantError
	return errors.As(err, &tenantErr) && tenantErr.Type == ErrorTypeForbidden
}

func IsInternalError(err error) bool {
	var tenantErr *TenantError
	return errors.As(err, &tenantErr) && tenantErr.Type == ErrorTypeInternal
}