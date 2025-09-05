package billing

import (
	"context"
	"fmt"
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

// Repository is an alias for BillingRepository
type Repository = BillingRepository

// NewRepository creates a new repository (alias for NewBillingRepository)
func NewRepository(db *gorm.DB) Repository {
	return NewBillingRepository(db)
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

// BeginTransaction starts a new database transaction
func (r *gormBillingRepository) BeginTransaction(ctx context.Context) (Transaction, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &gormTransaction{tx: tx}, nil
}

// GetNextInvoiceNumber generates the next invoice number
func (r *gormBillingRepository) GetNextInvoiceNumber(ctx context.Context) (string, error) {
	// Simple implementation - in production, this should be more sophisticated
	var count int64
	err := r.db.WithContext(ctx).Model(&Invoice{}).Count(&count).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("INV-%06d", count+1), nil
}

// CreateDunningAction creates a new dunning action
func (r *gormBillingRepository) CreateDunningAction(ctx context.Context, action *DunningAction) error {
	return r.db.WithContext(ctx).Create(action).Error
}

// CreateDunningProcess creates a new dunning process
func (r *gormBillingRepository) CreateDunningProcess(ctx context.Context, process *DunningProcess) error {
	return r.db.WithContext(ctx).Create(process).Error
}

// Stub implementations for missing methods to fix compilation
// TODO: Implement these methods properly

// Subscription methods
func (r *gormBillingRepository) CreateSubscription(ctx context.Context, subscription *TenantSubscription) error {
	return r.db.WithContext(ctx).Create(subscription).Error
}

func (r *gormBillingRepository) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*TenantSubscription, error) {
	var subscription TenantSubscription
	err := r.db.WithContext(ctx).First(&subscription, "id = ?", subscriptionID).Error
	return &subscription, err
}

func (r *gormBillingRepository) GetSubscriptionByTenantID(ctx context.Context, tenantID uuid.UUID) (*TenantSubscription, error) {
	var subscription TenantSubscription
	err := r.db.WithContext(ctx).First(&subscription, "tenant_id = ?", tenantID).Error
	return &subscription, err
}

func (r *gormBillingRepository) GetSubscriptions(ctx context.Context, filter SubscriptionFilter) ([]*TenantSubscription, error) {
	var subscriptions []*TenantSubscription
	err := r.db.WithContext(ctx).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *gormBillingRepository) GetSubscriptionsWithPendingChanges(ctx context.Context) ([]*TenantSubscription, error) {
	var subscriptions []*TenantSubscription
	err := r.db.WithContext(ctx).Where("pending_change_date IS NOT NULL").Find(&subscriptions).Error
	return subscriptions, err
}

func (r *gormBillingRepository) GetSubscriptionsDueForBilling(ctx context.Context, before time.Time) ([]*TenantSubscription, error) {
	var subscriptions []*TenantSubscription
	err := r.db.WithContext(ctx).Where("next_billing_date <= ?", before).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *gormBillingRepository) UpdateSubscription(ctx context.Context, subscription *TenantSubscription) error {
	return r.db.WithContext(ctx).Save(subscription).Error
}

func (r *gormBillingRepository) DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&TenantSubscription{}, "id = ?", subscriptionID).Error
}

// Usage methods
func (r *gormBillingRepository) CreateUsageRecord(ctx context.Context, usage *UsageRecord) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

func (r *gormBillingRepository) GetUsageRecords(ctx context.Context, filter UsageFilter) ([]*UsageRecord, error) {
	var records []*UsageRecord
	err := r.db.WithContext(ctx).Find(&records).Error
	return records, err
}

func (r *gormBillingRepository) GetUsageSummary(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (map[UsageType]int64, error) {
	return make(map[UsageType]int64), nil
}

func (r *gormBillingRepository) GetUsageByType(ctx context.Context, tenantID uuid.UUID, usageType UsageType, startDate, endDate time.Time) ([]*UsageRecord, error) {
	var records []*UsageRecord
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND usage_type = ? AND created_at BETWEEN ? AND ?", tenantID, usageType, startDate, endDate).Find(&records).Error
	return records, err
}

func (r *gormBillingRepository) UpdateUsageRecord(ctx context.Context, usage *UsageRecord) error {
	return r.db.WithContext(ctx).Save(usage).Error
}

func (r *gormBillingRepository) DeleteUsageRecord(ctx context.Context, usageID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&UsageRecord{}, "id = ?", usageID).Error
}

