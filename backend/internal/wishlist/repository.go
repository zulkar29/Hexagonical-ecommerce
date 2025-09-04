package wishlist

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for wishlist data operations
type Repository interface {
	// Wishlist operations
	Save(ctx context.Context, wishlist *Wishlist) error
	FindByID(ctx context.Context, tenantID, wishlistID uuid.UUID) (*Wishlist, error)
	FindByCustomerID(ctx context.Context, tenantID, customerID uuid.UUID) ([]Wishlist, error)
	FindDefaultByCustomerID(ctx context.Context, tenantID, customerID uuid.UUID) (*Wishlist, error)
	FindByShareToken(ctx context.Context, shareToken string) (*Wishlist, error)
	Update(ctx context.Context, wishlist *Wishlist) error
	Delete(ctx context.Context, tenantID, wishlistID uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, filter WishlistFilter, limit, offset int) ([]Wishlist, error)
	Count(ctx context.Context, tenantID uuid.UUID, filter WishlistFilter) (int64, error)
	
	// Wishlist item operations
	SaveItem(ctx context.Context, item *WishlistItem) error
	FindItemByID(ctx context.Context, tenantID, itemID uuid.UUID) (*WishlistItem, error)
	FindItemByProductAndWishlist(ctx context.Context, tenantID, wishlistID, productID uuid.UUID, variantID *uuid.UUID) (*WishlistItem, error)
	UpdateItem(ctx context.Context, item *WishlistItem) error
	DeleteItem(ctx context.Context, tenantID, itemID uuid.UUID) error
	ListItems(ctx context.Context, tenantID uuid.UUID, filter WishlistItemFilter, limit, offset int) ([]WishlistItem, error)
	CountItems(ctx context.Context, tenantID uuid.UUID, filter WishlistItemFilter) (int64, error)
	
	// Wishlist management
	SetDefaultWishlist(ctx context.Context, tenantID, customerID, wishlistID uuid.UUID) error
	ClearItems(ctx context.Context, tenantID, wishlistID uuid.UUID) error
	MoveItem(ctx context.Context, tenantID, itemID, targetWishlistID uuid.UUID) error
	CopyItem(ctx context.Context, tenantID, itemID, targetWishlistID uuid.UUID) error
	MergeWishlists(ctx context.Context, tenantID, sourceWishlistID, targetWishlistID uuid.UUID) error
	
	// Bulk operations
	BulkAddItems(ctx context.Context, items []WishlistItem) error
	BulkDeleteItems(ctx context.Context, tenantID uuid.UUID, itemIDs []uuid.UUID) error
	BulkUpdateItemPriority(ctx context.Context, tenantID uuid.UUID, updates map[uuid.UUID]int) error
	
	// Validation and checks
	ExistsByName(ctx context.Context, tenantID, customerID uuid.UUID, name string, excludeID *uuid.UUID) (bool, error)
	ExistsByShareToken(ctx context.Context, shareToken string, excludeID *uuid.UUID) (bool, error)
	CountWishlistsByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) (int64, error)
	CountItemsByWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID) (int64, error)
	
	// Statistics and analytics
	GetStats(ctx context.Context, tenantID uuid.UUID) (*WishlistStats, error)
	GetMostWishedProducts(ctx context.Context, tenantID uuid.UUID, limit int) ([]ProductWishCount, error)
	GetCustomerWishlistActivity(ctx context.Context, tenantID, customerID uuid.UUID, days int) ([]WishlistActivity, error)
	GetPopularWishlists(ctx context.Context, tenantID uuid.UUID, limit int) ([]Wishlist, error)
	
	// Cleanup operations
	CleanupEmptyWishlists(ctx context.Context, tenantID uuid.UUID, olderThanDays int) (int64, error)
	CleanupOrphanedItems(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// GormRepository implements Repository using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM-based repository
func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// WishlistActivity represents wishlist activity data
type WishlistActivity struct {
	Date       string `json:"date"`
	ItemsAdded int64  `json:"items_added"`
	ItemsRemoved int64 `json:"items_removed"`
}

// Wishlist operations

// Save creates or updates a wishlist
func (r *GormRepository) Save(ctx context.Context, wishlist *Wishlist) error {
	return r.db.WithContext(ctx).Create(wishlist).Error
}

// FindByID finds a wishlist by ID
func (r *GormRepository) FindByID(ctx context.Context, tenantID, wishlistID uuid.UUID) (*Wishlist, error) {
	var wishlist Wishlist
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Variant").
		Where("tenant_id = ? AND id = ?", tenantID, wishlistID).
		First(&wishlist).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}
	
	return &wishlist, nil
}

