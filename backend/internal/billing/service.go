package billing

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
)

// BillingService handles all billing operations
type BillingService interface {
	// Billing Plans
	CreateBillingPlan(ctx context.Context, plan *BillingPlan) error
	GetBillingPlan(ctx context.Context, planID uuid.UUID) (*BillingPlan, error)
	GetBillingPlans(ctx context.Context, filter PlanFilter) ([]*BillingPlan, error)
	UpdateBillingPlan(ctx context.Context, plan *BillingPlan) error
	DeleteBillingPlan(ctx context.Context, planID uuid.UUID) error
	
	// Subscription management
	CreateSubscription(ctx context.Context, tenantID, planID uuid.UUID, paymentMethodID *string) (*TenantSubscription, error)
	GetSubscription(ctx context.Context, tenantID uuid.UUID) (*TenantSubscription, error)
	UpdateSubscription(ctx context.Context, tenantID uuid.UUID, updates SubscriptionUpdate) (*TenantSubscription, error)
	CancelSubscription(ctx context.Context, tenantID uuid.UUID, reason string, cancelImmediately bool) error
	
	// Usage Tiers
	CreateUsageTier(ctx context.Context, tier *UsageTier) error
	GetUsageTiersByPlan(ctx context.Context, planID uuid.UUID) ([]*UsageTier, error)
	UpdateUsageTier(ctx context.Context, tier *UsageTier) error
	DeleteUsageTier(ctx context.Context, tierID uuid.UUID) error
	
	// Plan changes
	UpgradePlan(ctx context.Context, tenantID, newPlanID uuid.UUID) (*TenantSubscription, error)
	DowngradePlan(ctx context.Context, tenantID, newPlanID uuid.UUID) (*TenantSubscription, error)
	ProcessPendingPlanChanges(ctx context.Context) error
	
	// Usage tracking
	RecordUsage(ctx context.Context, tenantID uuid.UUID, usageType UsageType, quantity int64, metadata map[string]interface{}) error
	GetUsageSummary(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (map[UsageType]int64, error)
	CheckUsageLimits(ctx context.Context, tenantID uuid.UUID) (*UsageLimitStatus, error)
	
	// Invoicing
	GenerateInvoice(ctx context.Context, subscriptionID uuid.UUID, periodStart, periodEnd time.Time) (*Invoice, error)
	GetInvoices(ctx context.Context, tenantID uuid.UUID, filter InvoiceFilter) ([]*Invoice, int64, error)
	GetInvoice(ctx context.Context, tenantID, invoiceID uuid.UUID) (*Invoice, error)
	ProcessRecurringBilling(ctx context.Context) error
	
	// Payment processing
	ProcessPayment(ctx context.Context, invoiceID uuid.UUID, paymentMethodID string) (*PaymentAttempt, error)
	RetryFailedPayments(ctx context.Context) error
	RefundPayment(ctx context.Context, invoiceID uuid.UUID, amount float64, reason string) error
	
	// Dunning management
	StartDunningProcess(ctx context.Context, invoiceID uuid.UUID) (*DunningProcess, error)
	ProcessDunning(ctx context.Context) error
	SuspendService(ctx context.Context, tenantID uuid.UUID, reason string) error
	ReactivateService(ctx context.Context, tenantID uuid.UUID) error
	
	// Analytics and reporting
	GetBillingAnalytics(ctx context.Context, filter AnalyticsFilter) (*BillingAnalytics, error)
	GetRevenueReport(ctx context.Context, filter RevenueReportFilter) (*RevenueReport, error)
	GetChurnAnalysis(ctx context.Context, period time.Duration) (*ChurnAnalysis, error)
}

// SubscriptionUpdate contains fields that can be updated for a subscription
type SubscriptionUpdate struct {
	PaymentMethodID *string
	BillingCycle    *BillingCycle
	Metadata        map[string]interface{}
}

// UsageLimitStatus represents the status of usage limits for a tenant
type UsageLimitStatus struct {
	TenantID      uuid.UUID                    `json:"tenant_id"`
	PlanLimits    map[string]interface{}       `json:"plan_limits"`
	CurrentUsage  map[UsageType]int64          `json:"current_usage"`
	LimitWarnings []UsageLimitWarning          `json:"limit_warnings"`
	Overages      map[UsageType]int64          `json:"overages"`
	EstimatedCost float64                      `json:"estimated_cost"`
}

// UsageLimitWarning represents a warning about approaching usage limits
type UsageLimitWarning struct {
	UsageType      UsageType `json:"usage_type"`
	CurrentUsage   int64     `json:"current_usage"`
	Limit          int64     `json:"limit"`
	PercentageUsed float64   `json:"percentage_used"`
	Warning        string    `json:"warning"`
}

// InvoiceFilter contains filtering options for invoice queries
type InvoiceFilter struct {
	Status      *InvoiceStatus `json:"status,omitempty"`
	StartDate   *time.Time     `json:"start_date,omitempty"`
	EndDate     *time.Time     `json:"end_date,omitempty"`
	MinAmount   *float64       `json:"min_amount,omitempty"`
	MaxAmount   *float64       `json:"max_amount,omitempty"`
	Limit       int            `json:"limit"`
	Offset      int            `json:"offset"`
	SortBy      string         `json:"sort_by"`
	SortDesc    bool           `json:"sort_desc"`
}

// BillingAnalytics contains billing analytics data
type BillingAnalytics struct {
	Period              time.Duration         `json:"period"`
	TotalRevenue        float64               `json:"total_revenue"`
	RecurringRevenue    float64               `json:"recurring_revenue"`
	NewRevenue          float64               `json:"new_revenue"`
	ChurnedRevenue      float64               `json:"churned_revenue"`
	ActiveSubscriptions int64                 `json:"active_subscriptions"`
	NewSubscriptions    int64                 `json:"new_subscriptions"`
	ChurnedSubscriptions int64                `json:"churned_subscriptions"`
	AverageRevenuePer   map[string]float64    `json:"average_revenue_per"`
	RevenueByPlan       map[string]float64    `json:"revenue_by_plan"`
	PaymentFailureRate  float64               `json:"payment_failure_rate"`
	DunningRecoveryRate float64               `json:"dunning_recovery_rate"`
}

// AnalyticsFilter contains filtering options for analytics
type AnalyticsFilter struct {
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	PlanIDs   []uuid.UUID `json:"plan_ids,omitempty"`
	Currency  *string     `json:"currency,omitempty"`
}

// RevenueReport contains detailed revenue reporting
type RevenueReport struct {
	Period           time.Duration            `json:"period"`
	TotalRevenue     float64                  `json:"total_revenue"`
	RevenueByMonth   []MonthlyRevenue         `json:"revenue_by_month"`
	RevenueByPlan    []PlanRevenue            `json:"revenue_by_plan"`
	RevenueByCountry []CountryRevenue         `json:"revenue_by_country"`
	UsageRevenue     float64                  `json:"usage_revenue"`
	SubscriptionRevenue float64               `json:"subscription_revenue"`
	Refunds          float64                  `json:"refunds"`
	NetRevenue       float64                  `json:"net_revenue"`
}

// RevenueReportFilter contains filtering options for revenue reports
type RevenueReportFilter struct {
	StartDate    time.Time    `json:"start_date"`
	EndDate      time.Time    `json:"end_date"`
	GroupBy      string       `json:"group_by"` // month, plan, country
	Currency     *string      `json:"currency,omitempty"`
	IncludeUsage bool         `json:"include_usage"`
}

// MonthlyRevenue represents revenue for a specific month
type MonthlyRevenue struct {
	Month   time.Time `json:"month"`
	Revenue float64   `json:"revenue"`
}

// PlanRevenue represents revenue by plan
type PlanRevenue struct {
	PlanID      uuid.UUID `json:"plan_id"`
	PlanName    string    `json:"plan_name"`
	Revenue     float64   `json:"revenue"`
	Subscribers int64     `json:"subscribers"`
}

// CountryRevenue represents revenue by country
type CountryRevenue struct {
	Country     string  `json:"country"`
	Revenue     float64 `json:"revenue"`
	Subscribers int64   `json:"subscribers"`
}

// ChurnAnalysis contains churn analysis data
type ChurnAnalysis struct {
	Period              time.Duration `json:"period"`
	ChurnRate           float64       `json:"churn_rate"`
	RevenueChurnRate    float64       `json:"revenue_churn_rate"`
	ChurnedCustomers    int64         `json:"churned_customers"`
	ChurnReasons        map[string]int64 `json:"churn_reasons"`
	ChurnByPlan         map[string]float64 `json:"churn_by_plan"`
	CustomerLifetimeValue float64     `json:"customer_lifetime_value"`
}

// service implements BillingService
type service struct {
	repo            BillingRepository
	paymentProvider PaymentProvider
	emailService    EmailService
	analyticsService AnalyticsService
}

// PaymentProvider interface for payment processing
type PaymentProvider interface {
	CreateCharge(amount float64, currency, paymentMethodID string, metadata map[string]interface{}) (*PaymentResult, error)
	RefundCharge(chargeID string, amount float64, reason string) (*RefundResult, error)
	GetPaymentMethod(paymentMethodID string) (*PaymentMethod, error)
}

// PaymentResult represents the result of a payment attempt
type PaymentResult struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	Amount      float64                `json:"amount"`
	Currency    string                 `json:"currency"`
	Metadata    map[string]interface{} `json:"metadata"`
	FailureReason *string              `json:"failure_reason,omitempty"`
}

