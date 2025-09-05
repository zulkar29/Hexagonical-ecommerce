package wishlist

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for wishlist operations
type Handler struct {
	service Service
}

// NewHandler creates a new wishlist handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers wishlist routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Wishlist routes
	wishlists := router.Group("/wishlists")
	{
		wishlists.POST("", h.CreateWishlist)
		wishlists.GET("", h.ListWishlists)
		wishlists.GET("/stats", h.GetWishlistStats)
		wishlists.GET("/popular", h.GetPopularWishlists)
		wishlists.POST("/cleanup/empty", h.CleanupEmptyWishlists)
		wishlists.GET("/share/:shareToken", h.GetWishlistByShareToken)
		wishlists.GET("/:wishlistID", h.GetWishlist)
		wishlists.PUT("/:wishlistID", h.UpdateWishlist)
		wishlists.DELETE("/:wishlistID", h.DeleteWishlist)
		wishlists.POST("/:wishlistID/share", h.ShareWishlist)
		wishlists.GET("/:wishlistID/items", h.ListWishlistItems)
		wishlists.POST("/:wishlistID/items", h.AddItem)
		wishlists.POST("/:wishlistID/items/bulk", h.BulkAddItems)
		wishlists.POST("/:wishlistID/items/reorder", h.ReorderItems)
		wishlists.POST("/:wishlistID/clear", h.ClearWishlist)
		wishlists.POST("/:sourceID/merge/:targetID", h.MergeWishlists)
	}
	
	// Customer wishlist routes
	customers := router.Group("/customers")
	{
		customers.GET("/:customerID/wishlists", h.GetCustomerWishlists)
		customers.POST("/:customerID/wishlists", h.CreateCustomerWishlist)
		customers.GET("/:customerID/wishlists/default", h.GetDefaultWishlist)
		customers.POST("/:customerID/wishlists/:wishlistID/default", h.SetDefaultWishlist)
		customers.GET("/:customerID/wishlist-activity", h.GetCustomerActivity)
	}
	
	// Wishlist item routes
	wishlistItems := router.Group("/wishlist-items")
	{
		wishlistItems.GET("", h.ListItems)
		wishlistItems.DELETE("/bulk", h.BulkRemoveItems)
		wishlistItems.PUT("/bulk/priority", h.BulkUpdateItemPriority)
		wishlistItems.POST("/cleanup/orphaned", h.CleanupOrphanedItems)
		wishlistItems.GET("/:itemID", h.GetItem)
		wishlistItems.PUT("/:itemID", h.UpdateItem)
		wishlistItems.DELETE("/:itemID", h.RemoveItem)
		wishlistItems.POST("/:itemID/move", h.MoveItem)
		wishlistItems.POST("/:itemID/copy", h.CopyItem)
	}
	
	// Products routes
	products := router.Group("/products")
	{
		products.GET("/most-wished", h.GetMostWishedProducts)
	}
}

// Wishlist operations

// CreateWishlist creates a new wishlist
func (h *Handler) CreateWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	var req CreateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	req.TenantID = tenantID
	
	response, err := h.service.CreateWishlist(ctx, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, response)
}

// CreateCustomerWishlist creates a new wishlist for a specific customer
func (h *Handler) CreateCustomerWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	customerID, err := h.getUUIDParam(c, "customerID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
		return
	}
	
	var req CreateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	req.TenantID = tenantID
	req.CustomerID = customerID
	
	response, err := h.service.CreateWishlist(ctx, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, response)
}

