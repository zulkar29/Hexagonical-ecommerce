package analytics

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ecommerce-saas/internal/shared/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	// Event tracking routes
	tracking := r.Group("/track")
	{
		tracking.POST("/event", h.TrackEvent)
		tracking.POST("/page-view", h.TrackPageView)
		tracking.POST("/product-view", h.TrackProductView)
		tracking.POST("/purchase", h.TrackPurchase)
	}

	// Analytics dashboard routes
	r.GET("/dashboard", h.GetDashboardStats)
	r.GET("/traffic", h.GetTrafficStats)
	r.GET("/sales", h.GetSalesStats)
	r.GET("/realtime", h.GetRealTimeStats)

	// Top performers routes
	top := r.Group("/top")
	{
		top.GET("/products", h.GetTopProducts)
		top.GET("/pages", h.GetTopPages)
		top.GET("/referrers", h.GetTopReferrers)
	}

	// Advanced analytics routes
	advanced := r.Group("/advanced")
	{
		advanced.GET("/cohorts", h.GetCohortAnalysis)
		advanced.GET("/funnel", h.GetFunnelAnalysis)
		advanced.GET("/clv", h.GetCustomerLifetimeValue)
		advanced.GET("/retention", h.GetRetentionRate)
	}

	// Reports routes
	reports := r.Group("/reports")
	{
		reports.POST("/generate", h.GenerateReport)
		reports.POST("/schedule", h.ScheduleReport)
		reports.GET("/scheduled", h.GetScheduledReports)
		reports.DELETE("/scheduled/:id", h.DeleteScheduledReport)
	}

	// Removed data export routes - GDPR functionality not needed
}

