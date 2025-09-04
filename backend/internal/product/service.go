package product

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

// Request/Response DTOs
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
	repo      *Repository
	validator *validator.Validate
}

// NewService creates a new product service
func NewService(repo *Repository) *Service {
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
		Status:            StatusDraft,
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
	// Simple slug generation - in production, use a proper slug library
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters
	// TODO: Implement proper slug generation with regex
	return slug
}

func (s *Service) generateUniqueSlug(tenantID uuid.UUID, baseSlug string) string {
	// TODO: Implement unique slug generation by appending numbers
	return baseSlug + "-" + uuid.New().String()[:8]
}

func (s *Service) generateUniqueCategorySlug(tenantID uuid.UUID, baseSlug string) string {
	// TODO: Implement unique category slug generation
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
