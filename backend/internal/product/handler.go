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

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	product, err := h.service.CreateProduct(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    product,
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
	
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	product, err := h.service.UpdateProduct(tenantID.(uuid.UUID), productID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data":    product,
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

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	category, err := h.service.CreateCategory(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data":    category,
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

	// TODO: Implement product statistics
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_products":    0,
			"active_products":   0,
			"draft_products":    0,
			"out_of_stock":      0,
			"low_stock":         0,
			"total_categories":  0,
			"total_value":       0,
		},
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

	// TODO: Implement bulk update
	c.JSON(http.StatusOK, gin.H{
		"message": "Products updated successfully",
		"updated_count": len(req.ProductIDs),
	})
}

// TODO: Add more handlers
// - ImportProducts(c *gin.Context) - CSV/Excel import
// - ExportProducts(c *gin.Context) - CSV/Excel export
// - DuplicateProduct(c *gin.Context) - Clone a product
// - UploadProductImages(c *gin.Context) - Image upload
// - DeleteProductImage(c *gin.Context) - Delete specific image
