package product

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
)

func TestProductValidation(t *testing.T) {
	tenantID := uuid.New()
	product := &Product{
		ID:       uuid.New(),
		TenantID: tenantID,
		Name:     "Test Product",
		Slug:     "test-product",
		Type:     TypePhysical,
		Status:   StatusActive,
		Price:    99.99,
	}

	// Test valid product
	if product.Name == "" {
		t.Error("Product should have name")
	}

	if product.Slug == "" {
		t.Error("Product should have slug")
	}

	if product.Type == "" {
		t.Error("Product should have type")
	}

	if product.Status == "" {
		t.Error("Product should have status")
	}

	if product.Price <= 0 {
		t.Error("Product should have valid price")
	}
}

func TestProductStatus(t *testing.T) {
	testCases := []struct {
		status   ProductStatus
		expected string
	}{
		{StatusDraft, "draft"},
		{StatusActive, "active"},
		{StatusInactive, "inactive"},
		{StatusArchived, "archived"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.status), func(t *testing.T) {
			if string(tc.status) != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, string(tc.status))
			}
		})
	}
}

func TestProductType(t *testing.T) {
	testCases := []struct {
		pType    ProductType
		expected string
	}{
		{TypePhysical, "physical"},
		{TypeDigital, "digital"},
		{TypeService, "service"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.pType), func(t *testing.T) {
			if string(tc.pType) != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, string(tc.pType))
			}
		})
	}
}

func TestProductCreation(t *testing.T) {
	tenantID := uuid.New()
	product := &Product{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Name:         "Test Product",
		Slug:         "test-product",
		Description:  "A test product",
		Type:         TypePhysical,
		Status:       StatusActive,
		Price:        99.99,
		ComparePrice: 120.00,
		CostPrice:    50.00,
	}

	if product.ID == uuid.Nil {
		t.Error("Product ID should not be nil")
	}

	if product.TenantID == uuid.Nil {
		t.Error("Tenant ID should not be nil")
	}

	if product.Name == "" {
		t.Error("Name should not be empty")
	}

	if product.Slug == "" {
		t.Error("Slug should not be empty")
	}

	if product.Type == "" {
		t.Error("Type should not be empty")
	}

	if product.Status == "" {
		t.Error("Status should not be empty")
	}

	if product.Price <= 0 {
		t.Error("Price should be positive")
	}

	// Test pricing logic
	if product.ComparePrice <= product.Price {
		t.Error("Compare price should be higher than selling price for discount")
	}

	if product.CostPrice >= product.Price {
		t.Error("Cost price should be lower than selling price for profit")
	}
}

func TestProductProfitCalculation(t *testing.T) {
	product := &Product{
		Price:     100.00,
		CostPrice: 60.00,
	}

	profit := product.Price - product.CostPrice
	expectedProfit := 40.00

	if profit != expectedProfit {
		t.Errorf("Expected profit %.2f, got %.2f", expectedProfit, profit)
	}

	// Test profit margin percentage
	margin := (profit / product.Price) * 100
	expectedMargin := 40.0

	if margin != expectedMargin {
		t.Errorf("Expected margin %.1f%%, got %.1f%%", expectedMargin, margin)
	}
}

func TestProductDiscountCalculation(t *testing.T) {
	product := &Product{
		Price:        80.00,
		ComparePrice: 100.00,
	}

	discount := product.ComparePrice - product.Price
	expectedDiscount := 20.00

	if discount != expectedDiscount {
		t.Errorf("Expected discount %.2f, got %.2f", expectedDiscount, discount)
	}

	// Test discount percentage
	discountPercent := (discount / product.ComparePrice) * 100
	expectedPercent := 20.0

	if discountPercent != expectedPercent {
		t.Errorf("Expected discount percent %.1f%%, got %.1f%%", expectedPercent, discountPercent)
	}
}

// Benchmark tests
func BenchmarkProductCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &Product{
			ID:          uuid.New(),
			TenantID:    uuid.New(),
			Name:        "Test Product",
			Slug:        "test-product",
			Description: "A test product",
			Type:        TypePhysical,
			Status:      StatusActive,
			Price:       99.99,
		}
	}
}

func BenchmarkProductProfitCalculation(b *testing.B) {
	product := &Product{
		Price:     100.00,
		CostPrice: 60.00,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = product.Price - product.CostPrice
	}
}

