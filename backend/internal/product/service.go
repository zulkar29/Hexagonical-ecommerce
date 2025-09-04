package product

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

// Request/Response Structures
type CreateProductRequest struct {
	Name        string      `json:"name" validate:"required,min=2,max=255"`
	Description string      `json:"description,omitempty" validate:"max=5000"`
	Type        ProductType `json:"type" validate:"required"`
	Price       float64     `json:"price" validate:"required,min=0"`
	ComparePrice float64    `json:"compare_price,omitempty" validate:"min=0"`
	CostPrice   float64     `json:"cost_price,omitempty" validate:"min=0"`
	
	SKU               string  `json:"sku,omitempty" validate:"max=100"`
	Barcode           string  `json:"barcode,omitempty" validate:"max=100"`
	InventoryQuantity int     `json:"inventory_quantity" validate:"min=0"`
	TrackQuantity     bool    `json:"track_quantity"`
	AllowBackorder    bool    `json:"allow_backorder"`
	
	Weight float64 `json:"weight,omitempty" validate:"min=0"`
	Length float64 `json:"length,omitempty" validate:"min=0"`
	Width  float64 `json:"width,omitempty" validate:"min=0"`
	Height float64 `json:"height,omitempty" validate:"min=0"`
	
	MetaTitle       string `json:"meta_title,omitempty" validate:"max=255"`
	MetaDescription string `json:"meta_description,omitempty" validate:"max=500"`
	MetaKeywords    string `json:"meta_keywords,omitempty" validate:"max=255"`
	
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	Images     []string   `json:"images,omitempty"`
}

type UpdateProductRequest struct {
	Name        string      `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Description string      `json:"description,omitempty" validate:"max=5000"`
	Type        ProductType `json:"type,omitempty"`
	Price       float64     `json:"price,omitempty" validate:"omitempty,min=0"`
	ComparePrice float64    `json:"compare_price,omitempty" validate:"min=0"`
	CostPrice   float64     `json:"cost_price,omitempty" validate:"min=0"`
	
	SKU               string `json:"sku,omitempty" validate:"max=100"`
	Barcode           string `json:"barcode,omitempty" validate:"max=100"`
	InventoryQuantity *int   `json:"inventory_quantity,omitempty" validate:"omitempty,min=0"`
	TrackQuantity     *bool  `json:"track_quantity,omitempty"`
	AllowBackorder    *bool  `json:"allow_backorder,omitempty"`
	
	Weight *float64 `json:"weight,omitempty" validate:"omitempty,min=0"`
	Length *float64 `json:"length,omitempty" validate:"omitempty,min=0"`
	Width  *float64 `json:"width,omitempty" validate:"omitempty,min=0"`
	Height *float64 `json:"height,omitempty" validate:"omitempty,min=0"`
	
	MetaTitle       string `json:"meta_title,omitempty" validate:"max=255"`
	MetaDescription string `json:"meta_description,omitempty" validate:"max=500"`
	MetaKeywords    string `json:"meta_keywords,omitempty" validate:"max=255"`
	
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	Images     []string   `json:"images,omitempty"`
	Status     ProductStatus `json:"status,omitempty"`
}

type CreateCategoryRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=100"`
	Description string     `json:"description,omitempty" validate:"max=500"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder   int        `json:"sort_order"`
	MetaTitle   string     `json:"meta_title,omitempty" validate:"max=255"`
	MetaDescription string `json:"meta_description,omitempty" validate:"max=500"`
}

type UpdateCategoryRequest struct {
	Name            string     `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description     string     `json:"description,omitempty" validate:"max=500"`
	ParentID        *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder       int        `json:"sort_order,omitempty"`
	IsActive        *bool      `json:"is_active,omitempty"`
	MetaTitle       string     `json:"meta_title,omitempty" validate:"max=255"`
	MetaDescription string     `json:"meta_description,omitempty" validate:"max=500"`
}

type CreateVariantRequest struct {
	Name              string            `json:"name" validate:"required,min=1,max=255"`
	SKU               string            `json:"sku,omitempty" validate:"max=100"`
	Price             float64           `json:"price" validate:"min=0"`
	InventoryQuantity int               `json:"inventory_quantity" validate:"min=0"`
	AllowBackorder    bool              `json:"allow_backorder"`
	Options           map[string]string `json:"options,omitempty"`
	Images            []string          `json:"images,omitempty"`
}

