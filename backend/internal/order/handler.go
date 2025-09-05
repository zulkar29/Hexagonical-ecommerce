package order

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles order HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new order handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateOrder creates a new order
// @Summary Create a new order
// @Description Create a new order with items and calculate totals
// @Tags orders
// @Accept json
// @Produce json
// @Param order body Order true "Order data"
// @Success 201 {object} Order
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get tenant and user from context
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context required"})
		return
	}

	// Set tenant and user context
	order.TenantID = tenantID.(uuid.UUID)
	order.UserID = userID.(uuid.UUID)

	createdOrder, err := h.service.CreateOrder(c.Request.Context(), tenantID.(uuid.UUID), &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
}

// GetOrder retrieves an order by ID
// @Summary Get order by ID
// @Description Get order details including items
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} Order
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [get]
func (h *Handler) GetOrder(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	orderID := c.Param("id")
	order, err := h.service.GetOrder(tenantID.(uuid.UUID), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrderByNumber retrieves an order by order number
// @Summary Get order by order number
// @Description Get order details by order number (for customer tracking)
// @Tags orders
// @Produce json
// @Param number path string true "Order Number"
// @Success 200 {object} Order
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/number/{number} [get]
func (h *Handler) GetOrderByNumber(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	orderNumber := c.Param("number")
	order, err := h.service.GetOrderByNumber(tenantID.(uuid.UUID), orderNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListOrders lists orders with filtering and pagination
// @Summary List orders
// @Description Get paginated list of orders with optional filtering
// @Tags orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param status query string false "Order status filter"
// @Param payment_status query string false "Payment status filter"
// @Param fulfillment_status query string false "Fulfillment status filter"
// @Param customer_email query string false "Customer email filter"
// @Param order_number query string false "Order number filter"
// @Param search query string false "General search term"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /orders [get]
func (h *Handler) ListOrders(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse filters
	filter := OrderFilter{
		CustomerEmail: c.Query("customer_email"),
		OrderNumber:   c.Query("order_number"),
		Search:        c.Query("search"),
	}

	// Parse status filters
	if status := c.Query("status"); status != "" {
		orderStatus := OrderStatus(status)
		filter.Status = &orderStatus
	}

	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		payStatus := PaymentStatus(paymentStatus)
		filter.PaymentStatus = &payStatus
	}

	if fulfillmentStatus := c.Query("fulfillment_status"); fulfillmentStatus != "" {
		fulStatus := FulfillmentStatus(fulfillmentStatus)
		filter.FulfillmentStatus = &fulStatus
	}

	orders, total, err := h.service.ListOrders(tenantID.(uuid.UUID), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"orders":      orders,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}

	c.JSON(http.StatusOK, response)
}

// UpdateOrderStatus updates the status of an order
// @Summary Update order status
// @Description Update order status and related tracking information
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body map[string]interface{} true "Status update data"
// @Success 200 {object} Order
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/status [patch]
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	var req struct {
		Status         string `json:"status" binding:"required"`
		TrackingNumber string `json:"tracking_number"`
		TrackingURL    string `json:"tracking_url"`
		Notes          string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.service.UpdateOrderStatus(c.Request.Context(), tenantID.(uuid.UUID), orderID, OrderStatus(req.Status), req.TrackingNumber, req.TrackingURL, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CancelOrder cancels an order
// @Summary Cancel order
// @Description Cancel an order if it's in a cancellable state
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param body body map[string]string false "Cancellation reason"
// @Success 200 {object} Order
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/cancel [post]
func (h *Handler) CancelOrder(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	orderID := c.Param("id")
	order, err := h.service.CancelOrder(tenantID.(uuid.UUID), orderID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ProcessPayment processes payment for an order
// @Summary Process payment
// @Description Process payment for an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param payment body map[string]interface{} true "Payment data"
// @Success 200 {object} Order
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/payment [post]
func (h *Handler) ProcessPayment(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	var req struct {
		PaymentID       string `json:"payment_id" binding:"required"`
		PaymentMethodID string `json:"payment_method_id" binding:"required"`
		Confirmation    string `json:"confirmation"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	_, err = uuid.Parse(req.PaymentMethodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment method ID"})
		return
	}

	order, err := h.service.ProcessPayment(tenantID.(uuid.UUID), paymentID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// RefundOrder refunds an order
// @Summary Refund order
// @Description Refund an order and update status
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param refund body map[string]interface{} true "Refund data"
// @Success 200 {object} Order
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/refund [post]
func (h *Handler) RefundOrder(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	var req struct {
		PaymentID string  `json:"payment_id" binding:"required"`
		Amount    float64 `json:"amount" binding:"required"`
		Reason    string  `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	payment, err := h.service.RefundOrder(c.Request.Context(), tenantID.(uuid.UUID), uuid.Nil, paymentID.String(), req.Amount, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// TrackOrder provides tracking information for an order
// @Summary Track order
// @Description Get tracking information for an order by order number
// @Tags orders
// @Produce json
// @Param number path string true "Order Number"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/track/{number} [get]
func (h *Handler) TrackOrder(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	orderNumber := c.Param("number")
	tracking, err := h.service.TrackOrder(tenantID.(uuid.UUID), orderNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracking)
}

// GetOrderStats retrieves order statistics
// @Summary Get order statistics
// @Description Get various order statistics for analytics
// @Tags orders
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /orders/stats [get]
func (h *Handler) GetOrderStats(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	stats, err := h.service.GetOrderStats(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetCustomerOrders retrieves orders for the current customer
// @Summary Get customer orders
// @Description Get all orders for the authenticated customer
// @Tags orders
// @Produce json
// @Success 200 {object} []Order
// @Failure 400 {object} map[string]interface{}
// @Router /orders/my-orders [get]
func (h *Handler) GetCustomerOrders(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context required"})
		return
	}

	orders, err := h.service.GetCustomerOrders(tenantID.(uuid.UUID), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderInvoice generates an invoice for an order
// @Summary Get order invoice
// @Description Generate and return order invoice data
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/invoice [get]
func (h *Handler) GetOrderInvoice(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	orderID := c.Param("id")
	order, err := h.service.GetOrder(tenantID.(uuid.UUID), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Generate invoice data
	invoice := map[string]interface{}{
		"order":           order,
		"invoice_number":  "INV-" + order.OrderNumber,
		"invoice_date":    order.CreatedAt,
		"due_date":        order.CreatedAt.AddDate(0, 0, 30), // 30 days from order
		"company_info": map[string]string{
			"name":    "Your Company Name",
			"address": "Company Address",
			"phone":   "Company Phone",
			"email":   "company@example.com",
		},
		"items":      order.Items,
		"subtotal":   order.SubtotalAmount,
		"tax":        order.TaxAmount,
		"shipping":   order.ShippingAmount,
		"discount":   order.DiscountAmount,
		"total":      order.TotalAmount,
		"currency":   order.Currency,
	}

	c.JSON(http.StatusOK, invoice)
}

// RegisterRoutes registers all order routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		// Order CRUD
		orders.POST("", h.CreateOrder)
		orders.GET("", h.ListOrders)
		orders.GET("/stats", h.GetOrderStats)
		orders.GET("/my-orders", h.GetCustomerOrders)
		
		// Export/Import
		orders.GET("/export", h.ExportOrders)
		orders.POST("/import", h.ImportOrders)
		
		// Bulk operations
		orders.PATCH("/bulk", h.BulkUpdateOrders)
		
		orders.GET("/:id", h.GetOrder)
		orders.PATCH("/:id/status", h.UpdateOrderStatus)
		orders.POST("/:id/cancel", h.CancelOrder)
		orders.POST("/:id/payment", h.ProcessPayment)
		orders.POST("/:id/refund", h.RefundOrder)
		orders.GET("/:id/invoice", h.GetOrderInvoice)
		
		// Order tracking (by order number)
		orders.GET("/number/:number", h.GetOrderByNumber)
		orders.GET("/track/:number", h.TrackOrder)
	}
}

// ExportOrders exports orders to CSV or Excel format
// @Summary Export orders
// @Description Export orders to CSV or Excel format with optional filters
// @Tags orders
// @Produce application/octet-stream
// @Param format query string true "Export format (csv or excel)"
// @Param status query string false "Filter by status"
// @Param payment_status query string false "Filter by payment status"
// @Param from_date query string false "Filter from date (YYYY-MM-DD)"
// @Param to_date query string false "Filter to date (YYYY-MM-DD)"
// @Success 200 {file} file "Exported file"
// @Failure 400 {object} map[string]interface{}
// @Router /orders/export [get]
func (h *Handler) ExportOrders(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	// Get export format
	format := c.DefaultQuery("format", "csv")
	if format != "csv" && format != "excel" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format must be 'csv' or 'excel'"})
		return
	}

	// Build filters
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		filters["payment_status"] = paymentStatus
	}
	if fromDate := c.Query("from_date"); fromDate != "" {
		filters["from_date"] = fromDate
	}
	if toDate := c.Query("to_date"); toDate != "" {
		filters["to_date"] = toDate
	}

	// Export orders
	data, filename, err := h.service.ExportOrders(tenantID.(uuid.UUID), format, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set appropriate content type and headers
	var contentType string
	if format == "csv" {
		contentType = "text/csv"
	} else {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
	c.Data(http.StatusOK, contentType, data)
}

// ImportOrders imports orders from CSV
// @Summary Import orders
// @Description Import orders from CSV file
// @Tags orders
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /orders/import [post]
func (h *Handler) ImportOrders(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	totalRecords, successfulImports, failedImports, errors, err := h.service.ImportOrders(c.Request.Context(), tenantID.(uuid.UUID), src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := map[string]interface{}{
		"total_records":      totalRecords,
		"successful_imports": successfulImports,
		"failed_imports":     failedImports,
		"errors":             errors,
	}

	c.JSON(http.StatusOK, result)
}

// BulkUpdateOrders updates multiple orders
// @Summary Bulk update orders
// @Description Update multiple orders with the same action
// @Tags orders
// @Accept json
// @Produce json
// @Param bulk body map[string]interface{} true "Bulk update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /orders/bulk [patch]
func (h *Handler) BulkUpdateOrders(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	var req struct {
		OrderIDs []string               `json:"order_ids" binding:"required"`
		Action   string                 `json:"action" binding:"required"`
		Data     map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse order IDs
	orderIDs := make([]uuid.UUID, len(req.OrderIDs))
	for i, idStr := range req.OrderIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid order ID: %s", idStr)})
			return
		}
		orderIDs[i] = id
	}

	successfulUpdates, failedUpdates, errors, err := h.service.BulkUpdateOrders(c.Request.Context(), tenantID.(uuid.UUID), orderIDs, req.Action, req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := map[string]interface{}{
		"successful_updates": successfulUpdates,
		"failed_updates":     failedUpdates,
		"errors":             errors,
	}

	c.JSON(http.StatusOK, result)
}

// TODO: Add more handlers
// - GetOrderTimeline(c *gin.Context) - Get order status history
// - SendOrderNotification(c *gin.Context) - Send custom notifications