// RefundResult represents the result of a refund
type RefundResult struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Last4    string `json:"last4"`
	Brand    string `json:"brand"`
	ExpMonth int    `json:"exp_month"`
	ExpYear  int    `json:"exp_year"`
}

// EmailService interface for sending billing emails
type EmailService interface {
	SendInvoice(ctx context.Context, tenantID uuid.UUID, invoice *Invoice) error
	SendPaymentFailedEmail(ctx context.Context, tenantID uuid.UUID, invoice *Invoice) error
	SendDunningEmail(ctx context.Context, tenantID uuid.UUID, step int, invoice *Invoice) error
	SendSubscriptionCanceledEmail(ctx context.Context, tenantID uuid.UUID, subscription *TenantSubscription) error
}

// AnalyticsService interface for tracking billing events
type AnalyticsService interface {
	TrackBillingEvent(ctx context.Context, tenantID uuid.UUID, eventType string, properties map[string]interface{}) error
}

// NewBillingService creates a new billing service
func NewBillingService(repo BillingRepository, paymentProvider PaymentProvider, emailService EmailService, analyticsService AnalyticsService) BillingService {
	return &service{
		repo:            repo,
		paymentProvider: paymentProvider,
		emailService:    emailService,
		analyticsService: analyticsService,
	}
}

// Billing Plans implementations
func (s *service) CreateBillingPlan(ctx context.Context, plan *BillingPlan) error {
	return s.repo.CreateBillingPlan(ctx, plan)
}

func (s *service) GetBillingPlan(ctx context.Context, planID uuid.UUID) (*BillingPlan, error) {
	return s.repo.GetBillingPlan(ctx, planID)
}

func (s *service) GetBillingPlans(ctx context.Context, filter PlanFilter) ([]*BillingPlan, error) {
	return s.repo.GetBillingPlans(ctx, filter)
}

func (s *service) UpdateBillingPlan(ctx context.Context, plan *BillingPlan) error {
	return s.repo.UpdateBillingPlan(ctx, plan)
}

func (s *service) DeleteBillingPlan(ctx context.Context, planID uuid.UUID) error {
	return s.repo.DeleteBillingPlan(ctx, planID)
}

// Usage Tiers implementations (duplicates removed)

// Usage Tiers implementations
func (s *service) CreateUsageTier(ctx context.Context, tier *UsageTier) error {
	return s.repo.CreateUsageTier(ctx, tier)
}

func (s *service) GetUsageTiersByPlan(ctx context.Context, planID uuid.UUID) ([]*UsageTier, error) {
	return s.repo.GetUsageTiersByPlan(ctx, planID)
}

