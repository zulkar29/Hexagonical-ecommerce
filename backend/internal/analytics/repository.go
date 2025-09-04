package analytics

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Event operations
	CreateEvent(ctx context.Context, event *AnalyticsEvent) error
	CreatePageView(ctx context.Context, pageView *PageView) error
	CreateProductView(ctx context.Context, productView *ProductView) error
	CreatePurchase(ctx context.Context, purchase *Purchase) error

	// Basic analytics
	GetDashboardStats(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*AnalyticsStats, error)
	GetTrafficStats(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*TrafficStats, error)
	GetSalesStats(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*SalesStats, error)

	// Top performers
	GetTopProducts(ctx context.Context, tenantID uuid.UUID, dateRange DateRange, limit int) ([]*ProductStats, error)
	GetTopPages(ctx context.Context, tenantID uuid.UUID, dateRange DateRange, limit int) ([]*PageStats, error)
	GetTopReferrers(ctx context.Context, tenantID uuid.UUID, dateRange DateRange, limit int) ([]*ReferrerStats, error)

	// Advanced analytics
	GetCohortAnalysis(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*CohortAnalysis, error)
	GetFunnelAnalysis(ctx context.Context, tenantID uuid.UUID, funnelSteps []string, dateRange DateRange) (*FunnelAnalysis, error)
	GetCustomerLifetimeValue(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (float64, error)
	GetRetentionRate(ctx context.Context, tenantID uuid.UUID, days int) (float64, error)

	// Real-time analytics
	GetRealTimeStats(ctx context.Context, tenantID uuid.UUID) (*RealTimeStats, error)
	GetActiveUsers(ctx context.Context, tenantID uuid.UUID) (int64, error)

	// Scheduled reports
	CreateScheduledReport(ctx context.Context, tenantID uuid.UUID, request ScheduleReportRequest) (*ScheduledReport, error)
	GetScheduledReports(ctx context.Context, tenantID uuid.UUID) ([]*ScheduledReport, error)
	UpdateScheduledReport(ctx context.Context, report *ScheduledReport) error
	DeleteScheduledReport(ctx context.Context, tenantID, reportID uuid.UUID) error

	// Data export
	ExportData(ctx context.Context, tenantID uuid.UUID, request ExportRequest) ([]byte, string, error)
}

