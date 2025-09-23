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

// AffiliateType represents the type of affiliate
type AffiliateType string

const (
	AffiliateTypeCustomer      AffiliateType = "customer"
	AffiliateTypeInfluencer    AffiliateType = "influencer"
	AffiliateTypePartner       AffiliateType = "partner"
	AffiliateTypeEmployee      AffiliateType = "employee"
	AffiliateTypeDigitalMarketer AffiliateType = "digital_marketer"
)

// CommissionStatus represents the status of a commission
type CommissionStatus string

const (
	CommissionStatusPending CommissionStatus = "pending"
	CommissionStatusPaid    CommissionStatus = "paid"
	CommissionStatusFailed  CommissionStatus = "failed"
	CommissionStatusVoided  CommissionStatus = "voided"
)

// Referral represents a referral relationship between users (also serves as affiliate program)
type Referral struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID       uuid.UUID      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ReferrerID     uuid.UUID      `json:"referrer_id" gorm:"type:uuid;not null;index"`
	RefereeID      *uuid.UUID     `json:"referee_id" gorm:"type:uuid;index"`
	ReferralCode   string         `json:"referral_code" gorm:"type:varchar(50);unique;not null;index"`
	Status         ReferralStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	CommissionRate float64        `json:"commission_rate" gorm:"type:decimal(5,4);not null;default:0.1"`
	AffiliateType  AffiliateType  `json:"affiliate_type" gorm:"type:varchar(20);not null;default:'customer'"`
	TrackingData   map[string]interface{} `json:"tracking_data,omitempty" gorm:"type:jsonb"`
	PayoutThreshold float64       `json:"payout_threshold" gorm:"type:decimal(10,2);default:50.0"`
	ExpiresAt      *time.Time     `json:"expires_at"`
	CompletedAt    *time.Time     `json:"completed_at"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// ReferralCommission represents a commission earned from a referral/affiliate
type ReferralCommission struct {
	ID               uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID         uuid.UUID        `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ReferralID       uuid.UUID        `json:"referral_id" gorm:"type:uuid;not null;index"`
	ReferrerID       uuid.UUID        `json:"referrer_id" gorm:"type:uuid;not null;index"`
	OrderID          *uuid.UUID       `json:"order_id,omitempty" gorm:"type:uuid;index"`
	SubscriptionID   *uuid.UUID       `json:"subscription_id,omitempty" gorm:"type:uuid;index"`
	Amount           float64          `json:"amount" gorm:"type:decimal(10,2);not null"`
	Currency         string           `json:"currency" gorm:"type:varchar(3);not null;default:'USD'"`
	CommissionRate   float64          `json:"commission_rate" gorm:"type:decimal(5,4);not null"`
	OrderAmount      float64          `json:"order_amount" gorm:"type:decimal(10,2);default:0"`
	SubscriptionAmount float64        `json:"subscription_amount" gorm:"type:decimal(10,2);default:0"`
	ConversionType   string           `json:"conversion_type" gorm:"type:varchar(20);not null;default:'order'"` // order, subscription, signup
	ClickData        map[string]interface{} `json:"click_data,omitempty" gorm:"type:jsonb"`
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

