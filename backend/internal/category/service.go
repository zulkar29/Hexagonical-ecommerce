package category

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Service defines the interface for category business logic
type Service interface {
	// Category CRUD operations
	CreateCategory(ctx context.Context, tenantID uuid.UUID, req CreateCategoryRequest) (*CategoryResponse, error)
	GetCategory(ctx context.Context, tenantID, categoryID uuid.UUID) (*CategoryResponse, error)
	GetCategoryBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*CategoryResponse, error)
	UpdateCategory(ctx context.Context, tenantID, categoryID uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error)
	DeleteCategory(ctx context.Context, tenantID, categoryID uuid.UUID) error
	
	// Category listing and filtering
	ListCategories(ctx context.Context, tenantID uuid.UUID, filter CategoryFilter, limit, offset int) ([]CategoryResponse, int64, error)
	GetCategoryTree(ctx context.Context, tenantID uuid.UUID, parentID *uuid.UUID) ([]CategoryTreeResponse, error)
	GetCategoryPath(ctx context.Context, tenantID, categoryID uuid.UUID) ([]CategoryResponse, error)
	
	// Category management
	MoveCategory(ctx context.Context, tenantID, categoryID uuid.UUID, newParentID *uuid.UUID) error
	ReorderCategories(ctx context.Context, tenantID uuid.UUID, categoryOrders map[uuid.UUID]int) error
	BulkUpdateStatus(ctx context.Context, tenantID uuid.UUID, categoryIDs []uuid.UUID, status CategoryStatus) error
	
	// Product associations
	AddProductToCategory(ctx context.Context, tenantID, categoryID, productID uuid.UUID) error
	RemoveProductFromCategory(ctx context.Context, tenantID, categoryID, productID uuid.UUID) error
	GetCategoryProducts(ctx context.Context, tenantID, categoryID uuid.UUID, limit, offset int) ([]Product, error)
	
	// Statistics and analytics
	GetCategoryStats(ctx context.Context, tenantID uuid.UUID) (*CategoryStats, error)
	GetFeaturedCategories(ctx context.Context, tenantID uuid.UUID, limit int) ([]CategoryResponse, error)
	GetPopularCategories(ctx context.Context, tenantID uuid.UUID, limit int) ([]CategoryResponse, error)
	
	// Utility operations
	ValidateSlug(ctx context.Context, tenantID uuid.UUID, slug string, excludeID *uuid.UUID) error
	GenerateUniqueSlug(ctx context.Context, tenantID uuid.UUID, name string, excludeID *uuid.UUID) (string, error)
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new category service
func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

const (
	MaxCategoryDepth = 5
	MaxSlugLength    = 100
)

// CreateCategory creates a new category
func (s *ServiceImpl) CreateCategory(ctx context.Context, tenantID uuid.UUID, req CreateCategoryRequest) (*CategoryResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(ctx, tenantID, req); err != nil {
		return nil, err
	}
	
	// Create category entity
	category := &Category{
		ID:              uuid.New(),
		TenantID:        tenantID,
		Name:            req.Name,
		Slug:            req.Slug,
		Description:     req.Description,
		Image:           req.Image,
		Icon:            req.Icon,
		ParentID:        req.ParentID,
		SortOrder:       req.SortOrder,
		Status:          StatusActive,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		MetaKeywords:    req.MetaKeywords,
		IsFeatured:      req.IsFeatured,
		ShowInMenu:      req.ShowInMenu,
	}
	
	// Generate slug if not provided
	if category.Slug == "" {
		slug, err := s.GenerateUniqueSlug(ctx, tenantID, category.Name, nil)
		if err != nil {
			return nil, err
		}
		category.Slug = slug
	}
	
	// Set hierarchy information
	if err := s.setHierarchyInfo(ctx, tenantID, category); err != nil {
		return nil, err
	}
	
	// Save category
	if err := s.repo.Save(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	
	return s.buildCategoryResponse(category), nil
}

// GetCategory retrieves a category by ID
func (s *ServiceImpl) GetCategory(ctx context.Context, tenantID, categoryID uuid.UUID) (*CategoryResponse, error) {
	category, err := s.repo.FindByID(ctx, tenantID, categoryID)
	if err != nil {
		return nil, err
	}
	
	return s.buildCategoryResponse(category), nil
}

// GetCategoryBySlug retrieves a category by slug
func (s *ServiceImpl) GetCategoryBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*CategoryResponse, error) {
	category, err := s.repo.FindBySlug(ctx, tenantID, slug)
	if err != nil {
		return nil, err
	}
	
	return s.buildCategoryResponse(category), nil
}