type repository struct {
	db *gorm.DB
	// TODO: Add cache client for real-time data
	// cache redis.Client
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Event operations
func (r *repository) CreateEvent(ctx context.Context, event *AnalyticsEvent) error {
	// TODO: Implement event creation with proper indexing
	// - Validate event structure
	// - Add tenant isolation
	// - Handle batch inserts for high volume
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *repository) CreatePageView(ctx context.Context, pageView *PageView) error {
	// TODO: Implement page view creation
	// - Update previous page view duration if exists
	// - Create new page view record
	return r.db.WithContext(ctx).Create(pageView).Error
}

func (r *repository) CreateProductView(ctx context.Context, productView *ProductView) error {
	// TODO: Implement product view creation
	return r.db.WithContext(ctx).Create(productView).Error
}

func (r *repository) CreatePurchase(ctx context.Context, purchase *Purchase) error {
	// TODO: Implement purchase event creation
	return r.db.WithContext(ctx).Create(purchase).Error
}

// Basic analytics
func (r *repository) GetDashboardStats(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*AnalyticsStats, error) {
	// TODO: Implement comprehensive dashboard stats query
	// This would be a complex query aggregating data from multiple tables
	stats := &AnalyticsStats{
		TenantID: tenantID,
		Date:     dateRange.Start,
	}

	// TODO: Execute complex aggregation queries
	// - Count page views in date range
	// - Count unique visitors
	// - Calculate bounce rate
	// - Get revenue and conversion metrics
	// - Get top performers

	return stats, nil
}

func (r *repository) GetTrafficStats(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*TrafficStats, error) {
	// TODO: Implement traffic analytics query
	stats := &TrafficStats{
		TenantID:  tenantID,
		DateRange: dateRange,
	}

	// Example query structure (to be implemented)
	// SELECT COUNT(*) as page_views, 
	//        COUNT(DISTINCT anonymous_id, user_id) as unique_visitors,
	//        COUNT(DISTINCT session_id) as sessions
	// FROM page_views 
	// WHERE tenant_id = ? AND timestamp BETWEEN ? AND ?

	var pageViews, uniqueVisitors, sessions int64
	err := r.db.WithContext(ctx).
		Model(&PageView{}).
		Where("tenant_id = ? AND timestamp BETWEEN ? AND ?", tenantID, dateRange.Start, dateRange.End).
		Select("COUNT(*) as page_views").
		Row().
		Scan(&pageViews)
	if err != nil {
		return nil, err
	}

	stats.PageViews = pageViews
	stats.UniqueVisitors = uniqueVisitors // TODO: Implement proper unique visitor count
	stats.Sessions = sessions             // TODO: Implement session count

	// TODO: Calculate bounce rate, avg session time, etc.

	return stats, nil
}

func (r *repository) GetSalesStats(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*SalesStats, error) {
	// TODO: Implement sales analytics query
	stats := &SalesStats{
		TenantID:  tenantID,
		DateRange: dateRange,
	}

	// Query purchase events for sales metrics
	var totalRevenue float64
	var totalOrders int64

	err := r.db.WithContext(ctx).
		Model(&Purchase{}).
		Where("tenant_id = ? AND timestamp BETWEEN ? AND ?", tenantID, dateRange.Start, dateRange.End).
		Select("SUM(total_amount) as total_revenue, COUNT(*) as total_orders").
		Row().
		Scan(&totalRevenue, &totalOrders)
	if err != nil {
		return nil, err
	}

	stats.TotalRevenue = totalRevenue
	stats.TotalOrders = totalOrders

	if totalOrders > 0 {
		stats.AvgOrderValue = totalRevenue / float64(totalOrders)
	}

	// TODO: Calculate conversion rate (requires page views and purchases correlation)
	// TODO: Get top products, sales by category, daily sales breakdown

	return stats, nil
}

// Top performers
func (r *repository) GetTopProducts(ctx context.Context, tenantID uuid.UUID, dateRange DateRange, limit int) ([]*ProductStats, error) {
	// TODO: Implement top products query
	// This would join product_views with purchases to get comprehensive product stats
	var products []*ProductStats

	// Example query structure:
	// SELECT pv.product_id, p.name, COUNT(pv.id) as views, COUNT(pu.id) as sales, SUM(pu.total_amount) as revenue
	// FROM product_views pv
	// LEFT JOIN purchases pu ON pv.product_id = pu.product_id AND pu.tenant_id = pv.tenant_id
	// LEFT JOIN products p ON pv.product_id = p.id
	// WHERE pv.tenant_id = ? AND pv.timestamp BETWEEN ? AND ?
	// GROUP BY pv.product_id, p.name
	// ORDER BY views DESC
	// LIMIT ?

	return products, nil
}

func (r *repository) GetTopPages(ctx context.Context, tenantID uuid.UUID, dateRange DateRange, limit int) ([]*PageStats, error) {
	// TODO: Implement top pages query
	var pages []*PageStats

	// Query page views grouped by path
	rows, err := r.db.WithContext(ctx).
		Model(&PageView{}).
		Where("tenant_id = ? AND timestamp BETWEEN ? AND ?", tenantID, dateRange.Start, dateRange.End).
		Select("path, COUNT(*) as views, COUNT(DISTINCT COALESCE(user_id::text, anonymous_id)) as unique_views, AVG(duration_seconds) as avg_time").
		Group("path").
		Order("views DESC").
		Limit(limit).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var page PageStats
		if err := rows.Scan(&page.Path, &page.Views, &page.UniqueViews, &page.AvgTime); err != nil {
			return nil, err
		}
		pages = append(pages, &page)
	}

	// TODO: Calculate bounce rate per page

	return pages, nil
}

