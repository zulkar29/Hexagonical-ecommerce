package product

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"

	sharedErrors "ecommerce-saas/internal/shared/errors"
)

type ProductListFilter struct {
	Status     ProductStatus `json:"status,omitempty"`
	Type       ProductType   `json:"type,omitempty"`
	CategoryID *uuid.UUID    `json:"category_id,omitempty"`
	MinPrice   *float64      `json:"min_price,omitempty"`
	MaxPrice   *float64      `json:"max_price,omitempty"`
	InStock    *bool         `json:"in_stock,omitempty"`
	Search     string        `json:"search,omitempty"`
}

// Service handles product business logic
type Service struct {
	repo      Repository
	validator *validator.Validate
}

// NewService creates a new product service
func NewService(repo Repository) *Service {
	return &Service{
		repo:      repo,
		validator: validator.New(),
	}
}

// CreateProduct creates a new product
func (s *Service) CreateProduct(tenantID uuid.UUID, product *Product) (*Product, error) {
	// Set tenant ID and generate new ID
	product.TenantID = tenantID
	product.ID = uuid.New()

	// Validate product data
	if validateErr := product.ValidateProductData(); validateErr != nil {
		return nil, validateErr
	}

	// Validate struct tags
	if structErr := s.validator.Struct(product); structErr != nil {
		return nil, sharedErrors.NewValidationError("Product validation failed", structErr.Error())
	}

	// Generate slug
	slug := s.generateSlug(product.Name)

	// Check if slug exists for this tenant
	if exists, slugErr := s.repo.ProductSlugExists(tenantID, slug); slugErr != nil {
		return nil, slugErr
	} else if exists {
		return nil, sharedErrors.ErrProductSlugTaken
	}
	product.Slug = slug

	// Validate category if provided
	if product.CategoryID != uuid.Nil {
		if exists, catErr := s.repo.CategoryExists(tenantID, product.CategoryID); catErr != nil {
			return nil, catErr
		} else if !exists {
			return nil, errors.New("category not found")
		}
	}

	// Trim string fields
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)
	product.SKU = strings.TrimSpace(product.SKU)
	product.Barcode = strings.TrimSpace(product.Barcode)
	product.MetaTitle = strings.TrimSpace(product.MetaTitle)
	product.MetaDescription = strings.TrimSpace(product.MetaDescription)
	product.MetaKeywords = strings.TrimSpace(product.MetaKeywords)

	// Set default status if not provided
	if product.Status == "" {
		product.Status = StatusDraft
	}

	// Set featured image
	if len(product.Images) > 0 {
		product.FeaturedImage = product.Images[0]
	}

	return s.repo.CreateProduct(product)
}

// GetProduct retrieves a product by ID
func (s *Service) GetProduct(tenantID uuid.UUID, id string) (*Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, sharedErrors.NewBadRequestError("Invalid product ID format")
	}

	product, err := s.repo.GetProductByID(tenantID, productID)
	if err != nil {
		return nil, sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to retrieve product", 500)
	}
	if product == nil {
		return nil, sharedErrors.ErrProductNotFound
	}
	return product, nil
}

// GetProductBySlug retrieves a product by slug
func (s *Service) GetProductBySlug(tenantID uuid.UUID, slug string) (*Product, error) {
	if slug == "" {
		return nil, sharedErrors.NewBadRequestError("Product slug is required")
	}

	product, err := s.repo.GetProductBySlug(tenantID, slug)
	if err != nil {
		return nil, sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to retrieve product by slug", 500)
	}
	if product == nil {
		return nil, sharedErrors.ErrProductNotFound
	}
	return product, nil
}

