package referral

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for referral/affiliate data operations
type Repository interface {
	// Referral operations
	CreateReferral(ctx context.Context, referral *Referral) error
	GetReferralByID(ctx context.Context, tenantID, referralID uuid.UUID) (*Referral, error)
	GetReferralByCode(ctx context.Context, tenantID uuid.UUID, code string) (*Referral, error)
	GetReferralsByReferrer(ctx context.Context, tenantID, referrerID uuid.UUID, limit, offset int) ([]*Referral, error)
	UpdateReferral(ctx context.Context, referral *Referral) error
	DeleteReferral(ctx context.Context, tenantID, referralID uuid.UUID) error

	// Commission operations
	CreateCommission(ctx context.Context, commission *ReferralCommission) error
	GetCommissionsByReferrer(ctx context.Context, tenantID, referrerID uuid.UUID, limit, offset int) ([]*ReferralCommission, error)
	GetCommissionsByReferral(ctx context.Context, tenantID, referralID uuid.UUID) ([]*ReferralCommission, error)
	GetPendingCommissionsByReferrer(ctx context.Context, tenantID, referrerID uuid.UUID) ([]*ReferralCommission, error)
	UpdateCommission(ctx context.Context, commission *ReferralCommission) error
	GetPendingCommissions(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*ReferralCommission, error)

	// Affiliate click tracking
	CreateAffiliateClick(ctx context.Context, click *AffiliateClick) error
	GetAffiliateClicksByReferrer(ctx context.Context, tenantID, referrerID uuid.UUID, limit, offset int) ([]*AffiliateClick, error)
	GetAffiliateClicksByDateRange(ctx context.Context, tenantID, referrerID uuid.UUID, startDate, endDate time.Time) ([]*AffiliateClick, error)

	// Performance analytics
	GetAffiliatePerformance(ctx context.Context, tenantID, referrerID uuid.UUID) (*AffiliatePerformance, error)
	GetTopPerformingAffiliates(ctx context.Context, tenantID uuid.UUID, limit int) ([]*AffiliatePerformance, error)

	// Payout batch operations
	CreatePayoutBatch(ctx context.Context, batch *AffiliatePayoutBatch) error
	GetPayoutBatch(ctx context.Context, tenantID, batchID uuid.UUID) (*AffiliatePayoutBatch, error)
	UpdatePayoutBatch(ctx context.Context, batch *AffiliatePayoutBatch) error
	GetPayoutBatches(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*AffiliatePayoutBatch, error)

	// Statistics operations
	GetReferralStats(ctx context.Context, tenantID, referrerID uuid.UUID) (*ReferralStats, error)
	GetTenantReferralStats(ctx context.Context, tenantID uuid.UUID) (*ReferralStats, error)

	// Utility operations
	IsReferralCodeUnique(ctx context.Context, tenantID uuid.UUID, code string) (bool, error)
	ExpireOldReferrals(ctx context.Context, tenantID uuid.UUID) error
}

// GormRepository implements Repository using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM repository
func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// CreateReferral creates a new referral
func (r *GormRepository) CreateReferral(ctx context.Context, referral *Referral) error {
	return r.db.WithContext(ctx).Create(referral).Error
}

// GetReferralByID retrieves a referral by ID
func (r *GormRepository) GetReferralByID(ctx context.Context, tenantID, referralID uuid.UUID) (*Referral, error) {
	var referral Referral
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, referralID).
		First(&referral).Error
	if err != nil {
		return nil, err
	}
	return &referral, nil
}

// GetReferralByCode retrieves a referral by code
func (r *GormRepository) GetReferralByCode(ctx context.Context, tenantID uuid.UUID, code string) (*Referral, error) {
	var referral Referral
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND referral_code = ?", tenantID, code).
		First(&referral).Error
	if err != nil {
		return nil, err
	}
	return &referral, nil
}

// GetReferralsByReferrer retrieves referrals by referrer ID
func (r *GormRepository) GetReferralsByReferrer(ctx context.Context, tenantID, referrerID uuid.UUID, limit, offset int) ([]*Referral, error) {
	var referrals []*Referral
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND referrer_id = ?", tenantID, referrerID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&referrals).Error
	return referrals, err
}

// UpdateReferral updates a referral
func (r *GormRepository) UpdateReferral(ctx context.Context, referral *Referral) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", referral.TenantID, referral.ID).
		Updates(referral).Error
}

// DeleteReferral deletes a referral
func (r *GormRepository) DeleteReferral(ctx context.Context, tenantID, referralID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, referralID).
		Delete(&Referral{}).Error
}

