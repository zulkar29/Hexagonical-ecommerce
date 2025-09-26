package security

import (
	"fmt"
	"net/http"
	"strconv"

	"ecommerce-saas/internal/shared/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler provides HTTP handlers for security operations
type Handler struct {
	service *SecurityService
}

// NewHandler creates a new security handler
func NewHandler(service *SecurityService) *Handler {
	return &Handler{service: service}
}

// PasswordValidationRequest represents a password validation request
type PasswordValidationRequest struct {
	Password string    `json:"password" binding:"required"`
	UserID   uuid.UUID `json:"user_id,omitempty"`
}

// SetupRoutes sets up the security routes
func (h *Handler) SetupRoutes(rg *gin.RouterGroup) {
	// Password validation
	rg.POST("/validate-password", h.ValidatePassword)

	// Account lockout management
	rg.GET("/account-lockout/:user_id", h.CheckAccountLockout)
	rg.POST("/failed-login", h.RecordFailedLogin)
	rg.POST("/successful-login", h.RecordSuccessfulLogin)
	rg.DELETE("/failed-logins/:user_id", h.ResetFailedLogins)

	// Two-Factor Authentication
	rg.POST("/2fa/setup/:user_id", h.Setup2FA)
	rg.POST("/2fa/verify", h.Verify2FA)
	rg.POST("/2fa/enable", h.Enable2FA)
	rg.DELETE("/2fa/:user_id", h.Disable2FA)
	rg.POST("/2fa/backup-code", h.VerifyBackupCode)

	// Password history
	rg.POST("/password-history", h.SavePasswordHistory)
	rg.POST("/check-password-reuse", h.CheckPasswordReuse)

	// Security settings and logs
	rg.GET("/settings/:user_id", h.GetSecuritySettings)
	rg.PUT("/settings/:user_id", h.UpdateSecuritySettings)
	rg.GET("/logs/:user_id", h.GetSecurityLogs)
}

// ValidatePassword validates a password against security policies
func (h *Handler) ValidatePassword(c *gin.Context) {
	var req PasswordValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	// Extract tenant ID from context or headers
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	var tenantID *uuid.UUID
	if tenantIDStr != "" {
		if parsed, err := uuid.Parse(tenantIDStr); err == nil {
			tenantID = &parsed
		}
	}

	ctx := c.Request.Context()
	err := h.service.ValidatePassword(ctx, tenantID, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"valid": false, "error": err.Error()})
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"valid": true})
}

// GetSecuritySettings retrieves security settings for a user
func (h *Handler) GetSecuritySettings(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid user ID"))
		return
	}

	ctx := c.Request.Context()
	settings, err := h.service.GetSecuritySettings(ctx, userID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, settings)
}

// GetSecurityLogs retrieves security logs for a user
func (h *Handler) GetSecurityLogs(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid user ID"))
		return
	}

	// Get limit from query parameter, default to 50
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, limitError := strconv.Atoi(limitStr); limitError == nil {
			limit = parsedLimit
		}
	}

	ctx := c.Request.Context()
	logs, err := h.service.GetSecurityLogs(ctx, userID, limit)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, logs)
}

// CheckAccountLockout checks if an account is currently locked
func (h *Handler) CheckAccountLockout(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid user ID"))
		return
	}

	// Extract tenant ID from context or headers
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	var tenantID *uuid.UUID
	if tenantIDStr != "" {
		if parsed, tenantIDError := uuid.Parse(tenantIDStr); tenantIDError == nil {
			tenantID = &parsed
		}
	}

	ctx := c.Request.Context()
	isLocked, err := h.service.IsAccountLocked(ctx, userID, tenantID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"locked": isLocked})
}

// LoginAttemptRequest represents a login attempt request
type LoginAttemptRequest struct {
	UserID    uuid.UUID  `json:"user_id" binding:"required"`
	TenantID  *uuid.UUID `json:"tenant_id,omitempty"`
	IPAddress string     `json:"ip_address" binding:"required"`
	UserAgent string     `json:"user_agent" binding:"required"`
}

// RecordFailedLogin records a failed login attempt
func (h *Handler) RecordFailedLogin(c *gin.Context) {
	var req LoginAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.service.RecordFailedLogin(ctx, req.UserID, req.TenantID, req.IPAddress, req.UserAgent)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Failed login recorded"})
}

// RecordSuccessfulLogin records a successful login attempt
func (h *Handler) RecordSuccessfulLogin(c *gin.Context) {
	var req LoginAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.service.RecordSuccessfulLogin(ctx, req.UserID, req.TenantID, req.IPAddress, req.UserAgent)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Successful login recorded"})
}

// ResetFailedLogins resets failed login attempts for a user
func (h *Handler) ResetFailedLogins(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid user ID"))
		return
	}

	ctx := c.Request.Context()
	err = h.service.ResetFailedLogins(ctx, userID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Failed logins reset"})
}