func (r *repository) GetTopReferrers(ctx context.Context, tenantID uuid.UUID, dateRange DateRange, limit int) ([]*ReferrerStats, error) {
	// TODO: Implement top referrers query
	var referrers []*ReferrerStats

	rows, err := r.db.WithContext(ctx).
		Model(&PageView{}).
		Where("tenant_id = ? AND timestamp BETWEEN ? AND ? AND referrer != ''", tenantID, dateRange.Start, dateRange.End).
		Select("referrer, COUNT(*) as visits").
		Group("referrer").
		Order("visits DESC").
		Limit(limit).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totalVisits := int64(0)
	for rows.Next() {
		var referrer ReferrerStats
		if err := rows.Scan(&referrer.Referrer, &referrer.Visits); err != nil {
			return nil, err
		}
		totalVisits += referrer.Visits
		referrers = append(referrers, &referrer)
	}

	// Calculate percentages
	for _, ref := range referrers {
		ref.Percentage = float64(ref.Visits) / float64(totalVisits) * 100
	}

	return referrers, nil
}

// Advanced analytics
func (r *repository) GetCohortAnalysis(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*CohortAnalysis, error) {
	cohorts := &CohortAnalysis{
		TenantID: tenantID,
		Cohorts:  make(map[string]map[string]float64),
	}

	// Get first purchase month for each customer
	type customerCohort struct {
		CustomerID string
		CohortMonth string
		PurchaseMonth string
	}

	var cohortData []customerCohort
	rows, err := r.db.WithContext(ctx).Raw(`
		WITH first_purchases AS (
			SELECT 
				customer_id,
				DATE_TRUNC('month', MIN(timestamp)) as cohort_month
			FROM purchases 
			WHERE tenant_id = ? AND timestamp BETWEEN ? AND ?
			GROUP BY customer_id
		),
		all_purchases AS (
			SELECT 
				p.customer_id,
				fp.cohort_month,
				DATE_TRUNC('month', p.timestamp) as purchase_month
			FROM purchases p
			JOIN first_purchases fp ON p.customer_id = fp.customer_id
			WHERE p.tenant_id = ? AND p.timestamp BETWEEN ? AND ?
		)
		SELECT 
			customer_id,
			cohort_month::text,
			purchase_month::text
		FROM all_purchases
		ORDER BY cohort_month, purchase_month
	`, tenantID, dateRange.Start, dateRange.End, tenantID, dateRange.Start, dateRange.End).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cd customerCohort
		if err := rows.Scan(&cd.CustomerID, &cd.CohortMonth, &cd.PurchaseMonth); err != nil {
			return nil, err
		}
		cohortData = append(cohortData, cd)
	}

	// Calculate retention rates
	cohortCounts := make(map[string]map[string]int)
	cohortSizes := make(map[string]int)

	for _, cd := range cohortData {
		if cohortCounts[cd.CohortMonth] == nil {
			cohortCounts[cd.CohortMonth] = make(map[string]int)
		}
		cohortCounts[cd.CohortMonth][cd.PurchaseMonth]++
		cohortSizes[cd.CohortMonth]++
	}

	// Convert to percentages
	for cohortMonth, purchases := range cohortCounts {
		if cohorts.Cohorts[cohortMonth] == nil {
			cohorts.Cohorts[cohortMonth] = make(map[string]float64)
		}
		for purchaseMonth, count := range purchases {
			retentionRate := float64(count) / float64(cohortSizes[cohortMonth]) * 100
			cohorts.Cohorts[cohortMonth][purchaseMonth] = retentionRate
		}
	}

	return cohorts, nil
}

