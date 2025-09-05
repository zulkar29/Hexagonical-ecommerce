package product

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the product repository interface
type Repository interface {
	// Product operations
	SaveProduct(product *Product) (*Product, error)
	FindProductByID(tenantID, productID uuid.UUID) (*Product, error)
	FindProductBySlug(tenantID uuid.UUID, slug string) (*Product, error)
	UpdateProduct(product *Product) (*Product, error)
	DeleteProduct(tenantID, productID uuid.UUID) error
	ListProducts(tenantID uuid.UUID, filter ProductListFilter, offset, limit int) ([]*Product, int64, error)
	SlugExists(tenantID uuid.UUID, slug string) (bool, error)
	ProductExists(tenantID, productID uuid.UUID) (bool, error)
	UpdateInventory(tenantID, productID uuid.UUID, quantity int) error
	GetProductsByCategoryID(tenantID, categoryID uuid.UUID, offset, limit int) ([]*Product, int64, error)
	GetLowStockProducts(tenantID uuid.UUID, threshold int) ([]*Product, error)
	BulkUpdateProducts(tenantID uuid.UUID, productIDs []uuid.UUID, updates map[string]interface{}) error
	GetPopularProducts(tenantID uuid.UUID, limit int) ([]*Product, error)
	GetNewProducts(tenantID uuid.UUID, limit int) ([]*Product, error)
	GetFeaturedProducts(tenantID uuid.UUID, limit int) ([]*Product, error)
	GetRelatedProducts(tenantID uuid.UUID, productID uuid.UUID, categoryID *uuid.UUID, limit int) ([]*Product, error)
	BulkImportProducts(products []*Product) error

	// Category operations
	SaveCategory(category *Category) (*Category, error)
	FindCategoryByID(tenantID, categoryID uuid.UUID) (*Category, error)
	UpdateCategory(category *Category) (*Category, error)
	DeleteCategory(tenantID, categoryID uuid.UUID) error
	ListCategories(tenantID uuid.UUID) ([]*Category, error)
	CategoryExists(tenantID, categoryID uuid.UUID) (bool, error)
	
	// GetCategory retrieves a category by ID
	GetCategory(tenantID, categoryID uuid.UUID) (*Category, error)
	CategorySlugExists(tenantID uuid.UUID, slug string) (bool, error)
	GetRootCategories(tenantID uuid.UUID) ([]*Category, error)
	GetCategoryChildren(tenantID, parentID uuid.UUID) ([]*Category, error)

	// Product variant operations
	SaveProductVariant(variant *ProductVariant) (*ProductVariant, error)
	FindProductVariants(tenantID, productID uuid.UUID) ([]*ProductVariant, error)
	GetProductVariant(tenantID, variantID uuid.UUID) (*ProductVariant, error)
	UpdateProductVariant(variant *ProductVariant) (*ProductVariant, error)
	DeleteProductVariant(tenantID, variantID uuid.UUID) error

	// Statistics and aggregations
	GetProductStats(tenantID uuid.UUID) (*ProductStats, error)
	SearchProducts(tenantID uuid.UUID, query string, offset, limit int) ([]*Product, int64, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new product repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Product operations

// SaveProduct creates a new product
func (r *repository) SaveProduct(product *Product) (*Product, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

// FindProductByID retrieves a product by ID
func (r *repository) FindProductByID(tenantID, productID uuid.UUID) (*Product, error) {
	var product Product
	err := r.db.Preload("Variants").Preload("Category").
		First(&product, "id = ? AND tenant_id = ?", productID, tenantID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindProductBySlug retrieves a product by slug
func (r *repository) FindProductBySlug(tenantID uuid.UUID, slug string) (*Product, error) {
	var product Product
	err := r.db.Preload("Variants").Preload("Category").
		First(&product, "slug = ? AND tenant_id = ?", slug, tenantID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateProduct updates an existing product
func (r *repository) UpdateProduct(product *Product) (*Product, error) {
	if err := r.db.Save(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

// DeleteProduct soft deletes a product
func (r *repository) DeleteProduct(tenantID, productID uuid.UUID) error {
	return r.db.Where("tenant_id = ?", tenantID).Delete(&Product{}, productID).Error
}

// ListProducts returns paginated products with filters
func (r *repository) ListProducts(tenantID uuid.UUID, filter ProductListFilter, offset, limit int) ([]*Product, int64, error) {
	var products []*Product
	var total int64

	query := r.db.Model(&Product{}).Where("tenant_id = ?", tenantID)

	// Apply filters
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}
	if filter.InStock != nil {
		if *filter.InStock {
			query = query.Where("(track_quantity = false OR inventory_quantity > 0)")
		} else {
			query = query.Where("track_quantity = true AND inventory_quantity <= 0")
		}
	}
	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?",
			"%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with preloads
	if err := query.Preload("Category").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// SlugExists checks if a product slug exists for a tenant
func (r *repository) SlugExists(tenantID uuid.UUID, slug string) (bool, error) {
	var count int64
	err := r.db.Model(&Product{}).Where("tenant_id = ? AND slug = ?", tenantID, slug).Count(&count).Error
	return count > 0, err
}

// ProductExists checks if a product exists
func (r *repository) ProductExists(tenantID, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&Product{}).Where("tenant_id = ? AND id = ?", tenantID, productID).Count(&count).Error
	return count > 0, err
}

// UpdateInventory updates product inventory quantity
func (r *repository) UpdateInventory(tenantID, productID uuid.UUID, quantity int) error {
	return r.db.Model(&Product{}).
		Where("id = ? AND tenant_id = ?", productID, tenantID).
		Update("inventory_quantity", quantity).Error
}

// GetProductsByCategoryID returns products in a specific category
func (r *repository) GetProductsByCategoryID(tenantID, categoryID uuid.UUID, offset, limit int) ([]*Product, int64, error) {
	var products []*Product
	var total int64

	query := r.db.Model(&Product{}).Where("tenant_id = ? AND category_id = ?", tenantID, categoryID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetLowStockProducts returns products with low inventory
func (r *repository) GetLowStockProducts(tenantID uuid.UUID, threshold int) ([]*Product, error) {
	var products []*Product
	err := r.db.Where("tenant_id = ? AND track_quantity = true AND inventory_quantity <= ?", 
		tenantID, threshold).Find(&products).Error
	return products, err
}

// BulkUpdateProducts updates multiple products
func (r *repository) BulkUpdateProducts(tenantID uuid.UUID, productIDs []uuid.UUID, updates map[string]interface{}) error {
	return r.db.Model(&Product{}).
		Where("tenant_id = ? AND id IN ?", tenantID, productIDs).
		Updates(updates).Error
}

// Category operations

// SaveCategory creates a new category
func (r *repository) SaveCategory(category *Category) (*Category, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// FindCategoryByID retrieves a category by ID
func (r *repository) FindCategoryByID(tenantID, categoryID uuid.UUID) (*Category, error) {
	var category Category
	err := r.db.Preload("Parent").Preload("Children").
		First(&category, "id = ? AND tenant_id = ?", categoryID, tenantID).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory updates an existing category
func (r *repository) UpdateCategory(category *Category) (*Category, error) {
	if err := r.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// DeleteCategory soft deletes a category
func (r *repository) DeleteCategory(tenantID, categoryID uuid.UUID) error {
	return r.db.Where("tenant_id = ?", tenantID).Delete(&Category{}, categoryID).Error
}

// ListCategories returns all categories for a tenant
func (r *repository) ListCategories(tenantID uuid.UUID) ([]*Category, error) {
	var categories []*Category
	err := r.db.Where("tenant_id = ? AND is_active = true", tenantID).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// CategoryExists checks if a category exists
func (r *repository) CategoryExists(tenantID, categoryID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&Category{}).Where("tenant_id = ? AND id = ?", tenantID, categoryID).Count(&count).Error
	return count > 0, err
}

// GetCategory retrieves a category by ID
func (r *repository) GetCategory(tenantID, categoryID uuid.UUID) (*Category, error) {
	var category Category
	err := r.db.Preload("Parent").Preload("Children").
		First(&category, "id = ? AND tenant_id = ?", categoryID, tenantID).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// CategorySlugExists checks if a category slug exists for a tenant
func (r *repository) CategorySlugExists(tenantID uuid.UUID, slug string) (bool, error) {
	var count int64
	err := r.db.Model(&Category{}).Where("tenant_id = ? AND slug = ?", tenantID, slug).Count(&count).Error
	return count > 0, err
}

// GetRootCategories returns top-level categories
func (r *repository) GetRootCategories(tenantID uuid.UUID) ([]*Category, error) {
	var categories []*Category
	err := r.db.Where("tenant_id = ? AND parent_id IS NULL AND is_active = true", tenantID).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetCategoryChildren returns child categories
func (r *repository) GetCategoryChildren(tenantID, parentID uuid.UUID) ([]*Category, error) {
	var categories []*Category
	err := r.db.Where("tenant_id = ? AND parent_id = ? AND is_active = true", tenantID, parentID).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// Product Variant operations

// SaveProductVariant creates a new product variant
func (r *repository) SaveProductVariant(variant *ProductVariant) (*ProductVariant, error) {
	if err := r.db.Create(variant).Error; err != nil {
		return nil, err
	}
	return variant, nil
}

// FindProductVariants returns all variants for a product
func (r *repository) FindProductVariants(tenantID, productID uuid.UUID) ([]*ProductVariant, error) {
	var variants []*ProductVariant
	// Join with products table to filter by tenant_id since variants don't have tenant_id directly
	err := r.db.Joins("JOIN products ON product_variants.product_id = products.id").
		Where("products.tenant_id = ? AND product_variants.product_id = ?", tenantID, productID).
		Order("product_variants.created_at ASC").
		Find(&variants).Error
	return variants, err
}

// GetProductVariant returns a specific product variant
func (r *repository) GetProductVariant(tenantID, variantID uuid.UUID) (*ProductVariant, error) {
	var variant ProductVariant
	// Join with products table to filter by tenant_id since variants don't have tenant_id directly
	err := r.db.Joins("JOIN products ON product_variants.product_id = products.id").
		Where("products.tenant_id = ? AND product_variants.id = ?", tenantID, variantID).
		First(&variant).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

// UpdateProductVariant updates a product variant
func (r *repository) UpdateProductVariant(variant *ProductVariant) (*ProductVariant, error) {
	if err := r.db.Save(variant).Error; err != nil {
		return nil, err
	}
	return variant, nil
}

// DeleteProductVariant deletes a product variant
func (r *repository) DeleteProductVariant(tenantID, variantID uuid.UUID) error {
	return r.db.Delete(&ProductVariant{}, "id = ? AND tenant_id = ?", variantID, tenantID).Error
}

// Statistics and aggregations

// GetProductStats returns product statistics for a tenant
func (r *repository) GetProductStats(tenantID uuid.UUID) (*ProductStats, error) {
	stats := &ProductStats{}

	// Total products
	r.db.Model(&Product{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalProducts)

	// Products by status
	r.db.Model(&Product{}).Where("tenant_id = ? AND status = ?", tenantID, ProductStatusActive).Count(&stats.ActiveProducts)
	r.db.Model(&Product{}).Where("tenant_id = ? AND status = ?", tenantID, ProductStatusDraft).Count(&stats.DraftProducts)

	// Out of stock products
	r.db.Model(&Product{}).Where("tenant_id = ? AND track_quantity = true AND inventory_quantity <= 0", tenantID).Count(&stats.OutOfStock)

	// Low stock products (< 10)
	r.db.Model(&Product{}).Where("tenant_id = ? AND track_quantity = true AND inventory_quantity > 0 AND inventory_quantity < 10", tenantID).Count(&stats.LowStock)

	// Total categories
	r.db.Model(&Category{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalCategories)

	// Total inventory value
	r.db.Model(&Product{}).Where("tenant_id = ?", tenantID).Select("COALESCE(SUM(price * inventory_quantity), 0)").Scan(&stats.TotalValue)

	return stats, nil
}

// Search operations

// SearchProducts performs full-text search on products
func (r *repository) SearchProducts(tenantID uuid.UUID, query string, offset, limit int) ([]*Product, int64, error) {
	var products []*Product
	var total int64

	searchQuery := r.db.Model(&Product{}).Where("tenant_id = ?", tenantID).Where(
		"name ILIKE ? OR description ILIKE ? OR sku ILIKE ? OR meta_keywords ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%",
	)

	// Get total count
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := searchQuery.Preload("Category").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetPopularProducts returns products sorted by popularity (order count)
func (r *repository) GetPopularProducts(tenantID uuid.UUID, limit int) ([]*Product, error) {
	var products []*Product
	err := r.db.Where("tenant_id = ? AND status = ?", tenantID, ProductStatusActive).
		Order("order_count DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// GetNewProducts returns recently created products
func (r *repository) GetNewProducts(tenantID uuid.UUID, limit int) ([]*Product, error) {
	var products []*Product
	err := r.db.Where("tenant_id = ? AND status = ?", tenantID, ProductStatusActive).
		Order("created_at DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// GetFeaturedProducts returns featured products
func (r *repository) GetFeaturedProducts(tenantID uuid.UUID, limit int) ([]*Product, error) {
	var products []*Product
	err := r.db.Where("tenant_id = ? AND status = ? AND is_featured = ?", tenantID, ProductStatusActive, true).
		Order("created_at DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// GetRelatedProducts returns products related to a given product (same category)
func (r *repository) GetRelatedProducts(tenantID uuid.UUID, productID uuid.UUID, categoryID *uuid.UUID, limit int) ([]*Product, error) {
	var products []*Product
	query := r.db.Where("tenant_id = ? AND status = ? AND id != ?", tenantID, ProductStatusActive, productID)
	
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	
	err := query.Order("created_at DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// BulkImportProducts imports multiple products at once
func (r *repository) BulkImportProducts(products []*Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, product := range products {
			if err := tx.Create(product).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
