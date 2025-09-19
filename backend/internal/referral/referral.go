package referral

import (
	"time"

	"github.com/google/uuid"
)

// ReferralStatus represents the status of a referral
type ReferralStatus string

const (
	ReferralStatusPending   ReferralStatus = "pending"
	ReferralStatusActive    ReferralStatus = "active"
	ReferralStatusCompleted ReferralStatus = "completed"
	ReferralStatusExpired   ReferralStatus = "expired"
	ReferralStatusCancelled ReferralStatus = "cancelled"
)

// CommissionStatus represents the status of a commission
type CommissionStatus string

const (
	CommissionStatusPending CommissionStatus = "pending"
	CommissionStatusPaid    CommissionStatus = "paid"
	CommissionStatusFailed  CommissionStatus = "failed"
	CommissionStatusVoided  CommissionStatus = "voided"
)

// Referral represents a referral relationship between users
type Referral struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID     uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ReferrerID   uuid.UUID      `json:"referrer_id" gorm:"type:uuid;not null;index"`
	RefereeID    *uuid.UUID     `json:"referee_id" gorm:"type:uuid;index"`
	ReferralCode string         `json:"referral_code" gorm:"type:varchar(50);unique;not null;index"`
	Status       ReferralStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	CommissionRate float64      `json:"commission_rate" gorm:"type:decimal(5,4);not null;default:0.1"`
	ExpiresAt    *time.Time     `json:"expires_at"`
	CompletedAt  *time.Time     `json:"completed_at"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// ReferralCommission represents a commission earned from a referral
type ReferralCommission struct {
	ID               uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID         uuid.UUID        `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ReferralID       uuid.UUID        `json:"referral_id" gorm:"type:uuid;not null;index"`
	ReferrerID       uuid.UUID        `json:"referrer_id" gorm:"type:uuid;not null;index"`
	SubscriptionID   uuid.UUID        `json:"subscription_id" gorm:"type:uuid;not null;index"`
	Amount           float64          `json:"amount" gorm:"type:decimal(10,2);not null"`
	Currency         string           `json:"currency" gorm:"type:varchar(3);not null;default:'USD'"`
	CommissionRate   float64          `json:"commission_rate" gorm:"type:decimal(5,4);not null"`
	SubscriptionAmount float64        `json:"subscription_amount" gorm:"type:decimal(10,2);not null"`
	Status           CommissionStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	PaidAt           *time.Time       `json:"paid_at"`
	CreatedAt        time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time        `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Referral *Referral `json:"referral,omitempty" gorm:"foreignKey:ReferralID"`
}

// ReferralStats represents aggregated referral statistics
type ReferralStats struct {
	TotalReferrals     int64   `json:"total_referrals"`
	ActiveReferrals    int64   `json:"active_referrals"`
	CompletedReferrals int64   `json:"completed_referrals"`
	TotalCommissions   float64 `json:"total_commissions"`
	PendingCommissions float64 `json:"pending_commissions"`
	PaidCommissions    float64 `json:"paid_commissions"`
}

// TableName returns the table name for Referral
func (Referral) TableName() string {
	return "referrals"
}

// TableName returns the table name for ReferralCommission
func (ReferralCommission) TableName() string {
	return "referral_commissions"
}

// IsExpired checks if the referral has expired
func (r *Referral) IsExpired() bool {
	if r.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*r.ExpiresAt)
}

// CanBeUsed checks if the referral can be used
func (r *Referral) CanBeUsed() bool {
	return r.Status == ReferralStatusActive && !r.IsExpired()
}

// MarkAsCompleted marks the referral as completed
func (r *Referral) MarkAsCompleted() {
	r.Status = ReferralStatusCompleted
	now := time.Now()
	r.CompletedAt = &now
}

// MarkAsExpired marks the referral as expired
func (r *Referral) MarkAsExpired() {
	r.Status = ReferralStatusExpired
}

// CalculateCommission calculates the commission amount based on subscription amount
func (r *Referral) CalculateCommission(subscriptionAmount float64) float64 {
	return subscriptionAmount * r.CommissionRate
}

// MarkAsPaid marks the commission as paid
func (rc *ReferralCommission) MarkAsPaid() {
	rc.Status = CommissionStatusPaid
	now := time.Now()
	rc.PaidAt = &now
}

// MarkAsFailed marks the commission as failed
func (rc *ReferralCommission) MarkAsFailed() {
	rc.Status = CommissionStatusFailed
}

// MarkAsVoided marks the commission as voided
func (rc *ReferralCommission) MarkAsVoided() {
	rc.Status = CommissionStatusVoided
}