// UpdateCategory updates an existing category
func (s *ServiceImpl) UpdateCategory(ctx context.Context, tenantID, categoryID uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error) {
	// Get existing category
	category, err := s.repo.FindByID(ctx, tenantID, categoryID)
	if err != nil {
		return nil, err
	}
	
	// Validate update request
	if err := s.validateUpdateRequest(ctx, tenantID, categoryID, req); err != nil {
		return nil, err
	}
	
	// Update fields
	s.updateCategoryFields(category, req)
	
	// Handle parent change
	if req.ParentID != nil && (category.ParentID == nil || *category.ParentID != *req.ParentID) {
		if err := s.setHierarchyInfo(ctx, tenantID, category); err != nil {
			return nil, err
		}
	}
	
	// Save updated category
	if err := s.repo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}
	
	// Update children paths if parent changed
	if req.ParentID != nil {
		if err := s.updateChildrenPaths(ctx, tenantID, category); err != nil {
			return nil, err
		}
	}
	
	return s.buildCategoryResponse(category), nil
}

// DeleteCategory deletes a category
func (s *ServiceImpl) DeleteCategory(ctx context.Context, tenantID, categoryID uuid.UUID) error {
	// Get category
	category, err := s.repo.FindByID(ctx, tenantID, categoryID)
	if err != nil {
		return err
	}
	
	// Check if category can be deleted
	if !category.CanDelete() {
		if category.HasChildren() {
			return ErrCategoryHasChildren
		}
		if category.HasProducts() {
			return ErrCategoryHasProducts
		}
	}
	
	// Delete category
	return s.repo.Delete(ctx, tenantID, categoryID)
}

// ListCategories returns paginated list of categories
func (s *ServiceImpl) ListCategories(ctx context.Context, tenantID uuid.UUID, filter CategoryFilter, limit, offset int) ([]CategoryResponse, int64, error) {
	categories, err := s.repo.List(ctx, tenantID, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	count, err := s.repo.Count(ctx, tenantID, filter)
	if err != nil {
		return nil, 0, err
	}
	
	responses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *s.buildCategoryResponse(&category)
	}
	
	return responses, count, nil
}

// GetCategoryTree returns hierarchical category tree
func (s *ServiceImpl) GetCategoryTree(ctx context.Context, tenantID uuid.UUID, parentID *uuid.UUID) ([]CategoryTreeResponse, error) {
	categories, err := s.repo.GetCategoryTree(ctx, tenantID, parentID)
	if err != nil {
		return nil, err
	}
	
	return s.buildCategoryTreeResponse(categories), nil
}

// GetCategoryPath returns the path from root to the specified category
func (s *ServiceImpl) GetCategoryPath(ctx context.Context, tenantID, categoryID uuid.UUID) ([]CategoryResponse, error) {
	path, err := s.repo.GetCategoryPath(ctx, tenantID, categoryID)
	if err != nil {
		return nil, err
	}
	
	responses := make([]CategoryResponse, len(path))
	for i, category := range path {
		responses[i] = *s.buildCategoryResponse(&category)
	}
	
	return responses, nil
}

// MoveCategory moves a category to a new parent
func (s *ServiceImpl) MoveCategory(ctx context.Context, tenantID, categoryID uuid.UUID, newParentID *uuid.UUID) error {
	// Validate parent if provided
	if newParentID != nil {
		if err := s.repo.ValidateParent(ctx, tenantID, categoryID, *newParentID); err != nil {
			return err
		}
	}
	
	// Get category
	category, err := s.repo.FindByID(ctx, tenantID, categoryID)
	if err != nil {
		return err
	}
	
	// Update parent
	category.ParentID = newParentID
	
	// Update hierarchy info
	if err := s.setHierarchyInfo(ctx, tenantID, category); err != nil {
		return err
	}
	
	// Save category
	if err := s.repo.Update(ctx, category); err != nil {
		return err
	}
	
	// Update children paths
	return s.updateChildrenPaths(ctx, tenantID, category)
}

