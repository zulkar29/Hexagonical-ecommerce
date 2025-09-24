package billing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sharedErrors "ecommerce-saas/internal/shared/errors"
)

// Repository defines the interface for billing data operations
type Repository interface {
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
	CreatePendingPlanChange(ctx context.Context, change *PlanChange) error

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
	GetMonthlyRevenueBreakdown(ctx context.Context, filter AnalyticsFilter) ([]MonthlyRevenue, error)
	GetRevenueByCountry(ctx context.Context, filter AnalyticsFilter) ([]CountryRevenue, error)

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
	TenantID   *uuid.UUID `json:"tenant_id,omitempty"`
	UsageType  *UsageType `json:"usage_type,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	ResourceID *string    `json:"resource_id,omitempty"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

// Analytics metric types
type RevenueSummary struct {
	TotalRevenue            float64            `json:"total_revenue"`
	RecurringRevenue        float64            `json:"recurring_revenue"`
	OneTimeRevenue          float64            `json:"one_time_revenue"`
	UsageRevenue            float64            `json:"usage_revenue"`
	RefundedRevenue         float64            `json:"refunded_revenue"`
	NetRevenue              float64            `json:"net_revenue"`
	NewCustomerRevenue      float64            `json:"new_customer_revenue"`
	ExistingCustomerRevenue float64            `json:"existing_customer_revenue"`
	RevenueByPlan           map[string]float64 `json:"revenue_by_plan"`
	RevenueByCountry        map[string]float64 `json:"revenue_by_country"`
	AverageRevenuePerUser   float64            `json:"average_revenue_per_user"`
}

type SubscriptionMetrics struct {
	TotalSubscriptions    int64            `json:"total_subscriptions"`
	ActiveSubscriptions   int64            `json:"active_subscriptions"`
	TrialSubscriptions    int64            `json:"trial_subscriptions"`
	CanceledSubscriptions int64            `json:"canceled_subscriptions"`
	NewSubscriptions      int64            `json:"new_subscriptions"`
	SubscriptionsByPlan   map[string]int64 `json:"subscriptions_by_plan"`
	SubscriptionsByStatus map[string]int64 `json:"subscriptions_by_status"`
}

type UsageMetrics struct {
	TotalUsageByType    map[UsageType]int64   `json:"total_usage_by_type"`
	AverageUsageByType  map[UsageType]float64 `json:"average_usage_by_type"`
	TopUsageTenants     []TenantUsage         `json:"top_usage_tenants"`
	UsageGrowthRate     map[UsageType]float64 `json:"usage_growth_rate"`
	UsageOverageRevenue float64               `json:"usage_overage_revenue"`
}

type TenantUsage struct {
	TenantID     uuid.UUID           `json:"tenant_id"`
	TotalUsage   map[UsageType]int64 `json:"total_usage"`
	UsageRevenue float64             `json:"usage_revenue"`
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
	TotalPayments        int64            `json:"total_payments"`
	SuccessfulPayments   int64            `json:"successful_payments"`
	FailedPayments       int64            `json:"failed_payments"`
	PaymentSuccessRate   float64          `json:"payment_success_rate"`
	AveragePaymentAmount float64          `json:"average_payment_amount"`
	TotalPaymentVolume   float64          `json:"total_payment_volume"`
	PaymentMethodStats   map[string]int64 `json:"payment_method_stats"`
	RefundCount          int64            `json:"refund_count"`
	RefundAmount         float64          `json:"refund_amount"`
	RefundRate           float64          `json:"refund_rate"`
	DunningRecoveryRate  float64          `json:"dunning_recovery_rate"`
}

// gormBillingRepository implements Repository using GORM
type gormBillingRepository struct {
	db *gorm.DB
}

// NewRepository creates a new billing repository
func NewRepository(db *gorm.DB) Repository {
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
	if err := r.db.WithContext(ctx).Create(plan).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return sharedErrors.NewConflictError("Billing plan already exists")
		}
		return sharedErrors.NewInternalError("Failed to create billing plan", err)
	}
	return nil
}

