package cart

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveCart(cart *Cart) (*Cart, error) {
	args := m.Called(cart)
	return args.Get(0).(*Cart), args.Error(1)
}

func (m *MockRepository) FindCartByID(tenantID, cartID uuid.UUID) (*Cart, error) {
	args := m.Called(tenantID, cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Cart), args.Error(1)
}

func (m *MockRepository) FindCartByCustomerID(tenantID, customerID uuid.UUID) (*Cart, error) {
	args := m.Called(tenantID, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Cart), args.Error(1)
}

func (m *MockRepository) FindCartBySessionID(tenantID uuid.UUID, sessionID string) (*Cart, error) {
	args := m.Called(tenantID, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Cart), args.Error(1)
}

func (m *MockRepository) UpdateCart(cart *Cart) (*Cart, error) {
	args := m.Called(cart)
	return args.Get(0).(*Cart), args.Error(1)
}

func (m *MockRepository) DeleteCart(tenantID, cartID uuid.UUID) error {
	args := m.Called(tenantID, cartID)
	return args.Error(0)
}

func (m *MockRepository) ListCarts(tenantID uuid.UUID, filter CartListFilter, offset, limit int) ([]*Cart, int64, error) {
	args := m.Called(tenantID, filter, offset, limit)
	return args.Get(0).([]*Cart), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetAbandonedCarts(tenantID uuid.UUID, since time.Time) ([]*Cart, error) {
	args := m.Called(tenantID, since)
	return args.Get(0).([]*Cart), args.Error(1)
}

func (m *MockRepository) GetExpiredCarts(tenantID uuid.UUID) ([]*Cart, error) {
	args := m.Called(tenantID)
	return args.Get(0).([]*Cart), args.Error(1)
}

func (m *MockRepository) CleanupExpiredCarts(tenantID uuid.UUID) error {
	args := m.Called(tenantID)
	return args.Error(0)
}

func (m *MockRepository) MergeGuestCartToCustomer(tenantID uuid.UUID, sessionID string, customerID uuid.UUID) error {
	args := m.Called(tenantID, sessionID, customerID)
	return args.Error(0)
}

func (m *MockRepository) AddCartItem(item *CartItem) (*CartItem, error) {
	args := m.Called(item)
	return args.Get(0).(*CartItem), args.Error(1)
}

func (m *MockRepository) UpdateCartItem(item *CartItem) (*CartItem, error) {
	args := m.Called(item)
	return args.Get(0).(*CartItem), args.Error(1)
}

func (m *MockRepository) RemoveCartItem(tenantID, cartID, itemID uuid.UUID) error {
	args := m.Called(tenantID, cartID, itemID)
	return args.Error(0)
}

func (m *MockRepository) FindCartItem(tenantID, cartID, itemID uuid.UUID) (*CartItem, error) {
	args := m.Called(tenantID, cartID, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CartItem), args.Error(1)
}

func (m *MockRepository) ClearCartItems(tenantID, cartID uuid.UUID) error {
	args := m.Called(tenantID, cartID)
	return args.Error(0)
}

func (m *MockRepository) GetCartStats(tenantID uuid.UUID) (*CartStats, error) {
	args := m.Called(tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CartStats), args.Error(1)
}

func (m *MockRepository) GetAbandonmentRate(tenantID uuid.UUID, days int) (float64, error) {
	args := m.Called(tenantID, days)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRepository) GetAverageCartValue(tenantID uuid.UUID, days int) (float64, error) {
	args := m.Called(tenantID, days)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRepository) GetTopAbandonedProducts(tenantID uuid.UUID, limit int) ([]*AbandonedProductStats, error) {
	args := m.Called(tenantID, limit)
	return args.Get(0).([]*AbandonedProductStats), args.Error(1)
}

// MockProductService mocks the product service dependency
type MockProductService struct {
	mock.Mock
}

// MockDiscountService mocks the discount service dependency
type MockDiscountService struct {
	mock.Mock
}

// MockShippingService mocks the shipping service dependency
type MockShippingService struct {
	mock.Mock
}

// Integration tests with real database
func TestCartIntegration_CartLifecycle(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&Cart{}, &CartItem{})
	require.NoError(t, err)

	// Setup repository
	repo := NewRepository(testDB.DB)

	t.Run("Complete cart lifecycle", func(t *testing.T) {
		// Step 1: Create cart
		tenantID := uuid.New()
		customerID := uuid.New()

		cart := &Cart{
			ID:         uuid.New(),
			TenantID:   tenantID,
			CustomerID: &customerID,
			Status:     StatusActive,
			Currency:   "USD",
			Subtotal:   0.00,
			Total:      0.00,
		}

		savedCart, err := repo.SaveCart(cart)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, savedCart.ID)
		assert.Equal(t, StatusActive, savedCart.Status)

		// Step 2: Add cart item
		productID := uuid.New()
		item := &CartItem{
			ID:          uuid.New(),
			CartID:      savedCart.ID,
			ProductID:   productID,
			ProductName: "Test Product",
			SKU:         "TEST-001",
			Price:       25.99,
			Quantity:    2,
			LineTotal:   51.98,
		}

		savedItem, err := repo.AddCartItem(item)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, savedItem.ID)
		assert.Equal(t, 2, savedItem.Quantity)

		// Step 3: Update cart totals
		savedCart.Subtotal = 51.98
		savedCart.Total = 51.98
		updatedCart, err := repo.UpdateCart(savedCart)
		require.NoError(t, err)
		assert.Equal(t, 51.98, updatedCart.Total)

		// Step 4: Update item quantity
		savedItem.Quantity = 3
		savedItem.LineTotal = 77.97
		updatedItem, err := repo.UpdateCartItem(savedItem)
		require.NoError(t, err)
		assert.Equal(t, 3, updatedItem.Quantity)
		assert.Equal(t, 77.97, updatedItem.LineTotal)

		// Step 5: Get cart by ID
		retrievedCart, err := repo.FindCartByID(tenantID, savedCart.ID)
		require.NoError(t, err)
		assert.Equal(t, savedCart.ID, retrievedCart.ID)
		assert.Equal(t, StatusActive, retrievedCart.Status)

		// Step 6: Get cart by customer
		customerCart, err := repo.FindCartByCustomerID(tenantID, customerID)
		require.NoError(t, err)
		assert.Equal(t, savedCart.ID, customerCart.ID)

		// Step 7: Remove cart item
		err = repo.RemoveCartItem(tenantID, savedCart.ID, savedItem.ID)
		require.NoError(t, err)

		// Step 8: Clear cart
		err = repo.ClearCartItems(tenantID, savedCart.ID)
		require.NoError(t, err)

		// Step 9: Delete cart
		err = repo.DeleteCart(tenantID, savedCart.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.FindCartByID(tenantID, savedCart.ID)
		assert.Error(t, err)
	})

	t.Run("Multi-tenant isolation", func(t *testing.T) {
		// Create carts for different tenants
		tenant1 := uuid.New()
		tenant2 := uuid.New()
		customer1 := uuid.New()
		customer2 := uuid.New()

		cart1 := &Cart{
			ID:         uuid.New(),
			TenantID:   tenant1,
			CustomerID: &customer1,
			Status:     StatusActive,
			Currency:   "USD",
			Total:      100.00,
		}

		cart2 := &Cart{
			ID:         uuid.New(),
			TenantID:   tenant2,
			CustomerID: &customer2,
			Status:     StatusActive,
			Currency:   "USD",
			Total:      200.00,
		}

		// Save carts
		savedCart1, err := repo.SaveCart(cart1)
		require.NoError(t, err)

		savedCart2, err := repo.SaveCart(cart2)
		require.NoError(t, err)

		// Verify tenant isolation
		tenant1Cart, err := repo.FindCartByID(tenant1, savedCart1.ID)
		require.NoError(t, err)
		assert.Equal(t, savedCart1.ID, tenant1Cart.ID)

		tenant2Cart, err := repo.FindCartByID(tenant2, savedCart2.ID)
		require.NoError(t, err)
		assert.Equal(t, savedCart2.ID, tenant2Cart.ID)

		// Verify cross-tenant access is blocked
		_, err = repo.FindCartByID(tenant1, savedCart2.ID)
		assert.Error(t, err) // Should not find cart2 in tenant1

		_, err = repo.FindCartByID(tenant2, savedCart1.ID)
		assert.Error(t, err) // Should not find cart1 in tenant2
	})

	t.Run("Guest cart handling", func(t *testing.T) {
		tenantID := uuid.New()
		sessionID := "guest_session_12345"

		// Create guest cart
		guestCart := &Cart{
			ID:        uuid.New(),
			TenantID:  tenantID,
			SessionID: sessionID,
			Status:    StatusActive,
			Currency:  "USD",
			Total:     75.50,
		}

		savedGuestCart, err := repo.SaveCart(guestCart)
		require.NoError(t, err)

		// Find cart by session ID
		foundCart, err := repo.FindCartBySessionID(tenantID, sessionID)
		require.NoError(t, err)
		assert.Equal(t, savedGuestCart.ID, foundCart.ID)
		assert.Equal(t, sessionID, foundCart.SessionID)

		// Convert guest cart to customer cart
		customerID := uuid.New()
		err = repo.MergeGuestCartToCustomer(tenantID, sessionID, customerID)
		require.NoError(t, err)

		// Verify conversion
		customerCart, err := repo.FindCartByCustomerID(tenantID, customerID)
		require.NoError(t, err)
		assert.Equal(t, savedGuestCart.ID, customerCart.ID)
		assert.Equal(t, &customerID, customerCart.CustomerID)
	})

	t.Run("Cart status management", func(t *testing.T) {
		tenantID := uuid.New()
		customerID := uuid.New()

		// Test different cart statuses
		statuses := []CartStatus{StatusActive, StatusAbandoned, StatusConverted, StatusExpired}

		for _, status := range statuses {
			cart := &Cart{
				ID:         uuid.New(),
				TenantID:   tenantID,
				CustomerID: &customerID,
				Status:     status,
				Currency:   "USD",
				Total:      50.00,
			}

			savedCart, err := repo.SaveCart(cart)
			require.NoError(t, err)
			assert.Equal(t, status, savedCart.Status)

			// Verify status in retrieval
			retrievedCart, err := repo.FindCartByID(tenantID, savedCart.ID)
			require.NoError(t, err)
			assert.Equal(t, status, retrievedCart.Status)
		}
	})
}
