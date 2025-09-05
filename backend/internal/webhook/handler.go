package webhook

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Helper function to extract tenant ID from context
func (h *Handler) getTenantID(c *gin.Context) (uuid.UUID, error) {
	tenantIDStr, exists := c.Get("tenant_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("tenant ID not found in context")
	}
	
	tenantID, ok := tenantIDStr.(uuid.UUID)
	if !ok {
		// Try to parse as string
		if str, ok := tenantIDStr.(string); ok {
			return uuid.Parse(str)
		}
		return uuid.Nil, fmt.Errorf("invalid tenant ID format")
	}
	
	return tenantID, nil
}

// Webhook Endpoint Management Handlers

type CreateEndpointRequest struct {
	URL         string   `json:"url" binding:"required,url"`
	Events      []string `json:"events" binding:"required"`
	Description string   `json:"description"`
	IsActive    bool     `json:"is_active"`
}

type UpdateEndpointRequest struct {
	URL         string   `json:"url" binding:"omitempty,url"`
	Events      []string `json:"events"`
	Description string   `json:"description"`
	IsActive    *bool    `json:"is_active"`
}

func (h *Handler) CreateEndpoint(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	var req CreateEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Convert string events to WebhookEvent type
	events := make([]WebhookEvent, len(req.Events))
	for i, event := range req.Events {
		events[i] = WebhookEvent(event)
	}
	
	endpoint := &WebhookEndpoint{
		TenantID:    tenantID,
		URL:         req.URL,
		Events:      events,
		Description: req.Description,
		IsActive:    req.IsActive,
	}
	
	createdEndpoint, err := h.service.CreateEndpoint(tenantID, endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, createdEndpoint)
}

func (h *Handler) UpdateEndpoint(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	endpointID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endpoint ID"})
		return
	}
	
	var req UpdateEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get existing endpoint
	endpoint, err := h.service.GetEndpoint(tenantID, endpointID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
		return
	}
	
	// Update fields
	if req.URL != "" {
		endpoint.URL = req.URL
	}
	if req.Events != nil {
		// Convert string events to WebhookEvent type
		events := make([]WebhookEvent, len(req.Events))
		for i, event := range req.Events {
			events[i] = WebhookEvent(event)
		}
		endpoint.Events = events
	}
	if req.Description != "" {
		endpoint.Description = req.Description
	}
	if req.IsActive != nil {
		endpoint.IsActive = *req.IsActive
	}
	
	updatedEndpoint, err := h.service.UpdateEndpoint(tenantID, endpointID, endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, updatedEndpoint)
}

func (h *Handler) DeleteEndpoint(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	endpointID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endpoint ID"})
		return
	}
	
	if err := h.service.DeleteEndpoint(tenantID, endpointID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) GetEndpoints(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	endpoints, err := h.service.GetEndpoints(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"endpoints": endpoints})
}

func (h *Handler) GetEndpoint(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	endpointID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endpoint ID"})
		return
	}
	
	endpoint, err := h.service.GetEndpoint(tenantID, endpointID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
		return
	}
	
	c.JSON(http.StatusOK, endpoint)
}

func (h *Handler) TestEndpoint(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	endpointID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endpoint ID"})
		return
	}
	
	result, err := h.service.TestEndpoint(tenantID, endpointID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetDeliveries(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	// Parse query parameters
	endpointIDStr := c.Query("endpoint_id")
	var endpointID uuid.UUID
	if endpointIDStr != "" {
		endpointID, err = uuid.Parse(endpointIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endpoint ID"})
			return
		}
	}
	
	deliveries, err := h.service.GetDeliveries(tenantID, endpointID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
}

func (h *Handler) GetDelivery(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	deliveryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery ID"})
		return
	}
	
	delivery, err := h.service.GetDelivery(tenantID, deliveryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		return
	}
	
	c.JSON(http.StatusOK, delivery)
}

func (h *Handler) RetryDelivery(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	deliveryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery ID"})
		return
	}
	
	delivery, err := h.service.GetDelivery(tenantID, deliveryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		return
	}
	
	if err := h.service.ScheduleRetry(delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Delivery retry scheduled"})
}

func (h *Handler) GetEndpointStats(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	// Parse date range
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	
	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)
	
	stats, err := h.service.GetDeliveryStats(tenantID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetWebhookLogs(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)
	
	failedDeliveries, err := h.service.GetFailedDeliveries(tenantID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"logs": failedDeliveries})
}

// Incoming Webhook Handlers

