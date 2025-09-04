package category

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for category data operations
type Repository interface {
	// Category CRUD operations
	Save(ctx context.Context, category *Category) error
	FindByID(ctx context.Context, tenantID, categoryID uuid.UUID) (*Category, error)
	FindBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, tenantID, categoryID uuid.UUID) error
	
	// Category listing and filtering
	List(ctx context.Context, tenantID uuid.UUID, filter CategoryFilter, limit, offset int) ([]Category, error)
	Count(ctx context.Context, tenantID uuid.UUID, filter CategoryFilter) (int64, error)
	
	// Hierarchical operations
	FindChildren(ctx context.Context, tenantID, parentID uuid.UUID) ([]Category, error)
	FindByParent(ctx context.Context, tenantID uuid.UUID, parentID *uuid.UUID) ([]Category, error)
	FindRootCategories(ctx context.Context, tenantID uuid.UUID) ([]Category, error)
	GetCategoryTree(ctx context.Context, tenantID uuid.UUID, parentID *uuid.UUID) ([]Category, error)
	GetCategoryPath(ctx context.Context, tenantID, categoryID uuid.UUID) ([]Category, error)
	
	// Category validation
	ExistsBySlug(ctx context.Context, tenantID uuid.UUID, slug string, excludeID *uuid.UUID) (bool, error)
	HasChildren(ctx context.Context, tenantID, categoryID uuid.UUID) (bool, error)
	HasProducts(ctx context.Context, tenantID, categoryID uuid.UUID) (bool, error)
	ValidateParent(ctx context.Context, tenantID, categoryID, parentID uuid.UUID) error
	
	// Product associations
	AddProduct(ctx context.Context, tenantID, categoryID, productID uuid.UUID) error
	RemoveProduct(ctx context.Context, tenantID, categoryID, productID uuid.UUID) error
	GetCategoryProducts(ctx context.Context, tenantID, categoryID uuid.UUID, limit, offset int) ([]Product, error)
	UpdateProductCount(ctx context.Context, tenantID, categoryID uuid.UUID) error
	
	// Bulk operations
	UpdateChildrenPath(ctx context.Context, tenantID, parentID uuid.UUID, newPath string) error
	BulkUpdateStatus(ctx context.Context, tenantID uuid.UUID, categoryIDs []uuid.UUID, status CategoryStatus) error
	ReorderCategories(ctx context.Context, tenantID uuid.UUID, categoryOrders map[uuid.UUID]int) error
	
	// Statistics and analytics
	GetStats(ctx context.Context, tenantID uuid.UUID) (*CategoryStats, error)
	GetFeaturedCategories(ctx context.Context, tenantID uuid.UUID, limit int) ([]Category, error)
	GetPopularCategories(ctx context.Context, tenantID uuid.UUID, limit int) ([]Category, error)
	GetCategoriesByLevel(ctx context.Context, tenantID uuid.UUID, level int) ([]Category, error)
}

// GormRepository implements Repository using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM repository
func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// Save creates or updates a category
func (r *GormRepository) Save(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// FindByID finds a category by ID
func (r *GormRepository) FindByID(ctx context.Context, tenantID, categoryID uuid.UUID) (*Category, error) {
	var category Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("tenant_id = ? AND id = ?", tenantID, categoryID).
		First(&category).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	
	return &category, nil
}

// FindBySlug finds a category by slug
func (r *GormRepository) FindBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Category, error) {
	var category Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("tenant_id = ? AND slug = ?", tenantID, slug).
		First(&category).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	
	return &category, nil
}

// Update updates a category
func (r *GormRepository) Update(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete deletes a category
func (r *GormRepository) Delete(ctx context.Context, tenantID, categoryID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, categoryID).
		Delete(&Category{}).Error
}

// List returns categories with filtering and pagination
func (r *GormRepository) List(ctx context.Context, tenantID uuid.UUID, filter CategoryFilter, limit, offset int) ([]Category, error) {
	var categories []Category
	
	query := r.db.WithContext(ctx).
		Preload("Parent").
		Where("tenant_id = ?", tenantID)
	
	// Apply filters
	query = r.applyFilters(query, filter)
	
	err := query.
		Order("sort_order ASC, name ASC").
		Limit(limit).
		Offset(offset).
		Find(&categories).Error
	
	return categories, err
}

// Count returns the total count of categories matching the filter
func (r *GormRepository) Count(ctx context.Context, tenantID uuid.UUID, filter CategoryFilter) (int64, error) {
	var count int64
	
	query := r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ?", tenantID)
	
	// Apply filters
	query = r.applyFilters(query, filter)
	
	err := query.Count(&count).Error
	return count, err
}

