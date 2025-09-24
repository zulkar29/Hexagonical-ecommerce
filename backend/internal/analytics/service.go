package analytics

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type TrafficStats struct {
	TenantID       uuid.UUID `json:"tenant_id"`
	DateRange      DateRange `json:"date_range"`
	PageViews      int64     `json:"page_views"`
	UniqueVisitors int64     `json:"unique_visitors"`
	Sessions       int64     `json:"sessions"`
	BounceRate     float64   `json:"bounce_rate"`
	AvgSessionTime float64   `json:"avg_session_time"`
	NewVisitors    int64     `json:"new_visitors"`
	ReturningVisitors int64  `json:"returning_visitors"`
	TopPages       []PageStats `json:"top_pages"`
	TopReferrers   []ReferrerStats `json:"top_referrers"`
}

type SalesStats struct {
	TenantID         uuid.UUID `json:"tenant_id"`
	DateRange        DateRange `json:"date_range"`
	TotalRevenue     float64   `json:"total_revenue"`
	TotalOrders      int64     `json:"total_orders"`
	AvgOrderValue    float64   `json:"avg_order_value"`
	ConversionRate   float64   `json:"conversion_rate"`
	RefundRate       float64   `json:"refund_rate"`
	TopProducts      []ProductStats `json:"top_products"`
	SalesByCategory  []CategoryStats `json:"sales_by_category"`
	SalesByDay       []DailySales   `json:"sales_by_day"`
}

type PageStats struct {
	Path        string  `json:"path"`
	Views       int64   `json:"views"`
	UniqueViews int64   `json:"unique_views"`
	BounceRate  float64 `json:"bounce_rate"`
	AvgTime     float64 `json:"avg_time"`
}

type ProductStats struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Views       int64     `json:"views"`
	Sales       int64     `json:"sales"`
	Revenue     float64   `json:"revenue"`
	ConversionRate float64 `json:"conversion_rate"`
}

type ReferrerStats struct {
	Referrer string `json:"referrer"`
	Visits   int64  `json:"visits"`
	Percentage float64 `json:"percentage"`
}

type CategoryStats struct {
	Category string  `json:"category"`
	Revenue  float64 `json:"revenue"`
	Orders   int64   `json:"orders"`
}

type DailySales struct {
	Date    time.Time `json:"date"`
	Revenue float64   `json:"revenue"`
	Orders  int64     `json:"orders"`
}

type CohortAnalysis struct {
	TenantID uuid.UUID `json:"tenant_id"`
	Cohorts  []CohortData `json:"cohorts"`
}

type CohortData struct {
	CohortMonth time.Time `json:"cohort_month"`
	Customers   int64     `json:"customers"`
	Retention   []float64 `json:"retention"` // Month 0, 1, 2, 3...
}

type FunnelAnalysis struct {
	TenantID uuid.UUID   `json:"tenant_id"`
	Steps    []FunnelStep `json:"steps"`
	OverallConversion float64 `json:"overall_conversion"`
}

type FunnelStep struct {
	Step        string  `json:"step"`
	Users       int64   `json:"users"`
	Conversion  float64 `json:"conversion"`
	Dropoff     float64 `json:"dropoff"`
}

type ReportFormat string

const (
	ReportFormatJSON ReportFormat = "json"
	ReportFormatCSV  ReportFormat = "csv"
	ReportFormatPDF  ReportFormat = "pdf"
)

type ReportRequest struct {
	ReportType string       `json:"report_type"`
	DateRange  DateRange    `json:"date_range"`
	Format     ReportFormat `json:"format"`
	Filters    map[string]interface{} `json:"filters,omitempty"`
}

type Service interface {
	// Event tracking
	TrackEvent(ctx context.Context, tenantID uuid.UUID, event *AnalyticsEvent) error
	TrackPageView(ctx context.Context, tenantID uuid.UUID, pageView *PageView) error
	TrackProductView(ctx context.Context, tenantID uuid.UUID, productView *ProductView) error
	TrackPurchase(ctx context.Context, tenantID uuid.UUID, purchase *Purchase) error

	// Basic stats
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

	// Reports
	GenerateReport(ctx context.Context, tenantID uuid.UUID, request ReportRequest) ([]byte, string, error)
	ScheduleReport(ctx context.Context, tenantID uuid.UUID, request ScheduleReportRequest) (*ScheduledReport, error)
	GetScheduledReports(ctx context.Context, tenantID uuid.UUID) ([]*ScheduledReport, error)
	DeleteScheduledReport(ctx context.Context, tenantID uuid.UUID, reportID uuid.UUID) error

	// Removed ExportData - GDPR data export functionality not needed
}