type UpdateVariantRequest struct {
	Name              string            `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	SKU               string            `json:"sku,omitempty" validate:"max=100"`
	Price             float64           `json:"price,omitempty" validate:"min=0"`
	InventoryQuantity *int              `json:"inventory_quantity,omitempty" validate:"omitempty,min=0"`
	AllowBackorder    *bool             `json:"allow_backorder,omitempty"`
	Options           map[string]string `json:"options,omitempty"`
	Images            []string          `json:"images,omitempty"`
}

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
func (s *Service) CreateProduct(tenantID uuid.UUID, req CreateProductRequest) (*Product, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Validate business rules
	if req.ComparePrice > 0 && req.Price >= req.ComparePrice {
		return nil, errors.New("compare price must be higher than selling price")
	}

	// Generate slug
	slug := s.generateSlug(req.Name)
	
	// Check if slug exists for this tenant
	if exists, err := s.repo.SlugExists(tenantID, slug); err != nil {
		return nil, err
	} else if exists {
		slug = s.generateUniqueSlug(tenantID, slug)
	}

	// Validate category if provided
	if req.CategoryID != nil {
		if exists, err := s.repo.CategoryExists(tenantID, *req.CategoryID); err != nil {
			return nil, err
		} else if !exists {
			return nil, errors.New("category not found")
		}
	}

	// Create product
	product := &Product{
		ID:                uuid.New(),
		TenantID:          tenantID,
		Name:              strings.TrimSpace(req.Name),
		Slug:              slug,
		Description:       strings.TrimSpace(req.Description),
		Type:              req.Type,
		Status:            ProductStatusDraft,
		Price:             req.Price,
		ComparePrice:      req.ComparePrice,
		CostPrice:         req.CostPrice,
		SKU:               strings.TrimSpace(req.SKU),
		Barcode:           strings.TrimSpace(req.Barcode),
		InventoryQuantity: req.InventoryQuantity,
		TrackQuantity:     req.TrackQuantity,
		AllowBackorder:    req.AllowBackorder,
		Weight:            req.Weight,
		Length:            req.Length,
		Width:             req.Width,
		Height:            req.Height,
		MetaTitle:         strings.TrimSpace(req.MetaTitle),
		MetaDescription:   strings.TrimSpace(req.MetaDescription),
		MetaKeywords:      strings.TrimSpace(req.MetaKeywords),
		CategoryID:        getUUIDValue(req.CategoryID),
		Tags:              req.Tags,
		Images:            req.Images,
	}

	// Set featured image
	if len(req.Images) > 0 {
		product.FeaturedImage = req.Images[0]
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
func (s *Service) UpdateProduct(tenantID uuid.UUID, id string, req UpdateProductRequest) (*Product, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	// Get existing product
	product, err := s.repo.FindProductByID(tenantID, productID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		product.Name = strings.TrimSpace(req.Name)
		// Regenerate slug if name changed
		newSlug := s.generateSlug(req.Name)
		if newSlug != product.Slug {
			if exists, err := s.repo.SlugExists(tenantID, newSlug); err != nil {
				return nil, err
			} else if exists {
				newSlug = s.generateUniqueSlug(tenantID, newSlug)
			}
			product.Slug = newSlug
		}
	}

	if req.Description != "" {
		product.Description = strings.TrimSpace(req.Description)
	}
	if req.Type != "" {
		product.Type = req.Type
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.ComparePrice >= 0 {
		product.ComparePrice = req.ComparePrice
	}
	if req.CostPrice >= 0 {
		product.CostPrice = req.CostPrice
	}
	if req.SKU != "" {
		product.SKU = strings.TrimSpace(req.SKU)
	}
	if req.Barcode != "" {
		product.Barcode = strings.TrimSpace(req.Barcode)
	}
	if req.InventoryQuantity != nil {
		product.InventoryQuantity = *req.InventoryQuantity
	}
	if req.TrackQuantity != nil {
		product.TrackQuantity = *req.TrackQuantity
	}
	if req.AllowBackorder != nil {
		product.AllowBackorder = *req.AllowBackorder
	}
	if req.Weight != nil {
		product.Weight = *req.Weight
	}
	if req.Length != nil {
		product.Length = *req.Length
	}
	if req.Width != nil {
		product.Width = *req.Width
	}
	if req.Height != nil {
		product.Height = *req.Height
	}
	if req.MetaTitle != "" {
		product.MetaTitle = strings.TrimSpace(req.MetaTitle)
	}
	if req.MetaDescription != "" {
		product.MetaDescription = strings.TrimSpace(req.MetaDescription)
	}
	if req.MetaKeywords != "" {
		product.MetaKeywords = strings.TrimSpace(req.MetaKeywords)
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Tags != nil {
		product.Tags = req.Tags
	}
	if req.Images != nil {
		product.Images = req.Images
		if len(req.Images) > 0 {
			product.FeaturedImage = req.Images[0]
		}
	}
	if req.Status != "" {
		product.Status = req.Status
	}

	// Validate business rules
	if product.ComparePrice > 0 && product.Price >= product.ComparePrice {
		return nil, errors.New("compare price must be higher than selling price")
	}

	return s.repo.UpdateProduct(product)
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
func (s *Service) CreateCategory(tenantID uuid.UUID, req CreateCategoryRequest) (*Category, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Generate slug
	slug := s.generateSlug(req.Name)
	
	// Check if slug exists for this tenant
	if exists, err := s.repo.CategorySlugExists(tenantID, slug); err != nil {
		return nil, err
	} else if exists {
		slug = s.generateUniqueCategorySlug(tenantID, slug)
	}

	// Validate parent category if provided
	if req.ParentID != nil {
		if exists, err := s.repo.CategoryExists(tenantID, *req.ParentID); err != nil {
			return nil, err
		} else if !exists {
			return nil, errors.New("parent category not found")
		}
	}

	// Create category
	category := &Category{
		ID:              uuid.New(),
		TenantID:        tenantID,
		Name:            strings.TrimSpace(req.Name),
		Slug:            slug,
		Description:     strings.TrimSpace(req.Description),
		ParentID:        req.ParentID,
		SortOrder:       req.SortOrder,
		IsActive:        true,
		MetaTitle:       strings.TrimSpace(req.MetaTitle),
		MetaDescription: strings.TrimSpace(req.MetaDescription),
	}

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

func getUUIDValue(ptr *uuid.UUID) uuid.UUID {
	if ptr == nil {
		return uuid.Nil
	}
	return *ptr
}

// TODO: Add more service methods
// - BulkUpdateProducts(tenantID uuid.UUID, productIDs []uuid.UUID, updates map[string]interface{}) error
// - ExportProducts(tenantID uuid.UUID, format string) ([]byte, error)
// - ImportProducts(tenantID uuid.UUID, data []byte) error
// - GetProductStats(tenantID uuid.UUID) (*ProductStats, error)
// - DuplicateProduct(tenantID uuid.UUID, productID uuid.UUID) (*Product, error)

// GetProductStats returns product statistics for a tenant
func (s *Service) GetProductStats(tenantID uuid.UUID) (map[string]interface{}, error) {
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
func (s *Service) CreateProductVariant(tenantID uuid.UUID, productIDStr string, req CreateVariantRequest) (*ProductVariant, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	// Verify product exists and belongs to tenant
	_, err = s.repo.FindProductByID(tenantID, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Create variant
	variant := &ProductVariant{
		ID:                uuid.New(),
		ProductID:         productID,
		Name:              strings.TrimSpace(req.Name),
		SKU:               strings.TrimSpace(req.SKU),
		Price:             req.Price,
		InventoryQuantity: req.InventoryQuantity,
		AllowBackorder:    req.AllowBackorder,
		Options:           req.Options,
		Images:            req.Images,
	}

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

	return s.repo.FindProductVariants(productID)
}

// UpdateProductVariant updates a product variant
func (s *Service) UpdateProductVariant(tenantID uuid.UUID, productIDStr, variantIDStr string, req UpdateVariantRequest) (*ProductVariant, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		return nil, errors.New("invalid variant ID")
	}

	// Verify product exists and belongs to tenant
	_, err = s.repo.FindProductByID(tenantID, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Get existing variant
	variants, err := s.repo.FindProductVariants(productID)
	if err != nil {
		return nil, err
	}

	var variant *ProductVariant
	for _, v := range variants {
		if v.ID == variantID {
			variant = v
			break
		}
	}

	if variant == nil {
		return nil, errors.New("variant not found")
	}

	// Update fields
	if req.Name != "" {
		variant.Name = strings.TrimSpace(req.Name)
	}
	if req.SKU != "" {
		variant.SKU = strings.TrimSpace(req.SKU)
	}
	if req.Price >= 0 {
		variant.Price = req.Price
	}
	if req.InventoryQuantity != nil {
		variant.InventoryQuantity = *req.InventoryQuantity
	}
	if req.AllowBackorder != nil {
		variant.AllowBackorder = *req.AllowBackorder
	}
	if req.Options != nil {
		variant.Options = req.Options
	}
	if req.Images != nil {
		variant.Images = req.Images
	}

	return s.repo.UpdateProductVariant(variant)
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

	return s.repo.DeleteProductVariant(variantID)
}

// Category management methods

// UpdateCategory updates an existing category
func (s *Service) UpdateCategory(tenantID uuid.UUID, categoryIDStr string, req UpdateCategoryRequest) (*Category, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	// Get existing category
	category, err := s.repo.FindCategoryByID(tenantID, categoryID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		category.Name = strings.TrimSpace(req.Name)
		// Regenerate slug if name changed
		newSlug := s.generateSlug(req.Name)
		if newSlug != category.Slug {
			if exists, err := s.repo.CategorySlugExists(tenantID, newSlug); err != nil {
				return nil, err
			} else if exists {
				newSlug = s.generateUniqueCategorySlug(tenantID, newSlug)
			}
			category.Slug = newSlug
		}
	}

	if req.Description != "" {
		category.Description = strings.TrimSpace(req.Description)
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.SortOrder != 0 {
		category.SortOrder = req.SortOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if req.MetaTitle != "" {
		category.MetaTitle = strings.TrimSpace(req.MetaTitle)
	}
	if req.MetaDescription != "" {
		category.MetaDescription = strings.TrimSpace(req.MetaDescription)
	}

	return s.repo.UpdateCategory(category)
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
