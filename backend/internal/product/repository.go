package product

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles product data operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new product repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Product operations

// SaveProduct creates a new product
func (r *Repository) SaveProduct(product *Product) (*Product, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

// FindProductByID retrieves a product by ID
func (r *Repository) FindProductByID(tenantID, productID uuid.UUID) (*Product, error) {
	var product Product
	err := r.db.Preload("Variants").Preload("Category").
		First(&product, "id = ? AND tenant_id = ?", productID, tenantID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindProductBySlug retrieves a product by slug
func (r *Repository) FindProductBySlug(tenantID uuid.UUID, slug string) (*Product, error) {
	var product Product
	err := r.db.Preload("Variants").Preload("Category").
		First(&product, "slug = ? AND tenant_id = ?", slug, tenantID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateProduct updates an existing product
func (r *Repository) UpdateProduct(product *Product) (*Product, error) {
	if err := r.db.Save(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

// DeleteProduct soft deletes a product
func (r *Repository) DeleteProduct(tenantID, productID uuid.UUID) error {
	return r.db.Where("tenant_id = ?", tenantID).Delete(&Product{}, productID).Error
}

// ListProducts returns paginated products with filters
func (r *Repository) ListProducts(tenantID uuid.UUID, filter ProductListFilter, offset, limit int) ([]*Product, int64, error) {
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
func (r *Repository) SlugExists(tenantID uuid.UUID, slug string) (bool, error) {
	var count int64
	err := r.db.Model(&Product{}).Where("tenant_id = ? AND slug = ?", tenantID, slug).Count(&count).Error
	return count > 0, err
}

// UpdateInventory updates product inventory quantity
func (r *Repository) UpdateInventory(tenantID, productID uuid.UUID, quantity int) error {
	return r.db.Model(&Product{}).
		Where("id = ? AND tenant_id = ?", productID, tenantID).
		Update("inventory_quantity", quantity).Error
}

// GetProductsByCategoryID returns products in a specific category
func (r *Repository) GetProductsByCategoryID(tenantID, categoryID uuid.UUID, offset, limit int) ([]*Product, int64, error) {
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
func (r *Repository) GetLowStockProducts(tenantID uuid.UUID, threshold int) ([]*Product, error) {
	var products []*Product
	err := r.db.Where("tenant_id = ? AND track_quantity = true AND inventory_quantity <= ?", 
		tenantID, threshold).Find(&products).Error
	return products, err
}

// BulkUpdateProducts updates multiple products
func (r *Repository) BulkUpdateProducts(tenantID uuid.UUID, productIDs []uuid.UUID, updates map[string]interface{}) error {
	return r.db.Model(&Product{}).
		Where("tenant_id = ? AND id IN ?", tenantID, productIDs).
		Updates(updates).Error
}

// Category operations

// SaveCategory creates a new category
func (r *Repository) SaveCategory(category *Category) (*Category, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// FindCategoryByID retrieves a category by ID
func (r *Repository) FindCategoryByID(tenantID, categoryID uuid.UUID) (*Category, error) {
	var category Category
	err := r.db.Preload("Parent").Preload("Children").
		First(&category, "id = ? AND tenant_id = ?", categoryID, tenantID).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory updates an existing category
func (r *Repository) UpdateCategory(category *Category) (*Category, error) {
	if err := r.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// DeleteCategory soft deletes a category
func (r *Repository) DeleteCategory(tenantID, categoryID uuid.UUID) error {
	return r.db.Where("tenant_id = ?", tenantID).Delete(&Category{}, categoryID).Error
}

// ListCategories returns all categories for a tenant
func (r *Repository) ListCategories(tenantID uuid.UUID) ([]*Category, error) {
	var categories []*Category
	err := r.db.Where("tenant_id = ? AND is_active = true", tenantID).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// CategoryExists checks if a category exists
func (r *Repository) CategoryExists(tenantID, categoryID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&Category{}).Where("tenant_id = ? AND id = ?", tenantID, categoryID).Count(&count).Error
	return count > 0, err
}

// CategorySlugExists checks if a category slug exists for a tenant
func (r *Repository) CategorySlugExists(tenantID uuid.UUID, slug string) (bool, error) {
	var count int64
	err := r.db.Model(&Category{}).Where("tenant_id = ? AND slug = ?", tenantID, slug).Count(&count).Error
	return count > 0, err
}

// GetRootCategories returns top-level categories
func (r *Repository) GetRootCategories(tenantID uuid.UUID) ([]*Category, error) {
	var categories []*Category
	err := r.db.Where("tenant_id = ? AND parent_id IS NULL AND is_active = true", tenantID).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetCategoryChildren returns child categories
func (r *Repository) GetCategoryChildren(tenantID, parentID uuid.UUID) ([]*Category, error) {
	var categories []*Category
	err := r.db.Where("tenant_id = ? AND parent_id = ? AND is_active = true", tenantID, parentID).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// Product Variant operations

// SaveProductVariant creates a new product variant
func (r *Repository) SaveProductVariant(variant *ProductVariant) (*ProductVariant, error) {
	if err := r.db.Create(variant).Error; err != nil {
		return nil, err
	}
	return variant, nil
}

// FindProductVariants returns all variants for a product
func (r *Repository) FindProductVariants(productID uuid.UUID) ([]*ProductVariant, error) {
	var variants []*ProductVariant
	err := r.db.Where("product_id = ?", productID).Order("created_at ASC").Find(&variants).Error
	return variants, err
}

// UpdateProductVariant updates a product variant
func (r *Repository) UpdateProductVariant(variant *ProductVariant) (*ProductVariant, error) {
	if err := r.db.Save(variant).Error; err != nil {
		return nil, err
	}
	return variant, nil
}

// DeleteProductVariant deletes a product variant
func (r *Repository) DeleteProductVariant(variantID uuid.UUID) error {
	return r.db.Delete(&ProductVariant{}, variantID).Error
}

// Statistics and aggregations

// GetProductStats returns product statistics for a tenant
func (r *Repository) GetProductStats(tenantID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total products
	var totalProducts int64
	r.db.Model(&Product{}).Where("tenant_id = ?", tenantID).Count(&totalProducts)
	stats["total_products"] = totalProducts

	// Products by status
	var activeProducts int64
	r.db.Model(&Product{}).Where("tenant_id = ? AND status = ?", tenantID, StatusActive).Count(&activeProducts)
	stats["active_products"] = activeProducts

	var draftProducts int64
	r.db.Model(&Product{}).Where("tenant_id = ? AND status = ?", tenantID, StatusDraft).Count(&draftProducts)
	stats["draft_products"] = draftProducts

	// Out of stock products
	var outOfStock int64
	r.db.Model(&Product{}).Where("tenant_id = ? AND track_quantity = true AND inventory_quantity <= 0", tenantID).Count(&outOfStock)
	stats["out_of_stock"] = outOfStock

	// Low stock products (< 10)
	var lowStock int64
	r.db.Model(&Product{}).Where("tenant_id = ? AND track_quantity = true AND inventory_quantity > 0 AND inventory_quantity < 10", tenantID).Count(&lowStock)
	stats["low_stock"] = lowStock

	// Total categories
	var totalCategories int64
	r.db.Model(&Category{}).Where("tenant_id = ?", tenantID).Count(&totalCategories)
	stats["total_categories"] = totalCategories

	// Total inventory value
	var totalValue float64
	r.db.Model(&Product{}).Where("tenant_id = ?", tenantID).Select("COALESCE(SUM(price * inventory_quantity), 0)").Scan(&totalValue)
	stats["total_value"] = totalValue

	return stats, nil
}

// Search operations

// SearchProducts performs full-text search on products
func (r *Repository) SearchProducts(tenantID uuid.UUID, query string, offset, limit int) ([]*Product, int64, error) {
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

// TODO: Add more repository methods
// - GetPopularProducts(tenantID uuid.UUID, limit int) ([]*Product, error)
// - GetNewProducts(tenantID uuid.UUID, days int, limit int) ([]*Product, error)
// - GetFeaturedProducts(tenantID uuid.UUID, limit int) ([]*Product, error)
// - GetRelatedProducts(productID uuid.UUID, limit int) ([]*Product, error)
// - BulkImportProducts(tenantID uuid.UUID, products []*Product) error