func (r *repository) GetFunnelAnalysis(ctx context.Context, tenantID uuid.UUID, funnelSteps []string, dateRange DateRange) (*FunnelAnalysis, error) {
	funnel := &FunnelAnalysis{
		TenantID: tenantID,
		Steps:    make([]FunnelStep, 0, len(funnelSteps)),
	}

	if len(funnelSteps) == 0 {
		// Default funnel steps
		funnelSteps = []string{"page_view", "product_view", "add_to_cart", "checkout", "purchase"}
	}

	// Count users at each step
	for i, step := range funnelSteps {
		var userCount int64
		var conversionRate float64

		switch step {
		case "page_view":
			err := r.db.WithContext(ctx).
				Model(&PageView{}).
				Where("tenant_id = ? AND timestamp BETWEEN ? AND ?", tenantID, dateRange.Start, dateRange.End).
				Select("COUNT(DISTINCT COALESCE(user_id::text, anonymous_id))").
				Row().
				Scan(&userCount)
			if err != nil {
				return nil, err
			}
		case "product_view":
			err := r.db.WithContext(ctx).
				Model(&ProductView{}).
				Where("tenant_id = ? AND timestamp BETWEEN ? AND ?", tenantID, dateRange.Start, dateRange.End).
				Select("COUNT(DISTINCT COALESCE(user_id::text, anonymous_id))").
				Row().
				Scan(&userCount)
			if err != nil {
				return nil, err
			}
		case "add_to_cart":
			err := r.db.WithContext(ctx).
				Model(&Event{}).
				Where("tenant_id = ? AND event_type = ? AND timestamp BETWEEN ? AND ?", tenantID, "add_to_cart", dateRange.Start, dateRange.End).
				Select("COUNT(DISTINCT COALESCE(user_id::text, anonymous_id))").
				Row().
				Scan(&userCount)
			if err != nil {
				return nil, err
			}
		case "checkout":
			err := r.db.WithContext(ctx).
				Model(&Event{}).
				Where("tenant_id = ? AND event_type = ? AND timestamp BETWEEN ? AND ?", tenantID, "checkout_started", dateRange.Start, dateRange.End).
				Select("COUNT(DISTINCT COALESCE(user_id::text, anonymous_id))").
				Row().
				Scan(&userCount)
			if err != nil {
				return nil, err
			}
		case "purchase":
			err := r.db.WithContext(ctx).
				Model(&Purchase{}).
				Where("tenant_id = ? AND timestamp BETWEEN ? AND ?", tenantID, dateRange.Start, dateRange.End).
				Select("COUNT(DISTINCT customer_id)").
				Row().
				Scan(&userCount)
			if err != nil {
				return nil, err
			}
		default:
			// Custom event type
			err := r.db.WithContext(ctx).
				Model(&Event{}).
				Where("tenant_id = ? AND event_type = ? AND timestamp BETWEEN ? AND ?", tenantID, step, dateRange.Start, dateRange.End).
				Select("COUNT(DISTINCT COALESCE(user_id::text, anonymous_id))").
				Row().
				Scan(&userCount)
			if err != nil {
				return nil, err
			}
		}

		// Calculate conversion rate from previous step
		if i > 0 && len(funnel.Steps) > 0 {
			prevStepUsers := funnel.Steps[i-1].Users
			if prevStepUsers > 0 {
				conversionRate = float64(userCount) / float64(prevStepUsers) * 100
			}
		} else {
			conversionRate = 100.0 // First step is 100%
		}

		funnelStep := FunnelStep{
			Step:           step,
			Users:          userCount,
			ConversionRate: conversionRate,
		}

		funnel.Steps = append(funnel.Steps, funnelStep)
	}

	return funnel, nil
}

