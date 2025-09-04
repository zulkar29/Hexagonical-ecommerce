package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// SendNotification handles POST /notifications
func (h *Handler) SendNotification(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	response, err := h.service.SendNotification(tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// SendEmail handles POST /notifications/email
func (h *Handler) SendEmail(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.SendEmail(tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

// SendSMS handles POST /notifications/sms
func (h *Handler) SendSMS(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.SendSMS(tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMS sent successfully"})
}

// GetNotification handles GET /notifications/:id
func (h *Handler) GetNotification(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	notificationID := c.Param("id")
	
	notification, err := h.service.GetNotification(tenantID.(uuid.UUID), notificationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notification})
}

// ListNotifications handles GET /notifications
func (h *Handler) ListNotifications(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	// Parse query parameters
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")
	userIDStr := c.Query("user_id")

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)

	var userID *uuid.UUID
	if userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			userID = &id
		}
	}

	notifications, total, err := h.service.ListNotifications(tenantID.(uuid.UUID), userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": notifications,
		"total": total,
		"offset": offset,
		"limit": limit,
	})
}

// MarkAsRead handles PUT /notifications/:id/read
func (h *Handler) MarkAsRead(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	notificationID := c.Param("id")
	
	err := h.service.MarkAsRead(tenantID.(uuid.UUID), notificationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// CreateTemplate handles POST /notifications/templates
func (h *Handler) CreateTemplate(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	template, err := h.service.CreateTemplate(tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": template})
}

// UpdateTemplate handles PUT /notifications/templates/:id
func (h *Handler) UpdateTemplate(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	templateID := c.Param("id")
	
	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.UpdateTemplate(tenantID.(uuid.UUID), templateID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template updated successfully"})
}

// GetTemplate handles GET /notifications/templates/:id
func (h *Handler) GetTemplate(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	templateID := c.Param("id")
	
	template, err := h.service.GetTemplate(tenantID.(uuid.UUID), templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": template})
}

// ListTemplates handles GET /notifications/templates
func (h *Handler) ListTemplates(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	notificationType := c.Query("type")
	channel := c.Query("channel")

	templates, err := h.service.ListTemplates(tenantID.(uuid.UUID), notificationType, channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": templates})
}

// GetPreferences handles GET /notifications/preferences
func (h *Handler) GetPreferences(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	preferences, err := h.service.GetPreferences(tenantID.(uuid.UUID), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Preferences not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": preferences})
}

// UpdatePreferences handles PUT /notifications/preferences
func (h *Handler) UpdatePreferences(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var req NotificationPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.UpdatePreferences(tenantID.(uuid.UUID), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preferences updated successfully"})
}

// GetStats handles GET /notifications/stats
func (h *Handler) GetStats(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	stats, err := h.service.GetStats(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// RegisterRoutes registers all notification routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	notifications := r.Group("/notifications")
	{
		// Notification management
		notifications.POST("", h.SendNotification)
		notifications.GET("", h.ListNotifications)
		notifications.GET("/:id", h.GetNotification)
		notifications.PUT("/:id/read", h.MarkAsRead)
		
		// Specific notification types
		notifications.POST("/email", h.SendEmail)
		notifications.POST("/sms", h.SendSMS)
		
		// Template management
		notifications.POST("/templates", h.CreateTemplate)
		notifications.GET("/templates", h.ListTemplates)
		notifications.GET("/templates/:id", h.GetTemplate)
		notifications.PUT("/templates/:id", h.UpdateTemplate)
		
		// User preferences
		notifications.GET("/preferences", h.GetPreferences)
		notifications.PUT("/preferences", h.UpdatePreferences)
		
		// Statistics
		notifications.GET("/stats", h.GetStats)
	}
}