func (s *service) UpdateUsageTier(ctx context.Context, tier *UsageTier) error {
	return s.repo.UpdateUsageTier(ctx, tier)
}

func (s *service) DeleteUsageTier(ctx context.Context, tierID uuid.UUID) error {
	return s.repo.DeleteUsageTier(ctx, tierID)
}

// Subscription management implementations
func (s *service) CreateSubscription(ctx context.Context, tenantID, planID uuid.UUID, paymentMethodID *string) (*TenantSubscription, error) {
	// TODO: Implement subscription creation
	// 1. Get billing plan details
	// 2. Calculate trial period if applicable
	// 3. Set initial billing period
	// 4. Create subscription record
	// 5. Track analytics event
	
	plan, err := s.repo.GetBillingPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing plan: %w", err)
	}
	
	subscription := &TenantSubscription{
		ID:       uuid.New(),
		TenantID: tenantID,
		PlanID:   planID,
		Status:   SubscriptionStatusActive,
		BillingCycle: plan.BillingCycle,
		BaseAmount: plan.BasePrice,
		Currency: plan.Currency,
		PaymentMethodID: paymentMethodID,
		CurrentPeriodStart: time.Now(),
		CurrentPeriodEnd: time.Now().AddDate(0, 0, plan.BillingCycle.GetDays()),
		NextBillingDate: time.Now().AddDate(0, 0, plan.BillingCycle.GetDays()),
	}
	
	// Set trial end if applicable
	if plan.TrialPeriodDays > 0 {
		trialEnd := time.Now().AddDate(0, 0, plan.TrialPeriodDays)
		subscription.TrialEnd = &trialEnd
		subscription.NextBillingDate = trialEnd
	}
	
	err = s.repo.CreateSubscription(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}
	
	// Track analytics
	s.analyticsService.TrackBillingEvent(ctx, tenantID, "subscription_created", map[string]interface{}{
		"plan_id": planID,
		"trial_days": plan.TrialPeriodDays,
		"base_amount": plan.BasePrice,
	})
	
	return subscription, nil
}

func (s *service) GetSubscription(ctx context.Context, tenantID uuid.UUID) (*TenantSubscription, error) {
	// TODO: Implement get subscription with caching
	return s.repo.GetSubscriptionByTenantID(ctx, tenantID)
}

func (s *service) UpdateSubscription(ctx context.Context, tenantID uuid.UUID, updates SubscriptionUpdate) (*TenantSubscription, error) {
	// TODO: Implement subscription updates
	// 1. Get current subscription
	// 2. Validate updates
	// 3. Apply updates
	// 4. Track analytics
	
	subscription, err := s.repo.GetSubscriptionByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}
	
	// Apply updates
	if updates.PaymentMethodID != nil {
		subscription.PaymentMethodID = updates.PaymentMethodID
	}
	if updates.BillingCycle != nil {
		subscription.BillingCycle = *updates.BillingCycle
	}
	
	err = s.repo.UpdateSubscription(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}
	
	return subscription, nil
}

func (s *service) CancelSubscription(ctx context.Context, tenantID uuid.UUID, reason string, cancelImmediately bool) error {
	// TODO: Implement subscription cancellation
	// 1. Get subscription
	// 2. Set cancellation date (immediate or end of period)
	// 3. Update status
	// 4. Send cancellation email
	// 5. Track analytics
	
	subscription, err := s.repo.GetSubscriptionByTenantID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}
	
	now := time.Now()
	subscription.CanceledAt = &now
	subscription.CancellationReason = reason
	
	if cancelImmediately {
		subscription.Status = SubscriptionStatusCanceled
		subscription.CurrentPeriodEnd = now
	} else {
		// Cancel at end of current period
		subscription.Status = SubscriptionStatusActive // Keep active until period end
		// The status will be updated to canceled when the period ends
	}
	
	err = s.repo.UpdateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}
	
	// Send cancellation email
	s.emailService.SendSubscriptionCanceledEmail(ctx, tenantID, subscription)
	
	// Track analytics
	s.analyticsService.TrackBillingEvent(ctx, tenantID, "subscription_canceled", map[string]interface{}{
		"reason": reason,
		"immediate": cancelImmediately,
	})
	
	return nil
}

// Plan change implementations
func (s *service) UpgradePlan(ctx context.Context, tenantID, newPlanID uuid.UUID) (*TenantSubscription, error) {
	// TODO: Implement plan upgrade with proration
	return s.changePlan(ctx, tenantID, newPlanID, "upgrade")
}

func (s *service) DowngradePlan(ctx context.Context, tenantID, newPlanID uuid.UUID) (*TenantSubscription, error) {
	// TODO: Implement plan downgrade (effective at period end)
	return s.changePlan(ctx, tenantID, newPlanID, "downgrade")
}

func (s *service) changePlan(ctx context.Context, tenantID, newPlanID uuid.UUID, changeType string) (*TenantSubscription, error) {
	subscription, err := s.repo.GetSubscriptionByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}
	
	if !subscription.CanUpgrade() && changeType == "upgrade" {
		return nil, fmt.Errorf("subscription cannot be upgraded at this time")
	}
	
	newPlan, err := s.repo.GetBillingPlan(ctx, newPlanID)
	if err != nil {
		return nil, fmt.Errorf("new plan not found: %w", err)
	}
	
	// Calculate proration for upgrades (immediate) or schedule for downgrades (end of period)
	var effectiveDate time.Time
	var prorationAmount float64
	
	if changeType == "upgrade" {
		// Immediate upgrade with proration
		effectiveDate = time.Now()
		daysUsed := int(time.Since(subscription.CurrentPeriodStart).Hours() / 24)
		totalDays := subscription.BillingCycle.GetDays()
		prorationAmount = CalculateProration(subscription.BaseAmount, newPlan.BasePrice, daysUsed, totalDays)
	} else {
		// Downgrade at end of current period
		effectiveDate = subscription.CurrentPeriodEnd
		prorationAmount = 0 // No immediate charge for downgrades
	}
	
	// Set pending plan change
	subscription.PendingPlanChange = &PlanChange{
		NewPlanID:       newPlanID,
		EffectiveDate:   effectiveDate,
		ProrationAmount: prorationAmount,
		ChangeReason:    changeType,
	}
	
	err = s.repo.UpdateSubscription(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to schedule plan change: %w", err)
	}
	
	// For upgrades, process immediately
	if changeType == "upgrade" {
		err = s.processPlanChange(ctx, subscription)
		if err != nil {
			return nil, fmt.Errorf("failed to process plan upgrade: %w", err)
		}
	}
	
	// Track analytics
	s.analyticsService.TrackBillingEvent(ctx, tenantID, fmt.Sprintf("plan_%s_scheduled", changeType), map[string]interface{}{
		"old_plan_id": subscription.PlanID,
		"new_plan_id": newPlanID,
		"proration_amount": prorationAmount,
		"effective_date": effectiveDate,
	})
	
	return subscription, nil
}