func (r *repository) GetCustomerLifetimeValue(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (float64, error) {
	// Formula: CLV = (Average Order Value) × (Purchase Frequency) × (Customer Lifespan)
	var clv float64

	// Get average order value
	var avgOrderValue float64
	err := r.db.WithContext(ctx).
		Model(&Purchase{}).
		Where("tenant_id = ? AND timestamp BETWEEN ? AND ?", tenantID, dateRange.Start, dateRange.End).
		Select("AVG(total_amount)").
		Row().
		Scan(&avgOrderValue)
	if err != nil {
		return 0, err
	}

	// Calculate purchase frequency (orders per customer)
	var totalOrders, uniqueCustomers int64
	err = r.db.WithContext(ctx).
		Model(&Purchase{}).
		Where("tenant_id = ? AND timestamp BETWEEN ? AND ?", tenantID, dateRange.Start, dateRange.End).
		Select("COUNT(*), COUNT(DISTINCT customer_id)").
		Row().
		Scan(&totalOrders, &uniqueCustomers)
	if err != nil || uniqueCustomers == 0 {
		return avgOrderValue, err
	}

	purchaseFrequency := float64(totalOrders) / float64(uniqueCustomers)

	// Calculate average customer lifespan in months
	type customerLifespan struct {
		CustomerID string
		FirstPurchase time.Time
		LastPurchase time.Time
	}

	var lifespans []customerLifespan
	rows, err := r.db.WithContext(ctx).
		Model(&Purchase{}).
		Where("tenant_id = ? AND timestamp BETWEEN ? AND ?", tenantID, dateRange.Start, dateRange.End).
		Select("customer_id, MIN(timestamp) as first_purchase, MAX(timestamp) as last_purchase").
		Group("customer_id").
		Rows()
	if err != nil {
		return avgOrderValue * purchaseFrequency, err
	}
	defer rows.Close()

	var totalLifespanMonths float64
	customerCount := 0
	for rows.Next() {
		var ls customerLifespan
		if err := rows.Scan(&ls.CustomerID, &ls.FirstPurchase, &ls.LastPurchase); err != nil {
			continue
		}
		lifespanMonths := ls.LastPurchase.Sub(ls.FirstPurchase).Hours() / (24 * 30) // Approximate months
		if lifespanMonths < 1 {
			lifespanMonths = 1 // Minimum 1 month
		}
		totalLifespanMonths += lifespanMonths
		customerCount++
	}

	avgLifespanMonths := float64(1) // Default to 1 month
	if customerCount > 0 {
		avgLifespanMonths = totalLifespanMonths / float64(customerCount)
	}

	// Calculate CLV
	clv = avgOrderValue * purchaseFrequency * avgLifespanMonths

	return clv, nil
}

func (r *repository) GetRetentionRate(ctx context.Context, tenantID uuid.UUID, days int) (float64, error) {
	// Calculate percentage of customers who return within X days
	var retentionRate float64

	// Get customers who made their first purchase in the period
	startDate := time.Now().AddDate(0, 0, -days*2) // Look back twice the retention period
	midDate := time.Now().AddDate(0, 0, -days)     // Retention period start
	endDate := time.Now()                          // Current time

	// Get customers who made first purchase in the first period
	type firstPurchaseCustomer struct {
		CustomerID string
		FirstPurchase time.Time
	}

	var firstPurchaseCustomers []firstPurchaseCustomer
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT 
			customer_id,
			MIN(timestamp) as first_purchase
		FROM purchases 
		WHERE tenant_id = ? AND timestamp BETWEEN ? AND ?
		GROUP BY customer_id
		HAVING MIN(timestamp) BETWEEN ? AND ?
	`, tenantID, startDate, endDate, startDate, midDate).Rows()

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var fpc firstPurchaseCustomer
		if err := rows.Scan(&fpc.CustomerID, &fpc.FirstPurchase); err != nil {
			continue
		}
		firstPurchaseCustomers = append(firstPurchaseCustomers, fpc)
	}

	if len(firstPurchaseCustomers) == 0 {
		return 0, nil
	}

	// Check how many of these customers made another purchase within the retention period
	returnedCustomers := 0
	for _, fpc := range firstPurchaseCustomers {
		retentionStart := fpc.FirstPurchase
		retentionEnd := fpc.FirstPurchase.AddDate(0, 0, days)

		var hasReturnPurchase int64
		err := r.db.WithContext(ctx).
			Model(&Purchase{}).
			Where("tenant_id = ? AND customer_id = ? AND timestamp > ? AND timestamp <= ?", 
				tenantID, fpc.CustomerID, retentionStart, retentionEnd).
			Count(&hasReturnPurchase)
		if err != nil {
			continue
		}

		if hasReturnPurchase > 0 {
			returnedCustomers++
		}
	}

	// Calculate retention rate
	retentionRate = float64(returnedCustomers) / float64(len(firstPurchaseCustomers)) * 100

	return retentionRate, nil
}

// Real-time analytics
func (r *repository) GetRealTimeStats(ctx context.Context, tenantID uuid.UUID) (*RealTimeStats, error) {
	stats := &RealTimeStats{
		LastUpdated: time.Now(),
	}

	// Get active users (last 5 minutes)
	activeUsers, err := r.GetActiveUsers(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	stats.ActiveUsers = activeUsers

	// Get page views in last hour
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	var pageViewsLastHour int64
	err = r.db.WithContext(ctx).
		Model(&PageView{}).
		Where("tenant_id = ? AND timestamp >= ?", tenantID, oneHourAgo).
		Count(&pageViewsLastHour)
	if err != nil {
		return nil, err
	}
	stats.PageViewsLastHour = pageViewsLastHour

	// Get active pages (most viewed in last hour)
	rows, err := r.db.WithContext(ctx).
		Model(&PageView{}).
		Where("tenant_id = ? AND timestamp >= ?", tenantID, oneHourAgo).
		Select("path, COUNT(*) as views").
		Group("path").
		Order("views DESC").
		Limit(10).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var page ActivePageStats
		if err := rows.Scan(&page.Path, &page.Views); err != nil {
			continue
		}
		stats.ActivePages = append(stats.ActivePages, page)
	}

	// Get recent conversions (last hour)
	var conversionsLastHour int64
	err = r.db.WithContext(ctx).
		Model(&Purchase{}).
		Where("tenant_id = ? AND timestamp >= ?", tenantID, oneHourAgo).
		Count(&conversionsLastHour)
	if err != nil {
		return nil, err
	}
	stats.ConversionsLastHour = conversionsLastHour

	return stats, nil
}

func (r *repository) GetActiveUsers(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	// Count users active in last 5 minutes
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	var activeUsers int64

	err := r.db.WithContext(ctx).
		Model(&PageView{}).
		Where("tenant_id = ? AND timestamp >= ?", tenantID, fiveMinutesAgo).
		Select("COUNT(DISTINCT COALESCE(user_id::text, anonymous_id))").
		Row().
		Scan(&activeUsers)

	return activeUsers, err
}

// Scheduled reports
func (r *repository) CreateScheduledReport(ctx context.Context, tenantID uuid.UUID, request ScheduleReportRequest) (*ScheduledReport, error) {
	report := &ScheduledReport{
		ID:         uuid.New(),
		TenantID:   tenantID,
		Name:       request.Name,
		ReportType: request.ReportType,
		Format:     request.Format,
		Frequency:  request.Frequency,
		Recipients: request.Recipients,
		Filters:    request.Filters,
		IsActive:   request.IsActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Calculate next run time based on frequency
	now := time.Now()
	switch request.Frequency {
	case "daily":
		report.NextRun = now.AddDate(0, 0, 1)
	case "weekly":
		report.NextRun = now.AddDate(0, 0, 7)
	case "monthly":
		report.NextRun = now.AddDate(0, 1, 0)
	default:
		report.NextRun = now.AddDate(0, 0, 1) // Default to daily
	}

	err := r.db.WithContext(ctx).Create(report).Error
	return report, err
}

func (r *repository) GetScheduledReports(ctx context.Context, tenantID uuid.UUID) ([]*ScheduledReport, error) {
	var reports []*ScheduledReport
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&reports).Error
	return reports, err
}

func (r *repository) UpdateScheduledReport(ctx context.Context, report *ScheduledReport) error {
	report.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *repository) DeleteScheduledReport(ctx context.Context, tenantID, reportID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, reportID).
		Delete(&ScheduledReport{}).Error
}

// Data export
func (r *repository) ExportData(ctx context.Context, tenantID uuid.UUID, request ExportRequest) ([]byte, string, error) {
	// Based on request.DataType, query appropriate tables and export data
	switch request.DataType {
	case "events":
		return r.exportEvents(ctx, tenantID, request)
	case "pageviews":
		return r.exportPageViews(ctx, tenantID, request)
	case "purchases":
		return r.exportPurchases(ctx, tenantID, request)
	default:
		return nil, "", fmt.Errorf("unsupported data type: %s", request.DataType)
	}
}

func (r *repository) exportEvents(ctx context.Context, tenantID uuid.UUID, request ExportRequest) ([]byte, string, error) {
	var events []Event
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	// Apply date range filter if provided
	if !request.DateRange.Start.IsZero() && !request.DateRange.End.IsZero() {
		query = query.Where("timestamp BETWEEN ? AND ?", request.DateRange.Start, request.DateRange.End)
	}

	// Apply limit if provided
	if request.Limit > 0 {
		query = query.Limit(request.Limit)
	}

	err := query.Order("timestamp DESC").Find(&events).Error
	if err != nil {
		return nil, "", err
	}

	switch request.Format {
	case "csv":
		return r.eventsToCSV(events)
	default:
		return r.eventsToJSON(events)
	}
}

func (r *repository) exportPageViews(ctx context.Context, tenantID uuid.UUID, request ExportRequest) ([]byte, string, error) {
	var pageViews []PageView
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	// Apply date range filter if provided
	if !request.DateRange.Start.IsZero() && !request.DateRange.End.IsZero() {
		query = query.Where("timestamp BETWEEN ? AND ?", request.DateRange.Start, request.DateRange.End)
	}

	// Apply limit if provided
	if request.Limit > 0 {
		query = query.Limit(request.Limit)
	}

	err := query.Order("timestamp DESC").Find(&pageViews).Error
	if err != nil {
		return nil, "", err
	}

	switch request.Format {
	case "csv":
		return r.pageViewsToCSV(pageViews)
	default:
		return r.pageViewsToJSON(pageViews)
	}
}

func (r *repository) exportPurchases(ctx context.Context, tenantID uuid.UUID, request ExportRequest) ([]byte, string, error) {
	var purchases []Purchase
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	// Apply date range filter if provided
	if !request.DateRange.Start.IsZero() && !request.DateRange.End.IsZero() {
		query = query.Where("timestamp BETWEEN ? AND ?", request.DateRange.Start, request.DateRange.End)
	}

	// Apply limit if provided
	if request.Limit > 0 {
		query = query.Limit(request.Limit)
	}

	err := query.Order("timestamp DESC").Find(&purchases).Error
	if err != nil {
		return nil, "", err
	}

	switch request.Format {
	case "csv":
		return r.purchasesToCSV(purchases)
	default:
		return r.purchasesToJSON(purchases)
	}
}

// Helper methods for data conversion
func (r *repository) eventsToJSON(events []Event) ([]byte, string, error) {
	data, err := json.Marshal(events)
	return data, "application/json", err
}

func (r *repository) eventsToCSV(events []Event) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "TenantID", "EventType", "UserID", "AnonymousID", "Properties", "Timestamp"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	// Write data
	for _, event := range events {
		properties, _ := json.Marshal(event.Properties)
		userID := ""
		if event.UserID != nil {
			userID = event.UserID.String()
		}
		row := []string{
			event.ID.String(),
			event.TenantID.String(),
			event.EventType,
			userID,
			event.AnonymousID,
			string(properties),
			event.Timestamp.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "text/csv", writer.Error()
}

func (r *repository) pageViewsToJSON(pageViews []PageView) ([]byte, string, error) {
	data, err := json.Marshal(pageViews)
	return data, "application/json", err
}

func (r *repository) pageViewsToCSV(pageViews []PageView) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "TenantID", "Path", "UserID", "AnonymousID", "Referrer", "UserAgent", "DurationSeconds", "Timestamp"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	// Write data
	for _, pv := range pageViews {
		userID := ""
		if pv.UserID != nil {
			userID = pv.UserID.String()
		}
		row := []string{
			pv.ID.String(),
			pv.TenantID.String(),
			pv.Path,
			userID,
			pv.AnonymousID,
			pv.Referrer,
			pv.UserAgent,
			fmt.Sprintf("%.2f", pv.DurationSeconds),
			pv.Timestamp.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "text/csv", writer.Error()
}

func (r *repository) purchasesToJSON(purchases []Purchase) ([]byte, string, error) {
	data, err := json.Marshal(purchases)
	return data, "application/json", err
}

func (r *repository) purchasesToCSV(purchases []Purchase) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "TenantID", "CustomerID", "OrderID", "ProductID", "Quantity", "UnitPrice", "TotalAmount", "Currency", "Timestamp"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	// Write data
	for _, purchase := range purchases {
		row := []string{
			purchase.ID.String(),
			purchase.TenantID.String(),
			purchase.CustomerID,
			purchase.OrderID,
			purchase.ProductID,
			fmt.Sprintf("%d", purchase.Quantity),
			fmt.Sprintf("%.2f", purchase.UnitPrice),
			fmt.Sprintf("%.2f", purchase.TotalAmount),
			purchase.Currency,
			purchase.Timestamp.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "text/csv", writer.Error()
}
