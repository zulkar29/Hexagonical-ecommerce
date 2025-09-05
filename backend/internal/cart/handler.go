package cart

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for cart operations
type Handler struct {
	service *Service
}

// NewHandler creates a new cart handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers cart routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	cartGroup := router.Group("/carts")
	{
		// Cart management
		cartGroup.POST("", h.CreateCart)
		cartGroup.GET("/:cart_id", h.GetCart)
		cartGroup.GET("/customer/:customer_id", h.GetCartByCustomer)
		cartGroup.GET("/session/:session_id", h.GetCartBySession)
		cartGroup.GET("/:cart_id/summary", h.GetCartSummary)
		cartGroup.GET("", h.ListCarts)
		
		// Cart items
		cartGroup.POST("/:cart_id/items", h.AddItem)
		cartGroup.PUT("/:cart_id/items/:item_id", h.UpdateItem)
		cartGroup.DELETE("/:cart_id/items/:item_id", h.RemoveItem)
		cartGroup.DELETE("/:cart_id/items", h.ClearCart)
		
		// Cart operations
		cartGroup.POST("/:cart_id/coupon", h.ApplyCoupon)
		cartGroup.DELETE("/:cart_id/coupon", h.RemoveCoupon)
		cartGroup.PUT("/:cart_id/address", h.UpdateAddress)
		cartGroup.PUT("/:cart_id/shipping", h.UpdateShipping)
		cartGroup.POST("/merge", h.MergeGuestCart)
		cartGroup.POST("/:cart_id/abandon", h.AbandonCart)
		cartGroup.POST("/:cart_id/convert", h.ConvertCart)
	}
}

// CreateCart creates a new cart
func (h *Handler) CreateCart(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	var req CreateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.CreateCart(tenantID.(uuid.UUID), req)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": cart})
}

