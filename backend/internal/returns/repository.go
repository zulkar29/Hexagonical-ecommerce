package returns

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for return data operations
type Repository interface {
	// Return operations
	CreateReturn(ctx context.Context, return_ *Return) error
	GetReturnByID(ctx context.Context, tenantID, returnID uuid.UUID) (*Return, error)
	GetReturnByNumber(ctx context.Context, tenantID uuid.UUID, returnNumber string) (*Return, error)
	ListReturns(ctx context.Context, tenantID uuid.UUID, filter ReturnFilter) ([]*Return, int64, error)
	UpdateReturn(ctx context.Context, return_ *Return) error
	DeleteReturn(ctx context.Context, tenantID, returnID uuid.UUID) error
	
	// Return item operations
	CreateReturnItem(ctx context.Context, item *ReturnItem) error
	GetReturnItemByID(ctx context.Context, itemID uuid.UUID) (*ReturnItem, error)
	GetReturnItems(ctx context.Context, returnID uuid.UUID) ([]*ReturnItem, error)
	UpdateReturnItem(ctx context.Context, item *ReturnItem) error
	DeleteReturnItem(ctx context.Context, itemID uuid.UUID) error
	
	// Return reason operations
	CreateReturnReason(ctx context.Context, reason *ReturnReason) error
	GetReturnReasonByID(ctx context.Context, tenantID, reasonID uuid.UUID) (*ReturnReason, error)
	ListReturnReasons(ctx context.Context, tenantID uuid.UUID, activeOnly bool) ([]*ReturnReason, error)
	UpdateReturnReason(ctx context.Context, reason *ReturnReason) error
	DeleteReturnReason(ctx context.Context, tenantID, reasonID uuid.UUID) error
	
	// Statistics and analytics
	GetReturnStats(ctx context.Context, tenantID uuid.UUID, filter StatsFilter) (*ReturnStats, error)
	GetReturnsByCustomer(ctx context.Context, tenantID, customerID uuid.UUID, limit int) ([]*Return, error)
	GetReturnsByOrder(ctx context.Context, tenantID, orderID uuid.UUID) ([]*Return, error)
}

// ReturnFilter represents filtering options for returns
type ReturnFilter struct {
	Status     []ReturnStatus `json:"status,omitempty"`
	Type       []ReturnType   `json:"type,omitempty"`
	CustomerID *uuid.UUID     `json:"customer_id,omitempty"`
	OrderID    *uuid.UUID     `json:"order_id,omitempty"`
	ReasonID   *uuid.UUID     `json:"reason_id,omitempty"`
	
	// Date filters
	CreatedAfter  *time.Time `json:"created_after,omitempty"`
	CreatedBefore *time.Time `json:"created_before,omitempty"`
	
	// Amount filters
	MinRefundAmount *float64 `json:"min_refund_amount,omitempty"`
	MaxRefundAmount *float64 `json:"max_refund_amount,omitempty"`
	
	// Search
	Search string `json:"search,omitempty"`
	
	// Pagination
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
	
	// Sorting
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// StatsFilter represents filtering options for return statistics
type StatsFilter struct {
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	GroupBy  string     `json:"group_by,omitempty"` // day, week, month, year
}

// ReturnStats represents return statistics
type ReturnStats struct {
	TotalReturns       int64   `json:"total_returns"`
	TotalRefundAmount  float64 `json:"total_refund_amount"`
	AverageRefundAmount float64 `json:"average_refund_amount"`
	
	// Status breakdown
	PendingReturns   int64 `json:"pending_returns"`
	ApprovedReturns  int64 `json:"approved_returns"`
	RejectedReturns  int64 `json:"rejected_returns"`
	ProcessingReturns int64 `json:"processing_returns"`
	CompletedReturns int64 `json:"completed_returns"`
	
	// Type breakdown
	RefundReturns   int64 `json:"refund_returns"`
	ExchangeReturns int64 `json:"exchange_returns"`
	
	// Time-based stats
	ReturnsByPeriod []PeriodStats `json:"returns_by_period,omitempty"`
	
	// Top reasons
	TopReasons []ReasonStats `json:"top_reasons,omitempty"`
}

// PeriodStats represents statistics for a specific time period
type PeriodStats struct {
	Period      string  `json:"period"`
	Count       int64   `json:"count"`
	RefundAmount float64 `json:"refund_amount"`
}

// ReasonStats represents statistics for return reasons
type ReasonStats struct {
	ReasonID    uuid.UUID `json:"reason_id"`
	ReasonName  string    `json:"reason_name"`
	Count       int64     `json:"count"`
	Percentage  float64   `json:"percentage"`
}

// gormRepository implements the Repository interface using GORM
type gormRepository struct {
	db *gorm.DB
}

// NewRepository creates a new return repository
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// CreateReturn creates a new return
func (r *gormRepository) CreateReturn(ctx context.Context, return_ *Return) error {
	return r.db.WithContext(ctx).Create(return_).Error
}

// GetReturnByID retrieves a return by ID
func (r *gormRepository) GetReturnByID(ctx context.Context, tenantID, returnID uuid.UUID) (*Return, error) {
	var return_ Return
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Reason").
		Where("tenant_id = ? AND id = ?", tenantID, returnID).
		First(&return_).Error
	
	if err != nil {
		return nil, err
	}
	return &return_, nil
}

// GetReturnByNumber retrieves a return by return number
func (r *gormRepository) GetReturnByNumber(ctx context.Context, tenantID uuid.UUID, returnNumber string) (*Return, error) {
	var return_ Return
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Reason").
		Where("tenant_id = ? AND return_number = ?", tenantID, returnNumber).
		First(&return_).Error
	
	if err != nil {
		return nil, err
	}
	return &return_, nil
}

// ListReturns retrieves returns with filtering and pagination
func (r *gormRepository) ListReturns(ctx context.Context, tenantID uuid.UUID, filter ReturnFilter) ([]*Return, int64, error) {
	query := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Reason").
		Where("tenant_id = ?", tenantID)
	
	// Apply filters
	query = r.applyReturnFilters(query, filter)
	
	// Count total records
	var total int64
	countQuery := r.db.WithContext(ctx).Model(&Return{}).Where("tenant_id = ?", tenantID)
	countQuery = r.applyReturnFilters(countQuery, filter)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply sorting
	if filter.SortBy != "" {
		order := "ASC"
		if filter.SortOrder == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", filter.SortBy, order))
	} else {
		query = query.Order("created_at DESC")
	}
	
	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	var returns []*Return
	err := query.Find(&returns).Error
	return returns, total, err
}

// UpdateReturn updates a return
func (r *gormRepository) UpdateReturn(ctx context.Context, return_ *Return) error {
	return r.db.WithContext(ctx).Save(return_).Error
}

// DeleteReturn deletes a return
func (r *gormRepository) DeleteReturn(ctx context.Context, tenantID, returnID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, returnID).
		Delete(&Return{}).Error
}