func (s *service) ProcessPendingPlanChanges(ctx context.Context) error {
	// TODO: Process all pending plan changes that are due
	subscriptions, err := s.repo.GetSubscriptionsWithPendingChanges(ctx)
	if err != nil {
		return fmt.Errorf("failed to get subscriptions with pending changes: %w", err)
	}
	
	for _, subscription := range subscriptions {
		if subscription.PendingPlanChange != nil && time.Now().After(subscription.PendingPlanChange.EffectiveDate) {
			err := s.processPlanChange(ctx, subscription)
			if err != nil {
				// Log error but continue with other subscriptions
				// TODO: Implement proper error logging
				continue
			}
		}
	}
	
	return nil
}

func (s *service) processPlanChange(ctx context.Context, subscription *TenantSubscription) error {
	if subscription.PendingPlanChange == nil {
		return nil
	}
	
	change := subscription.PendingPlanChange
	
	// Update subscription to new plan
	subscription.PlanID = change.NewPlanID
	
	// Get new plan details
	newPlan, err := s.repo.GetBillingPlan(ctx, change.NewPlanID)
	if err != nil {
		return fmt.Errorf("failed to get new plan: %w", err)
	}
	
	subscription.BaseAmount = newPlan.BasePrice
	subscription.BillingCycle = newPlan.BillingCycle
	
	// Clear pending change
	subscription.PendingPlanChange = nil
	
	// If there's a proration amount, create an immediate invoice
	if change.ProrationAmount != 0 {
		// TODO: Create proration invoice
		// This would create an invoice for the prorated amount
	}
	
	err = s.repo.UpdateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to complete plan change: %w", err)
	}
	
	return nil
}

// Usage tracking implementations
func (s *service) RecordUsage(ctx context.Context, tenantID uuid.UUID, usageType UsageType, quantity int64, metadata map[string]interface{}) error {
	// TODO: Implement usage recording
	// 1. Determine current billing period
	// 2. Create usage record
	// 3. Check if usage limits are exceeded
	// 4. Track analytics
	
	subscription, err := s.repo.GetSubscriptionByTenantID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}
	
	usage := &UsageRecord{
		ID:                 uuid.New(),
		TenantID:          tenantID,
		UsageType:         usageType,
		Quantity:          quantity,
		Units:             string(usageType),
		BillingPeriodStart: subscription.CurrentPeriodStart,
		BillingPeriodEnd:  subscription.CurrentPeriodEnd,
		Metadata:          metadata,
		RecordedAt:        time.Now(),
	}
	
	err = s.repo.CreateUsageRecord(ctx, usage)
	if err != nil {
		return fmt.Errorf("failed to record usage: %w", err)
	}
	
	// Check usage limits asynchronously
	go func() {
		s.CheckUsageLimits(context.Background(), tenantID)
	}()
	
	return nil
}

func (s *service) GetUsageSummary(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (map[UsageType]int64, error) {
	// TODO: Implement usage summary aggregation
	return s.repo.GetUsageSummary(ctx, tenantID, startDate, endDate)
}

func (s *service) CheckUsageLimits(ctx context.Context, tenantID uuid.UUID) (*UsageLimitStatus, error) {
	// TODO: Implement usage limit checking
	// 1. Get subscription and plan limits
	// 2. Get current usage for billing period
	// 3. Compare against limits
	// 4. Generate warnings and overages
	// 5. Calculate estimated overage costs
	
	subscription, err := s.repo.GetSubscriptionByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}
	
	plan, err := s.repo.GetBillingPlan(ctx, subscription.PlanID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}
	
	currentUsage, err := s.repo.GetUsageSummary(ctx, tenantID, subscription.CurrentPeriodStart, subscription.CurrentPeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage summary: %w", err)
	}
	
	status := &UsageLimitStatus{
		TenantID:     tenantID,
		PlanLimits:   plan.Limits,
		CurrentUsage: currentUsage,
		Overages:     make(map[UsageType]int64),
	}
	
	// TODO: Implement limit checking logic
	// This would compare currentUsage against plan.Limits
	// and generate warnings and calculate overages
	
	return status, nil
}

