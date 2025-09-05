package category

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for category operations
type Handler struct {
	service Service
}

// NewHandler creates a new category handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers category routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	categories := router.Group("/categories")
	{
		categories.POST("", h.CreateCategory)
		categories.GET("", h.ListCategories)
		categories.GET("/tree", h.GetCategoryTree)
		categories.GET("/stats", h.GetCategoryStats)
		categories.GET("/featured", h.GetFeaturedCategories)
		categories.GET("/popular", h.GetPopularCategories)
		categories.PUT("/reorder", h.ReorderCategories)
		categories.PUT("/bulk/status", h.BulkUpdateStatus)
		categories.GET("/slug/:slug", h.GetCategoryBySlug)
		categories.GET("/:id", h.GetCategory)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
		categories.GET("/:id/path", h.GetCategoryPath)
		categories.PUT("/:id/move", h.MoveCategory)
		categories.GET("/:id/products", h.GetCategoryProducts)
		categories.POST("/:id/products/:product_id", h.AddProductToCategory)
		categories.DELETE("/:id/products/:product_id", h.RemoveProductFromCategory)
	}
}

// CreateCategory handles POST /categories
func (h *Handler) CreateCategory(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID from context
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Parse request body
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	// Create category
	category, err := h.service.CreateCategory(ctx, tenantID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, category)
}

// GetCategory handles GET /categories/{id}
func (h *Handler) GetCategory(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID", "details": err.Error()})
		return
	}
	
	// Get category
	category, err := h.service.GetCategory(ctx, tenantID, categoryID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, category)
}

// GetCategoryBySlug handles GET /categories/slug/{slug}
func (h *Handler) GetCategoryBySlug(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Extract slug
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug is required"})
		return
	}
	
	// Get category by slug
	category, err := h.service.GetCategoryBySlug(ctx, tenantID, slug)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, category)
}

// UpdateCategory handles PUT /categories/{id}
func (h *Handler) UpdateCategory(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID", "details": err.Error()})
		return
	}
	
	// Parse request body
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	// Update category
	category, err := h.service.UpdateCategory(ctx, tenantID, categoryID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, category)
}

// DeleteCategory handles DELETE /categories/{id}
func (h *Handler) DeleteCategory(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID", "details": err.Error()})
		return
	}
	
	// Delete category
	if err := h.service.DeleteCategory(ctx, tenantID, categoryID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ListCategories handles GET /categories
func (h *Handler) ListCategories(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Parse query parameters
	filter := h.parseCategoryFilter(c)
	limit, offset := h.parsePagination(c)
	
	// List categories
	categories, total, err := h.service.ListCategories(ctx, tenantID, filter, limit, offset)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	// Build response
	response := gin.H{
		"categories": categories,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	}
	
	c.JSON(http.StatusOK, response)
}

// GetCategoryTree handles GET /categories/tree
func (h *Handler) GetCategoryTree(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Parse parent ID if provided
	var parentID *uuid.UUID
	if parentIDStr := c.Query("parent_id"); parentIDStr != "" {
		parsedID, err := uuid.Parse(parentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent ID", "details": err.Error()})
			return
		}
		parentID = &parsedID
	}
	
	// Get category tree
	tree, err := h.service.GetCategoryTree(ctx, tenantID, parentID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, tree)
}

// GetCategoryPath handles GET /categories/{id}/path
func (h *Handler) GetCategoryPath(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID", "details": err.Error()})
		return
	}
	
	// Get category path
	path, err := h.service.GetCategoryPath(ctx, tenantID, categoryID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, path)
}

