package security

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for security operations
type Handler struct {
	service SecurityService
}

// NewHandler creates a new security handler
func NewHandler(service SecurityService) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers security routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// 🔐 PASSWORD MANAGEMENT ENDPOINTS
	passwords := router.Group("/passwords")
	{
		passwords.POST("/validate", h.ValidatePassword)          // ValidatePassword
		passwords.POST("/check-compromised", h.CheckCompromised) // IsPasswordCompromised
		passwords.GET("/policy", h.GetPasswordPolicy)            // Get password policy
	}

	// 🔐 LOGIN SECURITY ENDPOINTS
	logins := router.Group("/login-security")
	{
		logins.POST("/attempts", h.RecordLoginAttempt)   // RecordLoginAttempt
		logins.POST("/validate", h.ValidateLoginAttempt) // ValidateLoginAttempt
		logins.GET("/attempts", h.GetLoginAttempts)      // Get login attempts
	}

	// 🔐 ACCOUNT LOCKOUT ENDPOINTS
	lockouts := router.Group("/lockouts")
	{
		lockouts.GET("/status/:user_id", h.GetLockoutStatus) // CheckAccountLockout
		lockouts.POST("/lock", h.LockAccount)                // LockAccount
		lockouts.POST("/unlock", h.UnlockAccount)            // UnlockAccount
		lockouts.GET("", h.GetAccountLockouts)               // Get all lockouts
	}

	// 🔐 TRUSTED DEVICE ENDPOINTS
	devices := router.Group("/devices")
	{
		devices.POST("/register", h.RegisterTrustedDevice)   // RegisterTrustedDevice
		devices.POST("/validate", h.ValidateDevice)          // ValidateDevice
		devices.GET("/user/:user_id", h.GetUserDevices)      // GetUserDevices
		devices.DELETE("/:device_id", h.RevokeTrustedDevice) // RevokeTrustedDevice
	}

	// 🔐 SECURITY EVENT ENDPOINTS
	events := router.Group("/events")
	{
		events.POST("", h.LogSecurityEvent)                      // LogSecurityEvent
		events.GET("", h.GetSecurityEvents)                      // Get security events
		events.PUT("/:event_id/resolve", h.ResolveSecurityEvent) // ResolveSecurityEvent
	}

	// 🔐 THREAT DETECTION ENDPOINTS
	threats := router.Group("/threats")
	{
		threats.POST("/analyze", h.AnalyzeThreatLevel)      // AnalyzeThreatLevel
		threats.POST("/detect", h.DetectSuspiciousActivity) // DetectSuspiciousActivity
	}

	// 🔐 SECURITY ANALYTICS ENDPOINTS
	analytics := router.Group("/analytics")
	{
		analytics.GET("/dashboard", h.GetSecurityDashboard)   // GetSecurityDashboard
		analytics.GET("/report", h.GetSecurityReport)         // GetSecurityReport
		analytics.GET("/risk-score/:user_id", h.GetRiskScore) // GetRiskScore
	}

	// 🔐 AUDIT LOG ENDPOINTS
	audit := router.Group("/audit")
	{
		audit.GET("/logs", h.GetAuditLogs)                    // Get audit logs
		audit.POST("/logs", h.CreateAuditLog)                 // Create audit log
		audit.GET("/logs/:log_id", h.GetAuditLogDetails)      // Get audit log details
		audit.POST("/logs/export", h.ExportAuditLogs)         // Export audit logs
	}

	// 🔐 FRAUD DETECTION ENDPOINTS
	fraud := router.Group("/fraud")
	{
		fraud.POST("/detect", h.DetectFraud)                  // Detect fraud
		fraud.GET("/alerts", h.GetFraudAlerts)               // Get fraud alerts
		fraud.PUT("/alerts/:alert_id", h.UpdateFraudAlert)   // Update fraud alert
		fraud.GET("/patterns", h.GetFraudPatterns)           // Get fraud patterns
	}
}

// Password Management Handlers