// CreateCommission creates a new commission
func (r *GormRepository) CreateCommission(ctx context.Context, commission *ReferralCommission) error {
	return r.db.WithContext(ctx).Create(commission).Error
}

// GetCommissionsByReferrer retrieves commissions by referrer ID
func (r *GormRepository) GetCommissionsByReferrer(ctx context.Context, tenantID, referrerID uuid.UUID, limit, offset int) ([]*ReferralCommission, error) {
	var commissions []*ReferralCommission
	err := r.db.WithContext(ctx).
		Preload("Referral").
		Where("tenant_id = ? AND referrer_id = ?", tenantID, referrerID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&commissions).Error
	return commissions, err
}

// GetCommissionsByReferral retrieves commissions by referral ID
func (r *GormRepository) GetCommissionsByReferral(ctx context.Context, tenantID, referralID uuid.UUID) ([]*ReferralCommission, error) {
	var commissions []*ReferralCommission
	err := r.db.WithContext(ctx).
		Preload("Referral").
		Where("tenant_id = ? AND referral_id = ?", tenantID, referralID).
		Order("created_at DESC").
		Find(&commissions).Error
	return commissions, err
}

// UpdateCommission updates a commission
func (r *GormRepository) UpdateCommission(ctx context.Context, commission *ReferralCommission) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", commission.TenantID, commission.ID).
		Updates(commission).Error
}

// GetPendingCommissions retrieves pending commissions
func (r *GormRepository) GetPendingCommissions(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*ReferralCommission, error) {
	var commissions []*ReferralCommission
	err := r.db.WithContext(ctx).
		Preload("Referral").
		Where("tenant_id = ? AND status = ?", tenantID, CommissionStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&commissions).Error
	return commissions, err
}

// GetReferralStats retrieves referral statistics for a specific referrer
func (r *GormRepository) GetReferralStats(ctx context.Context, tenantID, referrerID uuid.UUID) (*ReferralStats, error) {
	stats := &ReferralStats{}

	// Count referrals by status
	var totalReferrals, activeReferrals, completedReferrals int64

	r.db.WithContext(ctx).
		Model(&Referral{}).
		Where("tenant_id = ? AND referrer_id = ?", tenantID, referrerID).
		Count(&totalReferrals)

	r.db.WithContext(ctx).
		Model(&Referral{}).
		Where("tenant_id = ? AND referrer_id = ? AND status = ?", tenantID, referrerID, ReferralStatusActive).
		Count(&activeReferrals)

	r.db.WithContext(ctx).
		Model(&Referral{}).
		Where("tenant_id = ? AND referrer_id = ? AND status = ?", tenantID, referrerID, ReferralStatusCompleted).
		Count(&completedReferrals)

	// Calculate commission totals
	var totalCommissions, pendingCommissions, paidCommissions float64

	r.db.WithContext(ctx).
		Model(&ReferralCommission{}).
		Where("tenant_id = ? AND referrer_id = ?", tenantID, referrerID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalCommissions)

	r.db.WithContext(ctx).
		Model(&ReferralCommission{}).
		Where("tenant_id = ? AND referrer_id = ? AND status = ?", tenantID, referrerID, CommissionStatusPending).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&pendingCommissions)

	r.db.WithContext(ctx).
		Model(&ReferralCommission{}).
		Where("tenant_id = ? AND referrer_id = ? AND status = ?", tenantID, referrerID, CommissionStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&paidCommissions)

	stats.TotalReferrals = totalReferrals
	stats.ActiveReferrals = activeReferrals
	stats.CompletedReferrals = completedReferrals
	stats.TotalCommissions = totalCommissions
	stats.PendingCommissions = pendingCommissions
	stats.PaidCommissions = paidCommissions

	return stats, nil
}