// UpdateProduct updates an existing product
func (s *Service) UpdateProduct(tenantID, productID uuid.UUID, product *Product) (*Product, error) {
	// Get existing product
	existingProduct, err := s.repo.GetProductByID(tenantID, productID)
	if err != nil {
		return nil, sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to find product", 500)
	}
	if existingProduct == nil {
		return nil, sharedErrors.ErrProductNotFound
	}

	// Validate product data
	if validateErr := product.ValidateProductData(); validateErr != nil {
		return nil, validateErr
	}

	// Validate struct tags
	if structErr := s.validator.Struct(product); structErr != nil {
		return nil, structErr
	}

	// Update fields from provided product
	if product.Name != "" {
		existingProduct.Name = strings.TrimSpace(product.Name)
		// Regenerate slug if name changed
		slug := s.generateSlug(existingProduct.Name)
		if slug != existingProduct.Slug {
			if exists, slugErr := s.repo.ProductSlugExists(tenantID, slug); slugErr != nil {
				return nil, sharedErrors.Wrap(slugErr, sharedErrors.CodeInternal, "Failed to check slug existence", 500)
			} else if exists {
				slug = s.generateUniqueSlug(tenantID, slug)
			}
			existingProduct.Slug = slug
		}
	}

	if product.Description != "" {
		existingProduct.Description = strings.TrimSpace(product.Description)
	}
	if product.Type != "" {
		existingProduct.Type = product.Type
	}
	if product.Status != "" {
		existingProduct.Status = product.Status
	}
	if product.Price > 0 {
		existingProduct.Price = product.Price
	}
	if product.ComparePrice > 0 {
		existingProduct.ComparePrice = product.ComparePrice
	}
	if product.CostPrice > 0 {
		existingProduct.CostPrice = product.CostPrice
	}
	if product.SKU != "" {
		existingProduct.SKU = strings.TrimSpace(product.SKU)
	}
	if product.Barcode != "" {
		existingProduct.Barcode = strings.TrimSpace(product.Barcode)
	}
	if product.InventoryQuantity >= 0 {
		existingProduct.InventoryQuantity = product.InventoryQuantity
	}
	existingProduct.TrackQuantity = product.TrackQuantity
	existingProduct.AllowBackorder = product.AllowBackorder
	if product.Weight > 0 {
		existingProduct.Weight = product.Weight
	}
	if product.Length > 0 {
		existingProduct.Length = product.Length
	}
	if product.Width > 0 {
		existingProduct.Width = product.Width
	}
	if product.Height > 0 {
		existingProduct.Height = product.Height
	}
	if product.MetaTitle != "" {
		existingProduct.MetaTitle = strings.TrimSpace(product.MetaTitle)
	}
	if product.MetaDescription != "" {
		existingProduct.MetaDescription = strings.TrimSpace(product.MetaDescription)
	}
	if product.MetaKeywords != "" {
		existingProduct.MetaKeywords = strings.TrimSpace(product.MetaKeywords)
	}
	if product.CategoryID != uuid.Nil {
		// Validate category if provided
		if exists, catErr := s.repo.CategoryExists(tenantID, product.CategoryID); catErr != nil {
			return nil, sharedErrors.Wrap(catErr, sharedErrors.CodeInternal, "Failed to validate category", 500)
		} else if !exists {
			return nil, sharedErrors.ErrCategoryNotFound
		}
		existingProduct.CategoryID = product.CategoryID
	}
	if product.Tags != nil {
		existingProduct.Tags = product.Tags
	}
	if product.Images != nil {
		existingProduct.Images = product.Images
		// Update featured image
		if len(product.Images) > 0 {
			existingProduct.FeaturedImage = product.Images[0]
		} else {
			existingProduct.FeaturedImage = ""
		}
	}

	existingProduct.UpdatedAt = time.Now()

	return s.repo.UpdateProduct(existingProduct)
}

// ListProducts returns a paginated list of products
func (s *Service) ListProducts(tenantID uuid.UUID, filter ProductListFilter, offset, limit int) ([]*Product, int64, error) {
	return s.repo.ListProducts(tenantID, filter, offset, limit)
}

// DeleteProduct soft deletes a product
func (s *Service) DeleteProduct(tenantID uuid.UUID, id string) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return sharedErrors.NewBadRequestError("Invalid product ID format")
	}

	if err := s.repo.DeleteProduct(tenantID, productID); err != nil {
		return sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to delete product", 500)
	}
	return nil
}

// UpdateInventory updates product inventory
func (s *Service) UpdateInventory(tenantID uuid.UUID, id string, quantity int) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return sharedErrors.NewBadRequestError("Invalid product ID format")
	}

	if err := s.repo.UpdateProductInventory(tenantID, productID, quantity); err != nil {
		return sharedErrors.Wrap(err, sharedErrors.CodeInternal, "Failed to update inventory", 500)
	}
	return nil
}