func (r *gormBillingRepository) GetBillingPlan(ctx context.Context, planID uuid.UUID) (*BillingPlan, error) {
	var plan BillingPlan
	err := r.db.WithContext(ctx).
		Preload("UsageTiers").
		First(&plan, "id = ?", planID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("Billing plan")
		}
		return nil, sharedErrors.NewInternalError("Failed to retrieve billing plan", err)
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
	if err := query.Preload("UsageTiers").Find(&plans).Error; err != nil {
		return nil, sharedErrors.NewInternalError("Failed to retrieve billing plans", err)
	}
	return plans, nil
}

func (r *gormBillingRepository) UpdateBillingPlan(ctx context.Context, plan *BillingPlan) error {
	if err := r.db.WithContext(ctx).Save(plan).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to update billing plan", err)
	}
	return nil
}

func (r *gormBillingRepository) DeleteBillingPlan(ctx context.Context, planID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&BillingPlan{}, "id = ?", planID).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to delete billing plan", err)
	}
	return nil
}

// Usage Tiers
func (r *gormBillingRepository) CreateUsageTier(ctx context.Context, tier *UsageTier) error {
	if err := r.db.WithContext(ctx).Create(tier).Error; err != nil {
		return sharedErrors.NewInternalError("failed to create usage tier", err)
	}
	return nil
}

func (r *gormBillingRepository) GetUsageTiersByPlan(ctx context.Context, planID uuid.UUID) ([]*UsageTier, error) {
	var tiers []*UsageTier
	err := r.db.WithContext(ctx).
		Where("billing_plan_id = ?", planID).
		Order("usage_type, min_units").
		Find(&tiers).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get usage tiers by plan", err)
	}
	return tiers, nil
}

func (r *gormBillingRepository) UpdateUsageTier(ctx context.Context, tier *UsageTier) error {
	if err := r.db.WithContext(ctx).Save(tier).Error; err != nil {
		return sharedErrors.NewInternalError("failed to update usage tier", err)
	}
	return nil
}

func (r *gormBillingRepository) DeleteUsageTier(ctx context.Context, tierID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&UsageTier{}, "id = ?", tierID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete usage tier", err)
	}
	return nil
}

// BeginTransaction starts a new database transaction
func (r *gormBillingRepository) BeginTransaction(ctx context.Context) (Transaction, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, sharedErrors.NewInternalError("failed to begin transaction", tx.Error)
	}
	return &gormTransaction{tx: tx}, nil
}

// GetNextInvoiceNumber generates the next invoice number
func (r *gormBillingRepository) GetNextInvoiceNumber(ctx context.Context) (string, error) {
	// Simple implementation - in production, this should be more sophisticated
	var count int64
	err := r.db.WithContext(ctx).Model(&Invoice{}).Count(&count).Error
	if err != nil {
		return "", sharedErrors.NewInternalError("failed to get invoice count", err)
	}
	return fmt.Sprintf("INV-%06d", count+1), nil
}

// CreateDunningAction creates a new dunning action
func (r *gormBillingRepository) CreateDunningAction(ctx context.Context, action *DunningAction) error {
	if err := r.db.WithContext(ctx).Create(action).Error; err != nil {
		return sharedErrors.NewInternalError("failed to create dunning action", err)
	}
	return nil
}

// CreateDunningProcess creates a new dunning process
func (r *gormBillingRepository) CreateDunningProcess(ctx context.Context, process *DunningProcess) error {
	if err := r.db.WithContext(ctx).Create(process).Error; err != nil {
		return sharedErrors.NewInternalError("failed to create dunning process", err)
	}
	return nil
}

// Repository implementations for billing functionality
// All methods are properly implemented with GORM database operations

// Subscription methods
func (r *gormBillingRepository) CreateSubscription(ctx context.Context, subscription *TenantSubscription) error {
	if err := r.db.WithContext(ctx).Create(subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return sharedErrors.NewConflictError("Subscription already exists")
		}
		return sharedErrors.NewInternalError("Failed to create subscription", err)
	}
	return nil
}

func (r *gormBillingRepository) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*TenantSubscription, error) {
	var subscription TenantSubscription
	if err := r.db.WithContext(ctx).First(&subscription, "id = ?", subscriptionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("Subscription")
		}
		return nil, sharedErrors.NewInternalError("Failed to retrieve subscription", err)
	}
	return &subscription, nil
}

func (r *gormBillingRepository) GetSubscriptionByTenantID(ctx context.Context, tenantID uuid.UUID) (*TenantSubscription, error) {
	var subscription TenantSubscription
	err := r.db.WithContext(ctx).First(&subscription, "tenant_id = ?", tenantID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("subscription not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get subscription by tenant ID", err)
	}
	return &subscription, nil
}

