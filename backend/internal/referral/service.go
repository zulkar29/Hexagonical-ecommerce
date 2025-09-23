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

// Service defines the interface for referral/affiliate business logic
type Service interface {
	// Referral/Affiliate management
	GenerateReferralCode(ctx context.Context, tenantID, referrerID uuid.UUID, commissionRate float64, expiresAt *time.Time) (*Referral, error)
	CreateAffiliateAccount(ctx context.Context, tenantID, userID uuid.UUID, affiliateType AffiliateType, commissionRate float64, payoutThreshold float64) (*Referral, error)
	ApplyReferralCode(ctx context.Context, tenantID uuid.UUID, referralCode string, refereeID uuid.UUID) (*Referral, error)
	GetReferralByCode(ctx context.Context, tenantID uuid.UUID, code string) (*Referral, error)
	GetUserReferrals(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]*Referral, error)
	UpdateAffiliateSettings(ctx context.Context, tenantID, referralID uuid.UUID, settings map[string]interface{}) error
	DeactivateReferral(ctx context.Context, tenantID, referralID uuid.UUID) error

	// Click tracking
	TrackAffiliateClick(ctx context.Context, tenantID uuid.UUID, referralCode string, clickData *AffiliateClick) (*AffiliateClick, error)
	GetAffiliateClicks(ctx context.Context, tenantID, referrerID uuid.UUID, limit, offset int) ([]*AffiliateClick, error)
	GetClicksByDateRange(ctx context.Context, tenantID, referrerID uuid.UUID, startDate, endDate time.Time) ([]*AffiliateClick, error)

	// Commission management
	CreateCommission(ctx context.Context, tenantID, referralID, subscriptionID uuid.UUID, subscriptionAmount float64) (*ReferralCommission, error)
	CreateOrderCommission(ctx context.Context, tenantID, referralID, orderID uuid.UUID, orderAmount float64) (*ReferralCommission, error)
	GetUserCommissions(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]*ReferralCommission, error)
	ProcessPendingCommissions(ctx context.Context, tenantID uuid.UUID, limit int) ([]*ReferralCommission, error)
	MarkCommissionAsPaid(ctx context.Context, tenantID, commissionID uuid.UUID) error

	// Affiliate performance & analytics
	GetAffiliatePerformance(ctx context.Context, tenantID, referrerID uuid.UUID) (*AffiliatePerformance, error)
	GetTopPerformingAffiliates(ctx context.Context, tenantID uuid.UUID, limit int) ([]*AffiliatePerformance, error)
	GetUserReferralStats(ctx context.Context, tenantID, userID uuid.UUID) (*ReferralStats, error)
	GetTenantReferralStats(ctx context.Context, tenantID uuid.UUID) (*ReferralStats, error)

	// Payout management
	CreatePayoutBatch(ctx context.Context, tenantID uuid.UUID, affiliateIDs []uuid.UUID) (*AffiliatePayoutBatch, error)
	ProcessPayoutBatch(ctx context.Context, tenantID, batchID uuid.UUID) error
	GetPayoutBatches(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*AffiliatePayoutBatch, error)

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
	return s.CreateAffiliateAccount(ctx, tenantID, referrerID, AffiliateTypeCustomer, commissionRate, 50.0)
}

