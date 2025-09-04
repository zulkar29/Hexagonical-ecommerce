package discount

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the discount repository interface
type Repository interface {
	// Discount operations
	CreateDiscount(ctx context.Context, discount *Discount) error
	GetDiscountByID(ctx context.Context, tenantID, discountID uuid.UUID) (*Discount, error)
	GetDiscountByCode(ctx context.Context, tenantID uuid.UUID, code string) (*Discount, error)
	GetDiscounts(ctx context.Context, tenantID uuid.UUID, filter DiscountFilter) ([]Discount, error)
	UpdateDiscount(ctx context.Context, tenantID, discountID uuid.UUID, updates map[string]interface{}) error
	DeleteDiscount(ctx context.Context, tenantID, discountID uuid.UUID) error
	IncrementUsageCount(ctx context.Context, tenantID, discountID uuid.UUID) error
	DecrementUsageCount(ctx context.Context, tenantID, discountID uuid.UUID) error
	
	// Discount usage operations
	CreateDiscountUsage(ctx context.Context, usage *DiscountUsage) error
	GetDiscountUsage(ctx context.Context, tenantID, discountID uuid.UUID, filter UsageFilter) ([]DiscountUsage, error)
	GetCustomerDiscountUsageCount(ctx context.Context, tenantID uuid.UUID, customerEmail string, discountID uuid.UUID) (int, error)
	DeleteDiscountUsage(ctx context.Context, tenantID uuid.UUID, orderID uuid.UUID) error
	
	// Gift card operations
	CreateGiftCard(ctx context.Context, giftCard *GiftCard) error
	GetGiftCardByID(ctx context.Context, tenantID, giftCardID uuid.UUID) (*GiftCard, error)
	GetGiftCardByCode(ctx context.Context, tenantID uuid.UUID, code string) (*GiftCard, error)
	GetGiftCards(ctx context.Context, tenantID uuid.UUID, filter GiftCardFilter) ([]GiftCard, error)
	UpdateGiftCard(ctx context.Context, tenantID, giftCardID uuid.UUID, updates map[string]interface{}) error
	DeleteGiftCard(ctx context.Context, tenantID, giftCardID uuid.UUID) error
	
	// Gift card transaction operations
	CreateGiftCardTransaction(ctx context.Context, transaction *GiftCardTransaction) error
	GetGiftCardTransactions(ctx context.Context, tenantID, giftCardID uuid.UUID) ([]GiftCardTransaction, error)
	
	// Store credit operations
	CreateStoreCredit(ctx context.Context, storeCredit *StoreCredit) error
	GetStoreCredit(ctx context.Context, tenantID, customerID uuid.UUID) (*StoreCredit, error)
	UpdateStoreCredit(ctx context.Context, tenantID, customerID uuid.UUID, updates map[string]interface{}) error
	
	// Store credit transaction operations
	CreateStoreCreditTransaction(ctx context.Context, transaction *StoreCreditTransaction) error
	GetStoreCreditTransactions(ctx context.Context, tenantID, customerID uuid.UUID, filter StoreCreditFilter) ([]StoreCreditTransaction, error)
	
	// Analytics operations
	GetDiscountUsageStats(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*DiscountStats, error)
	GetDiscountPerformance(ctx context.Context, tenantID uuid.UUID, limit int) ([]DiscountPerformance, error)
	GetRevenueImpact(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*RevenueImpact, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new discount repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Discount operations
func (r *repository) CreateDiscount(ctx context.Context, discount *Discount) error {
	return r.db.WithContext(ctx).Create(discount).Error
}

func (r *repository) GetDiscountByID(ctx context.Context, tenantID, discountID uuid.UUID) (*Discount, error) {
	var discount Discount
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", discountID, tenantID).
		Preload("Usages").
		First(&discount).Error
	return &discount, err
}

func (r *repository) GetDiscountByCode(ctx context.Context, tenantID uuid.UUID, code string) (*Discount, error) {
	var discount Discount
	err := r.db.WithContext(ctx).
		Where("code = ? AND tenant_id = ?", code, tenantID).
		Preload("Usages").
		First(&discount).Error
	return &discount, err
}

func (r *repository) GetDiscounts(ctx context.Context, tenantID uuid.UUID, filter DiscountFilter) ([]Discount, error) {
	var discounts []Discount
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}
	
	if len(filter.Type) > 0 {
		query = query.Where("type IN ?", filter.Type)
	}
	
	if len(filter.Target) > 0 {
		query = query.Where("target IN ?", filter.Target)
	}
	
	if filter.Search != "" {
		query = query.Where("code ILIKE ? OR title ILIKE ? OR description ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	if filter.IsExpired != nil {
		now := time.Now()
		if *filter.IsExpired {
			query = query.Where("expires_at IS NOT NULL AND expires_at < ?", now)
		} else {
			query = query.Where("expires_at IS NULL OR expires_at >= ?", now)
		}
	}
	
	if filter.IsActive != nil {
		if *filter.IsActive {
			now := time.Now()
			query = query.Where("status = ? AND (starts_at IS NULL OR starts_at <= ?) AND (expires_at IS NULL OR expires_at >= ?)", 
				StatusActive, now, now)
		} else {
			query = query.Where("status != ?", StatusActive)
		}
	}
	
	if filter.CreatedBy != nil {
		query = query.Where("created_by = ?", *filter.CreatedBy)
	}
	
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}
	
	// Sorting
	sortBy := "created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "code", "title", "usage_count", "value", "created_at", "expires_at":
			sortBy = filter.SortBy
		}
	}
	
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	
	query = query.Order(sortBy + " " + sortOrder)
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Find(&discounts).Error
	return discounts, err
}