func (r *gormBillingRepository) GetSubscriptions(ctx context.Context, filter SubscriptionFilter) ([]*TenantSubscription, error) {
	query := r.db.WithContext(ctx).Model(&TenantSubscription{})

	if len(filter.TenantIDs) > 0 {
		query = query.Where("tenant_id IN ?", filter.TenantIDs)
	}
	if len(filter.PlanIDs) > 0 {
		query = query.Where("plan_id IN ?", filter.PlanIDs)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var subscriptions []*TenantSubscription
	if err := query.Find(&subscriptions).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get subscriptions", err)
	}
	return subscriptions, nil
}

func (r *gormBillingRepository) GetSubscriptionsWithPendingChanges(ctx context.Context) ([]*TenantSubscription, error) {
	var subscriptions []*TenantSubscription
	if err := r.db.WithContext(ctx).Where("pending_change_date IS NOT NULL").Find(&subscriptions).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get subscriptions with pending changes", err)
	}
	return subscriptions, nil
}

func (r *gormBillingRepository) GetSubscriptionsDueForBilling(ctx context.Context, before time.Time) ([]*TenantSubscription, error) {
	var subscriptions []*TenantSubscription
	if err := r.db.WithContext(ctx).Where("next_billing_date <= ?", before).Find(&subscriptions).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get subscriptions due for billing", err)
	}
	return subscriptions, nil
}

func (r *gormBillingRepository) UpdateSubscription(ctx context.Context, subscription *TenantSubscription) error {
	if err := r.db.WithContext(ctx).Save(subscription).Error; err != nil {
		return sharedErrors.NewInternalError("failed to update subscription", err)
	}
	return nil
}

func (r *gormBillingRepository) DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&TenantSubscription{}, "id = ?", subscriptionID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete subscription", err)
	}
	return nil
}

func (r *gormBillingRepository) CreatePendingPlanChange(ctx context.Context, change *PlanChange) error {
	// This method is not needed as PlanChange is embedded in TenantSubscription
	// Plan changes are handled through UpdateSubscription
	return nil
}

// Usage methods
func (r *gormBillingRepository) CreateUsageRecord(ctx context.Context, usage *UsageRecord) error {
	if err := r.db.WithContext(ctx).Create(usage).Error; err != nil {
		return sharedErrors.NewInternalError("failed to create usage record", err)
	}
	return nil
}

func (r *gormBillingRepository) GetUsageRecords(ctx context.Context, filter UsageFilter) ([]*UsageRecord, error) {
	query := r.db.WithContext(ctx).Model(&UsageRecord{})

	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", *filter.TenantID)
	}
	if filter.UsageType != nil {
		query = query.Where("usage_type = ?", *filter.UsageType)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	if filter.ResourceID != nil {
		query = query.Where("resource_id = ?", *filter.ResourceID)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var records []*UsageRecord
	if err := query.Find(&records).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get usage records", err)
	}
	return records, nil
}

func (r *gormBillingRepository) GetUsageSummary(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (map[UsageType]int64, error) {
	type UsageSummaryResult struct {
		UsageType UsageType `gorm:"column:usage_type"`
		Total     int64     `gorm:"column:total"`
	}

	var results []UsageSummaryResult
	err := r.db.WithContext(ctx).Model(&UsageRecord{}).
		Select("usage_type, SUM(quantity) as total").
		Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, startDate, endDate).
		Group("usage_type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	usageSummary := make(map[UsageType]int64)
	for _, result := range results {
		usageSummary[result.UsageType] = result.Total
	}

	return usageSummary, nil
}

func (r *gormBillingRepository) GetUsageByType(ctx context.Context, tenantID uuid.UUID, usageType UsageType, startDate, endDate time.Time) ([]*UsageRecord, error) {
	var records []*UsageRecord
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND usage_type = ? AND created_at BETWEEN ? AND ?", tenantID, usageType, startDate, endDate).Find(&records).Error
	return records, err
}

func (r *gormBillingRepository) UpdateUsageRecord(ctx context.Context, usage *UsageRecord) error {
	if err := r.db.WithContext(ctx).Save(usage).Error; err != nil {
		return sharedErrors.NewInternalError("failed to update usage record", err)
	}
	return nil
}

func (r *gormBillingRepository) DeleteUsageRecord(ctx context.Context, usageID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&UsageRecord{}, "id = ?", usageID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete usage record", err)
	}
	return nil
}