// Invoice methods
func (r *gormBillingRepository) CreateInvoice(ctx context.Context, invoice *Invoice) error {
	return r.db.WithContext(ctx).Create(invoice).Error
}

func (r *gormBillingRepository) GetInvoice(ctx context.Context, invoiceID uuid.UUID) (*Invoice, error) {
	var invoice Invoice
	err := r.db.WithContext(ctx).Preload("LineItems").First(&invoice, "id = ?", invoiceID).Error
	return &invoice, err
}

func (r *gormBillingRepository) GetInvoices(ctx context.Context, filter InvoiceFilter) ([]*Invoice, int64, error) {
	var invoices []*Invoice
	var total int64
	query := r.db.WithContext(ctx).Model(&Invoice{})
	
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = query.Find(&invoices).Error
	return invoices, total, err
}

func (r *gormBillingRepository) GetInvoicesByTenant(ctx context.Context, tenantID uuid.UUID, filter InvoiceFilter) ([]*Invoice, int64, error) {
	var invoices []*Invoice
	var total int64
	query := r.db.WithContext(ctx).Model(&Invoice{}).Where("tenant_id = ?", tenantID)
	
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = query.Find(&invoices).Error
	return invoices, total, err
}

func (r *gormBillingRepository) UpdateInvoice(ctx context.Context, invoice *Invoice) error {
	return r.db.WithContext(ctx).Save(invoice).Error
}

func (r *gormBillingRepository) DeleteInvoice(ctx context.Context, invoiceID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Invoice{}, "id = ?", invoiceID).Error
}

func (r *gormBillingRepository) GetOverdueInvoices(ctx context.Context, before time.Time) ([]*Invoice, error) {
	var invoices []*Invoice
	err := r.db.WithContext(ctx).Where("due_date < ? AND status = 'outstanding'", before).Find(&invoices).Error
	return invoices, err
}

// Invoice Line Item methods
func (r *gormBillingRepository) CreateInvoiceLineItem(ctx context.Context, lineItem *InvoiceLineItem) error {
	return r.db.WithContext(ctx).Create(lineItem).Error
}

func (r *gormBillingRepository) UpdateInvoiceLineItem(ctx context.Context, lineItem *InvoiceLineItem) error {
	return r.db.WithContext(ctx).Save(lineItem).Error
}

func (r *gormBillingRepository) DeleteInvoiceLineItem(ctx context.Context, lineItemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&InvoiceLineItem{}, "id = ?", lineItemID).Error
}

func (r *gormBillingRepository) GetInvoiceLineItems(ctx context.Context, invoiceID uuid.UUID) ([]*InvoiceLineItem, error) {
	var lineItems []*InvoiceLineItem
	err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Find(&lineItems).Error
	return lineItems, err
}

// Payment Attempt methods
func (r *gormBillingRepository) CreatePaymentAttempt(ctx context.Context, attempt *PaymentAttempt) error {
	return r.db.WithContext(ctx).Create(attempt).Error
}

func (r *gormBillingRepository) UpdatePaymentAttempt(ctx context.Context, attempt *PaymentAttempt) error {
	return r.db.WithContext(ctx).Save(attempt).Error
}

func (r *gormBillingRepository) DeletePaymentAttempt(ctx context.Context, attemptID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&PaymentAttempt{}, "id = ?", attemptID).Error
}

func (r *gormBillingRepository) GetPaymentAttempt(ctx context.Context, attemptID uuid.UUID) (*PaymentAttempt, error) {
	var attempt PaymentAttempt
	err := r.db.WithContext(ctx).First(&attempt, "id = ?", attemptID).Error
	return &attempt, err
}

