package referral

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
	"ecommerce-saas/internal/tenant"
	"ecommerce-saas/internal/user"
)

func TestReferralIntegration(t *testing.T) {
	// Setup test database without migrations
	testDB := testhelpers.SetupTestDatabaseWithOptions(t, testhelpers.TestDatabaseOptions{
		RunMigrations: false,
	})
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Create test tenant
	testTenant := &tenant.Tenant{
		ID:        uuid.New(),
		Name:      "Test Tenant",
		Subdomain: "test-referral",
		Status:    "active",
		Plan:      "professional",
	}
	err := testDB.DB.Create(testTenant).Error
	require.NoError(t, err)

	// Create test users
	referrer := &user.User{
		ID:        uuid.New(),
		TenantID:  &testTenant.ID,
		Email:     "referrer@test.com",
		FirstName: "Test",
		LastName:  "Referrer",
		Password:  "password123",
		Status:    "active",
	}
	err = testDB.DB.Create(referrer).Error
	require.NoError(t, err)

	referee := &user.User{
		ID:        uuid.New(),
		TenantID:  &testTenant.ID,
		Email:     "referee@test.com",
		FirstName: "Test",
		LastName:  "Referee",
		Password:  "password123",
		Status:    "active",
	}
	err = testDB.DB.Create(referee).Error
	require.NoError(t, err)

	// Setup referral module
	referralModule := NewModule(testDB.DB)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add middleware to set tenant_id and user_id in context
	router.Use(func(c *gin.Context) {
		c.Set("tenant_id", testTenant.ID)
		c.Set("user_id", referrer.ID)
		c.Set("is_admin", true)
		c.Next()
	})
	
	api := router.Group("/api")
	referralModule.RegisterRoutes(api)

	t.Run("Generate Referral Code", func(t *testing.T) {
		expiresAt := time.Now().Add(30 * 24 * time.Hour)
		req := GenerateReferralCodeRequest{
			CommissionRate: 0.15, // 15%
			ExpiresAt:      &expiresAt,
		}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/referrals/generate", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-Tenant-ID", testTenant.ID.String())

		router.ServeHTTP(w, httpReq)

		// Debug: print response body if test fails
		if w.Code != http.StatusCreated {
			t.Logf("Expected status %d, got %d. Response body: %s", http.StatusCreated, w.Code, w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)

		var response GenerateReferralResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.ReferralCode)
		assert.Equal(t, referrer.ID, response.ReferrerID)
		assert.Equal(t, 0.15, response.CommissionRate)
		assert.NotNil(t, response.ExpiresAt)
	})

	t.Run("Apply Referral Code", func(t *testing.T) {
		// First generate a referral code
		referralCode := "TEST123"
		referral := &Referral{
			ID:             uuid.New(),
			TenantID:       testTenant.ID,
			ReferrerID:     referrer.ID,
			ReferralCode:   referralCode,
			Status:         ReferralStatusActive,
			CommissionRate: 0.10,
			ExpiresAt:      &[]time.Time{time.Now().Add(30 * 24 * time.Hour)}[0],
		}
		err := testDB.DB.Create(referral).Error
		require.NoError(t, err)

		req := ApplyReferralCodeRequest{
			ReferralCode: referralCode,
			RefereeID:    referee.ID,
		}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/referrals/apply", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-Tenant-ID", testTenant.ID.String())

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response ApplyReferralResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Equal(t, referralCode, response.ReferralCode)
		assert.Equal(t, 0.10, response.CommissionRate)

		// Verify referral was updated
		var updatedReferral Referral
		err = testDB.DB.Where("id = ?", referral.ID).First(&updatedReferral).Error
		require.NoError(t, err)
		assert.Equal(t, ReferralStatusCompleted, updatedReferral.Status)
		assert.Equal(t, &referee.ID, updatedReferral.RefereeID)
	})

	t.Run("Get Referral Stats", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/referrals/stats?user_id=%s", referrer.ID), nil)
		httpReq.Header.Set("X-Tenant-ID", testTenant.ID.String())

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response ReferralStatsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Based on the actual referrals created in previous tests:
		// - 1 referral from "Generate Referral Code" test
		// - 1 referral from "Apply Referral Code" test (completed)
		assert.Equal(t, 2, response.TotalReferrals)
		assert.Equal(t, 1, response.SuccessfulReferrals)
		assert.Equal(t, 1, response.PendingReferrals)
		// No commissions created yet in stats test
		assert.Equal(t, 0.0, response.TotalCommissionEarned)
		assert.Equal(t, 0.0, response.TotalCommissionPaid)
		assert.Equal(t, 0.0, response.TotalCommissionPending)
	})

	t.Run("Get Commission History", func(t *testing.T) {
		// Create test referral first
		testReferral := &Referral{
			ID:             uuid.New(),
			TenantID:       testTenant.ID,
			ReferrerID:     referrer.ID,
			ReferralCode:   "COMMISSION123",
			Status:         ReferralStatusActive,
			CommissionRate: 0.10,
			ExpiresAt:      &[]time.Time{time.Now().Add(24 * time.Hour)}[0],
		}
		err := testDB.DB.Create(testReferral).Error
		require.NoError(t, err)

		// Create test commission
		commissionID := uuid.New()
		subscriptionID := uuid.New()
		commission := &ReferralCommission{
			ID:                 commissionID,
			TenantID:           testTenant.ID,
			ReferralID:         testReferral.ID,
			ReferrerID:         referrer.ID,
			SubscriptionID:     &subscriptionID,
			Amount:             25.00,
			Currency:           "USD",
			CommissionRate:     0.10,
			SubscriptionAmount: 250.00,
			Status:             CommissionStatusPaid,
			PaidAt:             &time.Time{},
		}
		now := time.Now()
		commission.PaidAt = &now
		err = testDB.DB.Create(commission).Error
		require.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/api/commissions/my", nil)
		httpReq.Header.Set("X-Tenant-ID", testTenant.ID.String())

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response CommissionHistoryResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Commissions, 1)
		assert.Equal(t, 25.00, response.Commissions[0].Amount)
		assert.Equal(t, "USD", response.Commissions[0].Currency)
		assert.Equal(t, 0.10, response.Commissions[0].CommissionRate)
		assert.Equal(t, 250.00, response.Commissions[0].SubscriptionAmount)
		assert.Equal(t, CommissionStatusPaid, response.Commissions[0].Status)
	})

	t.Run("Invalid Referral Code", func(t *testing.T) {
		req := ApplyReferralCodeRequest{
			ReferralCode: "INVALID123",
			RefereeID:    referee.ID,
		}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/referrals/apply", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-Tenant-ID", testTenant.ID.String())

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Expired Referral Code", func(t *testing.T) {
		// Create expired referral
		expiredCode := "EXPIRED123"
		expiredReferral := &Referral{
			ID:             uuid.New(),
			TenantID:       testTenant.ID,
			ReferrerID:     referrer.ID,
			ReferralCode:   expiredCode,
			Status:         ReferralStatusActive,
			CommissionRate: 0.10,
			ExpiresAt:      &[]time.Time{time.Now().Add(-24 * time.Hour)}[0], // Expired yesterday
		}
		err := testDB.DB.Create(expiredReferral).Error
		require.NoError(t, err)

		req := ApplyReferralCodeRequest{
			ReferralCode: expiredCode,
			RefereeID:    referee.ID,
		}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/referrals/apply", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-Tenant-ID", testTenant.ID.String())

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestReferralService(t *testing.T) {
	// Setup test database without migrations
	testDB := testhelpers.SetupTestDatabaseWithOptions(t, testhelpers.TestDatabaseOptions{
		RunMigrations: false,
	})
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Create test data
	testTenant := &tenant.Tenant{
		ID:        uuid.New(),
		Name:      "Test Tenant",
		Subdomain: "test-service",
		Status:    "active",
	}
	err := testDB.DB.Create(testTenant).Error
	require.NoError(t, err)

	referrer := &user.User{
		ID:        uuid.New(),
		TenantID:  &testTenant.ID,
		Email:     "referrer@test.com",
		FirstName: "Test",
		LastName:  "Referrer",
		Password:  "password123",
		Status:    "active",
	}
	err = testDB.DB.Create(referrer).Error
	require.NoError(t, err)

	// Setup service
	repo := NewGormRepository(testDB.DB)
	service := NewService(repo)

	t.Run("Generate Referral Code", func(t *testing.T) {
		ctx := context.Background()
		expiresAt := time.Now().Add(30 * 24 * time.Hour)
		referral, err := service.GenerateReferralCode(ctx, testTenant.ID, referrer.ID, 0.15, &expiresAt)
		require.NoError(t, err)
		assert.NotNil(t, referral)
		assert.Equal(t, testTenant.ID, referral.TenantID)
		assert.Equal(t, referrer.ID, referral.ReferrerID)
		assert.Equal(t, 0.15, referral.CommissionRate)
		assert.Equal(t, ReferralStatusActive, referral.Status)
		assert.NotEmpty(t, referral.ReferralCode)
		assert.True(t, referral.ExpiresAt.After(time.Now()))
	})

	t.Run("Calculate Commission", func(t *testing.T) {
		// Test commission calculation through referral
		testReferral := &Referral{
			CommissionRate: 0.15,
		}
		commissionAmount := testReferral.CalculateCommission(100.0)
		assert.Equal(t, 15.0, commissionAmount) // 100 * 0.15
	})

	t.Run("Validate Referral Code", func(t *testing.T) {
		// Create a valid referral
		referral := &Referral{
			ID:             uuid.New(),
			TenantID:       testTenant.ID,
			ReferrerID:     referrer.ID,
			ReferralCode:   "VALID123",
			Status:         ReferralStatusActive,
			CommissionRate: 0.10,
			ExpiresAt:      &[]time.Time{time.Now().Add(24 * time.Hour)}[0],
		}
		err := testDB.DB.Create(referral).Error
		require.NoError(t, err)

		// Test valid code
		ctx := context.Background()
		foundReferral, err := service.ValidateReferralCode(ctx, testTenant.ID, "VALID123")
		require.NoError(t, err)
		assert.Equal(t, referral.ID, foundReferral.ID)

		// Test invalid code
		_, err = service.ValidateReferralCode(ctx, testTenant.ID, "INVALID123")
		assert.Error(t, err)
	})
}