// ReserveStock reserves inventory for a product (decrements inventory)
func (s *Service) ReserveStock(tenantID uuid.UUID, productID uuid.UUID, quantity int) error {
	// Get the product first
	product, err := s.repo.GetProductByID(tenantID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedErrors.NewNotFoundError("product not found")
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Check if we can decrement inventory
	if !product.CanDecrementInventory(quantity) {
		return fmt.Errorf("insufficient inventory for product %s. Available: %d, Requested: %d", product.Name, product.InventoryQuantity, quantity)
	}

	// Decrement inventory
	if decrementErr := product.DecrementInventory(quantity); decrementErr != nil {
		return decrementErr
	}

	// Save the updated product
	_, err = s.repo.UpdateProduct(product)
	return err
}

// RestoreStock restores inventory for a product (increments inventory)
func (s *Service) RestoreStock(tenantID uuid.UUID, productID uuid.UUID, quantity int) error {
	// Get the product first
	product, err := s.repo.GetProductByID(tenantID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedErrors.NewNotFoundError("product not found")
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Increment inventory
	product.IncrementInventory(quantity)

	// Save the updated product
	_, err = s.repo.UpdateProduct(product)
	return err
}

// CheckAvailability checks if sufficient inventory is available
func (s *Service) CheckAvailability(tenantID uuid.UUID, productID uuid.UUID, variantID *uuid.UUID, quantity int) (bool, error) {
	// For now, we only handle product-level inventory (not variant-level)
	product, err := s.repo.GetProductByID(tenantID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, sharedErrors.NewNotFoundError("product not found")
		}
		return false, fmt.Errorf("failed to get product: %w", err)
	}

	return product.CanDecrementInventory(quantity), nil
}

// CreateCategory creates a new category
func (s *Service) CreateCategory(tenantID uuid.UUID, category *Category) (*Category, error) {
	// Set tenant ID and generate new ID
	category.TenantID = tenantID
	category.ID = uuid.New()

	// Validate struct tags
	if structErr := s.validator.Struct(category); structErr != nil {
		return nil, sharedErrors.NewValidationError("Category validation failed", structErr.Error())
	}

	// Generate slug
	slug := s.generateSlug(category.Name)

	// Check if slug exists for this tenant
	if exists, slugErr := s.repo.CategorySlugExists(tenantID, slug); slugErr != nil {
		return nil, sharedErrors.Wrap(slugErr, sharedErrors.CodeInternal, "Failed to check category slug existence", 500)
	} else if exists {
		return nil, sharedErrors.NewConflictError("Category slug already exists")
	}
	category.Slug = slug

	// Validate parent category if provided
	if category.ParentID != nil {
		if exists, parentErr := s.repo.CategoryExists(tenantID, *category.ParentID); parentErr != nil {
			return nil, sharedErrors.Wrap(parentErr, sharedErrors.CodeInternal, "Failed to validate parent category", 500)
		} else if !exists {
			return nil, sharedErrors.ErrCategoryNotFound
		}
	}

	// Trim string fields
	category.Name = strings.TrimSpace(category.Name)
	category.Description = strings.TrimSpace(category.Description)
	category.MetaTitle = strings.TrimSpace(category.MetaTitle)
	category.MetaDescription = strings.TrimSpace(category.MetaDescription)

	return s.repo.CreateCategory(category)
}

// GetCategory retrieves a category by ID
func (s *Service) GetCategory(tenantID uuid.UUID, id string) (*Category, error) {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, sharedErrors.NewBadRequestError("Invalid category ID format")
	}

	category, err := s.repo.GetCategoryByID(tenantID, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return category, nil
}

// ListCategories returns all categories for a tenant
func (s *Service) ListCategories(tenantID uuid.UUID) ([]*Category, error) {
	categories, err := s.repo.ListCategories(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	return categories, nil
}

// Private helper methods

func (s *Service) generateSlug(name string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters using regex-like approach
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	slug = result.String()

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Ensure slug is not empty
	if slug == "" {
		slug = "product"
	}

	return slug
}

func (s *Service) generateUniqueSlug(tenantID uuid.UUID, baseSlug string) string {
	// Try appending numbers 1, 2, 3, etc.
	for i := 1; i <= 100; i++ {
		newSlug := fmt.Sprintf("%s-%d", baseSlug, i)
		if exists, err := s.repo.ProductSlugExists(tenantID, newSlug); err == nil && !exists {
			return newSlug
		}
	}

	// Fallback to UUID if all numbers are taken
	return baseSlug + "-" + uuid.New().String()[:8]
}

func (s *Service) generateUniqueCategorySlug(tenantID uuid.UUID, baseSlug string) string {
	// Try appending numbers 1, 2, 3, etc.
	for i := 1; i <= 100; i++ {
		newSlug := fmt.Sprintf("%s-%d", baseSlug, i)
		if exists, err := s.repo.CategorySlugExists(tenantID, newSlug); err == nil && !exists {
			return newSlug
		}
	}

	// Fallback to UUID if all numbers are taken
	return baseSlug + "-" + uuid.New().String()[:8]
}

// BulkDeleteProducts deletes multiple products at once
func (s *Service) BulkDeleteProducts(tenantID uuid.UUID, productIDs []uuid.UUID) error {
	if len(productIDs) == 0 {
		return sharedErrors.NewValidationError("No product IDs provided", "")
	}

	if err := s.repo.BulkDeleteProducts(tenantID, productIDs); err != nil {
		return fmt.Errorf("failed to bulk delete products: %w", err)
	}
	return nil
}

// GetProductAnalytics returns analytics data for a specific product
func (s *Service) GetProductAnalytics(tenantID uuid.UUID, productID, analyticsType string) (interface{}, error) {
	prodID, err := uuid.Parse(productID)
	if err != nil {
		return nil, sharedErrors.NewBadRequestError("Invalid product ID format")
	}

	// Verify product exists
	_, err = s.repo.GetProductByID(tenantID, prodID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	switch analyticsType {
	case "performance":
		// Get product performance analytics
		product, err := s.repo.GetProductByID(tenantID, prodID)
		if err != nil {
			return nil, err
		}

		// Calculate basic performance metrics
		conversionRate := 0.0
		if product.InventoryQuantity > 0 {
			// Simple conversion rate calculation based on stock movement
			conversionRate = float64(product.InventoryQuantity) / 100.0
		}

		return map[string]interface{}{
			"views":           150,                         // Placeholder - would come from analytics service
			"sales":           25,                          // Placeholder - would come from order service
			"revenue":         float64(product.Price) * 25, // Calculated from price and sales
			"conversion_rate": conversionRate,
		}, nil
	case "inventory":
		// Get inventory analytics
		product, err := s.repo.GetProductByID(tenantID, prodID)
		if err != nil {
			return nil, err
		}

		lowStockAlert := 0
		if product.TrackQuantity && product.InventoryQuantity < 10 {
			lowStockAlert = 1
		}

		return map[string]interface{}{
			"current_stock": product.InventoryQuantity,
			"stock_movements": []interface{}{
				map[string]interface{}{
					"date":     time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
					"type":     "sale",
					"quantity": -2,
					"reason":   "Product sold",
				},
				map[string]interface{}{
					"date":     time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
					"type":     "restock",
					"quantity": 50,
					"reason":   "Inventory replenishment",
				},
			},
			"low_stock_alerts": lowStockAlert,
		}, nil
	case "sales":
		// Get sales analytics
		product, err := s.repo.GetProductByID(tenantID, prodID)
		if err != nil {
			return nil, err
		}

		// Generate mock monthly sales data
		monthlySales := []interface{}{}
		for i := 5; i >= 0; i-- {
			month := time.Now().AddDate(0, -i, 0)
			sales := 10 + (i * 3) // Mock increasing sales
			monthlySales = append(monthlySales, map[string]interface{}{
				"month":   month.Format("2006-01"),
				"sales":   sales,
				"revenue": float64(product.Price) * float64(sales),
			})
		}

		// Generate top variants data
		topVariants := []interface{}{}
		for i, variant := range product.Variants {
			if i < 3 { // Top 3 variants
				topVariants = append(topVariants, map[string]interface{}{
					"variant_id": variant.ID,
					"name":       variant.GetDisplayName(),
					"sales":      15 - (i * 3), // Mock decreasing sales
					"revenue":    float64(variant.Price) * float64(15-(i*3)),
				})
			}
		}

		return map[string]interface{}{
			"total_sales":   75, // Mock total sales
			"monthly_sales": monthlySales,
			"top_variants":  topVariants,
		}, nil
	default:
		return nil, errors.New("invalid analytics type")
	}
}

// Additional service methods for import/export functionality
// These would be implemented when file processing capabilities are added:
// - ExportProducts(tenantID uuid.UUID, format string) ([]byte, error)
// - ImportProducts(tenantID uuid.UUID, data []byte) error

// GetProductStats returns product statistics for a tenant
func (s *Service) GetProductStats(tenantID uuid.UUID) (*ProductStats, error) {
	stats, err := s.repo.GetProductStats(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product stats: %w", err)
	}
	return stats, nil
}

// BulkUpdateProducts updates multiple products at once
func (s *Service) BulkUpdateProducts(tenantID uuid.UUID, productIDs []string, updates map[string]interface{}) error {
	// Parse product IDs
	uuidIDs := make([]uuid.UUID, 0, len(productIDs))
	for _, idStr := range productIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			uuidIDs = append(uuidIDs, id)
		}
	}

	if len(uuidIDs) == 0 {
		return sharedErrors.NewValidationError("No valid product IDs provided", "")
	}

	// Validate updates - only allow specific fields
	allowedFields := map[string]bool{
		"status":        true,
		"price":         true,
		"compare_price": true,
		"category_id":   true,
		"tags":          true,
	}

	validUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			validUpdates[key] = value
		}
	}

	if len(validUpdates) == 0 {
		return sharedErrors.NewValidationError("No valid update fields provided", "")
	}

	if err := s.repo.BulkUpdateProducts(tenantID, uuidIDs, validUpdates); err != nil {
		return fmt.Errorf("failed to bulk update products: %w", err)
	}
	return nil
}

