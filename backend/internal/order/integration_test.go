package order

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
	"ecommerce-saas/internal/tenant"
	"ecommerce-saas/internal/user"
)

// Integration tests for order functionality

func TestOrderIntegration_BasicOrderOperations(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownSimpleTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&tenant.Tenant{}, &user.User{}, &Order{}, &OrderItem{})
	require.NoError(t, err)

	// Setup services
	tenantRepo := tenant.NewRepository(testDB.DB)
	tenantService := tenant.NewService(tenantRepo)
	userRepo := user.NewRepository(testDB.DB)
	orderRepo := NewRepository(testDB.DB)

	t.Run("Create and retrieve order", func(t *testing.T) {
		// Create tenant
		testTenant, err := tenantService.CreateTenant(tenant.CreateTenantRequest{
			Name:      "Order Test Store",
			Subdomain: "ordertest",
			Email:     "admin@ordertest.com",
		})
		require.NoError(t, err)

		// Create customer
		customer := &user.User{
			ID:        uuid.New(),
			TenantID:  &testTenant.ID,
			Email:     "customer@ordertest.com",
			Password:  "hashedpassword123",
			FirstName: "Test",
			LastName:  "Customer",
			Role:      user.RoleCustomer,
			Status:    user.StatusActive,
		}
		createdCustomer, err := userRepo.CreateUser(context.Background(), customer)
		require.NoError(t, err)

		// Create order manually
		order := &Order{
			ID:            uuid.New(),
			TenantID:      testTenant.ID,
			UserID:        createdCustomer.ID,
			OrderNumber:   "TEST-001",
			Status:        StatusPending,
			CustomerEmail: createdCustomer.Email,
			CustomerPhone: "+8801234567890",
			ShippingAddress: Address{
				FirstName:  "Test",
				LastName:   "Customer",
				Address1:   "123 Test Street",
				City:       "Dhaka",
				Country:    "Bangladesh",
				PostalCode: "1000",
			},
			SubtotalAmount:    199.98,
			ShippingAmount:    10.00,
			TotalAmount:       209.98,
			Currency:          "BDT",
			PaymentStatus:     PaymentPending,
			FulfillmentStatus: FulfillmentPending,
		}

		// Save order
		createdOrder, err := orderRepo.CreateOrder(order)
		require.NoError(t, err)
		assert.NotNil(t, createdOrder)
		assert.Equal(t, order.OrderNumber, createdOrder.OrderNumber)
		assert.Equal(t, order.TotalAmount, createdOrder.TotalAmount)

		// Retrieve order
		retrievedOrder, err := orderRepo.GetOrderByID(testTenant.ID, createdOrder.ID)
		require.NoError(t, err)
		assert.Equal(t, createdOrder.ID, retrievedOrder.ID)
		assert.Equal(t, createdOrder.TenantID, retrievedOrder.TenantID)
	})
}

func TestOrderIntegration_OrderItems(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownSimpleTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&tenant.Tenant{}, &user.User{}, &Order{}, &OrderItem{})
	require.NoError(t, err)

	// Setup repositories
	tenantRepo := tenant.NewRepository(testDB.DB)
	tenantService := tenant.NewService(tenantRepo)
	userRepo := user.NewRepository(testDB.DB)
	orderRepo := NewRepository(testDB.DB)

	t.Run("Order with multiple items", func(t *testing.T) {
		// Setup test data
		testTenant, _ := tenantService.CreateTenant(tenant.CreateTenantRequest{
			Name:      "Items Test Store",
			Subdomain: "itemstest",
			Email:     "admin@itemstest.com",
		})

		customer := &user.User{
			ID:        uuid.New(),
			TenantID:  &testTenant.ID,
			Email:     "customer@itemstest.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "Customer",
			Role:      user.RoleCustomer,
			Status:    user.StatusActive,
		}
		createdCustomer, _ := userRepo.CreateUser(context.Background(), customer)

		// Create order with items
		order := &Order{
			ID:            uuid.New(),
			TenantID:      testTenant.ID,
			UserID:        createdCustomer.ID,
			OrderNumber:   "ITEMS-001",
			Status:        StatusPending,
			CustomerEmail: createdCustomer.Email,
			ShippingAddress: Address{
				FirstName: "Test",
				LastName:  "Customer",
				Address1:  "123 Test Street",
				City:      "Dhaka",
				Country:   "Bangladesh",
			},
			SubtotalAmount:    299.97,
			TotalAmount:       299.97,
			Currency:          "BDT",
			PaymentStatus:     PaymentPending,
			FulfillmentStatus: FulfillmentPending,
			Items: []OrderItem{
				{
					ID:          uuid.New(),
					ProductID:   uuid.New(),
					ProductName: "Test Product 1",
					ProductSKU:  "TEST-001",
					Quantity:    2,
					UnitPrice:   99.99,
					TotalPrice:  199.98,
				},
				{
					ID:          uuid.New(),
					ProductID:   uuid.New(),
					ProductName: "Test Product 2",
					ProductSKU:  "TEST-002",
					Quantity:    1,
					UnitPrice:   99.99,
					TotalPrice:  99.99,
				},
			},
		}

		// Save order with items
		createdOrder, err := orderRepo.CreateOrder(order)
		require.NoError(t, err)
		assert.NotNil(t, createdOrder)
		assert.Len(t, createdOrder.Items, 2)

		// Verify order items were saved correctly
		retrievedOrder, err := orderRepo.GetOrderByID(testTenant.ID, createdOrder.ID)
		require.NoError(t, err)
		assert.Len(t, retrievedOrder.Items, 2)
		assert.Equal(t, "Test Product 1", retrievedOrder.Items[0].ProductName)
		assert.Equal(t, 2, retrievedOrder.Items[0].Quantity)
	})
}