// CreateReturnItem creates a new return item
func (r *gormRepository) CreateReturnItem(ctx context.Context, item *ReturnItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// GetReturnItemByID retrieves a return item by ID
func (r *gormRepository) GetReturnItemByID(ctx context.Context, itemID uuid.UUID) (*ReturnItem, error) {
	var item ReturnItem
	err := r.db.WithContext(ctx).
		Where("id = ?", itemID).
		First(&item).Error
	
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetReturnItems retrieves return items for a return
func (r *gormRepository) GetReturnItems(ctx context.Context, returnID uuid.UUID) ([]*ReturnItem, error) {
	var items []*ReturnItem
	err := r.db.WithContext(ctx).
		Where("return_id = ?", returnID).
		Find(&items).Error
	return items, err
}

// UpdateReturnItem updates a return item
func (r *gormRepository) UpdateReturnItem(ctx context.Context, item *ReturnItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// DeleteReturnItem deletes a return item
func (r *gormRepository) DeleteReturnItem(ctx context.Context, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&ReturnItem{}, itemID).Error
}

// CreateReturnReason creates a new return reason
func (r *gormRepository) CreateReturnReason(ctx context.Context, reason *ReturnReason) error {
	return r.db.WithContext(ctx).Create(reason).Error
}

// GetReturnReasonByID retrieves a return reason by ID
func (r *gormRepository) GetReturnReasonByID(ctx context.Context, tenantID, reasonID uuid.UUID) (*ReturnReason, error) {
	var reason ReturnReason
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, reasonID).
		First(&reason).Error
	
	if err != nil {
		return nil, err
	}
	return &reason, nil
}

// ListReturnReasons retrieves return reasons
func (r *gormRepository) ListReturnReasons(ctx context.Context, tenantID uuid.UUID, activeOnly bool) ([]*ReturnReason, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	
	var reasons []*ReturnReason
	err := query.Order("display_order ASC, name ASC").Find(&reasons).Error
	return reasons, err
}

// UpdateReturnReason updates a return reason
func (r *gormRepository) UpdateReturnReason(ctx context.Context, reason *ReturnReason) error {
	return r.db.WithContext(ctx).Save(reason).Error
}

// DeleteReturnReason deletes a return reason
func (r *gormRepository) DeleteReturnReason(ctx context.Context, tenantID, reasonID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, reasonID).
		Delete(&ReturnReason{}).Error
}