// ValidatePassword validates a password against policy
func (h *Handler) ValidatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}

	var req struct {
		Password string    `json:"password" binding:"required"`
		UserID   uuid.UUID `json:"user_id" binding:"required"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	result, err := h.service.ValidatePassword(ctx, req.Password, req.UserID, &tenantID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// CheckCompromised checks if a password is compromised
func (h *Handler) CheckCompromised(c *gin.Context) {
	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	isCompromised, err := h.service.IsPasswordCompromised(req.Password)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_compromised": isCompromised})
}

// GetPasswordPolicy gets the password policy
func (h *Handler) GetPasswordPolicy(c *gin.Context) {
	// TODO: Implement get password policy
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// Login Security Handlers

// RecordLoginAttempt records a login attempt
func (h *Handler) RecordLoginAttempt(c *gin.Context) {
	ctx := c.Request.Context()

	var req LoginAttemptRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	attempt, err := h.service.RecordLoginAttempt(ctx, &req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, attempt)
}

// ValidateLoginAttempt validates a login attempt
func (h *Handler) ValidateLoginAttempt(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		UserID    uuid.UUID `json:"user_id" binding:"required"`
		Email     string    `json:"email" binding:"required"`
		IPAddress string    `json:"ip_address" binding:"required"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	result, err := h.service.ValidateLoginAttempt(ctx, req.UserID, req.Email, req.IPAddress)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetLoginAttempts gets login attempts
func (h *Handler) GetLoginAttempts(c *gin.Context) {
	// TODO: Implement get login attempts
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// Account Lockout Handlers

// GetLockoutStatus gets account lockout status
func (h *Handler) GetLockoutStatus(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := h.getUUIDParam(c, "user_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "details": err.Error()})
		return
	}

	status, err := h.service.CheckAccountLockout(ctx, userID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, status)
}

// LockAccount locks an account
func (h *Handler) LockAccount(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		UserID   uuid.UUID      `json:"user_id" binding:"required"`
		Reason   string         `json:"reason" binding:"required"`
		Duration *time.Duration `json:"duration,omitempty"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	if err := h.service.LockAccount(ctx, req.UserID, req.Reason, req.Duration); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// UnlockAccount unlocks an account
func (h *Handler) UnlockAccount(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		UserID  uuid.UUID  `json:"user_id" binding:"required"`
		AdminID *uuid.UUID `json:"admin_id,omitempty"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	if err := h.service.UnlockAccount(ctx, req.UserID, req.AdminID); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAccountLockouts gets account lockouts
func (h *Handler) GetAccountLockouts(c *gin.Context) {
	// TODO: Implement get account lockouts
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// Trusted Device Handlers

// RegisterTrustedDevice registers a trusted device
func (h *Handler) RegisterTrustedDevice(c *gin.Context) {
	ctx := c.Request.Context()

	var req TrustedDeviceRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	device, err := h.service.RegisterTrustedDevice(ctx, &req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, device)
}

// ValidateDevice validates a device
func (h *Handler) ValidateDevice(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		UserID      uuid.UUID `json:"user_id" binding:"required"`
		Fingerprint string    `json:"fingerprint" binding:"required"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	result, err := h.service.ValidateDevice(ctx, req.UserID, req.Fingerprint)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserDevices gets user devices
func (h *Handler) GetUserDevices(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := h.getUUIDParam(c, "user_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "details": err.Error()})
		return
	}

	devices, err := h.service.GetUserDevices(ctx, userID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, devices)
}

// RevokeTrustedDevice revokes a trusted device
func (h *Handler) RevokeTrustedDevice(c *gin.Context) {
	ctx := c.Request.Context()
	deviceID, err := h.getUUIDParam(c, "device_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID", "details": err.Error()})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	if err := h.service.RevokeTrustedDevice(ctx, deviceID, req.Reason); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Security Event Handlers

// LogSecurityEvent logs a security event
func (h *Handler) LogSecurityEvent(c *gin.Context) {
	ctx := c.Request.Context()

	var req SecurityEventRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	if err := h.service.LogSecurityEvent(ctx, &req); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// GetSecurityEvents gets security events
func (h *Handler) GetSecurityEvents(c *gin.Context) {
	// TODO: Implement get security events
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// ResolveSecurityEvent resolves a security event
func (h *Handler) ResolveSecurityEvent(c *gin.Context) {
	ctx := c.Request.Context()
	eventID, err := h.getUUIDParam(c, "event_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID", "details": err.Error()})
		return
	}

	var req struct {
		Resolution string     `json:"resolution" binding:"required"`
		AdminID    *uuid.UUID `json:"admin_id,omitempty"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	if err := h.service.ResolveSecurityEvent(ctx, eventID, req.Resolution, req.AdminID); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Threat Detection Handlers

// AnalyzeThreatLevel analyzes threat level
func (h *Handler) AnalyzeThreatLevel(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		UserID    uuid.UUID `json:"user_id" binding:"required"`
		IPAddress string    `json:"ip_address" binding:"required"`
		UserAgent string    `json:"user_agent" binding:"required"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	threatLevel, err := h.service.AnalyzeThreatLevel(ctx, req.UserID, req.IPAddress, req.UserAgent)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"threat_level": threatLevel})
}

// DetectSuspiciousActivity detects suspicious activity
func (h *Handler) DetectSuspiciousActivity(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		UserID   uuid.UUID        `json:"user_id" binding:"required"`
		Activity *ActivityContext `json:"activity" binding:"required"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	assessment, err := h.service.DetectSuspiciousActivity(ctx, req.UserID, req.Activity)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, assessment)
}

// Security Analytics Handlers

// GetSecurityDashboard gets security dashboard
func (h *Handler) GetSecurityDashboard(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}

	// Parse period parameter
	periodStr := c.DefaultQuery("period", "24h")
	period, err := time.ParseDuration(periodStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period format", "details": err.Error()})
		return
	}

	dashboard, err := h.service.GetSecurityDashboard(ctx, &tenantID, period)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// GetSecurityReport gets security report
func (h *Handler) GetSecurityReport(c *gin.Context) {
	// TODO: Implement get security report
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetRiskScore gets risk score
func (h *Handler) GetRiskScore(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := h.getUUIDParam(c, "user_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "details": err.Error()})
		return
	}

	riskScore, err := h.service.GetRiskScore(ctx, userID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, riskScore)
}

// Helper methods

// getTenantID extracts tenant ID from context or headers
func (h *Handler) getTenantID(c *gin.Context) (uuid.UUID, error) {
	// Try to get from context first (set by middleware)
	if tenantID, exists := c.Get("tenant_id"); exists {
		if id, ok := tenantID.(uuid.UUID); ok {
			return id, nil
		}
	}

	// Try to get from header
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		return uuid.Nil, nil // Return nil UUID for optional tenant ID
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return tenantID, nil
}

// getUUIDParam extracts UUID parameter from URL
func (h *Handler) getUUIDParam(c *gin.Context, param string) (uuid.UUID, error) {
	idStr := c.Param(param)
	return uuid.Parse(idStr)
}

// handleServiceError handles service layer errors
func (h *Handler) handleServiceError(c *gin.Context, err error) {
	// TODO: Implement proper error handling based on error types
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// Audit Log Handlers

// GetAuditLogs gets audit logs
func (h *Handler) GetAuditLogs(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}

	// Parse query parameters
	userID := c.Query("user_id")
	action := c.Query("action")
	resource := c.Query("resource")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "50")

	logs, total, err := h.service.GetAuditLogs(ctx, tenantID, userID, action, resource, startDate, endDate, page, limit)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// CreateAuditLog creates an audit log
func (h *Handler) CreateAuditLog(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		UserID     uuid.UUID   `json:"user_id" binding:"required"`
		Action     string      `json:"action" binding:"required"`
		Resource   string      `json:"resource" binding:"required"`
		ResourceID *uuid.UUID  `json:"resource_id,omitempty"`
		Details    interface{} `json:"details,omitempty"`
		IPAddress  string      `json:"ip_address,omitempty"`
		UserAgent  string      `json:"user_agent,omitempty"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	log, err := h.service.CreateAuditLog(ctx, &req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, log)
}

// GetAuditLogDetails gets audit log details
func (h *Handler) GetAuditLogDetails(c *gin.Context) {
	ctx := c.Request.Context()
	logID, err := h.getUUIDParam(c, "log_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID", "details": err.Error()})
		return
	}

	log, err := h.service.GetAuditLogDetails(ctx, logID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, log)
}

// ExportAuditLogs exports audit logs
func (h *Handler) ExportAuditLogs(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}

	var req struct {
		Format    string `json:"format" binding:"required"`
		StartDate string `json:"start_date,omitempty"`
		EndDate   string `json:"end_date,omitempty"`
		Filters   map[string]interface{} `json:"filters,omitempty"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	result, err := h.service.ExportAuditLogs(ctx, tenantID, req.Format, req.StartDate, req.EndDate, req.Filters)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// Fraud Detection Handlers

// DetectFraud detects fraud
func (h *Handler) DetectFraud(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		UserID      uuid.UUID   `json:"user_id" binding:"required"`
		Transaction interface{} `json:"transaction" binding:"required"`
		Context     interface{} `json:"context,omitempty"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	result, err := h.service.DetectFraud(ctx, req.UserID, req.Transaction, req.Context)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetFraudAlerts gets fraud alerts
func (h *Handler) GetFraudAlerts(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}

	status := c.Query("status")
	severity := c.Query("severity")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "50")

	alerts, total, err := h.service.GetFraudAlerts(ctx, tenantID, status, severity, page, limit)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// UpdateFraudAlert updates a fraud alert
func (h *Handler) UpdateFraudAlert(c *gin.Context) {
	ctx := c.Request.Context()
	alertID, err := h.getUUIDParam(c, "alert_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID", "details": err.Error()})
		return
	}

	var req struct {
		Status     string     `json:"status" binding:"required"`
		Resolution string     `json:"resolution,omitempty"`
		AdminID    *uuid.UUID `json:"admin_id,omitempty"`
	}

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": bindErr.Error()})
		return
	}

	alert, err := h.service.UpdateFraudAlert(ctx, alertID, req.Status, req.Resolution, req.AdminID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, alert)
}

// GetFraudPatterns gets fraud patterns
func (h *Handler) GetFraudPatterns(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}

	patternType := c.Query("type")
	period := c.DefaultQuery("period", "30d")

	patterns, err := h.service.GetFraudPatterns(ctx, tenantID, patternType, period)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, patterns)
}