func TestOrderIntegration_MultiTenantIsolation(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownSimpleTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&tenant.Tenant{}, &user.User{}, &Order{}, &OrderItem{})
	require.NoError(t, err)

	// Setup repositories
	tenantRepo := tenant.NewRepository(testDB.DB)
	tenantService := tenant.NewService(tenantRepo)
	userRepo := user.NewRepository(testDB.DB)
	orderRepo := NewRepository(testDB.DB)

	t.Run("Orders are isolated between tenants", func(t *testing.T) {
		// Create two tenants
		tenant1, _ := tenantService.CreateTenant(tenant.CreateTenantRequest{
			Name:      "Store 1",
			Subdomain: "store1orders",
			Email:     "admin@store1orders.com",
		})

		tenant2, _ := tenantService.CreateTenant(tenant.CreateTenantRequest{
			Name:      "Store 2",
			Subdomain: "store2orders",
			Email:     "admin@store2orders.com",
		})

		// Create customers for each tenant
		customer1 := &user.User{
			ID:        uuid.New(),
			TenantID:  &tenant1.ID,
			Email:     "customer1@store1.com",
			Password:  "password123",
			FirstName: "Customer",
			LastName:  "One",
			Role:      user.RoleCustomer,
			Status:    user.StatusActive,
		}
		createdCustomer1, _ := userRepo.CreateUser(context.Background(), customer1)

		customer2 := &user.User{
			ID:        uuid.New(),
			TenantID:  &tenant2.ID,
			Email:     "customer2@store2.com",
			Password:  "password123",
			FirstName: "Customer",
			LastName:  "Two",
			Role:      user.RoleCustomer,
			Status:    user.StatusActive,
		}
		createdCustomer2, _ := userRepo.CreateUser(context.Background(), customer2)

		// Create orders for each tenant
		order1 := &Order{
			ID:            uuid.New(),
			TenantID:      tenant1.ID,
			UserID:        createdCustomer1.ID,
			OrderNumber:   "T1-001",
			Status:        StatusPending,
			CustomerEmail: createdCustomer1.Email,
			ShippingAddress: Address{
				FirstName: "Customer",
				LastName:  "One",
				Address1:  "Address 1",
				City:      "Dhaka",
				Country:   "Bangladesh",
			},
			SubtotalAmount:    100.00,
			TotalAmount:       100.00,
			Currency:          "BDT",
			PaymentStatus:     PaymentPending,
			FulfillmentStatus: FulfillmentPending,
		}

		order2 := &Order{
			ID:            uuid.New(),
			TenantID:      tenant2.ID,
			UserID:        createdCustomer2.ID,
			OrderNumber:   "T2-001",
			Status:        StatusPending,
			CustomerEmail: createdCustomer2.Email,
			ShippingAddress: Address{
				FirstName: "Customer",
				LastName:  "Two",
				Address1:  "Address 2",
				City:      "Dhaka",
				Country:   "Bangladesh",
			},
			SubtotalAmount:    200.00,
			TotalAmount:       200.00,
			Currency:          "BDT",
			PaymentStatus:     PaymentPending,
			FulfillmentStatus: FulfillmentPending,
		}

		// Save orders
		createdOrder1, _ := orderRepo.CreateOrder(order1)
		createdOrder2, _ := orderRepo.CreateOrder(order2)

		// Verify tenant isolation - each tenant can only access their own orders
		retrievedOrder1, err := orderRepo.GetOrderByID(tenant1.ID, createdOrder1.ID)
		require.NoError(t, err)
		assert.Equal(t, createdOrder1.ID, retrievedOrder1.ID)

		retrievedOrder2, err := orderRepo.GetOrderByID(tenant2.ID, createdOrder2.ID)
		require.NoError(t, err)
		assert.Equal(t, createdOrder2.ID, retrievedOrder2.ID)

		// Cross-tenant access should fail or return empty
		_, err = orderRepo.GetOrderByID(tenant1.ID, createdOrder2.ID)
		assert.Error(t, err) // Should fail - order2 belongs to tenant2

		_, err = orderRepo.GetOrderByID(tenant2.ID, createdOrder1.ID)
		assert.Error(t, err) // Should fail - order1 belongs to tenant1
	})
}

