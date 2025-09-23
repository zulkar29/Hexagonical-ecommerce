package platform

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"ecommerce-saas/internal/shared/testhelpers"
)

// Integration tests for platform service - critical for SaaS management

func TestPlatformService_PlatformAdminManagement(t *testing.T) {
	// Setup test database with migrations
	testDB := testhelpers.SetupTestDatabaseWithOptions(t, testhelpers.TestDatabaseOptions{
		RunMigrations: true,
	})
	defer testDB.TeardownTestDatabase(t)

	// Setup services
	platformRepo := NewRepository(testDB.DB)
	platformService := NewService(platformRepo)

	ctx := context.Background()

	t.Run("Platform admin lifecycle", func(t *testing.T) {
		// Create platform admin
		createReq := &CreatePlatformAdminRequest{
			Email:       "super@platform.com",
			Name:        "Super Administrator",
			Password:    "securepassword123",
			Role:        "super_admin",
			Permissions: []string{"*"}, // Full access
		}

		admin, err := platformService.CreatePlatformAdmin(ctx, createReq)
		require.NoError(t, err)
		assert.Equal(t, createReq.Email, admin.Email)
		assert.Equal(t, createReq.Name, admin.Name)
		assert.Equal(t, createReq.Role, admin.Role)
		assert.Equal(t, "active", admin.Status)
		assert.NotEmpty(t, admin.Password)

		// Verify password is hashed
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(createReq.Password))
		require.NoError(t, err)

		// List platform admins
		admins, err := platformService.ListPlatformAdmins(ctx, "")
		require.NoError(t, err)
		assert.Len(t, admins, 1)
		assert.Equal(t, admin.ID, admins[0].ID)

		// Get specific admin
		retrievedAdmin, err := platformService.GetPlatformAdmin(ctx, admin.ID)
		require.NoError(t, err)
		assert.Equal(t, admin.ID, retrievedAdmin.ID)
		assert.Equal(t, admin.Email, retrievedAdmin.Email)

		// Update platform admin
		updateReq := &UpdatePlatformAdminRequest{
			Name:        stringPtr("Senior Super Administrator"),
			Role:        stringPtr("super_admin"),
			Permissions: []string{"platform.*", "tenant.*", "billing.*"},
			Status:      stringPtr("active"),
		}

		updatedAdmin, err := platformService.UpdatePlatformAdmin(ctx, admin.ID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, *updateReq.Name, updatedAdmin.Name)
		assert.Equal(t, *updateReq.Role, updatedAdmin.Role)
		assert.Equal(t, *updateReq.Status, updatedAdmin.Status)

		// Delete platform admin
		err = platformService.DeletePlatformAdmin(ctx, admin.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = platformService.GetPlatformAdmin(ctx, admin.ID)
		assert.Error(t, err)
	})

	t.Run("Platform admin role filtering", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create admins with different roles
		admins := []*CreatePlatformAdminRequest{
			{
				Email:    "super@platform.com",
				Name:     "Super Admin",
				Password: "password123",
				Role:     "super_admin",
			},
			{
				Email:    "platform@platform.com",
				Name:     "Platform Admin",
				Password: "password123",
				Role:     "platform_admin",
			},
			{
				Email:    "support@platform.com",
				Name:     "Support Admin",
				Password: "password123",
				Role:     "support",
			},
		}

		for _, req := range admins {
			_, err := platformService.CreatePlatformAdmin(ctx, req)
			require.NoError(t, err)
		}

		// Filter by super_admin role
		superAdmins, err := platformService.ListPlatformAdmins(ctx, "super_admin")
		require.NoError(t, err)
		assert.Len(t, superAdmins, 1)
		assert.Equal(t, "super_admin", superAdmins[0].Role)

		// Get all admins
		allAdmins, err := platformService.ListPlatformAdmins(ctx, "")
		require.NoError(t, err)
		assert.Len(t, allAdmins, 3)
	})
}

func TestPlatformService_PlatformRoleManagement(t *testing.T) {
	// Setup test database with migrations
	testDB := testhelpers.SetupTestDatabaseWithOptions(t, testhelpers.TestDatabaseOptions{
		RunMigrations: true,
	})
	defer testDB.TeardownTestDatabase(t)

	// Setup services
	platformRepo := NewRepository(testDB.DB)
	platformService := NewService(platformRepo)

	ctx := context.Background()

	t.Run("Platform role lifecycle", func(t *testing.T) {
		// Create platform role
		createReq := &CreatePlatformRoleRequest{
			Name:        "Platform Support",
			Description: "Support role with limited platform access",
			Permissions: []string{
				"tenant.read",
				"tenant.list",
				"user.read",
				"support.all",
			},
		}

		role, err := platformService.CreatePlatformRole(ctx, createReq)
		require.NoError(t, err)
		assert.Equal(t, createReq.Name, role.Name)
		assert.Equal(t, createReq.Description, role.Description)
		assert.Equal(t, createReq.Permissions, role.Permissions)

		// List platform roles
		roles, err := platformService.ListPlatformRoles(ctx)
		require.NoError(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, role.ID, roles[0].ID)

		// Get specific role
		retrievedRole, err := platformService.GetPlatformRole(ctx, role.ID)
		require.NoError(t, err)
		assert.Equal(t, role.ID, retrievedRole.ID)
		assert.Equal(t, role.Name, retrievedRole.Name)

		// Update platform role
		updateReq := &UpdatePlatformRoleRequest{
			Name:        stringPtr("Senior Platform Support"),
			Description: stringPtr("Senior support role with extended permissions"),
			Permissions: []string{
				"tenant.*",
				"user.read",
				"user.update",
				"support.*",
				"billing.read",
			},
		}

		updatedRole, err := platformService.UpdatePlatformRole(ctx, role.ID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, *updateReq.Name, updatedRole.Name)
		assert.Equal(t, *updateReq.Description, updatedRole.Description)
		assert.Equal(t, updateReq.Permissions, updatedRole.Permissions)

		// Delete platform role
		err = platformService.DeletePlatformRole(ctx, role.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = platformService.GetPlatformRole(ctx, role.ID)
		assert.Error(t, err)
	})
}

func TestPlatformService_PlatformStats(t *testing.T) {
	// Setup test database with migrations
	testDB := testhelpers.SetupTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Seed test data
	testDB.SeedTestData(t)

	// Setup services
	platformRepo := NewRepository(testDB.DB)
	platformService := NewService(platformRepo)

	ctx := context.Background()

	t.Run("Platform statistics", func(t *testing.T) {
		// Get platform stats for different periods
		periods := []string{"day", "week", "month"}

		for _, period := range periods {
			stats, err := platformService.GetPlatformStats(ctx, period)
			require.NoError(t, err, "Failed to get stats for period: %s", period)
			assert.NotNil(t, stats)

			// Stats should have reasonable values
			assert.GreaterOrEqual(t, stats.TotalTenants, 0)
			assert.GreaterOrEqual(t, stats.ActiveTenants, 0)
			assert.GreaterOrEqual(t, stats.TotalRevenue, float64(0))
			assert.GreaterOrEqual(t, stats.TotalUsers, 0)
			assert.True(t, stats.CreatedAt.After(time.Now().Add(-time.Minute)))
		}

		// Get system status
		status, err := platformService.GetSystemStatus(ctx)
		require.NoError(t, err)
		assert.NotNil(t, status)
		assert.NotEmpty(t, status.Status)
		assert.True(t, status.CreatedAt.After(time.Now().Add(-time.Minute)))
	})
}

// Helper function
func stringPtr(s string) *string {
	return &s
}