// GetTenantReferralStats retrieves referral statistics for the entire tenant
func (r *GormRepository) GetTenantReferralStats(ctx context.Context, tenantID uuid.UUID) (*ReferralStats, error) {
	stats := &ReferralStats{}

	// Count referrals by status
	var totalReferrals, activeReferrals, completedReferrals int64

	r.db.WithContext(ctx).
		Model(&Referral{}).
		Where("tenant_id = ?", tenantID).
		Count(&totalReferrals)

	r.db.WithContext(ctx).
		Model(&Referral{}).
		Where("tenant_id = ? AND status = ?", tenantID, ReferralStatusActive).
		Count(&activeReferrals)

	r.db.WithContext(ctx).
		Model(&Referral{}).
		Where("tenant_id = ? AND status = ?", tenantID, ReferralStatusCompleted).
		Count(&completedReferrals)

	// Calculate commission totals
	var totalCommissions, pendingCommissions, paidCommissions float64

	r.db.WithContext(ctx).
		Model(&ReferralCommission{}).
		Where("tenant_id = ?", tenantID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalCommissions)

	r.db.WithContext(ctx).
		Model(&ReferralCommission{}).
		Where("tenant_id = ? AND status = ?", tenantID, CommissionStatusPending).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&pendingCommissions)

	r.db.WithContext(ctx).
		Model(&ReferralCommission{}).
		Where("tenant_id = ? AND status = ?", tenantID, CommissionStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&paidCommissions)

	stats.TotalReferrals = totalReferrals
	stats.ActiveReferrals = activeReferrals
	stats.CompletedReferrals = completedReferrals
	stats.TotalCommissions = totalCommissions
	stats.PendingCommissions = pendingCommissions
	stats.PaidCommissions = paidCommissions

	return stats, nil
}

// IsReferralCodeUnique checks if a referral code is unique within a tenant
func (r *GormRepository) IsReferralCodeUnique(ctx context.Context, tenantID uuid.UUID, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Referral{}).
		Where("tenant_id = ? AND referral_code = ?", tenantID, code).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// ExpireOldReferrals marks expired referrals as expired
func (r *GormRepository) ExpireOldReferrals(ctx context.Context, tenantID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&Referral{}).
		Where("tenant_id = ? AND status IN (?, ?) AND expires_at IS NOT NULL AND expires_at < ?",
			tenantID, ReferralStatusPending, ReferralStatusActive, now).
		Update("status", ReferralStatusExpired).Error
}

// GetPendingCommissionsByReferrer retrieves pending commissions for a specific referrer
func (r *GormRepository) GetPendingCommissionsByReferrer(ctx context.Context, tenantID, referrerID uuid.UUID) ([]*ReferralCommission, error) {
	var commissions []*ReferralCommission
	err := r.db.WithContext(ctx).
		Preload("Referral").
		Where("tenant_id = ? AND referrer_id = ? AND status = ?", tenantID, referrerID, CommissionStatusPending).
		Order("created_at ASC").
		Find(&commissions).Error
	return commissions, err
}

// CreateAffiliateClick creates a new affiliate click record
func (r *GormRepository) CreateAffiliateClick(ctx context.Context, click *AffiliateClick) error {
	return r.db.WithContext(ctx).Create(click).Error
}

// GetAffiliateClicksByReferrer retrieves clicks for a specific affiliate
func (r *GormRepository) GetAffiliateClicksByReferrer(ctx context.Context, tenantID, referrerID uuid.UUID, limit, offset int) ([]*AffiliateClick, error) {
	var clicks []*AffiliateClick
	err := r.db.WithContext(ctx).
		Preload("Referral").
		Where("tenant_id = ? AND referrer_id = ?", tenantID, referrerID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&clicks).Error
	return clicks, err
}

// GetAffiliateClicksByDateRange retrieves clicks within a date range
func (r *GormRepository) GetAffiliateClicksByDateRange(ctx context.Context, tenantID, referrerID uuid.UUID, startDate, endDate time.Time) ([]*AffiliateClick, error) {
	var clicks []*AffiliateClick
	err := r.db.WithContext(ctx).
		Preload("Referral").
		Where("tenant_id = ? AND referrer_id = ? AND created_at BETWEEN ? AND ?", tenantID, referrerID, startDate, endDate).
		Order("created_at DESC").
		Find(&clicks).Error
	return clicks, err
}