// MoveCategory handles PUT /categories/{id}/move
func (h *Handler) MoveCategory(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID", "details": err.Error()})
		return
	}
	
	// Parse request body
	var req struct {
		ParentID *uuid.UUID `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	// Move category
	if err := h.service.MoveCategory(ctx, tenantID, categoryID, req.ParentID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ReorderCategories handles PUT /categories/reorder
func (h *Handler) ReorderCategories(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Parse request body
	var req map[uuid.UUID]int
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	// Reorder categories
	if err := h.service.ReorderCategories(ctx, tenantID, req); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// BulkUpdateStatus handles PUT /categories/bulk/status
func (h *Handler) BulkUpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Parse request body
	var req struct {
		CategoryIDs []uuid.UUID   `json:"category_ids"`
		Status      CategoryStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	// Bulk update status
	if err := h.service.BulkUpdateStatus(ctx, tenantID, req.CategoryIDs, req.Status); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// AddProductToCategory handles POST /categories/{id}/products/{product_id}
func (h *Handler) AddProductToCategory(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID", "details": err.Error()})
		return
	}
	
	// Extract product ID
	productID, err := h.extractProductID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID", "details": err.Error()})
		return
	}
	
	// Add product to category
	if err := h.service.AddProductToCategory(ctx, tenantID, categoryID, productID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// RemoveProductFromCategory handles DELETE /categories/{id}/products/{product_id}
func (h *Handler) RemoveProductFromCategory(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID", "details": err.Error()})
		return
	}
	
	// Extract product ID
	productID, err := h.extractProductID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID", "details": err.Error()})
		return
	}
	
	// Remove product from category
	if err := h.service.RemoveProductFromCategory(ctx, tenantID, categoryID, productID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// GetCategoryProducts handles GET /categories/{id}/products
func (h *Handler) GetCategoryProducts(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID", "details": err.Error()})
		return
	}
	
	// Parse pagination
	limit, offset := h.parsePagination(c)
	
	// Get category products
	products, err := h.service.GetCategoryProducts(ctx, tenantID, categoryID, limit, offset)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, products)
}

// GetCategoryStats handles GET /categories/stats
func (h *Handler) GetCategoryStats(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Get category stats
	stats, err := h.service.GetCategoryStats(ctx, tenantID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// GetFeaturedCategories handles GET /categories/featured
func (h *Handler) GetFeaturedCategories(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Parse limit
	limit := 10 // default
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	// Get featured categories
	categories, err := h.service.GetFeaturedCategories(ctx, tenantID, limit)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, categories)
}

// GetPopularCategories handles GET /categories/popular
func (h *Handler) GetPopularCategories(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Parse limit
	limit := 10 // default
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	// Get popular categories
	categories, err := h.service.GetPopularCategories(ctx, tenantID, limit)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, categories)
}

// Helper methods

// extractTenantID extracts tenant ID from request context or headers
func (h *Handler) extractTenantID(c *gin.Context) (uuid.UUID, error) {
	// Try to get from context first (set by middleware)
	if tenantID := c.Request.Context().Value("tenant_id"); tenantID != nil {
		if id, ok := tenantID.(uuid.UUID); ok {
			return id, nil
		}
	}
	
	// Try to get from header
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		return uuid.Nil, fmt.Errorf("tenant ID not found")
	}
	
	return uuid.Parse(tenantIDStr)
}

// extractCategoryID extracts category ID from URL path
func (h *Handler) extractCategoryID(c *gin.Context) (uuid.UUID, error) {
	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		return uuid.Nil, fmt.Errorf("category ID not found")
	}
	
	return uuid.Parse(categoryIDStr)
}

// extractProductID extracts product ID from URL path
func (h *Handler) extractProductID(c *gin.Context) (uuid.UUID, error) {
	productIDStr := c.Param("product_id")
	if productIDStr == "" {
		return uuid.Nil, fmt.Errorf("product ID not found")
	}
	
	return uuid.Parse(productIDStr)
}

// parseCategoryFilter parses category filter from query parameters
func (h *Handler) parseCategoryFilter(c *gin.Context) CategoryFilter {
	filter := CategoryFilter{}
	
	if name := c.Query("name"); name != "" {
		filter.Name = name
	}
	
	if slug := c.Query("slug"); slug != "" {
		filter.Slug = slug
	}
	
	if status := c.Query("status"); status != "" {
		filter.Status = CategoryStatus(status)
	}
	
	if parentIDStr := c.Query("parent_id"); parentIDStr != "" {
		if parentID, err := uuid.Parse(parentIDStr); err == nil {
			filter.ParentID = &parentID
		}
	}
	
	if levelStr := c.Query("level"); levelStr != "" {
		if level, err := strconv.Atoi(levelStr); err == nil {
			filter.Level = &level
		}
	}
	
	if featuredStr := c.Query("is_featured"); featuredStr != "" {
		if featured, err := strconv.ParseBool(featuredStr); err == nil {
			filter.IsFeatured = &featured
		}
	}
	
	if menuStr := c.Query("show_in_menu"); menuStr != "" {
		if menu, err := strconv.ParseBool(menuStr); err == nil {
			filter.ShowInMenu = &menu
		}
	}
	
	return filter
}

// parsePagination parses pagination parameters from query string
func (h *Handler) parsePagination(c *gin.Context) (limit, offset int) {
	limit = 20 // default
	offset = 0 // default
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100 // max limit
			}
		}
	}
	
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	return limit, offset
}

// writeJSONResponse writes a JSON response (deprecated - use c.JSON directly)
func (h *Handler) writeJSONResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// writeErrorResponse writes an error response (deprecated - use c.JSON directly)
func (h *Handler) writeErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	errorResponse := gin.H{
		"error":  message,
		"status": statusCode,
	}
	
	// Include error details in development mode
	if err != nil {
		// In production, you might want to log this instead
		errorResponse["details"] = err.Error()
	}
	
	c.JSON(statusCode, errorResponse)
}

// handleServiceError handles service layer errors and maps them to HTTP status codes
func (h *Handler) handleServiceError(c *gin.Context, err error) {
	switch err {
	case ErrCategoryNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found", "details": err.Error()})
	case ErrCategoryExists:
		c.JSON(http.StatusConflict, gin.H{"error": "Category already exists", "details": err.Error()})
	case ErrInvalidParent:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent category", "details": err.Error()})
	case ErrCategoryHasChildren:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category has children", "details": err.Error()})
	case ErrCategoryHasProducts:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category has products", "details": err.Error()})
	case ErrMaxDepthExceeded:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum category depth exceeded", "details": err.Error()})
	case ErrInvalidSlug:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slug", "details": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
	}
}