// Invoice generation and management
func (s *service) GenerateInvoice(ctx context.Context, subscriptionID uuid.UUID, periodStart, periodEnd time.Time) (*Invoice, error) {
	subscription, err := s.repo.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	// Generate unique invoice number
	invoiceNumber, err := s.repo.GetNextInvoiceNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	invoice := &Invoice{
		ID:             uuid.New(),
		TenantID:       subscription.TenantID,
		SubscriptionID: subscriptionID,
		InvoiceNumber:  invoiceNumber,
		Status:         InvoiceStatusDraft,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		Currency:       subscription.Currency,
		DueDate:        time.Now().AddDate(0, 0, 30), // 30 days payment terms
	}

	// Add subscription base charge
	baseLineItem := &InvoiceLineItem{
		ID:          uuid.New(),
		InvoiceID:   invoice.ID,
		Description: fmt.Sprintf("Subscription to %s (%s)", subscription.Plan.Name, subscription.BillingCycle),
		Quantity:    1,
		UnitPrice:   subscription.BaseAmount,
		TotalPrice:  subscription.BaseAmount,
		ItemType:    "subscription",
		PeriodStart: &periodStart,
		PeriodEnd:   &periodEnd,
	}

	invoice.LineItems = append(invoice.LineItems, *baseLineItem)
	invoice.SubtotalAmount = subscription.BaseAmount

	// Add usage charges
	usageCharges, err := s.calculateUsageCharges(ctx, subscription, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate usage charges: %w", err)
	}

	for _, usageCharge := range usageCharges {
		invoice.LineItems = append(invoice.LineItems, usageCharge)
		invoice.SubtotalAmount += usageCharge.TotalPrice
	}

	// Calculate tax (if applicable)
	// TODO: Implement tax calculation based on tenant location
	invoice.TaxAmount = 0

	invoice.TotalAmount = invoice.SubtotalAmount + invoice.TaxAmount

	// Start transaction
	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create invoice
	err = s.repo.CreateInvoice(tx.GetContext(), invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Create line items
	for i := range invoice.LineItems {
		err = s.repo.CreateInvoiceLineItem(tx.GetContext(), &invoice.LineItems[i])
		if err != nil {
			return nil, fmt.Errorf("failed to create line item: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Update invoice status to pending
	invoice.Status = InvoiceStatusPending
	err = s.repo.UpdateInvoice(ctx, invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice status: %w", err)
	}

	// Send invoice email
	go s.emailService.SendInvoice(ctx, subscription.TenantID, invoice)

	// Track analytics
	s.analyticsService.TrackBillingEvent(ctx, subscription.TenantID, "invoice_generated", map[string]interface{}{
		"invoice_id":     invoice.ID,
		"amount":         invoice.TotalAmount,
		"billing_period": fmt.Sprintf("%s to %s", periodStart.Format("2006-01-02"), periodEnd.Format("2006-01-02")),
	})

	return invoice, nil
}

func (s *service) calculateUsageCharges(ctx context.Context, subscription *TenantSubscription, periodStart, periodEnd time.Time) ([]InvoiceLineItem, error) {
	var lineItems []InvoiceLineItem

	// Get usage summary for the billing period
	usageSummary, err := s.repo.GetUsageSummary(ctx, subscription.TenantID, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage summary: %w", err)
	}

	// Get usage tiers for the plan
	usageTiers, err := s.repo.GetUsageTiersByPlan(ctx, subscription.PlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage tiers: %w", err)
	}

	// Group tiers by usage type
	tiersByType := make(map[UsageType][]*UsageTier)
	for _, tier := range usageTiers {
		tiersByType[tier.UsageType] = append(tiersByType[tier.UsageType], tier)
	}

	// Calculate charges for each usage type
	for usageType, totalUsage := range usageSummary {
		if totalUsage == 0 {
			continue
		}

		tiers, exists := tiersByType[usageType]
		if !exists {
			continue // No pricing tiers for this usage type
		}

		// Calculate tiered pricing
		remainingUsage := totalUsage
		totalCharge := 0.0

		for _, tier := range tiers {
			if remainingUsage <= 0 {
				break
			}

			// Determine units to charge in this tier
			var unitsInTier int64
			if tier.MaxUnits == nil {
				// Unlimited tier - charge all remaining usage
				unitsInTier = remainingUsage
			} else {
				// Limited tier
				tierCapacity := *tier.MaxUnits - tier.MinUnits
				if remainingUsage <= tierCapacity {
					unitsInTier = remainingUsage
				} else {
					unitsInTier = tierCapacity
				}
			}

			if unitsInTier > 0 {
				tierCharge := float64(unitsInTier) * tier.PricePerUnit
				totalCharge += tierCharge
				remainingUsage -= unitsInTier
			}
		}

		if totalCharge > 0 {
			lineItem := InvoiceLineItem{
				ID:          uuid.New(),
				Description: fmt.Sprintf("Usage: %s (%d %s)", usageType, totalUsage, string(usageType)),
				Quantity:    totalUsage,
				UnitPrice:   totalCharge / float64(totalUsage), // Average unit price
				TotalPrice:  totalCharge,
				ItemType:    "usage",
				UsageType:   &usageType,
				PeriodStart: &periodStart,
				PeriodEnd:   &periodEnd,
			}
			lineItems = append(lineItems, lineItem)
		}
	}

	return lineItems, nil
}

func (s *service) GetInvoices(ctx context.Context, tenantID uuid.UUID, filter InvoiceFilter) ([]*Invoice, int64, error) {
	return s.repo.GetInvoicesByTenant(ctx, tenantID, filter)
}

func (s *service) GetInvoice(ctx context.Context, tenantID, invoiceID uuid.UUID) (*Invoice, error) {
	invoice, err := s.repo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	// Verify tenant ownership
	if invoice.TenantID != tenantID {
		return nil, fmt.Errorf("invoice not found for tenant")
	}

	return invoice, nil
}

func (s *service) ProcessRecurringBilling(ctx context.Context) error {
	// Get subscriptions due for billing
	subscriptions, err := s.repo.GetSubscriptionsDueForBilling(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to get subscriptions due for billing: %w", err)
	}

	for _, subscription := range subscriptions {
		err := s.processSubscriptionBilling(ctx, subscription)
		if err != nil {
			// Log error but continue with other subscriptions
			// TODO: Implement proper error logging
			continue
		}
	}

	return nil
}

func (s *service) processSubscriptionBilling(ctx context.Context, subscription *TenantSubscription) error {
	// Generate invoice for current period
	invoice, err := s.GenerateInvoice(ctx, subscription.ID, 
		subscription.CurrentPeriodStart, subscription.CurrentPeriodEnd)
	if err != nil {
		return fmt.Errorf("failed to generate invoice: %w", err)
	}

	// Attempt payment if payment method is available
	if subscription.PaymentMethodID != nil {
		_, err = s.ProcessPayment(ctx, invoice.ID, *subscription.PaymentMethodID)
		if err != nil {
			// Payment failed, but invoice is created
			// Dunning process will handle failed payments
		}
	}

	// Update subscription for next billing period
	subscription.CurrentPeriodStart = subscription.CurrentPeriodEnd
	subscription.CurrentPeriodEnd = subscription.CurrentPeriodEnd.AddDate(0, 0, subscription.BillingCycle.GetDays())
	subscription.NextBillingDate = subscription.CurrentPeriodEnd

	err = s.repo.UpdateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// Payment processing
func (s *service) ProcessPayment(ctx context.Context, invoiceID uuid.UUID, paymentMethodID string) (*PaymentAttempt, error) {
	invoice, err := s.repo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	if invoice.Status != InvoiceStatusPending {
		return nil, fmt.Errorf("invoice is not in pending status")
	}

	// Create payment attempt record
	attempt := &PaymentAttempt{
		ID:            uuid.New(),
		InvoiceID:     invoiceID,
		TenantID:      invoice.TenantID,
		Status:        PaymentStatusPending,
		Amount:        invoice.RemainingAmount(),
		Currency:      invoice.Currency,
		PaymentMethod: paymentMethodID,
		AttemptedAt:   time.Now(),
		RetryCount:    0,
	}

	err = s.repo.CreatePaymentAttempt(ctx, attempt)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment attempt: %w", err)
	}

	// Process payment with payment provider
	paymentResult, err := s.paymentProvider.CreateCharge(
		attempt.Amount,
		attempt.Currency,
		paymentMethodID,
		map[string]interface{}{
			"invoice_id":     invoiceID,
			"tenant_id":      invoice.TenantID,
			"attempt_id":     attempt.ID,
		},
	)

	// Update payment attempt based on result
	if err != nil || paymentResult.Status != "succeeded" {
		attempt.Status = PaymentStatusFailed
		if paymentResult != nil && paymentResult.FailureReason != nil {
			attempt.FailureReason = paymentResult.FailureReason
		} else {
			reason := err.Error()
			attempt.FailureReason = &reason
		}
		
		// Schedule retry
		if attempt.RetryCount < 3 {
			nextRetry := time.Now().Add(time.Duration(attempt.RetryCount+1) * 24 * time.Hour)
			attempt.NextRetryAt = &nextRetry
		}
	} else {
		attempt.Status = PaymentStatusSuccess
		attempt.ProviderChargeID = &paymentResult.ID
		completedAt := time.Now()
		attempt.CompletedAt = &completedAt

		// Update invoice
		invoice.Status = InvoiceStatusPaid
		invoice.PaidAmount = invoice.TotalAmount
		invoice.PaidAt = &completedAt
		invoice.PaymentMethod = &paymentMethodID

		err = s.repo.UpdateInvoice(ctx, invoice)
		if err != nil {
			return nil, fmt.Errorf("failed to update invoice: %w", err)
		}
	}

	if paymentResult != nil {
		attempt.ProviderResponse = paymentResult.Metadata
	}

	err = s.repo.UpdatePaymentAttempt(ctx, attempt)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment attempt: %w", err)
	}

	// Track analytics
	eventType := "payment_success"
	if attempt.Status == PaymentStatusFailed {
		eventType = "payment_failed"
		
		// Send payment failed email
		go s.emailService.SendPaymentFailedEmail(ctx, invoice.TenantID, invoice)
		
		// Start dunning process if this is the first failure
		if attempt.RetryCount == 0 {
			go s.StartDunningProcess(ctx, invoiceID)
		}
	}

	s.analyticsService.TrackBillingEvent(ctx, invoice.TenantID, eventType, map[string]interface{}{
		"invoice_id":     invoiceID,
		"amount":         attempt.Amount,
		"payment_method": paymentMethodID,
		"retry_count":    attempt.RetryCount,
	})

	return attempt, nil
}

func (s *service) RetryFailedPayments(ctx context.Context) error {
	// Get failed payment attempts that are due for retry
	attempts, err := s.repo.GetFailedPaymentAttempts(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to get failed payment attempts: %w", err)
	}

	for _, attempt := range attempts {
		// Retry the payment
		_, err := s.ProcessPayment(ctx, attempt.InvoiceID, attempt.PaymentMethod)
		if err != nil {
			// Log error but continue with other attempts
			continue
		}
	}

	return nil
}

func (s *service) RefundPayment(ctx context.Context, invoiceID uuid.UUID, amount float64, reason string) error {
	invoice, err := s.repo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return fmt.Errorf("invoice not found: %w", err)
	}

	if invoice.Status != InvoiceStatusPaid {
		return fmt.Errorf("invoice is not paid")
	}

	if amount > invoice.PaidAmount {
		return fmt.Errorf("refund amount exceeds paid amount")
	}

	// Get the successful payment attempt
	attempts, err := s.repo.GetPaymentAttemptsByInvoice(ctx, invoiceID)
	if err != nil {
		return fmt.Errorf("failed to get payment attempts: %w", err)
	}

	var successfulAttempt *PaymentAttempt
	for _, attempt := range attempts {
		if attempt.Status == PaymentStatusSuccess {
			successfulAttempt = attempt
			break
		}
	}

	if successfulAttempt == nil || successfulAttempt.ProviderChargeID == nil {
		return fmt.Errorf("no successful payment found to refund")
	}

	// Process refund with payment provider
	refundResult, err := s.paymentProvider.RefundCharge(*successfulAttempt.ProviderChargeID, amount, reason)
	if err != nil {
		return fmt.Errorf("failed to process refund: %w", err)
	}

	// Update invoice
	invoice.Status = InvoiceStatusRefunded
	invoice.PaidAmount -= amount

	err = s.repo.UpdateInvoice(ctx, invoice)
	if err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	// Track analytics
	s.analyticsService.TrackBillingEvent(ctx, invoice.TenantID, "payment_refunded", map[string]interface{}{
		"invoice_id":   invoiceID,
		"amount":       amount,
		"reason":       reason,
		"refund_id":    refundResult.ID,
	})

	return nil
}

// Dunning management
func (s *service) StartDunningProcess(ctx context.Context, invoiceID uuid.UUID) (*DunningProcess, error) {
	invoice, err := s.repo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	// Check if dunning process already exists
	existing, err := s.repo.GetDunningProcessByInvoice(ctx, invoiceID)
	if err == nil && existing != nil {
		return existing, nil // Already exists
	}

	process := &DunningProcess{
		ID:             uuid.New(),
		TenantID:       invoice.TenantID,
		InvoiceID:      invoiceID,
		SubscriptionID: invoice.SubscriptionID,
		CurrentStep:    1,
		TotalSteps:     5, // Configurable dunning steps
		StartedAt:      time.Now(),
		NextActionAt:   &time.Time{}, // Will be set when scheduling first action
	}

	err = s.repo.CreateDunningProcess(ctx, process)
	if err != nil {
		return nil, fmt.Errorf("failed to create dunning process: %w", err)
	}

	// Schedule first dunning action (immediate email)
	firstAction := &DunningAction{
		ID:               uuid.New(),
		DunningProcessID: process.ID,
		ActionType:       "email",
		StepNumber:       1,
		Description:      "Send payment reminder email",
		Status:           "pending",
		ScheduledAt:      time.Now(),
	}

	err = s.repo.CreateDunningAction(ctx, firstAction)
	if err != nil {
		return nil, fmt.Errorf("failed to create dunning action: %w", err)
	}

	nextActionAt := time.Now()
	process.NextActionAt = &nextActionAt

	err = s.repo.UpdateDunningProcess(ctx, process)
	if err != nil {
		return nil, fmt.Errorf("failed to update dunning process: %w", err)
	}

	return process, nil
}

func (s *service) ProcessDunning(ctx context.Context) error {
	// Get dunning processes due for action
	processes, err := s.repo.GetDunningProcessesDueForAction(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to get dunning processes: %w", err)
	}

	for _, process := range processes {
		err := s.processDunningStep(ctx, process)
		if err != nil {
			// Log error but continue with other processes
			continue
		}
	}

	// Also process individual actions that are due
	actions, err := s.repo.GetDunningActionsDueForExecution(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to get dunning actions: %w", err)
	}

	for _, action := range actions {
		err := s.executeDunningAction(ctx, action)
		if err != nil {
			// Log error but continue with other actions
			continue
		}
	}

	return nil
}

func (s *service) processDunningStep(ctx context.Context, process *DunningProcess) error {
	// Determine next action based on current step
	var nextAction *DunningAction

	switch process.CurrentStep {
	case 1:
		// Step 1: Send reminder email
		nextAction = &DunningAction{
			ID:               uuid.New(),
			DunningProcessID: process.ID,
			ActionType:       "email",
			StepNumber:       1,
			Description:      "Send payment reminder email",
			Status:           "pending",
			ScheduledAt:      time.Now(),
		}
	case 2:
		// Step 2: Send second reminder after 3 days
		nextAction = &DunningAction{
			ID:               uuid.New(),
			DunningProcessID: process.ID,
			ActionType:       "email",
			StepNumber:       2,
			Description:      "Send second payment reminder",
			Status:           "pending",
			ScheduledAt:      time.Now().AddDate(0, 0, 3),
		}
	case 3:
		// Step 3: Send final notice after 7 days
		nextAction = &DunningAction{
			ID:               uuid.New(),
			DunningProcessID: process.ID,
			ActionType:       "email",
			StepNumber:       3,
			Description:      "Send final payment notice",
			Status:           "pending",
			ScheduledAt:      time.Now().AddDate(0, 0, 7),
		}
	case 4:
		// Step 4: Suspend service after 14 days
		nextAction = &DunningAction{
			ID:               uuid.New(),
			DunningProcessID: process.ID,
			ActionType:       "suspend",
			StepNumber:       4,
			Description:      "Suspend service due to non-payment",
			Status:           "pending",
			ScheduledAt:      time.Now().AddDate(0, 0, 14),
		}
	case 5:
		// Step 5: Cancel subscription after 30 days
		nextAction = &DunningAction{
			ID:               uuid.New(),
			DunningProcessID: process.ID,
			ActionType:       "cancel",
			StepNumber:       5,
			Description:      "Cancel subscription due to non-payment",
			Status:           "pending",
			ScheduledAt:      time.Now().AddDate(0, 0, 30),
		}
	default:
		// Process complete
		process.IsCompleted = true
		completedAt := time.Now()
		process.CompletedAt = &completedAt
		return s.repo.UpdateDunningProcess(ctx, process)
	}

	if nextAction != nil {
		err := s.repo.CreateDunningAction(ctx, nextAction)
		if err != nil {
			return fmt.Errorf("failed to create dunning action: %w", err)
		}

		process.CurrentStep++
		process.NextActionAt = &nextAction.ScheduledAt

		err = s.repo.UpdateDunningProcess(ctx, process)
		if err != nil {
			return fmt.Errorf("failed to update dunning process: %w", err)
		}
	}

	return nil
}

func (s *service) executeDunningAction(ctx context.Context, action *DunningAction) error {
	switch action.ActionType {
	case "email":
		err := s.sendDunningEmail(ctx, action)
		if err != nil {
			action.Status = "failed"
			errMsg := err.Error()
			action.ErrorMessage = &errMsg
		} else {
			action.Status = "completed"
		}
	case "suspend":
		err := s.suspendServiceForDunning(ctx, action)
		if err != nil {
			action.Status = "failed"
			errMsg := err.Error()
			action.ErrorMessage = &errMsg
		} else {
			action.Status = "completed"
		}
	case "cancel":
		err := s.cancelSubscriptionForDunning(ctx, action)
		if err != nil {
			action.Status = "failed"
			errMsg := err.Error()
			action.ErrorMessage = &errMsg
		} else {
			action.Status = "completed"
		}
	}

	executedAt := time.Now()
	action.ExecutedAt = &executedAt

	return s.repo.UpdateDunningAction(ctx, action)
}

func (s *service) sendDunningEmail(ctx context.Context, action *DunningAction) error {
	process, err := s.repo.GetDunningProcess(ctx, action.DunningProcessID)
	if err != nil {
		return fmt.Errorf("failed to get dunning process: %w", err)
	}

	err = s.emailService.SendDunningEmail(ctx, process.TenantID, action.StepNumber, &process.Invoice)
	if err != nil {
		return fmt.Errorf("failed to send dunning email: %w", err)
	}

	// Update process email count
	process.EmailsSent++
	return s.repo.UpdateDunningProcess(ctx, process)
}

func (s *service) suspendServiceForDunning(ctx context.Context, action *DunningAction) error {
	process, err := s.repo.GetDunningProcess(ctx, action.DunningProcessID)
	if err != nil {
		return fmt.Errorf("failed to get dunning process: %w", err)
	}

	err = s.SuspendService(ctx, process.TenantID, "Non-payment of invoice")
	if err != nil {
		return fmt.Errorf("failed to suspend service: %w", err)
	}

	// Update process suspension status
	process.ServiceSuspended = true
	suspendedAt := time.Now()
	process.ServiceSuspendedAt = &suspendedAt

	return s.repo.UpdateDunningProcess(ctx, process)
}

func (s *service) cancelSubscriptionForDunning(ctx context.Context, action *DunningAction) error {
	process, err := s.repo.GetDunningProcess(ctx, action.DunningProcessID)
	if err != nil {
		return fmt.Errorf("failed to get dunning process: %w", err)
	}

	err = s.CancelSubscription(ctx, process.TenantID, "Non-payment of invoice", true)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	// Mark process as completed
	process.IsCompleted = true
	completedAt := time.Now()
	process.CompletedAt = &completedAt

	return s.repo.UpdateDunningProcess(ctx, process)
}

func (s *service) SuspendService(ctx context.Context, tenantID uuid.UUID, reason string) error {
	subscription, err := s.repo.GetSubscriptionByTenantID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	subscription.Status = SubscriptionStatusSuspended
	err = s.repo.UpdateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to suspend subscription: %w", err)
	}

	// Track analytics
	s.analyticsService.TrackBillingEvent(ctx, tenantID, "service_suspended", map[string]interface{}{
		"reason": reason,
	})

	return nil
}

func (s *service) ReactivateService(ctx context.Context, tenantID uuid.UUID) error {
	subscription, err := s.repo.GetSubscriptionByTenantID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	subscription.Status = SubscriptionStatusActive
	err = s.repo.UpdateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to reactivate subscription: %w", err)
	}

	// Track analytics
	s.analyticsService.TrackBillingEvent(ctx, tenantID, "service_reactivated", map[string]interface{}{})

	return nil
}

