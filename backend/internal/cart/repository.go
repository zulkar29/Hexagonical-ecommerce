package cart

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the cart repository interface
type Repository interface {
	// Cart operations
	SaveCart(cart *Cart) (*Cart, error)
	FindCartByID(tenantID, cartID uuid.UUID) (*Cart, error)
	FindCartByCustomerID(tenantID, customerID uuid.UUID) (*Cart, error)
	FindCartBySessionID(tenantID uuid.UUID, sessionID string) (*Cart, error)
	UpdateCart(cart *Cart) (*Cart, error)
	DeleteCart(tenantID, cartID uuid.UUID) error
	ListCarts(tenantID uuid.UUID, filter CartListFilter, offset, limit int) ([]*Cart, int64, error)
	GetAbandonedCarts(tenantID uuid.UUID, since time.Time) ([]*Cart, error)
	GetExpiredCarts(tenantID uuid.UUID) ([]*Cart, error)
	CleanupExpiredCarts(tenantID uuid.UUID) error
	MergeGuestCartToCustomer(tenantID uuid.UUID, sessionID string, customerID uuid.UUID) error
	
	// Cart item operations
	AddCartItem(item *CartItem) (*CartItem, error)
	UpdateCartItem(item *CartItem) (*CartItem, error)
	RemoveCartItem(tenantID, cartID, itemID uuid.UUID) error
	FindCartItem(tenantID, cartID, itemID uuid.UUID) (*CartItem, error)
	ClearCartItems(tenantID, cartID uuid.UUID) error
	
	// Statistics and analytics
	GetCartStats(tenantID uuid.UUID) (*CartStats, error)
	GetAbandonmentRate(tenantID uuid.UUID, days int) (float64, error)
	GetAverageCartValue(tenantID uuid.UUID, days int) (float64, error)
	GetTopAbandonedProducts(tenantID uuid.UUID, limit int) ([]*AbandonedProductStats, error)
}

// CartListFilter represents filters for listing carts
type CartListFilter struct {
	Status     CartStatus `json:"status,omitempty"`
	CustomerID *uuid.UUID `json:"customer_id,omitempty"`
	MinTotal   *float64   `json:"min_total,omitempty"`
	MaxTotal   *float64   `json:"max_total,omitempty"`
	CreatedAfter  *time.Time `json:"created_after,omitempty"`
	CreatedBefore *time.Time `json:"created_before,omitempty"`
	HasCoupon     *bool      `json:"has_coupon,omitempty"`
}

// CartStats represents cart statistics
type CartStats struct {
	TotalCarts       int64   `json:"total_carts"`
	ActiveCarts      int64   `json:"active_carts"`
	AbandonedCarts   int64   `json:"abandoned_carts"`
	ConvertedCarts   int64   `json:"converted_carts"`
	AverageCartValue float64 `json:"average_cart_value"`
	AbandonmentRate  float64 `json:"abandonment_rate"`
	TotalRevenue     float64 `json:"total_revenue"`
}

