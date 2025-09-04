package billing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BillingRepository defines the interface for billing data operations
type BillingRepository interface {
	// Billing Plans
	CreateBillingPlan(ctx context.Context, plan *BillingPlan) error
	GetBillingPlan(ctx context.Context, planID uuid.UUID) (*BillingPlan, error)
	GetBillingPlans(ctx context.Context, filter PlanFilter) ([]*BillingPlan, error)
	UpdateBillingPlan(ctx context.Context, plan *BillingPlan) error
	DeleteBillingPlan(ctx context.Context, planID uuid.UUID) error

	// Usage Tiers
	CreateUsageTier(ctx context.Context, tier *UsageTier) error
	GetUsageTiersByPlan(ctx context.Context, planID uuid.UUID) ([]*UsageTier, error)
	UpdateUsageTier(ctx context.Context, tier *UsageTier) error
	DeleteUsageTier(ctx context.Context, tierID uuid.UUID) error

	// Tenant Subscriptions
	CreateSubscription(ctx context.Context, subscription *TenantSubscription) error
	GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*TenantSubscription, error)
	GetSubscriptionByTenantID(ctx context.Context, tenantID uuid.UUID) (*TenantSubscription, error)
	GetSubscriptions(ctx context.Context, filter SubscriptionFilter) ([]*TenantSubscription, error)
	GetSubscriptionsWithPendingChanges(ctx context.Context) ([]*TenantSubscription, error)
	GetSubscriptionsDueForBilling(ctx context.Context, before time.Time) ([]*TenantSubscription, error)
	UpdateSubscription(ctx context.Context, subscription *TenantSubscription) error
	DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error

	// Usage Records
	CreateUsageRecord(ctx context.Context, usage *UsageRecord) error
	GetUsageRecords(ctx context.Context, filter UsageFilter) ([]*UsageRecord, error)
	GetUsageSummary(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (map[UsageType]int64, error)
	GetUsageByType(ctx context.Context, tenantID uuid.UUID, usageType UsageType, startDate, endDate time.Time) ([]*UsageRecord, error)
	UpdateUsageRecord(ctx context.Context, usage *UsageRecord) error
	DeleteUsageRecord(ctx context.Context, usageID uuid.UUID) error

	// Invoices
	CreateInvoice(ctx context.Context, invoice *Invoice) error
	GetInvoice(ctx context.Context, invoiceID uuid.UUID) (*Invoice, error)
	GetInvoices(ctx context.Context, filter InvoiceFilter) ([]*Invoice, int64, error)
	GetInvoicesByTenant(ctx context.Context, tenantID uuid.UUID, filter InvoiceFilter) ([]*Invoice, int64, error)
	GetOverdueInvoices(ctx context.Context, before time.Time) ([]*Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *Invoice) error
	DeleteInvoice(ctx context.Context, invoiceID uuid.UUID) error

	// Invoice Line Items
	CreateInvoiceLineItem(ctx context.Context, lineItem *InvoiceLineItem) error
	GetInvoiceLineItems(ctx context.Context, invoiceID uuid.UUID) ([]*InvoiceLineItem, error)
	UpdateInvoiceLineItem(ctx context.Context, lineItem *InvoiceLineItem) error
	DeleteInvoiceLineItem(ctx context.Context, lineItemID uuid.UUID) error

	// Payment Attempts
	CreatePaymentAttempt(ctx context.Context, attempt *PaymentAttempt) error
	GetPaymentAttempt(ctx context.Context, attemptID uuid.UUID) (*PaymentAttempt, error)
	GetPaymentAttemptsByInvoice(ctx context.Context, invoiceID uuid.UUID) ([]*PaymentAttempt, error)
	GetFailedPaymentAttempts(ctx context.Context, retryBefore time.Time) ([]*PaymentAttempt, error)
	UpdatePaymentAttempt(ctx context.Context, attempt *PaymentAttempt) error
	DeletePaymentAttempt(ctx context.Context, attemptID uuid.UUID) error

	// Dunning Process
	CreateDunningProcess(ctx context.Context, process *DunningProcess) error
	GetDunningProcess(ctx context.Context, processID uuid.UUID) (*DunningProcess, error)
	GetDunningProcessByInvoice(ctx context.Context, invoiceID uuid.UUID) (*DunningProcess, error)
	GetActiveDunningProcesses(ctx context.Context) ([]*DunningProcess, error)
	GetDunningProcessesDueForAction(ctx context.Context, before time.Time) ([]*DunningProcess, error)
	UpdateDunningProcess(ctx context.Context, process *DunningProcess) error
	DeleteDunningProcess(ctx context.Context, processID uuid.UUID) error

	// Dunning Actions
	CreateDunningAction(ctx context.Context, action *DunningAction) error
	GetDunningActionsByProcess(ctx context.Context, processID uuid.UUID) ([]*DunningAction, error)
	GetDunningActionsDueForExecution(ctx context.Context, before time.Time) ([]*DunningAction, error)
	UpdateDunningAction(ctx context.Context, action *DunningAction) error
	DeleteDunningAction(ctx context.Context, actionID uuid.UUID) error

	// Analytics and Reporting
	GetRevenueSummary(ctx context.Context, filter AnalyticsFilter) (*RevenueSummary, error)
	GetSubscriptionMetrics(ctx context.Context, filter AnalyticsFilter) (*SubscriptionMetrics, error)
	GetUsageMetrics(ctx context.Context, filter AnalyticsFilter) (*UsageMetrics, error)
	GetChurnMetrics(ctx context.Context, filter AnalyticsFilter) (*ChurnMetrics, error)
	GetPaymentMetrics(ctx context.Context, filter AnalyticsFilter) (*PaymentMetrics, error)

	// Utility methods
	GetNextInvoiceNumber(ctx context.Context) (string, error)
	BeginTransaction(ctx context.Context) (Transaction, error)
}

