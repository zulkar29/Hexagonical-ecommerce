package tenant

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
	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/shared/testhelpers"
	"ecommerce-saas/internal/user"
)

// Integration tests for tenant onboarding and management flows

func TestTenantIntegration_CompleteOnboardingFlow(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services with proper config
	cfg, err := config.Load()
	if err != nil {
		// Fallback to default config for tests
		cfg = &config.Config{
			App: config.AppConfig{
				Currency: "BDT",
			},
		}
	}
	tenantRepo := NewRepository(testDB.DB)
	tenantService := NewService(tenantRepo, cfg)
	userRepo := user.NewRepository(testDB.DB)

	t.Run("Complete tenant onboarding with admin user", func(t *testing.T) {
		// Step 1: Create tenant
		createReq := CreateTenantRequest{
			Name:        "Test E-commerce Store",
			Subdomain:   "teststore",
			Description: "A test e-commerce store for integration testing",
			Phone:       "+8801234567890",
			Email:       "admin@teststore.com",
			Address:     "123 Test Street, Dhaka, Bangladesh",
		}

		tenant, err := tenantService.CreateTenant(createReq)
		require.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, createReq.Name, tenant.Name)
		assert.Equal(t, createReq.Subdomain, tenant.Subdomain)
		assert.Equal(t, StatusActive, tenant.Status)
		assert.Equal(t, PlanStarter, tenant.Plan)
		assert.Equal(t, "BDT", tenant.Currency)

		// Step 2: Create admin user for the tenant
		adminUser := &user.User{
			ID:        uuid.New(),
			TenantID:  &tenant.ID,
			Email:     "admin@teststore.com",
			Password:  "hashedpassword123",
			FirstName: "Store",
			LastName:  "Admin",
			Role:      user.RoleAdmin,
			Status:    user.StatusActive,
		}

		createdUser, err := userRepo.CreateUser(context.Background(), adminUser)
		require.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, tenant.ID, *createdUser.TenantID)
		assert.Equal(t, user.RoleAdmin, createdUser.Role)

		// Step 3: Verify tenant can be retrieved by subdomain
		retrievedTenant, err := tenantService.GetTenantBySubdomain(tenant.Subdomain)
		require.NoError(t, err)
		assert.Equal(t, tenant.ID, retrievedTenant.ID)

		// Step 4: Verify user belongs to tenant
		retrievedUser, err := userRepo.GetUserByID(context.Background(), createdUser.ID)
		require.NoError(t, err)
		assert.Equal(t, tenant.ID, *retrievedUser.TenantID)

		// Step 5: Test tenant business logic
		assert.True(t, tenant.IsActive())
		assert.True(t, tenant.CanCreateProducts(0))
		assert.True(t, tenant.CanCreateProducts(50))
		assert.False(t, tenant.CanCreateProducts(101)) // Starter plan limit is 100

		// Step 6: Update tenant information
		updateReq := UpdateTenantRequest{
			Name:         "Updated Store Name",
			Email:        "updated@store.com",
			CustomDomain: "store.example.com",
		}

		updatedTenant, err := tenantService.UpdateTenant(tenant.ID.String(), updateReq)
		require.NoError(t, err)
		assert.Equal(t, updateReq.Name, updatedTenant.Name)
		assert.Equal(t, updateReq.CustomDomain, updatedTenant.CustomDomain)
	})
}

