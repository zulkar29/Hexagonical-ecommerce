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

	createdProduct, err := h.service.CreateProduct(tenantID.(uuid.UUID), &product)
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

	productIDStr := c.Param("id")
	// Note: Service already handles string to UUID conversion
	product, err := h.service.GetProduct(tenantID.(uuid.UUID), productIDStr)
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
		if err := c.ShouldBindJSON(&req); err != nil {
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
		if err := c.ShouldBindJSON(&req); err != nil {
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
		product, err := h.service.DuplicateProduct(tenantID.(uuid.UUID), productID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	if err := c.ShouldBindJSON(&product); err != nil {
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
				"products": products,
				"threshold": threshold,
				"count":    len(products),
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

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	
	category, err := h.service.GetCategory(tenantID.(uuid.UUID), categoryID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// ListCategories handles GET /api/categories with hierarchy support
func (h *Handler) ListCategories(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	// Check for hierarchy queries
	hierarchy := c.Query("hierarchy")
	switch hierarchy {
	case "root":
		categories, err := h.service.GetRootCategories(tenantID.(uuid.UUID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Root categories retrieved successfully",
			"data":    categories,
		})
		return
	case "children":
		parentID := c.Query("parent_id")
		if parentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent_id is required for children hierarchy"})
			return
		}
		categories, err := h.service.GetCategoryChildren(tenantID.(uuid.UUID), parentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Child categories retrieved successfully",
			"data":    categories,
		})
		return
	}

	// Regular category listing
	categories, err := h.service.ListCategories(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categories retrieved successfully",
		"data":    categories,
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

	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
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
	
	err = h.service.DeleteProductVariant(tenantID.(uuid.UUID), productID.String(), variantID.String())
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

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	updatedCategory, err := h.service.UpdateCategory(tenantID.(uuid.UUID), categoryID, &category)
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

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	
	err = h.service.DeleteCategory(tenantID.(uuid.UUID), categoryID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
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

	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID is required"})
		return
	}

	children, err := h.service.GetCategoryChildren(tenantID.(uuid.UUID), categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": children,
	})
}



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
			ProductIDs []string `json:"product_ids" binding:"required"`
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
		products.GET("", h.ListProducts) // Supports ?type=stats|low-stock|search, ?q=query
		products.POST("/operations", h.HandleProductOperations) // Supports ?operation=bulk_update|import|export
		products.GET("/:id", h.GetProduct)
		products.PUT("/:id", h.UpdateProduct) // Supports ?action=update_inventory|update_status|duplicate
		products.DELETE("/:id", h.DeleteProduct)
		products.GET("/slug/:slug", h.GetProductBySlug)

		// Product variant routes
		products.POST("/:id/variants", h.CreateProductVariant)
		products.GET("/:id/variants", h.GetProductVariants)
		products.PUT("/:id/variants/:variant_id", h.UpdateProductVariant)
		products.DELETE("/:id/variants/:variant_id", h.DeleteProductVariant)
	}

	// Category routes
	categories := router.Group("/categories")
	{
		categories.POST("", h.CreateCategory)
		categories.GET("", h.ListCategories) // Supports ?hierarchy=root|children, ?parent_id=id
		categories.GET("/:id", h.GetCategory)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}

	// TODO: Add routes for:
	// - Product image uploads
	// - Category image uploads
	// - Advanced search and filtering
}

// TODO: Add more handlers
// - ImportProducts(c *gin.Context) - CSV/Excel import
// - ExportProducts(c *gin.Context) - CSV/Excel export
// - UploadProductImages(c *gin.Context) - Image upload
// - DeleteProductImage(c *gin.Context) - Delete specific image
