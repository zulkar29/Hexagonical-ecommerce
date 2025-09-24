package admin

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
)

// Integration tests for admin service - critical for tenant management

func TestAdminService_DashboardOperations(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	adminRepo := NewRepository(testDB.DB)
	adminService := NewService(adminRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Dashboard statistics", func(t *testing.T) {
		// Get dashboard stats
		dashboardReq := DashboardRequest{
			Period:  "month",
			Metrics: []string{"sales", "orders", "customers"},
		}

		stats, err := adminService.GetDashboardStats(ctx, &tenantID, dashboardReq)
		require.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, "month", stats.Period)
		assert.True(t, stats.GeneratedAt.After(time.Now().Add(-time.Minute)))

		// Quick stats
		quickStats, err := adminService.GetQuickStats(ctx, &tenantID)
		require.NoError(t, err)
		assert.NotNil(t, quickStats)
		assert.True(t, quickStats.GeneratedAt.After(time.Now().Add(-time.Minute)))
	})

	t.Run("System health monitoring", func(t *testing.T) {
		health, err := adminService.GetSystemHealth(ctx)
		require.NoError(t, err)
		assert.NotNil(t, health)
		assert.NotEmpty(t, health.Status)
		// Note: CheckedAt field may not be available in SystemHealth model
	})
}

func TestAdminService_StaffManagement(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	adminRepo := NewRepository(testDB.DB)
	adminService := NewService(adminRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Staff lifecycle management", func(t *testing.T) {
		// Create staff member
		staffReq := StaffRequest{
			Email:     "manager@teststore.com",
			FirstName: "John",
			LastName:  "Manager",
			Role:      "manager",
		}

		staff, err := adminService.CreateStaff(ctx, &tenantID, staffReq)
		require.NoError(t, err)
		assert.Equal(t, staffReq.Email, staff.Email)
		assert.Equal(t, staffReq.FirstName, staff.FirstName)
		assert.Equal(t, staffReq.Role, staff.Role)
		assert.Equal(t, tenantID, *staff.TenantID)

		// List staff members
		staffList, err := adminService.ListStaff(ctx, &tenantID, "", "")
		require.NoError(t, err)
		assert.Len(t, staffList, 1)
		assert.Equal(t, staff.ID, staffList[0].ID)

		// Update staff member
		updateReq := StaffRequest{
			Email:     "john.manager@teststore.com",
			FirstName: "John",
			LastName:  "Manager Senior",
			Role:      "senior_manager",
		}

		updatedStaff, err := adminService.UpdateStaff(ctx, &tenantID, staff.ID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, updateReq.Email, updatedStaff.Email)
		assert.Equal(t, updateReq.LastName, updatedStaff.LastName)
		assert.Equal(t, updateReq.Role, updatedStaff.Role)

		// Change staff status
		err = adminService.ChangeStaffStatus(ctx, &tenantID, staff.ID, "inactive")
		require.NoError(t, err)

		// Delete staff member
		err = adminService.DeleteStaff(ctx, &tenantID, staff.ID)
		require.NoError(t, err)

		// Verify deletion
		finalStaffList, err := adminService.ListStaff(ctx, &tenantID, "", "")
		require.NoError(t, err)
		assert.Len(t, finalStaffList, 0)
	})

	t.Run("Staff filtering", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create multiple staff members
		adminStaff := StaffRequest{
			Email:     "admin@teststore.com",
			FirstName: "Admin",
			LastName:  "User",
			Role:      "admin",
		}

		managerStaff := StaffRequest{
			Email:     "manager@teststore.com",
			FirstName: "Manager",
			LastName:  "User",
			Role:      "manager",
		}

		_, err := adminService.CreateStaff(ctx, &tenantID, adminStaff)
		require.NoError(t, err)

		_, err = adminService.CreateStaff(ctx, &tenantID, managerStaff)
		require.NoError(t, err)

		// Filter by role
		adminStaffList, err := adminService.ListStaff(ctx, &tenantID, "admin", "")
		require.NoError(t, err)
		assert.Len(t, adminStaffList, 1)
		assert.Equal(t, "admin", adminStaffList[0].Role)

		// Get all staff
		allStaff, err := adminService.ListStaff(ctx, &tenantID, "", "")
		require.NoError(t, err)
		assert.Len(t, allStaff, 2)
	})
}

