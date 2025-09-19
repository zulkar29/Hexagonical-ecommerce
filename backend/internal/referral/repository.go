package referral

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for referral data operations
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
	UpdateCommission(ctx context.Context, commission *ReferralCommission) error
	GetPendingCommissions(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*ReferralCommission, error)

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