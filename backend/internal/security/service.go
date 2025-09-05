package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// SecurityService defines the interface for security business logic
type SecurityService interface {
	// Password Management
	ValidatePassword(ctx context.Context, password string, userID uuid.UUID, tenantID *uuid.UUID) (*PasswordValidationResult, error)
	HashPassword(password string) (string, error)
	IsPasswordCompromised(password string) (bool, error)
	EnforcePasswordExpiry(ctx context.Context, userID uuid.UUID) (*PasswordPolicy, error)
	StorePasswordHistory(ctx context.Context, userID uuid.UUID, hashedPassword string) error
	
	// Login Security
	RecordLoginAttempt(ctx context.Context, request *LoginAttemptRequest) (*LoginAttempt, error)
	ValidateLoginAttempt(ctx context.Context, userID uuid.UUID, email, ipAddress string) (*LoginValidationResult, error)
	ProcessSuccessfulLogin(ctx context.Context, userID uuid.UUID, ipAddress, userAgent string) error
	ProcessFailedLogin(ctx context.Context, email, ipAddress, userAgent, reason string) error
	
	// Account Lockout Management
	CheckAccountLockout(ctx context.Context, userID uuid.UUID) (*AccountLockoutStatus, error)
	LockAccount(ctx context.Context, userID uuid.UUID, reason string, duration *time.Duration) error
	UnlockAccount(ctx context.Context, userID uuid.UUID, adminID *uuid.UUID) error
	ProcessAutomaticUnlocks(ctx context.Context) error
	
	// Device Management
	RegisterTrustedDevice(ctx context.Context, request *TrustedDeviceRequest) (*TrustedDevice, error)
	ValidateDevice(ctx context.Context, userID uuid.UUID, fingerprint string) (*DeviceValidationResult, error)
	UpdateDeviceActivity(ctx context.Context, userID uuid.UUID, fingerprint, ipAddress string) error
	RevokeTrustedDevice(ctx context.Context, deviceID uuid.UUID, reason string) error
	GetUserDevices(ctx context.Context, userID uuid.UUID) ([]*TrustedDevice, error)
	
	// Security Event Handling
	LogSecurityEvent(ctx context.Context, event *SecurityEventRequest) error
	ProcessSecurityEvents(ctx context.Context, threatLevel ThreatLevel) ([]*SecurityEvent, error)
	ResolveSecurityEvent(ctx context.Context, eventID uuid.UUID, resolution string, adminID *uuid.UUID) error
	
	// Threat Detection
	AnalyzeThreatLevel(ctx context.Context, userID uuid.UUID, ipAddress, userAgent string) (ThreatLevel, error)
	DetectSuspiciousActivity(ctx context.Context, userID uuid.UUID, activity *ActivityContext) (*ThreatAssessment, error)
	ProcessBruteForceDetection(ctx context.Context, email, ipAddress string) error
	
	// Encryption Management
	GetEncryptionKey(ctx context.Context, keyName string, tenantID *uuid.UUID) (*EncryptionKey, error)
	RotateEncryptionKey(ctx context.Context, keyName string, tenantID *uuid.UUID) (*EncryptionKey, error)
	EncryptData(data []byte, keyID uuid.UUID) ([]byte, error)
	DecryptData(encryptedData []byte, keyID uuid.UUID) ([]byte, error)
	
	// Security Analytics
	GetSecurityDashboard(ctx context.Context, tenantID *uuid.UUID, period time.Duration) (*SecurityDashboard, error)
	GetSecurityReport(ctx context.Context, filter SecurityReportFilter) (*SecurityReport, error)
	GetRiskScore(ctx context.Context, userID uuid.UUID) (*RiskScore, error)
}

// Request/Response types
type LoginAttemptRequest struct {
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Email      string     `json:"email"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	Success    bool       `json:"success"`
	FailReason string     `json:"fail_reason,omitempty"`
	Location   *Location  `json:"location,omitempty"`
}

type Location struct {
	Country    string  `json:"country"`
	City       string  `json:"city"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	ISP        string  `json:"isp"`
	Timezone   string  `json:"timezone"`
}