// GetWishlist retrieves a wishlist by ID
func (h *Handler) GetWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	wishlistID, err := h.getUUIDParam(c, "wishlistID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID", "details": err.Error()})
		return
	}
	
	response, err := h.service.GetWishlist(ctx, tenantID, wishlistID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// GetWishlistByShareToken retrieves a public wishlist by share token
func (h *Handler) GetWishlistByShareToken(c *gin.Context) {
	ctx := c.Request.Context()
	
	shareToken := c.Param("shareToken")
	if shareToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Share token is required"})
		return
	}
	
	response, err := h.service.GetWishlistByShareToken(ctx, shareToken)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// UpdateWishlist updates an existing wishlist
func (h *Handler) UpdateWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	wishlistID, err := h.getUUIDParam(c, "wishlistID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID", "details": err.Error()})
		return
	}
	
	var req UpdateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	response, err := h.service.UpdateWishlist(ctx, tenantID, wishlistID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// DeleteWishlist deletes a wishlist
func (h *Handler) DeleteWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	wishlistID, err := h.getUUIDParam(c, "wishlistID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID", "details": err.Error()})
		return
	}
	
	if err := h.service.DeleteWishlist(ctx, tenantID, wishlistID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ListWishlists returns paginated wishlists
func (h *Handler) ListWishlists(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	// Parse query parameters
	filter := h.parseWishlistFilter(c)
	limit, offset := h.parsePagination(c)
	
	wishlists, total, err := h.service.ListWishlists(ctx, tenantID, filter, limit, offset)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	response := map[string]interface{}{
		"wishlists": wishlists,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	}
	
	c.JSON(http.StatusOK, response)
}

// ShareWishlist makes a wishlist public/private
func (h *Handler) ShareWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	wishlistID, err := h.getUUIDParam(c, "wishlistID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID", "details": err.Error()})
		return
	}
	
	var req struct {
		IsPublic bool `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	shareURL, err := h.service.ShareWishlist(ctx, tenantID, wishlistID, req.IsPublic)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	response := map[string]interface{}{
		"is_public": req.IsPublic,
		"share_url": shareURL,
	}
	
	c.JSON(http.StatusOK, response)
}

// Customer wishlist operations

// GetCustomerWishlists returns all wishlists for a customer
func (h *Handler) GetCustomerWishlists(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	customerID, err := h.getUUIDParam(c, "customerID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
		return
	}
	
	wishlists, err := h.service.GetCustomerWishlists(ctx, tenantID, customerID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, wishlists)
}

// GetDefaultWishlist returns the default wishlist for a customer
func (h *Handler) GetDefaultWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	customerID, err := h.getUUIDParam(c, "customerID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
		return
	}
	
	wishlist, err := h.service.GetDefaultWishlist(ctx, tenantID, customerID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, wishlist)
}

// SetDefaultWishlist sets a wishlist as default for a customer
func (h *Handler) SetDefaultWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	customerID, err := h.getUUIDParam(c, "customerID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
		return
	}
	
	wishlistID, err := h.getUUIDParam(c, "wishlistID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID", "details": err.Error()})
		return
	}
	
	if err := h.service.SetDefaultWishlist(ctx, tenantID, customerID, wishlistID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// Wishlist item operations

// AddItem adds an item to a wishlist
func (h *Handler) AddItem(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	wishlistID, err := h.getUUIDParam(c, "wishlistID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID", "details": err.Error()})
		return
	}
	
	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	req.WishlistID = wishlistID
	
	response, err := h.service.AddItem(ctx, tenantID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, response)
}

// GetItem retrieves a wishlist item by ID
func (h *Handler) GetItem(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	itemID, err := h.getUUIDParam(c, "itemID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID", "details": err.Error()})
		return
	}
	
	response, err := h.service.GetItem(ctx, tenantID, itemID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// UpdateItem updates a wishlist item
func (h *Handler) UpdateItem(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	itemID, err := h.getUUIDParam(c, "itemID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID", "details": err.Error()})
		return
	}
	
	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	response, err := h.service.UpdateItem(ctx, tenantID, itemID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// RemoveItem removes an item from a wishlist
func (h *Handler) RemoveItem(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	itemID, err := h.getUUIDParam(c, "itemID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID", "details": err.Error()})
		return
	}
	
	if err := h.service.RemoveItem(ctx, tenantID, itemID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ListItems returns paginated wishlist items
func (h *Handler) ListItems(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	// Parse query parameters
	filter := h.parseWishlistItemFilter(c)
	limit, offset := h.parsePagination(c)
	
	items, total, err := h.service.ListItems(ctx, tenantID, filter, limit, offset)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	response := map[string]interface{}{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
	
	c.JSON(http.StatusOK, response)
}

// ListWishlistItems returns items for a specific wishlist
func (h *Handler) ListWishlistItems(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	wishlistID, err := h.getUUIDParam(c, "wishlistID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID", "details": err.Error()})
		return
	}
	
	// Parse query parameters
	filter := h.parseWishlistItemFilter(c)
	filter.WishlistID = &wishlistID
	limit, offset := h.parsePagination(c)
	
	items, total, err := h.service.ListItems(ctx, tenantID, filter, limit, offset)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	response := map[string]interface{}{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
	
	c.JSON(http.StatusOK, response)
}

// MoveItem moves an item to another wishlist
func (h *Handler) MoveItem(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	itemID, err := h.getUUIDParam(c, "itemID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID", "details": err.Error()})
		return
	}
	
	var req struct {
		TargetWishlistID uuid.UUID `json:"target_wishlist_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	if err := h.service.MoveItem(ctx, tenantID, itemID, req.TargetWishlistID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// CopyItem copies an item to another wishlist
func (h *Handler) CopyItem(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	itemID, err := h.getUUIDParam(c, "itemID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID", "details": err.Error()})
		return
	}
	
	var req struct {
		TargetWishlistID uuid.UUID `json:"target_wishlist_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	if err := h.service.CopyItem(ctx, tenantID, itemID, req.TargetWishlistID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ClearWishlist removes all items from a wishlist
func (h *Handler) ClearWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	wishlistID, err := h.getUUIDParam(c, "wishlistID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID", "details": err.Error()})
		return
	}
	
	if err := h.service.ClearWishlist(ctx, tenantID, wishlistID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// Bulk operations

// BulkAddItems adds multiple items to a wishlist
func (h *Handler) BulkAddItems(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	wishlistID, err := h.getUUIDParam(c, "wishlistID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID", "details": err.Error()})
		return
	}
	
	var req struct {
		Items []AddItemRequest `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	// Set wishlist ID for all items
	for i := range req.Items {
		req.Items[i].WishlistID = wishlistID
	}
	
	responses, err := h.service.BulkAddItems(ctx, tenantID, req.Items)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, responses)
}

