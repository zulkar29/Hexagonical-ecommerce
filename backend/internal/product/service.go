package product

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
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
	if err := product.ValidateProductData(); err != nil {
		return nil, err
	}

	// Validate struct tags
	if err := s.validator.Struct(product); err != nil {
		return nil, err
	}

	// Generate slug
	slug := s.generateSlug(product.Name)
	
	// Check if slug exists for this tenant
	if exists, err := s.repo.SlugExists(tenantID, slug); err != nil {
		return nil, err
	} else if exists {
		slug = s.generateUniqueSlug(tenantID, slug)
	}
	product.Slug = slug

	// Validate category if provided
	if product.CategoryID != uuid.Nil {
		if exists, err := s.repo.CategoryExists(tenantID, product.CategoryID); err != nil {
			return nil, err
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

	return s.repo.SaveProduct(product)
}

// GetProduct retrieves a product by ID
func (s *Service) GetProduct(tenantID uuid.UUID, id string) (*Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	return s.repo.FindProductByID(tenantID, productID)
}

// GetProductBySlug retrieves a product by slug
func (s *Service) GetProductBySlug(tenantID uuid.UUID, slug string) (*Product, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	return s.repo.FindProductBySlug(tenantID, slug)
}

// UpdateProduct updates an existing product
func (s *Service) UpdateProduct(tenantID, productID uuid.UUID, product *Product) (*Product, error) {
	// Get existing product
	existingProduct, err := s.repo.FindProductByID(tenantID, productID)
	if err != nil {
		return nil, err
	}
	if existingProduct == nil {
		return nil, errors.New("product not found")
	}

	// Validate product data
	if err := product.ValidateProductData(); err != nil {
		return nil, err
	}

	// Validate struct tags
	if err := s.validator.Struct(product); err != nil {
		return nil, err
	}

	// Update fields from provided product
	if product.Name != "" {
		existingProduct.Name = strings.TrimSpace(product.Name)
		// Regenerate slug if name changed
		slug := s.generateSlug(existingProduct.Name)
		if slug != existingProduct.Slug {
			if exists, err := s.repo.SlugExists(tenantID, slug); err != nil {
				return nil, err
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
		if exists, err := s.repo.CategoryExists(tenantID, product.CategoryID); err != nil {
			return nil, err
		} else if !exists {
			return nil, errors.New("category not found")
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

	return s.repo.SaveProduct(existingProduct)
}

// ListProducts returns a paginated list of products
func (s *Service) ListProducts(tenantID uuid.UUID, filter ProductListFilter, offset, limit int) ([]*Product, int64, error) {
	return s.repo.ListProducts(tenantID, filter, offset, limit)
}

// DeleteProduct soft deletes a product
func (s *Service) DeleteProduct(tenantID uuid.UUID, id string) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	return s.repo.DeleteProduct(tenantID, productID)
}

// UpdateInventory updates product inventory
func (s *Service) UpdateInventory(tenantID uuid.UUID, id string, quantity int) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	return s.repo.UpdateInventory(tenantID, productID, quantity)
}

// CreateCategory creates a new category
func (s *Service) CreateCategory(tenantID uuid.UUID, category *Category) (*Category, error) {
	// Set tenant ID and generate new ID
	category.TenantID = tenantID
	category.ID = uuid.New()

	// Validate struct tags
	if err := s.validator.Struct(category); err != nil {
		return nil, err
	}

	// Generate slug
	slug := s.generateSlug(category.Name)
	
	// Check if slug exists for this tenant
	if exists, err := s.repo.CategorySlugExists(tenantID, slug); err != nil {
		return nil, err
	} else if exists {
		slug = s.generateUniqueCategorySlug(tenantID, slug)
	}
	category.Slug = slug

	// Validate parent category if provided
	if category.ParentID != nil {
		if exists, err := s.repo.CategoryExists(tenantID, *category.ParentID); err != nil {
			return nil, err
		} else if !exists {
			return nil, errors.New("parent category not found")
		}
	}

	// Trim string fields
	category.Name = strings.TrimSpace(category.Name)
	category.Description = strings.TrimSpace(category.Description)
	category.MetaTitle = strings.TrimSpace(category.MetaTitle)
	category.MetaDescription = strings.TrimSpace(category.MetaDescription)

	return s.repo.SaveCategory(category)
}

// GetCategory retrieves a category by ID
func (s *Service) GetCategory(tenantID uuid.UUID, id string) (*Category, error) {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	return s.repo.FindCategoryByID(tenantID, categoryID)
}

// ListCategories returns all categories for a tenant
func (s *Service) ListCategories(tenantID uuid.UUID) ([]*Category, error) {
	return s.repo.ListCategories(tenantID)
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
		if exists, err := s.repo.SlugExists(tenantID, newSlug); err == nil && !exists {
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



// TODO: Add more service methods
// - BulkUpdateProducts(tenantID uuid.UUID, productIDs []uuid.UUID, updates map[string]interface{}) error
// - ExportProducts(tenantID uuid.UUID, format string) ([]byte, error)
// - ImportProducts(tenantID uuid.UUID, data []byte) error
// - GetProductStats(tenantID uuid.UUID) (*ProductStats, error)
// - DuplicateProduct(tenantID uuid.UUID, productID uuid.UUID) (*Product, error)

// GetProductStats returns product statistics for a tenant
func (s *Service) GetProductStats(tenantID uuid.UUID) (*ProductStats, error) {
	return s.repo.GetProductStats(tenantID)
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
		return errors.New("no valid product IDs provided")
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
		return errors.New("no valid update fields provided")
	}

	return s.repo.BulkUpdateProducts(tenantID, uuidIDs, validUpdates)
}

// DuplicateProduct creates a copy of an existing product
func (s *Service) DuplicateProduct(tenantID uuid.UUID, productIDStr string) (*Product, error) {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	// Get original product
	original, err := s.repo.FindProductByID(tenantID, productID)
	if err != nil {
		return nil, err
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

	return s.repo.SaveProduct(duplicate)
}

// SearchProducts performs search across products
func (s *Service) SearchProducts(tenantID uuid.UUID, query string, offset, limit int) ([]*Product, int64, error) {
	if query == "" {
		return nil, 0, errors.New("search query is required")
	}

	return s.repo.SearchProducts(tenantID, query, offset, limit)
}

// GetLowStockProducts returns products with low inventory
func (s *Service) GetLowStockProducts(tenantID uuid.UUID, threshold int) ([]*Product, error) {
	if threshold <= 0 {
		threshold = 10 // Default threshold
	}

	return s.repo.GetLowStockProducts(tenantID, threshold)
}

// UpdateProductStatus updates product status
func (s *Service) UpdateProductStatus(tenantID uuid.UUID, productIDStr string, status ProductStatus) error {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return errors.New("invalid product ID")
	}

	// Validate status
	validStatuses := map[ProductStatus]bool{
		ProductStatusDraft:     true,
		ProductStatusActive:    true,
		ProductStatusInactive:  true,
		ProductStatusArchived:  true,
	}

	if !validStatuses[status] {
		return errors.New("invalid product status")
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
	if err := s.validator.Struct(variant); err != nil {
		return nil, err
	}

	// Validate product exists
	if exists, err := s.repo.ProductExists(tenantID, productID); err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New("product not found")
	}

	// Validate business rules
	if variant.ComparePrice > 0 && variant.Price >= variant.ComparePrice {
		return nil, errors.New("compare price must be higher than selling price")
	}

	// Trim string fields
	variant.SKU = strings.TrimSpace(variant.SKU)
	variant.Barcode = strings.TrimSpace(variant.Barcode)

	return s.repo.SaveProductVariant(variant)
}

// GetProductVariants returns all variants for a product
func (s *Service) GetProductVariants(tenantID uuid.UUID, productIDStr string) ([]*ProductVariant, error) {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	// Verify product exists and belongs to tenant
	_, err = s.repo.FindProductByID(tenantID, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return s.repo.FindProductVariants(tenantID, productID)
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
	if err := s.validator.Struct(variant); err != nil {
		return nil, err
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

	return s.repo.SaveProductVariant(existingVariant)
}

// DeleteProductVariant deletes a product variant
func (s *Service) DeleteProductVariant(tenantID uuid.UUID, productIDStr, variantIDStr string) error {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return errors.New("invalid product ID")
	}

	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		return errors.New("invalid variant ID")
	}

	// Verify product exists and belongs to tenant
	_, err = s.repo.FindProductByID(tenantID, productID)
	if err != nil {
		return errors.New("product not found")
	}

	return s.repo.DeleteProductVariant(tenantID, variantID)
}

// Category management methods

// UpdateCategory updates an existing category
func (s *Service) UpdateCategory(tenantID, categoryID uuid.UUID, category *Category) (*Category, error) {
	// Get existing category
	existingCategory, err := s.repo.GetCategory(tenantID, categoryID)
	if err != nil {
		return nil, err
	}
	if existingCategory == nil {
		return nil, errors.New("category not found")
	}

	// Validate struct tags
	if err := s.validator.Struct(category); err != nil {
		return nil, err
	}

	// Update fields from provided category
	if category.Name != "" {
		existingCategory.Name = strings.TrimSpace(category.Name)
		// Regenerate slug if name changed
		slug := s.generateSlug(existingCategory.Name)
		if slug != existingCategory.Slug {
			if exists, err := s.repo.CategorySlugExists(tenantID, slug); err != nil {
				return nil, err
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
		if exists, err := s.repo.CategoryExists(tenantID, *category.ParentID); err != nil {
			return nil, err
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

	return s.repo.SaveCategory(existingCategory)
}

// DeleteCategory deletes a category
func (s *Service) DeleteCategory(tenantID uuid.UUID, categoryIDStr string) error {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return errors.New("invalid category ID")
	}

	// Check if category has products
	products, _, err := s.repo.GetProductsByCategoryID(tenantID, categoryID, 0, 1)
	if err != nil {
		return err
	}

	if len(products) > 0 {
		return errors.New("cannot delete category with products")
	}

	// Check if category has children
	children, err := s.repo.GetCategoryChildren(tenantID, categoryID)
	if err != nil {
		return err
	}

	if len(children) > 0 {
		return errors.New("cannot delete category with subcategories")
	}

	return s.repo.DeleteCategory(tenantID, categoryID)
}

// GetRootCategories returns top-level categories
func (s *Service) GetRootCategories(tenantID uuid.UUID) ([]*Category, error) {
	return s.repo.GetRootCategories(tenantID)
}

// GetCategoryChildren returns child categories
func (s *Service) GetCategoryChildren(tenantID uuid.UUID, categoryIDStr string) ([]*Category, error) {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	return s.repo.GetCategoryChildren(tenantID, categoryID)
}