// Transaction interface for database transactions
type Transaction interface {
	Commit() error
	Rollback() error
	GetContext() context.Context
}

// Filter types for repository queries
type PlanFilter struct {
	IsActive *bool   `json:"is_active,omitempty"`
	IsPublic *bool   `json:"is_public,omitempty"`
	Currency *string `json:"currency,omitempty"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
}

type SubscriptionFilter struct {
	TenantIDs []uuid.UUID         `json:"tenant_ids,omitempty"`
	PlanIDs   []uuid.UUID         `json:"plan_ids,omitempty"`
	Status    *SubscriptionStatus `json:"status,omitempty"`
	Limit     int                 `json:"limit"`
	Offset    int                 `json:"offset"`
}

type UsageFilter struct {
	TenantID    *uuid.UUID `json:"tenant_id,omitempty"`
	UsageType   *UsageType `json:"usage_type,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	ResourceID  *string    `json:"resource_id,omitempty"`
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}

// Analytics metric types
type RevenueSummary struct {
	TotalRevenue        float64            `json:"total_revenue"`
	RecurringRevenue    float64            `json:"recurring_revenue"`
	UsageRevenue        float64            `json:"usage_revenue"`
	RefundedRevenue     float64            `json:"refunded_revenue"`
	NetRevenue          float64            `json:"net_revenue"`
	RevenueByPlan       map[string]float64 `json:"revenue_by_plan"`
	RevenueByCountry    map[string]float64 `json:"revenue_by_country"`
	AverageRevenuePerUser float64          `json:"average_revenue_per_user"`
}

type SubscriptionMetrics struct {
	TotalSubscriptions    int64              `json:"total_subscriptions"`
	ActiveSubscriptions   int64              `json:"active_subscriptions"`
	TrialSubscriptions    int64              `json:"trial_subscriptions"`
	CanceledSubscriptions int64              `json:"canceled_subscriptions"`
	NewSubscriptions      int64              `json:"new_subscriptions"`
	SubscriptionsByPlan   map[string]int64   `json:"subscriptions_by_plan"`
	SubscriptionsByStatus map[string]int64   `json:"subscriptions_by_status"`
}

type UsageMetrics struct {
	TotalUsageByType     map[UsageType]int64    `json:"total_usage_by_type"`
	AverageUsageByType   map[UsageType]float64  `json:"average_usage_by_type"`
	TopUsageTenants      []TenantUsage          `json:"top_usage_tenants"`
	UsageGrowthRate      map[UsageType]float64  `json:"usage_growth_rate"`
	UsageOverageRevenue  float64                `json:"usage_overage_revenue"`
}