// AbandonedProductStats represents statistics for abandoned products
type AbandonedProductStats struct {
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductSlug  string    `json:"product_slug"`
	AbandonCount int64     `json:"abandon_count"`
	TotalValue   float64   `json:"total_value"`
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new cart repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Cart operations

// SaveCart creates a new cart
func (r *repository) SaveCart(cart *Cart) (*Cart, error) {
	if err := r.db.Create(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

// FindCartByID retrieves a cart by ID with items
func (r *repository) FindCartByID(tenantID, cartID uuid.UUID) (*Cart, error) {
	var cart Cart
	err := r.db.Preload("Items").
		First(&cart, "id = ? AND tenant_id = ?", cartID, tenantID).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// FindCartByCustomerID retrieves active cart for a customer
func (r *repository) FindCartByCustomerID(tenantID, customerID uuid.UUID) (*Cart, error) {
	var cart Cart
	err := r.db.Preload("Items").
		Where("tenant_id = ? AND customer_id = ? AND status = ?", tenantID, customerID, StatusActive).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// FindCartBySessionID retrieves cart for a guest session
func (r *repository) FindCartBySessionID(tenantID uuid.UUID, sessionID string) (*Cart, error) {
	var cart Cart
	err := r.db.Preload("Items").
		Where("tenant_id = ? AND session_id = ? AND status = ?", tenantID, sessionID, StatusActive).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// UpdateCart updates an existing cart
func (r *repository) UpdateCart(cart *Cart) (*Cart, error) {
	if err := r.db.Save(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

// DeleteCart soft deletes a cart
func (r *repository) DeleteCart(tenantID, cartID uuid.UUID) error {
	return r.db.Where("tenant_id = ?", tenantID).Delete(&Cart{}, cartID).Error
}

// ListCarts returns paginated carts with filters
func (r *repository) ListCarts(tenantID uuid.UUID, filter CartListFilter, offset, limit int) ([]*Cart, int64, error) {
	var carts []*Cart
	var total int64

	query := r.db.Model(&Cart{}).Where("tenant_id = ?", tenantID)

	// Apply filters
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.MinTotal != nil {
		query = query.Where("total >= ?", *filter.MinTotal)
	}
	if filter.MaxTotal != nil {
		query = query.Where("total <= ?", *filter.MaxTotal)
	}
	if filter.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filter.CreatedAfter)
	}
	if filter.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filter.CreatedBefore)
	}
	if filter.HasCoupon != nil {
		if *filter.HasCoupon {
			query = query.Where("coupon_code IS NOT NULL AND coupon_code != ''")
		} else {
			query = query.Where("coupon_code IS NULL OR coupon_code = ''")
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with preloads
	if err := query.Preload("Items").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&carts).Error; err != nil {
		return nil, 0, err
	}

	return carts, total, nil
}

// GetAbandonedCarts returns carts abandoned since a specific time
func (r *repository) GetAbandonedCarts(tenantID uuid.UUID, since time.Time) ([]*Cart, error) {
	var carts []*Cart
	err := r.db.Preload("Items").
		Where("tenant_id = ? AND status = ? AND abandoned_at >= ?", tenantID, StatusAbandoned, since).
		Find(&carts).Error
	return carts, err
}

// GetExpiredCarts returns expired carts
func (r *repository) GetExpiredCarts(tenantID uuid.UUID) ([]*Cart, error) {
	var carts []*Cart
	now := time.Now()
	err := r.db.Where("tenant_id = ? AND expires_at IS NOT NULL AND expires_at <= ? AND status = ?", 
		tenantID, now, StatusActive).Find(&carts).Error
	return carts, err
}

// CleanupExpiredCarts marks expired carts as expired
func (r *repository) CleanupExpiredCarts(tenantID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&Cart{}).
		Where("tenant_id = ? AND expires_at IS NOT NULL AND expires_at <= ? AND status = ?", 
			tenantID, now, StatusActive).
		Update("status", StatusExpired).Error
}

// MergeGuestCartToCustomer merges guest cart to customer cart
func (r *repository) MergeGuestCartToCustomer(tenantID uuid.UUID, sessionID string, customerID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Find guest cart
		var guestCart Cart
		if err := tx.Where("tenant_id = ? AND session_id = ? AND status = ?", 
			tenantID, sessionID, StatusActive).First(&guestCart).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil // No guest cart to merge
			}
			return err
		}

		// Find or create customer cart
		var customerCart Cart
		err := tx.Where("tenant_id = ? AND customer_id = ? AND status = ?", 
			tenantID, customerID, StatusActive).First(&customerCart).Error
		
		if err == gorm.ErrRecordNotFound {
			// No existing customer cart, convert guest cart
			guestCart.CustomerID = &customerID
			guestCart.SessionID = ""
			return tx.Save(&guestCart).Error
		} else if err != nil {
			return err
		}

		// Merge guest cart items into customer cart
		var guestItems []CartItem
		if err := tx.Where("cart_id = ?", guestCart.ID).Find(&guestItems).Error; err != nil {
			return err
		}

		for _, item := range guestItems {
			// Check if item already exists in customer cart
			var existingItem CartItem
			err := tx.Where("cart_id = ? AND product_id = ? AND variant_id = ?", 
				customerCart.ID, item.ProductID, item.VariantID).First(&existingItem).Error
			
			if err == gorm.ErrRecordNotFound {
				// Item doesn't exist, move it to customer cart
				item.CartID = customerCart.ID
				if err := tx.Save(&item).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				// Item exists, merge quantities
				existingItem.Quantity += item.Quantity
				existingItem.CalculateLineTotal()
				if err := tx.Save(&existingItem).Error; err != nil {
					return err
				}
				// Delete the guest item
				if err := tx.Delete(&item).Error; err != nil {
					return err
				}
			}
		}

		// Delete guest cart
		return tx.Delete(&guestCart).Error
	})
}

// Cart item operations