// GetCart retrieves a cart by ID
func (h *Handler) GetCart(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	cart, err := h.service.GetCart(tenantID.(uuid.UUID), cartID)
	if err != nil {
		if err == ErrCartNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
			return
		}
		if err == ErrCartExpired {
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// GetCartByCustomer retrieves cart for a customer
func (h *Handler) GetCartByCustomer(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	customerIDStr := c.Param("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	cart, err := h.service.GetCartByCustomer(tenantID.(uuid.UUID), customerID)
	if err != nil {
		if err == ErrCartNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
			return
		}
		if err == ErrCartExpired {
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// GetCartBySession retrieves cart for a guest session
func (h *Handler) GetCartBySession(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	cart, err := h.service.GetCartBySession(tenantID.(uuid.UUID), sessionID)
	if err != nil {
		if err == ErrCartNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
			return
		}
		if err == ErrCartExpired {
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// GetCartSummary retrieves cart summary
func (h *Handler) GetCartSummary(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	summary, err := h.service.GetCartSummary(tenantID.(uuid.UUID), cartID)
	if err != nil {
		if err == ErrCartNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart summary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}

// ListCarts lists carts with filtering and pagination, or returns stats when type=stats
func (h *Handler) ListCarts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	// Check if stats are requested
	if c.Query("type") == "stats" {
		stats, err := h.service.GetCartStats(tenantID.(uuid.UUID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart statistics"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": stats})
		return
	}

	// Parse query parameters for regular listing
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	// Build filter
	filter := CartListFilter{
		Status: CartStatus(c.Query("status")),
	}

	// Parse min_total
	if minTotalStr := c.Query("min_total"); minTotalStr != "" {
		if minTotal, err := strconv.ParseFloat(minTotalStr, 64); err == nil {
			filter.MinTotal = &minTotal
		}
	}

	// Parse max_total
	if maxTotalStr := c.Query("max_total"); maxTotalStr != "" {
		if maxTotal, err := strconv.ParseFloat(maxTotalStr, 64); err == nil {
			filter.MaxTotal = &maxTotal
		}
	}

	if customerIDStr := c.Query("customer_id"); customerIDStr != "" {
		if customerID, err := uuid.Parse(customerIDStr); err == nil {
			filter.CustomerID = &customerID
		}
	}

	carts, total, err := h.service.ListCarts(tenantID.(uuid.UUID), filter, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list carts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": carts,
		"pagination": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
		},
	})
}

// AddItem adds an item to the cart
func (h *Handler) AddItem(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.AddItem(tenantID.(uuid.UUID), cartID, req)
	if err != nil {
		switch err {
		case ErrCartNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		case ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		case ErrInsufficientStock:
			c.JSON(http.StatusConflict, gin.H{"error": "Insufficient stock"})
		case ErrCartNotModifiable:
			c.JSON(http.StatusConflict, gin.H{"error": "Cart cannot be modified"})
		default:
			if strings.Contains(err.Error(), "validation") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// UpdateItem updates a cart item
func (h *Handler) UpdateItem(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	itemIDStr := c.Param("item_id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.UpdateItem(tenantID.(uuid.UUID), cartID, itemID, req)
	if err != nil {
		switch err {
		case ErrCartNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		case ErrItemNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		case ErrInsufficientStock:
			c.JSON(http.StatusConflict, gin.H{"error": "Insufficient stock"})
		case ErrCartNotModifiable:
			c.JSON(http.StatusConflict, gin.H{"error": "Cart cannot be modified"})
		default:
			if strings.Contains(err.Error(), "validation") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// RemoveItem removes an item from the cart
func (h *Handler) RemoveItem(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	itemIDStr := c.Param("item_id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	cart, err := h.service.RemoveItem(tenantID.(uuid.UUID), cartID, itemID)
	if err != nil {
		switch err {
		case ErrCartNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		case ErrItemNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		case ErrCartNotModifiable:
			c.JSON(http.StatusConflict, gin.H{"error": "Cart cannot be modified"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove cart item"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// ClearCart removes all items from the cart
func (h *Handler) ClearCart(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	cart, err := h.service.ClearCart(tenantID.(uuid.UUID), cartID)
	if err != nil {
		switch err {
		case ErrCartNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		case ErrCartNotModifiable:
			c.JSON(http.StatusConflict, gin.H{"error": "Cart cannot be modified"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// ApplyCoupon applies a coupon to the cart
func (h *Handler) ApplyCoupon(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	var req ApplyCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.ApplyCoupon(tenantID.(uuid.UUID), cartID, req)
	if err != nil {
		switch err {
		case ErrCartNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		case ErrInvalidCoupon:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired coupon"})
		case ErrCartNotModifiable:
			c.JSON(http.StatusConflict, gin.H{"error": "Cart cannot be modified"})
		default:
			if strings.Contains(err.Error(), "validation") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply coupon"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// RemoveCoupon removes coupon from the cart
func (h *Handler) RemoveCoupon(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	cart, err := h.service.RemoveCoupon(tenantID.(uuid.UUID), cartID)
	if err != nil {
		switch err {
		case ErrCartNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		case ErrCartNotModifiable:
			c.JSON(http.StatusConflict, gin.H{"error": "Cart cannot be modified"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove coupon"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// UpdateAddress updates shipping/billing address
func (h *Handler) UpdateAddress(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	var req UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.UpdateAddress(tenantID.(uuid.UUID), cartID, req)
	if err != nil {
		switch err {
		case ErrCartNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		case ErrCartNotModifiable:
			c.JSON(http.StatusConflict, gin.H{"error": "Cart cannot be modified"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// UpdateShipping updates shipping method
func (h *Handler) UpdateShipping(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	var req UpdateShippingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.UpdateShipping(tenantID.(uuid.UUID), cartID, req)
	if err != nil {
		switch err {
		case ErrCartNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		case ErrCartNotModifiable:
			c.JSON(http.StatusConflict, gin.H{"error": "Cart cannot be modified"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shipping"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// MergeGuestCart merges guest cart to customer cart
func (h *Handler) MergeGuestCart(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	var req struct {
		SessionID  string    `json:"session_id" binding:"required"`
		CustomerID uuid.UUID `json:"customer_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.MergeGuestCart(tenantID.(uuid.UUID), req.SessionID, req.CustomerID)
	if err != nil {
		if err == ErrCartNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Guest cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to merge cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// AbandonCart marks cart as abandoned
func (h *Handler) AbandonCart(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	if err := h.service.AbandonCart(tenantID.(uuid.UUID), cartID); err != nil {
		if err == ErrCartNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to abandon cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart marked as abandoned"})
}

// ConvertCart marks cart as converted to order
func (h *Handler) ConvertCart(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	if err := h.service.ConvertCart(tenantID.(uuid.UUID), cartID); err != nil {
		if err == ErrCartNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart marked as converted"})
}