func (r *gormBillingRepository) GetPaymentAttemptsByInvoice(ctx context.Context, invoiceID uuid.UUID) ([]*PaymentAttempt, error) {
	var attempts []*PaymentAttempt
	err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Find(&attempts).Error
	return attempts, err
}

func (r *gormBillingRepository) GetFailedPaymentAttempts(ctx context.Context, retryBefore time.Time) ([]*PaymentAttempt, error) {
	var attempts []*PaymentAttempt
	err := r.db.WithContext(ctx).Where("status = 'failed' AND next_retry_at <= ?", retryBefore).Find(&attempts).Error
	return attempts, err
}

// Additional Dunning methods
func (r *gormBillingRepository) GetDunningProcess(ctx context.Context, processID uuid.UUID) (*DunningProcess, error) {
	var process DunningProcess
	err := r.db.WithContext(ctx).First(&process, "id = ?", processID).Error
	return &process, err
}

func (r *gormBillingRepository) GetDunningProcessByInvoice(ctx context.Context, invoiceID uuid.UUID) (*DunningProcess, error) {
	var process DunningProcess
	err := r.db.WithContext(ctx).First(&process, "invoice_id = ?", invoiceID).Error
	return &process, err
}

func (r *gormBillingRepository) GetActiveDunningProcesses(ctx context.Context) ([]*DunningProcess, error) {
	var processes []*DunningProcess
	err := r.db.WithContext(ctx).Where("status = 'active'").Find(&processes).Error
	return processes, err
}

func (r *gormBillingRepository) GetDunningProcessesDueForAction(ctx context.Context, before time.Time) ([]*DunningProcess, error) {
	var processes []*DunningProcess
	err := r.db.WithContext(ctx).Where("next_action_date <= ?", before).Find(&processes).Error
	return processes, err
}

func (r *gormBillingRepository) UpdateDunningProcess(ctx context.Context, process *DunningProcess) error {
	return r.db.WithContext(ctx).Save(process).Error
}

func (r *gormBillingRepository) DeleteDunningProcess(ctx context.Context, processID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&DunningProcess{}, "id = ?", processID).Error
}

func (r *gormBillingRepository) UpdateDunningAction(ctx context.Context, action *DunningAction) error {
	return r.db.WithContext(ctx).Save(action).Error
}

func (r *gormBillingRepository) DeleteDunningAction(ctx context.Context, actionID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&DunningAction{}, "id = ?", actionID).Error
}

func (r *gormBillingRepository) GetDunningActionsByProcess(ctx context.Context, processID uuid.UUID) ([]*DunningAction, error) {
	var actions []*DunningAction
	err := r.db.WithContext(ctx).Where("dunning_process_id = ?", processID).Find(&actions).Error
	return actions, err
}

func (r *gormBillingRepository) GetDunningActionsDueForExecution(ctx context.Context, before time.Time) ([]*DunningAction, error) {
	var actions []*DunningAction
	err := r.db.WithContext(ctx).Where("scheduled_at <= ? AND status = 'pending'", before).Find(&actions).Error
	return actions, err
}

// Metrics and summary methods
func (r *gormBillingRepository) GetRevenueSummary(ctx context.Context, filter AnalyticsFilter) (*RevenueSummary, error) {
	return &RevenueSummary{}, nil
}

func (r *gormBillingRepository) GetSubscriptionMetrics(ctx context.Context, filter AnalyticsFilter) (*SubscriptionMetrics, error) {
	return &SubscriptionMetrics{}, nil
}

func (r *gormBillingRepository) GetUsageMetrics(ctx context.Context, filter AnalyticsFilter) (*UsageMetrics, error) {
	return &UsageMetrics{}, nil
}

func (r *gormBillingRepository) GetChurnMetrics(ctx context.Context, filter AnalyticsFilter) (*ChurnMetrics, error) {
	return &ChurnMetrics{}, nil
}

func (r *gormBillingRepository) GetPaymentMetrics(ctx context.Context, filter AnalyticsFilter) (*PaymentMetrics, error) {
	return &PaymentMetrics{}, nil
}
