package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	createdProduct, err := h.service.CreateProduct(tenantID.(uuid.UUID), product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    createdProduct,
	})
}

// GetProduct handles GET /api/products/:id
func (h *Handler) GetProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productID := c.Param("id")
	
	product, err := h.service.GetProduct(tenantID.(uuid.UUID), productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
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

// UpdateProduct handles PUT /api/products/:id
func (h *Handler) UpdateProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productID := c.Param("id")
	
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	updatedProduct, err := h.service.UpdateProduct(tenantID.(uuid.UUID), productID, product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data":    updatedProduct,
	})
}

// ListProducts handles GET /api/products
func (h *Handler) ListProducts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	// Parse query parameters
	var filter ProductListFilter
	
	if status := c.Query("status"); status != "" {
		filter.Status = ProductStatus(status)
	}
	if productType := c.Query("type"); productType != "" {
		filter.Type = ProductType(productType)
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := uuid.Parse(categoryID); err == nil {
			filter.CategoryID = &id
		}
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = &price
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = &price
		}
	}
	if inStock := c.Query("in_stock"); inStock != "" {
		if stock, err := strconv.ParseBool(inStock); err == nil {
			filter.InStock = &stock
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

	productID := c.Param("id")
	
	err := h.service.DeleteProduct(tenantID.(uuid.UUID), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

// UpdateInventory handles PATCH /api/products/:id/inventory
func (h *Handler) UpdateInventory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productID := c.Param("id")
	
	var req struct {
		Quantity int `json:"quantity" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.UpdateInventory(tenantID.(uuid.UUID), productID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory updated successfully",
	})
}

// Category Handlers

// CreateCategory handles POST /api/categories
func (h *Handler) CreateCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	createdCategory, err := h.service.CreateCategory(tenantID.(uuid.UUID), &category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data":    createdCategory,
	})
}

// GetCategory handles GET /api/categories/:id
func (h *Handler) GetCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categoryID := c.Param("id")
	
	category, err := h.service.GetCategory(tenantID.(uuid.UUID), categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// ListCategories handles GET /api/categories
func (h *Handler) ListCategories(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categories, err := h.service.ListCategories(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// GetProductStats handles GET /api/products/stats
func (h *Handler) GetProductStats(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	stats, err := h.service.GetProductStats(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// BulkUpdateProducts handles PATCH /api/products/bulk
func (h *Handler) BulkUpdateProducts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"message": "Products updated successfully",
		"updated_count": len(req.ProductIDs),
	})
}

// DuplicateProduct handles POST /api/products/:id/duplicate
func (h *Handler) DuplicateProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productID := c.Param("id")
	
	product, err := h.service.DuplicateProduct(tenantID.(uuid.UUID), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product duplicated successfully",
		"data":    product,
	})
}

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
}

// GetLowStockProducts handles GET /api/products/low-stock
func (h *Handler) GetLowStockProducts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

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
			"products": products,
			"threshold": threshold,
			"count":    len(products),
		},
	})
}

// UpdateProductStatus handles PATCH /api/products/:id/status
func (h *Handler) UpdateProductStatus(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	productID := c.Param("id")
	
	var req struct {
		Status ProductStatus `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.UpdateProductStatus(tenantID.(uuid.UUID), productID, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product status updated successfully",
	})
}

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
	if err := c.ShouldBindJSON(&variant); err != nil {
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

	productID := c.Param("id")
	
	variants, err := h.service.GetProductVariants(tenantID.(uuid.UUID), productID)
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
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	variantIDStr := c.Param("variantId")
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant ID"})
		return
	}
	
	var variant ProductVariant
	if err := c.ShouldBindJSON(&variant); err != nil {
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

	productID := c.Param("id")
	variantID := c.Param("variantId")
	
	err := h.service.DeleteProductVariant(tenantID.(uuid.UUID), productID, variantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product variant deleted successfully",
	})
}

// Enhanced Category Handlers

// UpdateCategory handles PUT /api/categories/:id
func (h *Handler) UpdateCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categoryID := c.Param("id")
	
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	updatedCategory, err := h.service.UpdateCategory(tenantID.(uuid.UUID), categoryID, category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"data":    updatedCategory,
	})
}

// DeleteCategory handles DELETE /api/categories/:id
func (h *Handler) DeleteCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categoryID := c.Param("id")
	
	err := h.service.DeleteCategory(tenantID.(uuid.UUID), categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}

// GetRootCategories handles GET /api/categories/root
func (h *Handler) GetRootCategories(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categories, err := h.service.GetRootCategories(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch root categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// GetCategoryChildren handles GET /api/categories/:id/children
func (h *Handler) GetCategoryChildren(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categoryID := c.Param("id")
	
	categories, err := h.service.GetCategoryChildren(tenantID.(uuid.UUID), categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// RegisterRoutes registers all product routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		// Product CRUD
		products.POST("", h.CreateProduct)
		products.GET("", h.ListProducts)
		products.GET("/search", h.SearchProducts)
		products.GET("/stats", h.GetProductStats)
		products.GET("/low-stock", h.GetLowStockProducts)
		products.PATCH("/bulk", h.BulkUpdateProducts)
		
		products.GET("/:id", h.GetProduct)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
		products.PATCH("/:id/status", h.UpdateProductStatus)
		products.PATCH("/:id/inventory", h.UpdateInventory)
		products.POST("/:id/duplicate", h.DuplicateProduct)
		
		// Product variants
		products.POST("/:id/variants", h.CreateProductVariant)
		products.GET("/:id/variants", h.GetProductVariants)
		products.PUT("/:id/variants/:variantId", h.UpdateProductVariant)
		products.DELETE("/:id/variants/:variantId", h.DeleteProductVariant)
		
		// Product by slug (for storefront)
		products.GET("/slug/:slug", h.GetProductBySlug)
	}

	categories := router.Group("/categories")
	{
		// Category CRUD
		categories.POST("", h.CreateCategory)
		categories.GET("", h.ListCategories)
		categories.GET("/root", h.GetRootCategories)
		
		categories.GET("/:id", h.GetCategory)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
		categories.GET("/:id/children", h.GetCategoryChildren)
	}
}

// TODO: Add more handlers
// - ImportProducts(c *gin.Context) - CSV/Excel import
// - ExportProducts(c *gin.Context) - CSV/Excel export
// - UploadProductImages(c *gin.Context) - Image upload
// - DeleteProductImage(c *gin.Context) - Delete specific image