// AddCartItem adds an item to cart
func (r *repository) AddCartItem(item *CartItem) (*CartItem, error) {
	if err := r.db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// UpdateCartItem updates a cart item
func (r *repository) UpdateCartItem(item *CartItem) (*CartItem, error) {
	if err := r.db.Save(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// RemoveCartItem removes an item from cart
func (r *repository) RemoveCartItem(tenantID, cartID, itemID uuid.UUID) error {
	return r.db.Where("id = ? AND cart_id IN (SELECT id FROM carts WHERE id = ? AND tenant_id = ?)", 
		itemID, cartID, tenantID).Delete(&CartItem{}).Error
}

// FindCartItem finds a specific cart item
func (r *repository) FindCartItem(tenantID, cartID, itemID uuid.UUID) (*CartItem, error) {
	var item CartItem
	err := r.db.Where("id = ? AND cart_id IN (SELECT id FROM carts WHERE id = ? AND tenant_id = ?)", 
		itemID, cartID, tenantID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ClearCartItems removes all items from a cart
func (r *repository) ClearCartItems(tenantID, cartID uuid.UUID) error {
	return r.db.Where("cart_id IN (SELECT id FROM carts WHERE id = ? AND tenant_id = ?)", 
		cartID, tenantID).Delete(&CartItem{}).Error
}

// Statistics and analytics

// GetCartStats returns cart statistics
func (r *repository) GetCartStats(tenantID uuid.UUID) (*CartStats, error) {
	stats := &CartStats{}
	
	// Total carts
	r.db.Model(&Cart{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalCarts)
	
	// Active carts
	r.db.Model(&Cart{}).Where("tenant_id = ? AND status = ?", tenantID, StatusActive).Count(&stats.ActiveCarts)
	
	// Abandoned carts
	r.db.Model(&Cart{}).Where("tenant_id = ? AND status = ?", tenantID, StatusAbandoned).Count(&stats.AbandonedCarts)
	
	// Converted carts
	r.db.Model(&Cart{}).Where("tenant_id = ? AND status = ?", tenantID, StatusConverted).Count(&stats.ConvertedCarts)
	
	// Average cart value
	var avgValue sql.NullFloat64
	r.db.Model(&Cart{}).Where("tenant_id = ? AND status = ?", tenantID, StatusActive).
		Select("AVG(total)").Scan(&avgValue)
	if avgValue.Valid {
		stats.AverageCartValue = avgValue.Float64
	}
	
	// Abandonment rate
	if stats.TotalCarts > 0 {
		stats.AbandonmentRate = float64(stats.AbandonedCarts) / float64(stats.TotalCarts) * 100
	}
	
	// Total revenue from converted carts
	var totalRevenue sql.NullFloat64
	r.db.Model(&Cart{}).Where("tenant_id = ? AND status = ?", tenantID, StatusConverted).
		Select("SUM(total)").Scan(&totalRevenue)
	if totalRevenue.Valid {
		stats.TotalRevenue = totalRevenue.Float64
	}
	
	return stats, nil
}

// GetAbandonmentRate calculates abandonment rate for the last N days
func (r *repository) GetAbandonmentRate(tenantID uuid.UUID, days int) (float64, error) {
	since := time.Now().AddDate(0, 0, -days)
	
	var totalCarts, abandonedCarts int64
	
	r.db.Model(&Cart{}).Where("tenant_id = ? AND created_at >= ?", tenantID, since).Count(&totalCarts)
	r.db.Model(&Cart{}).Where("tenant_id = ? AND status = ? AND created_at >= ?", 
		tenantID, StatusAbandoned, since).Count(&abandonedCarts)
	
	if totalCarts == 0 {
		return 0, nil
	}
	
	return float64(abandonedCarts) / float64(totalCarts) * 100, nil
}

// GetAverageCartValue calculates average cart value for the last N days
func (r *repository) GetAverageCartValue(tenantID uuid.UUID, days int) (float64, error) {
	since := time.Now().AddDate(0, 0, -days)
	
	var avgValue sql.NullFloat64
	r.db.Model(&Cart{}).Where("tenant_id = ? AND created_at >= ?", tenantID, since).
		Select("AVG(total)").Scan(&avgValue)
	
	if avgValue.Valid {
		return avgValue.Float64, nil
	}
	return 0, nil
}

// GetTopAbandonedProducts returns most abandoned products
func (r *repository) GetTopAbandonedProducts(tenantID uuid.UUID, limit int) ([]*AbandonedProductStats, error) {
	var stats []*AbandonedProductStats
	
	err := r.db.Table("cart_items ci").
		Select("ci.product_id, ci.product_name, ci.product_slug, COUNT(*) as abandon_count, SUM(ci.line_total) as total_value").
		Joins("JOIN carts c ON ci.cart_id = c.id").
		Where("c.tenant_id = ? AND c.status = ?", tenantID, StatusAbandoned).
		Group("ci.product_id, ci.product_name, ci.product_slug").
		Order("abandon_count DESC").
		Limit(limit).
		Scan(&stats).Error
	
	return stats, err
}