// Analytics and reporting
func (s *service) GetBillingAnalytics(ctx context.Context, filter AnalyticsFilter) (*BillingAnalytics, error) {
	// Get revenue metrics
	revenueSummary, err := s.repo.GetRevenueSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue summary: %w", err)
	}

	// Get subscription metrics
	subscriptionMetrics, err := s.repo.GetSubscriptionMetrics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription metrics: %w", err)
	}

	// Get payment metrics
	paymentMetrics, err := s.repo.GetPaymentMetrics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment metrics: %w", err)
	}

	// Get churn metrics
	churnMetrics, err := s.repo.GetChurnMetrics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get churn metrics: %w", err)
	}

	return &BillingAnalytics{
		Period:               filter.EndDate.Sub(filter.StartDate),
		TotalRevenue:         revenueSummary.TotalRevenue,
		RecurringRevenue:     revenueSummary.RecurringRevenue,
		NewRevenue:           0, // TODO: Calculate new customer revenue
		ChurnedRevenue:       churnMetrics.ChurnedRevenue,
		ActiveSubscriptions:  subscriptionMetrics.ActiveSubscriptions,
		NewSubscriptions:     subscriptionMetrics.NewSubscriptions,
		ChurnedSubscriptions: churnMetrics.ChurnedSubscriptions,
		RevenueByPlan:        revenueSummary.RevenueByPlan,
		PaymentFailureRate:   1.0 - paymentMetrics.PaymentSuccessRate,
		DunningRecoveryRate:  paymentMetrics.DunningRecoveryRate,
	}, nil
}