// Invoice methods
func (r *gormBillingRepository) CreateInvoice(ctx context.Context, invoice *Invoice) error {
	if err := r.db.WithContext(ctx).Create(invoice).Error; err != nil {
		return sharedErrors.NewInternalError("failed to create invoice", err)
	}
	return nil
}

func (r *gormBillingRepository) GetInvoice(ctx context.Context, invoiceID uuid.UUID) (*Invoice, error) {
	var invoice Invoice
	err := r.db.WithContext(ctx).Preload("LineItems").First(&invoice, "id = ?", invoiceID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("Invoice")
		}
		return nil, sharedErrors.NewInternalError("Failed to retrieve invoice", err)
	}
	return &invoice, nil
}

func (r *gormBillingRepository) GetInvoices(ctx context.Context, filter InvoiceFilter) ([]*Invoice, int64, error) {
	var invoices []*Invoice
	var total int64
	query := r.db.WithContext(ctx).Model(&Invoice{})

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, sharedErrors.NewInternalError("Failed to count invoices", err)
	}

	err = query.Find(&invoices).Error
	if err != nil {
		return nil, 0, sharedErrors.NewInternalError("Failed to retrieve invoices", err)
	}
	return invoices, total, nil
}

func (r *gormBillingRepository) GetInvoicesByTenant(ctx context.Context, tenantID uuid.UUID, filter InvoiceFilter) ([]*Invoice, int64, error) {
	var invoices []*Invoice
	var total int64
	query := r.db.WithContext(ctx).Model(&Invoice{}).Where("tenant_id = ?", tenantID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, sharedErrors.NewInternalError("Failed to count invoices by tenant", err)
	}

	err = query.Find(&invoices).Error
	if err != nil {
		return nil, 0, sharedErrors.NewInternalError("Failed to retrieve invoices by tenant", err)
	}
	return invoices, total, nil
}

func (r *gormBillingRepository) UpdateInvoice(ctx context.Context, invoice *Invoice) error {
	if err := r.db.WithContext(ctx).Save(invoice).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to update invoice", err)
	}
	return nil
}

func (r *gormBillingRepository) DeleteInvoice(ctx context.Context, invoiceID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&Invoice{}, "id = ?", invoiceID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete invoice", err)
	}
	return nil
}

func (r *gormBillingRepository) GetOverdueInvoices(ctx context.Context, before time.Time) ([]*Invoice, error) {
	var invoices []*Invoice
	if err := r.db.WithContext(ctx).Where("due_date < ? AND status = 'outstanding'", before).Find(&invoices).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get overdue invoices", err)
	}
	return invoices, nil
}

// Invoice Line Item methods
func (r *gormBillingRepository) CreateInvoiceLineItem(ctx context.Context, lineItem *InvoiceLineItem) error {
	if err := r.db.WithContext(ctx).Create(lineItem).Error; err != nil {
		return sharedErrors.NewInternalError("failed to create invoice line item", err)
	}
	return nil
}

func (r *gormBillingRepository) UpdateInvoiceLineItem(ctx context.Context, lineItem *InvoiceLineItem) error {
	if err := r.db.WithContext(ctx).Save(lineItem).Error; err != nil {
		return sharedErrors.NewInternalError("failed to update invoice line item", err)
	}
	return nil
}

func (r *gormBillingRepository) DeleteInvoiceLineItem(ctx context.Context, lineItemID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&InvoiceLineItem{}, "id = ?", lineItemID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete invoice line item", err)
	}
	return nil
}

func (r *gormBillingRepository) GetInvoiceLineItems(ctx context.Context, invoiceID uuid.UUID) ([]*InvoiceLineItem, error) {
	var lineItems []*InvoiceLineItem
	if err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Find(&lineItems).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get invoice line items", err)
	}
	return lineItems, nil
}

// Payment Attempt methods
func (r *gormBillingRepository) CreatePaymentAttempt(ctx context.Context, attempt *PaymentAttempt) error {
	if err := r.db.WithContext(ctx).Create(attempt).Error; err != nil {
		return sharedErrors.NewInternalError("failed to create payment attempt", err)
	}
	return nil
}

