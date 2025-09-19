package security

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserRepository interface for user operations needed by security service
type UserRepository interface {
	UpdateUser2FA(ctx context.Context, userID uuid.UUID, enabled bool) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (User, error) // Alias for consistency
}

// User interface that security service needs
type User interface {
	GetID() uuid.UUID
	GetTenantID() *uuid.UUID
	IsTwoFactorEnabled() bool
}

// SecuritySettings represents security configuration
type SecuritySettings struct {
	MaxFailedAttempts      int           `json:"max_failed_attempts" default:"5"`
	LockoutDuration        time.Duration `json:"lockout_duration" default:"15m"`
	PasswordMinLength      int           `json:"password_min_length" default:"8"`
	PasswordRequireUpper   bool          `json:"password_require_upper" default:"true"`
	PasswordRequireLower   bool          `json:"password_require_lower" default:"true"`
	PasswordRequireDigit   bool          `json:"password_require_digit" default:"true"`
	PasswordRequireSpecial bool          `json:"password_require_special" default:"true"`
	PasswordMaxAge         time.Duration `json:"password_max_age" default:"2160h"` // 90 days
	PasswordHistoryCount   int           `json:"password_history_count" default:"5"`
	SessionTimeout         time.Duration `json:"session_timeout" default:"24h"`
	Require2FA             bool          `json:"require_2fa" default:"false"`
}

// UserSecurityLog represents security-related events
type UserSecurityLog struct {
	ID        uuid.UUID              `json:"id" gorm:"primarykey"`
	UserID    uuid.UUID              `json:"user_id" gorm:"not null;index"`
	TenantID  *uuid.UUID             `json:"tenant_id,omitempty" gorm:"index"`
	EventType string                 `json:"event_type" gorm:"not null;index"` // login_success, login_failed, password_changed, etc.
	IPAddress string                 `json:"ip_address,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty" gorm:"type:jsonb"`
	CreatedAt time.Time              `json:"created_at"`
}

// UserAccountLockout represents account lockout information
type UserAccountLockout struct {
	ID             uuid.UUID  `json:"id" gorm:"primarykey"`
	UserID         uuid.UUID  `json:"user_id" gorm:"not null;unique;index"`
	TenantID       *uuid.UUID `json:"tenant_id,omitempty" gorm:"index"`
	FailedAttempts int        `json:"failed_attempts" gorm:"default:0"`
	LastFailedAt   *time.Time `json:"last_failed_at,omitempty"`
	LockedUntil    *time.Time `json:"locked_until,omitempty"`
	IsLocked       bool       `json:"is_locked" gorm:"default:false"`
	LockReason     string     `json:"lock_reason,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TwoFactorAuth represents 2FA configuration for a user
type TwoFactorAuth struct {
	ID          uuid.UUID  `json:"id" gorm:"primarykey"`
	UserID      uuid.UUID  `json:"user_id" gorm:"not null;unique;index"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty" gorm:"index"`
	Secret      string     `json:"-" gorm:"not null"`   // TOTP secret, encrypted
	BackupCodes []string   `json:"-" gorm:"type:jsonb"` // Encrypted backup codes
	IsEnabled   bool       `json:"is_enabled" gorm:"default:false"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// PasswordHistory represents user's password history
type PasswordHistory struct {
	ID           uuid.UUID  `json:"id" gorm:"primarykey"`
	UserID       uuid.UUID  `json:"user_id" gorm:"not null;index"`
	TenantID     *uuid.UUID `json:"tenant_id,omitempty" gorm:"index"`
	PasswordHash string     `json:"-" gorm:"not null"` // Hashed password
	CreatedAt    time.Time  `json:"created_at"`
}

// Business Logic Methods for UserAccountLockout

// IsCurrentlyLocked checks if the account is currently locked
func (l *UserAccountLockout) IsCurrentlyLocked() bool {
	if !l.IsLocked || l.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*l.LockedUntil)
}

// ShouldLockAccount checks if account should be locked based on failed attempts
func (l *UserAccountLockout) ShouldLockAccount(maxAttempts int) bool {
	return l.FailedAttempts >= maxAttempts
}

// ResetFailedAttempts resets the failed attempts counter
func (l *UserAccountLockout) ResetFailedAttempts() {
	l.FailedAttempts = 0
	l.LastFailedAt = nil
	l.IsLocked = false
	l.LockedUntil = nil
	l.LockReason = ""
}

// IncrementFailedAttempts increments the failed attempts counter
func (l *UserAccountLockout) IncrementFailedAttempts() {
	l.FailedAttempts++
	now := time.Now()
	l.LastFailedAt = &now
}

// LockAccount locks the account for the specified duration
func (l *UserAccountLockout) LockAccount(duration time.Duration, reason string) {
	l.IsLocked = true
	lockedUntil := time.Now().Add(duration)
	l.LockedUntil = &lockedUntil
	l.LockReason = reason
}

// Business Logic Methods for TwoFactorAuth

// IsSetupComplete checks if 2FA setup is complete
func (t *TwoFactorAuth) IsSetupComplete() bool {
	return t.Secret != "" && t.IsEnabled
}

// HasBackupCodes checks if backup codes are available
func (t *TwoFactorAuth) HasBackupCodes() bool {
	return len(t.BackupCodes) > 0
}

// Security Event Types
const (
	EventLoginSuccess    = "login_success"
	EventLoginFailed     = "login_failed"
	EventPasswordChanged = "password_changed"
	EventAccountLocked   = "account_locked"
	EventAccountUnlocked = "account_unlocked"
	Event2FAEnabled      = "2fa_enabled"
	Event2FADisabled     = "2fa_disabled"
	Event2FAUsed         = "2fa_used"
	Event2FAVerified     = "2fa_verified"
	EventBackupCodeUsed  = "backup_code_used"
	EventPasswordReset   = "password_reset"
	EventEmailChanged    = "email_changed"
	EventRoleChanged     = "role_changed"
	EventSessionCreated  = "session_created"
	EventSessionExpired  = "session_expired"
	EventSuspiciousLogin = "suspicious_login"
	EventAccountCreated  = "account_created"

	// Legacy constants for backward compatibility
	SecurityEventAccountCreated = EventAccountCreated
	SecurityEventFailedLogin    = EventLoginFailed
	SecurityEventPasswordChange = EventPasswordChanged
	SecurityEvent2FAVerified    = Event2FAVerified
	SecurityEvent2FAEnabled     = Event2FAEnabled
	SecurityEvent2FADisabled    = Event2FADisabled
	SecurityEventBackupCodeUsed = EventBackupCodeUsed
)

// Default security settings
var DefaultSecuritySettings = SecuritySettings{
	MaxFailedAttempts:      5,
	LockoutDuration:        15 * time.Minute,
	PasswordMinLength:      8,
	PasswordRequireUpper:   true,
	PasswordRequireLower:   true,
	PasswordRequireDigit:   true,
	PasswordRequireSpecial: true,
	PasswordMaxAge:         90 * 24 * time.Hour, // 90 days
	PasswordHistoryCount:   5,
	SessionTimeout:         24 * time.Hour,
	Require2FA:             false,
}

// GetDefaultSecuritySettings returns the default security settings
func GetDefaultSecuritySettings() SecuritySettings {
	return DefaultSecuritySettings
}