type PasswordValidationResult struct {
	IsValid      bool     `json:"is_valid"`
	Errors       []string `json:"errors"`
	Strength     int      `json:"strength"`
	IsCompromised bool    `json:"is_compromised"`
	Policy       *PasswordPolicy `json:"policy"`
}

type LoginValidationResult struct {
	IsAllowed      bool      `json:"is_allowed"`
	Reason         string    `json:"reason,omitempty"`
	ThreatLevel    ThreatLevel `json:"threat_level"`
	RequiresMFA    bool      `json:"requires_mfa"`
	AccountLocked  bool      `json:"account_locked"`
	UnlocksAt      *time.Time `json:"unlocks_at,omitempty"`
	RemainingAttempts int    `json:"remaining_attempts"`
}

type AccountLockoutStatus struct {
	IsLocked       bool       `json:"is_locked"`
	LockedAt       *time.Time `json:"locked_at,omitempty"`
	UnlocksAt      *time.Time `json:"unlocks_at,omitempty"`
	Reason         string     `json:"reason,omitempty"`
	AttemptCount   int        `json:"attempt_count"`
	MaxAttempts    int        `json:"max_attempts"`
}

type TrustedDeviceRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	Fingerprint string    `json:"fingerprint"`
	DeviceName  string    `json:"device_name"`
	DeviceType  string    `json:"device_type"`
	OS          string    `json:"os"`
	Browser     string    `json:"browser"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Location    *Location `json:"location,omitempty"`
}

type DeviceValidationResult struct {
	IsTrusted    bool       `json:"is_trusted"`
	TrustScore   float64    `json:"trust_score"`
	Device       *TrustedDevice `json:"device,omitempty"`
	RequiresMFA  bool       `json:"requires_mfa"`
	IsBlocked    bool       `json:"is_blocked"`
	LastSeen     *time.Time `json:"last_seen,omitempty"`
}

type SecurityEventRequest struct {
	UserID      *uuid.UUID  `json:"user_id,omitempty"`
	TenantID    *uuid.UUID  `json:"tenant_id,omitempty"`
	EventType   string      `json:"event_type"`
	Description string      `json:"description"`
	ThreatLevel ThreatLevel `json:"threat_level"`
	IPAddress   string      `json:"ip_address"`
	UserAgent   string      `json:"user_agent"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Location    *Location   `json:"location,omitempty"`
}

