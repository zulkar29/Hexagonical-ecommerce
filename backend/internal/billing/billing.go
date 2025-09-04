package billing

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BillingCycle represents billing frequency
type BillingCycle string

const (
	BillingCycleMonthly   BillingCycle = "monthly"
	BillingCycleQuarterly BillingCycle = "quarterly"
	BillingCycleYearly    BillingCycle = "yearly"
)

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusPending   SubscriptionStatus = "pending"
	SubscriptionStatusSuspended SubscriptionStatus = "suspended"
	SubscriptionStatusCanceled  SubscriptionStatus = "canceled"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
)

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusPending   InvoiceStatus = "pending"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
	InvoiceStatusVoided    InvoiceStatus = "voided"
	InvoiceStatusRefunded  InvoiceStatus = "refunded"
)

// PaymentStatus represents payment attempt status
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRetrying  PaymentStatus = "retrying"
	PaymentStatusAbandoned PaymentStatus = "abandoned"
)

// UsageType represents different types of billable usage
type UsageType string

const (
	UsageTypeAPIRequests     UsageType = "api_requests"
	UsageTypeStorageGB       UsageType = "storage_gb"
	UsageTypeBandwidthGB     UsageType = "bandwidth_gb"
	UsageTypeOrders          UsageType = "orders"
	UsageTypeProducts        UsageType = "products"
	UsageTypeUsers           UsageType = "users"
	UsageTypeEmailsSent      UsageType = "emails_sent"
	UsageTypeTransactions    UsageType = "transactions"
)