func TestOrderIntegration_OrderStatusManagement(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownSimpleTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&tenant.Tenant{}, &user.User{}, &Order{}, &OrderItem{})
	require.NoError(t, err)

	// Setup repositories
	tenantRepo := tenant.NewRepository(testDB.DB)
	tenantService := tenant.NewService(tenantRepo)
	userRepo := user.NewRepository(testDB.DB)
	orderRepo := NewRepository(testDB.DB)

	t.Run("Order status lifecycle", func(t *testing.T) {
		// Setup test data
		testTenant, _ := tenantService.CreateTenant(tenant.CreateTenantRequest{
			Name:      "Status Test Store",
			Subdomain: "statustest",
			Email:     "admin@statustest.com",
		})

		customer := &user.User{
			ID:        uuid.New(),
			TenantID:  &testTenant.ID,
			Email:     "customer@statustest.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "Customer",
			Role:      user.RoleCustomer,
			Status:    user.StatusActive,
		}
		createdCustomer, _ := userRepo.CreateUser(context.Background(), customer)

		// Create order
		order := &Order{
			ID:            uuid.New(),
			TenantID:      testTenant.ID,
			UserID:        createdCustomer.ID,
			OrderNumber:   "STATUS-001",
			Status:        StatusPending,
			CustomerEmail: createdCustomer.Email,
			ShippingAddress: Address{
				FirstName: "Test",
				LastName:  "Customer",
				Address1:  "123 Test Street",
				City:      "Dhaka",
				Country:   "Bangladesh",
			},
			SubtotalAmount:    99.99,
			TotalAmount:       99.99,
			Currency:          "BDT",
			PaymentStatus:     PaymentPending,
			FulfillmentStatus: FulfillmentPending,
		}

		createdOrder, _ := orderRepo.CreateOrder(order)

		// Update status to confirmed
		err = orderRepo.UpdateOrderStatus(testTenant.ID, createdOrder.ID, StatusConfirmed)
		require.NoError(t, err)

		// Verify status update
		updatedOrder, err := orderRepo.GetOrderByID(testTenant.ID, createdOrder.ID)
		require.NoError(t, err)
		assert.Equal(t, StatusConfirmed, updatedOrder.Status)

		// Update payment status - need to check if this method exists
		err = orderRepo.UpdateOrderFields(testTenant.ID, createdOrder.ID, map[string]interface{}{
			"payment_status": PaymentPaid,
		})
		require.NoError(t, err)

		// Verify payment status update
		paidOrder, err := orderRepo.GetOrderByID(testTenant.ID, createdOrder.ID)
		require.NoError(t, err)
		assert.Equal(t, PaymentPaid, paidOrder.PaymentStatus)
	})
}

// Benchmark test for order creation performance
func BenchmarkOrderIntegration_CreateOrder(b *testing.B) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(&testing.T{})
	defer testDB.TeardownSimpleTestDatabase(&testing.T{})

	// Migrate schemas
	testDB.DB.AutoMigrate(&tenant.Tenant{}, &user.User{}, &Order{}, &OrderItem{})

	// Setup repositories
	tenantRepo := tenant.NewRepository(testDB.DB)
	tenantService := tenant.NewService(tenantRepo)
	userRepo := user.NewRepository(testDB.DB)
	orderRepo := NewRepository(testDB.DB)

	// Create test data once
	testTenant, _ := tenantService.CreateTenant(tenant.CreateTenantRequest{
		Name:      "Benchmark Store",
		Subdomain: "benchmarkstore",
		Email:     "admin@benchmark.com",
	})

	customer := &user.User{
		ID:        uuid.New(),
		TenantID:  &testTenant.ID,
		Email:     "customer@benchmark.com",
		Password:  "password123",
		FirstName: "Benchmark",
		LastName:  "Customer",
		Role:      user.RoleCustomer,
		Status:    user.StatusActive,
	}
	createdCustomer, _ := userRepo.CreateUser(context.Background(), customer)

	// Reset timer after setup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		order := &Order{
			ID:            uuid.New(),
			TenantID:      testTenant.ID,
			UserID:        createdCustomer.ID,
			OrderNumber:   fmt.Sprintf("BENCH-%d", i),
			Status:        StatusPending,
			CustomerEmail: createdCustomer.Email,
			ShippingAddress: Address{
				FirstName: "Benchmark",
				LastName:  "Customer",
				Address1:  "123 Benchmark Street",
				City:      "Dhaka",
				Country:   "Bangladesh",
			},
			SubtotalAmount:    99.99,
			TotalAmount:       99.99,
			Currency:          "BDT",
			PaymentStatus:     PaymentPending,
			FulfillmentStatus: FulfillmentPending,
		}

		_, err := orderRepo.CreateOrder(order)
		if err != nil {
			b.Fatal(err)
		}
	}
}
