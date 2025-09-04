package security

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SecurityRepository defines the interface for security data operations
type SecurityRepository interface {
	// Password Policies
	CreatePasswordPolicy(ctx context.Context, policy *PasswordPolicy) error
	GetPasswordPolicy(ctx context.Context, tenantID *uuid.UUID) (*PasswordPolicy, error)
	UpdatePasswordPolicy(ctx context.Context, policy *PasswordPolicy) error
	
	// Login Attempts
	CreateLoginAttempt(ctx context.Context, attempt *LoginAttempt) error
	GetLoginAttempts(ctx context.Context, filter LoginAttemptFilter) ([]*LoginAttempt, error)
	GetFailedLoginCount(ctx context.Context, email string, since time.Time) (int, error)
	GetRecentLoginAttempts(ctx context.Context, userID uuid.UUID, limit int) ([]*LoginAttempt, error)
	
	// Trusted Devices
	CreateTrustedDevice(ctx context.Context, device *TrustedDevice) error
	GetTrustedDevice(ctx context.Context, userID uuid.UUID, fingerprint string) (*TrustedDevice, error)
	GetTrustedDevices(ctx context.Context, userID uuid.UUID) ([]*TrustedDevice, error)
	UpdateTrustedDevice(ctx context.Context, device *TrustedDevice) error
	RevokeTrustedDevice(ctx context.Context, deviceID uuid.UUID, reason string) error
	CleanupExpiredDevices(ctx context.Context, days int) error
	
	// Security Events
	CreateSecurityEvent(ctx context.Context, event *SecurityEvent) error
	GetSecurityEvents(ctx context.Context, filter SecurityEventFilter) ([]*SecurityEvent, int64, error)
	GetUnresolvedEvents(ctx context.Context, threatLevel ThreatLevel) ([]*SecurityEvent, error)
	UpdateSecurityEvent(ctx context.Context, event *SecurityEvent) error
	
	// Password History
	CreatePasswordHistory(ctx context.Context, history *PasswordHistory) error
	GetPasswordHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*PasswordHistory, error)
	CleanupOldPasswordHistory(ctx context.Context, userID uuid.UUID, keepCount int) error
	
	// Account Lockouts
	CreateAccountLockout(ctx context.Context, lockout *AccountLockout) error
	GetActiveAccountLockout(ctx context.Context, userID uuid.UUID) (*AccountLockout, error)
	GetAccountLockouts(ctx context.Context, userID uuid.UUID) ([]*AccountLockout, error)
	UpdateAccountLockout(ctx context.Context, lockout *AccountLockout) error
	ExpireAccountLockouts(ctx context.Context) error
	
	// Encryption Keys
	CreateEncryptionKey(ctx context.Context, key *EncryptionKey) error
	GetEncryptionKey(ctx context.Context, keyName string, tenantID *uuid.UUID) (*EncryptionKey, error)
	GetActiveEncryptionKeys(ctx context.Context, tenantID *uuid.UUID) ([]*EncryptionKey, error)
	UpdateEncryptionKey(ctx context.Context, key *EncryptionKey) error
	RevokeEncryptionKey(ctx context.Context, keyID uuid.UUID) error
	
	// Analytics
	GetSecurityMetrics(ctx context.Context, filter SecurityMetricsFilter) (*SecurityMetrics, error)
}

// Filter types
type LoginAttemptFilter struct {
	UserID        *uuid.UUID          `json:"user_id,omitempty"`
	Email         *string             `json:"email,omitempty"`
	Status        *LoginAttemptStatus `json:"status,omitempty"`
	IPAddress     *string             `json:"ip_address,omitempty"`
	ThreatLevel   *ThreatLevel        `json:"threat_level,omitempty"`
	StartTime     *time.Time          `json:"start_time,omitempty"`
	EndTime       *time.Time          `json:"end_time,omitempty"`
	Limit         int                 `json:"limit"`
	Offset        int                 `json:"offset"`
}

type SecurityEventFilter struct {
	UserID      *uuid.UUID   `json:"user_id,omitempty"`
	TenantID    *uuid.UUID   `json:"tenant_id,omitempty"`
	EventType   *string      `json:"event_type,omitempty"`
	ThreatLevel *ThreatLevel `json:"threat_level,omitempty"`
	IsResolved  *bool        `json:"is_resolved,omitempty"`
	StartTime   *time.Time   `json:"start_time,omitempty"`
	EndTime     *time.Time   `json:"end_time,omitempty"`
	Limit       int          `json:"limit"`
	Offset      int          `json:"offset"`
}