// CreateAffiliateAccount creates a new affiliate account with specific type and settings
func (s *ServiceImpl) CreateAffiliateAccount(ctx context.Context, tenantID, userID uuid.UUID, affiliateType AffiliateType, commissionRate float64, payoutThreshold float64) (*Referral, error) {
	// Validate commission rate
	if commissionRate < 0 || commissionRate > 1 {
		return nil, errors.New("commission rate must be between 0 and 1")
	}

	// Validate payout threshold
	if payoutThreshold < 0 {
		return nil, errors.New("payout threshold must be positive")
	}

	// Generate unique referral code
	code, err := s.generateUniqueCode(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate unique code: %w", err)
	}

	// Create referral/affiliate account
	referral := &Referral{
		ID:              uuid.New(),
		TenantID:        tenantID,
		ReferrerID:      userID,
		ReferralCode:    code,
		Status:          ReferralStatusActive,
		CommissionRate:  commissionRate,
		AffiliateType:   affiliateType,
		PayoutThreshold: payoutThreshold,
		TrackingData:    make(map[string]interface{}),
	}

	err = s.repo.CreateReferral(ctx, referral)
	if err != nil {
		return nil, fmt.Errorf("failed to create affiliate account: %w", err)
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
		SubscriptionID:     &subscriptionID,
		Amount:             commissionAmount,
		Currency:           "USD", // Default currency, could be configurable
		CommissionRate:     referral.CommissionRate,
		SubscriptionAmount: subscriptionAmount,
		ConversionType:     "subscription",
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

// UpdateAffiliateSettings updates affiliate-specific settings
func (s *ServiceImpl) UpdateAffiliateSettings(ctx context.Context, tenantID, referralID uuid.UUID, settings map[string]interface{}) error {
	referral, err := s.repo.GetReferralByID(ctx, tenantID, referralID)
	if err != nil {
		return fmt.Errorf("failed to get referral: %w", err)
	}

	// Update allowed settings
	if payoutThreshold, ok := settings["payout_threshold"].(float64); ok && payoutThreshold >= 0 {
		referral.PayoutThreshold = payoutThreshold
	}

	if affiliateType, ok := settings["affiliate_type"].(string); ok {
		referral.AffiliateType = AffiliateType(affiliateType)
	}

	// Update tracking data
	if trackingData, ok := settings["tracking_data"].(map[string]interface{}); ok {
		if referral.TrackingData == nil {
			referral.TrackingData = make(map[string]interface{})
		}
		for key, value := range trackingData {
			referral.TrackingData[key] = value
		}
	}

	return s.repo.UpdateReferral(ctx, referral)
}

// TrackAffiliateClick tracks a click on an affiliate link
func (s *ServiceImpl) TrackAffiliateClick(ctx context.Context, tenantID uuid.UUID, referralCode string, clickData *AffiliateClick) (*AffiliateClick, error) {
	// Validate referral code exists
	referral, err := s.ValidateReferralCode(ctx, tenantID, referralCode)
	if err != nil {
		return nil, err
	}

	// Set required fields
	clickData.ID = uuid.New()
	clickData.TenantID = tenantID
	clickData.ReferralID = referral.ID
	clickData.ReferrerID = referral.ReferrerID

	// Create click record
	err = s.repo.CreateAffiliateClick(ctx, clickData)
	if err != nil {
		return nil, fmt.Errorf("failed to create affiliate click: %w", err)
	}

	return clickData, nil
}

// GetAffiliateClicks retrieves clicks for an affiliate
func (s *ServiceImpl) GetAffiliateClicks(ctx context.Context, tenantID, referrerID uuid.UUID, limit, offset int) ([]*AffiliateClick, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetAffiliateClicksByReferrer(ctx, tenantID, referrerID, limit, offset)
}

// GetClicksByDateRange retrieves clicks within a date range
func (s *ServiceImpl) GetClicksByDateRange(ctx context.Context, tenantID, referrerID uuid.UUID, startDate, endDate time.Time) ([]*AffiliateClick, error) {
	return s.repo.GetAffiliateClicksByDateRange(ctx, tenantID, referrerID, startDate, endDate)
}

// CreateOrderCommission creates a commission for an order
func (s *ServiceImpl) CreateOrderCommission(ctx context.Context, tenantID, referralID, orderID uuid.UUID, orderAmount float64) (*ReferralCommission, error) {
	// Get the referral
	referral, err := s.repo.GetReferralByID(ctx, tenantID, referralID)
	if err != nil {
		return nil, fmt.Errorf("failed to get referral: %w", err)
	}

	// Validate referral status
	if referral.Status != ReferralStatusCompleted && referral.Status != ReferralStatusActive {
		return nil, errors.New("referral must be active or completed to generate commission")
	}

	// Calculate commission amount
	commissionAmount := referral.CalculateCommission(orderAmount)

	// Create commission
	commission := &ReferralCommission{
		ID:             uuid.New(),
		TenantID:       tenantID,
		ReferralID:     referralID,
		ReferrerID:     referral.ReferrerID,
		OrderID:        &orderID,
		Amount:         commissionAmount,
		Currency:       "USD",
		CommissionRate: referral.CommissionRate,
		OrderAmount:    orderAmount,
		ConversionType: "order",
		Status:         CommissionStatusPending,
	}

	err = s.repo.CreateCommission(ctx, commission)
	if err != nil {
		return nil, fmt.Errorf("failed to create order commission: %w", err)
	}

	return commission, nil
}

// GetAffiliatePerformance retrieves performance metrics for an affiliate
func (s *ServiceImpl) GetAffiliatePerformance(ctx context.Context, tenantID, referrerID uuid.UUID) (*AffiliatePerformance, error) {
	return s.repo.GetAffiliatePerformance(ctx, tenantID, referrerID)
}

// GetTopPerformingAffiliates retrieves top performing affiliates
func (s *ServiceImpl) GetTopPerformingAffiliates(ctx context.Context, tenantID uuid.UUID, limit int) ([]*AffiliatePerformance, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	return s.repo.GetTopPerformingAffiliates(ctx, tenantID, limit)
}

// CreatePayoutBatch creates a batch payout for multiple affiliates
func (s *ServiceImpl) CreatePayoutBatch(ctx context.Context, tenantID uuid.UUID, affiliateIDs []uuid.UUID) (*AffiliatePayoutBatch, error) {
	// Calculate total amount for batch
	var totalAmount float64
	for _, affiliateID := range affiliateIDs {
		commissions, err := s.repo.GetPendingCommissionsByReferrer(ctx, tenantID, affiliateID)
		if err != nil {
			return nil, fmt.Errorf("failed to get pending commissions for affiliate %s: %w", affiliateID, err)
		}
		for _, commission := range commissions {
			totalAmount += commission.Amount
		}
	}

	// Create payout batch
	batch := &AffiliatePayoutBatch{
		ID:             uuid.New(),
		TenantID:       tenantID,
		TotalAmount:    totalAmount,
		Currency:       "USD",
		AffiliateCount: len(affiliateIDs),
		Status:         "pending",
	}
	batch.BatchNumber = batch.GenerateBatchNumber()

	err := s.repo.CreatePayoutBatch(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("failed to create payout batch: %w", err)
	}

	return batch, nil
}

// ProcessPayoutBatch processes a payout batch
func (s *ServiceImpl) ProcessPayoutBatch(ctx context.Context, tenantID, batchID uuid.UUID) error {
	batch, err := s.repo.GetPayoutBatch(ctx, tenantID, batchID)
	if err != nil {
		return fmt.Errorf("failed to get payout batch: %w", err)
	}

	if batch.Status != "pending" {
		return errors.New("batch is not in pending status")
	}

	// Mark batch as processing
	batch.MarkAsProcessing()
	err = s.repo.UpdatePayoutBatch(ctx, batch)
	if err != nil {
		return fmt.Errorf("failed to update batch status: %w", err)
	}

	// Here you would integrate with payment provider
	// For now, we'll mark as completed
	batch.MarkAsCompleted()
	return s.repo.UpdatePayoutBatch(ctx, batch)
}

// GetPayoutBatches retrieves payout batches
func (s *ServiceImpl) GetPayoutBatches(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*AffiliatePayoutBatch, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetPayoutBatches(ctx, tenantID, limit, offset)
}