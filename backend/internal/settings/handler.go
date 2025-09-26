package settings

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ecommerce-saas/internal/shared/handlers"
)

// Handler handles HTTP requests for settings
type Handler struct {
	service Service
}

// NewHandler creates a new settings handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers settings routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Settings management routes (authenticated)
	settings := router.Group("/settings")
	{
		settings.GET("", h.GetSettings)      // GET /settings
		settings.PATCH("", h.UpdateSettings) // PATCH /settings
	}

	// Note: Public settings routes are registered separately in routes.go setupPublicProductRoutes
	// to avoid duplicate registration conflicts
}

// GetSettings handles GET /settings
func (h *Handler) GetSettings(c *gin.Context) {
	// Get tenant ID from context
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		handlers.HandleError(c, fmt.Errorf("tenant ID not found"))
		return
	}

	tenantIDUUID, ok := tenantID.(uuid.UUID)
	if !ok {
		// Try to convert from string
		if tenantIDStr, ok := tenantID.(string); ok {
			if parsedUUID, err := uuid.Parse(tenantIDStr); err == nil {
				tenantIDUUID = parsedUUID
			} else {
				handlers.HandleError(c, fmt.Errorf("invalid tenant ID format"))
				return
			}
		} else {
			handlers.HandleError(c, fmt.Errorf("invalid tenant ID type"))
			return
		}
	}

	// Parse query parameters
	var req GetSettingsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid query parameters: %w", err))
		return
	}

	// Call service
	response, err := h.service.GetSettings(c.Request.Context(), tenantIDUUID, &req)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("failed to get settings: %w", err))
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK,response)
}

// UpdateSettings handles PATCH /settings
func (h *Handler) UpdateSettings(c *gin.Context) {
	// Get tenant ID from context
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		handlers.HandleError(c, fmt.Errorf("tenant ID not found"))
		return
	}

	tenantIDUUID, ok := tenantID.(uuid.UUID)
	if !ok {
		// Try to convert from string
		if tenantIDStr, ok := tenantID.(string); ok {
			if parsedUUID, err := uuid.Parse(tenantIDStr); err == nil {
				tenantIDUUID = parsedUUID
			} else {
				handlers.HandleError(c, fmt.Errorf("invalid tenant ID format"))
				return
			}
		} else {
			handlers.HandleError(c, fmt.Errorf("invalid tenant ID type"))
			return
		}
	}

	// Parse request body
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	// Validate request
	if len(req.Settings) == 0 {
		handlers.HandleError(c, fmt.Errorf("no settings provided"))
		return
	}

	// Call service
	response, err := h.service.UpdateSettings(c.Request.Context(), tenantIDUUID, &req)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("failed to update settings: %w", err))
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK,response)
}

// GetPublicSettings handles GET /public/settings
func (h *Handler) GetPublicSettings(c *gin.Context) {
	// Get tenant ID from context (should be set by tenant resolution middleware)
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		handlers.HandleError(c, fmt.Errorf("tenant ID not found"))
		return
	}

	tenantIDUUID, ok := tenantID.(uuid.UUID)
	if !ok {
		// Try to convert from string
		if tenantIDStr, ok := tenantID.(string); ok {
			if parsedUUID, err := uuid.Parse(tenantIDStr); err == nil {
				tenantIDUUID = parsedUUID
			} else {
				handlers.HandleError(c, fmt.Errorf("invalid tenant ID format"))
				return
			}
		} else {
			handlers.HandleError(c, fmt.Errorf("invalid tenant ID type"))
			return
		}
	}

	// Call service
	response, err := h.service.GetPublicSettings(c.Request.Context(), tenantIDUUID)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("failed to get public settings: %w", err))
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK,response)
}

// GetHandler returns the handler instance
func (h *Handler) GetHandler() *Handler {
	return h
}