// GetAffiliatePerformance retrieves performance metrics for an affiliate
func (r *GormRepository) GetAffiliatePerformance(ctx context.Context, tenantID, referrerID uuid.UUID) (*AffiliatePerformance, error) {
	performance := &AffiliatePerformance{
		TenantID:   tenantID,
		ReferrerID: referrerID,
	}

	// Get click counts
	var totalClicks, uniqueClicks, conversions int64

	r.db.WithContext(ctx).
		Model(&AffiliateClick{}).
		Where("tenant_id = ? AND referrer_id = ?", tenantID, referrerID).
		Count(&totalClicks)

	r.db.WithContext(ctx).
		Model(&AffiliateClick{}).
		Where("tenant_id = ? AND referrer_id = ?", tenantID, referrerID).
		Distinct("ip_address").
		Count(&uniqueClicks)

	r.db.WithContext(ctx).
		Model(&AffiliateClick{}).
		Where("tenant_id = ? AND referrer_id = ? AND converted = ?", tenantID, referrerID, true).
		Count(&conversions)

	// Get earnings
	var earningsToday, earningsThisWeek, earningsThisMonth float64
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := startOfDay.AddDate(0, 0, -int(now.Weekday()))
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	r.db.WithContext(ctx).
		Model(&ReferralCommission{}).
		Where("tenant_id = ? AND referrer_id = ? AND created_at >= ?", tenantID, referrerID, startOfDay).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&earningsToday)

	r.db.WithContext(ctx).
		Model(&ReferralCommission{}).
		Where("tenant_id = ? AND referrer_id = ? AND created_at >= ?", tenantID, referrerID, startOfWeek).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&earningsThisWeek)

	r.db.WithContext(ctx).
		Model(&ReferralCommission{}).
		Where("tenant_id = ? AND referrer_id = ? AND created_at >= ?", tenantID, referrerID, startOfMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&earningsThisMonth)

	// Calculate click counts for different periods
	var clicksToday, clicksThisWeek, clicksThisMonth int64

	r.db.WithContext(ctx).
		Model(&AffiliateClick{}).
		Where("tenant_id = ? AND referrer_id = ? AND created_at >= ?", tenantID, referrerID, startOfDay).
		Count(&clicksToday)

	r.db.WithContext(ctx).
		Model(&AffiliateClick{}).
		Where("tenant_id = ? AND referrer_id = ? AND created_at >= ?", tenantID, referrerID, startOfWeek).
		Count(&clicksThisWeek)

	r.db.WithContext(ctx).
		Model(&AffiliateClick{}).
		Where("tenant_id = ? AND referrer_id = ? AND created_at >= ?", tenantID, referrerID, startOfMonth).
		Count(&clicksThisMonth)

	// Calculate conversion rate
	conversionRate := 0.0
	if totalClicks > 0 {
		conversionRate = float64(conversions) / float64(totalClicks) * 100
	}

	performance.TotalClicks = totalClicks
	performance.UniqueClicks = uniqueClicks
	performance.Conversions = conversions
	performance.ConversionRate = conversionRate
	performance.ClicksToday = clicksToday
	performance.ClicksThisWeek = clicksThisWeek
	performance.ClicksThisMonth = clicksThisMonth
	performance.EarningsToday = earningsToday
	performance.EarningsThisWeek = earningsThisWeek
	performance.EarningsThisMonth = earningsThisMonth

	return performance, nil
}

// GetTopPerformingAffiliates retrieves top performing affiliates
func (r *GormRepository) GetTopPerformingAffiliates(ctx context.Context, tenantID uuid.UUID, limit int) ([]*AffiliatePerformance, error) {
	var referrers []uuid.UUID

	// Get top referrers by total commissions
	err := r.db.WithContext(ctx).
		Model(&ReferralCommission{}).
		Select("referrer_id").
		Where("tenant_id = ?", tenantID).
		Group("referrer_id").
		Order("SUM(amount) DESC").
		Limit(limit).
		Pluck("referrer_id", &referrers).Error

	if err != nil {
		return nil, err
	}

	var performances []*AffiliatePerformance
	for _, referrerID := range referrers {
		performance, err := r.GetAffiliatePerformance(ctx, tenantID, referrerID)
		if err != nil {
			continue // Skip on error, don't fail entire operation
		}
		performances = append(performances, performance)
	}

	return performances, nil
}

// CreatePayoutBatch creates a new payout batch
func (r *GormRepository) CreatePayoutBatch(ctx context.Context, batch *AffiliatePayoutBatch) error {
	return r.db.WithContext(ctx).Create(batch).Error
}

// GetPayoutBatch retrieves a payout batch by ID
func (r *GormRepository) GetPayoutBatch(ctx context.Context, tenantID, batchID uuid.UUID) (*AffiliatePayoutBatch, error) {
	var batch AffiliatePayoutBatch
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, batchID).
		First(&batch).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// UpdatePayoutBatch updates a payout batch
func (r *GormRepository) UpdatePayoutBatch(ctx context.Context, batch *AffiliatePayoutBatch) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", batch.TenantID, batch.ID).
		Updates(batch).Error
}

// GetPayoutBatches retrieves payout batches
func (r *GormRepository) GetPayoutBatches(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*AffiliatePayoutBatch, error) {
	var batches []*AffiliatePayoutBatch
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&batches).Error
	return batches, err
}