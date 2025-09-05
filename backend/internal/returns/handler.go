package returns

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for returns
type Handler struct {
	service Service
}

// NewHandler creates a new returns handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers all return routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	returns := router.Group("/returns")
	{
		// Return management
		returns.POST("", h.CreateReturn)
		returns.GET("", h.ListReturns) // Supports ?type=stats&category=customer|order for analytics
		returns.GET("/:id", h.GetReturn)
		returns.PUT("/:id", h.UpdateReturn) // Supports ?action=approve|reject|process|complete for workflow
		returns.DELETE("/:id", h.DeleteReturn)
		
		// Return items (consolidated)
		returns.POST("/:id/items", h.AddReturnItem)
		returns.PUT("/:id/items/:itemId", h.UpdateReturnItem)
		returns.DELETE("/:id/items/:itemId", h.RemoveReturnItem)
		
		// Return operations (consolidated)
		returns.POST("/:id/operations", h.HandleReturnOperation) // Supports operation=label|track for shipping
		
		// Return lookup by number
		returns.GET("/number/:number", h.GetReturnByNumber)
	}
	
	// Return reasons
	reasons := router.Group("/return-reasons")
	{
		reasons.POST("", h.CreateReturnReason)
		reasons.GET("", h.ListReturnReasons)
		reasons.GET("/:id", h.GetReturnReason)
		reasons.PUT("/:id", h.UpdateReturnReason)
		reasons.DELETE("/:id", h.DeleteReturnReason)
	}
}

// CreateReturn creates a new return
// @Summary Create return
// @Description Create a new return request
// @Tags returns
// @Accept json
// @Produce json
// @Param return body Return true "Return creation request"
// @Success 201 {object} Return
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns [post]
func (h *Handler) CreateReturn(c *gin.Context) {
	var returnReq Return
	if err := c.ShouldBindJSON(&returnReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	return_, err := h.service.CreateReturn(c.Request.Context(), &returnReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create return", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, return_)
}

// ListReturns retrieves returns with filtering and pagination
// @Summary List returns
// @Description Get a list of returns with optional filtering and pagination
// @Tags returns
// @Accept json
// @Produce json
// @Param status query []string false "Filter by status" collectionFormat(multi)
// @Param type query []string false "Filter by type" collectionFormat(multi)
// @Param customer_id query string false "Filter by customer ID"
// @Param order_id query string false "Filter by order ID"
// @Param reason_id query string false "Filter by reason ID"
// @Param created_after query string false "Filter by created after date (RFC3339)"
// @Param created_before query string false "Filter by created before date (RFC3339)"
// @Param min_refund_amount query number false "Filter by minimum refund amount"
// @Param max_refund_amount query number false "Filter by maximum refund amount"
// @Param search query string false "Search in return number, reason text, or description"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_order query string false "Sort order (asc/desc)" default(desc)
// @Success 200 {object} PaginatedReturnsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns [get]
func (h *Handler) ListReturns(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)
	
	// Check for analytics type
	analyticsType := c.Query("type")
	switch analyticsType {
	case "stats":
		category := c.Query("category")
		switch category {
		case "customer":
			customerIDStr := c.Query("customer_id")
			if customerIDStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id is required"})
				return
			}
			customerID, err := uuid.Parse(customerIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
				return
			}
			returns, total, err := h.service.GetReturnsByCustomer(c.Request.Context(), tenantID, customerID, ReturnFilter{Limit: getIntQueryParam(c, "limit", 10)})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get returns by customer", "details": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": returns, "total": total})
			return
		case "order":
			orderIDStr := c.Query("order_id")
			if orderIDStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
				return
			}
			orderID, err := uuid.Parse(orderIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
				return
			}
			returns, err := h.service.GetReturnsByOrder(c.Request.Context(), tenantID, orderID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get returns by order", "details": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": returns})
			return
		default:
			filter := StatsFilter{
				GroupBy: c.DefaultQuery("group_by", "day"),
			}
			if dateFrom := c.Query("date_from"); dateFrom != "" {
				if date, err := time.Parse(time.RFC3339, dateFrom); err == nil {
					filter.DateFrom = &date
				}
			}
			if dateTo := c.Query("date_to"); dateTo != "" {
				if date, err := time.Parse(time.RFC3339, dateTo); err == nil {
					filter.DateTo = &date
				}
			}
			stats, err := h.service.GetReturnStats(c.Request.Context(), tenantID, filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get return stats", "details": err.Error()})
				return
			}
			c.JSON(http.StatusOK, stats)
			return
		}
	}
	
	// Parse query parameters
	filter := ReturnFilter{
		Page:      getIntQueryParam(c, "page", 1),
		Limit:     getIntQueryParam(c, "limit", 20),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
		Search:    c.Query("search"),
	}
	
	// Parse status filter
	if statusParams := c.QueryArray("status"); len(statusParams) > 0 {
		for _, status := range statusParams {
			filter.Status = append(filter.Status, ReturnStatus(status))
		}
	}
	
	// Parse type filter
	if typeParams := c.QueryArray("type"); len(typeParams) > 0 {
		for _, returnType := range typeParams {
			filter.Type = append(filter.Type, ReturnType(returnType))
		}
	}
	
	// Parse UUID filters
	if customerID := c.Query("customer_id"); customerID != "" {
		if id, err := uuid.Parse(customerID); err == nil {
			filter.CustomerID = &id
		}
	}
	
	if orderID := c.Query("order_id"); orderID != "" {
		if id, err := uuid.Parse(orderID); err == nil {
			filter.OrderID = &id
		}
	}
	
	if reasonID := c.Query("reason_id"); reasonID != "" {
		if id, err := uuid.Parse(reasonID); err == nil {
			filter.ReasonID = &id
		}
	}
	
	// Parse date filters
	if createdAfter := c.Query("created_after"); createdAfter != "" {
		if date, err := time.Parse(time.RFC3339, createdAfter); err == nil {
			filter.CreatedAfter = &date
		}
	}
	
	if createdBefore := c.Query("created_before"); createdBefore != "" {
		if date, err := time.Parse(time.RFC3339, createdBefore); err == nil {
			filter.CreatedBefore = &date
		}
	}
	
	// Parse amount filters
	if minAmount := c.Query("min_refund_amount"); minAmount != "" {
		if amount, err := strconv.ParseFloat(minAmount, 64); err == nil {
			filter.MinRefundAmount = &amount
		}
	}
	
	if maxAmount := c.Query("max_refund_amount"); maxAmount != "" {
		if amount, err := strconv.ParseFloat(maxAmount, 64); err == nil {
			filter.MaxRefundAmount = &amount
		}
	}
	
	returns, total, err := h.service.ListReturns(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list returns", "details": err.Error()})
		return
	}
	
	response := PaginatedReturnsResponse{
		Data:       returns,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: (total + int64(filter.Limit) - 1) / int64(filter.Limit),
	}
	
	c.JSON(http.StatusOK, response)
}