func (r *repository) UpdateDiscount(ctx context.Context, tenantID, discountID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&Discount{}).
		Where("id = ? AND tenant_id = ?", discountID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteDiscount(ctx context.Context, tenantID, discountID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete related usage records
		if err := tx.Where("discount_id = ? AND tenant_id = ?", discountID, tenantID).
			Delete(&DiscountUsage{}).Error; err != nil {
			return err
		}
		
		// Delete the discount
		return tx.Where("id = ? AND tenant_id = ?", discountID, tenantID).
			Delete(&Discount{}).Error
	})
}

func (r *repository) IncrementUsageCount(ctx context.Context, tenantID, discountID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&Discount{}).
		Where("id = ? AND tenant_id = ?", discountID, tenantID).
		Update("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

func (r *repository) DecrementUsageCount(ctx context.Context, tenantID, discountID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&Discount{}).
		Where("id = ? AND tenant_id = ?", discountID, tenantID).
		Update("usage_count", gorm.Expr("CASE WHEN usage_count > 0 THEN usage_count - 1 ELSE 0 END")).Error
}

// Discount usage operations
func (r *repository) CreateDiscountUsage(ctx context.Context, usage *DiscountUsage) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

func (r *repository) GetDiscountUsage(ctx context.Context, tenantID, discountID uuid.UUID, filter UsageFilter) ([]DiscountUsage, error) {
	var usages []DiscountUsage
	query := r.db.WithContext(ctx).
		Where("discount_id = ? AND tenant_id = ?", discountID, tenantID)
	
	// Apply filters
	if filter.StartDate != nil {
		query = query.Where("used_at >= ?", filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("used_at <= ?", filter.EndDate)
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Order("used_at DESC").Find(&usages).Error
	return usages, err
}

func (r *repository) GetCustomerDiscountUsageCount(ctx context.Context, tenantID uuid.UUID, customerEmail string, discountID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&DiscountUsage{}).
		Where("discount_id = ? AND tenant_id = ? AND customer_email = ?", discountID, tenantID, customerEmail).
		Count(&count).Error
	return int(count), err
}

func (r *repository) DeleteDiscountUsage(ctx context.Context, tenantID uuid.UUID, orderID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("order_id = ? AND tenant_id = ?", orderID, tenantID).
		Delete(&DiscountUsage{}).Error
}

// Gift card operations
func (r *repository) CreateGiftCard(ctx context.Context, giftCard *GiftCard) error {
	return r.db.WithContext(ctx).Create(giftCard).Error
}

func (r *repository) GetGiftCardByID(ctx context.Context, tenantID, giftCardID uuid.UUID) (*GiftCard, error) {
	var giftCard GiftCard
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", giftCardID, tenantID).
		Preload("Transactions").
		First(&giftCard).Error
	return &giftCard, err
}

func (r *repository) GetGiftCardByCode(ctx context.Context, tenantID uuid.UUID, code string) (*GiftCard, error) {
	var giftCard GiftCard
	err := r.db.WithContext(ctx).
		Where("code = ? AND tenant_id = ?", code, tenantID).
		Preload("Transactions").
		First(&giftCard).Error
	return &giftCard, err
}

func (r *repository) GetGiftCards(ctx context.Context, tenantID uuid.UUID, filter GiftCardFilter) ([]GiftCard, error) {
	var giftCards []GiftCard
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}
	
	if filter.Search != "" {
		query = query.Where("code ILIKE ? OR recipient_name ILIKE ? OR recipient_email ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	if filter.IsExpired != nil {
		now := time.Now()
		if *filter.IsExpired {
			query = query.Where("expires_at IS NOT NULL AND expires_at < ?", now)
		} else {
			query = query.Where("expires_at IS NULL OR expires_at >= ?", now)
		}
	}
	
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Order("created_at DESC").Find(&giftCards).Error
	return giftCards, err
}

func (r *repository) UpdateGiftCard(ctx context.Context, tenantID, giftCardID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&GiftCard{}).
		Where("id = ? AND tenant_id = ?", giftCardID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteGiftCard(ctx context.Context, tenantID, giftCardID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete related transactions
		if err := tx.Where("gift_card_id = ? AND tenant_id = ?", giftCardID, tenantID).
			Delete(&GiftCardTransaction{}).Error; err != nil {
			return err
		}
		
		// Delete the gift card
		return tx.Where("id = ? AND tenant_id = ?", giftCardID, tenantID).
			Delete(&GiftCard{}).Error
	})
}