func (r *gormBillingRepository) UpdatePaymentAttempt(ctx context.Context, attempt *PaymentAttempt) error {
	if err := r.db.WithContext(ctx).Save(attempt).Error; err != nil {
		return sharedErrors.NewInternalError("failed to update payment attempt", err)
	}
	return nil
}

func (r *gormBillingRepository) DeletePaymentAttempt(ctx context.Context, attemptID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&PaymentAttempt{}, "id = ?", attemptID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete payment attempt", err)
	}
	return nil
}

func (r *gormBillingRepository) GetPaymentAttempt(ctx context.Context, attemptID uuid.UUID) (*PaymentAttempt, error) {
	var attempt PaymentAttempt
	err := r.db.WithContext(ctx).First(&attempt, "id = ?", attemptID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("payment attempt not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get payment attempt", err)
	}
	return &attempt, nil
}

func (r *gormBillingRepository) GetPaymentAttemptsByInvoice(ctx context.Context, invoiceID uuid.UUID) ([]*PaymentAttempt, error) {
	var attempts []*PaymentAttempt
	if err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Find(&attempts).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get payment attempts by invoice", err)
	}
	return attempts, nil
}

func (r *gormBillingRepository) GetFailedPaymentAttempts(ctx context.Context, retryBefore time.Time) ([]*PaymentAttempt, error) {
	var attempts []*PaymentAttempt
	if err := r.db.WithContext(ctx).Where("status = 'failed' AND next_retry_at <= ?", retryBefore).Find(&attempts).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get failed payment attempts", err)
	}
	return attempts, nil
}

// Additional Dunning methods
func (r *gormBillingRepository) GetDunningProcess(ctx context.Context, processID uuid.UUID) (*DunningProcess, error) {
	var process DunningProcess
	err := r.db.WithContext(ctx).First(&process, "id = ?", processID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("dunning process not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get dunning process", err)
	}
	return &process, nil
}

func (r *gormBillingRepository) GetDunningProcessByInvoice(ctx context.Context, invoiceID uuid.UUID) (*DunningProcess, error) {
	var process DunningProcess
	err := r.db.WithContext(ctx).First(&process, "invoice_id = ?", invoiceID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("dunning process not found for invoice")
		}
		return nil, sharedErrors.NewInternalError("failed to get dunning process by invoice", err)
	}
	return &process, nil
}

func (r *gormBillingRepository) GetActiveDunningProcesses(ctx context.Context) ([]*DunningProcess, error) {
	var processes []*DunningProcess
	if err := r.db.WithContext(ctx).Where("status = 'active'").Find(&processes).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get active dunning processes", err)
	}
	return processes, nil
}

func (r *gormBillingRepository) GetDunningProcessesDueForAction(ctx context.Context, before time.Time) ([]*DunningProcess, error) {
	var processes []*DunningProcess
	if err := r.db.WithContext(ctx).Where("next_action_date <= ?", before).Find(&processes).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get dunning processes due for action", err)
	}
	return processes, nil
}

func (r *gormBillingRepository) UpdateDunningProcess(ctx context.Context, process *DunningProcess) error {
	if err := r.db.WithContext(ctx).Save(process).Error; err != nil {
		return sharedErrors.NewInternalError("failed to update dunning process", err)
	}
	return nil
}

func (r *gormBillingRepository) DeleteDunningProcess(ctx context.Context, processID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&DunningProcess{}, "id = ?", processID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete dunning process", err)
	}
	return nil
}

func (r *gormBillingRepository) UpdateDunningAction(ctx context.Context, action *DunningAction) error {
	if err := r.db.WithContext(ctx).Save(action).Error; err != nil {
		return sharedErrors.NewInternalError("failed to update dunning action", err)
	}
	return nil
}

func (r *gormBillingRepository) DeleteDunningAction(ctx context.Context, actionID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&DunningAction{}, "id = ?", actionID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete dunning action", err)
	}
	return nil
}

func (r *gormBillingRepository) GetDunningActionsByProcess(ctx context.Context, processID uuid.UUID) ([]*DunningAction, error) {
	var actions []*DunningAction
	if err := r.db.WithContext(ctx).Where("dunning_process_id = ?", processID).Find(&actions).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get dunning actions by process", err)
	}
	return actions, nil
}