type TenantUsage struct {
	TenantID    uuid.UUID           `json:"tenant_id"`
	TotalUsage  map[UsageType]int64 `json:"total_usage"`
	UsageRevenue float64            `json:"usage_revenue"`
}

type ChurnMetrics struct {
	ChurnRate             float64            `json:"churn_rate"`
	RevenueChurnRate      float64            `json:"revenue_churn_rate"`
	ChurnedSubscriptions  int64              `json:"churned_subscriptions"`
	ChurnedRevenue        float64            `json:"churned_revenue"`
	ChurnReasons          map[string]int64   `json:"churn_reasons"`
	ChurnByPlan           map[string]float64 `json:"churn_by_plan"`
	AverageLifetime       float64            `json:"average_lifetime_days"`
	CustomerLifetimeValue float64            `json:"customer_lifetime_value"`
}

type PaymentMetrics struct {
	TotalPayments        int64   `json:"total_payments"`
	SuccessfulPayments   int64   `json:"successful_payments"`
	FailedPayments       int64   `json:"failed_payments"`
	PaymentSuccessRate   float64 `json:"payment_success_rate"`
	AveragePaymentAmount float64 `json:"average_payment_amount"`
	TotalPaymentVolume   float64 `json:"total_payment_volume"`
	PaymentMethodStats   map[string]int64 `json:"payment_method_stats"`
	RefundCount          int64   `json:"refund_count"`
	RefundAmount         float64 `json:"refund_amount"`
	RefundRate           float64 `json:"refund_rate"`
	DunningRecoveryRate  float64 `json:"dunning_recovery_rate"`
}

// gormBillingRepository implements BillingRepository using GORM
type gormBillingRepository struct {
	db *gorm.DB
}

// NewBillingRepository creates a new billing repository
func NewBillingRepository(db *gorm.DB) BillingRepository {
	return &gormBillingRepository{db: db}
}

// Transaction implementation
type gormTransaction struct {
	tx *gorm.DB
}

func (t *gormTransaction) Commit() error {
	return t.tx.Commit().Error
}

func (t *gormTransaction) Rollback() error {
	return t.tx.Rollback().Error
}

func (t *gormTransaction) GetContext() context.Context {
	return t.tx.Statement.Context
}

// Billing Plans
func (r *gormBillingRepository) CreateBillingPlan(ctx context.Context, plan *BillingPlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *gormBillingRepository) GetBillingPlan(ctx context.Context, planID uuid.UUID) (*BillingPlan, error) {
	var plan BillingPlan
	err := r.db.WithContext(ctx).
		Preload("UsageTiers").
		First(&plan, "id = ?", planID).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *gormBillingRepository) GetBillingPlans(ctx context.Context, filter PlanFilter) ([]*BillingPlan, error) {
	query := r.db.WithContext(ctx).Model(&BillingPlan{})

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.IsPublic != nil {
		query = query.Where("is_public = ?", *filter.IsPublic)
	}
	if filter.Currency != nil {
		query = query.Where("currency = ?", *filter.Currency)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var plans []*BillingPlan
	err := query.Preload("UsageTiers").Find(&plans).Error
	return plans, err
}

func (r *gormBillingRepository) UpdateBillingPlan(ctx context.Context, plan *BillingPlan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

func (r *gormBillingRepository) DeleteBillingPlan(ctx context.Context, planID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&BillingPlan{}, "id = ?", planID).Error
}

// Usage Tiers
func (r *gormBillingRepository) CreateUsageTier(ctx context.Context, tier *UsageTier) error {
	return r.db.WithContext(ctx).Create(tier).Error
}

func (r *gormBillingRepository) GetUsageTiersByPlan(ctx context.Context, planID uuid.UUID) ([]*UsageTier, error) {
	var tiers []*UsageTier
	err := r.db.WithContext(ctx).
		Where("billing_plan_id = ?", planID).
		Order("usage_type, min_units").
		Find(&tiers).Error
	return tiers, err
}

func (r *gormBillingRepository) UpdateUsageTier(ctx context.Context, tier *UsageTier) error {
	return r.db.WithContext(ctx).Save(tier).Error
}

func (r *gormBillingRepository) DeleteUsageTier(ctx context.Context, tierID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&UsageTier{}, "id = ?", tierID).Error
}
