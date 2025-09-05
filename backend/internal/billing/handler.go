package billing

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BillingHandler handles billing HTTP requests
type BillingHandler struct {
	service BillingService
}

// NewBillingHandler creates a new billing handler
func NewBillingHandler(service BillingService) *BillingHandler {
	return &BillingHandler{service: service}
}

// RegisterRoutes registers billing routes
func (h *BillingHandler) RegisterRoutes(r *gin.RouterGroup) {
	billing := r.Group("/billing")
	{
		// Billing Plans
		billing.GET("/plans", h.GetBillingPlans)
		billing.GET("/plans/:planId", h.GetBillingPlan)
		billing.POST("/plans", h.CreateBillingPlan)
		billing.PUT("/plans/:planId", h.UpdateBillingPlan)
		billing.DELETE("/plans/:planId", h.DeleteBillingPlan)

		// Subscriptions
		billing.POST("/subscriptions", h.CreateSubscription)
		billing.GET("/subscriptions", h.GetSubscription)
		billing.PUT("/subscriptions", h.UpdateSubscription)
		billing.DELETE("/subscriptions", h.CancelSubscription)
		billing.POST("/subscriptions/upgrade", h.UpgradePlan)
		billing.POST("/subscriptions/downgrade", h.DowngradePlan)

		// Usage Tracking
		billing.POST("/usage", h.RecordUsage)
		billing.GET("/usage", h.GetUsageSummary)
		billing.GET("/usage/limits", h.CheckUsageLimits)

		// Invoices
		billing.GET("/invoices", h.GetInvoices)
		billing.GET("/invoices/:invoiceId", h.GetInvoice)
		billing.POST("/invoices/:invoiceId/payment", h.ProcessPayment)
		billing.POST("/invoices/:invoiceId/refund", h.RefundPayment)

		// Analytics and Reporting
		billing.GET("/analytics", h.GetBillingAnalytics)
		billing.GET("/reports/revenue", h.GetRevenueReport)
		billing.GET("/reports/churn", h.GetChurnAnalysis)

		// Admin endpoints
		admin := billing.Group("/admin")
		{
			admin.POST("/process-billing", h.ProcessRecurringBilling)
			admin.POST("/retry-payments", h.RetryFailedPayments)
			admin.POST("/process-dunning", h.ProcessDunning)
			admin.POST("/tenants/:tenantId/suspend", h.SuspendService)
			admin.POST("/tenants/:tenantId/reactivate", h.ReactivateService)
		}
	}
}

// Request/Response DTOs
type CreateSubscriptionRequest struct {
	PlanID          uuid.UUID `json:"plan_id" binding:"required"`
	PaymentMethodID *string   `json:"payment_method_id,omitempty"`
}

type UpdateSubscriptionRequest struct {
	PaymentMethodID *string       `json:"payment_method_id,omitempty"`
	BillingCycle    *BillingCycle `json:"billing_cycle,omitempty"`
}

type CancelSubscriptionRequest struct {
	Reason      string `json:"reason" binding:"required"`
	Immediately bool   `json:"immediately"`
}

type ChangePlanRequest struct {
	NewPlanID uuid.UUID `json:"new_plan_id" binding:"required"`
}

