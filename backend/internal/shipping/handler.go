package shipping

import (
	"net/http"
	"strconv"

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

// Shipping Zones

func (h *Handler) CreateShippingZone(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var req CreateShippingZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	zone, err := h.service.CreateShippingZone(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Shipping zone created successfully",
		"data":    zone,
	})
}

func (h *Handler) GetShippingZones(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	zones, err := h.service.GetShippingZones(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shipping zones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": zones})
}

func (h *Handler) GetShippingZone(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	zoneID := c.Param("id")
	zone, err := h.service.GetShippingZone(tenantID.(uuid.UUID), zoneID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping zone not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": zone})
}

func (h *Handler) UpdateShippingZone(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	zoneID := c.Param("id")
	var req UpdateShippingZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	zone, err := h.service.UpdateShippingZone(tenantID.(uuid.UUID), zoneID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shipping zone updated successfully",
		"data":    zone,
	})
}

func (h *Handler) DeleteShippingZone(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	zoneID := c.Param("id")
	err := h.service.DeleteShippingZone(tenantID.(uuid.UUID), zoneID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping zone deleted successfully"})
}

// Shipping Rates

func (h *Handler) GetShippingRates(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var req ShippingRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	rates, err := h.service.CalculateShippingRates(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rates})
}

func (h *Handler) CreateShippingRate(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var req CreateShippingRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	rate, err := h.service.CreateShippingRate(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Shipping rate created successfully",
		"data":    rate,
	})
}

// Shipping Labels

func (h *Handler) CreateShippingLabel(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var req CreateShippingLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	label, err := h.service.CreateShippingLabel(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Shipping label created successfully",
		"data":    label,
	})
}

func (h *Handler) GetShippingLabel(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	labelID := c.Param("id")
	label, err := h.service.GetShippingLabel(tenantID.(uuid.UUID), labelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping label not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": label})
}

func (h *Handler) GetShippingLabels(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	// Parse pagination
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	if limit > 100 {
		limit = 100
	}

	labels, total, err := h.service.GetShippingLabels(tenantID.(uuid.UUID), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shipping labels"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"labels": labels,
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	})
}

func (h *Handler) CancelShipment(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	labelID := c.Param("id")
	err := h.service.CancelShipment(tenantID.(uuid.UUID), labelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipment cancelled successfully"})
}

// Package Tracking

func (h *Handler) TrackPackage(c *gin.Context) {
	trackingNumber := c.Param("trackingNumber")
	
	tracking, err := h.service.TrackPackage(trackingNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tracking})
}

func (h *Handler) GetTrackingHistory(c *gin.Context) {
	trackingNumber := c.Param("trackingNumber")
	
	history, err := h.service.GetTrackingHistory(trackingNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tracking history not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}

// Address Validation

func (h *Handler) ValidateAddress(c *gin.Context) {
	var req AddressValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	result, err := h.service.ValidateAddress(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Delivery Estimates

func (h *Handler) GetDeliveryEstimate(c *gin.Context) {
	var req DeliveryEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	estimate, err := h.service.GetDeliveryEstimate(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": estimate})
}

// Provider Management

func (h *Handler) GetShippingProviders(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	providers, err := h.service.GetShippingProviders(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shipping providers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": providers})
}

func (h *Handler) ConfigureProvider(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	provider := c.Param("provider")
	var req ProviderConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.ConfigureProvider(tenantID.(uuid.UUID), provider, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider configured successfully"})
}

// Statistics

func (h *Handler) GetShippingStats(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	stats, err := h.service.GetShippingStats(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shipping statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handler) GetShippingHistory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	// Parse pagination
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	if limit > 100 {
		limit = 100
	}

	history, total, err := h.service.GetShippingHistory(tenantID.(uuid.UUID), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shipping history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"history": history,
			"total":   total,
			"offset":  offset,
			"limit":   limit,
		},
	})
}

// Provider Webhooks

func (h *Handler) PathaoWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	err := h.service.ProcessPathaoWebhook(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (h *Handler) RedXWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	err := h.service.ProcessRedXWebhook(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (h *Handler) PaperflyWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	err := h.service.ProcessPaperflyWebhook(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (h *Handler) DHLWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	err := h.service.ProcessDHLWebhook(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (h *Handler) FedExWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	err := h.service.ProcessFedExWebhook(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

// RegisterRoutes registers all shipping routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	shipping := router.Group("/shipping")
	{
		// Shipping Zones
		shipping.POST("/zones", h.CreateShippingZone)
		shipping.GET("/zones", h.GetShippingZones)
		shipping.GET("/zones/:id", h.GetShippingZone)
		shipping.PUT("/zones/:id", h.UpdateShippingZone)
		shipping.DELETE("/zones/:id", h.DeleteShippingZone)

		// Shipping Rates
		shipping.POST("/rates", h.GetShippingRates)
		shipping.POST("/rates/create", h.CreateShippingRate)

		// Shipping Labels
		shipping.POST("/labels", h.CreateShippingLabel)
		shipping.GET("/labels", h.GetShippingLabels)
		shipping.GET("/labels/:id", h.GetShippingLabel)
		shipping.DELETE("/labels/:id/cancel", h.CancelShipment)

		// Package Tracking (public)
		shipping.GET("/track/:trackingNumber", h.TrackPackage)
		shipping.GET("/track/:trackingNumber/history", h.GetTrackingHistory)

		// Address & Delivery
		shipping.POST("/validate-address", h.ValidateAddress)
		shipping.POST("/estimate", h.GetDeliveryEstimate)

		// Provider Management
		shipping.GET("/providers", h.GetShippingProviders)
		shipping.POST("/providers/:provider/configure", h.ConfigureProvider)

		// Statistics & History
		shipping.GET("/stats", h.GetShippingStats)
		shipping.GET("/history", h.GetShippingHistory)
	}

	// Webhook routes (these should be registered in a separate webhook router)
	webhooks := router.Group("/webhooks/shipping")
	{
		webhooks.POST("/pathao", h.PathaoWebhook)
		webhooks.POST("/redx", h.RedXWebhook)
		webhooks.POST("/paperfly", h.PaperflyWebhook)
		webhooks.POST("/dhl", h.DHLWebhook)
		webhooks.POST("/fedex", h.FedExWebhook)
	}
}