// GetReturn retrieves a return by ID
// @Summary Get return by ID
// @Description Get a specific return by its ID
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Success 200 {object} Return
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id} [get]
func (h *Handler) GetReturn(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	return_, err := h.service.GetReturn(c.Request.Context(), tenantID, returnID)
	if err != nil {
		if err.Error() == "return not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Return not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get return", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, return_)
}

// GetReturnByNumber retrieves a return by return number
// @Summary Get return by number
// @Description Get a specific return by its return number
// @Tags returns
// @Accept json
// @Produce json
// @Param number path string true "Return Number"
// @Success 200 {object} Return
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/number/{number} [get]
func (h *Handler) GetReturnByNumber(c *gin.Context) {
	returnNumber := c.Param("number")
	if returnNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Return number is required"})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	return_, err := h.service.GetReturnByNumber(c.Request.Context(), tenantID, returnNumber)
	if err != nil {
		if err.Error() == "return not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Return not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get return", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, return_)
}

// UpdateReturn updates a return
// @Summary Update return
// @Description Update a return's details
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Param return body Return true "Return update request"
// @Success 200 {object} Return
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id} [patch]
func (h *Handler) UpdateReturn(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	// Check for workflow action
	action := c.Query("action")
	switch action {
	case "approve":
		var req ApprovalRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}
		tenantID := getTenantIDFromContext(c)
		return_, err := h.service.ApproveReturn(c.Request.Context(), tenantID, returnID, req.ApprovedBy, req.ApprovalNote)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve return", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, return_)
		return
	case "reject":
		var req RejectionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}
		tenantID := getTenantIDFromContext(c)
		return_, err := h.service.RejectReturn(c.Request.Context(), tenantID, returnID, req.RejectedBy, req.RejectionReason)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject return", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, return_)
		return
	case "process":
		var req ProcessingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}
		tenantID := getTenantIDFromContext(c)
		return_, err := h.service.ProcessReturn(c.Request.Context(), tenantID, returnID, req.ProcessedBy, req.ProcessingNote, req.TrackingNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process return", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, return_)
		return
	case "complete":
		var req CompletionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}
		tenantID := getTenantIDFromContext(c)
		var refundAmount float64
		if req.RefundAmount != nil {
			refundAmount = *req.RefundAmount
		}
		return_, err := h.service.CompleteReturn(c.Request.Context(), tenantID, returnID, req.CompletedBy, req.CompletionNote, refundAmount, req.RefundMethod)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete return", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, return_)
		return
	}
	
	// Regular update
	var updatedReturn Return
	if err := c.ShouldBindJSON(&updatedReturn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	updatedReturn.ID = returnID
	updatedReturn.TenantID = tenantID
	return_, err := h.service.UpdateReturn(c.Request.Context(), &updatedReturn)
	if err != nil {
		if err.Error() == "return not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Return not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update return", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, return_)
}

