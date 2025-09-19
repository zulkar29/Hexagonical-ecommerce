package security

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
)

// Integration tests with real database
func TestSecurityIntegrationSecurityLifecycle(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(
		&PasswordPolicy{},
		&LoginAttempt{},
		&TrustedDevice{},
		&SecurityEvent{},
		&PasswordHistory{},
		&AccountLockout{},
		&EncryptionKey{},
	)
	require.NoError(t, err)

	// Setup repository
	repo := NewSecurityRepository(testDB.DB)

	t.Run("Password policy lifecycle", func(t *testing.T) {
		tenantID := uuid.New()

		// Create password policy
		policy := &PasswordPolicy{
			ID:                     uuid.New(),
			TenantID:               &tenantID,
			MinLength:              8,
			MaxLength:              128,
			RequireUppercase:       true,
			RequireLowercase:       true,
			RequireNumbers:         true,
			RequireSpecial:         true,
			PasswordHistoryCount:   5,
			ExpirationDays:         90,
			MaxFailedAttempts:      5,
			LockoutDurationMins:    30,
			PreventCommonPasswords: true,
			PreventUserInfo:        true,
			ForbiddenPatterns:      []string{"password", "123456"},
			IsActive:               true,
		}

		err := repo.CreatePasswordPolicy(context.Background(), policy)
		require.NoError(t, err)

		// Get password policy
		retrievedPolicy, err := repo.GetPasswordPolicy(context.Background(), &tenantID)
		require.NoError(t, err)
		assert.Equal(t, policy.ID, retrievedPolicy.ID)
		assert.Equal(t, policy.MinLength, retrievedPolicy.MinLength)
		assert.Equal(t, policy.MaxFailedAttempts, retrievedPolicy.MaxFailedAttempts)

		// Update password policy
		retrievedPolicy.MinLength = 10
		retrievedPolicy.MaxFailedAttempts = 3
		err = repo.UpdatePasswordPolicy(context.Background(), retrievedPolicy)
		require.NoError(t, err)

		// Verify update
		updatedPolicy, err := repo.GetPasswordPolicy(context.Background(), &tenantID)
		require.NoError(t, err)
		assert.Equal(t, 10, updatedPolicy.MinLength)
		assert.Equal(t, 3, updatedPolicy.MaxFailedAttempts)
	})

	t.Run("Login attempt tracking", func(t *testing.T) {
		userID := uuid.New()
		email := "test@example.com"
		ipAddress := "192.168.1.1"

		// Create successful login attempt
		successAttempt := &LoginAttempt{
			ID:                uuid.New(),
			UserID:            &userID,
			Email:             email,
			Status:            LoginAttemptSuccess,
			IPAddress:         ipAddress,
			UserAgent:         "Mozilla/5.0",
			Country:           "US",
			City:              "New York",
			DeviceFingerprint: "device123",
			ThreatLevel:       ThreatLevelLow,
			AttemptedAt:       time.Now(),
			ProcessedAt:       time.Now(),
		}

		err := repo.CreateLoginAttempt(context.Background(), successAttempt)
		require.NoError(t, err)

		// Create failed login attempt
		failedAttempt := &LoginAttempt{
			ID:                uuid.New(),
			UserID:            &userID,
			Email:             email,
			Status:            LoginAttemptFailed,
			IPAddress:         ipAddress,
			UserAgent:         "Mozilla/5.0",
			DeviceFingerprint: "device123",
			ThreatLevel:       ThreatLevelMedium,
			FailureReason:     "invalid_password",
			AttemptedAt:       time.Now(),
			ProcessedAt:       time.Now(),
		}

		err = repo.CreateLoginAttempt(context.Background(), failedAttempt)
		require.NoError(t, err)

		// Get recent login attempts
		attempts, err := repo.GetRecentLoginAttempts(context.Background(), userID, 10)
		require.NoError(t, err)
		assert.Len(t, attempts, 2)

		// Get failed login count
		failedCount, err := repo.GetFailedLoginCount(context.Background(), email, time.Now().Add(-time.Hour))
		require.NoError(t, err)
		assert.Equal(t, 1, failedCount)
	})

	t.Run("Trusted device management", func(t *testing.T) {
		userID := uuid.New()
		fingerprint := "device_fingerprint_123"

		// Create trusted device
		device := &TrustedDevice{
			ID:          uuid.New(),
			UserID:      userID,
			DeviceID:    "device_id_123",
			Fingerprint: fingerprint,
			Name:        "iPhone 13",
			DeviceType:  "mobile",
			OS:          "iOS 15",
			Browser:     "Safari",
			Status:      DeviceStatusTrusted,
			FirstSeenAt: time.Now(),
			LastSeenAt:  time.Now(),
			IPAddress:   "192.168.1.1",
			UserAgent:   "Mozilla/5.0",
			TrustScore:  1.0,
			Country:     "US",
			City:        "New York",
		}

		err := repo.CreateTrustedDevice(context.Background(), device)
		require.NoError(t, err)

		// Get trusted device
		retrievedDevice, err := repo.GetTrustedDevice(context.Background(), userID, fingerprint)
		require.NoError(t, err)
		assert.Equal(t, device.ID, retrievedDevice.ID)
		assert.Equal(t, device.Name, retrievedDevice.Name)
		assert.Equal(t, device.Status, retrievedDevice.Status)

		// Update device activity
		retrievedDevice.LastSeenAt = time.Now()
		retrievedDevice.AccessCount = 10
		err = repo.UpdateTrustedDevice(context.Background(), retrievedDevice)
		require.NoError(t, err)

		// Get all user devices
		devices, err := repo.GetTrustedDevices(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, devices, 1)
		assert.Equal(t, 10, devices[0].AccessCount)

		// Revoke device
		err = repo.RevokeTrustedDevice(context.Background(), device.ID, "user_requested")
		require.NoError(t, err)

		// Verify device is no longer found in active devices
		_, err = repo.GetTrustedDevice(context.Background(), userID, fingerprint)
		assert.Error(t, err) // Should return error since revoked devices are filtered out

		// Verify device is not in active devices list
		devices, err = repo.GetTrustedDevices(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, devices, 0) // Should be empty since device is revoked
	})

	t.Run("Security event tracking", func(t *testing.T) {
		userID := uuid.New()
		tenantID := uuid.New()

		// Create security event
		event := &SecurityEvent{
			ID:                uuid.New(),
			UserID:            &userID,
			TenantID:          &tenantID,
			EventType:         "login_failed",
			ThreatLevel:       ThreatLevelHigh,
			Description:       "Multiple failed login attempts detected",
			IPAddress:         "192.168.1.1",
			UserAgent:         "Mozilla/5.0",
			DeviceFingerprint: "device123",
			Country:           "US",
			City:              "New York",
			Metadata:          map[string]interface{}{"attempt_count": 5},
			IsResolved:        false,
			OccurredAt:        time.Now(),
		}

		err := repo.CreateSecurityEvent(context.Background(), event)
		require.NoError(t, err)

		// Get unresolved events
		unresolvedEvents, err := repo.GetUnresolvedEvents(context.Background(), ThreatLevelHigh)
		require.NoError(t, err)
		assert.Len(t, unresolvedEvents, 1)
		assert.Equal(t, event.ID, unresolvedEvents[0].ID)

		// Resolve security event
		adminID := uuid.New()
		event.IsResolved = true
		now := time.Now()
		event.ResolvedAt = &now
		event.ResolvedBy = &adminID
		event.ResolutionNotes = "False positive - user verified"

		err = repo.UpdateSecurityEvent(context.Background(), event)
		require.NoError(t, err)
	})

	t.Run("Multi-tenant isolation", func(t *testing.T) {
		// Create data for different tenants
		tenant1 := uuid.New()
		tenant2 := uuid.New()

		// Create password policies for different tenants
		policy1 := &PasswordPolicy{
			ID:        uuid.New(),
			TenantID:  &tenant1,
			MinLength: 8,
			IsActive:  true,
		}

		policy2 := &PasswordPolicy{
			ID:        uuid.New(),
			TenantID:  &tenant2,
			MinLength: 12,
			IsActive:  true,
		}

		err := repo.CreatePasswordPolicy(context.Background(), policy1)
		require.NoError(t, err)

		err = repo.CreatePasswordPolicy(context.Background(), policy2)
		require.NoError(t, err)

		// Verify tenant isolation
		tenant1Policy, err := repo.GetPasswordPolicy(context.Background(), &tenant1)
		require.NoError(t, err)
		assert.Equal(t, policy1.ID, tenant1Policy.ID)
		assert.Equal(t, 8, tenant1Policy.MinLength)

		tenant2Policy, err := repo.GetPasswordPolicy(context.Background(), &tenant2)
		require.NoError(t, err)
		assert.Equal(t, policy2.ID, tenant2Policy.ID)
		assert.Equal(t, 12, tenant2Policy.MinLength)
	})
}
