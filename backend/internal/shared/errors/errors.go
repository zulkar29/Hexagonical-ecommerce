package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Error codes
const (
	CodeValidation      = "VALIDATION_ERROR"
	CodeNotFound        = "NOT_FOUND"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeConflict        = "CONFLICT"
	CodeInternal        = "INTERNAL_ERROR"
	CodeBadRequest      = "BAD_REQUEST"
	CodeTooManyRequests = "TOO_MANY_REQUESTS"
)

// AppError represents a custom application error
type AppError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	StatusCode int         `json:"status_code"`
	Err        error       `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Predefined error constructors

// NewValidationError creates a validation error
func NewValidationError(message string, details interface{}) *AppError {
	return &AppError{
		Code:       CodeValidation,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       CodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Unauthorized access"
	}
	return &AppError{
		Code:       CodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "Access forbidden"
	}
	return &AppError{
		Code:       CodeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       CodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string, err error) *AppError {
	if message == "" {
		message = "Internal server error"
	}
	return &AppError{
		Code:       CodeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       CodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewTooManyRequestsError creates a rate limit error
func NewTooManyRequestsError(message string) *AppError {
	if message == "" {
		message = "Too many requests"
	}
	return &AppError{
		Code:       CodeTooManyRequests,
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

// Helper functions

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError extracts an AppError from an error
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// ToAppError converts any error to an AppError
func ToAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	if appErr := GetAppError(err); appErr != nil {
		return appErr
	}

	// Convert common errors
	if errors.Is(err, errors.New("record not found")) {
		return NewNotFoundError("Resource")
	}

	// Default to internal error
	return NewInternalError("Internal server error", err)
}

// Common domain errors

var (
	// Tenant errors
	ErrTenantNotFound       = NewNotFoundError("Tenant")
	ErrSubdomainTaken       = NewConflictError("Subdomain already taken")
	ErrCustomDomainTaken    = NewConflictError("Custom domain already taken")
	ErrTenantInactive       = NewForbiddenError("Tenant is inactive")

	// Product errors
	ErrProductNotFound      = NewNotFoundError("Product")
	ErrCategoryNotFound     = NewNotFoundError("Category")
	ErrInsufficientStock    = NewBadRequestError("Insufficient stock")
	ErrProductSlugTaken     = NewConflictError("Product slug already taken")

	// Order errors
	ErrOrderNotFound        = NewNotFoundError("Order")
	ErrOrderNotEditable     = NewBadRequestError("Order cannot be modified")
	ErrOrderNotCancellable  = NewBadRequestError("Order cannot be cancelled")
	ErrOrderAlreadyPaid     = NewConflictError("Order is already paid")

	// User errors
	ErrUserNotFound         = NewNotFoundError("User")
	ErrEmailTaken           = NewConflictError("Email already taken")
	ErrInvalidCredentials   = NewUnauthorizedError("Invalid email or password")
	ErrAccountInactive      = NewForbiddenError("Account is inactive")
	ErrInvalidToken         = NewUnauthorizedError("Invalid or expired token")

	// Permission errors
	ErrPermissionDenied     = NewForbiddenError("Permission denied")
	ErrInvalidRole          = NewBadRequestError("Invalid role")

	// Payment errors
	ErrPaymentFailed        = NewBadRequestError("Payment processing failed")
	ErrRefundFailed         = NewBadRequestError("Refund processing failed")
	ErrInvalidPaymentMethod = NewBadRequestError("Invalid payment method")
)

// TODO: Add more error handling utilities
// - ErrorResponse struct for consistent API responses
// - ErrorMiddleware for Gin
// - Error logging integration
// - Error metrics collection
