package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	sharedErrors "ecommerce-saas/internal/shared/errors"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

// HandleError provides centralized error handling for HTTP responses
func HandleError(c *gin.Context, err error) {
	switch {
	case err == gorm.ErrRecordNotFound:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Resource not found",
			Message: "The requested resource could not be found",
			Code:    "NOT_FOUND",
		})
	case sharedErrors.IsValidationError(err):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    "VALIDATION_ERROR",
		})
	case sharedErrors.IsUnauthorizedError(err):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: err.Error(),
			Code:    "UNAUTHORIZED",
		})
	case sharedErrors.IsForbiddenError(err):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Forbidden",
			Message: err.Error(),
			Code:    "FORBIDDEN",
		})
	case sharedErrors.IsConflictError(err):
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "Conflict",
			Message: err.Error(),
			Code:    "CONFLICT",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: "An unexpected error occurred",
			Code:    "INTERNAL_ERROR",
		})
	}
}

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *gin.Context, statusCode int, data any, message ...string) {
	response := SuccessResponse{
		Data: data,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	c.JSON(statusCode, response)
}

// RespondWithSuccessAndMeta sends a success response with metadata
func RespondWithSuccessAndMeta(c *gin.Context, statusCode int, data any, meta any, message ...string) {
	response := SuccessResponse{
		Data: data,
		Meta: meta,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	c.JSON(statusCode, response)
}

// RespondWithCreated sends a standardized creation success response
func RespondWithCreated(c *gin.Context, data any, message ...string) {
	msg := "Resource created successfully"
	if len(message) > 0 {
		msg = message[0]
	}
	RespondWithSuccess(c, http.StatusCreated, data, msg)
}

// RespondWithUpdated sends a standardized update success response
func RespondWithUpdated(c *gin.Context, data any, message ...string) {
	msg := "Resource updated successfully"
	if len(message) > 0 {
		msg = message[0]
	}
	RespondWithSuccess(c, http.StatusOK, data, msg)
}

// RespondWithDeleted sends a standardized deletion success response
func RespondWithDeleted(c *gin.Context, message ...string) {
	msg := "Resource deleted successfully"
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Message: msg,
	})
}

// RespondWithValidationError responds with validation error details
func RespondWithValidationError(c *gin.Context, fieldErrors map[string]string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "Validation failed",
		Message: "One or more fields contain invalid values",
		Code:    "VALIDATION_ERROR",
		Details: fieldErrors,
	})
}

// HandleValidationError is a helper for JSON binding validation errors
func HandleValidationError(c *gin.Context, err error) {
	RespondWithValidationError(c, map[string]string{
		"request": err.Error(),
	})
}

// ExtractTenantID extracts tenant ID from context without side effects
func ExtractTenantID(c *gin.Context) (uuid.UUID, error) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		return uuid.Nil, sharedErrors.NewBadRequestError("Tenant context is missing from request")
	}

	if id, ok := tenantID.(uuid.UUID); ok {
		return id, nil
	}

	return uuid.Nil, sharedErrors.NewBadRequestError("Tenant context format is invalid")
}

// ExtractUserID extracts user ID from context without side effects
func ExtractUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, sharedErrors.NewUnauthorizedError("Authentication required for this operation")
	}

	if id, ok := userID.(uuid.UUID); ok {
		return id, nil
	}

	return uuid.Nil, sharedErrors.NewUnauthorizedError("User context format is invalid")
}

// RequireTenantID extracts tenant ID and handles errors automatically
func RequireTenantID(c *gin.Context) (uuid.UUID, bool) {
	tenantID, err := ExtractTenantID(c)
	if err != nil {
		HandleError(c, err)
		return uuid.Nil, false
	}
	return tenantID, true
}

// RequireUserID extracts user ID and handles errors automatically
func RequireUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, err := ExtractUserID(c)
	if err != nil {
		HandleError(c, err)
		return uuid.Nil, false
	}
	return userID, true
}