// GetReturnStats retrieves return statistics
func (r *gormRepository) GetReturnStats(ctx context.Context, tenantID uuid.UUID, filter StatsFilter) (*ReturnStats, error) {
	stats := &ReturnStats{}
	
	// Base query
	query := r.db.WithContext(ctx).Model(&Return{}).Where("tenant_id = ?", tenantID)
	
	// Apply date filters
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", filter.DateTo)
	}
	
	// Get total counts and amounts
	var totalCount int64
	var totalRefund float64
	
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats.TotalReturns = totalCount
	
	if err := query.Select("COALESCE(SUM(total_refund), 0)").Scan(&totalRefund).Error; err != nil {
		return nil, err
	}
	stats.TotalRefundAmount = totalRefund
	
	if totalCount > 0 {
		stats.AverageRefundAmount = totalRefund / float64(totalCount)
	}
	
	// Get status breakdown
	var statusCounts []struct {
		Status ReturnStatus `json:"status"`
		Count  int64        `json:"count"`
	}
	
	if err := query.Select("status, COUNT(*) as count").Group("status").Scan(&statusCounts).Error; err != nil {
		return nil, err
	}
	
	for _, sc := range statusCounts {
		switch sc.Status {
		case StatusPending:
			stats.PendingReturns = sc.Count
		case StatusApproved:
			stats.ApprovedReturns = sc.Count
		case StatusRejected:
			stats.RejectedReturns = sc.Count
		case StatusProcessing:
			stats.ProcessingReturns = sc.Count
		case StatusCompleted:
			stats.CompletedReturns = sc.Count
		}
	}
	
	// Get type breakdown
	var typeCounts []struct {
		Type  ReturnType `json:"type"`
		Count int64      `json:"count"`
	}
	
	if err := query.Select("type, COUNT(*) as count").Group("type").Scan(&typeCounts).Error; err != nil {
		return nil, err
	}
	
	for _, tc := range typeCounts {
		switch tc.Type {
		case TypeRefund:
			stats.RefundReturns = tc.Count
		case TypeExchange:
			stats.ExchangeReturns = tc.Count
		}
	}
	
	// Get top reasons
	var reasonStats []ReasonStats
	if err := r.db.WithContext(ctx).
		Table("returns r").
		Select("rr.id as reason_id, rr.name as reason_name, COUNT(*) as count, (COUNT(*) * 100.0 / ?) as percentage", totalCount).
		Joins("LEFT JOIN return_reasons rr ON r.reason_id = rr.id").
		Where("r.tenant_id = ? AND r.reason_id IS NOT NULL", tenantID).
		Group("rr.id, rr.name").
		Order("count DESC").
		Limit(10).
		Scan(&reasonStats).Error; err != nil {
		return nil, err
	}
	stats.TopReasons = reasonStats
	
	return stats, nil
}

// GetReturnsByCustomer retrieves returns for a specific customer
func (r *gormRepository) GetReturnsByCustomer(ctx context.Context, tenantID, customerID uuid.UUID, limit int) ([]*Return, error) {
	query := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Reason").
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	var returns []*Return
	err := query.Find(&returns).Error
	return returns, err
}

// GetReturnsByOrder retrieves returns for a specific order
func (r *gormRepository) GetReturnsByOrder(ctx context.Context, tenantID, orderID uuid.UUID) ([]*Return, error) {
	var returns []*Return
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Reason").
		Where("tenant_id = ? AND order_id = ?", tenantID, orderID).
		Order("created_at DESC").
		Find(&returns).Error
	return returns, err
}

// applyReturnFilters applies filters to the query
func (r *gormRepository) applyReturnFilters(query *gorm.DB, filter ReturnFilter) *gorm.DB {
	// Status filter
	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}
	
	// Type filter
	if len(filter.Type) > 0 {
		query = query.Where("type IN ?", filter.Type)
	}
	
	// Customer filter
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	
	// Order filter
	if filter.OrderID != nil {
		query = query.Where("order_id = ?", *filter.OrderID)
	}
	
	// Reason filter
	if filter.ReasonID != nil {
		query = query.Where("reason_id = ?", *filter.ReasonID)
	}
	
	// Date filters
	if filter.CreatedAfter != nil {
		query = query.Where("created_at >= ?", filter.CreatedAfter)
	}
	if filter.CreatedBefore != nil {
		query = query.Where("created_at <= ?", filter.CreatedBefore)
	}
	
	// Amount filters
	if filter.MinRefundAmount != nil {
		query = query.Where("total_refund >= ?", *filter.MinRefundAmount)
	}
	if filter.MaxRefundAmount != nil {
		query = query.Where("total_refund <= ?", *filter.MaxRefundAmount)
	}
	
	// Search filter
	if filter.Search != "" {
		query = query.Where(
			"return_number ILIKE ? OR reason_text ILIKE ? OR description ILIKE ?",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
		)
	}
	
	return query
}