package payment

import (
	"fmt"
	"net/http"
	"strconv"

	"ecommerce-saas/internal/shared/handlers"
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
		handlers.HandleValidationError(c, err)
		return
	}

	response, err := h.service.CreatePayment(c.Request.Context(), &req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithCreated(c, gin.H{"data": response})
}

// GetPayment handles GET /payments/:id
func (h *Handler) GetPayment(c *gin.Context) {

	paymentID := c.Param("id")
	
	payment, err := h.service.GetPayment(c.Request.Context(), paymentID)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("payment not found"))
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"data": payment})
}

// ListPayments handles GET /payments
func (h *Handler) ListPayments(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		handlers.HandleError(c, fmt.Errorf("user not authenticated"))
		return
	}

	// Parse query parameters
	req := &ListPaymentsRequest{
		Status: c.Query("status"),
		Method: c.Query("method"),
		View:   c.Query("view"),
	}

	if offset, offsetErr := strconv.Atoi(c.DefaultQuery("offset", "0")); offsetErr == nil {
		req.Offset = offset
	}
	if limit, limitErr := strconv.Atoi(c.DefaultQuery("limit", "20")); limitErr == nil {
		req.Limit = limit
	}

	response, err := h.service.ListPayments(c.Request.Context(), req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, response)
}

// UpdatePayment handles PATCH /payments/:id
func (h *Handler) UpdatePayment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handlers.HandleError(c, fmt.Errorf("payment ID is required"))
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	payment, err := h.service.UpdatePayment(c.Request.Context(), id, updates)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, payment)
}

// GetPaymentMethods handles GET /payments/methods
func (h *Handler) GetPaymentMethods(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		handlers.HandleError(c, fmt.Errorf("user not authenticated"))
		return
	}

	methods, err := h.service.GetPaymentMethods(c.Request.Context(), userID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"payment_methods": methods})
}

// UpdatePaymentMethod handles PATCH /payments/methods/:id
func (h *Handler) UpdatePaymentMethod(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handlers.HandleError(c, fmt.Errorf("payment method ID is required"))
		return
	}

	var req UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	method, err := h.service.UpdatePaymentMethod(c.Request.Context(), id, &req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, method)
}

// PaymentWebhook handles POST /webhooks/payment/:provider
func (h *Handler) PaymentWebhook(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		handlers.HandleError(c, fmt.Errorf("provider is required"))
		return
	}

	switch provider {
	case "sslcommerz":
		h.handleSSLCommerzWebhook(c)

	default:
		handlers.HandleError(c, fmt.Errorf("unsupported payment provider"))
	}
}

func (h *Handler) handleSSLCommerzWebhook(c *gin.Context) {
	var ipnData SSLCommerzIPNResponse
	if err := c.ShouldBindJSON(&ipnData); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	if err := h.service.ValidateSSLCommerzPayment(c.Request.Context(), &ipnData); err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"status": "success"})
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
