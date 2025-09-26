package product

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ecommerce-saas/internal/shared/handlers"
)

// Handler handles HTTP requests for product operations
type Handler struct {
	service *Service
}

// NewHandler creates a new product handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateProduct handles POST /api/products
func (h *Handler) CreateProduct(c *gin.Context) {
	// Extract tenant ID from context (set by middleware)
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	var product Product
	if bindErr := c.ShouldBindJSON(&product); bindErr != nil {
		handlers.HandleValidationError(c, bindErr)
		return
	}

	createdProduct, err := h.service.CreateProduct(tenantID, &product)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithCreated(c, createdProduct, "Product created successfully")
}

// GetProduct handles GET /api/products/:id
func (h *Handler) GetProduct(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	productIDStr := c.Param("id")
	// Note: Service already handles string to UUID conversion
	product, err := h.service.GetProduct(tenantID, productIDStr)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, product)
}

// GetProductBySlug handles GET /api/products/slug/:slug
func (h *Handler) GetProductBySlug(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	slug := c.Param("slug")

	product, err := h.service.GetProductBySlug(tenantID.(uuid.UUID), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

// UpdateProduct handles PUT /api/products/:id with action support
func (h *Handler) UpdateProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Check for specific actions
	action := c.Query("action")
	switch action {
	case "update_inventory":
		var req struct {
			Quantity int `json:"quantity" binding:"required"`
		}
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		err = h.service.UpdateInventory(tenantID.(uuid.UUID), productID.String(), req.Quantity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Inventory updated successfully"})
		return
	case "update_status":
		var req struct {
			Status ProductStatus `json:"status" binding:"required"`
		}
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		err = h.service.UpdateProductStatus(tenantID.(uuid.UUID), productID.String(), req.Status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product status updated successfully"})
		return
	case "duplicate":
		product, duplicateErr := h.service.DuplicateProduct(tenantID.(uuid.UUID), productID.String())
		if duplicateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": duplicateErr.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Product duplicated successfully",
			"data":    product,
		})
		return
	}

	// Regular product update
	var product Product
	if bindErr := c.ShouldBindJSON(&product); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	updatedProduct, err := h.service.UpdateProduct(tenantID.(uuid.UUID), productID, &product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data":    updatedProduct,
	})
}

// ListProducts handles GET /api/products with analytics and export support
func (h *Handler) ListProducts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	// Check for analytics or export type
	queryType := c.Query("type")
	switch queryType {
	case "stats":
		stats, err := h.service.GetProductStats(tenantID.(uuid.UUID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product statistics"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": stats})
		return
	case "low-stock":
		thresholdStr := c.DefaultQuery("threshold", "10")
		threshold, err := strconv.Atoi(thresholdStr)
		if err != nil || threshold <= 0 {
			threshold = 10
		}
		products, err := h.service.GetLowStockProducts(tenantID.(uuid.UUID), threshold)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch low stock products"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"products":  products,
				"threshold": threshold,
				"count":     len(products),
			},
		})
		return
	case "search":
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
			return
		}
		// Parse pagination for search
		offsetStr := c.DefaultQuery("offset", "0")
		limitStr := c.DefaultQuery("limit", "20")
		offset, offsetErr := strconv.Atoi(offsetStr)
		if offsetErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
		limit, limitErr := strconv.Atoi(limitStr)
		if limitErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
		if limit > 100 {
			limit = 100
		}
		if limit < 1 {
			limit = 20
		}
		products, total, err := h.service.SearchProducts(tenantID.(uuid.UUID), query, offset, limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"products": products,
				"total":    total,
				"offset":   offset,
				"limit":    limit,
				"query":    query,
			},
		})
		return
	}

	// Regular product listing with filters
	var filter ProductListFilter

	if status := c.Query("status"); status != "" {
		filter.Status = ProductStatus(status)
	}
	if productType := c.Query("product_type"); productType != "" {
		filter.Type = ProductType(productType)
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, parseErr := uuid.Parse(categoryID); parseErr == nil {
			filter.CategoryID = &id
		}
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, parseErr := strconv.ParseFloat(minPrice, 64); parseErr == nil {
			filter.MinPrice = &price
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, parseErr := strconv.ParseFloat(maxPrice, 64); parseErr == nil {
			filter.MaxPrice = &price
		}
	}
	if inStock := c.Query("in_stock"); inStock != "" {
		if stock, parseErr := strconv.ParseBool(inStock); parseErr == nil {
			filter.InStock = &stock
		}
	}
	filter.Search = c.Query("search")

	// Parse pagination
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")

	offset, offsetErr := strconv.Atoi(offsetStr)
	if offsetErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	limit, limitErr := strconv.Atoi(limitStr)
	if limitErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	// Validate limit
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}

	products, total, err := h.service.ListProducts(tenantID.(uuid.UUID), filter, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"products": products,
			"total":    total,
			"offset":   offset,
			"limit":    limit,
			"filter":   filter,
		},
	})
}