// BillingPlan represents a subscription plan with pricing tiers
type BillingPlan struct {
	ID          uuid.UUID `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	
	// Pricing
	BasePrice     float64      `json:"base_price" gorm:"not null"`
	Currency      string       `json:"currency" gorm:"default:BDT"`
	BillingCycle  BillingCycle `json:"billing_cycle" gorm:"default:monthly"`
	
	// Plan limits
	Limits        map[string]interface{} `json:"limits" gorm:"serializer:json"`
	Features      []string               `json:"features" gorm:"serializer:json"`
	
	// Usage-based pricing
	UsageTiers    []UsageTier `json:"usage_tiers" gorm:"foreignKey:BillingPlanID"`
	
	// Plan settings
	TrialPeriodDays int  `json:"trial_period_days" gorm:"default:0"`
	IsActive        bool `json:"is_active" gorm:"default:true"`
	IsPublic        bool `json:"is_public" gorm:"default:true"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// UsageTier represents usage-based pricing tiers
type UsageTier struct {
	ID            uuid.UUID `json:"id" gorm:"primarykey"`
	BillingPlanID uuid.UUID `json:"billing_plan_id" gorm:"not null;index"`
	
	UsageType     UsageType `json:"usage_type" gorm:"not null"`
	MinUnits      int64     `json:"min_units" gorm:"not null"`
	MaxUnits      *int64    `json:"max_units,omitempty"` // null for unlimited
	PricePerUnit  float64   `json:"price_per_unit" gorm:"not null"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TenantSubscription represents a tenant's subscription to a billing plan
type TenantSubscription struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;uniqueIndex:idx_tenant_subscription"`
	PlanID   uuid.UUID `json:"plan_id" gorm:"not null;index"`
	
	// Subscription details
	Status            SubscriptionStatus `json:"status" gorm:"default:pending"`
	BillingCycle      BillingCycle       `json:"billing_cycle"`
	
	// Billing periods
	CurrentPeriodStart time.Time  `json:"current_period_start" gorm:"not null"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end" gorm:"not null"`
	TrialEnd           *time.Time `json:"trial_end,omitempty"`
	
	// Pricing
	BaseAmount        float64 `json:"base_amount" gorm:"not null"`
	Currency          string  `json:"currency" gorm:"default:BDT"`
	
	// Plan change management
	PendingPlanChange *PlanChange `json:"pending_plan_change,omitempty" gorm:"embedded;embeddedPrefix:pending_"`
	
	// Payment and billing
	PaymentMethodID   *string    `json:"payment_method_id,omitempty"`
	NextBillingDate   time.Time  `json:"next_billing_date" gorm:"not null"`
	CanceledAt        *time.Time `json:"canceled_at,omitempty"`
	CancellationReason string    `json:"cancellation_reason,omitempty"`
	
	// Relationships
	Plan    BillingPlan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	Invoices []Invoice   `json:"invoices,omitempty" gorm:"foreignKey:SubscriptionID"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// PlanChange represents a pending plan change
type PlanChange struct {
	NewPlanID       uuid.UUID `json:"new_plan_id"`
	EffectiveDate   time.Time `json:"effective_date"`
	ProrationAmount float64   `json:"proration_amount"`
	ChangeReason    string    `json:"change_reason"`
}

// UsageRecord represents billable usage by a tenant
type UsageRecord struct {
	ID         uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID   uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	UsageType  UsageType `json:"usage_type" gorm:"not null;index"`
	Quantity   int64     `json:"quantity" gorm:"not null"`
	Units      string    `json:"units" gorm:"not null"` // requests, GB, users, etc.
	
	// Billing period association
	BillingPeriodStart time.Time `json:"billing_period_start" gorm:"not null;index"`
	BillingPeriodEnd   time.Time `json:"billing_period_end" gorm:"not null;index"`
	
	// Usage metadata
	ResourceID   *string                `json:"resource_id,omitempty"` // ID of resource that generated usage
	Metadata     map[string]interface{} `json:"metadata,omitempty" gorm:"serializer:json"`
	
	RecordedAt time.Time `json:"recorded_at" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}

// Invoice represents a billing invoice
type Invoice struct {
	ID             uuid.UUID  `json:"id" gorm:"primarykey"`
	TenantID       uuid.UUID  `json:"tenant_id" gorm:"not null;index"`
	SubscriptionID uuid.UUID  `json:"subscription_id" gorm:"not null;index"`
	
	// Invoice details
	InvoiceNumber  string        `json:"invoice_number" gorm:"uniqueIndex;not null"`
	Status         InvoiceStatus `json:"status" gorm:"default:draft"`
	
	// Billing period
	PeriodStart time.Time `json:"period_start" gorm:"not null"`
	PeriodEnd   time.Time `json:"period_end" gorm:"not null"`
	
	// Amounts
	SubtotalAmount float64 `json:"subtotal_amount" gorm:"not null"`
	TaxAmount      float64 `json:"tax_amount" gorm:"default:0"`
	TotalAmount    float64 `json:"total_amount" gorm:"not null"`
	PaidAmount     float64 `json:"paid_amount" gorm:"default:0"`
	Currency       string  `json:"currency" gorm:"default:BDT"`
	
	// Payment details
	DueDate       time.Time  `json:"due_date" gorm:"not null"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	PaymentMethod *string    `json:"payment_method,omitempty"`
	
	// Invoice items
	LineItems []InvoiceLineItem `json:"line_items" gorm:"foreignKey:InvoiceID"`
	
	// Relationships
	Subscription TenantSubscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
	PaymentAttempts []PaymentAttempt `json:"payment_attempts,omitempty" gorm:"foreignKey:InvoiceID"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// InvoiceLineItem represents individual items on an invoice
type InvoiceLineItem struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	InvoiceID uuid.UUID `json:"invoice_id" gorm:"not null;index"`
	
	// Item details
	Description   string  `json:"description" gorm:"not null"`
	Quantity      int64   `json:"quantity" gorm:"not null"`
	UnitPrice     float64 `json:"unit_price" gorm:"not null"`
	TotalPrice    float64 `json:"total_price" gorm:"not null"`
	
	// Item type and metadata
	ItemType      string                 `json:"item_type" gorm:"not null"` // subscription, usage, addon, discount
	UsageType     *UsageType             `json:"usage_type,omitempty"`
	PeriodStart   *time.Time             `json:"period_start,omitempty"`
	PeriodEnd     *time.Time             `json:"period_end,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" gorm:"serializer:json"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaymentAttempt represents an attempt to collect payment for an invoice
type PaymentAttempt struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	InvoiceID uuid.UUID `json:"invoice_id" gorm:"not null;index"`
	TenantID  uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Payment details
	Status        PaymentStatus `json:"status" gorm:"default:pending"`
	Amount        float64       `json:"amount" gorm:"not null"`
	Currency      string        `json:"currency" gorm:"default:BDT"`
	PaymentMethod string        `json:"payment_method" gorm:"not null"`
	
	// External payment provider details
	ProviderID         *string                `json:"provider_id,omitempty"` // Stripe, bKash, etc.
	ProviderChargeID   *string                `json:"provider_charge_id,omitempty"`
	ProviderResponse   map[string]interface{} `json:"provider_response,omitempty" gorm:"serializer:json"`
	
	// Attempt details
	AttemptedAt time.Time  `json:"attempted_at" gorm:"not null"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	FailureReason *string  `json:"failure_reason,omitempty"`
	
	// Retry logic
	RetryCount    int        `json:"retry_count" gorm:"default:0"`
	NextRetryAt   *time.Time `json:"next_retry_at,omitempty"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DunningProcess represents the process of collecting overdue payments
type DunningProcess struct {
	ID             uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID       uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	InvoiceID      uuid.UUID `json:"invoice_id" gorm:"not null;index"`
	SubscriptionID uuid.UUID `json:"subscription_id" gorm:"not null;index"`
	
	// Dunning state
	CurrentStep   int       `json:"current_step" gorm:"default:1"`
	TotalSteps    int       `json:"total_steps" gorm:"not null"`
	IsCompleted   bool      `json:"is_completed" gorm:"default:false"`
	IsAbandoned   bool      `json:"is_abandoned" gorm:"default:false"`
	
	// Timeline
	StartedAt     time.Time  `json:"started_at" gorm:"not null"`
	NextActionAt  *time.Time `json:"next_action_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	
	// Actions taken
	EmailsSent         int `json:"emails_sent" gorm:"default:0"`
	PaymentAttempts    int `json:"payment_attempts" gorm:"default:0"`
	ServiceSuspended   bool `json:"service_suspended" gorm:"default:false"`
	ServiceSuspendedAt *time.Time `json:"service_suspended_at,omitempty"`
	
	// Relationships
	Invoice      Invoice            `json:"invoice,omitempty" gorm:"foreignKey:InvoiceID"`
	Subscription TenantSubscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
	Actions      []DunningAction    `json:"actions,omitempty" gorm:"foreignKey:DunningProcessID"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DunningAction represents individual actions taken during the dunning process
type DunningAction struct {
	ID               uuid.UUID `json:"id" gorm:"primarykey"`
	DunningProcessID uuid.UUID `json:"dunning_process_id" gorm:"not null;index"`
	
	// Action details
	ActionType   string                 `json:"action_type" gorm:"not null"` // email, suspend, retry_payment, cancel
	StepNumber   int                    `json:"step_number" gorm:"not null"`
	Description  string                 `json:"description" gorm:"not null"`
	
	// Action result
	Status       string                 `json:"status" gorm:"not null"` // pending, completed, failed
	Result       map[string]interface{} `json:"result,omitempty" gorm:"serializer:json"`
	ErrorMessage *string                `json:"error_message,omitempty"`
	
	// Timing
	ScheduledAt  time.Time  `json:"scheduled_at" gorm:"not null"`
	ExecutedAt   *time.Time `json:"executed_at,omitempty"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Business Logic Methods

// IsTrialActive checks if the subscription is in trial period
func (s *TenantSubscription) IsTrialActive() bool {
	if s.TrialEnd == nil {
		return false
	}
	return time.Now().Before(*s.TrialEnd)
}

// DaysUntilRenewal returns days until next billing date
func (s *TenantSubscription) DaysUntilRenewal() int {
	duration := time.Until(s.NextBillingDate)
	return int(duration.Hours() / 24)
}

// CanUpgrade checks if subscription can be upgraded
func (s *TenantSubscription) CanUpgrade() bool {
	return s.Status == SubscriptionStatusActive && s.PendingPlanChange == nil
}

// CanDowngrade checks if subscription can be downgraded
func (s *TenantSubscription) CanDowngrade() bool {
	return s.Status == SubscriptionStatusActive && s.PendingPlanChange == nil
}

// IsOverdue checks if an invoice is overdue
func (i *Invoice) IsOverdue() bool {
	return i.Status == InvoiceStatusPending && time.Now().After(i.DueDate)
}

// RemainingAmount returns the unpaid amount on an invoice
func (i *Invoice) RemainingAmount() float64 {
	return i.TotalAmount - i.PaidAmount
}

// CanRetry checks if a payment attempt can be retried
func (p *PaymentAttempt) CanRetry() bool {
	return p.Status == PaymentStatusFailed && p.RetryCount < 3 && 
		   (p.NextRetryAt == nil || time.Now().After(*p.NextRetryAt))
}

// ShouldSuspendService determines if service should be suspended based on dunning state
func (d *DunningProcess) ShouldSuspendService() bool {
	return d.CurrentStep >= 3 && !d.ServiceSuspended && !d.IsCompleted
}

// ShouldCancelSubscription determines if subscription should be canceled
func (d *DunningProcess) ShouldCancelSubscription() bool {
	return d.CurrentStep >= d.TotalSteps && !d.IsCompleted
}

// CalculateProration calculates prorated amount for plan changes
func CalculateProration(oldAmount, newAmount float64, daysUsed, totalDays int) float64 {
	if totalDays <= 0 {
		return 0
	}
	
	daysRemaining := totalDays - daysUsed
	if daysRemaining <= 0 {
		return 0
	}
	
	// Credit for unused portion of old plan
	oldCredit := oldAmount * float64(daysRemaining) / float64(totalDays)
	
	// Charge for new plan for remaining days
	newCharge := newAmount * float64(daysRemaining) / float64(totalDays)
	
	return newCharge - oldCredit
}

// GetBillingCycleDays returns the number of days in a billing cycle
func (bc BillingCycle) GetDays() int {
	switch bc {
	case BillingCycleMonthly:
		return 30
	case BillingCycleQuarterly:
		return 90
	case BillingCycleYearly:
		return 365
	default:
		return 30
	}
}

// GenerateInvoiceNumber generates a unique invoice number
func GenerateInvoiceNumber(tenantID uuid.UUID) string {
	// Format: INV-YYYYMMDD-TENANT_SHORT-XXXX
	now := time.Now()
	tenantShort := tenantID.String()[:8]
	timestamp := now.Format("20060102")
	sequence := now.Unix() % 10000
	
	return fmt.Sprintf("INV-%s-%s-%04d", timestamp, tenantShort, sequence)
}

// TODO: Add more business logic methods
// - CalculateUsageCharges(usageRecords []UsageRecord, usageTiers []UsageTier) float64
// - DetermineNextBillingDate(currentDate time.Time, cycle BillingCycle) time.Time
// - ValidateSubscriptionLimits(subscription *TenantSubscription, usage map[UsageType]int64) error
// - CalculateTax(amount float64, tenantLocation string) float64