func TestAdminService_RoleManagement(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	adminRepo := NewRepository(testDB.DB)
	adminService := NewService(adminRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Role lifecycle management", func(t *testing.T) {
		// Create role
		roleReq := RoleRequest{
			Name:        "Store Manager",
			Description: "Manages store operations",
		}

		role, err := adminService.CreateRole(ctx, &tenantID, roleReq)
		require.NoError(t, err)
		assert.Equal(t, roleReq.Name, role.Name)
		assert.Equal(t, roleReq.Description, role.Description)
		assert.Equal(t, tenantID, *role.TenantID)

		// List roles
		roles, err := adminService.ListRoles(ctx, &tenantID, false)
		require.NoError(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, role.ID, roles[0].ID)

		// Update role
		updateRoleReq := RoleRequest{
			Name:        "Senior Store Manager",
			Description: "Senior manager with extended permissions",
		}

		updatedRole, err := adminService.UpdateRole(ctx, &tenantID, role.ID, updateRoleReq)
		require.NoError(t, err)
		assert.Equal(t, updateRoleReq.Name, updatedRole.Name)
		assert.Equal(t, updateRoleReq.Description, updatedRole.Description)

		// Delete role
		err = adminService.DeleteRole(ctx, &tenantID, role.ID)
		require.NoError(t, err)

		// Verify deletion
		rolesAfterDelete, err := adminService.ListRoles(ctx, &tenantID, false)
		require.NoError(t, err)
		assert.Len(t, rolesAfterDelete, 0)
	})
}

func TestAdminService_MultiTenantIsolation(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	// Database schema is handled by raw SQL migrations in /migrations directory

	// Setup services
	adminRepo := NewRepository(testDB.DB)
	adminService := NewService(adminRepo)

	ctx := context.Background()
	tenant1ID := uuid.New()
	tenant2ID := uuid.New()

	t.Run("Staff isolation between tenants", func(t *testing.T) {
		// Create staff for tenant 1
		staff1Req := StaffRequest{
			Email:     "admin1@tenant1.com",
			FirstName: "Tenant1",
			LastName:  "Admin",
			Role:      "admin",
		}

		staff1, err := adminService.CreateStaff(ctx, &tenant1ID, staff1Req)
		require.NoError(t, err)

		// Create staff for tenant 2
		staff2Req := StaffRequest{
			Email:     "admin2@tenant2.com",
			FirstName: "Tenant2",
			LastName:  "Admin",
			Role:      "admin",
		}

		staff2, err := adminService.CreateStaff(ctx, &tenant2ID, staff2Req)
		require.NoError(t, err)

		// Verify tenant 1 can only see their staff
		tenant1Staff, err := adminService.ListStaff(ctx, &tenant1ID, "", "")
		require.NoError(t, err)
		assert.Len(t, tenant1Staff, 1)
		assert.Equal(t, staff1.ID, tenant1Staff[0].ID)
		assert.Equal(t, tenant1ID, *tenant1Staff[0].TenantID)

		// Verify tenant 2 can only see their staff
		tenant2Staff, err := adminService.ListStaff(ctx, &tenant2ID, "", "")
		require.NoError(t, err)
		assert.Len(t, tenant2Staff, 1)
		assert.Equal(t, staff2.ID, tenant2Staff[0].ID)
		assert.Equal(t, tenant2ID, *tenant2Staff[0].TenantID)
	})

	t.Run("Role isolation between tenants", func(t *testing.T) {
		// Create role for tenant 1
		role1Req := RoleRequest{
			Name:        "Tenant1 Manager",
			Description: "Manager role for tenant 1",
		}

		role1, err := adminService.CreateRole(ctx, &tenant1ID, role1Req)
		require.NoError(t, err)

		// Create role for tenant 2
		role2Req := RoleRequest{
			Name:        "Tenant2 Manager",
			Description: "Manager role for tenant 2",
		}

		role2, err := adminService.CreateRole(ctx, &tenant2ID, role2Req)
		require.NoError(t, err)

		// Verify tenant 1 can only see their roles
		tenant1Roles, err := adminService.ListRoles(ctx, &tenant1ID, false)
		require.NoError(t, err)
		assert.Len(t, tenant1Roles, 1)
		assert.Equal(t, role1.ID, tenant1Roles[0].ID)

		// Verify tenant 2 can only see their roles
		tenant2Roles, err := adminService.ListRoles(ctx, &tenant2ID, false)
		require.NoError(t, err)
		assert.Len(t, tenant2Roles, 1)
		assert.Equal(t, role2.ID, tenant2Roles[0].ID)
	})
}