// DeleteProduct handles DELETE /api/products/:id
func (h *Handler) DeleteProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = h.service.DeleteProduct(tenantID.(uuid.UUID), productID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

// Category Handlers

// CreateCategory handles POST /api/categories
// CreateCategory method removed - handled by category module

// GetCategory and ListCategories methods removed - handled by category module

// ListCategories method removed - handled by category module

// Product Variant Handlers

// CreateProductVariant handles POST /api/products/:id/variants
func (h *Handler) CreateProductVariant(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var variant ProductVariant
	if bindErr := c.ShouldBindJSON(&variant); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	createdVariant, err := h.service.CreateProductVariant(tenantID.(uuid.UUID), productID, &variant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product variant created successfully",
		"data":    createdVariant,
	})
}

// GetProductVariants handles GET /api/products/:id/variants
func (h *Handler) GetProductVariants(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productIDStr := c.Param("id")
	productID, parseErr := uuid.Parse(productIDStr)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	variants, err := h.service.GetProductVariants(tenantID.(uuid.UUID), productID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": variants,
	})
}

// UpdateProductVariant handles PUT /api/products/:id/variants/:variantId
func (h *Handler) UpdateProductVariant(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productIDStr := c.Param("id")
	productID, productParseErr := uuid.Parse(productIDStr)
	if productParseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	variantIDStr := c.Param("variantId")
	variantID, variantParseErr := uuid.Parse(variantIDStr)
	if variantParseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant ID"})
		return
	}

	var variant ProductVariant
	if bindErr := c.ShouldBindJSON(&variant); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	updatedVariant, err := h.service.UpdateProductVariant(tenantID.(uuid.UUID), productID, variantID, &variant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product variant updated successfully",
		"data":    updatedVariant,
	})
}

// DeleteProductVariant handles DELETE /api/products/:id/variants/:variantId
func (h *Handler) DeleteProductVariant(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productIDStr := c.Param("id")
	productID, productParseErr := uuid.Parse(productIDStr)
	if productParseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	variantIDStr := c.Param("variantId")
	variantID, variantParseErr := uuid.Parse(variantIDStr)
	if variantParseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant ID"})
		return
	}

	deleteErr := h.service.DeleteProductVariant(tenantID.(uuid.UUID), productID.String(), variantID.String())
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": deleteErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product variant deleted successfully",
	})
}

// Enhanced Category Handlers

// UpdateCategory and DeleteCategory methods removed - handled by category module

// SearchProducts handles GET /api/products/search
func (h *Handler) SearchProducts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Parse pagination
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	if limit > 100 {
		limit = 100
	}

	products, total, err := h.service.SearchProducts(tenantID.(uuid.UUID), query, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"products": products,
			"total":    total,
			"offset":   offset,
			"limit":    limit,
		},
	})
}

// GetRootCategories method removed - handled by category module

// GetCategoryChildren method removed - handled by category module