// FindChildren finds direct children of a category
func (r *GormRepository) FindChildren(ctx context.Context, tenantID, parentID uuid.UUID) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND parent_id = ?", tenantID, parentID).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	
	return categories, err
}

// FindByParent finds categories by parent ID (nil for root categories)
func (r *GormRepository) FindByParent(ctx context.Context, tenantID uuid.UUID, parentID *uuid.UUID) ([]Category, error) {
	var categories []Category
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	
	err := query.
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	
	return categories, err
}

// FindRootCategories finds all root categories (no parent)
func (r *GormRepository) FindRootCategories(ctx context.Context, tenantID uuid.UUID) ([]Category, error) {
	return r.FindByParent(ctx, tenantID, nil)
}

// GetCategoryTree returns hierarchical category tree
func (r *GormRepository) GetCategoryTree(ctx context.Context, tenantID uuid.UUID, parentID *uuid.UUID) ([]Category, error) {
	var categories []Category
	query := r.db.WithContext(ctx).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, name ASC")
		}).
		Where("tenant_id = ?", tenantID)
	
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	
	err := query.
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	
	return categories, err
}

// GetCategoryPath returns the path from root to the specified category
func (r *GormRepository) GetCategoryPath(ctx context.Context, tenantID, categoryID uuid.UUID) ([]Category, error) {
	category, err := r.FindByID(ctx, tenantID, categoryID)
	if err != nil {
		return nil, err
	}
	
	var path []Category
	current := category
	
	// Build path from current to root
	for current != nil {
		path = append([]Category{*current}, path...) // Prepend to maintain order
		
		if current.ParentID == nil {
			break
		}
		
		// Get parent
		parent, err := r.FindByID(ctx, tenantID, *current.ParentID)
		if err != nil {
			break
		}
		current = parent
	}
	
	return path, nil
}

// ExistsBySlug checks if a category with the given slug exists
func (r *GormRepository) ExistsBySlug(ctx context.Context, tenantID uuid.UUID, slug string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ? AND slug = ?", tenantID, slug)
	
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	err := query.Count(&count).Error
	return count > 0, err
}

// HasChildren checks if a category has child categories
func (r *GormRepository) HasChildren(ctx context.Context, tenantID, categoryID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ? AND parent_id = ?", tenantID, categoryID).
		Count(&count).Error
	
	return count > 0, err
}

// HasProducts checks if a category has associated products
func (r *GormRepository) HasProducts(ctx context.Context, tenantID, categoryID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("product_categories").
		Joins("JOIN categories ON categories.id = product_categories.category_id").
		Where("categories.tenant_id = ? AND product_categories.category_id = ?", tenantID, categoryID).
		Count(&count).Error
	
	return count > 0, err
}

// ValidateParent validates if a parent-child relationship is valid
func (r *GormRepository) ValidateParent(ctx context.Context, tenantID, categoryID, parentID uuid.UUID) error {
	// Check if parent exists
	parent, err := r.FindByID(ctx, tenantID, parentID)
	if err != nil {
		return ErrInvalidParent
	}
	
	// Check for circular reference by traversing up the parent chain
	current := parent
	for current != nil {
		if current.ID == categoryID {
			return ErrCircularReference
		}
		
		if current.ParentID == nil {
			break
		}
		
		// Get next parent
		nextParent, err := r.FindByID(ctx, tenantID, *current.ParentID)
		if err != nil {
			break
		}
		current = nextParent
	}
	
	return nil
}

// AddProduct associates a product with a category
func (r *GormRepository) AddProduct(ctx context.Context, tenantID, categoryID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO product_categories (product_id, category_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		productID, categoryID,
	).Error
}

// RemoveProduct removes a product from a category
func (r *GormRepository) RemoveProduct(ctx context.Context, tenantID, categoryID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM product_categories WHERE product_id = ? AND category_id = ?",
		productID, categoryID,
	).Error
}

// GetCategoryProducts returns products in a category
func (r *GormRepository) GetCategoryProducts(ctx context.Context, tenantID, categoryID uuid.UUID, limit, offset int) ([]Product, error) {
	var products []Product
	err := r.db.WithContext(ctx).
		Joins("JOIN product_categories ON products.id = product_categories.product_id").
		Where("products.tenant_id = ? AND product_categories.category_id = ?", tenantID, categoryID).
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	
	return products, err
}