// AffiliateClick represents a click on an affiliate link
type AffiliateClick struct {
	ID           uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID     uuid.UUID              `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ReferralID   uuid.UUID              `json:"referral_id" gorm:"type:uuid;not null;index"`
	ReferrerID   uuid.UUID              `json:"referrer_id" gorm:"type:uuid;not null;index"`
	IPAddress    string                 `json:"ip_address" gorm:"type:varchar(45)"`
	UserAgent    string                 `json:"user_agent" gorm:"type:text"`
	Referrer     string                 `json:"referrer" gorm:"type:text"`
	UTMSource    string                 `json:"utm_source" gorm:"type:varchar(100)"`
	UTMMedium    string                 `json:"utm_medium" gorm:"type:varchar(100)"`
	UTMCampaign  string                 `json:"utm_campaign" gorm:"type:varchar(100)"`
	UTMTerm      string                 `json:"utm_term" gorm:"type:varchar(100)"`
	UTMContent   string                 `json:"utm_content" gorm:"type:varchar(100)"`
	DeviceType   string                 `json:"device_type" gorm:"type:varchar(20)"`
	Country      string                 `json:"country" gorm:"type:varchar(2)"`
	City         string                 `json:"city" gorm:"type:varchar(100)"`
	Converted    bool                   `json:"converted" gorm:"default:false;index"`
	ConvertedAt  *time.Time             `json:"converted_at"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt    time.Time              `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Referral *Referral `json:"referral,omitempty" gorm:"foreignKey:ReferralID"`
}

// AffiliatePerformance represents performance metrics for an affiliate
type AffiliatePerformance struct {
	TenantID         uuid.UUID `json:"tenant_id"`
	ReferrerID       uuid.UUID `json:"referrer_id"`
	TotalClicks      int64     `json:"total_clicks"`
	UniqueClicks     int64     `json:"unique_clicks"`
	Conversions      int64     `json:"conversions"`
	ConversionRate   float64   `json:"conversion_rate"`
	ClicksToday      int64     `json:"clicks_today"`
	ClicksThisWeek   int64     `json:"clicks_this_week"`
	ClicksThisMonth  int64     `json:"clicks_this_month"`
	EarningsToday    float64   `json:"earnings_today"`
	EarningsThisWeek float64   `json:"earnings_this_week"`
	EarningsThisMonth float64  `json:"earnings_this_month"`
}

// AffiliatePayoutBatch represents a batch of payouts
type AffiliatePayoutBatch struct {
	ID             uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID       uuid.UUID              `json:"tenant_id" gorm:"type:uuid;not null;index"`
	BatchNumber    string                 `json:"batch_number" gorm:"type:varchar(50);not null;uniqueIndex:idx_tenant_batch_number"`
	TotalAmount    float64                `json:"total_amount" gorm:"type:decimal(15,2);not null"`
	Currency       string                 `json:"currency" gorm:"type:varchar(3);not null;default:'USD'"`
	AffiliateCount int                    `json:"affiliate_count" gorm:"not null"`
	Status         string                 `json:"status" gorm:"type:varchar(20);not null;default:'pending'"` // pending, processing, completed, failed
	ProcessedAt    *time.Time             `json:"processed_at"`
	CompletedAt    *time.Time             `json:"completed_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt      time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Referral
func (Referral) TableName() string {
	return "referrals"
}

// TableName returns the table name for ReferralCommission
func (ReferralCommission) TableName() string {
	return "referral_commissions"
}

// TableName returns the table name for AffiliateClick
func (AffiliateClick) TableName() string {
	return "affiliate_clicks"
}

// TableName returns the table name for AffiliatePayoutBatch
func (AffiliatePayoutBatch) TableName() string {
	return "affiliate_payout_batches"
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

// Affiliate business methods

// IsDigitalMarketer checks if this is a digital marketer affiliate
func (r *Referral) IsDigitalMarketer() bool {
	return r.AffiliateType == AffiliateTypeDigitalMarketer
}

// IsInfluencer checks if this is an influencer affiliate
func (r *Referral) IsInfluencer() bool {
	return r.AffiliateType == AffiliateTypeInfluencer
}

// GenerateTrackingURL generates a tracking URL for the affiliate link
func (r *Referral) GenerateTrackingURL(baseURL, productPath string) string {
	if productPath == "" {
		return baseURL + "?ref=" + r.ReferralCode
	}
	return baseURL + productPath + "?ref=" + r.ReferralCode
}

// CanReceivePayout checks if affiliate can receive payout based on threshold
func (r *Referral) CanReceivePayout(pendingAmount float64) bool {
	return pendingAmount >= r.PayoutThreshold
}

// UpdateTrackingData updates the tracking data with new information
func (r *Referral) UpdateTrackingData(key string, value interface{}) {
	if r.TrackingData == nil {
		r.TrackingData = make(map[string]interface{})
	}
	r.TrackingData[key] = value
}

// AffiliateClick business methods

// MarkAsConverted marks the click as converted
func (ac *AffiliateClick) MarkAsConverted() {
	ac.Converted = true
	now := time.Now()
	ac.ConvertedAt = &now
}

// IsFromMobile checks if the click came from a mobile device
func (ac *AffiliateClick) IsFromMobile() bool {
	return ac.DeviceType == "mobile"
}

// IsFromDesktop checks if the click came from a desktop device
func (ac *AffiliateClick) IsFromDesktop() bool {
	return ac.DeviceType == "desktop"
}

// HasUTMParameters checks if the click has UTM tracking parameters
func (ac *AffiliateClick) HasUTMParameters() bool {
	return ac.UTMSource != "" || ac.UTMMedium != "" || ac.UTMCampaign != ""
}

// ReferralCommission business methods for affiliate tracking

// CalculateNetCommission calculates commission after platform fees
func (rc *ReferralCommission) CalculateNetCommission(platformFeeRate float64) float64 {
	platformFee := rc.Amount * platformFeeRate
	return rc.Amount - platformFee
}

// IsOrderCommission checks if this commission is from an order
func (rc *ReferralCommission) IsOrderCommission() bool {
	return rc.ConversionType == "order" && rc.OrderID != nil
}

// IsSubscriptionCommission checks if this commission is from a subscription
func (rc *ReferralCommission) IsSubscriptionCommission() bool {
	return rc.ConversionType == "subscription" && rc.SubscriptionID != nil
}

// AffiliatePayoutBatch business methods

// GenerateBatchNumber generates a unique batch number
func (ap *AffiliatePayoutBatch) GenerateBatchNumber() string {
	return "BATCH-" + time.Now().Format("20060102") + "-" + ap.ID.String()[:8]
}

// MarkAsProcessing marks the batch as processing
func (ap *AffiliatePayoutBatch) MarkAsProcessing() {
	ap.Status = "processing"
	now := time.Now()
	ap.ProcessedAt = &now
}

// MarkAsCompleted marks the batch as completed
func (ap *AffiliatePayoutBatch) MarkAsCompleted() {
	ap.Status = "completed"
	now := time.Now()
	ap.CompletedAt = &now
}

// MarkAsFailed marks the batch as failed
func (ap *AffiliatePayoutBatch) MarkAsFailed() {
	ap.Status = "failed"
}