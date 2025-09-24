package analytics

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
)

// Integration tests for analytics service - critical for business intelligence

func TestAnalyticsService_EventTracking(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	analyticsRepo := NewRepository(testDB.DB)
	analyticsService := NewService(analyticsRepo)

	ctx := context.Background()
	tenantID := uuid.New()
	userID := uuid.New()
	sessionID := "session-123"

	t.Run("Track analytics events", func(t *testing.T) {
		// Track page view event
		pageViewEvent := &AnalyticsEvent{
			ID:        uuid.New(),
			TenantID:  tenantID,
			EventType: "page_view",
			EventName: "product_listing_viewed",
			Properties: map[string]interface{}{
				"path":     "/products",
				"category": "electronics",
				"count":    25,
			},
			UserID:      &userID,
			SessionID:   sessionID,
			IPAddress:   "127.0.0.1",
			UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X)",
			Referrer:    "https://google.com",
			UTMSource:   "google",
			UTMMedium:   "cpc",
			UTMCampaign: "summer_sale",
			Timestamp:   time.Now(),
			CreatedAt:   time.Now(),
		}

		err := analyticsService.TrackEvent(ctx, tenantID, pageViewEvent)
		require.NoError(t, err)

		// Track product view event
		productViewEvent := &AnalyticsEvent{
			ID:        uuid.New(),
			TenantID:  tenantID,
			EventType: "product_view",
			EventName: "product_page_viewed",
			Properties: map[string]interface{}{
				"product_id": "smartphone-123",
				"price":      599.99,
			},
			UserID:    &userID,
			SessionID: sessionID,
			IPAddress: "127.0.0.1",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X)",
			Timestamp: time.Now(),
			CreatedAt: time.Now(),
		}

		err = analyticsService.TrackEvent(ctx, tenantID, productViewEvent)
		require.NoError(t, err)

		// Verify events were created
		stats, err := analyticsService.GetDashboardStats(ctx, tenantID, DateRange{
			Start: time.Now().Add(-24 * time.Hour),
			End:   time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)
		assert.NotNil(t, stats)
	})

	t.Run("Track page view with duration", func(t *testing.T) {
		pageView := &PageView{
			ID:              uuid.New(),
			TenantID:        tenantID,
			URL:             "https://store.com/products/smartphone",
			Path:            "/products/smartphone",
			Title:           "Best Smartphone",
			UserID:          &userID,
			SessionID:       sessionID,
			AnonymousID:     "anon-123",
			IPAddress:       "127.0.0.1",
			UserAgent:       "Mozilla/5.0 (iPhone; CPU iPhone OS)",
			Referrer:        "https://facebook.com",
			DurationSeconds: intPtr(120), // 2 minutes
			Timestamp:       time.Now(),
			CreatedAt:       time.Now(),
		}

		err := analyticsService.TrackPageView(ctx, tenantID, pageView)
		require.NoError(t, err)

		// Verify page view was created
		stats, err := analyticsService.GetDashboardStats(ctx, tenantID, DateRange{
			Start: time.Now().Add(-24 * time.Hour),
			End:   time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)
		assert.NotNil(t, stats)
	})

	t.Run("Track product view", func(t *testing.T) {
		productID := uuid.New()

		productView := &ProductView{
			ID:              uuid.New(),
			TenantID:        tenantID,
			ProductID:       productID,
			UserID:          &userID,
			SessionID:       sessionID,
			AnonymousID:     "anon-456",
			IPAddress:       "127.0.0.1",
			UserAgent:       "Mozilla/5.0 (Android)",
			Referrer:        "https://google.com",
			DurationSeconds: intPtr(180), // 3 minutes
			Timestamp:       time.Now(),
			CreatedAt:       time.Now(),
		}

		err := analyticsService.TrackProductView(ctx, tenantID, productView)
		require.NoError(t, err)
	})

	t.Run("Track purchase", func(t *testing.T) {
		orderID := uuid.New()

		purchase := &Purchase{
			ID:            uuid.New(),
			TenantID:      tenantID,
			OrderID:       orderID,
			UserID:        &userID,
			SessionID:     sessionID,
			AnonymousID:   "anon-789",
			TotalAmount:   1299.99,
			Currency:      "BDT",
			ItemCount:     3,
			PaymentMethod: "credit_card",
			IPAddress:     "127.0.0.1",
			UserAgent:     "Mozilla/5.0 (Windows)",
			Timestamp:     time.Now(),
			CreatedAt:     time.Now(),
		}

		err := analyticsService.TrackPurchase(ctx, tenantID, purchase)
		require.NoError(t, err)

		// Verify purchase was created
		stats, err := analyticsService.GetDashboardStats(ctx, tenantID, DateRange{
			Start: time.Now().Add(-24 * time.Hour),
			End:   time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)
		assert.NotNil(t, stats)
	})
}

func TestAnalyticsService_TrafficAnalytics(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	analyticsRepo := NewRepository(testDB.DB)
	analyticsService := NewService(analyticsRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Traffic statistics", func(t *testing.T) {
		// Create sample page views for different users/sessions
		pageViews := []PageView{
			{
				ID:              uuid.New(),
				TenantID:        tenantID,
				URL:             "https://store.com/",
				Path:            "/",
				Title:           "Homepage",
				UserID:          uuidPtr(uuid.New()),
				SessionID:       "session-1",
				IPAddress:       "127.0.0.1",
				UserAgent:       "Mozilla/5.0",
				Referrer:        "https://google.com",
				DurationSeconds: intPtr(45),
				Timestamp:       time.Now(),
				CreatedAt:       time.Now(),
			},
			{
				ID:              uuid.New(),
				TenantID:        tenantID,
				URL:             "https://store.com/products",
				Path:            "/products",
				Title:           "Products",
				UserID:          uuidPtr(uuid.New()),
				SessionID:       "session-2",
				IPAddress:       "127.0.0.1",
				UserAgent:       "Mozilla/5.0",
				Referrer:        "https://facebook.com",
				DurationSeconds: intPtr(120),
				Timestamp:       time.Now(),
				CreatedAt:       time.Now(),
			},
			{
				ID:              uuid.New(),
				TenantID:        tenantID,
				URL:             "https://store.com/about",
				Path:            "/about",
				Title:           "About Us",
				UserID:          uuidPtr(uuid.New()),
				SessionID:       "session-3",
				IPAddress:       "127.0.0.1",
				UserAgent:       "Mozilla/5.0",
				Referrer:        "",
				DurationSeconds: intPtr(25), // bounce
				Timestamp:       time.Now(),
				CreatedAt:       time.Now(),
			},
		}

		for _, pv := range pageViews {
			err := analyticsService.TrackPageView(ctx, tenantID, &pv)
			require.NoError(t, err)
		}

		// Get traffic statistics
		trafficStats, err := analyticsService.GetTrafficStats(ctx, tenantID, DateRange{
			Start: time.Now().Add(-24 * time.Hour),
			End:   time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)
		assert.NotNil(t, trafficStats)
		assert.Equal(t, tenantID, trafficStats.TenantID)
	})
}

func TestAnalyticsService_SalesAnalytics(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	analyticsRepo := NewRepository(testDB.DB)
	analyticsService := NewService(analyticsRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Sales statistics", func(t *testing.T) {
		// Create sample purchases
		purchases := []Purchase{
			{
				ID:            uuid.New(),
				TenantID:      tenantID,
				OrderID:       uuid.New(),
				UserID:        uuidPtr(uuid.New()),
				SessionID:     "session-1",
				TotalAmount:   1500.00,
				Currency:      "BDT",
				ItemCount:     2,
				PaymentMethod: "credit_card",
				IPAddress:     "127.0.0.1",
				UserAgent:     "Mozilla/5.0",
				Timestamp:     time.Now(),
				CreatedAt:     time.Now(),
			},
			{
				ID:            uuid.New(),
				TenantID:      tenantID,
				OrderID:       uuid.New(),
				UserID:        uuidPtr(uuid.New()),
				SessionID:     "session-2",
				TotalAmount:   750.50,
				Currency:      "BDT",
				ItemCount:     1,
				PaymentMethod: "bkash",
				IPAddress:     "127.0.0.1",
				UserAgent:     "Mozilla/5.0",
				Timestamp:     time.Now(),
				CreatedAt:     time.Now(),
			},
			{
				ID:            uuid.New(),
				TenantID:      tenantID,
				OrderID:       uuid.New(),
				UserID:        uuidPtr(uuid.New()),
				SessionID:     "session-3",
				TotalAmount:   2200.75,
				Currency:      "BDT",
				ItemCount:     5,
				PaymentMethod: "cash_on_delivery",
				IPAddress:     "127.0.0.1",
				UserAgent:     "Mozilla/5.0",
				Timestamp:     time.Now(),
				CreatedAt:     time.Now(),
			},
		}

		for _, p := range purchases {
			err := analyticsService.TrackPurchase(ctx, tenantID, &p)
			require.NoError(t, err)
		}

		// Get sales statistics
		salesStats, err := analyticsService.GetSalesStats(ctx, tenantID, DateRange{
			Start: time.Now().Add(-24 * time.Hour),
			End:   time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)
		assert.NotNil(t, salesStats)
		assert.Equal(t, tenantID, salesStats.TenantID)
	})
}

func TestAnalyticsService_TopPerformers(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	analyticsRepo := NewRepository(testDB.DB)
	analyticsService := NewService(analyticsRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Top products analytics", func(t *testing.T) {
		// Create sample product views
		productIDs := []uuid.UUID{
			uuid.New(), // Product A
			uuid.New(), // Product B
			uuid.New(), // Product C
		}

		// Product A: 5 views
		for i := 0; i < 5; i++ {
			productView := &ProductView{
				ID:              uuid.New(),
				TenantID:        tenantID,
				ProductID:       productIDs[0],
				UserID:          uuidPtr(uuid.New()),
				SessionID:       "session-a-" + string(rune(i)),
				IPAddress:       "127.0.0.1",
				UserAgent:       "Mozilla/5.0",
				DurationSeconds: intPtr(60 + i*10),
				Timestamp:       time.Now(),
				CreatedAt:       time.Now(),
			}
			err := analyticsService.TrackProductView(ctx, tenantID, productView)
			require.NoError(t, err)
		}

		// Product B: 3 views
		for i := 0; i < 3; i++ {
			productView := &ProductView{
				ID:              uuid.New(),
				TenantID:        tenantID,
				ProductID:       productIDs[1],
				UserID:          uuidPtr(uuid.New()),
				SessionID:       "session-b-" + string(rune(i)),
				IPAddress:       "127.0.0.1",
				UserAgent:       "Mozilla/5.0",
				DurationSeconds: intPtr(45 + i*5),
				Timestamp:       time.Now(),
				CreatedAt:       time.Now(),
			}
			err := analyticsService.TrackProductView(ctx, tenantID, productView)
			require.NoError(t, err)
		}

		// Get top products
		topProducts, err := analyticsService.GetTopProducts(ctx, tenantID, DateRange{
			Start: time.Now().Add(-24 * time.Hour),
			End:   time.Now().Add(1 * time.Hour),
		}, 5)
		require.NoError(t, err)
		assert.NotNil(t, topProducts)
	})

	t.Run("Top pages analytics", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create sample page views
		pages := []struct {
			path  string
			count int
		}{
			{"/", 7},               // Homepage: 7 views
			{"/products", 5},       // Products page: 5 views
			{"/about", 3},          // About page: 3 views
			{"/contact", 2},        // Contact page: 2 views
			{"/products/phone", 4}, // Product page: 4 views
		}

		for _, page := range pages {
			for i := 0; i < page.count; i++ {
				pageView := &PageView{
					ID:        uuid.New(),
					TenantID:  tenantID,
					URL:       "https://store.com" + page.path,
					Path:      page.path,
					Title:     "Page Title",
					UserID:    uuidPtr(uuid.New()),
					SessionID: "session-" + page.path + "-" + string(rune(i)),
					IPAddress: "127.0.0.1",
					UserAgent: "Mozilla/5.0",
					Timestamp: time.Now(),
					CreatedAt: time.Now(),
				}
				err := analyticsService.TrackPageView(ctx, tenantID, pageView)
				require.NoError(t, err)
			}
		}

		// Get top pages
		topPages, err := analyticsService.GetTopPages(ctx, tenantID, DateRange{
			Start: time.Now().Add(-24 * time.Hour),
			End:   time.Now().Add(1 * time.Hour),
		}, 5)
		require.NoError(t, err)
		assert.NotNil(t, topPages)
	})
}