type RecordUsageRequest struct {
	UsageType UsageType              `json:"usage_type" binding:"required"`
	Quantity  int64                  `json:"quantity" binding:"required"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type ProcessPaymentRequest struct {
	PaymentMethodID string `json:"payment_method_id" binding:"required"`
}

type RefundPaymentRequest struct {
	Amount float64 `json:"amount" binding:"required"`
	Reason string  `json:"reason" binding:"required"`
}

type CreateBillingPlanRequest struct {
	Name            string                 `json:"name" binding:"required"`
	Description     string                 `json:"description"`
	BasePrice       float64                `json:"base_price" binding:"required"`
	Currency        string                 `json:"currency"`
	BillingCycle    BillingCycle           `json:"billing_cycle"`
	Limits          map[string]interface{} `json:"limits"`
	Features        []string               `json:"features"`
	TrialPeriodDays int                    `json:"trial_period_days"`
	IsActive        bool                   `json:"is_active"`
	IsPublic        bool                   `json:"is_public"`
	UsageTiers      []CreateUsageTierRequest `json:"usage_tiers"`
}

type CreateUsageTierRequest struct {
	UsageType    UsageType `json:"usage_type" binding:"required"`
	MinUnits     int64     `json:"min_units" binding:"required"`
	MaxUnits     *int64    `json:"max_units,omitempty"`
	PricePerUnit float64   `json:"price_per_unit" binding:"required"`
}

// Billing Plans endpoints
func (h *BillingHandler) GetBillingPlans(c *gin.Context) {
	filter := PlanFilter{
		Limit:  50,
		Offset: 0,
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	if activeStr := c.Query("active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			filter.IsActive = &active
		}
	}

	if publicStr := c.Query("public"); publicStr != "" {
		if public, err := strconv.ParseBool(publicStr); err == nil {
			filter.IsPublic = &public
		}
	}

	if currency := c.Query("currency"); currency != "" {
		filter.Currency = &currency
	}

	plans, err := h.service.GetBillingPlans(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get billing plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": plans})
}

func (h *BillingHandler) GetBillingPlan(c *gin.Context) {
	planIDStr := c.Param("planId")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	plan, err := h.service.GetBillingPlan(c.Request.Context(), planID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": plan})
}

func (h *BillingHandler) CreateBillingPlan(c *gin.Context) {
	var req CreateBillingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan := &BillingPlan{
		ID:              uuid.New(),
		Name:            req.Name,
		Description:     req.Description,
		BasePrice:       req.BasePrice,
		Currency:        req.Currency,
		BillingCycle:    req.BillingCycle,
		Limits:          req.Limits,
		Features:        req.Features,
		TrialPeriodDays: req.TrialPeriodDays,
		IsActive:        req.IsActive,
		IsPublic:        req.IsPublic,
	}

	if plan.Currency == "" {
		plan.Currency = "BDT"
	}
	if plan.BillingCycle == "" {
		plan.BillingCycle = BillingCycleMonthly
	}

	err := h.service.CreateBillingPlan(c.Request.Context(), plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create billing plan"})
		return
	}

	// Create usage tiers
	for _, tierReq := range req.UsageTiers {
		tier := &UsageTier{
			ID:            uuid.New(),
			BillingPlanID: plan.ID,
			UsageType:     tierReq.UsageType,
			MinUnits:      tierReq.MinUnits,
			MaxUnits:      tierReq.MaxUnits,
			PricePerUnit:  tierReq.PricePerUnit,
		}

		err := h.service.CreateUsageTier(c.Request.Context(), tier)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create usage tier"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": plan})
}

func (h *BillingHandler) UpdateBillingPlan(c *gin.Context) {
	planIDStr := c.Param("planId")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var req CreateBillingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.service.GetBillingPlan(c.Request.Context(), planID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	// Update plan fields
	plan.Name = req.Name
	plan.Description = req.Description
	plan.BasePrice = req.BasePrice
	if req.Currency != "" {
		plan.Currency = req.Currency
	}
	if req.BillingCycle != "" {
		plan.BillingCycle = req.BillingCycle
	}
	plan.Limits = req.Limits
	plan.Features = req.Features
	plan.TrialPeriodDays = req.TrialPeriodDays
	plan.IsActive = req.IsActive
	plan.IsPublic = req.IsPublic

	err = h.service.UpdateBillingPlan(c.Request.Context(), plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update billing plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": plan})
}

func (h *BillingHandler) DeleteBillingPlan(c *gin.Context) {
	planIDStr := c.Param("planId")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	err = h.service.DeleteBillingPlan(c.Request.Context(), planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete billing plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan deleted successfully"})
}

// Subscription endpoints
func (h *BillingHandler) CreateSubscription(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.service.CreateSubscription(c.Request.Context(), tenantID, req.PlanID, req.PaymentMethodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": subscription})
}

func (h *BillingHandler) GetSubscription(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	subscription, err := h.service.GetSubscription(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

func (h *BillingHandler) UpdateSubscription(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := SubscriptionUpdate{
		PaymentMethodID: req.PaymentMethodID,
		BillingCycle:    req.BillingCycle,
	}

	subscription, err := h.service.UpdateSubscription(c.Request.Context(), tenantID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

func (h *BillingHandler) CancelSubscription(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	var req CancelSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CancelSubscription(c.Request.Context(), tenantID, req.Reason, req.Immediately)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription canceled successfully"})
}

func (h *BillingHandler) UpgradePlan(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	var req ChangePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.service.UpgradePlan(c.Request.Context(), tenantID, req.NewPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

func (h *BillingHandler) DowngradePlan(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	var req ChangePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.service.DowngradePlan(c.Request.Context(), tenantID, req.NewPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

// Usage endpoints
func (h *BillingHandler) RecordUsage(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	var req RecordUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.RecordUsage(c.Request.Context(), tenantID, req.UsageType, req.Quantity, req.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usage recorded successfully"})
}

func (h *BillingHandler) GetUsageSummary(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
	} else {
		// Default to current month
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
	} else {
		// Default to end of current month
		endDate = startDate.AddDate(0, 1, -1)
	}

	summary, err := h.service.GetUsageSummary(c.Request.Context(), tenantID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
		"period": gin.H{
			"start": startDate.Format("2006-01-02"),
			"end":   endDate.Format("2006-01-02"),
		},
	})
}

func (h *BillingHandler) CheckUsageLimits(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	status, err := h.service.CheckUsageLimits(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": status})
}

// Invoice endpoints
func (h *BillingHandler) GetInvoices(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	filter := InvoiceFilter{
		Limit:  50,
		Offset: 0,
		SortBy: "created_at",
		SortDesc: true,
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	if status := c.Query("status"); status != "" {
		invoiceStatus := InvoiceStatus(status)
		filter.Status = &invoiceStatus
	}

	invoices, total, err := h.service.GetInvoices(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": invoices,
		"pagination": gin.H{
			"total":  total,
			"limit":  filter.Limit,
			"offset": filter.Offset,
		},
	})
}

func (h *BillingHandler) GetInvoice(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	invoiceIDStr := c.Param("invoiceId")
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	invoice, err := h.service.GetInvoice(c.Request.Context(), tenantID, invoiceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": invoice})
}

func (h *BillingHandler) ProcessPayment(c *gin.Context) {
	tenantID := h.getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID required"})
		return
	}

	invoiceIDStr := c.Param("invoiceId")
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	var req ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attempt, err := h.service.ProcessPayment(c.Request.Context(), invoiceID, req.PaymentMethodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": attempt})
}

func (h *BillingHandler) RefundPayment(c *gin.Context) {
	invoiceIDStr := c.Param("invoiceId")
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	var req RefundPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.RefundPayment(c.Request.Context(), invoiceID, req.Amount, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment refunded successfully"})
}

// Analytics endpoints
func (h *BillingHandler) GetBillingAnalytics(c *gin.Context) {
	filter := h.parseAnalyticsFilter(c)

	analytics, err := h.service.GetBillingAnalytics(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": analytics})
}

func (h *BillingHandler) GetRevenueReport(c *gin.Context) {
	filter := h.parseRevenueReportFilter(c)

	report, err := h.service.GetRevenueReport(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": report})
}

func (h *BillingHandler) GetChurnAnalysis(c *gin.Context) {
	periodStr := c.Query("period")
	if periodStr == "" {
		periodStr = "30d"
	}

	var period time.Duration
	switch periodStr {
	case "7d":
		period = 7 * 24 * time.Hour
	case "30d":
		period = 30 * 24 * time.Hour
	case "90d":
		period = 90 * 24 * time.Hour
	case "1y":
		period = 365 * 24 * time.Hour
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period. Use 7d, 30d, 90d, or 1y"})
		return
	}

	analysis, err := h.service.GetChurnAnalysis(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": analysis})
}

// Admin endpoints
func (h *BillingHandler) ProcessRecurringBilling(c *gin.Context) {
	err := h.service.ProcessRecurringBilling(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recurring billing processed successfully"})
}

func (h *BillingHandler) RetryFailedPayments(c *gin.Context) {
	err := h.service.RetryFailedPayments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Failed payments retry process completed"})
}

func (h *BillingHandler) ProcessDunning(c *gin.Context) {
	err := h.service.ProcessDunning(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dunning process completed"})
}

func (h *BillingHandler) SuspendService(c *gin.Context) {
	tenantIDStr := c.Param("tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	reason := c.Query("reason")
	if reason == "" {
		reason = "Administrative action"
	}

	err = h.service.SuspendService(c.Request.Context(), tenantID, reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service suspended successfully"})
}

func (h *BillingHandler) ReactivateService(c *gin.Context) {
	tenantIDStr := c.Param("tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	err = h.service.ReactivateService(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service reactivated successfully"})
}

// Helper methods
func (h *BillingHandler) getTenantID(c *gin.Context) uuid.UUID {
	// Extract tenant ID from JWT token or headers
	// This is a placeholder - implement based on your auth system
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		return uuid.Nil
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return uuid.Nil
	}

	return tenantID
}

func (h *BillingHandler) parseAnalyticsFilter(c *gin.Context) AnalyticsFilter {
	filter := AnalyticsFilter{
		StartDate: time.Now().AddDate(0, -1, 0), // Default to last month
		EndDate:   time.Now(),
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = endDate
		}
	}

	if currency := c.Query("currency"); currency != "" {
		filter.Currency = &currency
	}

	if planIDsStr := c.Query("plan_ids"); planIDsStr != "" {
		var planIDs []uuid.UUID
		if err := json.Unmarshal([]byte(planIDsStr), &planIDs); err == nil {
			filter.PlanIDs = planIDs
		}
	}

	return filter
}

func (h *BillingHandler) parseRevenueReportFilter(c *gin.Context) RevenueReportFilter {
	filter := RevenueReportFilter{
		StartDate:    time.Now().AddDate(0, -1, 0), // Default to last month
		EndDate:      time.Now(),
		GroupBy:      "month",
		IncludeUsage: true,
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = endDate
		}
	}

	if groupBy := c.Query("group_by"); groupBy != "" {
		filter.GroupBy = groupBy
	}

	if currency := c.Query("currency"); currency != "" {
		filter.Currency = &currency
	}

	if includeUsageStr := c.Query("include_usage"); includeUsageStr != "" {
		if includeUsage, err := strconv.ParseBool(includeUsageStr); err == nil {
			filter.IncludeUsage = includeUsage
		}
	}

	return filter
}
