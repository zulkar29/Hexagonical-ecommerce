package security

import (
	"net/http"
	"strconv"
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
		passwords.POST("/validate", h.ValidatePassword)        // ValidatePassword
		passwords.POST("/check-compromised", h.CheckCompromised) // IsPasswordCompromised
		passwords.GET("/policy", h.GetPasswordPolicy)          // Get password policy
	}

	// 🔐 LOGIN SECURITY ENDPOINTS
	logins := router.Group("/login-security")
	{
		logins.POST("/attempts", h.RecordLoginAttempt)     // RecordLoginAttempt
		logins.POST("/validate", h.ValidateLoginAttempt)   // ValidateLoginAttempt
		logins.GET("/attempts", h.GetLoginAttempts)        // Get login attempts
	}

	// 🔐 ACCOUNT LOCKOUT ENDPOINTS
	lockouts := router.Group("/lockouts")
	{
		lockouts.GET("/status/:user_id", h.GetLockoutStatus)  // CheckAccountLockout
		lockouts.POST("/lock", h.LockAccount)                 // LockAccount
		lockouts.POST("/unlock", h.UnlockAccount)             // UnlockAccount
		lockouts.GET("", h.GetAccountLockouts)                // Get all lockouts
	}

	// 🔐 TRUSTED DEVICE ENDPOINTS
	devices := router.Group("/devices")
	{
		devices.POST("/register", h.RegisterTrustedDevice)    // RegisterTrustedDevice
		devices.POST("/validate", h.ValidateDevice)           // ValidateDevice
		devices.GET("/user/:user_id", h.GetUserDevices)       // GetUserDevices
		devices.DELETE("/:device_id", h.RevokeTrustedDevice)  // RevokeTrustedDevice
	}

	// 🔐 SECURITY EVENT ENDPOINTS
	events := router.Group("/events")
	{
		events.POST("", h.LogSecurityEvent)                   // LogSecurityEvent
		events.GET("", h.GetSecurityEvents)                   // Get security events
		events.PUT("/:event_id/resolve", h.ResolveSecurityEvent) // ResolveSecurityEvent
	}

	// 🔐 THREAT DETECTION ENDPOINTS
	threats := router.Group("/threats")
	{
		threats.POST("/analyze", h.AnalyzeThreatLevel)        // AnalyzeThreatLevel
		threats.POST("/detect", h.DetectSuspiciousActivity)   // DetectSuspiciousActivity
	}

	// 🔐 SECURITY ANALYTICS ENDPOINTS
	analytics := router.Group("/analytics")
	{
		analytics.GET("/dashboard", h.GetSecurityDashboard)   // GetSecurityDashboard
		analytics.GET("/report", h.GetSecurityReport)         // GetSecurityReport
		analytics.GET("/risk-score/:user_id", h.GetRiskScore) // GetRiskScore
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	result, err := h.service.ValidatePassword(ctx, req.Password, req.UserID, tenantID)
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
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

	dashboard, err := h.service.GetSecurityDashboard(ctx, tenantID, period)
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
func (h *Handler) getTenantID(c *gin.Context) (*uuid.UUID, error) {
	// Try to get from context first (set by middleware)
	if tenantID, exists := c.Get("tenant_id"); exists {
		if id, ok := tenantID.(uuid.UUID); ok {
			return &id, nil
		}
	}

	// Try to get from header
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		return nil, nil // Optional tenant ID
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, err
	}

	return &tenantID, nil
}

// getUUIDParam extracts UUID parameter from URL
func (h *Handler) getUUIDParam(c *gin.Context, param string) (uuid.UUID, error) {
	idStr := c.Param(param)
	return uuid.Parse(idStr)
}

// getIntQuery extracts integer query parameter
func (h *Handler) getIntQuery(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// handleServiceError handles service layer errors
func (h *Handler) handleServiceError(c *gin.Context, err error) {
	// TODO: Implement proper error handling based on error types
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}