type ActivityContext struct {
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	Timestamp  time.Time              `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type ThreatAssessment struct {
	ThreatLevel ThreatLevel            `json:"threat_level"`
	Score       float64                `json:"score"`
	Factors     []string               `json:"factors"`
	Actions     []string               `json:"actions"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type SecurityDashboard struct {
	Period           time.Duration     `json:"period"`
	TotalLogins      int64             `json:"total_logins"`
	FailedLogins     int64             `json:"failed_logins"`
	BlockedAttempts  int64             `json:"blocked_attempts"`
	ActiveLockouts   int64             `json:"active_lockouts"`
	TrustedDevices   int64             `json:"trusted_devices"`
	SecurityEvents   int64             `json:"security_events"`
	HighThreatEvents int64             `json:"high_threat_events"`
	RiskScore        float64           `json:"risk_score"`
	TopThreats       []ThreatSummary   `json:"top_threats"`
	RecentEvents     []*SecurityEvent  `json:"recent_events"`
}

type ThreatSummary struct {
	Type        string  `json:"type"`
	Count       int64   `json:"count"`
	Severity    string  `json:"severity"`
	TrendChange float64 `json:"trend_change"`
}

type SecurityReportFilter struct {
	TenantID    *uuid.UUID `json:"tenant_id,omitempty"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	EventTypes  []string   `json:"event_types,omitempty"`
	ThreatLevel *ThreatLevel `json:"threat_level,omitempty"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
}

type SecurityReport struct {
	Filter      SecurityReportFilter `json:"filter"`
	Summary     *SecurityMetrics     `json:"summary"`
	Events      []*SecurityEvent     `json:"events"`
	Trends      []TrendData          `json:"trends"`
	Recommendations []string         `json:"recommendations"`
	GeneratedAt time.Time            `json:"generated_at"`
}

type TrendData struct {
	Period string  `json:"period"`
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type RiskScore struct {
	UserID          uuid.UUID            `json:"user_id"`
	OverallScore    float64              `json:"overall_score"`
	Level           string               `json:"level"`
	Factors         []RiskFactor         `json:"factors"`
	Recommendations []string             `json:"recommendations"`
	LastUpdated     time.Time            `json:"last_updated"`
}

type RiskFactor struct {
	Category    string  `json:"category"`
	Score       float64 `json:"score"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
}

// securityService implements SecurityService
type securityService struct {
	repo SecurityRepository
}

// NewSecurityService creates a new security service
func NewSecurityService(repo SecurityRepository) SecurityService {
	return &securityService{
		repo: repo,
	}
}

// Password Management
func (s *securityService) ValidatePassword(ctx context.Context, password string, userID uuid.UUID, tenantID *uuid.UUID) (*PasswordValidationResult, error) {
	result := &PasswordValidationResult{
		IsValid: true,
		Errors:  []string{},
	}
	
	// Get password policy
	policy, err := s.repo.GetPasswordPolicy(ctx, tenantID)
	if err != nil {
		// Use default policy if none found
		policy = &PasswordPolicy{
			MinLength:        8,
			RequireUppercase: true,
			RequireLowercase: true,
			RequireNumbers:   true,
			RequireSymbols:   true,
			MaxAge:          90,
			HistoryCount:    5,
		}
	}
	result.Policy = policy
	
	// Validate length
	if len(password) < policy.MinLength {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Password must be at least %d characters long", policy.MinLength))
	}
	
	if policy.MaxLength > 0 && len(password) > policy.MaxLength {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Password must be no more than %d characters long", policy.MaxLength))
	}
	
	// Character requirements
	if policy.RequireUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password must contain at least one uppercase letter")
	}
	
	if policy.RequireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password must contain at least one lowercase letter")
	}
	
	if policy.RequireNumbers && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password must contain at least one number")
	}
	
	if policy.RequireSymbols && !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password must contain at least one special character")
	}
	
	// Check forbidden patterns
	for _, pattern := range policy.ForbiddenPatterns {
		if matched, _ := regexp.MatchString(pattern, strings.ToLower(password)); matched {
			result.IsValid = false
			result.Errors = append(result.Errors, "Password contains forbidden pattern")
			break
		}
	}
	
	// Calculate strength
	result.Strength = s.calculatePasswordStrength(password)
	
	// Check if password is compromised
	compromised, err := s.IsPasswordCompromised(password)
	if err == nil && compromised {
		result.IsCompromised = true
		result.IsValid = false
		result.Errors = append(result.Errors, "Password has been found in data breaches")
	}
	
	// Check password history
	if userID != uuid.Nil {
		history, err := s.repo.GetPasswordHistory(ctx, userID, policy.HistoryCount)
		if err == nil {
			for _, h := range history {
				if bcrypt.CompareHashAndPassword([]byte(h.PasswordHash), []byte(password)) == nil {
					result.IsValid = false
					result.Errors = append(result.Errors, "Password has been used recently")
					break
				}
			}
		}
	}
	
	return result, nil
}

func (s *securityService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *securityService) IsPasswordCompromised(password string) (bool, error) {
	// Simple SHA-1 hash for HaveIBeenPwned API (not implemented here)
	// In real implementation, you would check against a compromised password database
	return false, nil
}

func (s *securityService) EnforcePasswordExpiry(ctx context.Context, userID uuid.UUID) (*PasswordPolicy, error) {
	policy, err := s.repo.GetPasswordPolicy(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	if policy.MaxAge <= 0 {
		return policy, nil // No expiry enforced
	}
	
	// Get latest password from history
	history, err := s.repo.GetPasswordHistory(ctx, userID, 1)
	if err != nil || len(history) == 0 {
		return policy, nil
	}
	
	latestPassword := history[0]
	expiryDate := latestPassword.CreatedAt.AddDate(0, 0, policy.MaxAge)
	
	if time.Now().After(expiryDate) {
		// Password has expired - this would trigger a password reset flow
		return policy, fmt.Errorf("password has expired")
	}
	
	return policy, nil
}

func (s *securityService) StorePasswordHistory(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	history := &PasswordHistory{
		ID:           uuid.New(),
		UserID:       userID,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}
	
	err := s.repo.CreatePasswordHistory(ctx, history)
	if err != nil {
		return err
	}
	
	// Cleanup old history
	return s.repo.CleanupOldPasswordHistory(ctx, userID, 10)
}

// Login Security
func (s *securityService) RecordLoginAttempt(ctx context.Context, request *LoginAttemptRequest) (*LoginAttempt, error) {
	var userID uuid.UUID
	if request.UserID != nil {
		userID = *request.UserID
	}
	threatLevel, _ := s.AnalyzeThreatLevel(ctx, userID, request.IPAddress, request.UserAgent)
	
	attempt := &LoginAttempt{
		ID:          uuid.New(),
		UserID:      request.UserID,
		Email:       request.Email,
		IPAddress:   request.IPAddress,
		UserAgent:   request.UserAgent,
		ThreatLevel: threatLevel,
		AttemptedAt: time.Now(),
	}
	
	if request.Success {
		attempt.Status = LoginAttemptSuccess
	} else {
		attempt.Status = LoginAttemptFailed
		attempt.FailureReason = request.FailReason
	}
	
	if request.Location != nil {
		attempt.Country = request.Location.Country
		attempt.City = request.Location.City
	}
	
	err := s.repo.CreateLoginAttempt(ctx, attempt)
	if err != nil {
		return nil, err
	}
	
	// Process failed login for brute force detection
	if !request.Success {
		_ = s.ProcessBruteForceDetection(ctx, request.Email, request.IPAddress)
	}
	
	return attempt, nil
}

func (s *securityService) ValidateLoginAttempt(ctx context.Context, userID uuid.UUID, email, ipAddress string) (*LoginValidationResult, error) {
	result := &LoginValidationResult{
		IsAllowed: true,
	}
	
	// Check account lockout
	lockoutStatus, err := s.CheckAccountLockout(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	if lockoutStatus.IsLocked {
		result.IsAllowed = false
		result.AccountLocked = true
		result.Reason = lockoutStatus.Reason
		result.UnlocksAt = lockoutStatus.UnlocksAt
		return result, nil
	}
	
	// Analyze threat level
	threatLevel, err := s.AnalyzeThreatLevel(ctx, userID, ipAddress, "")
	if err != nil {
		threatLevel = ThreatLevelLow
	}
	result.ThreatLevel = threatLevel
	
	// Check recent failed attempts
	failedCount, err := s.repo.GetFailedLoginCount(ctx, email, time.Now().Add(-15*time.Minute))
	if err == nil {
		maxAttempts := 5 // This could come from policy
		result.RemainingAttempts = maxAttempts - failedCount
		
		if failedCount >= maxAttempts {
			result.IsAllowed = false
			result.Reason = "Too many failed login attempts"
		}
	}
	
	// Require MFA for high threat levels
	if threatLevel >= ThreatLevelHigh {
		result.RequiresMFA = true
	}
	
	return result, nil
}

func (s *securityService) ProcessSuccessfulLogin(ctx context.Context, userID uuid.UUID, ipAddress, userAgent string) error {
	// Update device activity if it's a trusted device
	_ = s.UpdateDeviceActivity(ctx, userID, s.generateDeviceFingerprint(userAgent, ipAddress), ipAddress)
	
	// Log security event for successful login
	return s.LogSecurityEvent(ctx, &SecurityEventRequest{
		UserID:      &userID,
		EventType:   "login_success",
		Description: "User successfully logged in",
		ThreatLevel: ThreatLevelLow,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	})
}

func (s *securityService) ProcessFailedLogin(ctx context.Context, email, ipAddress, userAgent, reason string) error {
	// Log security event for failed login
	return s.LogSecurityEvent(ctx, &SecurityEventRequest{
		EventType:   "login_failed",
		Description: fmt.Sprintf("Failed login attempt: %s", reason),
		ThreatLevel: ThreatLevelMedium,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Metadata: map[string]interface{}{
			"email":  email,
			"reason": reason,
		},
	})
}

// Account Lockout Management
func (s *securityService) CheckAccountLockout(ctx context.Context, userID uuid.UUID) (*AccountLockoutStatus, error) {
	status := &AccountLockoutStatus{}
	
	lockout, err := s.repo.GetActiveAccountLockout(ctx, userID)
	if err != nil {
		return status, nil // No active lockout
	}
	
	status.IsLocked = true
	status.LockedAt = &lockout.LockedAt
	status.UnlocksAt = lockout.UnlocksAt
	status.Reason = lockout.Reason
	
	return status, nil
}

func (s *securityService) LockAccount(ctx context.Context, userID uuid.UUID, reason string, duration *time.Duration) error {
	lockout := &AccountLockout{
		ID:       uuid.New(),
		UserID:   userID,
		Reason:   reason,
		LockedAt: time.Now(),
		IsActive: true,
	}
	
	if duration != nil {
		unlocksAt := time.Now().Add(*duration)
		lockout.UnlocksAt = &unlocksAt
	}
	
	return s.repo.CreateAccountLockout(ctx, lockout)
}

func (s *securityService) UnlockAccount(ctx context.Context, userID uuid.UUID, adminID *uuid.UUID) error {
	lockout, err := s.repo.GetActiveAccountLockout(ctx, userID)
	if err != nil {
		return err
	}
	
	now := time.Now()
	lockout.IsActive = false
	lockout.UnlockedAt = &now
	lockout.UnlockedBy = adminID
	
	return s.repo.UpdateAccountLockout(ctx, lockout)
}

func (s *securityService) ProcessAutomaticUnlocks(ctx context.Context) error {
	return s.repo.ExpireAccountLockouts(ctx)
}

// Device Management
func (s *securityService) RegisterTrustedDevice(ctx context.Context, request *TrustedDeviceRequest) (*TrustedDevice, error) {
	device := &TrustedDevice{
		ID:           uuid.New(),
		UserID:       request.UserID,
		Fingerprint:  request.Fingerprint,
		DeviceName:   request.DeviceName,
		DeviceType:   request.DeviceType,
		OS:           request.OS,
		Browser:      request.Browser,
		IPAddress:    request.IPAddress,
		UserAgent:    request.UserAgent,
		Status:       DeviceStatusTrusted,
		TrustScore:   0.8, // Initial trust score
		FirstSeenAt:  time.Now(),
		LastSeenAt:   time.Now(),
	}
	
	if request.Location != nil {
		device.Country = request.Location.Country
		device.City = request.Location.City
	}
	
	err := s.repo.CreateTrustedDevice(ctx, device)
	if err != nil {
		return nil, err
	}
	
	return device, nil
}

func (s *securityService) ValidateDevice(ctx context.Context, userID uuid.UUID, fingerprint string) (*DeviceValidationResult, error) {
	result := &DeviceValidationResult{
		IsTrusted: false,
	}
	
	device, err := s.repo.GetTrustedDevice(ctx, userID, fingerprint)
	if err != nil {
		// Device not found - not trusted
		return result, nil
	}
	
	result.Device = device
	result.LastSeen = &device.LastSeenAt
	
	// Check device status
	if device.Status == DeviceStatusBlocked {
		result.IsBlocked = true
		return result, nil
	}
	
	if device.Status == DeviceStatusTrusted {
		result.IsTrusted = true
		result.TrustScore = device.TrustScore
		
		// Require MFA if trust score is low
		if device.TrustScore < 0.5 {
			result.RequiresMFA = true
		}
	}
	
	return result, nil
}

func (s *securityService) UpdateDeviceActivity(ctx context.Context, userID uuid.UUID, fingerprint, ipAddress string) error {
	device, err := s.repo.GetTrustedDevice(ctx, userID, fingerprint)
	if err != nil {
		return err
	}
	
	device.LastSeenAt = time.Now()
	device.IPAddress = ipAddress
	device.AccessCount++
	
	// Update trust score based on consistent usage
	if device.AccessCount > 10 && device.TrustScore < 1.0 {
		device.TrustScore = min(1.0, device.TrustScore+0.1)
	}
	
	return s.repo.UpdateTrustedDevice(ctx, device)
}

func (s *securityService) RevokeTrustedDevice(ctx context.Context, deviceID uuid.UUID, reason string) error {
	return s.repo.RevokeTrustedDevice(ctx, deviceID, reason)
}

func (s *securityService) GetUserDevices(ctx context.Context, userID uuid.UUID) ([]*TrustedDevice, error) {
	return s.repo.GetTrustedDevices(ctx, userID)
}

// Security Event Handling
func (s *securityService) LogSecurityEvent(ctx context.Context, request *SecurityEventRequest) error {
	event := &SecurityEvent{
		ID:          uuid.New(),
		UserID:      request.UserID,
		TenantID:    request.TenantID,
		EventType:   request.EventType,
		Description: request.Description,
		ThreatLevel: request.ThreatLevel,
		IPAddress:   request.IPAddress,
		UserAgent:   request.UserAgent,
		Metadata:    request.Metadata,
		OccurredAt:  time.Now(),
		IsResolved:  false,
	}
	
	if request.Location != nil {
		event.Country = request.Location.Country
		event.City = request.Location.City
	}
	
	return s.repo.CreateSecurityEvent(ctx, event)
}

func (s *securityService) ProcessSecurityEvents(ctx context.Context, threatLevel ThreatLevel) ([]*SecurityEvent, error) {
	return s.repo.GetUnresolvedEvents(ctx, threatLevel)
}

func (s *securityService) ResolveSecurityEvent(ctx context.Context, eventID uuid.UUID, resolution string, adminID *uuid.UUID) error {
	// Implementation would fetch the event and update it
	return nil
}

// Threat Detection
func (s *securityService) AnalyzeThreatLevel(ctx context.Context, userID uuid.UUID, ipAddress, userAgent string) (ThreatLevel, error) {
	score := 0.0
	
	// Check IP reputation (simplified)
	if s.isPrivateIP(ipAddress) {
		score -= 0.2 // Local IPs are generally safer
	} else {
		score += 0.1 // External IPs have slight risk
	}
	
	// Check for recent failed attempts from this IP
	failedCount, err := s.repo.GetFailedLoginCount(ctx, "", time.Now().Add(-1*time.Hour))
	if err == nil && failedCount > 5 {
		score += 0.3
	}
	
	// Determine threat level based on score
	if score >= 0.7 {
		return ThreatLevelCritical, nil
	} else if score >= 0.5 {
		return ThreatLevelHigh, nil
	} else if score >= 0.3 {
		return ThreatLevelMedium, nil
	}
	
	return ThreatLevelLow, nil
}

func (s *securityService) DetectSuspiciousActivity(ctx context.Context, userID uuid.UUID, activity *ActivityContext) (*ThreatAssessment, error) {
	assessment := &ThreatAssessment{
		ThreatLevel: ThreatLevelLow,
		Score:       0.0,
		Factors:     []string{},
		Actions:     []string{},
		Metadata:    make(map[string]interface{}),
	}
	
	// Analyze various factors
	// Geographic anomaly
	// Time-based anomaly
	// Behavioral anomaly
	// Device anomaly
	
	return assessment, nil
}

func (s *securityService) ProcessBruteForceDetection(ctx context.Context, email, ipAddress string) error {
	// Check failed attempts in the last hour
	failedCount, err := s.repo.GetFailedLoginCount(ctx, email, time.Now().Add(-1*time.Hour))
	if err != nil {
		return err
	}
	
	if failedCount >= 10 {
		// Log high-threat security event
		return s.LogSecurityEvent(ctx, &SecurityEventRequest{
			EventType:   "brute_force_detected",
			Description: fmt.Sprintf("Brute force attack detected from IP %s on email %s", ipAddress, email),
			ThreatLevel: ThreatLevelHigh,
			IPAddress:   ipAddress,
			Metadata: map[string]interface{}{
				"email":         email,
				"attempt_count": failedCount,
			},
		})
	}
	
	return nil
}

// Encryption Management
func (s *securityService) GetEncryptionKey(ctx context.Context, keyName string, tenantID *uuid.UUID) (*EncryptionKey, error) {
	return s.repo.GetEncryptionKey(ctx, keyName, tenantID)
}

func (s *securityService) RotateEncryptionKey(ctx context.Context, keyName string, tenantID *uuid.UUID) (*EncryptionKey, error) {
	// Implementation would create new key version
	return nil, nil
}

func (s *securityService) EncryptData(data []byte, keyID uuid.UUID) ([]byte, error) {
	// Implementation would use the encryption key to encrypt data
	return nil, nil
}

func (s *securityService) DecryptData(encryptedData []byte, keyID uuid.UUID) ([]byte, error) {
	// Implementation would use the encryption key to decrypt data
	return nil, nil
}

// Security Analytics
func (s *securityService) GetSecurityDashboard(ctx context.Context, tenantID *uuid.UUID, period time.Duration) (*SecurityDashboard, error) {
	// Implementation would aggregate security metrics
	return nil, nil
}

func (s *securityService) GetSecurityReport(ctx context.Context, filter SecurityReportFilter) (*SecurityReport, error) {
	// Implementation would generate comprehensive security report
	return nil, nil
}

func (s *securityService) GetRiskScore(ctx context.Context, userID uuid.UUID) (*RiskScore, error) {
	// Implementation would calculate user risk score
	return nil, nil
}

// Helper functions
func (s *securityService) calculatePasswordStrength(password string) int {
	score := 0
	
	if len(password) >= 8 {
		score += 25
	}
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		score += 25
	}
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		score += 25
	}
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		score += 25
	}
	if regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
		score += 25
	}
	
	// Bonus for length
	if len(password) >= 12 {
		score += 10
	}
	if len(password) >= 16 {
		score += 10
	}
	
	// Penalty for common patterns
	if regexp.MustCompile(`(.)\1{2,}`).MatchString(password) {
		score -= 10
	}
	if regexp.MustCompile(`(012|123|234|345|456|567|678|789|890|abc|bcd|cde)`).MatchString(strings.ToLower(password)) {
		score -= 10
	}
	
	return max(0, min(100.0, float64(score)))
}

func (s *securityService) generateDeviceFingerprint(userAgent, ipAddress string) string {
	hash := sha256.Sum256([]byte(userAgent + ipAddress))
	return hex.EncodeToString(hash[:])
}

func (s *securityService) isPrivateIP(ipAddress string) bool {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false
	}
	
	return ip.IsPrivate() || ip.IsLoopback()
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a float64, b float64) int {
	if a > b {
		return int(a)
	}
	return int(b)
}