// HandleProductOperations handles POST /api/products/operations for bulk operations
func (h *Handler) HandleProductOperations(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	operation := c.Query("operation")
	switch operation {
	case "bulk_update":
		var req struct {
			ProductIDs []string               `json:"product_ids" binding:"required"`
			Updates    map[string]interface{} `json:"updates" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		err := h.service.BulkUpdateProducts(tenantID.(uuid.UUID), req.ProductIDs, req.Updates)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Products updated successfully"})
		return
	case "import":
		// TODO: Implement product import functionality
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Import functionality not implemented"})
		return
	case "export":
		// TODO: Implement product export functionality
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Export functionality not implemented"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type"})
	}
}

// RegisterRoutes registers all product-related routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Product routes
	products := router.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.GET("", h.ListProducts)                        // Supports ?type=stats|low-stock|search, ?q=query
		products.POST("/operations", h.HandleProductOperations) // Supports ?operation=bulk_update|import|export
		products.GET("/:id", h.GetProduct)
		products.PUT("/:id", h.UpdateProduct) // Supports ?action=update_inventory|update_status|duplicate
		products.DELETE("/:id", h.DeleteProduct)
		products.GET("/slug/:slug", h.GetProductBySlug)

		// Bulk operations
		products.POST("/bulk", h.HandleProductBulk)

		// Product images
		products.POST("/:id/images", h.UploadProductImages)
		products.DELETE("/:id/images/:image-id", h.DeleteProductImage)

		// Product analytics
		products.GET("/:id/analytics", h.GetProductAnalytics)

		// Product variant routes
		products.POST("/:id/variants", h.CreateProductVariant)
		products.GET("/:id/variants", h.GetProductVariants)
		products.PUT("/:id/variants/:variant_id", h.UpdateProductVariant)
		products.DELETE("/:id/variants/:variant_id", h.DeleteProductVariant)
	}

	// Category routes
	// Category routes removed - handled by category module to avoid conflicts

	// Note: Public routes are registered separately in routes.go setupPublicProductRoutes
	// to avoid duplicate registration conflicts
}

// Missing Product Endpoints Implementation

// HandleProductBulk handles POST /api/products/bulk
func (h *Handler) HandleProductBulk(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	operation := c.Query("operation")
	switch operation {
	case "import":
		// TODO: Implement CSV/Excel import
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Import functionality not implemented"})
	case "export":
		// TODO: Implement CSV/Excel export
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Export functionality not implemented"})
	case "update":
		var req struct {
			ProductIDs []string               `json:"product_ids" binding:"required"`
			Updates    map[string]interface{} `json:"updates" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		err := h.service.BulkUpdateProducts(tenantID.(uuid.UUID), req.ProductIDs, req.Updates)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Products updated successfully"})
	case "delete":
		var req struct {
			ProductIDs []string `json:"product_ids" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		// Convert string IDs to UUIDs
		productUUIDs := make([]uuid.UUID, 0, len(req.ProductIDs))
		for _, idStr := range req.ProductIDs {
			if id, err := uuid.Parse(idStr); err == nil {
				productUUIDs = append(productUUIDs, id)
			}
		}

		if len(productUUIDs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No valid product IDs provided"})
			return
		}

		err := h.service.BulkDeleteProducts(tenantID.(uuid.UUID), productUUIDs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Products deleted successfully"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type"})
	}
}

// UploadProductImages handles POST /api/products/:id/images
func (h *Handler) UploadProductImages(c *gin.Context) {
	_, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productIDStr := c.Param("id")
	_, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// TODO: Implement image upload functionality
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Image upload functionality not implemented"})
}

// DeleteProductImage handles DELETE /api/products/:id/images/:image-id
func (h *Handler) DeleteProductImage(c *gin.Context) {
	_, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productIDStr := c.Param("id")
	_, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	imageID := c.Param("image-id")
	if imageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image ID is required"})
		return
	}

	// TODO: Implement image deletion functionality
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Image deletion functionality not implemented"})
}

// GetProductAnalytics handles GET /api/products/:id/analytics
func (h *Handler) GetProductAnalytics(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	analyticsType := c.DefaultQuery("type", "performance")
	analytics, err := h.service.GetProductAnalytics(tenantID.(uuid.UUID), productID.String(), analyticsType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": analytics,
		"type": analyticsType,
	})
}

// Public Product Endpoints

// GetPublicProducts handles GET /public/products
func (h *Handler) GetPublicProducts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	// Parse filters for public access
	var filter ProductListFilter
	filter.Status = "active" // Only show active products publicly

	if categoryID := c.Query("category"); categoryID != "" {
		if id, err := uuid.Parse(categoryID); err == nil {
			filter.CategoryID = &id
		}
	}
	filter.Search = c.Query("search")

	// Parse pagination
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}

	products, total, err := h.service.ListProducts(tenantID.(uuid.UUID), filter, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"products": products,
			"total":    total,
			"offset":   offset,
			"limit":    limit,
		},
	})
}

// GetPublicProduct handles GET /public/products/:id
func (h *Handler) GetPublicProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productIDStr := c.Param("id")
	product, err := h.service.GetProduct(tenantID.(uuid.UUID), productIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Only return active products for public access
	if product.Status != "active" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	include := c.Query("include")
	responseData := gin.H{"product": product}

	if strings.Contains(include, "variants") {
		variants, variantErr := h.service.GetProductVariants(tenantID.(uuid.UUID), productIDStr)
		if variantErr == nil {
			responseData["variants"] = variants
		} else {
			responseData["variants"] = []interface{}{} // Empty array if error occurs
		}
	}

	// TODO: Add reviews and related products when those modules are implemented
	if strings.Contains(include, "reviews") {
		responseData["reviews"] = []interface{}{} // Placeholder
	}

	if strings.Contains(include, "related") {
		responseData["related"] = []interface{}{} // Placeholder
	}

	c.JSON(http.StatusOK, gin.H{"data": responseData})
}

// GetPublicCategories method removed - handled by category module