// Setup2FARequest represents a 2FA setup request
type Setup2FARequest struct {
	Issuer      string `json:"issuer" binding:"required"`
	AccountName string `json:"account_name" binding:"required"`
}

// Setup2FA sets up 2FA for a user
func (h *Handler) Setup2FA(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid user ID"))
		return
	}

	var req Setup2FARequest
	if bindError := c.ShouldBindJSON(&req); bindError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindError.Error()})
		return
	}

	// Extract tenant ID from context or headers
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	var tenantID *uuid.UUID
	if tenantIDStr != "" {
		if parsed, tenantIDError := uuid.Parse(tenantIDStr); tenantIDError == nil {
			tenantID = &parsed
		}
	}

	ctx := c.Request.Context()
	twoFA, qrURL, err := h.service.Setup2FA(ctx, userID, tenantID, req.Issuer, req.AccountName)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{
		"two_factor_auth": twoFA,
		"qr_url":          qrURL,
	})
}

// Verify2FARequest represents a 2FA verification request
type Verify2FARequest struct {
	UserID   uuid.UUID  `json:"user_id" binding:"required"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
	Token    string     `json:"token" binding:"required"`
}

// Verify2FA verifies a 2FA token
func (h *Handler) Verify2FA(c *gin.Context) {
	var req Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	valid, err := h.service.Verify2FA(ctx, req.UserID, req.TenantID, req.Token)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"valid": valid})
}

// Enable2FA enables 2FA for a user after verification
func (h *Handler) Enable2FA(c *gin.Context) {
	var req Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.service.Enable2FA(ctx, req.UserID, req.TenantID, req.Token)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

// Disable2FA disables 2FA for a user
func (h *Handler) Disable2FA(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid user ID"))
		return
	}

	ctx := c.Request.Context()
	err = h.service.Disable2FA(ctx, userID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

// VerifyBackupCodeRequest represents a backup code verification request
type VerifyBackupCodeRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Code   string    `json:"code" binding:"required"`
}

// VerifyBackupCode verifies a 2FA backup code
func (h *Handler) VerifyBackupCode(c *gin.Context) {
	var req VerifyBackupCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	valid, err := h.service.VerifyBackupCode(ctx, req.UserID, req.Code)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"valid": valid})
}

// SavePasswordHistoryRequest represents a password history save request
type SavePasswordHistoryRequest struct {
	UserID       uuid.UUID `json:"user_id" binding:"required"`
	PasswordHash string    `json:"password_hash" binding:"required"`
}

// SavePasswordHistory saves password to history
func (h *Handler) SavePasswordHistory(c *gin.Context) {
	var req SavePasswordHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.service.SavePasswordHistory(ctx, req.UserID, req.PasswordHash)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Password history saved"})
}

// CheckPasswordReuseRequest represents a password reuse check request
type CheckPasswordReuseRequest struct {
	UserID   uuid.UUID  `json:"user_id" binding:"required"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
	Password string     `json:"password" binding:"required"`
}

// CheckPasswordReuse checks if password has been used recently
func (h *Handler) CheckPasswordReuse(c *gin.Context) {
	var req CheckPasswordReuseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	reused, err := h.service.IsPasswordReused(ctx, req.UserID, req.TenantID, req.Password)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"reused": reused})
}

// UpdateSecuritySettingsRequest represents a security settings update request
type UpdateSecuritySettingsRequest struct {
	MaxFailedAttempts      int  `json:"max_failed_attempts,omitempty"`
	PasswordMinLength      int  `json:"password_min_length,omitempty"`
	PasswordRequireUpper   bool `json:"password_require_upper"`
	PasswordRequireLower   bool `json:"password_require_lower"`
	PasswordRequireDigit   bool `json:"password_require_digit"`
	PasswordRequireSpecial bool `json:"password_require_special"`
	Require2FA             bool `json:"require_2fa"`
}

// UpdateSecuritySettings updates security settings for a user
func (h *Handler) UpdateSecuritySettings(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid user ID"))
		return
	}

	var req UpdateSecuritySettingsRequest
	if bindError := c.ShouldBindJSON(&req); bindError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindError.Error()})
		return
	}

	// Convert request to SecuritySettings - in a real implementation,
	// you'd want to only update the provided fields
	settings := SecuritySettings{
		MaxFailedAttempts:      req.MaxFailedAttempts,
		PasswordMinLength:      req.PasswordMinLength,
		PasswordRequireUpper:   req.PasswordRequireUpper,
		PasswordRequireLower:   req.PasswordRequireLower,
		PasswordRequireDigit:   req.PasswordRequireDigit,
		PasswordRequireSpecial: req.PasswordRequireSpecial,
		Require2FA:             req.Require2FA,
	}

	ctx := c.Request.Context()
	err = h.service.UpdateSecuritySettings(ctx, userID, settings)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Security settings updated successfully"})
}