// FindByCustomerID finds all wishlists for a customer
func (r *GormRepository) FindByCustomerID(ctx context.Context, tenantID, customerID uuid.UUID) ([]Wishlist, error) {
	var wishlists []Wishlist
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Order("is_default DESC, created_at ASC").
		Find(&wishlists).Error
	
	return wishlists, err
}

// FindDefaultByCustomerID finds the default wishlist for a customer
func (r *GormRepository) FindDefaultByCustomerID(ctx context.Context, tenantID, customerID uuid.UUID) (*Wishlist, error) {
	var wishlist Wishlist
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND customer_id = ? AND is_default = ?", tenantID, customerID, true).
		First(&wishlist).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}
	
	return &wishlist, nil
}

// FindByShareToken finds a wishlist by share token
func (r *GormRepository) FindByShareToken(ctx context.Context, shareToken string) (*Wishlist, error) {
	var wishlist Wishlist
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Variant").
		Where("share_token = ? AND is_public = ?", shareToken, true).
		First(&wishlist).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}
	
	return &wishlist, nil
}

// Update updates a wishlist
func (r *GormRepository) Update(ctx context.Context, wishlist *Wishlist) error {
	return r.db.WithContext(ctx).Save(wishlist).Error
}

// Delete deletes a wishlist
func (r *GormRepository) Delete(ctx context.Context, tenantID, wishlistID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, wishlistID).
		Delete(&Wishlist{}).Error
}

// List returns paginated wishlists with filters
func (r *GormRepository) List(ctx context.Context, tenantID uuid.UUID, filter WishlistFilter, limit, offset int) ([]Wishlist, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	query = r.applyWishlistFilters(query, filter)
	
	var wishlists []Wishlist
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&wishlists).Error
	return wishlists, err
}

// Count returns the total count of wishlists with filters
func (r *GormRepository) Count(ctx context.Context, tenantID uuid.UUID, filter WishlistFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&Wishlist{}).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	query = r.applyWishlistFilters(query, filter)
	
	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Wishlist item operations

// SaveItem creates a new wishlist item
func (r *GormRepository) SaveItem(ctx context.Context, item *WishlistItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the item
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		
		// Update wishlist item count
		return tx.Model(&Wishlist{}).
			Where("id = ?", item.WishlistID).
			Update("item_count", gorm.Expr("item_count + 1")).Error
	})
}

// FindItemByID finds a wishlist item by ID
func (r *GormRepository) FindItemByID(ctx context.Context, tenantID, itemID uuid.UUID) (*WishlistItem, error) {
	var item WishlistItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Variant").
		Where("tenant_id = ? AND id = ?", tenantID, itemID).
		First(&item).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrWishlistItemNotFound
		}
		return nil, err
	}
	
	return &item, nil
}

// FindItemByProductAndWishlist finds an item by product and wishlist
func (r *GormRepository) FindItemByProductAndWishlist(ctx context.Context, tenantID, wishlistID, productID uuid.UUID, variantID *uuid.UUID) (*WishlistItem, error) {
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND wishlist_id = ? AND product_id = ?", tenantID, wishlistID, productID)
	
	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}
	
	var item WishlistItem
	err := query.First(&item).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrWishlistItemNotFound
		}
		return nil, err
	}
	
	return &item, nil
}

// UpdateItem updates a wishlist item
func (r *GormRepository) UpdateItem(ctx context.Context, item *WishlistItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// DeleteItem deletes a wishlist item
func (r *GormRepository) DeleteItem(ctx context.Context, tenantID, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the item to find wishlist ID
		var item WishlistItem
		if err := tx.Where("tenant_id = ? AND id = ?", tenantID, itemID).First(&item).Error; err != nil {
			return err
		}
		
		// Delete the item
		if err := tx.Delete(&item).Error; err != nil {
			return err
		}
		
		// Update wishlist item count
		return tx.Model(&Wishlist{}).
			Where("id = ?", item.WishlistID).
			Update("item_count", gorm.Expr("item_count - 1")).Error
	})
}