// DuplicateProduct creates a copy of an existing product
func (s *Service) DuplicateProduct(tenantID uuid.UUID, productIDStr string) (*Product, error) {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, sharedErrors.NewBadRequestError("Invalid product ID format")
	}

	// Get original product
	original, err := s.repo.GetProductByID(tenantID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Create duplicate with modified name and slug
	duplicate := &Product{
		ID:                uuid.New(),
		TenantID:          tenantID,
		Name:              original.Name + " (Copy)",
		Slug:              s.generateUniqueSlug(tenantID, original.Slug+"-copy"),
		Description:       original.Description,
		Type:              original.Type,
		Status:            ProductStatusDraft, // Always create as draft
		Price:             original.Price,
		ComparePrice:      original.ComparePrice,
		CostPrice:         original.CostPrice,
		SKU:               "", // Clear SKU to avoid conflicts
		Barcode:           "", // Clear barcode to avoid conflicts
		InventoryQuantity: 0,  // Start with zero inventory
		TrackQuantity:     original.TrackQuantity,
		AllowBackorder:    original.AllowBackorder,
		Weight:            original.Weight,
		Length:            original.Length,
		Width:             original.Width,
		Height:            original.Height,
		MetaTitle:         original.MetaTitle,
		MetaDescription:   original.MetaDescription,
		MetaKeywords:      original.MetaKeywords,
		FeaturedImage:     original.FeaturedImage,
		Images:            original.Images,
		CategoryID:        original.CategoryID,
		Tags:              original.Tags,
	}

	return s.repo.CreateProduct(duplicate)
}

