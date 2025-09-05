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
	service Service
}

// NewHandler creates a new cart handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers cart routes according to API reference
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Cart Management (5 endpoints)
	router.GET("/cart", h.GetCart)                    // GET /cart
	router.POST("/cart/items", h.AddItem)             // POST /cart/items
	router.PATCH("/cart/items/:id", h.UpdateCartItem) // PATCH /cart/items/:id
	router.PATCH("/cart", h.UpdateCart)               // PATCH /cart
	router.POST("/cart/estimates", h.GetEstimates)    // POST /cart/estimates
	
	// Guest Cart & Checkout (3 endpoints)
	router.GET("/cart/guest", h.GetGuestCart)         // GET /cart/guest
	router.PATCH("/cart/guest", h.UpdateGuestCart)   // PATCH /cart/guest
	router.POST("/checkout/guest", h.ProcessGuestCheckout) // POST /checkout/guest
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

// AddItem adds an item to the cart
func (h *Handler) AddItem(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user authentication required"})
		return
	}

	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user's cart
	userUUID := userID.(uuid.UUID)
	cart, err := h.service.GetCartByCustomer(tenantID.(uuid.UUID), userUUID)
	if err != nil {
		if err == ErrCartNotFound {
			// Create new cart if none exists
			cart, err = h.service.CreateCart(tenantID.(uuid.UUID), CreateCartRequest{
				CustomerID: &userUUID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
			return
		}
	}

	item, err := h.service.AddItem(tenantID.(uuid.UUID), cart.ID, req)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "inventory") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

// GetCart retrieves current user's cart with optional includes
func (h *Handler) GetCart(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user authentication required"})
		return
	}

	// Parse include parameter
	include := c.Query("include")
	
	userUUID := userID.(uuid.UUID)
	cart, err := h.service.GetCartByCustomer(tenantID.(uuid.UUID), userUUID)
	if err != nil {
		if err == ErrCartNotFound {
			// Create new cart if none exists
			cart, err = h.service.CreateCart(tenantID.(uuid.UUID), CreateCartRequest{
				CustomerID: &userUUID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
			return
		}
	}

	// Handle include options
	response := gin.H{"data": cart}
	if strings.Contains(include, "summary") {
		summary, _ := h.service.GetCartSummary(tenantID.(uuid.UUID), cart.ID)
		response["summary"] = summary
	}
	if strings.Contains(include, "shipping_methods") {
		// Add shipping methods logic here
		response["shipping_methods"] = []interface{}{}
	}
	if strings.Contains(include, "taxes") {
		// Add tax calculation logic here
		response["taxes"] = gin.H{"total": 0}
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCartItem updates a specific item in the cart
func (h *Handler) UpdateCartItem(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user authentication required"})
		return
	}

	itemIDStr := c.Param("id")
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

	// Get user's cart
	cart, err := h.service.GetCartByCustomer(tenantID.(uuid.UUID), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	item, err := h.service.UpdateItem(tenantID.(uuid.UUID), cart.ID, itemID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// UpdateCart updates cart properties (shipping, billing, etc.)
func (h *Handler) UpdateCart(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user authentication required"})
		return
	}

	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user's cart
	cart, err := h.service.GetCartByCustomer(tenantID.(uuid.UUID), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	updatedCart, err := h.service.UpdateCart(tenantID.(uuid.UUID), cart.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedCart})
}

// GetEstimates calculates shipping and tax estimates
func (h *Handler) GetEstimates(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user authentication required"})
		return
	}

	var req EstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user's cart
	cart, err := h.service.GetCartByCustomer(tenantID.(uuid.UUID), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	estimates, err := h.service.GetEstimates(tenantID.(uuid.UUID), cart.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get estimates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": estimates})
}

// GetGuestCart retrieves guest cart by session
func (h *Handler) GetGuestCart(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	include := c.Query("include")

	cart, err := h.service.GetCartBySession(tenantID.(uuid.UUID), sessionID)
	if err != nil {
		if err == ErrCartNotFound {
			// Create new guest cart
			cart, err = h.service.CreateCart(tenantID.(uuid.UUID), CreateCartRequest{
				SessionID: sessionID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
			return
		}
	}

	// Handle include options
	response := gin.H{"data": cart}
	if strings.Contains(include, "summary") {
		summary, _ := h.service.GetCartSummary(tenantID.(uuid.UUID), cart.ID)
		response["summary"] = summary
	}

	c.JSON(http.StatusOK, response)
}

// UpdateGuestCart updates guest cart
func (h *Handler) UpdateGuestCart(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.GetCartBySession(tenantID.(uuid.UUID), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	updatedCart, err := h.service.UpdateCart(tenantID.(uuid.UUID), cart.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedCart})
}

// ProcessGuestCheckout processes guest checkout
func (h *Handler) ProcessGuestCheckout(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	var req GuestCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ProcessGuestCheckout(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process checkout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
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