func (s *service) GetRevenueReport(ctx context.Context, filter RevenueReportFilter) (*RevenueReport, error) {
	analyticsFilter := AnalyticsFilter{
		StartDate: filter.StartDate,
		EndDate:   filter.EndDate,
		Currency:  filter.Currency,
	}

	revenueSummary, err := s.repo.GetRevenueSummary(ctx, analyticsFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue summary: %w", err)
	}

	return &RevenueReport{
		Period:              filter.EndDate.Sub(filter.StartDate),
		TotalRevenue:        revenueSummary.TotalRevenue,
		UsageRevenue:        revenueSummary.UsageRevenue,
		SubscriptionRevenue: revenueSummary.RecurringRevenue,
		Refunds:            revenueSummary.RefundedRevenue,
		NetRevenue:         revenueSummary.NetRevenue,
		// TODO: Implement monthly/plan/country breakdowns
	}, nil
}

func (s *service) GetChurnAnalysis(ctx context.Context, period time.Duration) (*ChurnAnalysis, error) {
	endDate := time.Now()
	startDate := endDate.Add(-period)
	
	filter := AnalyticsFilter{
		StartDate: startDate,
		EndDate:   endDate,
	}

	churnMetrics, err := s.repo.GetChurnMetrics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get churn metrics: %w", err)
	}

	return &ChurnAnalysis{
		Period:                period,
		ChurnRate:            churnMetrics.ChurnRate,
		RevenueChurnRate:     churnMetrics.RevenueChurnRate,
		ChurnedCustomers:     churnMetrics.ChurnedSubscriptions,
		ChurnReasons:         churnMetrics.ChurnReasons,
		ChurnByPlan:          churnMetrics.ChurnByPlan,
		CustomerLifetimeValue: churnMetrics.CustomerLifetimeValue,
	}, nil
}