// SearchProducts performs search across products
func (s *Service) SearchProducts(tenantID uuid.UUID, query string, offset, limit int) ([]*Product, int64, error) {
	if query == "" {
		return nil, 0, sharedErrors.NewBadRequestError("Search query is required")
	}

	products, total, err := s.repo.SearchProducts(tenantID, query, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}
	return products, total, nil
}

// GetLowStockProducts returns products with low inventory
func (s *Service) GetLowStockProducts(tenantID uuid.UUID, threshold int) ([]*Product, error) {
	if threshold <= 0 {
		threshold = 10 // Default threshold
	}

	products, err := s.repo.GetLowStockProducts(tenantID, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}
	return products, nil
}

// UpdateProductStatus updates product status
func (s *Service) UpdateProductStatus(tenantID uuid.UUID, productIDStr string, status ProductStatus) error {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return sharedErrors.NewBadRequestError("Invalid product ID format")
	}

	// Validate status
	validStatuses := map[ProductStatus]bool{
		ProductStatusDraft:    true,
		ProductStatusActive:   true,
		ProductStatusInactive: true,
		ProductStatusArchived: true,
	}

	if !validStatuses[status] {
		return sharedErrors.NewValidationError("Invalid product status", "")
	}

	updates := map[string]interface{}{
		"status": status,
	}

	return s.repo.BulkUpdateProducts(tenantID, []uuid.UUID{productID}, updates)
}