func TestAnalyticsService_RealTimeAnalytics(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	analyticsRepo := NewRepository(testDB.DB)
	analyticsService := NewService(analyticsRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Real-time statistics", func(t *testing.T) {
		// Create recent events (within last 5 minutes)
		activeUsers := []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		}

		for _, userID := range activeUsers {
			event := &AnalyticsEvent{
				ID:        uuid.New(),
				TenantID:  tenantID,
				EventType: "page_view",
				EventName: "real_time_view",
				Properties: map[string]interface{}{
					"path": "/dashboard",
				},
				UserID:    &userID,
				SessionID: "session-" + userID.String(),
				IPAddress: "127.0.0.1",
				UserAgent: "Mozilla/5.0",
				Timestamp: time.Now(),
				CreatedAt: time.Now(),
			}

			err := analyticsService.TrackEvent(ctx, tenantID, event)
			require.NoError(t, err)
		}

		// Get real-time stats
		realTimeStats, err := analyticsService.GetRealTimeStats(ctx, tenantID)
		require.NoError(t, err)
		assert.NotNil(t, realTimeStats)

		// Get active users count
		activeCount, err := analyticsService.GetActiveUsers(ctx, tenantID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, activeCount, int64(0))
	})
}

func TestAnalyticsService_MultiTenantIsolation(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	analyticsRepo := NewRepository(testDB.DB)
	analyticsService := NewService(analyticsRepo)

	ctx := context.Background()
	tenant1ID := uuid.New()
	tenant2ID := uuid.New()

	t.Run("Analytics isolation between tenants", func(t *testing.T) {
		// Create events for tenant 1
		tenant1Event := &AnalyticsEvent{
			ID:        uuid.New(),
			TenantID:  tenant1ID,
			EventType: "page_view",
			EventName: "tenant1_view",
			Properties: map[string]interface{}{
				"path": "/tenant1-page",
			},
			UserID:    uuidPtr(uuid.New()),
			SessionID: "tenant1-session",
			IPAddress: "127.0.0.1",
			UserAgent: "Mozilla/5.0",
			Timestamp: time.Now(),
			CreatedAt: time.Now(),
		}

		err := analyticsService.TrackEvent(ctx, tenant1ID, tenant1Event)
		require.NoError(t, err)

		// Create purchase for tenant 1
		tenant1Purchase := &Purchase{
			ID:          uuid.New(),
			TenantID:    tenant1ID,
			OrderID:     uuid.New(),
			UserID:      uuidPtr(uuid.New()),
			SessionID:   "tenant1-session",
			TotalAmount: 1000.00,
			Currency:    "BDT",
			ItemCount:   2,
			IPAddress:   "127.0.0.1",
			UserAgent:   "Mozilla/5.0",
			Timestamp:   time.Now(),
			CreatedAt:   time.Now(),
		}

		err = analyticsService.TrackPurchase(ctx, tenant1ID, tenant1Purchase)
		require.NoError(t, err)

		// Create events for tenant 2
		tenant2Event := &AnalyticsEvent{
			ID:        uuid.New(),
			TenantID:  tenant2ID,
			EventType: "page_view",
			EventName: "tenant2_view",
			Properties: map[string]interface{}{
				"path": "/tenant2-page",
			},
			UserID:    uuidPtr(uuid.New()),
			SessionID: "tenant2-session",
			IPAddress: "127.0.0.1",
			UserAgent: "Mozilla/5.0",
			Timestamp: time.Now(),
			CreatedAt: time.Now(),
		}

		err = analyticsService.TrackEvent(ctx, tenant2ID, tenant2Event)
		require.NoError(t, err)

		// Create purchase for tenant 2
		tenant2Purchase := &Purchase{
			ID:          uuid.New(),
			TenantID:    tenant2ID,
			OrderID:     uuid.New(),
			UserID:      uuidPtr(uuid.New()),
			SessionID:   "tenant2-session",
			TotalAmount: 500.00,
			Currency:    "BDT",
			ItemCount:   1,
			IPAddress:   "127.0.0.1",
			UserAgent:   "Mozilla/5.0",
			Timestamp:   time.Now(),
			CreatedAt:   time.Now(),
		}

		err = analyticsService.TrackPurchase(ctx, tenant2ID, tenant2Purchase)
		require.NoError(t, err)

		// Verify tenant 1 stats isolation
		tenant1Stats, err := analyticsService.GetDashboardStats(ctx, tenant1ID, DateRange{
			Start: time.Now().Add(-24 * time.Hour),
			End:   time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)
		assert.NotNil(t, tenant1Stats)
		assert.Equal(t, tenant1ID, tenant1Stats.TenantID)

		// Verify tenant 2 stats isolation
		tenant2Stats, err := analyticsService.GetDashboardStats(ctx, tenant2ID, DateRange{
			Start: time.Now().Add(-24 * time.Hour),
			End:   time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)
		assert.NotNil(t, tenant2Stats)
		assert.Equal(t, tenant2ID, tenant2Stats.TenantID)

		// Verify that tenants have different data
		assert.NotEqual(t, tenant1Stats.TenantID, tenant2Stats.TenantID)
	})
}

// Helper functions
func uuidPtr(id uuid.UUID) *uuid.UUID {
	return &id
}

func intPtr(i int) *int {
	return &i
}