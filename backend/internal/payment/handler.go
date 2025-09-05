package payment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CreatePayment handles POST /payments
func (h *Handler) CreatePayment(c *gin.Context) {

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	response, err := h.service.CreatePayment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// GetPayment handles GET /payments/:id
func (h *Handler) GetPayment(c *gin.Context) {

	paymentID := c.Param("id")
	
	payment, err := h.service.GetPayment(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

// ListPayments handles GET /payments
func (h *Handler) ListPayments(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse query parameters
	req := &ListPaymentsRequest{
		Status: c.Query("status"),
		Method: c.Query("method"),
		View:   c.Query("view"),
	}

	if offset, err := strconv.Atoi(c.DefaultQuery("offset", "0")); err == nil {
		req.Offset = offset
	}
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		req.Limit = limit
	}

	response, err := h.service.ListPayments(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePayment handles PATCH /payments/:id
func (h *Handler) UpdatePayment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	payment, err := h.service.UpdatePayment(c.Request.Context(), id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetPaymentMethods handles GET /payments/methods
func (h *Handler) GetPaymentMethods(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	methods, err := h.service.GetPaymentMethods(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_methods": methods})
}

// UpdatePaymentMethod handles PATCH /payments/methods/:id
func (h *Handler) UpdatePaymentMethod(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment method ID is required"})
		return
	}

	var req UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	method, err := h.service.UpdatePaymentMethod(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, method)
}

// PaymentWebhook handles POST /webhooks/payment/:provider
func (h *Handler) PaymentWebhook(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider is required"})
		return
	}

	switch provider {
	case "sslcommerz":
		h.handleSSLCommerzWebhook(c)
	case "bkash":
		h.handleBkashWebhook(c)
	case "nagad":
		h.handleNagadWebhook(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported payment provider"})
	}
}

func (h *Handler) handleSSLCommerzWebhook(c *gin.Context) {
	var ipnData SSLCommerzIPNResponse
	if err := c.ShouldBindJSON(&ipnData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.ValidateSSLCommerzPayment(c.Request.Context(), &ipnData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *Handler) handleBkashWebhook(c *gin.Context) {
	// TODO: Implement bKash webhook handling
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *Handler) handleNagadWebhook(c *gin.Context) {
	// TODO: Implement Nagad webhook handling
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// RegisterRoutes registers all payment routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	paymentRoutes := router.Group("/payments")
	{
		paymentRoutes.POST("", h.CreatePayment)                    // POST /payments
		paymentRoutes.GET("", h.ListPayments)                     // GET /payments
		paymentRoutes.GET("/:id", h.GetPayment)                   // GET /payments/:id
		paymentRoutes.PATCH("/:id", h.UpdatePayment)              // PATCH /payments/:id
		paymentRoutes.GET("/methods", h.GetPaymentMethods)        // GET /payments/methods
		paymentRoutes.PATCH("/methods/:id", h.UpdatePaymentMethod) // PATCH /payments/methods/:id
	}

	webhookRoutes := router.Group("/webhooks")
	{
		webhookRoutes.POST("/payment/:provider", h.PaymentWebhook) // POST /webhooks/payment/:provider
	}
}