// ReorderCategories updates sort order for multiple categories
func (s *ServiceImpl) ReorderCategories(ctx context.Context, tenantID uuid.UUID, categoryOrders map[uuid.UUID]int) error {
	return s.repo.ReorderCategories(ctx, tenantID, categoryOrders)
}

// BulkUpdateStatus updates status for multiple categories
func (s *ServiceImpl) BulkUpdateStatus(ctx context.Context, tenantID uuid.UUID, categoryIDs []uuid.UUID, status CategoryStatus) error {
	return s.repo.BulkUpdateStatus(ctx, tenantID, categoryIDs, status)
}

// AddProductToCategory associates a product with a category
func (s *ServiceImpl) AddProductToCategory(ctx context.Context, tenantID, categoryID, productID uuid.UUID) error {
	// Verify category exists
	_, err := s.repo.FindByID(ctx, tenantID, categoryID)
	if err != nil {
		return err
	}
	
	// Add product to category
	if err := s.repo.AddProduct(ctx, tenantID, categoryID, productID); err != nil {
		return err
	}
	
	// Update product count
	return s.repo.UpdateProductCount(ctx, tenantID, categoryID)
}

// RemoveProductFromCategory removes a product from a category
func (s *ServiceImpl) RemoveProductFromCategory(ctx context.Context, tenantID, categoryID, productID uuid.UUID) error {
	// Remove product from category
	if err := s.repo.RemoveProduct(ctx, tenantID, categoryID, productID); err != nil {
		return err
	}
	
	// Update product count
	return s.repo.UpdateProductCount(ctx, tenantID, categoryID)
}

// GetCategoryProducts returns products in a category
func (s *ServiceImpl) GetCategoryProducts(ctx context.Context, tenantID, categoryID uuid.UUID, limit, offset int) ([]Product, error) {
	return s.repo.GetCategoryProducts(ctx, tenantID, categoryID, limit, offset)
}

// GetCategoryStats returns category statistics
func (s *ServiceImpl) GetCategoryStats(ctx context.Context, tenantID uuid.UUID) (*CategoryStats, error) {
	return s.repo.GetStats(ctx, tenantID)
}

// GetFeaturedCategories returns featured categories
func (s *ServiceImpl) GetFeaturedCategories(ctx context.Context, tenantID uuid.UUID, limit int) ([]CategoryResponse, error) {
	categories, err := s.repo.GetFeaturedCategories(ctx, tenantID, limit)
	if err != nil {
		return nil, err
	}
	
	responses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *s.buildCategoryResponse(&category)
	}
	
	return responses, nil
}

// GetPopularCategories returns popular categories
func (s *ServiceImpl) GetPopularCategories(ctx context.Context, tenantID uuid.UUID, limit int) ([]CategoryResponse, error) {
	categories, err := s.repo.GetPopularCategories(ctx, tenantID, limit)
	if err != nil {
		return nil, err
	}
	
	responses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *s.buildCategoryResponse(&category)
	}
	
	return responses, nil
}

// ValidateSlug validates if a slug is available
func (s *ServiceImpl) ValidateSlug(ctx context.Context, tenantID uuid.UUID, slug string, excludeID *uuid.UUID) error {
	if slug == "" {
		return ErrInvalidSlug
	}
	
	if len(slug) > MaxSlugLength {
		return ErrInvalidSlug
	}
	
	exists, err := s.repo.ExistsBySlug(ctx, tenantID, slug, excludeID)
	if err != nil {
		return err
	}
	
	if exists {
		return ErrCategoryExists
	}
	
	return nil
}

// GenerateUniqueSlug generates a unique slug for a category
func (s *ServiceImpl) GenerateUniqueSlug(ctx context.Context, tenantID uuid.UUID, name string, excludeID *uuid.UUID) (string, error) {
	baseSlug := s.generateSlugFromName(name)
	slug := baseSlug
	counter := 1
	
	for {
		err := s.ValidateSlug(ctx, tenantID, slug, excludeID)
		if err == nil {
			return slug, nil
		}
		
		if err != ErrCategoryExists {
			return "", err
		}
		
		// Try with counter
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
		
		if counter > 100 {
			return "", fmt.Errorf("unable to generate unique slug")
		}
	}
}

