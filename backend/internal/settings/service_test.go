package settings

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
)

// Integration tests for settings service - critical for tenant customization

func TestSettingsService_SettingsManagement(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	settingsRepo := NewRepository(testDB.DB)
	settingsService := NewService(settingsRepo)

	ctx := context.Background()
	// Use the fixed test tenant ID from fixtures
	tenantID, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	
	// Seed test data to create the tenant
	testhelpers.SeedMinimalTestData(t, testDB.DB)

	t.Run("General settings CRUD operations", func(t *testing.T) {
		// Step 1: Update general settings
		updateReq := &UpdateSettingsRequest{
			Settings: map[string]interface{}{
				"general": map[string]interface{}{
					"store_name":        "Test Store",
					"store_description": "A test store for e-commerce",
					"store_email":       "test@store.com",
					"currency":          "BDT",
					"timezone":          "Asia/Dhaka",
					"language":          "en",
				},
			},
		}

		updateResp, err := settingsService.UpdateSettings(ctx, tenantID, updateReq)
		require.NoError(t, err)
		assert.NotNil(t, updateResp)
		assert.NotEmpty(t, updateResp.Message)
		assert.Contains(t, updateResp.Updated, "general")

		// Step 2: Get general settings
		getReq := &GetSettingsRequest{
			Section: "general",
		}

		getResp, err := settingsService.GetSettings(ctx, tenantID, getReq)
		require.NoError(t, err)
		assert.NotNil(t, getResp)
		assert.Contains(t, getResp.Settings, "general")

		generalSettings := getResp.Settings["general"].(map[string]interface{})
		assert.Equal(t, "Test Store", generalSettings["store_name"])
		assert.Equal(t, "A test store for e-commerce", generalSettings["store_description"])
		assert.Equal(t, "test@store.com", generalSettings["store_email"])
		assert.Equal(t, "BDT", generalSettings["currency"])
	})

	t.Run("SEO settings management", func(t *testing.T) {
		// Update SEO settings
		seoUpdateReq := &UpdateSettingsRequest{
			Settings: map[string]interface{}{
				"seo": map[string]interface{}{
					"meta_title":         "Test Store - Best Products",
					"meta_description":   "Find the best products at Test Store",
					"meta_keywords":      "test, store, products, ecommerce",
					"google_analytics":   "GA-123456789",
					"facebook_pixel":     "FB-987654321",
					"google_tag_manager": "GTM-ABCDEFG",
				},
			},
		}

		seoUpdateResp, err := settingsService.UpdateSettings(ctx, tenantID, seoUpdateReq)
		require.NoError(t, err)
		assert.NotNil(t, seoUpdateResp)
		assert.Contains(t, seoUpdateResp.Updated, "seo")

		// Get SEO settings
		seoGetReq := &GetSettingsRequest{
			Section: "seo",
		}

		seoGetResp, err := settingsService.GetSettings(ctx, tenantID, seoGetReq)
		require.NoError(t, err)
		assert.Contains(t, seoGetResp.Settings, "seo")

		seoSettings := seoGetResp.Settings["seo"].(map[string]interface{})
		assert.Equal(t, "Test Store - Best Products", seoSettings["meta_title"])
		assert.Equal(t, "Find the best products at Test Store", seoSettings["meta_description"])
		assert.Equal(t, "GA-123456789", seoSettings["google_analytics"])
	})

	t.Run("Appearance settings management", func(t *testing.T) {
		// Update appearance settings
		appearanceUpdateReq := &UpdateSettingsRequest{
			Settings: map[string]interface{}{
				"appearance": map[string]interface{}{
					"logo":            "/uploads/logo.png",
					"favicon":         "/uploads/favicon.ico",
					"primary_color":   "#FF6B35",
					"secondary_color": "#004E89",
					"font_family":     "Roboto",
					"theme_mode":      "dark",
				},
			},
		}

		appearanceUpdateResp, err := settingsService.UpdateSettings(ctx, tenantID, appearanceUpdateReq)
		require.NoError(t, err)
		assert.NotNil(t, appearanceUpdateResp)
		assert.Contains(t, appearanceUpdateResp.Updated, "appearance")

		// Get appearance settings
		appearanceGetReq := &GetSettingsRequest{
			Section: "appearance",
		}

		appearanceGetResp, err := settingsService.GetSettings(ctx, tenantID, appearanceGetReq)
		require.NoError(t, err)
		assert.Contains(t, appearanceGetResp.Settings, "appearance")

		appearanceSettings := appearanceGetResp.Settings["appearance"].(map[string]interface{})
		assert.Equal(t, "/uploads/logo.png", appearanceSettings["logo"])
		assert.Equal(t, "#FF6B35", appearanceSettings["primary_color"])
		assert.Equal(t, "dark", appearanceSettings["theme_mode"])
	})

	t.Run("Integration settings management", func(t *testing.T) {
		// Update integration settings
		integrationUpdateReq := &UpdateSettingsRequest{
			Settings: map[string]interface{}{
				"integrations": map[string]interface{}{
					"stripe_publishable_key": "pk_test_123456789",
					"paypal_client_id":       "paypal_123456789",
					"mailchimp_api_key":      "mc_123456789",
					"facebook_app_id":        "fb_123456789",
					"google_client_id":       "google_123456789",
				},
			},
		}

		integrationUpdateResp, err := settingsService.UpdateSettings(ctx, tenantID, integrationUpdateReq)
		require.NoError(t, err)
		assert.NotNil(t, integrationUpdateResp)
		assert.Contains(t, integrationUpdateResp.Updated, "integrations")

		// Get integration settings
		integrationGetReq := &GetSettingsRequest{
			Section: "integrations",
		}

		integrationGetResp, err := settingsService.GetSettings(ctx, tenantID, integrationGetReq)
		require.NoError(t, err)
		assert.Contains(t, integrationGetResp.Settings, "integrations")

		integrationSettings := integrationGetResp.Settings["integrations"].(map[string]interface{})
		assert.Equal(t, "pk_test_123456789", integrationSettings["stripe_publishable_key"])
		assert.Equal(t, "paypal_123456789", integrationSettings["paypal_client_id"])
	})
}

