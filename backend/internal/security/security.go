package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LoginAttemptStatus represents the status of a login attempt
type LoginAttemptStatus string

const (
	LoginAttemptSuccess    LoginAttemptStatus = "success"
	LoginAttemptFailed     LoginAttemptStatus = "failed"
	LoginAttemptBlocked    LoginAttemptStatus = "blocked"
	LoginAttemptSuspicious LoginAttemptStatus = "suspicious"
)

// ThreatLevel represents the threat level of a security event
type ThreatLevel string

const (
	ThreatLevelLow      ThreatLevel = "low"
	ThreatLevelMedium   ThreatLevel = "medium"
	ThreatLevelHigh     ThreatLevel = "high"
	ThreatLevelCritical ThreatLevel = "critical"
)

// DeviceStatus represents the trust status of a device
type DeviceStatus string

const (
	DeviceStatusTrusted    DeviceStatus = "trusted"
	DeviceStatusUntrusted  DeviceStatus = "untrusted"
	DeviceStatusSuspicious DeviceStatus = "suspicious"
	DeviceStatusBlocked    DeviceStatus = "blocked"
)

// PasswordPolicy represents password policy configuration
type PasswordPolicy struct {
	ID               uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID         *uuid.UUID `json:"tenant_id,omitempty" gorm:"index"` // Null for global policy
	
	// Length requirements
	MinLength        int `json:"min_length" gorm:"default:8"`
	MaxLength        int `json:"max_length" gorm:"default:128"`
	
	// Character requirements
	RequireUppercase bool `json:"require_uppercase" gorm:"default:true"`
	RequireLowercase bool `json:"require_lowercase" gorm:"default:true"`
	RequireNumbers   bool `json:"require_numbers" gorm:"default:true"`
	RequireSpecial   bool `json:"require_special" gorm:"default:true"`
	RequireSymbols   bool `json:"require_symbols" gorm:"default:true"`
	
	// History and expiration
	PasswordHistoryCount int `json:"password_history_count" gorm:"default:5"`
	HistoryCount         int `json:"history_count" gorm:"default:5"`
	ExpirationDays       int `json:"expiration_days" gorm:"default:90"` // 0 for no expiration
	MaxAge               int `json:"max_age" gorm:"default:90"`
	
	// Lockout settings
	MaxFailedAttempts    int `json:"max_failed_attempts" gorm:"default:5"`
	LockoutDurationMins  int `json:"lockout_duration_mins" gorm:"default:30"`
	
	// Common password prevention
	PreventCommonPasswords bool     `json:"prevent_common_passwords" gorm:"default:true"`
	PreventUserInfo        bool     `json:"prevent_user_info" gorm:"default:true"`
	ForbiddenPatterns      []string `json:"forbidden_patterns" gorm:"serializer:json"`
	
	IsActive bool `json:"is_active" gorm:"default:true"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginAttempt tracks all login attempts for security monitoring
type LoginAttempt struct {
	ID           uuid.UUID          `json:"id" gorm:"primarykey"`
	UserID       *uuid.UUID         `json:"user_id,omitempty" gorm:"index"` // Null for failed attempts with non-existent users
	Email        string             `json:"email" gorm:"not null;index"`
	Status       LoginAttemptStatus `json:"status" gorm:"not null"`
	
	// Request details
	IPAddress    string `json:"ip_address" gorm:"not null;index"`
	UserAgent    string `json:"user_agent"`
	Country      string `json:"country,omitempty"`
	City         string `json:"city,omitempty"`
	
	// Security context
	DeviceFingerprint string      `json:"device_fingerprint" gorm:"index"`
	ThreatLevel       ThreatLevel `json:"threat_level" gorm:"default:low"`
	
	// Failure details
	FailureReason string `json:"failure_reason,omitempty"`
	BlockedReason string `json:"blocked_reason,omitempty"`
	
	// Timing
	AttemptedAt time.Time `json:"attempted_at" gorm:"not null"`
	ProcessedAt time.Time `json:"processed_at" gorm:"not null"`
	
	CreatedAt time.Time `json:"created_at"`
}

// TrustedDevice represents a device that has been verified by the user
type TrustedDevice struct {
	ID           uuid.UUID    `json:"id" gorm:"primarykey"`
	UserID       uuid.UUID    `json:"user_id" gorm:"not null;index"`
	DeviceID     string       `json:"device_id" gorm:"not null"` // Generated device identifier
	Fingerprint  string       `json:"fingerprint" gorm:"not null;index"`
	Name         string       `json:"name" gorm:"not null"` // User-friendly device name
	DeviceName   string       `json:"device_name"` // Alternative device name field
	DeviceType   string       `json:"device_type"` // mobile, desktop, tablet
	OS           string       `json:"os"`
	Browser      string       `json:"browser"`
	Status       DeviceStatus `json:"status" gorm:"default:trusted"`
	
	// Security tracking
	FirstSeenAt      time.Time  `json:"first_seen_at" gorm:"not null"`
	LastSeenAt       time.Time  `json:"last_seen_at" gorm:"not null"`
	LastIPAddress    string     `json:"last_ip_address"`
	IPAddress        string     `json:"ip_address"` // Current IP address
	UserAgent        string     `json:"user_agent"` // User agent string
	TrustScore       float64    `json:"trust_score" gorm:"default:1.0"` // 0.0-1.0
	AccessCount      int        `json:"access_count" gorm:"default:0"`
	Country          string     `json:"country,omitempty"`
	City             string     `json:"city,omitempty"`
	VerifiedAt       *time.Time `json:"verified_at,omitempty"`
	VerificationCode string     `json:"verification_code,omitempty"`
	
	// Revocation
	RevokedAt     *time.Time `json:"revoked_at,omitempty"`
	RevokedReason string     `json:"revoked_reason,omitempty"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SecurityEvent tracks security-related events for audit and monitoring
type SecurityEvent struct {
	ID       uuid.UUID   `json:"id" gorm:"primarykey"`
	UserID   *uuid.UUID  `json:"user_id,omitempty" gorm:"index"`
	TenantID *uuid.UUID  `json:"tenant_id,omitempty" gorm:"index"`
	
	// Event details
	EventType    string      `json:"event_type" gorm:"not null;index"` // login_failed, password_changed, device_added, etc.
	ThreatLevel  ThreatLevel `json:"threat_level" gorm:"not null"`
	Description  string      `json:"description" gorm:"not null"`
	
	// Context
	IPAddress         string                 `json:"ip_address,omitempty"`
	UserAgent         string                 `json:"user_agent,omitempty"`
	DeviceFingerprint string                 `json:"device_fingerprint,omitempty"`
	Country           string                 `json:"country,omitempty"`
	City              string                 `json:"city,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty" gorm:"serializer:json"`
	
	// Resolution
	IsResolved       bool       `json:"is_resolved" gorm:"default:false"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy       *uuid.UUID `json:"resolved_by,omitempty"`
	ResolutionNotes  string     `json:"resolution_notes,omitempty"`
	
	OccurredAt time.Time `json:"occurred_at" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PasswordHistory tracks user's previous passwords to prevent reuse
type PasswordHistory struct {
	ID           uuid.UUID `json:"id" gorm:"primarykey"`
	UserID       uuid.UUID `json:"user_id" gorm:"not null;index"`
	PasswordHash string    `json:"password_hash" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null"`
}

// AccountLockout tracks user account lockouts due to security violations
type AccountLockout struct {
	ID           uuid.UUID   `json:"id" gorm:"primarykey"`
	UserID       uuid.UUID   `json:"user_id" gorm:"not null;index"`
	LockoutType  string      `json:"lockout_type" gorm:"not null"` // failed_login, suspicious_activity, admin_action
	Reason       string      `json:"reason" gorm:"not null"`
	ThreatLevel  ThreatLevel `json:"threat_level" gorm:"not null"`
	
	// Lockout details
	LockedAt      time.Time  `json:"locked_at" gorm:"not null"`
	UnlocksAt     *time.Time `json:"unlocks_at,omitempty"`
	UnlockedAt    *time.Time `json:"unlocked_at,omitempty"`
	UnlockedBy    *uuid.UUID `json:"unlocked_by,omitempty"`
	
	// Context
	IPAddress         string                 `json:"ip_address,omitempty"`
	DeviceFingerprint string                 `json:"device_fingerprint,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty" gorm:"serializer:json"`
	
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EncryptionKey represents encryption keys for data protection
type EncryptionKey struct {
	ID          uuid.UUID  `json:"id" gorm:"primarykey"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty" gorm:"index"` // Null for system-wide keys
	KeyName     string     `json:"key_name" gorm:"not null"`
	KeyType     string     `json:"key_type" gorm:"not null"` // AES256, RSA2048, etc.
	KeyPurpose  string     `json:"key_purpose" gorm:"not null"` // data_encryption, token_signing, etc.
	
	// Key material (encrypted)
	EncryptedKey string `json:"encrypted_key" gorm:"not null"`
	KeyVersion   int    `json:"key_version" gorm:"default:1"`
	Algorithm    string `json:"algorithm" gorm:"not null"`
	
	// Lifecycle
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	ActivatedAt time.Time  `json:"activated_at" gorm:"not null"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Business Logic Methods

// IsPasswordExpired checks if the password has expired based on policy
func (u *PasswordPolicy) IsPasswordExpired(passwordChangedAt *time.Time) bool {
	if u.ExpirationDays == 0 || passwordChangedAt == nil {
		return false
	}
	
	expirationDate := passwordChangedAt.AddDate(0, 0, u.ExpirationDays)
	return time.Now().After(expirationDate)
}

// ValidatePassword validates a password against the policy
func (u *PasswordPolicy) ValidatePassword(password, userEmail, firstName, lastName string) error {
	if len(password) < u.MinLength {
		return fmt.Errorf("password must be at least %d characters long", u.MinLength)
	}
	
	if u.MaxLength > 0 && len(password) > u.MaxLength {
		return fmt.Errorf("password must be no more than %d characters long", u.MaxLength)
	}
	
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case (char >= 33 && char <= 47) || (char >= 58 && char <= 64) || 
			 (char >= 91 && char <= 96) || (char >= 123 && char <= 126):
			hasSpecial = true
		}
	}
	
	if u.RequireUppercase && !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	
	if u.RequireLowercase && !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	
	if u.RequireNumbers && !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	
	if u.RequireSpecial && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	
	// Check against user information
	if u.PreventUserInfo {
		if containsUserInfo(password, userEmail, firstName, lastName) {
			return fmt.Errorf("password cannot contain personal information")
		}
	}
	
	// Check against common passwords
	if u.PreventCommonPasswords {
		if isCommonPassword(password) {
			return fmt.Errorf("password is too common, please choose a stronger password")
		}
	}
	
	return nil
}

// ShouldLockAccount determines if an account should be locked based on failed attempts
func (u *PasswordPolicy) ShouldLockAccount(failedAttempts int) bool {
	return failedAttempts >= u.MaxFailedAttempts
}

// GetLockoutDuration returns the lockout duration
func (u *PasswordPolicy) GetLockoutDuration() time.Duration {
	return time.Duration(u.LockoutDurationMins) * time.Minute
}

// IsActive checks if the device is trusted and active
func (d *TrustedDevice) IsActive() bool {
	return d.Status == DeviceStatusTrusted && d.RevokedAt == nil
}

// ShouldUpdateTrustScore determines if trust score should be updated
func (d *TrustedDevice) ShouldUpdateTrustScore(ipAddress string, timesSeen int) float64 {
	newScore := d.TrustScore
	
	// Decrease trust if IP address changed frequently
	if d.IPAddress != ipAddress {
		newScore -= 5
	}
	
	// Increase trust based on usage frequency
	if timesSeen > 10 {
		newScore += 2
	}
	
	// Ensure score stays within bounds
	if newScore > 100 {
		newScore = 100
	}
	if newScore < 0 {
		newScore = 0
	}
	
	return newScore
}

// IsLocked checks if the lockout is currently active
func (l *AccountLockout) IsLocked() bool {
	if !l.IsActive || l.UnlockedAt != nil {
		return false
	}
	
	if l.UnlocksAt != nil && time.Now().After(*l.UnlocksAt) {
		return false
	}
	
	return true
}

// GetRemainingLockoutTime returns remaining lockout duration
func (l *AccountLockout) GetRemainingLockoutTime() time.Duration {
	if !l.IsLocked() {
		return 0
	}
	
	if l.UnlocksAt == nil {
		return time.Hour * 24 * 365 // Indefinite
	}
	
	remaining := time.Until(*l.UnlocksAt)
	if remaining < 0 {
		return 0
	}
	
	return remaining
}

// IsHigh checks if the event is high or critical threat level
func (e *SecurityEvent) IsHigh() bool {
	return e.ThreatLevel == ThreatLevelHigh || e.ThreatLevel == ThreatLevelCritical
}

// IsExpired checks if the encryption key has expired
func (k *EncryptionKey) IsExpired() bool {
	return k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt)
}

// IsRevoked checks if the encryption key has been revoked
func (k *EncryptionKey) IsRevoked() bool {
	return k.RevokedAt != nil
}

// IsUsable checks if the key is active and not expired or revoked
func (k *EncryptionKey) IsUsable() bool {
	return k.IsActive && !k.IsExpired() && !k.IsRevoked()
}

// Helper functions

func containsUserInfo(password, email, firstName, lastName string) bool {
	password = strings.ToLower(password)
	
	// Check against email parts
	if strings.Contains(password, strings.ToLower(strings.Split(email, "@")[0])) {
		return true
	}
	
	// Check against names
	if len(firstName) > 2 && strings.Contains(password, strings.ToLower(firstName)) {
		return true
	}
	
	if len(lastName) > 2 && strings.Contains(password, strings.ToLower(lastName)) {
		return true
	}
	
	return false
}

func isCommonPassword(password string) bool {
	// List of common passwords (simplified for example)
	commonPasswords := map[string]bool{
		"password":   true,
		"123456":     true,
		"password123": true,
		"admin":      true,
		"qwerty":     true,
		"letmein":    true,
		"welcome":    true,
		"monkey":     true,
		"dragon":     true,
		"master":     true,
	}
	
	return commonPasswords[strings.ToLower(password)]
}

// GenerateDeviceID generates a unique device identifier
func GenerateDeviceID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GenerateVerificationCode generates a random verification code
func GenerateVerificationCode() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return fmt.Sprintf("%06d", int(bytes[0])<<16|int(bytes[1])<<8|int(bytes[2]))
}

// ConstantTimeCompare performs constant-time string comparison
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
