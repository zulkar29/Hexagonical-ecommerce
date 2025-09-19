package security

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SecurityService handles user security operations
type SecurityService struct {
	db       *gorm.DB
	userRepo UserRepository
	settings SecuritySettings
}

// NewSecurityService creates a new security service
func NewSecurityService(db *gorm.DB, userRepo UserRepository) *SecurityService {
	return &SecurityService{
		db:       db,
		userRepo: userRepo,
		settings: GetDefaultSecuritySettings(),
	}
}

// ValidatePassword validates password against security policies
func (s *SecurityService) ValidatePassword(ctx context.Context, tenantID *uuid.UUID, password string) error {
	if len(password) < s.settings.PasswordMinLength {
		return fmt.Errorf("password must be at least %d characters long", s.settings.PasswordMinLength)
	}

	if s.settings.PasswordRequireUpper {
		if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
			return errors.New("password must contain at least one uppercase letter")
		}
	}

	if s.settings.PasswordRequireLower {
		if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
			return errors.New("password must contain at least one lowercase letter")
		}
	}

	if s.settings.PasswordRequireDigit {
		if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
			return errors.New("password must contain at least one digit")
		}
	}

	if s.settings.PasswordRequireSpecial {
		if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password); !matched {
			return errors.New("password must contain at least one special character")
		}
	}

	return nil
}

// CheckPasswordHistory validates password against history to prevent reuse
func (s *SecurityService) CheckPasswordHistory(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID, newPassword string, historyCount int) error {
	var histories []PasswordHistory
	query := s.db.WithContext(ctx).Where("user_id = ?", userID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.Order("created_at DESC").Limit(historyCount).Find(&histories).Error
	if err != nil {
		return err
	}

	for _, history := range histories {
		if err := bcrypt.CompareHashAndPassword([]byte(history.PasswordHash), []byte(newPassword)); err == nil {
			return errors.New("password cannot be the same as any of the last passwords used")
		}
	}

	return nil
}

// RecordFailedLogin records a failed login attempt and handles account lockout
func (s *SecurityService) RecordFailedLogin(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID, ipAddress, userAgent string) error {
	// Log the security event
	if err := s.LogSecurityEvent(ctx, userID, tenantID, EventLoginFailed, ipAddress, userAgent, nil); err != nil {
		return err
	}

	// Get or create lockout record
	lockout, err := s.getOrCreateLockout(ctx, userID, tenantID)
	if err != nil {
		return err
	}

	// Check if already locked
	if lockout.IsCurrentlyLocked() {
		return errors.New("account is currently locked")
	}

	// Increment failed attempts
	lockout.IncrementFailedAttempts()

	// Check if should lock account
	if lockout.ShouldLockAccount(s.settings.MaxFailedAttempts) {
		lockout.LockAccount(s.settings.LockoutDuration, "Too many failed login attempts")

		// Log account locked event
		if err := s.LogSecurityEvent(ctx, userID, tenantID, EventAccountLocked, ipAddress, userAgent, map[string]interface{}{
			"failed_attempts": lockout.FailedAttempts,
			"locked_until":    lockout.LockedUntil,
		}); err != nil {
			return err
		}
	}

	return s.db.WithContext(ctx).Save(lockout).Error
}

// RecordSuccessfulLogin records a successful login and resets failed attempts
func (s *SecurityService) RecordSuccessfulLogin(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID, ipAddress, userAgent string) error {
	// Log the security event
	if err := s.LogSecurityEvent(ctx, userID, tenantID, EventLoginSuccess, ipAddress, userAgent, nil); err != nil {
		return err
	}

	// Reset failed attempts
	lockout, err := s.getOrCreateLockout(ctx, userID, tenantID)
	if err != nil {
		return err
	}

	lockout.ResetFailedAttempts()
	return s.db.WithContext(ctx).Save(lockout).Error
}

// IsAccountLocked checks if an account is currently locked
func (s *SecurityService) IsAccountLocked(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) (bool, error) {
	lockout, err := s.getOrCreateLockout(ctx, userID, tenantID)
	if err != nil {
		return false, err
	}

	return lockout.IsCurrentlyLocked(), nil
}

// Setup2FA initializes 2FA for a user
func (s *SecurityService) Setup2FA(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID, issuer, accountName string) (*TwoFactorAuth, string, error) {
	// Generate secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return nil, "", err
	}

	// Generate backup codes
	backupCodes, err := s.generateBackupCodes(10)
	if err != nil {
		return nil, "", err
	}

	// Create 2FA record
	twoFA := &TwoFactorAuth{
		ID:          uuid.New(),
		UserID:      userID,
		TenantID:    tenantID,
		Secret:      key.Secret(),
		BackupCodes: backupCodes,
		IsEnabled:   false, // Will be enabled after verification
	}

	if err := s.db.WithContext(ctx).Create(twoFA).Error; err != nil {
		return nil, "", err
	}

	return twoFA, key.URL(), nil
}