// DeleteReturn deletes a return
// @Summary Delete return
// @Description Delete a return (only pending returns can be deleted)
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id} [delete]
func (h *Handler) DeleteReturn(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	err = h.service.DeleteReturn(c.Request.Context(), tenantID, returnID)
	if err != nil {
		if err.Error() == "return not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Return not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete return", "details": err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ApprovalRequest represents the request body for approving a return
type ApprovalRequest struct {
	ApprovedBy   uuid.UUID `json:"approved_by" binding:"required"`
	ApprovalNote string    `json:"approval_note,omitempty"`
}

// ApproveReturn approves a return
// @Summary Approve return
// @Description Approve a pending return
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Param approval body ApprovalRequest true "Approval request"
// @Success 200 {object} Return
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id}/approve [post]
func (h *Handler) ApproveReturn(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	var req ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	return_, err := h.service.ApproveReturn(c.Request.Context(), tenantID, returnID, req.ApprovedBy, req.ApprovalNote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve return", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, return_)
}

// RejectionRequest represents the request body for rejecting a return
type RejectionRequest struct {
	RejectedBy      uuid.UUID `json:"rejected_by" binding:"required"`
	RejectionReason string    `json:"rejection_reason" binding:"required"`
}

// RejectReturn rejects a return
// @Summary Reject return
// @Description Reject a pending return
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Param rejection body RejectionRequest true "Rejection request"
// @Success 200 {object} Return
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id}/reject [post]
func (h *Handler) RejectReturn(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	var req RejectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	return_, err := h.service.RejectReturn(c.Request.Context(), tenantID, returnID, req.RejectedBy, req.RejectionReason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject return", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, return_)
}

// ProcessingRequest represents the request body for processing a return
type ProcessingRequest struct {
	ProcessedBy    uuid.UUID `json:"processed_by" binding:"required"`
	ProcessingNote string    `json:"processing_note,omitempty"`
	TrackingNumber string    `json:"tracking_number,omitempty"`
}

// ProcessReturn processes a return
// @Summary Process return
// @Description Process an approved return
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Param processing body ProcessingRequest true "Processing request"
// @Success 200 {object} Return
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id}/process [post]
func (h *Handler) ProcessReturn(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	var req ProcessingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	return_, err := h.service.ProcessReturn(c.Request.Context(), tenantID, returnID, req.ProcessedBy, req.ProcessingNote, req.TrackingNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process return", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, return_)
}

// CompletionRequest represents the request body for completing a return
type CompletionRequest struct {
	CompletedBy     uuid.UUID `json:"completed_by" binding:"required"`
	CompletionNote  string    `json:"completion_note,omitempty"`
	RefundAmount    *float64  `json:"refund_amount,omitempty"`
	RefundMethod    string    `json:"refund_method,omitempty"`
}

// CompleteReturn completes a return
// @Summary Complete return
// @Description Complete a processed return
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Param completion body CompletionRequest true "Completion request"
// @Success 200 {object} Return
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id}/complete [post]
func (h *Handler) CompleteReturn(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	var req CompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	var refundAmount float64
	if req.RefundAmount != nil {
		refundAmount = *req.RefundAmount
	}
	return_, err := h.service.CompleteReturn(c.Request.Context(), tenantID, returnID, req.CompletedBy, req.CompletionNote, refundAmount, req.RefundMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete return", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, return_)
}

// AddReturnItem adds an item to a return
// @Summary Add return item
// @Description Add an item to an existing return
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Param item body ReturnItem true "Return item request"
// @Success 201 {object} ReturnItem
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id}/items [post]
func (h *Handler) AddReturnItem(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	var item ReturnItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	// Set the return ID from the URL parameter
	item.ReturnID = returnID
	
	createdItem, err := h.service.AddReturnItem(c.Request.Context(), &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add return item", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, createdItem)
}

// UpdateReturnItem updates a return item
// @Summary Update return item
// @Description Update an existing return item
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Param itemId path string true "Return Item ID"
// @Param item body ReturnItem true "Return item update request"
// @Success 200 {object} ReturnItem
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id}/items/{itemId} [put]
func (h *Handler) UpdateReturnItem(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}
	
	var item ReturnItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	// Set the IDs from the URL parameters
	item.ID = itemID
	item.ReturnID = returnID
	
	updatedItem, err := h.service.UpdateReturnItem(c.Request.Context(), &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update return item", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, updatedItem)
}

// RemoveReturnItem removes an item from a return
// @Summary Remove return item
// @Description Remove an item from a return
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Param item_id path string true "Return Item ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id}/items/{item_id} [delete]
func (h *Handler) RemoveReturnItem(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	err = h.service.RemoveReturnItem(c.Request.Context(), tenantID, returnID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove return item", "details": err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}

// GenerateReturnLabel generates a shipping label for return
// @Summary Generate return shipping label
// @Description Generate a shipping label for an approved return
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Success 200 {object} ShippingLabel
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id}/shipping-label [post]
func (h *Handler) GenerateReturnLabel(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	label, err := h.service.GenerateReturnLabel(c.Request.Context(), tenantID, returnID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate return label", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, label)
}

// TrackReturnShipment tracks a return shipment
// @Summary Track return shipment
// @Description Get tracking information for a return shipment
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Success 200 {object} ShipmentTracking
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/{id}/tracking [get]
func (h *Handler) TrackReturnShipment(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	// Get tracking number from query parameter or return record
	trackingNumber := c.Query("tracking_number")
	if trackingNumber == "" {
		// If no tracking number provided, get it from the return record
		returnRecord, err := h.service.GetReturn(c.Request.Context(), tenantID, returnID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Return not found"})
			return
		}
		trackingNumber = returnRecord.TrackingNumber
		if trackingNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No tracking number available for this return"})
			return
		}
	}
	
	tracking, err := h.service.TrackReturnShipment(c.Request.Context(), tenantID, returnID, trackingNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track return shipment", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, tracking)
}

// GetReturnStats retrieves return statistics
// @Summary Get return statistics
// @Description Get return statistics and analytics
// @Tags returns
// @Accept json
// @Produce json
// @Param date_from query string false "Start date for statistics (RFC3339)"
// @Param date_to query string false "End date for statistics (RFC3339)"
// @Param group_by query string false "Group statistics by (day, week, month, year)"
// @Success 200 {object} ReturnStats
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/stats [get]
func (h *Handler) GetReturnStats(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)
	
	filter := StatsFilter{
		GroupBy: c.DefaultQuery("group_by", "day"),
	}
	
	// Parse date filters
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if date, err := time.Parse(time.RFC3339, dateFrom); err == nil {
			filter.DateFrom = &date
		}
	}
	
	if dateTo := c.Query("date_to"); dateTo != "" {
		if date, err := time.Parse(time.RFC3339, dateTo); err == nil {
			filter.DateTo = &date
		}
	}
	
	stats, err := h.service.GetReturnStats(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get return stats", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// GetReturnsByCustomer retrieves returns for a specific customer
// @Summary Get returns by customer
// @Description Get returns for a specific customer
// @Tags returns
// @Accept json
// @Produce json
// @Param customer_id path string true "Customer ID"
// @Param limit query int false "Maximum number of returns to return" default(10)
// @Success 200 {array} Return
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/customer/{customer_id} [get]
func (h *Handler) GetReturnsByCustomer(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	
	// Create filter from query parameters
	filter := ReturnFilter{
		Limit: getIntQueryParam(c, "limit", 10),
	}
	
	returns, total, err := h.service.GetReturnsByCustomer(c.Request.Context(), tenantID, customerID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get returns by customer", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"returns": returns,
		"total":   total,
	})
}

// GetReturnsByOrder retrieves returns for a specific order
// @Summary Get returns by order
// @Description Get returns for a specific order
// @Tags returns
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {array} Return
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /returns/order/{order_id} [get]
func (h *Handler) GetReturnsByOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	returns, err := h.service.GetReturnsByOrder(c.Request.Context(), tenantID, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get returns by order", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, returns)
}

// CreateReturnReason creates a new return reason
// @Summary Create return reason
// @Description Create a new return reason
// @Tags return-reasons
// @Accept json
// @Produce json
// @Param reason body ReturnReason true "Return reason request"
// @Success 201 {object} ReturnReason
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /return-reasons [post]
func (h *Handler) CreateReturnReason(c *gin.Context) {
	var reason ReturnReason
	if err := c.ShouldBindJSON(&reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	reason.TenantID = tenantID
	
	createdReason, err := h.service.CreateReturnReason(c.Request.Context(), &reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create return reason", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, createdReason)
}

// ListReturnReasons retrieves return reasons
// @Summary List return reasons
// @Description Get a list of return reasons
// @Tags return-reasons
// @Accept json
// @Produce json
// @Param active_only query bool false "Filter to active reasons only" default(false)
// @Success 200 {array} ReturnReason
// @Failure 500 {object} ErrorResponse
// @Router /return-reasons [get]
func (h *Handler) ListReturnReasons(c *gin.Context) {
	activeOnly := c.DefaultQuery("active_only", "false") == "true"
	tenantID := getTenantIDFromContext(c)
	
	reasons, err := h.service.ListReturnReasons(c.Request.Context(), tenantID, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list return reasons", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, reasons)
}

// GetReturnReason retrieves a return reason by ID
// @Summary Get return reason by ID
// @Description Get a specific return reason by its ID
// @Tags return-reasons
// @Accept json
// @Produce json
// @Param id path string true "Return Reason ID"
// @Success 200 {object} ReturnReason
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /return-reasons/{id} [get]
func (h *Handler) GetReturnReason(c *gin.Context) {
	reasonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reason ID"})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	reason, err := h.service.GetReturnReason(c.Request.Context(), tenantID, reasonID)
	if err != nil {
		if err.Error() == "return reason not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Return reason not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get return reason", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, reason)
}

// UpdateReturnReason updates a return reason
// @Summary Update return reason
// @Description Update a return reason's details
// @Tags return-reasons
// @Accept json
// @Produce json
// @Param id path string true "Return Reason ID"
// @Param reason body ReturnReason true "Return reason update request"
// @Success 200 {object} ReturnReason
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /return-reasons/{id} [patch]
func (h *Handler) UpdateReturnReason(c *gin.Context) {
	reasonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reason ID"})
		return
	}
	
	var reason ReturnReason
	if err := c.ShouldBindJSON(&reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	// Set the ID from the URL parameter
	reason.ID = reasonID
	
	tenantID := getTenantIDFromContext(c)
	reason.TenantID = tenantID
	updatedReason, err := h.service.UpdateReturnReason(c.Request.Context(), &reason)
	if err != nil {
		if err.Error() == "return reason not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Return reason not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update return reason", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, updatedReason)
}

// DeleteReturnReason deletes a return reason
// @Summary Delete return reason
// @Description Delete a return reason
// @Tags return-reasons
// @Accept json
// @Produce json
// @Param id path string true "Return Reason ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /return-reasons/{id} [delete]
func (h *Handler) DeleteReturnReason(c *gin.Context) {
	reasonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reason ID"})
		return
	}
	
	tenantID := getTenantIDFromContext(c)
	err = h.service.DeleteReturnReason(c.Request.Context(), tenantID, reasonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete return reason", "details": err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}

// Response DTOs
type PaginatedReturnsResponse struct {
	Data       []*Return `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int64     `json:"total_pages"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// HandleReturnOperation handles POST /api/returns/:id/operations for shipping operations
func (h *Handler) HandleReturnOperation(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	returnIDStr := c.Param("id")
	returnID, err := uuid.Parse(returnIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}

	operationType := c.Query("type")
	switch operationType {
	case "generate_label":
		var req struct {
			ShippingAddress string `json:"shipping_address"`
			Carrier         string `json:"carrier"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		label, err := h.service.GenerateReturnLabel(c.Request.Context(), tenantID.(uuid.UUID), returnID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": label})
		return
	case "track_shipment":
		var req struct {
			TrackingNumber string `json:"tracking_number"`
			Carrier        string `json:"carrier"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		tracking, err := h.service.TrackReturnShipment(c.Request.Context(), tenantID.(uuid.UUID), returnID, req.TrackingNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": tracking})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type"})
	}
}

// Helper functions
func getTenantIDFromContext(c *gin.Context) uuid.UUID {
	// TODO: Extract tenant ID from context/middleware
	// This should be set by authentication middleware
	return uuid.New() // Placeholder
}

func getIntQueryParam(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}