// Product Variant methods

// CreateProductVariant creates a new product variant
func (s *Service) CreateProductVariant(tenantID, productID uuid.UUID, variant *ProductVariant) (*ProductVariant, error) {
	// Set product ID and generate new ID
	variant.ProductID = productID
	variant.ID = uuid.New()

	// Validate struct tags
	if structErr := s.validator.Struct(variant); structErr != nil {
		return nil, sharedErrors.NewValidationError("Variant validation failed", structErr.Error())
	}

	// Validate product exists
	if exists, productErr := s.repo.ProductExists(tenantID, productID); productErr != nil {
		return nil, sharedErrors.Wrap(productErr, sharedErrors.CodeInternal, "Failed to check product existence", 500)
	} else if !exists {
		return nil, sharedErrors.ErrProductNotFound
	}

	// Validate business rules
	if variant.ComparePrice > 0 && variant.Price >= variant.ComparePrice {
		return nil, sharedErrors.NewValidationError("Compare price must be higher than selling price", "")
	}

	// Trim string fields
	variant.SKU = strings.TrimSpace(variant.SKU)
	variant.Barcode = strings.TrimSpace(variant.Barcode)

	return s.repo.CreateProductVariant(variant)
}

// GetProductVariants returns all variants for a product
func (s *Service) GetProductVariants(tenantID uuid.UUID, productIDStr string) ([]*ProductVariant, error) {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, sharedErrors.NewBadRequestError("Invalid product ID format")
	}

	// Verify product exists and belongs to tenant
	_, err = s.repo.GetProductByID(tenantID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	variants, err := s.repo.GetProductVariants(tenantID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}
	return variants, nil
}

// UpdateProductVariant updates an existing product variant
func (s *Service) UpdateProductVariant(tenantID, productID, variantID uuid.UUID, variant *ProductVariant) (*ProductVariant, error) {
	// Get existing variant
	existingVariant, err := s.repo.GetProductVariant(tenantID, variantID)
	if err != nil {
		return nil, err
	}
	if existingVariant == nil {
		return nil, errors.New("variant not found")
	}

	// Validate struct tags
	if structErr := s.validator.Struct(variant); structErr != nil {
		return nil, structErr
	}

	// Update fields from provided variant
	if variant.SKU != "" {
		existingVariant.SKU = strings.TrimSpace(variant.SKU)
	}
	if variant.Barcode != "" {
		existingVariant.Barcode = strings.TrimSpace(variant.Barcode)
	}
	if variant.Price > 0 {
		existingVariant.Price = variant.Price
	}
	if variant.ComparePrice > 0 {
		existingVariant.ComparePrice = variant.ComparePrice
	}
	if variant.CostPrice > 0 {
		existingVariant.CostPrice = variant.CostPrice
	}
	if variant.InventoryQuantity >= 0 {
		existingVariant.InventoryQuantity = variant.InventoryQuantity
	}
	existingVariant.TrackQuantity = variant.TrackQuantity
	existingVariant.AllowBackorder = variant.AllowBackorder
	if variant.Weight > 0 {
		existingVariant.Weight = variant.Weight
	}
	if variant.Length > 0 {
		existingVariant.Length = variant.Length
	}
	if variant.Width > 0 {
		existingVariant.Width = variant.Width
	}
	if variant.Height > 0 {
		existingVariant.Height = variant.Height
	}
	if variant.Image != "" {
		existingVariant.Image = variant.Image
	}
	if variant.Options != nil {
		existingVariant.Options = variant.Options
	}
	existingVariant.IsDefault = variant.IsDefault

	// Validate business rules
	if existingVariant.ComparePrice > 0 && existingVariant.Price >= existingVariant.ComparePrice {
		return nil, errors.New("compare price must be higher than selling price")
	}

	existingVariant.UpdatedAt = time.Now()

	return s.repo.UpdateProductVariant(existingVariant)
}