// Helper methods

// validateCreateRequest validates category creation request
func (s *ServiceImpl) validateCreateRequest(ctx context.Context, tenantID uuid.UUID, req CreateCategoryRequest) error {
	if req.Name == "" {
		return fmt.Errorf("category name is required")
	}
	
	// Validate slug if provided
	if req.Slug != "" {
		if err := s.ValidateSlug(ctx, tenantID, req.Slug, nil); err != nil {
			return err
		}
	}
	
	// Validate parent if provided
	if req.ParentID != nil {
		parent, err := s.repo.FindByID(ctx, tenantID, *req.ParentID)
		if err != nil {
			return ErrInvalidParent
		}
		
		// Check depth limit
		if parent.Level >= MaxCategoryDepth-1 {
			return ErrMaxDepthExceeded
		}
	}
	
	return nil
}

// validateUpdateRequest validates category update request
func (s *ServiceImpl) validateUpdateRequest(ctx context.Context, tenantID, categoryID uuid.UUID, req UpdateCategoryRequest) error {
	// Validate slug if provided
	if req.Slug != "" {
		if err := s.ValidateSlug(ctx, tenantID, req.Slug, &categoryID); err != nil {
			return err
		}
	}
	
	// Validate parent if provided
	if req.ParentID != nil {
		if err := s.repo.ValidateParent(ctx, tenantID, categoryID, *req.ParentID); err != nil {
			return err
		}
		
		// Check depth limit
		parent, err := s.repo.FindByID(ctx, tenantID, *req.ParentID)
		if err != nil {
			return ErrInvalidParent
		}
		
		if parent.Level >= MaxCategoryDepth-1 {
			return ErrMaxDepthExceeded
		}
	}
	
	return nil
}

// updateCategoryFields updates category fields from request
func (s *ServiceImpl) updateCategoryFields(category *Category, req UpdateCategoryRequest) {
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Image != "" {
		category.Image = req.Image
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.Status != "" {
		category.Status = req.Status
	}
	if req.MetaTitle != "" {
		category.MetaTitle = req.MetaTitle
	}
	if req.MetaDescription != "" {
		category.MetaDescription = req.MetaDescription
	}
	if req.MetaKeywords != "" {
		category.MetaKeywords = req.MetaKeywords
	}
	if req.IsFeatured != nil {
		category.IsFeatured = *req.IsFeatured
	}
	if req.ShowInMenu != nil {
		category.ShowInMenu = *req.ShowInMenu
	}
}

// setHierarchyInfo sets hierarchy information for a category
func (s *ServiceImpl) setHierarchyInfo(ctx context.Context, tenantID uuid.UUID, category *Category) error {
	if category.ParentID == nil {
		category.Level = 0
		category.Path = ""
		return nil
	}
	
	parent, err := s.repo.FindByID(ctx, tenantID, *category.ParentID)
	if err != nil {
		return ErrInvalidParent
	}
	
	category.UpdatePath(parent)
	return nil
}

// updateChildrenPaths updates paths for all children of a category
func (s *ServiceImpl) updateChildrenPaths(ctx context.Context, tenantID uuid.UUID, category *Category) error {
	newPath := category.GetFullPath()
	return s.repo.UpdateChildrenPath(ctx, tenantID, category.ID, newPath)
}

// generateSlugFromName generates a slug from category name
func (s *ServiceImpl) generateSlugFromName(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "&", "and")
	slug = strings.ReplaceAll(slug, "'", "")
	
	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// buildCategoryResponse builds category response
func (s *ServiceImpl) buildCategoryResponse(category *Category) *CategoryResponse {
	return &CategoryResponse{
		Category:      category,
		ChildrenCount: len(category.Children),
		ProductCount:  category.ProductCount,
	}
}

// buildCategoryTreeResponse builds hierarchical category tree response
func (s *ServiceImpl) buildCategoryTreeResponse(categories []Category) []CategoryTreeResponse {
	responses := make([]CategoryTreeResponse, len(categories))
	
	for i, category := range categories {
		responses[i] = CategoryTreeResponse{
			Category: &category,
			Children: s.buildCategoryTreeResponse(category.Children),
		}
	}
	
	return responses
}