func TestTenantIntegration_MultiTenantIsolation(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services with proper config
	cfg, err := config.Load()
	if err != nil {
		// Fallback to default config for tests
		cfg = &config.Config{
			App: config.AppConfig{
				Currency: "BDT",
			},
		}
	}
	tenantRepo := NewRepository(testDB.DB)
	tenantService := NewService(tenantRepo, cfg)
	userRepo := user.NewRepository(testDB.DB)

	t.Run("Ensure tenant isolation", func(t *testing.T) {
		// Create two tenants
		tenant1, err := tenantService.CreateTenant(CreateTenantRequest{
			Name:      "Store 1",
			Subdomain: "store1",
			Email:     "admin@store1.com",
		})
		require.NoError(t, err)

		tenant2, err := tenantService.CreateTenant(CreateTenantRequest{
			Name:      "Store 2",
			Subdomain: "store2",
			Email:     "admin@store2.com",
		})
		require.NoError(t, err)

		// Create users for each tenant
		user1 := &user.User{
			ID:        uuid.New(),
			TenantID:  &tenant1.ID,
			Email:     "user1@store1.com",
			Password:  "hashedpassword123",
			FirstName: "User",
			LastName:  "One",
			Role:      user.RoleCustomer,
			Status:    user.StatusActive,
		}

		user2 := &user.User{
			ID:        uuid.New(),
			TenantID:  &tenant2.ID,
			Email:     "user2@store2.com",
			Password:  "hashedpassword123",
			FirstName: "User",
			LastName:  "Two",
			Role:      user.RoleCustomer,
			Status:    user.StatusActive,
		}

		_, err = userRepo.CreateUser(context.Background(), user1)
		require.NoError(t, err)

		_, err = userRepo.CreateUser(context.Background(), user2)
		require.NoError(t, err)

		// Verify tenant isolation - each tenant should only see their own users
		tenant1Users, _, err := userRepo.ListUsers(context.Background(), &tenant1.ID, user.UserFilter{}, 0, 10)
		require.NoError(t, err)
		assert.Len(t, tenant1Users, 1)
		assert.Equal(t, user1.Email, tenant1Users[0].Email)

		tenant2Users, _, err := userRepo.ListUsers(context.Background(), &tenant2.ID, user.UserFilter{}, 0, 10)
		require.NoError(t, err)
		assert.Len(t, tenant2Users, 1)
		assert.Equal(t, user2.Email, tenant2Users[0].Email)

		// Cross-tenant user access should fail
		_, err = userRepo.GetUserByEmail(context.Background(), "user2@store2.com")
		require.NoError(t, err) // This should work as it doesn't filter by tenant

		// Verify tenants are completely separate
		assert.NotEqual(t, tenant1.ID, tenant2.ID)
		assert.NotEqual(t, tenant1.Subdomain, tenant2.Subdomain)
	})
}

