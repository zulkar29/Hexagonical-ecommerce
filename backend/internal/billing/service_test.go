package billing

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
)

// Integration tests for billing service - critical for SaaS platform

func TestBillingService_BillingPlanManagement(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&BillingPlan{}, &TenantSubscription{}, &UsageTier{})
	require.NoError(t, err)

	// Setup services
	billingRepo := NewRepository(testDB.DB)
	billingService := NewBillingService(billingRepo, nil, nil, nil, nil)

	ctx := context.Background()

	t.Run("Billing plan CRUD operations", func(t *testing.T) {
		// Step 1: Create billing plan
		plan := &BillingPlan{
			ID:              uuid.New(),
			Name:            "Starter Plan",
			Description:     "Perfect for small businesses",
			BasePrice:       1990.00,
			Currency:        "BDT",
			BillingCycle:    BillingCycleMonthly,
			TrialPeriodDays: 14,
			Limits: map[string]interface{}{
				"max_products":   500,
				"max_staff":      1,
				"storage_gb":     5,
				"api_calls_hour": 10000,
			},
			Features: []string{"email_support", "basic_analytics", "ssl_certificate"},
			IsActive: true,
		}

		err := billingService.CreateBillingPlan(ctx, plan)
		require.NoError(t, err)

		// Step 2: Get billing plan
		retrievedPlan, err := billingService.GetBillingPlan(ctx, plan.ID)
		require.NoError(t, err)
		assert.Equal(t, plan.ID, retrievedPlan.ID)
		assert.Equal(t, plan.Name, retrievedPlan.Name)
		assert.Equal(t, plan.BasePrice, retrievedPlan.BasePrice)
		assert.Equal(t, plan.Currency, retrievedPlan.Currency)

		// Step 3: List billing plans
		filter := PlanFilter{
			IsActive: boolPtr(true),
		}
		plans, err := billingService.GetBillingPlans(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, plans, 1)
		assert.Equal(t, plan.ID, plans[0].ID)

		// Step 4: Update billing plan
		plan.BasePrice = 2490.00
		plan.Description = "Updated starter plan for growing businesses"
		err = billingService.UpdateBillingPlan(ctx, plan)
		require.NoError(t, err)

		// Verify update
		updatedPlan, err := billingService.GetBillingPlan(ctx, plan.ID)
		require.NoError(t, err)
		assert.Equal(t, 2490.00, updatedPlan.BasePrice)
		assert.Equal(t, "Updated starter plan for growing businesses", updatedPlan.Description)

		// Step 5: Delete billing plan
		err = billingService.DeleteBillingPlan(ctx, plan.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = billingService.GetBillingPlan(ctx, plan.ID)
		assert.Error(t, err) // Should not find deleted plan
	})
}

func TestBillingService_SubscriptionManagement(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&BillingPlan{}, &TenantSubscription{})
	require.NoError(t, err)

	// Setup services
	billingRepo := NewRepository(testDB.DB)
	billingService := NewBillingService(billingRepo, nil, nil, nil, nil)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Subscription lifecycle", func(t *testing.T) {
		// Step 1: Create billing plan
		plan := &BillingPlan{
			ID:              uuid.New(),
			Name:            "Test Plan",
			BasePrice:       1990.00,
			Currency:        "BDT",
			BillingCycle:    BillingCycleMonthly,
			TrialPeriodDays: 14,
			IsActive:        true,
		}

		err := billingService.CreateBillingPlan(ctx, plan)
		require.NoError(t, err)

		// Step 2: Create subscription
		subscription, err := billingService.CreateSubscription(ctx, tenantID, plan.ID, nil)
		require.NoError(t, err)
		assert.Equal(t, tenantID, subscription.TenantID)
		assert.Equal(t, plan.ID, subscription.PlanID)
		assert.NotEmpty(t, subscription.Status)

		// Step 3: Get subscription
		retrievedSub, err := billingService.GetSubscription(ctx, tenantID)
		require.NoError(t, err)
		assert.Equal(t, subscription.ID, retrievedSub.ID)
		assert.Equal(t, tenantID, retrievedSub.TenantID)

		// Step 4: Update subscription
		updates := SubscriptionUpdate{
			PaymentMethodID: stringPtr("test_payment_method"),
		}

		updatedSub, err := billingService.UpdateSubscription(ctx, tenantID, updates)
		require.NoError(t, err)
		assert.Equal(t, tenantID, updatedSub.TenantID)

		// Step 5: Cancel subscription
		err = billingService.CancelSubscription(ctx, tenantID, "User requested cancellation", false)
		require.NoError(t, err)

		// Verify cancellation
		cancelledSub, err := billingService.GetSubscription(ctx, tenantID)
		require.NoError(t, err)
		assert.Equal(t, SubscriptionStatusCanceled, cancelledSub.Status)
	})
}

