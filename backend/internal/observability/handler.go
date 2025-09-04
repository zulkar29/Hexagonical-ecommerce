package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ObservabilityHandler handles HTTP requests for observability endpoints
type ObservabilityHandler struct {
	service *ObservabilityService
}

// NewObservabilityHandler creates a new observability handler
func NewObservabilityHandler(service *ObservabilityService) *ObservabilityHandler {
	return &ObservabilityHandler{
		service: service,
	}
}

// RegisterRoutes registers observability routes
func (h *ObservabilityHandler) RegisterRoutes(router *gin.RouterGroup) {
	obs := router.Group("/observability")
	{
		// Health endpoints
		obs.GET("/health", h.GetHealth)
		obs.GET("/health/detailed", h.GetDetailedHealth)

		// Metrics endpoints
		obs.GET("/metrics", h.GetMetrics)
		obs.GET("/metrics/summary", h.GetMetricsSummary)

		// Logging endpoints
		obs.GET("/logs", h.GetLogs)
		obs.POST("/logs", h.CreateLog)

		// Tracing endpoints
		obs.GET("/traces", h.GetTraces)
		obs.GET("/traces/:traceId", h.GetTrace)

		// Alert endpoints
		obs.GET("/alerts", h.GetAlerts)
		obs.POST("/alerts", h.CreateAlert)
		obs.PUT("/alerts/:alertId/resolve", h.ResolveAlert)

		// System information
		obs.GET("/system/info", h.GetSystemInfo)
		obs.GET("/system/stats", h.GetSystemStats)
	}
}

// Health endpoints

// GetHealth returns basic health status
func (h *ObservabilityHandler) GetHealth(c *gin.Context) {
	status, err := h.service.GetHealthStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get health status",
			"details": err.Error(),
		})
		return
	}

	statusCode := http.StatusOK
	if status.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":    status.Status,
		"timestamp": status.Timestamp,
		"version":   status.Version,
	})
}

// GetDetailedHealth returns detailed health status including all services
func (h *ObservabilityHandler) GetDetailedHealth(c *gin.Context) {
	status, err := h.service.GetHealthStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get health status",
			"details": err.Error(),
		})
		return
	}

	statusCode := http.StatusOK
	if status.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, status)
}

// Metrics endpoints

// GetMetrics returns collected metrics
func (h *ObservabilityHandler) GetMetrics(c *gin.Context) {
	// Parse query parameters for filtering
	filters := make(map[string]string)
	if service := c.Query("service"); service != "" {
		filters["service"] = service
	}
	if metricType := c.Query("type"); metricType != "" {
		filters["type"] = metricType
	}

	metrics, err := h.service.GetMetrics(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get metrics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
		"count":   len(metrics),
	})
}

// GetMetricsSummary returns a summary of metrics
func (h *ObservabilityHandler) GetMetricsSummary(c *gin.Context) {
	metrics, err := h.service.GetMetrics(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get metrics",
			"details": err.Error(),
		})
		return
	}

	// Calculate summary statistics
	summary := map[string]interface{}{
		"total_metrics": len(metrics),
		"by_type":       make(map[string]int),
		"by_service":    make(map[string]int),
	}

	byType := summary["by_type"].(map[string]int)
	byService := summary["by_service"].(map[string]int)

	for _, metric := range metrics {
		byType[string(metric.Type)]++
		byService[metric.Service]++
	}

	c.JSON(http.StatusOK, summary)
}

// Logging endpoints

// GetLogs returns log entries
func (h *ObservabilityHandler) GetLogs(c *gin.Context) {
	// Parse query parameters
	limit := 100
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Build filters
	filters := make(map[string]interface{})
	if level := c.Query("level"); level != "" {
		filters["level"] = level
	}
	if service := c.Query("service"); service != "" {
		filters["service"] = service
	}
	if tenantID := c.Query("tenant_id"); tenantID != "" {
		filters["tenant_id"] = tenantID
	}
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}
	if traceID := c.Query("trace_id"); traceID != "" {
		filters["trace_id"] = traceID
	}
	if from := c.Query("from"); from != "" {
		if timestamp, err := time.Parse(time.RFC3339, from); err == nil {
			filters["from"] = timestamp
		}
	}
	if to := c.Query("to"); to != "" {
		if timestamp, err := time.Parse(time.RFC3339, to); err == nil {
			filters["to"] = timestamp
		}
	}

	// This would use the repository to get logs from database
	// For now, return empty array
	c.JSON(http.StatusOK, gin.H{
		"logs":   []interface{}{},
		"limit":  limit,
		"offset": offset,
		"total":  0,
	})
}