// Verify2FA verifies a 2FA token
func (s *SecurityService) Verify2FA(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID, token string) (bool, error) {
	var twoFA TwoFactorAuth
	query := s.db.WithContext(ctx).Where("user_id = ?", userID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.First(&twoFA).Error; err != nil {
		return false, err
	}

	// Verify TOTP token
	valid := totp.Validate(token, twoFA.Secret)
	if valid {
		// Update last used time
		now := time.Now()
		twoFA.LastUsedAt = &now
		s.db.WithContext(ctx).Save(&twoFA)

		// Log 2FA usage
		s.LogSecurityEvent(ctx, userID, tenantID, Event2FAUsed, "", "", nil)
		return true, nil
	}

	// Check backup codes
	for i, code := range twoFA.BackupCodes {
		if code == token {
			// Remove used backup code
			twoFA.BackupCodes = append(twoFA.BackupCodes[:i], twoFA.BackupCodes[i+1:]...)
			now := time.Now()
			twoFA.LastUsedAt = &now
			s.db.WithContext(ctx).Save(&twoFA)

			// Log backup code usage
			s.LogSecurityEvent(ctx, userID, tenantID, Event2FAUsed, "", "", map[string]interface{}{
				"backup_code_used": true,
			})
			return true, nil
		}
	}

	return false, nil
}

// Enable2FA enables 2FA for a user after verification
func (s *SecurityService) Enable2FA(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID, token string) error {
	valid, err := s.Verify2FA(ctx, userID, tenantID, token)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("invalid 2FA token")
	}

	// Enable 2FA
	query := s.db.WithContext(ctx).Model(&TwoFactorAuth{}).Where("user_id = ?", userID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Update("is_enabled", true).Error; err != nil {
		return err
	}

	// Update user record
	if err := s.userRepo.UpdateUser2FA(ctx, userID, true); err != nil {
		return err
	}

	// Log 2FA enabled
	return s.LogSecurityEvent(ctx, userID, tenantID, Event2FAEnabled, "", "", nil)
}

// LogSecurityEvent logs a security-related event
func (s *SecurityService) LogSecurityEvent(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID, eventType, ipAddress, userAgent string, details map[string]interface{}) error {
	log := &UserSecurityLog{
		ID:        uuid.New(),
		UserID:    userID,
		TenantID:  tenantID,
		EventType: eventType,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
	}

	return s.db.WithContext(ctx).Create(log).Error
}

// IsPasswordReused checks if password has been used recently
func (s *SecurityService) IsPasswordReused(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID, password string) (bool, error) {
	return s.CheckPasswordHistory(ctx, userID, tenantID, password, s.settings.PasswordHistoryCount) != nil, nil
}

// SavePasswordHistory saves password to history
func (s *SecurityService) SavePasswordHistory(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	// Get user to determine tenant_id
	currentUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	history := &PasswordHistory{
		ID:           uuid.New(),
		UserID:       userID,
		TenantID:     currentUser.GetTenantID(),
		PasswordHash: passwordHash,
	}

	return s.db.WithContext(ctx).Create(history).Error
}

// ResetFailedLogins resets failed login attempts for a user
func (s *SecurityService) ResetFailedLogins(ctx context.Context, userID uuid.UUID) error {
	// Get user to determine tenant_id
	currentUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	lockout, err := s.getOrCreateLockout(ctx, userID, currentUser.GetTenantID())
	if err != nil {
		return err
	}

	lockout.ResetFailedAttempts()
	return s.db.WithContext(ctx).Save(lockout).Error
}

// Disable2FA disables 2FA for a user
func (s *SecurityService) Disable2FA(ctx context.Context, userID uuid.UUID) error {
	// Get user to determine tenant_id
	currentUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Disable 2FA
	query := s.db.WithContext(ctx).Model(&TwoFactorAuth{}).Where("user_id = ?", userID)
	if currentUser.GetTenantID() != nil {
		query = query.Where("tenant_id = ?", *currentUser.GetTenantID())
	}

	if err := query.Update("is_enabled", false).Error; err != nil {
		return err
	}

	// Update user record
	if err := s.userRepo.UpdateUser2FA(ctx, userID, false); err != nil {
		return err
	}

	// Log 2FA disabled
	return s.LogSecurityEvent(ctx, userID, currentUser.GetTenantID(), Event2FADisabled, "", "", nil)
}

// VerifyBackupCode verifies a 2FA backup code
func (s *SecurityService) VerifyBackupCode(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	// Get user to determine tenant_id
	currentUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
	}

	var twoFA TwoFactorAuth
	query := s.db.WithContext(ctx).Where("user_id = ?", userID)
	if currentUser.GetTenantID() != nil {
		query = query.Where("tenant_id = ?", *currentUser.GetTenantID())
	}

	if err := query.First(&twoFA).Error; err != nil {
		return false, err
	}

	// Check backup codes
	for i, backupCode := range twoFA.BackupCodes {
		if backupCode == code {
			// Remove used backup code
			twoFA.BackupCodes = append(twoFA.BackupCodes[:i], twoFA.BackupCodes[i+1:]...)
			now := time.Now()
			twoFA.LastUsedAt = &now
			s.db.WithContext(ctx).Save(&twoFA)

			// Log backup code usage
			s.LogSecurityEvent(ctx, userID, currentUser.GetTenantID(), EventBackupCodeUsed, "", "", map[string]interface{}{
				"backup_code_used": true,
			})
			return true, nil
		}
	}

	return false, nil
}