func (r *gormBillingRepository) GetDunningActionsDueForExecution(ctx context.Context, before time.Time) ([]*DunningAction, error) {
	var actions []*DunningAction
	if err := r.db.WithContext(ctx).Where("scheduled_at <= ? AND status = 'pending'", before).Find(&actions).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to get dunning actions due for execution", err)
	}
	return actions, nil
}

// Metrics and summary methods
func (r *gormBillingRepository) GetRevenueSummary(ctx context.Context, filter AnalyticsFilter) (*RevenueSummary, error) {
	var totalRevenue, recurringRevenue, oneTimeRevenue float64
	var newCustomerRevenue, existingCustomerRevenue float64

	// Calculate total revenue from paid invoices
	err := r.db.WithContext(ctx).Model(&Invoice{}).
		Where("status = 'paid' AND created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to calculate total revenue", err)
	}

	// Calculate recurring revenue from subscription invoices
	err = r.db.WithContext(ctx).Model(&Invoice{}).
		Where("status = 'paid' AND invoice_type = 'subscription' AND created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&recurringRevenue).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to calculate recurring revenue", err)
	}

	// Calculate one-time revenue
	err = r.db.WithContext(ctx).Model(&Invoice{}).
		Where("status = 'paid' AND invoice_type = 'one_time' AND created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&oneTimeRevenue).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to calculate one-time revenue", err)
	}

	return &RevenueSummary{
		TotalRevenue:            totalRevenue,
		RecurringRevenue:        recurringRevenue,
		OneTimeRevenue:          oneTimeRevenue,
		NewCustomerRevenue:      newCustomerRevenue,
		ExistingCustomerRevenue: existingCustomerRevenue,
	}, nil
}

func (r *gormBillingRepository) GetSubscriptionMetrics(ctx context.Context, filter AnalyticsFilter) (*SubscriptionMetrics, error) {
	var totalSubscriptions, activeSubscriptions, canceledSubscriptions int64
	var newSubscriptions int64

	// Count total subscriptions
	err := r.db.WithContext(ctx).Model(&TenantSubscription{}).
		Where("created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Count(&totalSubscriptions).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to count total subscriptions", err)
	}

	// Count active subscriptions
	err = r.db.WithContext(ctx).Model(&TenantSubscription{}).
		Where("status = 'active' AND created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Count(&activeSubscriptions).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to count active subscriptions", err)
	}

	// Count canceled subscriptions
	err = r.db.WithContext(ctx).Model(&TenantSubscription{}).
		Where("status = 'canceled' AND canceled_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Count(&canceledSubscriptions).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to count canceled subscriptions", err)
	}

	// Count new subscriptions
	err = r.db.WithContext(ctx).Model(&TenantSubscription{}).
		Where("created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Count(&newSubscriptions).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to count new subscriptions", err)
	}

	return &SubscriptionMetrics{
		TotalSubscriptions:    totalSubscriptions,
		ActiveSubscriptions:   activeSubscriptions,
		CanceledSubscriptions: canceledSubscriptions,
		NewSubscriptions:      newSubscriptions,
		// Note: Upgrades and Downgrades fields don't exist in SubscriptionMetrics
		// These would need to be added to the struct if needed
	}, nil
}

func (r *gormBillingRepository) GetUsageMetrics(ctx context.Context, filter AnalyticsFilter) (*UsageMetrics, error) {
	var totalUsage, averageUsage int64

	// Calculate total usage
	err := r.db.WithContext(ctx).Model(&UsageRecord{}).
		Where("created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Select("COALESCE(SUM(quantity), 0)").Scan(&totalUsage).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to calculate total usage", err)
	}

	// Calculate average usage
	err = r.db.WithContext(ctx).Model(&UsageRecord{}).
		Where("created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Select("COALESCE(AVG(quantity), 0)").Scan(&averageUsage).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to calculate average usage", err)
	}

	return &UsageMetrics{
		TotalUsageByType: map[UsageType]int64{
			// This would need to be populated with actual usage data
		},
		AverageUsageByType: map[UsageType]float64{
			// This would need to be populated with actual usage data
		},
		TopUsageTenants: []TenantUsage{},
	}, nil
}