// CreateLog creates a new log entry
func (h *ObservabilityHandler) CreateLog(c *gin.Context) {
	var request struct {
		Level   string                 `json:"level" binding:"required"`
		Message string                 `json:"message" binding:"required"`
		Service string                 `json:"service"`
		Fields  map[string]interface{} `json:"fields"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Parse log level
	level := LogLevel(request.Level)
	if !isValidLogLevel(level) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid log level",
		})
		return
	}

	// Log the message
	h.service.LogEvent(c.Request.Context(), level, request.Message, request.Fields, nil)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Log entry created",
	})
}

// Tracing endpoints

// GetTraces returns traces
func (h *ObservabilityHandler) GetTraces(c *gin.Context) {
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	traces, err := h.service.GetTraces(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get traces",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"traces": traces,
		"count":  len(traces),
	})
}

// GetTrace returns a specific trace
func (h *ObservabilityHandler) GetTrace(c *gin.Context) {
	traceID := c.Param("traceId")
	if traceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Trace ID is required",
		})
		return
	}

	trace := h.service.tracer.GetTrace(traceID)
	if trace == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Trace not found",
		})
		return
	}

	c.JSON(http.StatusOK, trace)
}

// Alert endpoints

// GetAlerts returns alerts
func (h *ObservabilityHandler) GetAlerts(c *gin.Context) {
	alerts, err := h.service.GetAlerts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get alerts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"count":  len(alerts),
	})
}

// CreateAlert creates a new alert
func (h *ObservabilityHandler) CreateAlert(c *gin.Context) {
	var request struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Severity    string `json:"severity" binding:"required"`
		Service     string `json:"service" binding:"required"`
		Condition   string `json:"condition" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Validate severity
	severity := AlertSeverity(request.Severity)
	if !isValidAlertSeverity(severity) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid alert severity",
		})
		return
	}

	alert := &Alert{
		Title:       request.Title,
		Description: request.Description,
		Severity:    severity,
		Service:     request.Service,
		Condition:   request.Condition,
		Status:      AlertStatusActive,
	}

	err := h.service.CreateAlert(c.Request.Context(), alert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create alert",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, alert)
}

// ResolveAlert resolves an alert
func (h *ObservabilityHandler) ResolveAlert(c *gin.Context) {
	alertIDStr := c.Param("alertId")
	alertID, err := strconv.ParseInt(alertIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid alert ID",
		})
		return
	}

	// This would resolve the alert in the database
	// For now, just return success
	c.JSON(http.StatusOK, gin.H{
		"message":  "Alert resolved",
		"alert_id": alertID,
	})
}

// System information endpoints

// GetSystemInfo returns system information
func (h *ObservabilityHandler) GetSystemInfo(c *gin.Context) {
	info := map[string]interface{}{
		"service":     "esass-backend",
		"version":     "1.0.0",
		"environment": "development", // This should come from config
		"go_version":  "1.21",
		"uptime":      time.Since(time.Now().Add(-1*time.Hour)).String(), // Mock uptime
		"build_time":  "2024-01-01T00:00:00Z",                             // This should come from build
	}

	c.JSON(http.StatusOK, info)
}

// GetSystemStats returns system statistics
func (h *ObservabilityHandler) GetSystemStats(c *gin.Context) {
	stats := map[string]interface{}{
		"requests_total":      1000,    // Mock data
		"requests_per_second": 10.5,    // Mock data
		"response_time_avg":   "125ms",  // Mock data
		"error_rate":          "0.5%",   // Mock data
		"memory_usage":        "256MB",  // Mock data
		"cpu_usage":           "15%",    // Mock data
		"goroutines":          25,       // Mock data
	}

	c.JSON(http.StatusOK, stats)
}

// Helper functions

func isValidLogLevel(level LogLevel) bool {
	switch level {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError, LogLevelFatal:
		return true
	default:
		return false
	}
}

func isValidAlertSeverity(severity AlertSeverity) bool {
	switch severity {
	case AlertSeverityLow, AlertSeverityMedium, AlertSeverityHigh, AlertSeverityCritical:
		return true
	default:
		return false
	}
}
