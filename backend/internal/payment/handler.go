package payment

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

// CreatePayment handles POST /payments
func (h *Handler) CreatePayment(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	response, err := h.service.CreatePayment(tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// ProcessPayment handles POST /payments/:id/process
func (h *Handler) ProcessPayment(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	paymentID := c.Param("id")
	
	var req ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	req.PaymentID = paymentID

	err := h.service.ProcessPayment(tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}

// GetPayment handles GET /payments/:id
func (h *Handler) GetPayment(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	paymentID := c.Param("id")
	
	payment, err := h.service.GetPayment(tenantID.(uuid.UUID), paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

// ListPayments handles GET /payments
func (h *Handler) ListPayments(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	// Parse query parameters
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")
	orderIDStr := c.Query("order_id")

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)

	var orderID *uuid.UUID
	if orderIDStr != "" {
		if id, err := uuid.Parse(orderIDStr); err == nil {
			orderID = &id
		}
	}

	payments, total, err := h.service.ListPayments(tenantID.(uuid.UUID), orderID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": payments,
		"total": total,
		"offset": offset,
		"limit": limit,
	})
}

// RefundPayment handles POST /payments/:id/refund
func (h *Handler) RefundPayment(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	paymentID := c.Param("id")
	
	var req RefundPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	req.PaymentID = paymentID

	err := h.service.RefundPayment(tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Refund processed successfully"})
}

// SSLCommerzWebhook handles POST /webhooks/sslcommerz
func (h *Handler) SSLCommerzWebhook(c *gin.Context) {
	var ipnData SSLCommerzIPNResponse
	if err := c.ShouldBindJSON(&ipnData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook data"})
		return
	}

	err := h.service.ValidateSSLCommerzPayment(&ipnData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process the payment based on IPN data
	// This would typically update the payment status
	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

// RegisterRoutes registers all payment routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	payments := r.Group("/payments")
	{
		payments.POST("", h.CreatePayment)
		payments.GET("", h.ListPayments)
		payments.GET("/:id", h.GetPayment)
		payments.POST("/:id/process", h.ProcessPayment)
		payments.POST("/:id/refund", h.RefundPayment)
	}

	// Webhook routes (usually without authentication)
	r.POST("/webhooks/sslcommerz", h.SSLCommerzWebhook)
}