// GetSecuritySettings returns security settings for a user
func (s *SecurityService) GetSecuritySettings(ctx context.Context, userID uuid.UUID) (SecuritySettings, error) {
	// For now, return default settings. In the future, this could be user-specific
	return s.settings, nil
}

// UpdateSecuritySettings updates security settings for a user
func (s *SecurityService) UpdateSecuritySettings(ctx context.Context, userID uuid.UUID, settings SecuritySettings) error {
	// For now, this is a no-op. In the future, this could store user-specific settings
	return nil
}

// GetSecurityLogs returns security logs for a user
func (s *SecurityService) GetSecurityLogs(ctx context.Context, userID uuid.UUID, limit int) ([]UserSecurityLog, error) {
	// Get user to determine tenant_id
	currentUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var logs []UserSecurityLog
	query := s.db.WithContext(ctx).Where("user_id = ?", userID)
	if currentUser.GetTenantID() != nil {
		query = query.Where("tenant_id = ?", *currentUser.GetTenantID())
	}

	err = query.Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

// Helper methods

func (s *SecurityService) getOrCreateLockout(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) (*UserAccountLockout, error) {
	var lockout UserAccountLockout
	query := s.db.WithContext(ctx).Where("user_id = ?", userID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.First(&lockout).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new lockout record
			lockout = UserAccountLockout{
				ID:       uuid.New(),
				UserID:   userID,
				TenantID: tenantID,
			}
			if err := s.db.WithContext(ctx).Create(&lockout).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &lockout, nil
}

func (s *SecurityService) generateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		bytes := make([]byte, 5)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		codes[i] = strings.ToUpper(base32.StdEncoding.EncodeToString(bytes)[:8])
	}
	return codes, nil
}
