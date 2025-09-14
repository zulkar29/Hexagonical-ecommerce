package order

import (
	"fmt"
	"log"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListOrders lists orders with filtering and pagination, supports analytics queries
// @Summary List orders with analytics support
// @Description Get paginated list of orders with optional filtering, stats, customer orders, or tracking
// @Tags orders
// @Produce json
// @Param type query string false "Query type: stats, my-orders, track, export"
// @Param number query string false "Order number for tracking"
// @Param format query string false "Export format: csv, excel"
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

	// Handle different query types
	queryType := c.Query("type")
	switch queryType {
	case "stats":
		stats, err := h.service.GetOrderStats(tenantID.(uuid.UUID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order stats"})
			return
		}
		c.JSON(http.StatusOK, stats)
		return

	case "my-orders":
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user context required"})
			return
		}
		orders, err := h.service.GetCustomerOrders(tenantID.(uuid.UUID), userID.(uuid.UUID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customer orders"})
			return
		}
		c.JSON(http.StatusOK, orders)
		return

	case "track":
		orderNumber := c.Query("number")
		if orderNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order number required for tracking"})
			return
		}
		tracking, err := h.service.TrackOrder(tenantID.(uuid.UUID), orderNumber)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order tracking not found"})
			return
		}
		c.JSON(http.StatusOK, tracking)
		return

	case "export":
		// Handle export functionality
		format := c.DefaultQuery("format", "csv")
		if format != "csv" && format != "excel" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "format must be 'csv' or 'excel'"})
			return
		}

		// Build filters for export
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export orders"})
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
		return
	}

	// Default: Regular order listing with pagination and filtering
	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Printf("Invalid page parameter, using default: %v", err)
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		log.Printf("Invalid limit parameter, using default: %v", err)
		limit = 20
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders"})
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