func (h *Handler) StripeWebhook(c *gin.Context) {
	// Extract tenant ID from header or path
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		tenantIDStr = c.Param("tenant_id")
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	// Get signature from header
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}
	
	// Read body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	
	// Process webhook
	if err := h.service.ProcessStripeWebhook(tenantID, signature, body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *Handler) PayPalWebhook(c *gin.Context) {
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		tenantIDStr = c.Param("tenant_id")
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	signature := c.GetHeader("PAYPAL-TRANSMISSION-SIG")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	
	if err := h.service.ProcessPayPalWebhook(tenantID, signature, body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *Handler) BkashWebhook(c *gin.Context) {
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		tenantIDStr = c.Param("tenant_id")
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	signature := c.GetHeader("X-Bkash-Signature")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	
	if err := h.service.ProcessBkashWebhook(tenantID, signature, body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *Handler) NagadWebhook(c *gin.Context) {
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		tenantIDStr = c.Param("tenant_id")
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	signature := c.GetHeader("X-Nagad-Signature")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	
	if err := h.service.ProcessNagadWebhook(tenantID, signature, body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *Handler) PathaoWebhook(c *gin.Context) {
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		tenantIDStr = c.Param("tenant_id")
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	signature := c.GetHeader("X-Pathao-Signature")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	
	if err := h.service.ProcessPathaoWebhook(tenantID, signature, body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *Handler) RedXWebhook(c *gin.Context) {
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		tenantIDStr = c.Param("tenant_id")
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	signature := c.GetHeader("X-RedX-Signature")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	
	if err := h.service.ProcessRedXWebhook(tenantID, signature, body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *Handler) PaperflyWebhook(c *gin.Context) {
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		tenantIDStr = c.Param("tenant_id")
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	signature := c.GetHeader("X-Paperfly-Signature")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	
	if err := h.service.ProcessPaperflyWebhook(tenantID, signature, body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"received": true})
}

// Route Registration

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Webhook endpoint management
	webhooks := router.Group("/webhooks")
	{
		// Endpoint management
		webhooks.POST("/endpoints", h.CreateEndpoint)
		webhooks.PUT("/endpoints/:id", h.UpdateEndpoint)
		webhooks.DELETE("/endpoints/:id", h.DeleteEndpoint)
		webhooks.GET("/endpoints", h.GetEndpoints)
		webhooks.GET("/endpoints/:id", h.GetEndpoint)
		webhooks.POST("/endpoints/:id/test", h.TestEndpoint)
		
		// Delivery management
		webhooks.GET("/deliveries", h.GetDeliveries)
		webhooks.GET("/deliveries/:id", h.GetDelivery)
		webhooks.POST("/deliveries/:id/retry", h.RetryDelivery)
		
		// Analytics and monitoring
		webhooks.GET("/stats", h.GetEndpointStats)
		webhooks.GET("/logs", h.GetWebhookLogs)
	}
	
	// Incoming webhooks (public endpoints)
	incoming := router.Group("/webhooks/incoming")
	{
		// Payment providers
		incoming.POST("/stripe", h.StripeWebhook)
		incoming.POST("/paypal", h.PayPalWebhook)
		incoming.POST("/bkash", h.BkashWebhook)
		incoming.POST("/nagad", h.NagadWebhook)
		
		// Shipping providers
		incoming.POST("/pathao", h.PathaoWebhook)
		incoming.POST("/redx", h.RedXWebhook)
		incoming.POST("/paperfly", h.PaperflyWebhook)
	}
	
	// Tenant-specific incoming webhooks
	tenantIncoming := router.Group("/tenants/:tenant_id/webhooks/incoming")
	{
		// Payment providers
		tenantIncoming.POST("/stripe", h.StripeWebhook)
		tenantIncoming.POST("/paypal", h.PayPalWebhook)
		tenantIncoming.POST("/bkash", h.BkashWebhook)
		tenantIncoming.POST("/nagad", h.NagadWebhook)
		
		// Shipping providers
		tenantIncoming.POST("/pathao", h.PathaoWebhook)
		tenantIncoming.POST("/redx", h.RedXWebhook)
		tenantIncoming.POST("/paperfly", h.PaperflyWebhook)
	}
}

// Middleware functions

func (h *Handler) setupMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Add common headers
		c.Header("Content-Type", "application/json")
		c.Next()
	})
}

func (h *Handler) validateWebhookSignature() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Skip validation for test endpoints
		if c.Request.URL.Path == "/test" {
			c.Next()
			return
		}
		
		// Get signature from header
		signature := c.GetHeader("X-Webhook-Signature")
		if signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing webhook signature"})
			c.Abort()
			return
		}
		
		// Validate signature (implementation depends on provider)
		c.Next()
	})
}

func (h *Handler) rateLimitMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Get tenant ID from context or header
		tenantID, err := h.getTenantID(c)
		if err != nil {
			// Skip rate limiting if no tenant ID available
			c.Next()
			return
		}
		
		// For rate limiting, we'll use a default endpoint ID or skip if not available
		endpointIDStr := c.Param("id")
		if endpointIDStr == "" {
			// Skip rate limiting if no endpoint ID available
			c.Next()
			return
		}
		
		endpointID, err := uuid.Parse(endpointIDStr)
		if err != nil {
			// Skip rate limiting if invalid endpoint ID
			c.Next()
			return
		}
		
		// Check rate limit
		if h.service.IsRateLimited(tenantID, endpointID) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		
		// Increment rate limit counter
		h.service.IncrementRateLimit(tenantID, endpointID)
		c.Next()
	})
}

func (h *Handler) loggingMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Log request details
		start := time.Now()
		
		c.Next()
		
		// Log response details
		latency := time.Since(start)
		status := c.Writer.Status()
		
		// Log to service (simplified)
		fmt.Printf("[WEBHOOK] %s %s - %d (%v)\n", 
			c.Request.Method, 
			c.Request.URL.Path, 
			status, 
			latency)
	})
}
