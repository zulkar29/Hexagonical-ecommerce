package social

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles social commerce HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new social handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ConnectInstagram connects Instagram Business account
func (h *Handler) ConnectInstagram(c *gin.Context) {
	var req ConnectPlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	req.Platform = PlatformInstagram
	tenantID := getTenantIDFromContext(c)

	integration, err := h.service.ConnectPlatform(c.Request.Context(), tenantID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect Instagram", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, integration)
}

// ConnectFacebook connects Facebook Shop
func (h *Handler) ConnectFacebook(c *gin.Context) {
	var req ConnectPlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	req.Platform = PlatformFacebook
	tenantID := getTenantIDFromContext(c)

	integration, err := h.service.ConnectPlatform(c.Request.Context(), tenantID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect Facebook", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, integration)
}

// GetIntegrationStatus gets social integration status
func (h *Handler) GetIntegrationStatus(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	status, err := h.service.GetIntegrationStatus(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get integration status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"integrations": status})
}

// DisconnectPlatform disconnects social platform
func (h *Handler) DisconnectPlatform(c *gin.Context) {
	platformStr := c.Query("platform")
	if platformStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Platform parameter is required"})
		return
	}

	platform := Platform(platformStr)
	tenantID := getTenantIDFromContext(c)

	err := h.service.DisconnectPlatform(c.Request.Context(), tenantID, platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disconnect platform", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Platform disconnected successfully"})
}

// SyncProducts syncs products to social platforms
func (h *Handler) SyncProducts(c *gin.Context) {
	var req SyncProductsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Get platform from query parameter if not in body
	if req.Platform == "" {
		platformStr := c.Query("platform")
		if platformStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Platform is required"})
			return
		}
		req.Platform = Platform(platformStr)
	}

	// Check if syncing all platforms
	if c.Query("platform") == "all" {
		platforms := []Platform{PlatformInstagram, PlatformFacebook, PlatformTikTok}
		tenantID := getTenantIDFromContext(c)

		for _, platform := range platforms {
			req.Platform = platform
			// Don't fail if one platform fails
			h.service.SyncProducts(c.Request.Context(), tenantID, &req)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product sync initiated for all platforms"})
		return
	}

	tenantID := getTenantIDFromContext(c)

	err := h.service.SyncProducts(c.Request.Context(), tenantID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync products", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product sync initiated successfully"})
}

// GetSyncStatus gets product sync status
func (h *Handler) GetSyncStatus(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	var platform *Platform
	platformStr := c.Query("platform")
	if platformStr != "" {
		p := Platform(platformStr)
		platform = &p
	}

	status, err := h.service.GetSyncStatus(c.Request.Context(), tenantID, platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sync status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": status})
}

// UpdateProductSettings updates social product settings
func (h *Handler) UpdateProductSettings(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Determine action from query parameter
	action := c.Query("action")
	tenantID := getTenantIDFromContext(c)

	switch action {
	case "enable":
		h.handleEnableDisableProduct(c, tenantID, productID, true)
	case "disable":
		h.handleEnableDisableProduct(c, tenantID, productID, false)
	case "update_tags":
		h.handleUpdateTags(c, tenantID, productID)
	default:
		h.handleUpdateProduct(c, tenantID, productID)
	}
}

// handleUpdateProduct handles general product updates
func (h *Handler) handleUpdateProduct(c *gin.Context, tenantID, productID uuid.UUID) {
	platformStr := c.Query("platform")
	if platformStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Platform parameter is required"})
		return
	}

	platform := Platform(platformStr)

	var req UpdateSocialProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	socialProduct, err := h.service.UpdateSocialProduct(c.Request.Context(), tenantID, productID, platform, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product settings", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, socialProduct)
}

// handleEnableDisableProduct handles enabling/disabling products
func (h *Handler) handleEnableDisableProduct(c *gin.Context, tenantID, productID uuid.UUID, enabled bool) {
	platformStr := c.Query("platform")
	if platformStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Platform parameter is required"})
		return
	}

	platform := Platform(platformStr)
	req := UpdateSocialProductRequest{IsEnabled: &enabled}

	socialProduct, err := h.service.UpdateSocialProduct(c.Request.Context(), tenantID, productID, platform, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, socialProduct)
}

// handleUpdateTags handles updating product tags
func (h *Handler) handleUpdateTags(c *gin.Context, tenantID, productID uuid.UUID) {
	platformStr := c.Query("platform")
	if platformStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Platform parameter is required"})
		return
	}

	platform := Platform(platformStr)

	var tagsReq struct {
		Tags []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&tagsReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	req := UpdateSocialProductRequest{Tags: tagsReq.Tags}

	socialProduct, err := h.service.UpdateSocialProduct(c.Request.Context(), tenantID, productID, platform, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tags", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, socialProduct)
}

// GetAnalytics gets social commerce analytics
func (h *Handler) GetAnalytics(c *gin.Context) {
	platformStr := c.Query("platform")
	if platformStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Platform parameter is required"})
		return
	}

	platform := Platform(platformStr)
	tenantID := getTenantIDFromContext(c)

	// Parse date range
	dateFromStr := c.DefaultQuery("date_from", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	dateToStr := c.DefaultQuery("date_to", time.Now().Format("2006-01-02"))

	dateFrom, err := time.Parse("2006-01-02", dateFromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_from format. Use YYYY-MM-DD"})
		return
	}

	dateTo, err := time.Parse("2006-01-02", dateToStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_to format. Use YYYY-MM-DD"})
		return
	}

	analytics, err := h.service.GetAnalytics(c.Request.Context(), tenantID, platform, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// RegisterRoutes registers social commerce routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	social := r.Group("/integrations")
	{
		social.POST("/instagram/connect", h.ConnectInstagram)
		social.POST("/facebook/connect", h.ConnectFacebook)
		social.GET("/social/status", h.GetIntegrationStatus)
		social.POST("/social/disconnect", h.DisconnectPlatform)
		social.POST("/social/product-sync", h.SyncProducts)
		social.GET("/social/sync-status", h.GetSyncStatus)
		social.PATCH("/social/products/:id", h.UpdateProductSettings)
		social.GET("/social/analytics", h.GetAnalytics)
	}
}

// Helper function to get tenant ID from context
func getTenantIDFromContext(c *gin.Context) uuid.UUID {
	// TODO: Extract tenant ID from JWT token or middleware
	// For now, return a dummy UUID
	return uuid.New()
}