// UpdateOrder updates an existing order or performs status operations
// @Summary Update order or perform operations
// @Description Update order details or perform status operations (cancel, process-payment, refund, status)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param action query string false "Action to perform: cancel, process-payment, refund, status"
// @Param order body map[string]interface{} false "Order update data (varies by action)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [put]
func (h *Handler) UpdateOrder(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	// Handle different actions
	action := c.Query("action")
	switch action {
	case "cancel":
		var req struct {
			Reason string `json:"reason"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Failed to bind cancel request JSON: %v", err)
		}
		order, cancelErr := h.service.CancelOrder(tenantID.(uuid.UUID), orderID.String(), req.Reason)
		if cancelErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to cancel order"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"order": order})
		return

	case "process-payment":
		var req struct {
			PaymentID       string `json:"payment_id" binding:"required"`
			PaymentMethodID string `json:"payment_method_id" binding:"required"`
			Confirmation    string `json:"confirmation"`
		}
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment request data"})
			return
		}
		paymentID, parseErr := uuid.Parse(req.PaymentID)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
			return
		}
		_, parseErr = uuid.Parse(req.PaymentMethodID)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment method ID"})
			return
		}
		order, processErr := h.service.ProcessPayment(tenantID.(uuid.UUID), paymentID.String())
		if processErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process payment"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"order": order})
		return

	case "refund":
		var req struct {
			PaymentID string  `json:"payment_id" binding:"required"`
			Amount    float64 `json:"amount" binding:"required"`
			Reason    string  `json:"reason"`
		}
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refund request data"})
			return
		}
		paymentID, parseErr := uuid.Parse(req.PaymentID)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
			return
		}
		payment, refundErr := h.service.RefundOrder(c.Request.Context(), tenantID.(uuid.UUID), uuid.Nil, paymentID.String(), req.Amount, req.Reason)
		if refundErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process refund"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"payment": payment})
		return

	case "status":
		var req struct {
			Status         string `json:"status" binding:"required"`
			TrackingNumber string `json:"tracking_number"`
			TrackingURL    string `json:"tracking_url"`
			Notes          string `json:"notes"`
		}
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status update request data"})
			return
		}
		order, statusErr := h.service.UpdateOrderStatus(c.Request.Context(), tenantID.(uuid.UUID), orderID, OrderStatus(req.Status), req.TrackingNumber, req.TrackingURL, req.Notes)
		if statusErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update order status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"order": order})
		return
	}

	// Default: Regular order update (if no action specified)
	var req UpdateOrderRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order update data"})
		return
	}

	order, updateErr := h.service.UpdateOrder(c.Request.Context(), tenantID.(uuid.UUID), orderID.String(), &req)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// HandleOrderOperations handles bulk operations and import functionality
// @Summary Handle order operations
// @Description Handle bulk operations (update, delete) and import functionality
// @Tags orders
// @Accept json
// @Produce json
// @Param operation query string true "Operation type: bulk-update, bulk-delete, import"
// @Param format query string false "Import format: csv, excel (for import operation)"
// @Param body body map[string]interface{} true "Operation data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /orders/operations [post]
func (h *Handler) HandleOrderOperations(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant context required"})
		return
	}

	operation := c.Query("operation")
	switch operation {
	case "bulk-update":
		// Parse bulk update request
		var req struct {
			OrderIDs []string               `json:"order_ids" validate:"required,min=1"`
			Action   string                 `json:"action" validate:"required,oneof=update_status cancel refund"`
			Data     map[string]interface{} `json:"data,omitempty"`
		}

		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bulk update request data"})
			return
		}

		// Convert string IDs to UUIDs
		orderIDs := make([]uuid.UUID, len(req.OrderIDs))
		for i, idStr := range req.OrderIDs {
			id, parseErr := uuid.Parse(idStr)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid order ID: %s", idStr)})
				return
			}
			orderIDs[i] = id
		}

		// Perform bulk update
		successful, failed, errors, bulkErr := h.service.BulkUpdateOrders(c.Request.Context(), tenantID.(uuid.UUID), orderIDs, req.Action, req.Data)
		if bulkErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform bulk update"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"successful": successful,
			"failed":     failed,
			"errors":     errors,
		})
		return

	case "bulk-delete":
		var req struct {
			OrderIDs []string `json:"order_ids" binding:"required"`
			Reason   string   `json:"reason,omitempty"`
		}
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bulk delete request data"})
			return
		}

		// Perform bulk delete
		successful, failed, errors, deleteErr := h.service.BulkDeleteOrders(c.Request.Context(), tenantID.(uuid.UUID), req.OrderIDs, req.Reason)
		if deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform bulk delete"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"successful": successful,
			"failed":     failed,
			"errors":     errors,
		})
		return

	case "import":
		// Handle file import
		file, fileErr := c.FormFile("file")
		if fileErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}

		format := c.DefaultPostForm("format", "csv")
		if format != "csv" && format != "excel" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Supported: csv, excel"})
			return
		}

		// Open uploaded file
		src, openErr := file.Open()
		if openErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer src.Close()

		// Import orders
		totalRecords, successfulImports, failedImports, errors, importErr := h.service.ImportOrders(c.Request.Context(), tenantID.(uuid.UUID), src)
		if importErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": importErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"total_records":      totalRecords,
			"successful_imports": successfulImports,
			"failed_imports":     failedImports,
			"errors":             errors,
		})
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation type"})
	}
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
		"order":          order,
		"invoice_number": "INV-" + order.OrderNumber,
		"invoice_date":   order.CreatedAt,
		"due_date":       order.CreatedAt.AddDate(0, 0, 30), // 30 days from order
		"company_info": map[string]string{
			"name":    "Your Company Name",
			"address": "Company Address",
			"phone":   "Company Phone",
			"email":   "company@example.com",
		},
		"items":    order.Items,
		"subtotal": order.SubtotalAmount,
		"tax":      order.TaxAmount,
		"shipping": order.ShippingAmount,
		"discount": order.DiscountAmount,
		"total":    order.TotalAmount,
		"currency": order.Currency,
	}

	c.JSON(http.StatusOK, invoice)
}

// DeleteOrder deletes an order
// @Summary Delete order
// @Description Delete an order (soft delete)
// @Tags orders
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [delete]
func (h *Handler) DeleteOrder(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	err = h.service.DeleteOrder(c.Request.Context(), tenantID.(uuid.UUID), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// TrackOrder tracks an order by ID
// @Summary Track order
// @Description Get order tracking information
// @Tags orders
// @Param id path string true "Order ID"
// @Param public query bool false "Public access (no auth required)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/tracking [get]
func (h *Handler) TrackOrder(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	orderID := c.Param("id")
	tracking, err := h.service.TrackOrder(tenantID.(uuid.UUID), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracking)
}

// RegisterRoutes registers all order routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		// Order CRUD with analytics support via query parameters
		// GET /orders?type=stats|my-orders|track&number=xxx
		orders.POST("", h.CreateOrder)
		orders.GET("", h.ListOrders) // Supports type=stats,my-orders,track via query params

		// Data operations with operation type parameter
		// POST /orders/operations?type=import|bulk&action=xxx
		orders.POST("/operations", h.HandleOrderOperations) // Handles import, bulk operations

		// Individual order operations
		orders.GET("/:id", h.GetOrder)       // Supports include=invoice,timeline via query params
		orders.PATCH("/:id", h.UpdateOrder)  // Changed from PUT to PATCH to match API spec
		orders.DELETE("/:id", h.DeleteOrder) // Added missing DELETE endpoint

		// Order lookup and tracking
		orders.GET("/lookup/:number", h.GetOrderByNumber) // Changed from /number/ to /lookup/ to match API spec
		orders.GET("/:id/tracking", h.TrackOrder)         // Added missing tracking endpoint
	}
}

// TODO: Add more handlers
// - GetOrderTimeline(c *gin.Context) - Get order status history
// - SendOrderNotification(c *gin.Context) - Send custom notifications