func TestBillingService_UsageTracking(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas for usage tracking
	err := testDB.DB.AutoMigrate(&BillingPlan{}, &TenantSubscription{}, &UsageTier{})
	require.NoError(t, err)

	// Setup services
	billingRepo := NewRepository(testDB.DB)
	billingService := NewBillingService(billingRepo, nil, nil, nil, nil)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Usage recording and tracking", func(t *testing.T) {
		// Step 1: Create plan with usage tiers
		plan := &BillingPlan{
			ID:           uuid.New(),
			Name:         "Pro Plan",
			BasePrice:    7990.00,
			Currency:     "BDT",
			BillingCycle: BillingCycleMonthly,
			IsActive:     true,
		}

		err := billingService.CreateBillingPlan(ctx, plan)
		require.NoError(t, err)

		// Step 2: Create usage tier
		usageTier := &UsageTier{
			ID:            uuid.New(),
			BillingPlanID: plan.ID,
			UsageType:     UsageTypeAPIRequests,
			MinUnits:      0,
			MaxUnits:      int64Ptr(200000),
			PricePerUnit:  0.01, // 1 paisa per request
		}

		err = billingService.CreateUsageTier(ctx, usageTier)
		require.NoError(t, err)

		// Step 3: Create subscription
		_, err = billingService.CreateSubscription(ctx, tenantID, plan.ID, nil)
		require.NoError(t, err)

		// Step 4: Record usage
		metadata := map[string]interface{}{
			"source": "api_test",
			"period": time.Now().Format("2006-01"),
		}

		err = billingService.RecordUsage(ctx, tenantID, UsageTypeAPIRequests, 150000, metadata)
		require.NoError(t, err)

		// Step 5: Check usage summary
		startDate := time.Now().AddDate(0, 0, -1)
		endDate := time.Now().AddDate(0, 0, 1)

		usageSummary, err := billingService.GetUsageSummary(ctx, tenantID, startDate, endDate)
		require.NoError(t, err)
		assert.Equal(t, int64(150000), usageSummary[UsageTypeAPIRequests])

		// Step 6: Check usage limits
		limitStatus, err := billingService.CheckUsageLimits(ctx, tenantID)
		require.NoError(t, err)
		assert.Equal(t, tenantID, limitStatus.TenantID)
		assert.Equal(t, int64(150000), limitStatus.CurrentUsage[UsageTypeAPIRequests])
	})
}

func TestBillingService_MultiTenantIsolation(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&BillingPlan{}, &TenantSubscription{})
	require.NoError(t, err)

	// Setup services
	billingRepo := NewRepository(testDB.DB)
	billingService := NewBillingService(billingRepo, nil, nil, nil, nil)

	ctx := context.Background()
	tenant1ID := uuid.New()
	tenant2ID := uuid.New()

	t.Run("Subscription isolation between tenants", func(t *testing.T) {
		// Create plan
		plan := &BillingPlan{
			ID:           uuid.New(),
			Name:         "Test Plan",
			BasePrice:    1990.00,
			Currency:     "BDT",
			BillingCycle: BillingCycleMonthly,
			IsActive:     true,
		}

		err := billingService.CreateBillingPlan(ctx, plan)
		require.NoError(t, err)

		// Create subscriptions for both tenants
		sub1, err := billingService.CreateSubscription(ctx, tenant1ID, plan.ID, nil)
		require.NoError(t, err)

		sub2, err := billingService.CreateSubscription(ctx, tenant2ID, plan.ID, nil)
		require.NoError(t, err)

		// Verify tenant 1 can only see their subscription
		retrievedSub1, err := billingService.GetSubscription(ctx, tenant1ID)
		require.NoError(t, err)
		assert.Equal(t, tenant1ID, retrievedSub1.TenantID)
		assert.Equal(t, sub1.ID, retrievedSub1.ID)

		// Verify tenant 2 can only see their subscription
		retrievedSub2, err := billingService.GetSubscription(ctx, tenant2ID)
		require.NoError(t, err)
		assert.Equal(t, tenant2ID, retrievedSub2.TenantID)
		assert.Equal(t, sub2.ID, retrievedSub2.ID)

		// Verify subscriptions are different
		assert.NotEqual(t, sub1.ID, sub2.ID)
		assert.NotEqual(t, sub1.TenantID, sub2.TenantID)
	})
}

// Helper functions

func timePtr(t time.Time) *time.Time {
	return &t
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func int64Ptr(i int64) *int64 {
	return &i
}