// ListItems returns paginated wishlist items with filters
func (r *GormRepository) ListItems(ctx context.Context, tenantID uuid.UUID, filter WishlistItemFilter, limit, offset int) ([]WishlistItem, error) {
	query := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Variant").
		Where("tenant_id = ?", tenantID)
	
	// Apply filters
	query = r.applyWishlistItemFilters(query, filter)
	
	var items []WishlistItem
	err := query.Limit(limit).Offset(offset).Order("priority DESC, added_at DESC").Find(&items).Error
	return items, err
}

// CountItems returns the total count of wishlist items with filters
func (r *GormRepository) CountItems(ctx context.Context, tenantID uuid.UUID, filter WishlistItemFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&WishlistItem{}).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	query = r.applyWishlistItemFilters(query, filter)
	
	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Wishlist management operations

// SetDefaultWishlist sets a wishlist as default for a customer
func (r *GormRepository) SetDefaultWishlist(ctx context.Context, tenantID, customerID, wishlistID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Clear existing default
		if err := tx.Model(&Wishlist{}).
			Where("tenant_id = ? AND customer_id = ? AND is_default = ?", tenantID, customerID, true).
			Update("is_default", false).Error; err != nil {
			return err
		}
		
		// Set new default
		return tx.Model(&Wishlist{}).
			Where("tenant_id = ? AND id = ?", tenantID, wishlistID).
			Update("is_default", true).Error
	})
}

// ClearItems removes all items from a wishlist
func (r *GormRepository) ClearItems(ctx context.Context, tenantID, wishlistID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete all items
		if err := tx.Where("tenant_id = ? AND wishlist_id = ?", tenantID, wishlistID).
			Delete(&WishlistItem{}).Error; err != nil {
			return err
		}
		
		// Reset item count
		return tx.Model(&Wishlist{}).
			Where("tenant_id = ? AND id = ?", tenantID, wishlistID).
			Update("item_count", 0).Error
	})
}

// MoveItem moves an item to another wishlist
func (r *GormRepository) MoveItem(ctx context.Context, tenantID, itemID, targetWishlistID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the item
		var item WishlistItem
		if err := tx.Where("tenant_id = ? AND id = ?", tenantID, itemID).First(&item).Error; err != nil {
			return err
		}
		
		oldWishlistID := item.WishlistID
		
		// Update item wishlist
		if err := tx.Model(&item).Update("wishlist_id", targetWishlistID).Error; err != nil {
			return err
		}
		
		// Update old wishlist count
		if err := tx.Model(&Wishlist{}).
			Where("id = ?", oldWishlistID).
			Update("item_count", gorm.Expr("item_count - 1")).Error; err != nil {
			return err
		}
		
		// Update new wishlist count
		return tx.Model(&Wishlist{}).
			Where("id = ?", targetWishlistID).
			Update("item_count", gorm.Expr("item_count + 1")).Error
	})
}

// CopyItem copies an item to another wishlist
func (r *GormRepository) CopyItem(ctx context.Context, tenantID, itemID, targetWishlistID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the original item
		var originalItem WishlistItem
		if err := tx.Where("tenant_id = ? AND id = ?", tenantID, itemID).First(&originalItem).Error; err != nil {
			return err
		}
		
		// Create new item
		newItem := WishlistItem{
			ID:         uuid.New(),
			TenantID:   originalItem.TenantID,
			WishlistID: targetWishlistID,
			ProductID:  originalItem.ProductID,
			VariantID:  originalItem.VariantID,
			Quantity:   originalItem.Quantity,
			Notes:      originalItem.Notes,
			Priority:   originalItem.Priority,
		}
		
		// Save new item
		if err := tx.Create(&newItem).Error; err != nil {
			return err
		}
		
		// Update target wishlist count
		return tx.Model(&Wishlist{}).
			Where("id = ?", targetWishlistID).
			Update("item_count", gorm.Expr("item_count + 1")).Error
	})
}

// MergeWishlists merges source wishlist into target wishlist
func (r *GormRepository) MergeWishlists(ctx context.Context, tenantID, sourceWishlistID, targetWishlistID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Move all items from source to target
		if err := tx.Model(&WishlistItem{}).
			Where("tenant_id = ? AND wishlist_id = ?", tenantID, sourceWishlistID).
			Update("wishlist_id", targetWishlistID).Error; err != nil {
			return err
		}
		
		// Get item count from source
		var sourceWishlist Wishlist
		if err := tx.Where("tenant_id = ? AND id = ?", tenantID, sourceWishlistID).First(&sourceWishlist).Error; err != nil {
			return err
		}
		
		// Update target wishlist count
		if err := tx.Model(&Wishlist{}).
			Where("tenant_id = ? AND id = ?", tenantID, targetWishlistID).
			Update("item_count", gorm.Expr("item_count + ?", sourceWishlist.ItemCount)).Error; err != nil {
			return err
		}
		
		// Delete source wishlist
		return tx.Where("tenant_id = ? AND id = ?", tenantID, sourceWishlistID).Delete(&Wishlist{}).Error
	})
}