type SecurityMetricsFilter struct {
	TenantID  *uuid.UUID `json:"tenant_id,omitempty"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
}

// Metrics types
type SecurityMetrics struct {
	Period                time.Duration            `json:"period"`
	TotalLoginAttempts    int64                    `json:"total_login_attempts"`
	SuccessfulLogins      int64                    `json:"successful_logins"`
	FailedLogins          int64                    `json:"failed_logins"`
	BlockedAttempts       int64                    `json:"blocked_attempts"`
	LoginSuccessRate      float64                  `json:"login_success_rate"`
	UniqueUsers           int64                    `json:"unique_users"`
	SuspiciousActivities  int64                    `json:"suspicious_activities"`
	AccountLockouts       int64                    `json:"account_lockouts"`
	TrustedDevices        int64                    `json:"trusted_devices"`
	SecurityEvents        int64                    `json:"security_events"`
	ThreatsByLevel        map[ThreatLevel]int64    `json:"threats_by_level"`
	AttacksByCountry      map[string]int64         `json:"attacks_by_country"`
	TopAttackerIPs        []IPThreatInfo           `json:"top_attacker_ips"`
	DeviceTrustDistribution map[DeviceStatus]int64 `json:"device_trust_distribution"`
}

type IPThreatInfo struct {
	IPAddress    string `json:"ip_address"`
	AttemptCount int64  `json:"attempt_count"`
	Country      string `json:"country"`
	LastSeen     time.Time `json:"last_seen"`
}

// gormSecurityRepository implements SecurityRepository using GORM
type gormSecurityRepository struct {
	db *gorm.DB
}

// NewSecurityRepository creates a new security repository
func NewSecurityRepository(db *gorm.DB) SecurityRepository {
	return &gormSecurityRepository{db: db}
}

// Password Policies
func (r *gormSecurityRepository) CreatePasswordPolicy(ctx context.Context, policy *PasswordPolicy) error {
	return r.db.WithContext(ctx).Create(policy).Error
}

func (r *gormSecurityRepository) GetPasswordPolicy(ctx context.Context, tenantID *uuid.UUID) (*PasswordPolicy, error) {
	var policy PasswordPolicy
	query := r.db.WithContext(ctx).Where("is_active = ?", true)
	
	if tenantID != nil {
		// Try tenant-specific policy first
		err := query.Where("tenant_id = ?", *tenantID).First(&policy).Error
		if err == nil {
			return &policy, nil
		}
	}
	
	// Fall back to global policy
	err := query.Where("tenant_id IS NULL").First(&policy).Error
	if err != nil {
		return nil, err
	}
	
	return &policy, nil
}

func (r *gormSecurityRepository) UpdatePasswordPolicy(ctx context.Context, policy *PasswordPolicy) error {
	return r.db.WithContext(ctx).Save(policy).Error
}

// Login Attempts
func (r *gormSecurityRepository) CreateLoginAttempt(ctx context.Context, attempt *LoginAttempt) error {
	return r.db.WithContext(ctx).Create(attempt).Error
}

func (r *gormSecurityRepository) GetLoginAttempts(ctx context.Context, filter LoginAttemptFilter) ([]*LoginAttempt, error) {
	query := r.db.WithContext(ctx).Model(&LoginAttempt{})
	
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Email != nil {
		query = query.Where("email = ?", *filter.Email)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.IPAddress != nil {
		query = query.Where("ip_address = ?", *filter.IPAddress)
	}
	if filter.ThreatLevel != nil {
		query = query.Where("threat_level = ?", *filter.ThreatLevel)
	}
	if filter.StartTime != nil {
		query = query.Where("attempted_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("attempted_at <= ?", *filter.EndTime)
	}
	
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	
	var attempts []*LoginAttempt
	err := query.Order("attempted_at DESC").Find(&attempts).Error
	return attempts, err
}

func (r *gormSecurityRepository) GetFailedLoginCount(ctx context.Context, email string, since time.Time) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&LoginAttempt{}).
		Where("email = ? AND status = ? AND attempted_at >= ?", email, LoginAttemptFailed, since).
		Count(&count).Error
	return int(count), err
}

func (r *gormSecurityRepository) GetRecentLoginAttempts(ctx context.Context, userID uuid.UUID, limit int) ([]*LoginAttempt, error) {
	var attempts []*LoginAttempt
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("attempted_at DESC").
		Limit(limit).
		Find(&attempts).Error
	return attempts, err
}

// Trusted Devices
func (r *gormSecurityRepository) CreateTrustedDevice(ctx context.Context, device *TrustedDevice) error {
	return r.db.WithContext(ctx).Create(device).Error
}

func (r *gormSecurityRepository) GetTrustedDevice(ctx context.Context, userID uuid.UUID, fingerprint string) (*TrustedDevice, error) {
	var device TrustedDevice
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND fingerprint = ? AND revoked_at IS NULL", userID, fingerprint).
		First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *gormSecurityRepository) GetTrustedDevices(ctx context.Context, userID uuid.UUID) ([]*TrustedDevice, error) {
	var devices []*TrustedDevice
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Order("last_seen_at DESC").
		Find(&devices).Error
	return devices, err
}

func (r *gormSecurityRepository) UpdateTrustedDevice(ctx context.Context, device *TrustedDevice) error {
	return r.db.WithContext(ctx).Save(device).Error
}

func (r *gormSecurityRepository) RevokeTrustedDevice(ctx context.Context, deviceID uuid.UUID, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&TrustedDevice{}).
		Where("id = ?", deviceID).
		Updates(map[string]interface{}{
			"revoked_at":     &now,
			"revoked_reason": reason,
			"status":         DeviceStatusBlocked,
		}).Error
}

func (r *gormSecurityRepository) CleanupExpiredDevices(ctx context.Context, days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.WithContext(ctx).
		Where("last_seen_at < ?", cutoff).
		Delete(&TrustedDevice{}).Error
}

// Security Events
func (r *gormSecurityRepository) CreateSecurityEvent(ctx context.Context, event *SecurityEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *gormSecurityRepository) GetSecurityEvents(ctx context.Context, filter SecurityEventFilter) ([]*SecurityEvent, int64, error) {
	query := r.db.WithContext(ctx).Model(&SecurityEvent{})
	
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", *filter.TenantID)
	}
	if filter.EventType != nil {
		query = query.Where("event_type = ?", *filter.EventType)
	}
	if filter.ThreatLevel != nil {
		query = query.Where("threat_level = ?", *filter.ThreatLevel)
	}
	if filter.IsResolved != nil {
		query = query.Where("is_resolved = ?", *filter.IsResolved)
	}
	if filter.StartTime != nil {
		query = query.Where("occurred_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("occurred_at <= ?", *filter.EndTime)
	}
	
	// Count total records
	var total int64
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	
	var events []*SecurityEvent
	err := query.Order("occurred_at DESC").Find(&events).Error
	return events, total, err
}

func (r *gormSecurityRepository) GetUnresolvedEvents(ctx context.Context, threatLevel ThreatLevel) ([]*SecurityEvent, error) {
	var events []*SecurityEvent
	err := r.db.WithContext(ctx).
		Where("is_resolved = ? AND threat_level >= ?", false, threatLevel).
		Order("occurred_at DESC").
		Find(&events).Error
	return events, err
}

func (r *gormSecurityRepository) UpdateSecurityEvent(ctx context.Context, event *SecurityEvent) error {
	return r.db.WithContext(ctx).Save(event).Error
}

// Password History
func (r *gormSecurityRepository) CreatePasswordHistory(ctx context.Context, history *PasswordHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *gormSecurityRepository) GetPasswordHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*PasswordHistory, error) {
	var history []*PasswordHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}

func (r *gormSecurityRepository) CleanupOldPasswordHistory(ctx context.Context, userID uuid.UUID, keepCount int) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND id NOT IN (SELECT id FROM password_histories WHERE user_id = ? ORDER BY created_at DESC LIMIT ?)", 
			userID, userID, keepCount).
		Delete(&PasswordHistory{}).Error
}

// Account Lockouts
func (r *gormSecurityRepository) CreateAccountLockout(ctx context.Context, lockout *AccountLockout) error {
	return r.db.WithContext(ctx).Create(lockout).Error
}

func (r *gormSecurityRepository) GetActiveAccountLockout(ctx context.Context, userID uuid.UUID) (*AccountLockout, error) {
	var lockout AccountLockout
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ? AND unlocked_at IS NULL", userID, true).
		Where("unlocks_at IS NULL OR unlocks_at > ?", time.Now()).
		First(&lockout).Error
	if err != nil {
		return nil, err
	}
	return &lockout, nil
}

func (r *gormSecurityRepository) GetAccountLockouts(ctx context.Context, userID uuid.UUID) ([]*AccountLockout, error) {
	var lockouts []*AccountLockout
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("locked_at DESC").
		Find(&lockouts).Error
	return lockouts, err
}

func (r *gormSecurityRepository) UpdateAccountLockout(ctx context.Context, lockout *AccountLockout) error {
	return r.db.WithContext(ctx).Save(lockout).Error
}

func (r *gormSecurityRepository) ExpireAccountLockouts(ctx context.Context) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&AccountLockout{}).
		Where("is_active = ? AND unlocks_at IS NOT NULL AND unlocks_at <= ? AND unlocked_at IS NULL", true, now).
		Updates(map[string]interface{}{
			"is_active":   false,
			"unlocked_at": &now,
		}).Error
}

// Encryption Keys
func (r *gormSecurityRepository) CreateEncryptionKey(ctx context.Context, key *EncryptionKey) error {
	return r.db.WithContext(ctx).Create(key).Error
}

func (r *gormSecurityRepository) GetEncryptionKey(ctx context.Context, keyName string, tenantID *uuid.UUID) (*EncryptionKey, error) {
	var key EncryptionKey
	query := r.db.WithContext(ctx).Where("key_name = ? AND is_active = ?", keyName, true)
	
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}
	
	err := query.Order("key_version DESC").First(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (r *gormSecurityRepository) GetActiveEncryptionKeys(ctx context.Context, tenantID *uuid.UUID) ([]*EncryptionKey, error) {
	var keys []*EncryptionKey
	query := r.db.WithContext(ctx).Where("is_active = ? AND revoked_at IS NULL", true)
	
	if tenantID != nil {
		query = query.Where("tenant_id = ? OR tenant_id IS NULL", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}
	
	err := query.Order("key_name, key_version DESC").Find(&keys).Error
	return keys, err
}

func (r *gormSecurityRepository) UpdateEncryptionKey(ctx context.Context, key *EncryptionKey) error {
	return r.db.WithContext(ctx).Save(key).Error
}

func (r *gormSecurityRepository) RevokeEncryptionKey(ctx context.Context, keyID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&EncryptionKey{}).
		Where("id = ?", keyID).
		Updates(map[string]interface{}{
			"is_active":  false,
			"revoked_at": &now,
		}).Error
}

// Analytics
func (r *gormSecurityRepository) GetSecurityMetrics(ctx context.Context, filter SecurityMetricsFilter) (*SecurityMetrics, error) {
	metrics := &SecurityMetrics{
		Period:          filter.EndTime.Sub(filter.StartTime),
		ThreatsByLevel:  make(map[ThreatLevel]int64),
		AttacksByCountry: make(map[string]int64),
		DeviceTrustDistribution: make(map[DeviceStatus]int64),
	}
	
	// Login attempts metrics
	var loginMetrics struct {
		Total      int64
		Successful int64
		Failed     int64
		Blocked    int64
	}
	
	err := r.db.WithContext(ctx).Model(&LoginAttempt{}).
		Select(`
			COUNT(*) as total,
			SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as successful,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed,
			SUM(CASE WHEN status = 'blocked' THEN 1 ELSE 0 END) as blocked
		`).
		Where("attempted_at >= ? AND attempted_at <= ?", filter.StartTime, filter.EndTime).
		Scan(&loginMetrics).Error
	if err != nil {
		return nil, err
	}
	
	metrics.TotalLoginAttempts = loginMetrics.Total
	metrics.SuccessfulLogins = loginMetrics.Successful
	metrics.FailedLogins = loginMetrics.Failed
	metrics.BlockedAttempts = loginMetrics.Blocked
	
	if loginMetrics.Total > 0 {
		metrics.LoginSuccessRate = float64(loginMetrics.Successful) / float64(loginMetrics.Total)
	}
	
	// Unique users
	err = r.db.WithContext(ctx).Model(&LoginAttempt{}).
		Where("attempted_at >= ? AND attempted_at <= ? AND user_id IS NOT NULL", filter.StartTime, filter.EndTime).
		Distinct("user_id").
		Count(&metrics.UniqueUsers).Error
	if err != nil {
		return nil, err
	}
	
	// Security events by threat level
	var threatLevelCounts []struct {
		ThreatLevel ThreatLevel
		Count       int64
	}
	
	err = r.db.WithContext(ctx).Model(&SecurityEvent{}).
		Select("threat_level, COUNT(*) as count").
		Where("occurred_at >= ? AND occurred_at <= ?", filter.StartTime, filter.EndTime).
		Group("threat_level").
		Scan(&threatLevelCounts).Error
	if err != nil {
		return nil, err
	}
	
	for _, tc := range threatLevelCounts {
		metrics.ThreatsByLevel[tc.ThreatLevel] = tc.Count
	}
	
	return metrics, nil
}
