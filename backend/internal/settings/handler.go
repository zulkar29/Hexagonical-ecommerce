package settings

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		settings.GET("", h.GetSettings)           // GET /settings
		settings.PATCH("", h.UpdateSettings)      // PATCH /settings
	}

	// Note: Public settings routes are registered separately in routes.go setupPublicProductRoutes
	// to avoid duplicate registration conflicts
}

// GetSettings handles GET /settings
func (h *Handler) GetSettings(c *gin.Context) {
	// Get tenant ID from context
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID not found"})
		return
	}

	tenantIDUint, ok := tenantID.(uint)
	if !ok {
		// Try to convert from string
		if tenantIDStr, ok := tenantID.(string); ok {
			if id, err := strconv.ParseUint(tenantIDStr, 10, 32); err == nil {
				tenantIDUint = uint(id)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID format"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID type"})
			return
		}
	}

	// Parse query parameters
	var req GetSettingsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	// Call service
	response, err := h.service.GetSettings(c.Request.Context(), tenantIDUint, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get settings", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSettings handles PATCH /settings
func (h *Handler) UpdateSettings(c *gin.Context) {
	// Get tenant ID from context
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID not found"})
		return
	}

	tenantIDUint, ok := tenantID.(uint)
	if !ok {
		// Try to convert from string
		if tenantIDStr, ok := tenantID.(string); ok {
			if id, err := strconv.ParseUint(tenantIDStr, 10, 32); err == nil {
				tenantIDUint = uint(id)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID format"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID type"})
			return
		}
	}

	// Parse request body
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Validate request
	if len(req.Settings) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No settings provided"})
		return
	}

	// Call service
	response, err := h.service.UpdateSettings(c.Request.Context(), tenantIDUint, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPublicSettings handles GET /public/settings
func (h *Handler) GetPublicSettings(c *gin.Context) {
	// Get tenant ID from context (should be set by tenant resolution middleware)
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID not found"})
		return
	}

	tenantIDUint, ok := tenantID.(uint)
	if !ok {
		// Try to convert from string
		if tenantIDStr, ok := tenantID.(string); ok {
			if id, err := strconv.ParseUint(tenantIDStr, 10, 32); err == nil {
				tenantIDUint = uint(id)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID format"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID type"})
			return
		}
	}

	// Call service
	response, err := h.service.GetPublicSettings(c.Request.Context(), tenantIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get public settings", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetHandler returns the handler instance
func (h *Handler) GetHandler() *Handler {
	return h
}