// Integration tests with real database
func createDefaultCategory(t *testing.T, repo Repository, tenantID uuid.UUID) *Category {
	defaultCategory := &Category{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        "Default Category",
		Slug:        "default-category-" + uuid.New().String()[:8], // Make slug unique
		Description: "Default category for testing",
		IsActive:    true,
	}

	createdCategory, err := repo.CreateCategory(defaultCategory)
	require.NoError(t, err)
	return createdCategory
}

func TestProductIntegration_ProductLifecycle(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(
		&Product{},
		&ProductVariant{},
		&Category{},
	)
	require.NoError(t, err)

	// Setup repository
	repo := NewRepository(testDB.DB)

	t.Run("Complete product lifecycle", func(t *testing.T) {
		tenantID := uuid.New()

		// Create a default category first to avoid foreign key constraint
		defaultCategory := &Category{
			ID:          uuid.New(),
			TenantID:    tenantID,
			Name:        "Default Category",
			Slug:        "default-category",
			Description: "Default category for testing",
			IsActive:    true,
		}

		createdCategory, err := repo.CreateCategory(defaultCategory)
		require.NoError(t, err)

		// Step 1: Create product
		product := &Product{
			ID:                uuid.New(),
			TenantID:          tenantID,
			Name:              "Integration Test Product",
			Slug:              "integration-test-product",
			Description:       "A product for integration testing",
			Type:              TypePhysical,
			Status:            StatusActive,
			Price:             99.99,
			ComparePrice:      129.99,
			CostPrice:         60.00,
			InventoryQuantity: 100,
			TrackQuantity:     true,
			SKU:               "TEST-001",
			Weight:            500.0,
			Length:            10.0,
			Width:             5.0,
			Height:            3.0,
			Tags:              []string{"test", "electronics", "gadget"},
			CategoryID:        createdCategory.ID, // Use valid category ID
		}

		createdProduct, err := repo.CreateProduct(product)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, createdProduct.ID)
		assert.Equal(t, product.Name, createdProduct.Name)
		assert.Equal(t, product.Status, createdProduct.Status)
		assert.Equal(t, product.Price, createdProduct.Price)

		// Step 2: Get product by ID
		retrievedProduct, err := repo.GetProductByID(tenantID, createdProduct.ID)
		require.NoError(t, err)
		assert.Equal(t, createdProduct.ID, retrievedProduct.ID)
		assert.Equal(t, createdProduct.Name, retrievedProduct.Name)
		assert.Equal(t, createdProduct.SKU, retrievedProduct.SKU)

		// Step 3: Get product by slug
		slugProduct, err := repo.GetProductBySlug(tenantID, product.Slug)
		require.NoError(t, err)
		assert.Equal(t, createdProduct.ID, slugProduct.ID)

		// Step 4: Update product
		retrievedProduct.Name = "Updated Product Name"
		retrievedProduct.Price = 119.99
		retrievedProduct.Status = StatusInactive
		retrievedProduct.InventoryQuantity = 75

		updatedProduct, err := repo.UpdateProduct(retrievedProduct)
		require.NoError(t, err)
		assert.Equal(t, "Updated Product Name", updatedProduct.Name)
		assert.Equal(t, 119.99, updatedProduct.Price)
		assert.Equal(t, StatusInactive, updatedProduct.Status)
		assert.Equal(t, 75, updatedProduct.InventoryQuantity)

		// Step 5: Check if product exists
		exists, err := repo.ProductExists(tenantID, createdProduct.ID)
		require.NoError(t, err)
		assert.True(t, exists)

		// Step 6: Update inventory
		err = repo.UpdateProductInventory(tenantID, createdProduct.ID, 50)
		require.NoError(t, err)

		// Verify inventory update
		inventoryProduct, err := repo.GetProductByID(tenantID, createdProduct.ID)
		require.NoError(t, err)
		assert.Equal(t, 50, inventoryProduct.InventoryQuantity)

		// Step 7: Delete product
		err = repo.DeleteProduct(tenantID, createdProduct.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetProductByID(tenantID, createdProduct.ID)
		assert.Error(t, err)
	})

	t.Run("Product category management", func(t *testing.T) {
		tenantID := uuid.New()

		// Create parent category
		parentCategory := &Category{
			ID:          uuid.New(),
			TenantID:    tenantID,
			Name:        "Electronics",
			Slug:        "electronics",
			Description: "Electronic products",
			IsActive:    true,
		}

		createdParent, err := repo.CreateCategory(parentCategory)
		require.NoError(t, err)

		// Create child category
		childCategory := &Category{
			ID:          uuid.New(),
			TenantID:    tenantID,
			ParentID:    &createdParent.ID,
			Name:        "Smartphones",
			Slug:        "smartphones",
			Description: "Mobile phones and smartphones",
			IsActive:    true,
		}

		createdChild, err := repo.CreateCategory(childCategory)
		require.NoError(t, err)

		// Create product with category
		product := &Product{
			ID:          uuid.New(),
			TenantID:    tenantID,
			Name:        "iPhone 15",
			Slug:        "iphone-15",
			Description: "Latest iPhone model",
			Type:        TypePhysical,
			Status:      StatusActive,
			Price:       999.99,
			CategoryID:  createdChild.ID,
		}

		createdProduct, err := repo.CreateProduct(product)
		require.NoError(t, err)

		// Get products by category
		categoryProducts, total, err := repo.GetProductsByCategoryID(tenantID, createdChild.ID, 0, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, categoryProducts, 1)
		assert.Equal(t, createdProduct.ID, categoryProducts[0].ID)

		// Get root categories
		rootCategories, err := repo.GetRootCategories(tenantID)
		require.NoError(t, err)
		assert.Len(t, rootCategories, 1)
		assert.Equal(t, createdParent.ID, rootCategories[0].ID)

		// Get category children
		children, err := repo.GetCategoryChildren(tenantID, createdParent.ID)
		require.NoError(t, err)
		assert.Len(t, children, 1)
		assert.Equal(t, createdChild.ID, children[0].ID)
	})

	t.Run("Product variants management", func(t *testing.T) {
		tenantID := uuid.New()

		// Create default category first
		defaultCategory := createDefaultCategory(t, repo, tenantID)

		// Create base product
		product := &Product{
			ID:          uuid.New(),
			TenantID:    tenantID,
			Name:        "T-Shirt",
			Slug:        "t-shirt",
			Description: "Cotton t-shirt",
			Type:        TypePhysical,
			Status:      StatusActive,
			Price:       25.99,
			CategoryID:  defaultCategory.ID, // Use valid category ID
		}

		createdProduct, err := repo.CreateProduct(product)
		require.NoError(t, err)

		// Create variants
		variants := []*ProductVariant{
			{
				ID:        uuid.New(),
				TenantID:  tenantID,
				ProductID: createdProduct.ID,
				Name:      "Small Red",
				SKU:       "TSHIRT-SM-RED",
				Price:     25.99,
				Options: map[string]string{
					"size":  "Small",
					"color": "Red",
				},
				InventoryQuantity: 50,
				IsDefault:         true,
			},
			{
				ID:        uuid.New(),
				TenantID:  tenantID,
				ProductID: createdProduct.ID,
				Name:      "Large Blue",
				SKU:       "TSHIRT-LG-BLUE",
				Price:     27.99,
				Options: map[string]string{
					"size":  "Large",
					"color": "Blue",
				},
				InventoryQuantity: 30,
				IsDefault:         false,
			},
		}

		// Create variants
		for _, variant := range variants {
			createdVariant, err := repo.CreateProductVariant(variant)
			require.NoError(t, err)
			assert.Equal(t, variant.Name, createdVariant.Name)
			assert.Equal(t, variant.SKU, createdVariant.SKU)
		}

		// Get all variants for product
		productVariants, err := repo.GetProductVariants(tenantID, createdProduct.ID)
		require.NoError(t, err)
		assert.Len(t, productVariants, 2)

		// Get specific variant
		specificVariant, err := repo.GetProductVariant(tenantID, variants[0].ID)
		require.NoError(t, err)
		assert.Equal(t, variants[0].Name, specificVariant.Name)
		assert.True(t, specificVariant.IsDefault)

		// Update variant
		specificVariant.Price = 24.99
		specificVariant.InventoryQuantity = 45
		updatedVariant, err := repo.UpdateProductVariant(specificVariant)
		require.NoError(t, err)
		assert.Equal(t, 24.99, updatedVariant.Price)
		assert.Equal(t, 45, updatedVariant.InventoryQuantity)

		// Delete variant
		err = repo.DeleteProductVariant(tenantID, variants[1].ID)
		require.NoError(t, err)

		// Verify deletion
		remainingVariants, err := repo.GetProductVariants(tenantID, createdProduct.ID)
		require.NoError(t, err)
		assert.Len(t, remainingVariants, 1)
	})

	t.Run("Product search and filtering", func(t *testing.T) {
		tenantID := uuid.New()

		// Create default category first
		defaultCategory := createDefaultCategory(t, repo, tenantID)

		// Create test products
		products := []*Product{
			{
				ID:          uuid.New(),
				TenantID:    tenantID,
				Name:        "Apple iPhone 15",
				Slug:        "apple-iphone-15",
				Description: "Latest Apple smartphone",
				Type:        TypePhysical,
				Status:      StatusActive,
				Price:       999.99,
				Tags:        []string{"apple", "smartphone", "tech"},
				CategoryID:  defaultCategory.ID, // Use valid category ID
			},
			{
				ID:          uuid.New(),
				TenantID:    tenantID,
				Name:        "Samsung Galaxy S24",
				Slug:        "samsung-galaxy-s24",
				Description: "Android smartphone by Samsung",
				Type:        TypePhysical,
				Status:      StatusActive,
				Price:       849.99,
				Tags:        []string{"samsung", "android", "smartphone"},
				CategoryID:  defaultCategory.ID, // Use valid category ID
			},
			{
				ID:          uuid.New(),
				TenantID:    tenantID,
				Name:        "Apple MacBook Pro",
				Slug:        "apple-macbook-pro",
				Description: "Professional laptop by Apple",
				Type:        TypePhysical,
				Status:      StatusDraft,
				Price:       1999.99,
				Tags:        []string{"apple", "laptop", "professional"},
				CategoryID:  defaultCategory.ID, // Use valid category ID
			},
		}

		// Create all products
		for _, product := range products {
			_, err := repo.CreateProduct(product)
			require.NoError(t, err)
		}

		// Test search functionality
		searchResults, total, err := repo.SearchProducts(tenantID, "apple", 0, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total) // Should find iPhone and MacBook
		assert.Len(t, searchResults, 2)

		// Test list products with filter
		filter := ProductListFilter{
			Status: StatusActive,
		}
		activeProducts, activeTotal, err := repo.ListProducts(tenantID, filter, 0, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), activeTotal) // Only iPhone and Galaxy are active
		assert.Len(t, activeProducts, 2)

		// Test price range filtering
		minPrice := 800.0
		maxPrice := 1000.0
		filter = ProductListFilter{
			MinPrice: &minPrice,
			MaxPrice: &maxPrice,
		}
		priceRangeProducts, priceTotal, err := repo.ListProducts(tenantID, filter, 0, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), priceTotal) // iPhone and Galaxy
		assert.Len(t, priceRangeProducts, 2)
	})

	t.Run("Multi-tenant isolation", func(t *testing.T) {
		// Create products for different tenants
		tenant1 := uuid.New()
		tenant2 := uuid.New()

		// Create default categories for both tenants
		defaultCategory1 := createDefaultCategory(t, repo, tenant1)
		defaultCategory2 := createDefaultCategory(t, repo, tenant2)

		product1 := &Product{
			ID:          uuid.New(),
			TenantID:    tenant1,
			Name:        "Tenant 1 Product",
			Slug:        "tenant-1-product",
			Description: "Product for tenant 1",
			Type:        TypePhysical,
			Status:      StatusActive,
			Price:       50.00,
			CategoryID:  defaultCategory1.ID, // Use valid category ID
		}

		product2 := &Product{
			ID:          uuid.New(),
			TenantID:    tenant2,
			Name:        "Tenant 2 Product",
			Slug:        "tenant-2-product",
			Description: "Product for tenant 2",
			Type:        TypePhysical,
			Status:      StatusActive,
			Price:       75.00,
			CategoryID:  defaultCategory2.ID, // Use valid category ID
		}

		// Create products
		createdProduct1, err := repo.CreateProduct(product1)
		require.NoError(t, err)

		createdProduct2, err := repo.CreateProduct(product2)
		require.NoError(t, err)

		// Test tenant isolation in product listing
		filter := ProductListFilter{}

		// Get products for tenant1
		tenant1Products, total1, err := repo.ListProducts(tenant1, filter, 0, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total1)
		assert.Len(t, tenant1Products, 1)
		assert.Equal(t, createdProduct1.ID, tenant1Products[0].ID)

		// Get products for tenant2
		tenant2Products, total2, err := repo.ListProducts(tenant2, filter, 0, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total2)
		assert.Len(t, tenant2Products, 1)
		assert.Equal(t, createdProduct2.ID, tenant2Products[0].ID)

		// Verify cross-tenant access fails
		_, err = repo.GetProductByID(tenant1, createdProduct2.ID)
		assert.Error(t, err) // Should not find product from different tenant

		_, err = repo.GetProductByID(tenant2, createdProduct1.ID)
		assert.Error(t, err) // Should not find product from different tenant
	})

	t.Run("Low stock and bulk operations", func(t *testing.T) {
		tenantID := uuid.New()

		// Create default category first
		defaultCategory := createDefaultCategory(t, repo, tenantID)

		// Create products with different stock levels
		products := []*Product{
			{
				ID:                uuid.New(),
				TenantID:          tenantID,
				Name:              "High Stock Product",
				Slug:              "high-stock-product",
				Type:              TypePhysical,
				Status:            StatusActive,
				Price:             100.00,
				InventoryQuantity: 100,
				CategoryID:        defaultCategory.ID, // Use valid category ID
			},
			{
				ID:                uuid.New(),
				TenantID:          tenantID,
				Name:              "Low Stock Product",
				Slug:              "low-stock-product",
				Type:              TypePhysical,
				Status:            StatusActive,
				Price:             150.00,
				InventoryQuantity: 5,
				CategoryID:        defaultCategory.ID, // Use valid category ID
			},
			{
				ID:                uuid.New(),
				TenantID:          tenantID,
				Name:              "Out of Stock Product",
				Slug:              "out-of-stock-product",
				Type:              TypePhysical,
				Status:            StatusActive,
				Price:             200.00,
				InventoryQuantity: 0,
				CategoryID:        defaultCategory.ID, // Use valid category ID
			},
		}

		var productIDs []uuid.UUID
		// Create all products
		for _, product := range products {
			createdProduct, err := repo.CreateProduct(product)
			require.NoError(t, err)
			productIDs = append(productIDs, createdProduct.ID)
		}

		// Test low stock detection
		lowStockProducts, err := repo.GetLowStockProducts(tenantID, 10)
		require.NoError(t, err)
		assert.Len(t, lowStockProducts, 2) // Low stock and out of stock products

		// Test bulk update
		updates := map[string]interface{}{
			"status": StatusInactive,
		}
		err = repo.BulkUpdateProducts(tenantID, productIDs[:2], updates)
		require.NoError(t, err)

		// Verify bulk update
		for i := 0; i < 2; i++ {
			product, err := repo.GetProductByID(tenantID, productIDs[i])
			require.NoError(t, err)
			assert.Equal(t, StatusInactive, product.Status)
		}

		// Test bulk delete
		err = repo.BulkDeleteProducts(tenantID, productIDs)
		require.NoError(t, err)

		// Verify bulk deletion
		for _, productID := range productIDs {
			_, err := repo.GetProductByID(tenantID, productID)
			assert.Error(t, err)
		}
	})

	t.Run("Product statistics", func(t *testing.T) {
		tenantID := uuid.New()

		// Create default category first
		defaultCategory := createDefaultCategory(t, repo, tenantID)

		// Create products for statistics
		products := []*Product{
			{
				ID:                uuid.New(),
				TenantID:          tenantID,
				Name:              "Active Product 1",
				Slug:              "active-product-1",
				Type:              TypePhysical,
				Status:            StatusActive,
				Price:             100.00,
				InventoryQuantity: 10,                 // Add inventory for value calculation
				CategoryID:        defaultCategory.ID, // Use valid category ID
			},
			{
				ID:                uuid.New(),
				TenantID:          tenantID,
				Name:              "Active Product 2",
				Slug:              "active-product-2",
				Type:              TypeDigital,
				Status:            StatusActive,
				Price:             50.00,
				InventoryQuantity: 5,                  // Add inventory for value calculation
				CategoryID:        defaultCategory.ID, // Use valid category ID
			},
			{
				ID:                uuid.New(),
				TenantID:          tenantID,
				Name:              "Draft Product",
				Slug:              "draft-product",
				Type:              TypePhysical,
				Status:            StatusDraft,
				Price:             75.00,
				InventoryQuantity: 8,                  // Add inventory for value calculation
				CategoryID:        defaultCategory.ID, // Use valid category ID
			},
		}

		// Create all products
		for _, product := range products {
			_, err := repo.CreateProduct(product)
			require.NoError(t, err)
		}

		// Get product statistics
		stats, err := repo.GetProductStats(tenantID)
		require.NoError(t, err)

		assert.Equal(t, int64(3), stats.TotalProducts)
		assert.Equal(t, int64(2), stats.ActiveProducts)
		assert.Equal(t, int64(1), stats.DraftProducts)
		// Expected total value: (100*10) + (50*5) + (75*8) = 1000 + 250 + 600 = 1850
		assert.Equal(t, 1850.0, stats.TotalValue)
	})
}