// DeleteProductVariant deletes a product variant
func (s *Service) DeleteProductVariant(tenantID uuid.UUID, productIDStr, variantIDStr string) error {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return sharedErrors.NewBadRequestError("Invalid product ID format")
	}

	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		return sharedErrors.NewBadRequestError("Invalid variant ID format")
	}

	// Verify product exists and belongs to tenant
	_, err = s.repo.GetProductByID(tenantID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedErrors.NewNotFoundError("product not found")
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	if err := s.repo.DeleteProductVariant(tenantID, variantID); err != nil {
		return fmt.Errorf("failed to delete product variant: %w", err)
	}
	return nil
}

// Category management methods

// UpdateCategory updates an existing category
func (s *Service) UpdateCategory(tenantID, categoryID uuid.UUID, category *Category) (*Category, error) {
	// Get existing category
	existingCategory, err := s.repo.GetCategoryByID(tenantID, categoryID)
	if err != nil {
		return nil, err
	}
	if existingCategory == nil {
		return nil, errors.New("category not found")
	}

	// Validate struct tags
	if structErr := s.validator.Struct(category); structErr != nil {
		return nil, structErr
	}

	// Update fields from provided category
	if category.Name != "" {
		existingCategory.Name = strings.TrimSpace(category.Name)
		// Regenerate slug if name changed
		slug := s.generateSlug(existingCategory.Name)
		if slug != existingCategory.Slug {
			if exists, slugErr := s.repo.CategorySlugExists(tenantID, slug); slugErr != nil {
				return nil, slugErr
			} else if exists {
				slug = s.generateUniqueCategorySlug(tenantID, slug)
			}
			existingCategory.Slug = slug
		}
	}

	if category.Description != "" {
		existingCategory.Description = strings.TrimSpace(category.Description)
	}
	if category.Image != "" {
		existingCategory.Image = category.Image
	}
	if category.ParentID != nil && *category.ParentID != uuid.Nil {
		// Validate parent category if provided
		if exists, parentErr := s.repo.CategoryExists(tenantID, *category.ParentID); parentErr != nil {
			return nil, parentErr
		} else if !exists {
			return nil, errors.New("parent category not found")
		}
		existingCategory.ParentID = category.ParentID
	}
	if category.SortOrder > 0 {
		existingCategory.SortOrder = category.SortOrder
	}
	existingCategory.IsActive = category.IsActive
	if category.MetaTitle != "" {
		existingCategory.MetaTitle = strings.TrimSpace(category.MetaTitle)
	}
	if category.MetaDescription != "" {
		existingCategory.MetaDescription = strings.TrimSpace(category.MetaDescription)
	}

	existingCategory.UpdatedAt = time.Now()

	return s.repo.UpdateCategory(existingCategory)
}

// DeleteCategory deletes a category
func (s *Service) DeleteCategory(tenantID uuid.UUID, categoryIDStr string) error {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return sharedErrors.NewBadRequestError("Invalid category ID format")
	}

	// Check if category has products
	products, _, err := s.repo.GetProductsByCategoryID(tenantID, categoryID, 0, 1)
	if err != nil {
		return fmt.Errorf("failed to check category products: %w", err)
	}

	if len(products) > 0 {
		return sharedErrors.NewConflictError("Cannot delete category with products")
	}

	// Check if category has children
	children, err := s.repo.GetCategoryChildren(tenantID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to check category children: %w", err)
	}

	if len(children) > 0 {
		return sharedErrors.NewConflictError("Cannot delete category with subcategories")
	}

	if err := s.repo.DeleteCategory(tenantID, categoryID); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

// GetRootCategories returns top-level categories
func (s *Service) GetRootCategories(tenantID uuid.UUID) ([]*Category, error) {
	categories, err := s.repo.GetRootCategories(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get root categories: %w", err)
	}
	return categories, nil
}

// GetCategoryChildren returns child categories
func (s *Service) GetCategoryChildren(tenantID uuid.UUID, categoryIDStr string) ([]*Category, error) {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return nil, sharedErrors.NewBadRequestError("Invalid category ID format")
	}

	children, err := s.repo.GetCategoryChildren(tenantID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category children: %w", err)
	}
	return children, nil
}