func TestSettingsService_PublicSettings(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	settingsRepo := NewRepository(testDB.DB)
	settingsService := NewService(settingsRepo)

	ctx := context.Background()
	// Use the fixed test tenant ID from fixtures
	tenantID, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	
	// Seed test data to create the tenant
	testhelpers.SeedMinimalTestData(t, testDB.DB)

	t.Run("Public settings retrieval", func(t *testing.T) {
		// First, set up some public settings
		updateReq := &UpdateSettingsRequest{
			Settings: map[string]interface{}{
				"general": map[string]interface{}{
					"store_name":        "Public Test Store",
					"store_description": "A publicly accessible store",
					"store_email":       "public@store.com",
				},
				"appearance": map[string]interface{}{
					"logo":          "/uploads/public-logo.png",
					"primary_color": "#00AA00",
					"theme_mode":    "light",
				},
				"seo": map[string]interface{}{
					"meta_title":       "Public Store",
					"meta_description": "Welcome to our public store",
				},
			},
		}

		updateResp, err := settingsService.UpdateSettings(ctx, tenantID, updateReq)
		require.NoError(t, err)
		assert.NotNil(t, updateResp)

		// Get public settings
		publicResp, err := settingsService.GetPublicSettings(ctx, tenantID)
		require.NoError(t, err)
		assert.NotNil(t, publicResp)
		assert.Equal(t, "Public Test Store", publicResp.StoreName)
		assert.Equal(t, "/uploads/public-logo.png", publicResp.Logo)
		assert.NotNil(t, publicResp.Contact)
		assert.NotNil(t, publicResp.Theme)
		assert.NotNil(t, publicResp.SEO)
	})
}

func TestSettingsService_MultiTenantIsolation(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	settingsRepo := NewRepository(testDB.DB)
	settingsService := NewService(settingsRepo)

	ctx := context.Background()
	// Use fixed test tenant IDs
	tenant1ID, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	tenant2ID, _ := uuid.Parse("22222222-2222-2222-2222-222222222222")
	
	// Seed test data to create the tenants
	testhelpers.SeedMinimalTestData(t, testDB.DB)
	
	// Create second tenant
	result := testDB.DB.Exec(`
		INSERT INTO tenants (id, name, subdomain, status, plan, currency, language, timezone, product_limit, storage_limit, bandwidth_limit, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
	`, tenant2ID, "Test Store 2", "test2", "active", "starter", "BDT", "en", "Asia/Dhaka", 100, 1024, 10240)
	require.NoError(t, result.Error)

	t.Run("Settings isolation between tenants", func(t *testing.T) {
		// Update settings for tenant 1
		tenant1UpdateReq := &UpdateSettingsRequest{
			Settings: map[string]interface{}{
				"general": map[string]interface{}{
					"store_name":  "Tenant 1 Store",
					"store_email": "tenant1@store.com",
					"currency":    "USD",
				},
			},
		}

		updateResp1, err := settingsService.UpdateSettings(ctx, tenant1ID, tenant1UpdateReq)
		require.NoError(t, err)
		assert.NotNil(t, updateResp1)

		// Update settings for tenant 2
		tenant2UpdateReq := &UpdateSettingsRequest{
			Settings: map[string]interface{}{
				"general": map[string]interface{}{
					"store_name":  "Tenant 2 Store",
					"store_email": "tenant2@store.com",
					"currency":    "BDT",
				},
			},
		}

		updateResp2, err := settingsService.UpdateSettings(ctx, tenant2ID, tenant2UpdateReq)
		require.NoError(t, err)
		assert.NotNil(t, updateResp2)

		// Verify tenant 1 can only see their settings
		getReq := &GetSettingsRequest{
			Section: "general",
		}

		tenant1Settings, err := settingsService.GetSettings(ctx, tenant1ID, getReq)
		require.NoError(t, err)
		assert.Contains(t, tenant1Settings.Settings, "general")

		tenant1General := tenant1Settings.Settings["general"].(map[string]interface{})
		assert.Equal(t, "Tenant 1 Store", tenant1General["store_name"])
		assert.Equal(t, "tenant1@store.com", tenant1General["store_email"])
		assert.Equal(t, "USD", tenant1General["currency"])

		// Verify tenant 2 can only see their settings
		tenant2Settings, err := settingsService.GetSettings(ctx, tenant2ID, getReq)
		require.NoError(t, err)
		assert.Contains(t, tenant2Settings.Settings, "general")

		tenant2General := tenant2Settings.Settings["general"].(map[string]interface{})
		assert.Equal(t, "Tenant 2 Store", tenant2General["store_name"])
		assert.Equal(t, "tenant2@store.com", tenant2General["store_email"])
		assert.Equal(t, "BDT", tenant2General["currency"])

		// Verify settings are truly isolated
		assert.NotEqual(t, tenant1General["store_name"], tenant2General["store_name"])
		assert.NotEqual(t, tenant1General["store_email"], tenant2General["store_email"])
		assert.NotEqual(t, tenant1General["currency"], tenant2General["currency"])
	})
}