// Bulk operations

// BulkAddItems adds multiple items to wishlists
func (r *GormRepository) BulkAddItems(ctx context.Context, items []WishlistItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create all items
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		
		// Update wishlist counts
		wishlistCounts := make(map[uuid.UUID]int)
		for _, item := range items {
			wishlistCounts[item.WishlistID]++
		}
		
		for wishlistID, count := range wishlistCounts {
			if err := tx.Model(&Wishlist{}).
				Where("id = ?", wishlistID).
				Update("item_count", gorm.Expr("item_count + ?", count)).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

// BulkDeleteItems deletes multiple wishlist items
func (r *GormRepository) BulkDeleteItems(ctx context.Context, tenantID uuid.UUID, itemIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get items to find wishlist IDs
		var items []WishlistItem
		if err := tx.Where("tenant_id = ? AND id IN ?", tenantID, itemIDs).Find(&items).Error; err != nil {
			return err
		}
		
		// Count items per wishlist
		wishlistCounts := make(map[uuid.UUID]int)
		for _, item := range items {
			wishlistCounts[item.WishlistID]++
		}
		
		// Delete items
		if err := tx.Where("tenant_id = ? AND id IN ?", tenantID, itemIDs).Delete(&WishlistItem{}).Error; err != nil {
			return err
		}
		
		// Update wishlist counts
		for wishlistID, count := range wishlistCounts {
			if err := tx.Model(&Wishlist{}).
				Where("id = ?", wishlistID).
				Update("item_count", gorm.Expr("item_count - ?", count)).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

// BulkUpdateItemPriority updates priority for multiple items
func (r *GormRepository) BulkUpdateItemPriority(ctx context.Context, tenantID uuid.UUID, updates map[uuid.UUID]int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for itemID, priority := range updates {
			if err := tx.Model(&WishlistItem{}).
				Where("tenant_id = ? AND id = ?", tenantID, itemID).
				Update("priority", priority).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Validation and checks

// ExistsByName checks if a wishlist with the given name exists for a customer
func (r *GormRepository) ExistsByName(ctx context.Context, tenantID, customerID uuid.UUID, name string, excludeID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).Model(&Wishlist{}).
		Where("tenant_id = ? AND customer_id = ? AND name = ?", tenantID, customerID, name)
	
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// ExistsByShareToken checks if a wishlist with the given share token exists
func (r *GormRepository) ExistsByShareToken(ctx context.Context, shareToken string, excludeID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).Model(&Wishlist{}).Where("share_token = ?", shareToken)
	
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// CountWishlistsByCustomer counts wishlists for a customer
func (r *GormRepository) CountWishlistsByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Wishlist{}).
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Count(&count).Error
	return count, err
}

// CountItemsByWishlist counts items in a wishlist
func (r *GormRepository) CountItemsByWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&WishlistItem{}).
		Where("tenant_id = ? AND wishlist_id = ?", tenantID, wishlistID).
		Count(&count).Error
	return count, err
}

// Statistics and analytics

// GetStats returns wishlist statistics
func (r *GormRepository) GetStats(ctx context.Context, tenantID uuid.UUID) (*WishlistStats, error) {
	stats := &WishlistStats{}
	
	// Total wishlists
	r.db.WithContext(ctx).Model(&Wishlist{}).
		Where("tenant_id = ?", tenantID).
		Count(&stats.TotalWishlists)
	
	// Total items
	r.db.WithContext(ctx).Model(&WishlistItem{}).
		Where("tenant_id = ?", tenantID).
		Count(&stats.TotalItems)
	
	// Average items per wishlist
	if stats.TotalWishlists > 0 {
		stats.AverageItemsPerWishlist = float64(stats.TotalItems) / float64(stats.TotalWishlists)
	}
	
	// Public/Private wishlists
	r.db.WithContext(ctx).Model(&Wishlist{}).
		Where("tenant_id = ? AND is_public = ?", tenantID, true).
		Count(&stats.PublicWishlists)
	
	stats.PrivateWishlists = stats.TotalWishlists - stats.PublicWishlists
	
	// Empty wishlists
	r.db.WithContext(ctx).Model(&Wishlist{}).
		Where("tenant_id = ? AND item_count = 0", tenantID).
		Count(&stats.EmptyWishlists)
	
	// Most wished products
	mostWished, err := r.GetMostWishedProducts(ctx, tenantID, 10)
	if err != nil {
		return nil, err
	}
	stats.MostWishedProducts = mostWished
	
	return stats, nil
}

// GetMostWishedProducts returns the most wished products
func (r *GormRepository) GetMostWishedProducts(ctx context.Context, tenantID uuid.UUID, limit int) ([]ProductWishCount, error) {
	var results []ProductWishCount
	
	err := r.db.WithContext(ctx).
		Model(&WishlistItem{}).
		Select("product_id, COUNT(*) as wish_count").
		Where("tenant_id = ?", tenantID).
		Group("product_id").
		Order("wish_count DESC").
		Limit(limit).
		Scan(&results).Error
	
	return results, err
}

// GetCustomerWishlistActivity returns customer wishlist activity
func (r *GormRepository) GetCustomerWishlistActivity(ctx context.Context, tenantID, customerID uuid.UUID, days int) ([]WishlistActivity, error) {
	var results []WishlistActivity
	
	// This is a simplified implementation
	// In a real application, you might want to track activity in a separate table
	query := `
		SELECT 
			DATE(added_at) as date,
			COUNT(*) as items_added,
			0 as items_removed
		FROM wishlist_items wi
		JOIN wishlists w ON wi.wishlist_id = w.id
		WHERE wi.tenant_id = ? AND w.customer_id = ? AND wi.added_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY DATE(added_at)
		ORDER BY date DESC
	`
	
	err := r.db.WithContext(ctx).Raw(query, tenantID, customerID, days).Scan(&results).Error
	return results, err
}

// GetPopularWishlists returns popular public wishlists
func (r *GormRepository) GetPopularWishlists(ctx context.Context, tenantID uuid.UUID, limit int) ([]Wishlist, error) {
	var wishlists []Wishlist
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_public = ?", tenantID, true).
		Order("item_count DESC, created_at DESC").
		Limit(limit).
		Find(&wishlists).Error
	
	return wishlists, err
}

// Cleanup operations

// CleanupEmptyWishlists removes empty wishlists older than specified days
func (r *GormRepository) CleanupEmptyWishlists(ctx context.Context, tenantID uuid.UUID, olderThanDays int) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND item_count = 0 AND is_default = false AND created_at < DATE_SUB(NOW(), INTERVAL ? DAY)", tenantID, olderThanDays).
		Delete(&Wishlist{})
	
	return result.RowsAffected, result.Error
}

// CleanupOrphanedItems removes items for non-existent products
func (r *GormRepository) CleanupOrphanedItems(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	// This would require joining with the products table
	// For now, return 0 as this would need product service integration
	return 0, nil
}

// Helper methods

// applyWishlistFilters applies filters to wishlist queries
func (r *GormRepository) applyWishlistFilters(query *gorm.DB, filter WishlistFilter) *gorm.DB {
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	
	if filter.IsDefault != nil {
		query = query.Where("is_default = ?", *filter.IsDefault)
	}
	
	if filter.IsPublic != nil {
		query = query.Where("is_public = ?", *filter.IsPublic)
	}
	
	if filter.IsEmpty != nil {
		if *filter.IsEmpty {
			query = query.Where("item_count = 0")
		} else {
			query = query.Where("item_count > 0")
		}
	}
	
	return query
}

// applyWishlistItemFilters applies filters to wishlist item queries
func (r *GormRepository) applyWishlistItemFilters(query *gorm.DB, filter WishlistItemFilter) *gorm.DB {
	if filter.WishlistID != nil {
		query = query.Where("wishlist_id = ?", *filter.WishlistID)
	}
	
	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	
	if filter.VariantID != nil {
		query = query.Where("variant_id = ?", *filter.VariantID)
	}
	
	if filter.MinPriority != nil {
		query = query.Where("priority >= ?", *filter.MinPriority)
	}
	
	if filter.MaxPriority != nil {
		query = query.Where("priority <= ?", *filter.MaxPriority)
	}
	
	// Note: IsAvailable and HasDiscount filters would require joining with product tables
	// These would be better handled in the service layer
	
	return query
}