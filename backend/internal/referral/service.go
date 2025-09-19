package referral

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service defines the interface for referral business logic
type Service interface {
	// Referral management
	GenerateReferralCode(ctx context.Context, tenantID, referrerID uuid.UUID, commissionRate float64, expiresAt *time.Time) (*Referral, error)
	ApplyReferralCode(ctx context.Context, tenantID uuid.UUID, referralCode string, refereeID uuid.UUID) (*Referral, error)
	GetReferralByCode(ctx context.Context, tenantID uuid.UUID, code string) (*Referral, error)
	GetUserReferrals(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]*Referral, error)
	DeactivateReferral(ctx context.Context, tenantID, referralID uuid.UUID) error

	// Commission management
	CreateCommission(ctx context.Context, tenantID, referralID, subscriptionID uuid.UUID, subscriptionAmount float64) (*ReferralCommission, error)
	GetUserCommissions(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]*ReferralCommission, error)
	ProcessPendingCommissions(ctx context.Context, tenantID uuid.UUID, limit int) ([]*ReferralCommission, error)
	MarkCommissionAsPaid(ctx context.Context, tenantID, commissionID uuid.UUID) error

	// Statistics and analytics
	GetUserReferralStats(ctx context.Context, tenantID, userID uuid.UUID) (*ReferralStats, error)
	GetTenantReferralStats(ctx context.Context, tenantID uuid.UUID) (*ReferralStats, error)

	// Maintenance
	ExpireOldReferrals(ctx context.Context, tenantID uuid.UUID) error
	ValidateReferralCode(ctx context.Context, tenantID uuid.UUID, code string) (*Referral, error)
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new referral service
func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

// GenerateReferralCode generates a new referral code for a user
func (s *ServiceImpl) GenerateReferralCode(ctx context.Context, tenantID, referrerID uuid.UUID, commissionRate float64, expiresAt *time.Time) (*Referral, error) {
	// Validate commission rate
	if commissionRate < 0 || commissionRate > 1 {
		return nil, errors.New("commission rate must be between 0 and 1")
	}

	// Generate unique referral code
	code, err := s.generateUniqueCode(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate unique code: %w", err)
	}

	// Create referral
	referral := &Referral{
		ID:             uuid.New(),
		TenantID:       tenantID,
		ReferrerID:     referrerID,
		ReferralCode:   code,
		Status:         ReferralStatusActive,
		CommissionRate: commissionRate,
		ExpiresAt:      expiresAt,
	}

	err = s.repo.CreateReferral(ctx, referral)
	if err != nil {
		return nil, fmt.Errorf("failed to create referral: %w", err)
	}

	return referral, nil
}

// ApplyReferralCode applies a referral code when a user signs up
func (s *ServiceImpl) ApplyReferralCode(ctx context.Context, tenantID uuid.UUID, referralCode string, refereeID uuid.UUID) (*Referral, error) {
	// Validate referral code
	referral, err := s.ValidateReferralCode(ctx, tenantID, referralCode)
	if err != nil {
		return nil, err
	}

	// Check if user is trying to refer themselves
	if referral.ReferrerID == refereeID {
		return nil, errors.New("users cannot refer themselves")
	}

	// Check if referral already has a referee
	if referral.RefereeID != nil {
		return nil, errors.New("referral code has already been used")
	}

	// Apply the referral
	referral.RefereeID = &refereeID
	referral.Status = ReferralStatusCompleted
	now := time.Now()
	referral.CompletedAt = &now

	err = s.repo.UpdateReferral(ctx, referral)
	if err != nil {
		return nil, fmt.Errorf("failed to apply referral code: %w", err)
	}

	return referral, nil
}

// GetReferralByCode retrieves a referral by its code
func (s *ServiceImpl) GetReferralByCode(ctx context.Context, tenantID uuid.UUID, code string) (*Referral, error) {
	return s.repo.GetReferralByCode(ctx, tenantID, code)
}

// GetUserReferrals retrieves all referrals created by a user
func (s *ServiceImpl) GetUserReferrals(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]*Referral, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetReferralsByReferrer(ctx, tenantID, userID, limit, offset)
}

// DeactivateReferral deactivates a referral
func (s *ServiceImpl) DeactivateReferral(ctx context.Context, tenantID, referralID uuid.UUID) error {
	referral, err := s.repo.GetReferralByID(ctx, tenantID, referralID)
	if err != nil {
		return fmt.Errorf("failed to get referral: %w", err)
	}

	referral.Status = ReferralStatusCancelled

	err = s.repo.UpdateReferral(ctx, referral)
	if err != nil {
		return fmt.Errorf("failed to deactivate referral: %w", err)
	}

	return nil
}