// Gift card transaction operations
func (r *repository) CreateGiftCardTransaction(ctx context.Context, transaction *GiftCardTransaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *repository) GetGiftCardTransactions(ctx context.Context, tenantID, giftCardID uuid.UUID) ([]GiftCardTransaction, error) {
	var transactions []GiftCardTransaction
	err := r.db.WithContext(ctx).
		Where("gift_card_id = ? AND tenant_id = ?", giftCardID, tenantID).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

// Store credit operations
func (r *repository) CreateStoreCredit(ctx context.Context, storeCredit *StoreCredit) error {
	return r.db.WithContext(ctx).Create(storeCredit).Error
}

func (r *repository) GetStoreCredit(ctx context.Context, tenantID, customerID uuid.UUID) (*StoreCredit, error) {
	var storeCredit StoreCredit
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Preload("Transactions").
		First(&storeCredit).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create default store credit record
		storeCredit = StoreCredit{
			ID:             uuid.New(),
			TenantID:       tenantID,
			CustomerID:     customerID,
			CurrentBalance: 0,
			Currency:       "BDT", // TODO: Get from tenant settings
		}
		if createErr := r.CreateStoreCredit(ctx, &storeCredit); createErr != nil {
			return nil, createErr
		}
	}
	
	return &storeCredit, err
}

func (r *repository) UpdateStoreCredit(ctx context.Context, tenantID, customerID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&StoreCredit{}).
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Updates(updates).Error
}

// Store credit transaction operations
func (r *repository) CreateStoreCreditTransaction(ctx context.Context, transaction *StoreCreditTransaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *repository) GetStoreCreditTransactions(ctx context.Context, tenantID, customerID uuid.UUID, filter StoreCreditFilter) ([]StoreCreditTransaction, error) {
	var transactions []StoreCreditTransaction
	
	// First get the store credit ID
	var storeCredit StoreCredit
	err := r.db.WithContext(ctx).
		Select("id").
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		First(&storeCredit).Error
	
	if err != nil {
		return transactions, err
	}
	
	query := r.db.WithContext(ctx).
		Where("store_credit_id = ? AND tenant_id = ?", storeCredit.ID, tenantID)
	
	// Apply filters
	if len(filter.Type) > 0 {
		query = query.Where("type IN ?", filter.Type)
	}
	
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err = query.Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

// Analytics operations
func (r *repository) GetDiscountUsageStats(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*DiscountStats, error) {
	// TODO: Implement comprehensive discount analytics
	return nil, fmt.Errorf("TODO: implement GetDiscountUsageStats")
}

func (r *repository) GetDiscountPerformance(ctx context.Context, tenantID uuid.UUID, limit int) ([]DiscountPerformance, error) {
	var performance []DiscountPerformance
	
	// Query to get top performing discounts
	err := r.db.WithContext(ctx).
		Table("discounts d").
		Select(`
			d.id as discount_id,
			d.code,
			d.title,
			d.usage_count,
			COALESCE(SUM(du.discount_amount), 0) as total_savings,
			COALESCE(SUM(du.order_amount), 0) as revenue
		`).
		Joins("LEFT JOIN discount_usages du ON d.id = du.discount_id").
		Where("d.tenant_id = ?", tenantID).
		Group("d.id, d.code, d.title, d.usage_count").
		Order("d.usage_count DESC").
		Limit(limit).
		Find(&performance).Error
	
	// Calculate conversion rates (would need additional data)
	for i := range performance {
		performance[i].ConversionRate = 0 // TODO: Calculate actual conversion rate
	}
	
	return performance, err
}

func (r *repository) GetRevenueImpact(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*RevenueImpact, error) {
	// TODO: Implement revenue impact analysis with order data
	return nil, fmt.Errorf("TODO: implement GetRevenueImpact")
}