// BulkRemoveItems removes multiple items from wishlists
func (h *Handler) BulkRemoveItems(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	var req struct {
		ItemIDs []uuid.UUID `json:"item_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	if err := h.service.BulkRemoveItems(ctx, tenantID, req.ItemIDs); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// BulkUpdateItemPriority updates priority for multiple items
func (h *Handler) BulkUpdateItemPriority(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	var req struct {
		Updates map[uuid.UUID]int `json:"updates"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	if err := h.service.BulkUpdateItemPriority(ctx, tenantID, req.Updates); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ReorderItems reorders items in a wishlist
func (h *Handler) ReorderItems(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	wishlistID, err := h.getUUIDParam(c, "wishlistID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID", "details": err.Error()})
		return
	}
	
	var req struct {
		ItemOrder []uuid.UUID `json:"item_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	if err := h.service.ReorderItems(ctx, tenantID, wishlistID, req.ItemOrder); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// Wishlist management

// MergeWishlists merges source wishlist into target wishlist
func (h *Handler) MergeWishlists(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	sourceID, err := h.getUUIDParam(c, "sourceID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source wishlist ID", "details": err.Error()})
		return
	}
	
	targetID, err := h.getUUIDParam(c, "targetID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target wishlist ID", "details": err.Error()})
		return
	}
	
	if err := h.service.MergeWishlists(ctx, tenantID, sourceID, targetID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// Analytics and statistics

// GetWishlistStats returns wishlist statistics
func (h *Handler) GetWishlistStats(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	stats, err := h.service.GetWishlistStats(ctx, tenantID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// GetMostWishedProducts returns the most wished products
func (h *Handler) GetMostWishedProducts(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	limit := h.getIntParam(c, "limit", 10)
	
	products, err := h.service.GetMostWishedProducts(ctx, tenantID, limit)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, products)
}

// GetCustomerActivity returns customer wishlist activity
func (h *Handler) GetCustomerActivity(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	customerID, err := h.getUUIDParam(c, "customerID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
		return
	}
	
	days := h.getIntParam(c, "days", 30)
	
	activity, err := h.service.GetCustomerActivity(ctx, tenantID, customerID, days)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, activity)
}

// GetPopularWishlists returns popular public wishlists
func (h *Handler) GetPopularWishlists(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	limit := h.getIntParam(c, "limit", 10)
	
	wishlists, err := h.service.GetPopularWishlists(ctx, tenantID, limit)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, wishlists)
}

// Maintenance operations

// CleanupEmptyWishlists removes empty wishlists older than specified days
func (h *Handler) CleanupEmptyWishlists(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	days := h.getIntParam(c, "days", 30)
	
	count, err := h.service.CleanupEmptyWishlists(ctx, tenantID, days)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	response := map[string]interface{}{
		"deleted_count": count,
	}
	
	c.JSON(http.StatusOK, response)
}

// CleanupOrphanedItems removes items for non-existent products
func (h *Handler) CleanupOrphanedItems(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := h.getTenantID(c)
	
	count, err := h.service.CleanupOrphanedItems(ctx, tenantID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	response := map[string]interface{}{
		"deleted_count": count,
	}
	
	c.JSON(http.StatusOK, response)
}

// Helper methods

// getTenantID extracts tenant ID from request context or headers
func (h *Handler) getTenantID(c *gin.Context) uuid.UUID {
	// This would typically come from JWT token or request context
	// For now, return a placeholder
	return uuid.New()
}

// getUUIDParam extracts UUID parameter from URL
func (h *Handler) getUUIDParam(c *gin.Context, param string) (uuid.UUID, error) {
	value := c.Param(param)
	return uuid.Parse(value)
}

// getIntParam extracts integer parameter from query string
func (h *Handler) getIntParam(c *gin.Context, param string, defaultValue int) int {
	value := c.Query(param)
	if value == "" {
		return defaultValue
	}
	
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	
	return defaultValue
}

// parsePagination extracts pagination parameters
func (h *Handler) parsePagination(c *gin.Context) (limit, offset int) {
	limit = h.getIntParam(c, "limit", 20)
	offset = h.getIntParam(c, "offset", 0)
	
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	
	return limit, offset
}

// parseWishlistFilter extracts wishlist filter parameters
func (h *Handler) parseWishlistFilter(c *gin.Context) WishlistFilter {
	filter := WishlistFilter{}
	
	if customerIDStr := c.Query("customer_id"); customerIDStr != "" {
		if customerID, err := uuid.Parse(customerIDStr); err == nil {
			filter.CustomerID = &customerID
		}
	}
	
	if name := c.Query("name"); name != "" {
		filter.Name = name
	}
	
	if isDefaultStr := c.Query("is_default"); isDefaultStr != "" {
		if isDefault, err := strconv.ParseBool(isDefaultStr); err == nil {
			filter.IsDefault = &isDefault
		}
	}
	
	if isPublicStr := c.Query("is_public"); isPublicStr != "" {
		if isPublic, err := strconv.ParseBool(isPublicStr); err == nil {
			filter.IsPublic = &isPublic
		}
	}
	
	if isEmptyStr := c.Query("is_empty"); isEmptyStr != "" {
		if isEmpty, err := strconv.ParseBool(isEmptyStr); err == nil {
			filter.IsEmpty = &isEmpty
		}
	}
	
	return filter
}

// parseWishlistItemFilter extracts wishlist item filter parameters
func (h *Handler) parseWishlistItemFilter(c *gin.Context) WishlistItemFilter {
	filter := WishlistItemFilter{}
	
	if wishlistIDStr := c.Query("wishlist_id"); wishlistIDStr != "" {
		if wishlistID, err := uuid.Parse(wishlistIDStr); err == nil {
			filter.WishlistID = &wishlistID
		}
	}
	
	if productIDStr := c.Query("product_id"); productIDStr != "" {
		if productID, err := uuid.Parse(productIDStr); err == nil {
			filter.ProductID = &productID
		}
	}
	
	if variantIDStr := c.Query("variant_id"); variantIDStr != "" {
		if variantID, err := uuid.Parse(variantIDStr); err == nil {
			filter.VariantID = &variantID
		}
	}
	
	if minPriorityStr := c.Query("min_priority"); minPriorityStr != "" {
		if minPriority, err := strconv.Atoi(minPriorityStr); err == nil {
			filter.MinPriority = &minPriority
		}
	}
	
	if maxPriorityStr := c.Query("max_priority"); maxPriorityStr != "" {
		if maxPriority, err := strconv.Atoi(maxPriorityStr); err == nil {
			filter.MaxPriority = &maxPriority
		}
	}
	
	return filter
}

// writeJSON writes JSON response
func (h *Handler) writeJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// writeError writes error response
func (h *Handler) writeError(c *gin.Context, status int, message string, err error) {
	response := map[string]interface{}{
		"error": message,
	}
	
	if err != nil {
		response["details"] = err.Error()
	}
	
	c.JSON(status, response)
}

// handleServiceError handles service layer errors
func (h *Handler) handleServiceError(c *gin.Context, err error) {
	switch err {
	case ErrWishlistNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Wishlist not found", "details": err.Error()})
	case ErrWishlistItemNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Wishlist item not found", "details": err.Error()})
	case ErrWishlistNameExists:
		c.JSON(http.StatusConflict, gin.H{"error": "Wishlist name already exists", "details": err.Error()})
	case ErrWishlistLimitExceeded:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wishlist limit exceeded", "details": err.Error()})
	case ErrWishlistFull:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wishlist is full", "details": err.Error()})
	case ErrCannotDeleteDefaultWishlist:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete default wishlist", "details": err.Error()})
	case ErrInvalidTenantID, ErrInvalidCustomerID, ErrInvalidWishlistID, ErrInvalidProductID:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID provided", "details": err.Error()})
	case ErrInvalidWishlistName, ErrWishlistNameTooLong, ErrWishlistDescriptionTooLong:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist data", "details": err.Error()})
	case ErrInvalidQuantity, ErrInvalidPriority, ErrItemNotesTooLong:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item data", "details": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
	}
}