type RealTimeStats struct {
	ActiveUsers    int64            `json:"active_users"`
	PageViews      int64            `json:"page_views"`
	ActivePages    []ActivePageStats `json:"active_pages"`
	Conversions    int64            `json:"conversions"`
	Revenue        float64          `json:"revenue"`
	LastUpdated    time.Time        `json:"last_updated"`
}

type ActivePageStats struct {
	Path        string `json:"path"`
	ActiveUsers int64  `json:"active_users"`
	Views       int64  `json:"views"`
}

type ScheduleReportRequest struct {
	Name        string       `json:"name"`
	ReportType  string       `json:"report_type"`
	Format      ReportFormat `json:"format"`
	Frequency   string       `json:"frequency"` // daily, weekly, monthly
	Recipients  []string     `json:"recipients"`
	Filters     map[string]interface{} `json:"filters,omitempty"`
	IsActive    bool         `json:"is_active"`
}

type ScheduledReport struct {
	ID          uuid.UUID    `json:"id"`
	TenantID    uuid.UUID    `json:"tenant_id"`
	Name        string       `json:"name"`
	ReportType  string       `json:"report_type"`
	Format      ReportFormat `json:"format"`
	Frequency   string       `json:"frequency"`
	Recipients  []string     `json:"recipients"`
	Filters     map[string]interface{} `json:"filters"`
	IsActive    bool         `json:"is_active"`
	NextRun     time.Time    `json:"next_run"`
	LastRun     *time.Time   `json:"last_run,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Removed ExportRequest struct - GDPR data export functionality not needed

type service struct {
	repo Repository
	// TODO: Add cache service for real-time stats
	// cacheService cache.Service
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Event tracking implementations
func (s *service) TrackEvent(ctx context.Context, tenantID uuid.UUID, event *AnalyticsEvent) error {
	// Validate event structure
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}
	if event.EventType == "" {
		return fmt.Errorf("event type is required")
	}
	if event.TenantID == uuid.Nil {
		event.TenantID = tenantID
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Store in repository
	err := s.repo.CreateEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	// TODO: Update real-time cache when cache service is implemented
	return nil
}

func (s *service) TrackPageView(ctx context.Context, tenantID uuid.UUID, pageView *PageView) error {
	// Validate page view data
	if pageView == nil {
		return fmt.Errorf("page view cannot be nil")
	}
	if pageView.Path == "" {
		return fmt.Errorf("page path is required")
	}
	if pageView.TenantID == uuid.Nil {
		pageView.TenantID = tenantID
	}
	if pageView.Timestamp.IsZero() {
		pageView.Timestamp = time.Now()
	}

	// Store page view
	err := s.repo.CreatePageView(ctx, pageView)
	if err != nil {
		return fmt.Errorf("failed to create page view: %w", err)
	}

	// TODO: Update session information and calculate time on page
	return nil
}

func (s *service) TrackProductView(ctx context.Context, tenantID uuid.UUID, productView *ProductView) error {
	// Validate product view data
	if productView == nil {
		return fmt.Errorf("product view cannot be nil")
	}
	if productView.ProductID == uuid.Nil {
		return fmt.Errorf("product ID is required")
	}
	if productView.TenantID == uuid.Nil {
		productView.TenantID = tenantID
	}
	if productView.Timestamp.IsZero() {
		productView.Timestamp = time.Now()
	}

	// Store product view
	err := s.repo.CreateProductView(ctx, productView)
	if err != nil {
		return fmt.Errorf("failed to create product view: %w", err)
	}

	// TODO: Validate product exists and update popularity metrics
	return nil
}

func (s *service) TrackPurchase(ctx context.Context, tenantID uuid.UUID, purchase *Purchase) error {
	// Validate purchase data
	if purchase == nil {
		return fmt.Errorf("purchase cannot be nil")
	}
	if purchase.TotalAmount <= 0 {
		return fmt.Errorf("purchase total amount must be greater than 0")
	}
	if purchase.TenantID == uuid.Nil {
		purchase.TenantID = tenantID
	}
	if purchase.Timestamp.IsZero() {
		purchase.Timestamp = time.Now()
	}

	// Store purchase event
	err := s.repo.CreatePurchase(ctx, purchase)
	if err != nil {
		return fmt.Errorf("failed to create purchase: %w", err)
	}

	// TODO: Update conversion metrics in real-time
	return nil
}

// Basic stats implementations
func (s *service) GetDashboardStats(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*AnalyticsStats, error) {
	// Get traffic metrics
	trafficStats, err := s.GetTrafficStats(ctx, tenantID, dateRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get traffic stats: %w", err)
	}

	// Get sales metrics
	salesStats, err := s.GetSalesStats(ctx, tenantID, dateRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales stats: %w", err)
	}

	// Get top performers
	topProducts, err := s.GetTopProducts(ctx, tenantID, dateRange, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top products: %w", err)
	}

	topPages, err := s.GetTopPages(ctx, tenantID, dateRange, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top pages: %w", err)
	}

	// Aggregate dashboard stats
	stats := &AnalyticsStats{
		TenantID:        tenantID,
		Date:            dateRange.Start,
		PageViews:       trafficStats.PageViews,
		UniqueVisitors:  trafficStats.UniqueVisitors,
		Sessions:        trafficStats.Sessions,
		BounceRate:      trafficStats.BounceRate,
		ConversionRate:  salesStats.ConversionRate,
		Revenue:         salesStats.TotalRevenue,
		Orders:          salesStats.TotalOrders,
		AvgOrderValue:   salesStats.AvgOrderValue,
		ProductViews:    0, // TODO: Get from traffic stats
		AvgSessionTime:  0, // TODO: Calculate from session data
		TopProducts:     make([]string, len(topProducts)),
		TopPages:        make([]string, len(topPages)),
		TopReferrers:    []string{}, // TODO: Get from traffic stats
	}

	// Convert ProductStats to string slice
	for i, product := range topProducts {
		stats.TopProducts[i] = product.ProductName // Assuming ProductStats has ProductName field
	}

	// Convert PageStats to string slice
	for i, page := range topPages {
		stats.TopPages[i] = page.Path // Assuming PageStats has Path field
	}

	return stats, nil
}

func (s *service) GetTrafficStats(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*TrafficStats, error) {
	// Get basic traffic stats from repository
	stats, err := s.repo.GetTrafficStats(ctx, tenantID, dateRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get traffic stats from repository: %w", err)
	}

	// Get top pages and referrers
	topPages, err := s.GetTopPages(ctx, tenantID, dateRange, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top pages: %w", err)
	}

	topReferrers, err := s.GetTopReferrers(ctx, tenantID, dateRange, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top referrers: %w", err)
	}

	// Convert to required format
	stats.TopPages = make([]PageStats, len(topPages))
	for i, page := range topPages {
		stats.TopPages[i] = *page
	}

	stats.TopReferrers = make([]ReferrerStats, len(topReferrers))
	for i, referrer := range topReferrers {
		stats.TopReferrers[i] = *referrer
	}

	return stats, nil
}

func (s *service) GetSalesStats(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*SalesStats, error) {
	// Get basic sales stats from repository
	stats, err := s.repo.GetSalesStats(ctx, tenantID, dateRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales stats from repository: %w", err)
	}

	// Get top products for sales analysis
	topProducts, err := s.GetTopProducts(ctx, tenantID, dateRange, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top products: %w", err)
	}

	// Convert to required format
	stats.TopProducts = make([]ProductStats, len(topProducts))
	for i, product := range topProducts {
		stats.TopProducts[i] = *product
	}

	// Calculate conversion rate if we have traffic data
	trafficStats, err := s.repo.GetTrafficStats(ctx, tenantID, dateRange)
	if err == nil && trafficStats.UniqueVisitors > 0 {
		stats.ConversionRate = float64(stats.TotalOrders) / float64(trafficStats.UniqueVisitors) * 100
	}

	return stats, nil
}

// Top performers implementations
func (s *service) GetTopProducts(ctx context.Context, tenantID uuid.UUID, dateRange DateRange, limit int) ([]*ProductStats, error) {
	// TODO: Implement top products calculation
	return s.repo.GetTopProducts(ctx, tenantID, dateRange, limit)
}

func (s *service) GetTopPages(ctx context.Context, tenantID uuid.UUID, dateRange DateRange, limit int) ([]*PageStats, error) {
	// TODO: Implement top pages calculation
	return s.repo.GetTopPages(ctx, tenantID, dateRange, limit)
}

func (s *service) GetTopReferrers(ctx context.Context, tenantID uuid.UUID, dateRange DateRange, limit int) ([]*ReferrerStats, error) {
	// TODO: Implement top referrers calculation
	return s.repo.GetTopReferrers(ctx, tenantID, dateRange, limit)
}

// Advanced analytics implementations
func (s *service) GetCohortAnalysis(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (*CohortAnalysis, error) {
	// TODO: Implement cohort analysis
	// - Group users by signup month
	// - Calculate retention rates for each cohort
	// - Return cohort data with retention percentages
	return s.repo.GetCohortAnalysis(ctx, tenantID, dateRange)
}

func (s *service) GetFunnelAnalysis(ctx context.Context, tenantID uuid.UUID, funnelSteps []string, dateRange DateRange) (*FunnelAnalysis, error) {
	// TODO: Implement funnel analysis
	// - Track user progression through defined steps
	// - Calculate conversion rates between steps
	// - Identify drop-off points
	return s.repo.GetFunnelAnalysis(ctx, tenantID, funnelSteps, dateRange)
}

func (s *service) GetCustomerLifetimeValue(ctx context.Context, tenantID uuid.UUID, dateRange DateRange) (float64, error) {
	// TODO: Implement CLV calculation
	// - Calculate average order value
	// - Calculate purchase frequency
	// - Calculate customer lifespan
	// - Return CLV = AOV * Frequency * Lifespan
	return s.repo.GetCustomerLifetimeValue(ctx, tenantID, dateRange)
}

func (s *service) GetRetentionRate(ctx context.Context, tenantID uuid.UUID, days int) (float64, error) {
	// TODO: Implement retention rate calculation
	return s.repo.GetRetentionRate(ctx, tenantID, days)
}

// Real-time analytics implementations
func (s *service) GetRealTimeStats(ctx context.Context, tenantID uuid.UUID) (*RealTimeStats, error) {
	// Get active users (last 5 minutes)
	activeUsers, err := s.GetActiveUsers(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	// Get real-time stats from repository
	stats, err := s.repo.GetRealTimeStats(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get real-time stats: %w", err)
	}

	// Update with calculated active users
	stats.ActiveUsers = activeUsers
	stats.LastUpdated = time.Now()

	return stats, nil
}

func (s *service) GetActiveUsers(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	// Get users active in the last 5 minutes
	// TODO: Implement GetActiveUsers in repository
	// For now, return a placeholder value
	return 0, nil
}

// Reports implementations
func (s *service) GenerateReport(ctx context.Context, tenantID uuid.UUID, request ReportRequest) ([]byte, string, error) {
	// TODO: Implement comprehensive report generation
	switch request.ReportType {
	case "traffic":
		return s.generateTrafficReport(ctx, tenantID, request)
	case "sales":
		return s.generateSalesReport(ctx, tenantID, request)
	case "products":
		return s.generateProductReport(ctx, tenantID, request)
	case "cohort":
		return s.generateCohortReport(ctx, tenantID, request)
	default:
		return nil, "", fmt.Errorf("unsupported report type: %s", request.ReportType)
	}
}

func (s *service) generateTrafficReport(ctx context.Context, tenantID uuid.UUID, request ReportRequest) ([]byte, string, error) {
	// TODO: Generate traffic report in requested format
	stats, err := s.GetTrafficStats(ctx, tenantID, request.DateRange)
	if err != nil {
		return nil, "", err
	}

	switch request.Format {
	case ReportFormatJSON:
		data, err := json.Marshal(stats)
		return data, "application/json", err
	case ReportFormatCSV:
		return s.trafficStatsToCSV(stats)
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", request.Format)
	}
}

func (s *service) generateSalesReport(ctx context.Context, tenantID uuid.UUID, request ReportRequest) ([]byte, string, error) {
	// TODO: Generate sales report in requested format
	stats, err := s.GetSalesStats(ctx, tenantID, request.DateRange)
	if err != nil {
		return nil, "", err
	}

	switch request.Format {
	case ReportFormatJSON:
		data, err := json.Marshal(stats)
		return data, "application/json", err
	case ReportFormatCSV:
		return s.salesStatsToCSV(stats)
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", request.Format)
	}
}

func (s *service) generateProductReport(ctx context.Context, tenantID uuid.UUID, request ReportRequest) ([]byte, string, error) {
	// TODO: Generate product performance report
	products, err := s.GetTopProducts(ctx, tenantID, request.DateRange, 100)
	if err != nil {
		return nil, "", err
	}

	switch request.Format {
	case ReportFormatJSON:
		data, err := json.Marshal(products)
		return data, "application/json", err
	case ReportFormatCSV:
		return s.productStatsToCSV(products)
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", request.Format)
	}
}

func (s *service) generateCohortReport(ctx context.Context, tenantID uuid.UUID, request ReportRequest) ([]byte, string, error) {
	// TODO: Generate cohort analysis report
	cohorts, err := s.GetCohortAnalysis(ctx, tenantID, request.DateRange)
	if err != nil {
		return nil, "", err
	}

	switch request.Format {
	case ReportFormatJSON:
		data, err := json.Marshal(cohorts)
		return data, "application/json", err
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", request.Format)
	}
}

// CSV conversion helpers
func (s *service) trafficStatsToCSV(stats *TrafficStats) ([]byte, string, error) {
	var result strings.Builder
	w := csv.NewWriter(&result)
	
	// Write headers
	headers := []string{"Metric", "Value"}
	if err := w.Write(headers); err != nil {
		return nil, "", fmt.Errorf("failed to write CSV headers: %w", err)
	}
	
	// Write data
	data := [][]string{
		{"Page Views", fmt.Sprintf("%d", stats.PageViews)},
		{"Unique Visitors", fmt.Sprintf("%d", stats.UniqueVisitors)},
		{"Sessions", fmt.Sprintf("%d", stats.Sessions)},
		{"Bounce Rate", fmt.Sprintf("%.2f%%", stats.BounceRate*100)},
		{"Avg Session Time", fmt.Sprintf("%.2f minutes", stats.AvgSessionTime/60)},
		{"New Visitors", fmt.Sprintf("%d", stats.NewVisitors)},
		{"Returning Visitors", fmt.Sprintf("%d", stats.ReturningVisitors)},
	}
	
	for _, row := range data {
		if err := w.Write(row); err != nil {
			return nil, "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	
	// Write top pages section
	if len(stats.TopPages) > 0 {
		w.Write([]string{"", ""}) // Empty row
		w.Write([]string{"Top Pages", ""})
		w.Write([]string{"Page", "Views", "Unique Views", "Bounce Rate", "Avg Time"})
		for _, page := range stats.TopPages {
			w.Write([]string{
				page.Path,
				fmt.Sprintf("%d", page.Views),
				fmt.Sprintf("%d", page.UniqueViews),
				fmt.Sprintf("%.2f%%", page.BounceRate*100),
				fmt.Sprintf("%.2f minutes", page.AvgTime/60),
			})
		}
	}
	
	// Write top referrers section
	if len(stats.TopReferrers) > 0 {
		w.Write([]string{"", ""}) // Empty row
		w.Write([]string{"Top Referrers", ""})
		w.Write([]string{"Referrer", "Visits", "Percentage"})
		for _, referrer := range stats.TopReferrers {
			w.Write([]string{
				referrer.Referrer,
				fmt.Sprintf("%d", referrer.Visits),
				fmt.Sprintf("%.2f%%", referrer.Percentage),
			})
		}
	}
	
	w.Flush()
	
	return []byte(result.String()), "text/csv", nil
}

func (s *service) salesStatsToCSV(stats *SalesStats) ([]byte, string, error) {
	var result strings.Builder
	w := csv.NewWriter(&result)
	
	headers := []string{"Metric", "Value"}
	if err := w.Write(headers); err != nil {
		return nil, "", fmt.Errorf("failed to write CSV headers: %w", err)
	}
	
	data := [][]string{
		{"Total Revenue", fmt.Sprintf("$%.2f", stats.TotalRevenue)},
		{"Total Orders", fmt.Sprintf("%d", stats.TotalOrders)},
		{"Average Order Value", fmt.Sprintf("$%.2f", stats.AvgOrderValue)},
		{"Conversion Rate", fmt.Sprintf("%.2f%%", stats.ConversionRate)},
		{"Refund Rate", fmt.Sprintf("%.2f%%", stats.RefundRate)},
	}
	
	for _, row := range data {
		if err := w.Write(row); err != nil {
			return nil, "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	
	// Write top products section
	if len(stats.TopProducts) > 0 {
		w.Write([]string{"", ""}) // Empty row
		w.Write([]string{"Top Products", ""})
		w.Write([]string{"Product Name", "Views", "Sales", "Revenue", "Conversion Rate"})
		for _, product := range stats.TopProducts {
			w.Write([]string{
				product.ProductName,
				fmt.Sprintf("%d", product.Views),
				fmt.Sprintf("%d", product.Sales),
				fmt.Sprintf("$%.2f", product.Revenue),
				fmt.Sprintf("%.2f%%", product.ConversionRate),
			})
		}
	}
	
	// Write sales by category section
	if len(stats.SalesByCategory) > 0 {
		w.Write([]string{"", ""}) // Empty row
		w.Write([]string{"Sales by Category", ""})
		w.Write([]string{"Category", "Revenue", "Orders"})
		for _, category := range stats.SalesByCategory {
			w.Write([]string{
				category.Category,
				fmt.Sprintf("$%.2f", category.Revenue),
				fmt.Sprintf("%d", category.Orders),
			})
		}
	}
	
	// Write daily sales section
	if len(stats.SalesByDay) > 0 {
		w.Write([]string{"", ""}) // Empty row
		w.Write([]string{"Daily Sales", ""})
		w.Write([]string{"Date", "Revenue", "Orders"})
		for _, daily := range stats.SalesByDay {
			w.Write([]string{
				daily.Date.Format("2006-01-02"),
				fmt.Sprintf("$%.2f", daily.Revenue),
				fmt.Sprintf("%d", daily.Orders),
			})
		}
	}
	
	w.Flush()
	
	return []byte(result.String()), "text/csv", nil
}

func (s *service) productStatsToCSV(products []*ProductStats) ([]byte, string, error) {
	var result strings.Builder
	w := csv.NewWriter(&result)
	
	headers := []string{"Product ID", "Product Name", "Views", "Sales", "Revenue", "Conversion Rate"}
	if err := w.Write(headers); err != nil {
		return nil, "", fmt.Errorf("failed to write CSV headers: %w", err)
	}
	
	for _, product := range products {
		row := []string{
			product.ProductID.String(),
			product.ProductName,
			fmt.Sprintf("%d", product.Views),
			fmt.Sprintf("%d", product.Sales),
			fmt.Sprintf("$%.2f", product.Revenue),
			fmt.Sprintf("%.2f%%", product.ConversionRate),
		}
		if err := w.Write(row); err != nil {
			return nil, "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	w.Flush()
	
	return []byte(result.String()), "text/csv", nil
}

// Scheduled reports
func (s *service) ScheduleReport(ctx context.Context, tenantID uuid.UUID, request ScheduleReportRequest) (*ScheduledReport, error) {
	// TODO: Implement scheduled report creation
	// - Validate request
	// - Calculate next run time
	// - Store scheduled report
	return s.repo.CreateScheduledReport(ctx, tenantID, request)
}

func (s *service) GetScheduledReports(ctx context.Context, tenantID uuid.UUID) ([]*ScheduledReport, error) {
	// TODO: Get all scheduled reports for tenant
	return s.repo.GetScheduledReports(ctx, tenantID)
}

func (s *service) DeleteScheduledReport(ctx context.Context, tenantID uuid.UUID, reportID uuid.UUID) error {
	// TODO: Delete scheduled report
	// - Validate report belongs to tenant
	// - Remove from repository
	return s.repo.DeleteScheduledReport(ctx, tenantID, reportID)
}

// Removed ExportData - GDPR data export functionality not needed