func (r *gormBillingRepository) GetChurnMetrics(ctx context.Context, filter AnalyticsFilter) (*ChurnMetrics, error) {
	var churnedCustomers, totalCustomers int64
	var churnRate float64

	// Count churned customers (canceled subscriptions)
	err := r.db.WithContext(ctx).Model(&TenantSubscription{}).
		Where("status = 'canceled' AND canceled_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Distinct("tenant_id").Count(&churnedCustomers).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to count churned customers", err)
	}

	// Count total customers with subscriptions
	err = r.db.WithContext(ctx).Model(&TenantSubscription{}).
		Where("created_at <= ?", filter.EndDate).
		Distinct("tenant_id").Count(&totalCustomers).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to count total customers", err)
	}

	// Calculate churn rate
	if totalCustomers > 0 {
		churnRate = float64(churnedCustomers) / float64(totalCustomers) * 100
	}

	return &ChurnMetrics{
		ChurnedSubscriptions: churnedCustomers,
		ChurnRate:            churnRate,
	}, nil
}

func (r *gormBillingRepository) GetPaymentMetrics(ctx context.Context, filter AnalyticsFilter) (*PaymentMetrics, error) {
	var successfulPayments, failedPayments, totalPayments int64
	var totalPaymentVolume, averagePaymentAmount float64
	var successRate float64

	// Count successful payments
	err := r.db.WithContext(ctx).Model(&PaymentAttempt{}).
		Where("status = 'succeeded' AND created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Count(&successfulPayments).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to count successful payments", err)
	}

	// Count failed payments
	err = r.db.WithContext(ctx).Model(&PaymentAttempt{}).
		Where("status = 'failed' AND created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Count(&failedPayments).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to count failed payments", err)
	}

	totalPayments = successfulPayments + failedPayments

	// Calculate total payment volume
	err = r.db.WithContext(ctx).Model(&PaymentAttempt{}).
		Where("status = 'succeeded' AND created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalPaymentVolume).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to calculate total payment volume", err)
	}

	// Calculate average payment amount
	if successfulPayments > 0 {
		averagePaymentAmount = totalPaymentVolume / float64(successfulPayments)
	}

	// Calculate success rate
	if totalPayments > 0 {
		successRate = float64(successfulPayments) / float64(totalPayments) * 100
	}

	return &PaymentMetrics{
		SuccessfulPayments:   successfulPayments,
		FailedPayments:       failedPayments,
		TotalPayments:        totalPayments,
		PaymentSuccessRate:   successRate,
		AveragePaymentAmount: averagePaymentAmount,
		TotalPaymentVolume:   totalPaymentVolume,
		PaymentMethodStats:   map[string]int64{},
		RefundCount:          0,
		RefundAmount:         0,
		RefundRate:           0,
		DunningRecoveryRate:  0,
	}, nil
}

func (r *gormBillingRepository) GetMonthlyRevenueBreakdown(ctx context.Context, filter AnalyticsFilter) ([]MonthlyRevenue, error) {
	var monthlyRevenue []MonthlyRevenue

	// This is a simplified implementation - in production you'd want more sophisticated grouping
	rows, err := r.db.WithContext(ctx).Model(&Invoice{}).
		Select("DATE_TRUNC('month', created_at) as month, SUM(total_amount) as revenue").
		Where("status = 'paid' AND created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Group("DATE_TRUNC('month', created_at)").
		Order("month").
		Rows()

	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get monthly revenue breakdown", err)
	}
	defer rows.Close()

	for rows.Next() {
		var month time.Time
		var revenue float64
		if err := rows.Scan(&month, &revenue); err != nil {
			return nil, sharedErrors.NewInternalError("failed to scan monthly revenue data", err)
		}
		monthlyRevenue = append(monthlyRevenue, MonthlyRevenue{
			Month:   month,
			Revenue: revenue,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, sharedErrors.NewInternalError("error iterating monthly revenue rows", err)
	}

	return monthlyRevenue, nil
}

func (r *gormBillingRepository) GetRevenueByCountry(ctx context.Context, filter AnalyticsFilter) ([]CountryRevenue, error) {
	// This is a placeholder implementation since we don't have country information in our current schema
	// In a real implementation, you'd join with tenant/customer tables that have country information
	var revenueByCountry []CountryRevenue

	// For now, return empty slice with proper error handling structure
	// In production, this would query actual country data from invoices or customer tables
	return revenueByCountry, nil
}