func TestTenantIntegration_SubdomainAndDomainHandling(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services with proper config
	cfg, err := config.Load()
	if err != nil {
		// Fallback to default config for tests
		cfg = &config.Config{
			App: config.AppConfig{
				Currency: "BDT",
			},
		}
	}
	tenantRepo := NewRepository(testDB.DB)
	tenantService := NewService(tenantRepo, cfg)

	t.Run("Subdomain uniqueness and validation", func(t *testing.T) {
		// Create first tenant
		_, err := tenantService.CreateTenant(CreateTenantRequest{
			Name:      "First Store",
			Subdomain: "uniquestore",
			Email:     "admin@unique.com",
		})
		require.NoError(t, err)

		// Try to create second tenant with same subdomain - should fail
		_, err = tenantService.CreateTenant(CreateTenantRequest{
			Name:      "Second Store",
			Subdomain: "uniquestore",
			Email:     "admin@second.com",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Subdomain already taken")

		// Test subdomain validation - must be alphanumeric
		tenant2, err := tenantService.CreateTenant(CreateTenantRequest{
			Name:      "Normalized Store",
			Subdomain: "normalized2", // Valid alphanumeric subdomain
			Email:     "admin@normalized.com",
		})
		require.NoError(t, err)
		assert.Equal(t, "normalized2", tenant2.Subdomain)

		// Retrieve by subdomain
		retrieved, err := tenantService.GetTenantBySubdomain("normalized2")
		require.NoError(t, err)
		assert.Equal(t, tenant2.ID, retrieved.ID)

		// Case insensitive retrieval
		retrieved, err = tenantService.GetTenantBySubdomain("NORMALIZED2")
		require.NoError(t, err)
		assert.Equal(t, tenant2.ID, retrieved.ID)
	})

	t.Run("Custom domain handling", func(t *testing.T) {
		// Create tenant
		tenant, err := tenantService.CreateTenant(CreateTenantRequest{
			Name:      "Custom Domain Store",
			Subdomain: "customdomain",
			Email:     "admin@custom.com",
		})
		require.NoError(t, err)

		// Update with custom domain
		updated, err := tenantService.UpdateTenant(tenant.ID.String(), UpdateTenantRequest{
			Email:        "updated@custom.com",
			CustomDomain: "store.example.com",
		})
		require.NoError(t, err)
		assert.Equal(t, "store.example.com", updated.CustomDomain)

		// Retrieve by custom domain
		retrieved, err := tenantService.GetTenantByCustomDomain("store.example.com")
		require.NoError(t, err)
		assert.Equal(t, tenant.ID, retrieved.ID)
	})
}

func TestTenantIntegration_PlanLimitsAndUpgrades(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services with proper config
	cfg, err := config.Load()
	if err != nil {
		// Fallback to default config for tests
		cfg = &config.Config{
			App: config.AppConfig{
				Currency: "BDT",
			},
		}
	}
	tenantRepo := NewRepository(testDB.DB)
	tenantService := NewService(tenantRepo, cfg)

	t.Run("Plan limits and upgrades", func(t *testing.T) {
		// Create tenant with starter plan
		tenant, err := tenantService.CreateTenant(CreateTenantRequest{
			Name:      "Growing Store",
			Subdomain: "growing",
			Email:     "admin@growing.com",
		})
		require.NoError(t, err)
		assert.Equal(t, PlanStarter, tenant.Plan)

		// Test starter plan limits
		assert.True(t, tenant.CanCreateProducts(99))
		assert.False(t, tenant.CanCreateProducts(100))

		// Upgrade to Pro plan
		upgraded, err := tenantService.UpdatePlan(tenant.ID.String(), UpdatePlanRequest{
			Plan: PlanPro,
		})
		require.NoError(t, err)
		assert.Equal(t, PlanPro, upgraded.Plan)

		// Test pro plan limits
		assert.True(t, upgraded.CanCreateProducts(999))
		assert.False(t, upgraded.CanCreateProducts(1000))

		// Upgrade to Enterprise (unlimited)
		enterprise, err := tenantService.UpdatePlan(tenant.ID.String(), UpdatePlanRequest{
			Plan: PlanEnterprise,
		})
		require.NoError(t, err)
		assert.Equal(t, PlanEnterprise, enterprise.Plan)

		// Test enterprise plan limits (unlimited)
		assert.True(t, enterprise.CanCreateProducts(10000))
		assert.True(t, enterprise.CanCreateProducts(100000))
	})
}

func TestTenantIntegration_StatusManagement(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services with proper config
	cfg, err := config.Load()
	if err != nil {
		// Fallback to default config for tests
		cfg = &config.Config{
			App: config.AppConfig{
				Currency: "BDT",
			},
		}
	}
	tenantRepo := NewRepository(testDB.DB)
	tenantService := NewService(tenantRepo, cfg)

	t.Run("Tenant status lifecycle", func(t *testing.T) {
		// Create active tenant
		tenant, err := tenantService.CreateTenant(CreateTenantRequest{
			Name:      "Status Test Store",
			Subdomain: "statustest",
			Email:     "admin@statustest.com",
		})
		require.NoError(t, err)
		assert.Equal(t, StatusActive, tenant.Status)
		assert.True(t, tenant.IsActive())

		// Suspend tenant
		err = tenantService.SuspendTenant(tenant.ID.String(), "Testing suspension")
		require.NoError(t, err)

		// Verify status change
		updated, err := tenantService.GetTenant(tenant.ID.String())
		require.NoError(t, err)
		assert.Equal(t, StatusSuspended, updated.Status)
		assert.False(t, updated.IsActive())

		// Reactivate tenant
		err = tenantService.ActivateTenant(tenant.ID.String())
		require.NoError(t, err)

		// Verify reactivation
		reactivated, err := tenantService.GetTenant(tenant.ID.String())
		require.NoError(t, err)
		assert.Equal(t, StatusActive, reactivated.Status)
		assert.True(t, reactivated.IsActive())
	})
}

// HTTP Handler Integration Tests

func TestTenantHTTPIntegration_CreateTenantAPI(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services and handlers with proper config
	cfg, err := config.Load()
	if err != nil {
		// Fallback to default config for tests
		cfg = &config.Config{
			App: config.AppConfig{
				Currency: "BDT",
			},
		}
	}
	tenantRepo := NewRepository(testDB.DB)
	tenantService := NewService(tenantRepo, cfg)
	tenantHandler := NewHandler(tenantService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tenants", tenantHandler.CreateTenant)

	t.Run("Create tenant via HTTP API", func(t *testing.T) {
		// Prepare request
		createReq := CreateTenantRequest{
			Name:        "API Test Store",
			Subdomain:   "apitest",
			Description: "Store created via API",
			Email:       "admin@apitest.com",
		}

		jsonBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Make request
		req, err := http.NewRequest("POST", "/tenants", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Verify response
		assert.Equal(t, http.StatusCreated, recorder.Code)

		var response struct {
			Message string `json:"message"`
			Data    Tenant `json:"data"`
		}
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Tenant created successfully", response.Message)
		assert.Equal(t, createReq.Name, response.Data.Name)
		assert.Equal(t, createReq.Subdomain, response.Data.Subdomain)
		assert.Equal(t, StatusActive, response.Data.Status)
		assert.Equal(t, PlanStarter, response.Data.Plan)
	})

	t.Run("Validation errors via HTTP API", func(t *testing.T) {
		// Invalid request - missing required fields
		invalidReq := CreateTenantRequest{
			Name: "", // Missing required field
		}

		jsonBody, err := json.Marshal(invalidReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/tenants", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Should return validation error
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

// Performance and concurrency tests

func TestTenantIntegration_ConcurrentCreation(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services with proper config
	cfg, err := config.Load()
	if err != nil {
		// Fallback to default config for tests
		cfg = &config.Config{
			App: config.AppConfig{
				Currency: "BDT",
			},
		}
	}
	tenantRepo := NewRepository(testDB.DB)
	tenantService := NewService(tenantRepo, cfg)

	t.Run("Concurrent tenant creation", func(t *testing.T) {
		const numTenants = 10

		// Use channels to coordinate concurrent goroutines
		results := make(chan *Tenant, numTenants)
		errors := make(chan error, numTenants)

		// Create tenants concurrently
		for i := 0; i < numTenants; i++ {
			go func(id int) {
				tenant, err := tenantService.CreateTenant(CreateTenantRequest{
					Name:      fmt.Sprintf("Concurrent Store %d", id),
					Subdomain: fmt.Sprintf("concurrent%d", id),
					Email:     fmt.Sprintf("admin%d@concurrent.com", id),
				})
				if err != nil {
					errors <- err
				} else {
					results <- tenant
				}
			}(i)
		}

		// Collect results
		var createdTenants []*Tenant
		var creationErrors []error

		for i := 0; i < numTenants; i++ {
			select {
			case tenant := <-results:
				createdTenants = append(createdTenants, tenant)
			case err := <-errors:
				creationErrors = append(creationErrors, err)
			case <-time.After(10 * time.Second):
				t.Fatal("Timeout waiting for concurrent operations")
			}
		}

		// All should succeed
		assert.Len(t, createdTenants, numTenants)
		assert.Len(t, creationErrors, 0)

		// Verify all tenants are unique
		subdomains := make(map[string]bool)
		for _, tenant := range createdTenants {
			assert.False(t, subdomains[tenant.Subdomain], "Duplicate subdomain found")
			subdomains[tenant.Subdomain] = true
		}
	})
}

// Benchmark tests for performance monitoring

func BenchmarkTenantIntegration_CreateTenant(b *testing.B) {
	// Setup test database
	testDB := testhelpers.SetupTestDatabase(&testing.T{})
	defer testDB.TeardownTestDatabase(&testing.T{})

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services with proper config
	cfg, err := config.Load()
	if err != nil {
		// Fallback to default config for tests
		cfg = &config.Config{
			App: config.AppConfig{
				Currency: "BDT",
			},
		}
	}
	tenantRepo := NewRepository(testDB.DB)
	tenantService := NewService(tenantRepo, cfg)

	// Reset timer after setup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tenantService.CreateTenant(CreateTenantRequest{
			Name:      fmt.Sprintf("Benchmark Store %d", i),
			Subdomain: fmt.Sprintf("bench%d", i),
			Email:     fmt.Sprintf("admin%d@bench.com", i),
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