// Event tracking handlers
func (h *Handler) TrackEvent(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var event AnalyticsEvent
	if bindErr := c.ShouldBindJSON(&event); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	event.TenantID = tenantID
	event.Timestamp = time.Now()

	if err := h.service.TrackEvent(c.Request.Context(), tenantID, &event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event tracked successfully"})
}

func (h *Handler) TrackPageView(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var pageView PageView
	if bindErr := c.ShouldBindJSON(&pageView); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	pageView.TenantID = tenantID
	pageView.Timestamp = time.Now()

	if err := h.service.TrackPageView(c.Request.Context(), tenantID, &pageView); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Page view tracked successfully"})
}

func (h *Handler) TrackProductView(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var productView ProductView
	if bindErr := c.ShouldBindJSON(&productView); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	productView.TenantID = tenantID
	productView.Timestamp = time.Now()

	if err := h.service.TrackProductView(c.Request.Context(), tenantID, &productView); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product view tracked successfully"})
}

func (h *Handler) TrackPurchase(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var purchase Purchase
	if bindErr := c.ShouldBindJSON(&purchase); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	purchase.TenantID = tenantID
	purchase.Timestamp = time.Now()

	if err := h.service.TrackPurchase(c.Request.Context(), tenantID, &purchase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Purchase tracked successfully"})
}

// Analytics dashboard handlers
func (h *Handler) GetDashboardStats(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	dateRange := h.parseDateRange(c)
	stats, err := h.service.GetDashboardStats(c.Request.Context(), tenantID, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetTrafficStats(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	dateRange := h.parseDateRange(c)
	stats, err := h.service.GetTrafficStats(c.Request.Context(), tenantID, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetSalesStats(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	dateRange := h.parseDateRange(c)
	stats, err := h.service.GetSalesStats(c.Request.Context(), tenantID, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetRealTimeStats(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	stats, err := h.service.GetRealTimeStats(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Top performers handlers
func (h *Handler) GetTopProducts(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	dateRange := h.parseDateRange(c)
	limit := h.parseLimit(c, 10)

	products, err := h.service.GetTopProducts(c.Request.Context(), tenantID, dateRange, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *Handler) GetTopPages(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	dateRange := h.parseDateRange(c)
	limit := h.parseLimit(c, 10)

	pages, err := h.service.GetTopPages(c.Request.Context(), tenantID, dateRange, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pages})
}

func (h *Handler) GetTopReferrers(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	dateRange := h.parseDateRange(c)
	limit := h.parseLimit(c, 10)

	referrers, err := h.service.GetTopReferrers(c.Request.Context(), tenantID, dateRange, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": referrers})
}

// Advanced analytics handlers
func (h *Handler) GetCohortAnalysis(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	dateRange := h.parseDateRange(c)
	cohorts, err := h.service.GetCohortAnalysis(c.Request.Context(), tenantID, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cohorts)
}

func (h *Handler) GetFunnelAnalysis(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	dateRange := h.parseDateRange(c)
	
	// Parse funnel steps from query params or request body
	funnelSteps := []string{"page_view", "add_to_cart", "checkout", "purchase"} // Default funnel
	if steps := c.Query("steps"); steps != "" {
		// Parse comma-separated funnel steps
		parsedSteps := strings.Split(steps, ",")
		for i, step := range parsedSteps {
			parsedSteps[i] = strings.TrimSpace(step)
		}
		if len(parsedSteps) > 0 {
			funnelSteps = parsedSteps
		}
	}

	funnel, err := h.service.GetFunnelAnalysis(c.Request.Context(), tenantID, funnelSteps, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, funnel)
}

func (h *Handler) GetCustomerLifetimeValue(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	dateRange := h.parseDateRange(c)
	clv, err := h.service.GetCustomerLifetimeValue(c.Request.Context(), tenantID, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"customer_lifetime_value": clv})
}

func (h *Handler) GetRetentionRate(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	days := 30 // Default retention period
	if d := c.Query("days"); d != "" {
		if parsed, parseErr := strconv.Atoi(d); parseErr == nil {
			days = parsed
		}
	}

	retentionRate, err := h.service.GetRetentionRate(c.Request.Context(), tenantID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"retention_rate": retentionRate})
}

// Reports handlers
func (h *Handler) GenerateReport(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var request ReportRequest
	if bindErr := c.ShouldBindJSON(&request); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	data, contentType, err := h.service.GenerateReport(c.Request.Context(), tenantID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set appropriate headers for file download
	filename := fmt.Sprintf("%s_report_%d.%s", request.ReportType, time.Now().Unix(), h.getFileExtension(request.Format))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, contentType, data)
}

func (h *Handler) ScheduleReport(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var request ScheduleReportRequest
	if bindErr := c.ShouldBindJSON(&request); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	report, err := h.service.ScheduleReport(c.Request.Context(), tenantID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

func (h *Handler) GetScheduledReports(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	reports, err := h.service.GetScheduledReports(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reports})
}

func (h *Handler) DeleteScheduledReport(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	reportID, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	err = h.service.DeleteScheduledReport(c.Request.Context(), tenantID, reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Removed ExportData handler - GDPR data export functionality not needed

// Helper methods
func (h *Handler) parseDateRange(c *gin.Context) DateRange {
	// Default to last 30 days
	end := time.Now()
	start := end.AddDate(0, 0, -30)

	if startStr := c.Query("start"); startStr != "" {
		if parsed, parseErr := time.Parse(time.RFC3339, startStr); parseErr == nil {
			start = parsed
		}
	}

	if endStr := c.Query("end"); endStr != "" {
		if parsed, parseErr := time.Parse(time.RFC3339, endStr); parseErr == nil {
			end = parsed
		}
	}

	return DateRange{
		Start: start,
		End:   end,
	}
}

func (h *Handler) parseLimit(c *gin.Context, defaultLimit int) int {
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, parseErr := strconv.Atoi(limitStr); parseErr == nil && limit > 0 {
			return limit
		}
	}
	return defaultLimit
}

func (h *Handler) getFileExtension(format ReportFormat) string {
	switch format {
	case ReportFormatJSON:
		return "json"
	case ReportFormatCSV:
		return "csv"
	case ReportFormatPDF:
		return "pdf"
	default:
		return "json"
	}
}