// UpdateProductCount updates the product count for a category
func (r *GormRepository) UpdateProductCount(ctx context.Context, tenantID, categoryID uuid.UUID) error {
	var count int64
	err := r.db.WithContext(ctx).
		Table("product_categories").
		Joins("JOIN products ON products.id = product_categories.product_id").
		Where("products.tenant_id = ? AND product_categories.category_id = ?", tenantID, categoryID).
		Count(&count).Error
	
	if err != nil {
		return err
	}
	
	return r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ? AND id = ?", tenantID, categoryID).
		Update("product_count", count).Error
}

// UpdateChildrenPath updates the path for all children of a category
func (r *GormRepository) UpdateChildrenPath(ctx context.Context, tenantID, parentID uuid.UUID, newPath string) error {
	return r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ? AND parent_id = ?", tenantID, parentID).
		Update("path", newPath).Error
}

// BulkUpdateStatus updates status for multiple categories
func (r *GormRepository) BulkUpdateStatus(ctx context.Context, tenantID uuid.UUID, categoryIDs []uuid.UUID, status CategoryStatus) error {
	return r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ? AND id IN ?", tenantID, categoryIDs).
		Update("status", status).Error
}

// ReorderCategories updates sort order for multiple categories
func (r *GormRepository) ReorderCategories(ctx context.Context, tenantID uuid.UUID, categoryOrders map[uuid.UUID]int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for categoryID, sortOrder := range categoryOrders {
			err := tx.Model(&Category{}).
				Where("tenant_id = ? AND id = ?", tenantID, categoryID).
				Update("sort_order", sortOrder).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// GetStats returns category statistics
func (r *GormRepository) GetStats(ctx context.Context, tenantID uuid.UUID) (*CategoryStats, error) {
	stats := &CategoryStats{}
	
	// Total categories
	err := r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ?", tenantID).
		Count(&stats.TotalCategories).Error
	if err != nil {
		return nil, err
	}
	
	// Active categories
	err = r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ? AND status = ?", tenantID, StatusActive).
		Count(&stats.ActiveCategories).Error
	if err != nil {
		return nil, err
	}
	
	// Root categories
	err = r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ? AND parent_id IS NULL", tenantID).
		Count(&stats.RootCategories).Error
	if err != nil {
		return nil, err
	}
	
	// Featured categories
	err = r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ? AND is_featured = ?", tenantID, true).
		Count(&stats.FeaturedCategories).Error
	if err != nil {
		return nil, err
	}
	
	// Max depth
	var maxLevel int
	err = r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ?", tenantID).
		Select("COALESCE(MAX(level), 0)").
		Scan(&maxLevel).Error
	if err != nil {
		return nil, err
	}
	stats.MaxDepth = maxLevel + 1 // Convert level to depth
	
	// Average products per category
	var avgProducts float64
	err = r.db.WithContext(ctx).
		Model(&Category{}).
		Where("tenant_id = ?", tenantID).
		Select("COALESCE(AVG(product_count), 0)").
		Scan(&avgProducts).Error
	if err != nil {
		return nil, err
	}
	stats.AvgProductsPerCategory = avgProducts
	
	return stats, nil
}

// GetFeaturedCategories returns featured categories
func (r *GormRepository) GetFeaturedCategories(ctx context.Context, tenantID uuid.UUID, limit int) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_featured = ? AND status = ?", tenantID, true, StatusActive).
		Order("sort_order ASC, name ASC").
		Limit(limit).
		Find(&categories).Error
	
	return categories, err
}

// GetPopularCategories returns categories with most products
func (r *GormRepository) GetPopularCategories(ctx context.Context, tenantID uuid.UUID, limit int) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ?", tenantID, StatusActive).
		Order("product_count DESC, name ASC").
		Limit(limit).
		Find(&categories).Error
	
	return categories, err
}

// GetCategoriesByLevel returns categories at a specific level
func (r *GormRepository) GetCategoriesByLevel(ctx context.Context, tenantID uuid.UUID, level int) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND level = ?", tenantID, level).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	
	return categories, err
}

// applyFilters applies filters to the query
func (r *GormRepository) applyFilters(query *gorm.DB, filter CategoryFilter) *gorm.DB {
	if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
	}
	
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	
	if filter.Level != nil {
		query = query.Where("level = ?", *filter.Level)
	}
	
	if filter.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filter.IsFeatured)
	}
	
	if filter.ShowInMenu != nil {
		query = query.Where("show_in_menu = ?", *filter.ShowInMenu)
	}
	
	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", 
			fmt.Sprintf("%%%s%%", filter.Search),
			fmt.Sprintf("%%%s%%", filter.Search))
	}
	
	return query
}