// CreateCommission creates a commission when a referred user subscribes
func (s *ServiceImpl) CreateCommission(ctx context.Context, tenantID, referralID, subscriptionID uuid.UUID, subscriptionAmount float64) (*ReferralCommission, error) {
	// Get the referral
	referral, err := s.repo.GetReferralByID(ctx, tenantID, referralID)
	if err != nil {
		return nil, fmt.Errorf("failed to get referral: %w", err)
	}

	// Validate referral status
	if referral.Status != ReferralStatusCompleted {
		return nil, errors.New("referral must be completed to generate commission")
	}

	// Calculate commission amount
	commissionAmount := referral.CalculateCommission(subscriptionAmount)

	// Create commission
	commission := &ReferralCommission{
		ID:                 uuid.New(),
		TenantID:           tenantID,
		ReferralID:         referralID,
		ReferrerID:         referral.ReferrerID,
		SubscriptionID:     subscriptionID,
		Amount:             commissionAmount,
		Currency:           "USD", // Default currency, could be configurable
		CommissionRate:     referral.CommissionRate,
		SubscriptionAmount: subscriptionAmount,
		Status:             CommissionStatusPending,
	}

	err = s.repo.CreateCommission(ctx, commission)
	if err != nil {
		return nil, fmt.Errorf("failed to create commission: %w", err)
	}

	return commission, nil
}

// GetUserCommissions retrieves all commissions for a user
func (s *ServiceImpl) GetUserCommissions(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]*ReferralCommission, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetCommissionsByReferrer(ctx, tenantID, userID, limit, offset)
}

// ProcessPendingCommissions retrieves pending commissions for processing
func (s *ServiceImpl) ProcessPendingCommissions(ctx context.Context, tenantID uuid.UUID, limit int) ([]*ReferralCommission, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	return s.repo.GetPendingCommissions(ctx, tenantID, limit, 0)
}

// MarkCommissionAsPaid marks a commission as paid
func (s *ServiceImpl) MarkCommissionAsPaid(ctx context.Context, tenantID, commissionID uuid.UUID) error {
	// Note: In a real implementation, you'd first get the commission
	// For now, we'll create a commission object with the ID and mark it as paid
	commission := &ReferralCommission{
		ID:       commissionID,
		TenantID: tenantID,
	}
	commission.MarkAsPaid()

	err := s.repo.UpdateCommission(ctx, commission)
	if err != nil {
		return fmt.Errorf("failed to mark commission as paid: %w", err)
	}

	return nil
}

// GetUserReferralStats retrieves referral statistics for a user
func (s *ServiceImpl) GetUserReferralStats(ctx context.Context, tenantID, userID uuid.UUID) (*ReferralStats, error) {
	return s.repo.GetReferralStats(ctx, tenantID, userID)
}

// GetTenantReferralStats retrieves referral statistics for the entire tenant
func (s *ServiceImpl) GetTenantReferralStats(ctx context.Context, tenantID uuid.UUID) (*ReferralStats, error) {
	return s.repo.GetTenantReferralStats(ctx, tenantID)
}

// ExpireOldReferrals expires old referrals
func (s *ServiceImpl) ExpireOldReferrals(ctx context.Context, tenantID uuid.UUID) error {
	return s.repo.ExpireOldReferrals(ctx, tenantID)
}

// ValidateReferralCode validates a referral code
func (s *ServiceImpl) ValidateReferralCode(ctx context.Context, tenantID uuid.UUID, code string) (*Referral, error) {
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("referral code cannot be empty")
	}

	referral, err := s.repo.GetReferralByCode(ctx, tenantID, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid referral code")
		}
		return nil, fmt.Errorf("failed to get referral: %w", err)
	}

	// Check if referral can be used
	if !referral.CanBeUsed() {
		if referral.IsExpired() {
			return nil, errors.New("referral code has expired")
		}
		if referral.Status == ReferralStatusCompleted {
			return nil, errors.New("referral code has already been used")
		}
		if referral.Status == ReferralStatusCancelled {
			return nil, errors.New("referral code has been cancelled")
		}
		return nil, errors.New("referral code is not active")
	}

	return referral, nil
}

// generateUniqueCode generates a unique referral code
func (s *ServiceImpl) generateUniqueCode(ctx context.Context, tenantID uuid.UUID) (string, error) {
	const maxAttempts = 10

	for i := 0; i < maxAttempts; i++ {
		code := s.generateRandomCode()
		isUnique, err := s.repo.IsReferralCodeUnique(ctx, tenantID, code)
		if err != nil {
			return "", err
		}
		if isUnique {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique referral code after multiple attempts")
}

// generateRandomCode generates a random referral code
func (s *ServiceImpl) generateRandomCode() string {
	bytes := make([]byte, 4) // 8 character hex string
	rand.Read(bytes)
	return strings.ToUpper(hex.EncodeToString(bytes))
}