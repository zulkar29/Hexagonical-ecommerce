package category

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
)

// Integration tests for category service - critical for product organization

func TestCategoryService_CategoryCRUD(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&Category{})
	require.NoError(t, err)

	// Setup services
	categoryRepo := NewGormRepository(testDB.DB)
	categoryService := NewService(categoryRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Create and retrieve category", func(t *testing.T) {
		// Create root category
		createReq := CreateCategoryRequest{
			Name:        "Electronics",
			Description: "Electronic devices and gadgets",
			Image:       "/images/electronics.jpg",
			Icon:        "fas fa-laptop",
			SortOrder:   1,
			IsFeatured:  true,
			ShowInMenu:  true,
			MetaTitle:   "Electronics - Best Devices",
			MetaDescription: "Shop the latest electronic devices",
			MetaKeywords: "electronics, devices, gadgets",
		}

		category, err := categoryService.CreateCategory(ctx, tenantID, createReq)
		require.NoError(t, err)
		assert.Equal(t, createReq.Name, category.Name)
		assert.Equal(t, "electronics", category.Slug) // Should be auto-generated
		assert.Equal(t, createReq.Description, category.Description)
		assert.Equal(t, createReq.Image, category.Image)
		assert.Equal(t, createReq.Icon, category.Icon)
		assert.Equal(t, 0, category.Level) // Root category
		assert.Equal(t, "/electronics", category.Path)
		assert.Equal(t, StatusActive, category.Status)
		assert.True(t, category.IsFeatured)

		// Get category by ID
		retrievedCategory, err := categoryService.GetCategory(ctx, tenantID, category.ID)
		require.NoError(t, err)
		assert.Equal(t, category.ID, retrievedCategory.ID)
		assert.Equal(t, category.Name, retrievedCategory.Name)

		// Get category by slug
		categoryBySlug, err := categoryService.GetCategoryBySlug(ctx, tenantID, "electronics")
		require.NoError(t, err)
		assert.Equal(t, category.ID, categoryBySlug.ID)
		assert.Equal(t, category.Name, categoryBySlug.Name)
	})

	t.Run("Update category", func(t *testing.T) {
		// Create category first
		createReq := CreateCategoryRequest{
			Name:        "Clothing",
			Description: "Fashion and apparel",
		}

		category, err := categoryService.CreateCategory(ctx, tenantID, createReq)
		require.NoError(t, err)

		// Update category
		updateReq := UpdateCategoryRequest{
			Name:        "Fashion & Clothing",
			Description: "Latest fashion trends and clothing",
			Image:       "/images/fashion.jpg",
			IsFeatured:  &[]bool{true}[0],
			MetaTitle:   "Fashion Store - Latest Trends",
		}

		updatedCategory, err := categoryService.UpdateCategory(ctx, tenantID, category.ID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, updateReq.Name, updatedCategory.Name)
		assert.Equal(t, "fashion-clothing", updatedCategory.Slug) // Should be updated
		assert.Equal(t, updateReq.Description, updatedCategory.Description)
		assert.Equal(t, updateReq.Image, updatedCategory.Image)
		assert.True(t, updatedCategory.IsFeatured)
		assert.Equal(t, updateReq.MetaTitle, updatedCategory.MetaTitle)
	})

	t.Run("Delete category", func(t *testing.T) {
		// Create category first
		createReq := CreateCategoryRequest{
			Name:   "Temporary Category",
		}

		category, err := categoryService.CreateCategory(ctx, tenantID, createReq)
		require.NoError(t, err)

		// Delete category
		err = categoryService.DeleteCategory(ctx, tenantID, category.ID)
		require.NoError(t, err)

		// Verify deletion - should return error when trying to get deleted category
		_, err = categoryService.GetCategory(ctx, tenantID, category.ID)
		assert.Error(t, err)
	})
}

func TestCategoryService_HierarchicalCategories(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&Category{})
	require.NoError(t, err)

	// Setup services
	categoryRepo := NewGormRepository(testDB.DB)
	categoryService := NewService(categoryRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Create hierarchical categories", func(t *testing.T) {
		// Create root category
		rootReq := CreateCategoryRequest{
			Name:       "Electronics",
			SortOrder:  1,
			ShowInMenu: true,
		}

		rootCategory, err := categoryService.CreateCategory(ctx, tenantID, rootReq)
		require.NoError(t, err)
		assert.Equal(t, 0, rootCategory.Level)
		assert.Equal(t, "/electronics", rootCategory.Path)

		// Create subcategory
		subReq := CreateCategoryRequest{
			Name:      "Computers",
			ParentID:  &rootCategory.ID,
			SortOrder: 1,
		}

		subCategory, err := categoryService.CreateCategory(ctx, tenantID, subReq)
		require.NoError(t, err)
		assert.Equal(t, rootCategory.ID, *subCategory.ParentID)
		assert.Equal(t, 1, subCategory.Level)
		assert.Equal(t, "/electronics/computers", subCategory.Path)

		// Create sub-subcategory
		subSubReq := CreateCategoryRequest{
			Name:      "Laptops",
			ParentID:  &subCategory.ID,
			SortOrder: 1,
		}

		subSubCategory, err := categoryService.CreateCategory(ctx, tenantID, subSubReq)
		require.NoError(t, err)
		assert.Equal(t, subCategory.ID, *subSubCategory.ParentID)
		assert.Equal(t, 2, subSubCategory.Level)
		assert.Equal(t, "/electronics/computers/laptops", subSubCategory.Path)

		// Get category tree
		tree, err := categoryService.GetCategoryTree(ctx, tenantID, nil)
		require.NoError(t, err)
		assert.NotEmpty(t, tree)

		// Get category path
		path, err := categoryService.GetCategoryPath(ctx, tenantID, subSubCategory.ID)
		require.NoError(t, err)
		assert.Len(t, path, 3) // Root -> Sub -> SubSub
		assert.Equal(t, rootCategory.Name, path[0].Name)
		assert.Equal(t, subCategory.Name, path[1].Name)
		assert.Equal(t, subSubCategory.Name, path[2].Name)
	})

	t.Run("Move category to different parent", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create two root categories
		root1Req := CreateCategoryRequest{
			Name:   "Electronics",
		}
		root1, err := categoryService.CreateCategory(ctx, tenantID, root1Req)
		require.NoError(t, err)

		root2Req := CreateCategoryRequest{
			Name:   "Home & Garden",
		}
		root2, err := categoryService.CreateCategory(ctx, tenantID, root2Req)
		require.NoError(t, err)

		// Create subcategory under root1
		subReq := CreateCategoryRequest{
			Name:     "Smart Devices",
			ParentID: &root1.ID,
		}
		subCategory, err := categoryService.CreateCategory(ctx, tenantID, subReq)
		require.NoError(t, err)
		assert.Equal(t, "/electronics/smart-devices", subCategory.Path)

		// Move subcategory to root2
		err = categoryService.MoveCategory(ctx, tenantID, subCategory.ID, &root2.ID)
		require.NoError(t, err)

		// Verify the move
		movedCategory, err := categoryService.GetCategory(ctx, tenantID, subCategory.ID)
		require.NoError(t, err)
		assert.Equal(t, root2.ID, *movedCategory.ParentID)
		assert.Equal(t, "/home-garden/smart-devices", movedCategory.Path)
	})
}

func TestCategoryService_CategoryFiltering(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&Category{})
	require.NoError(t, err)

	// Setup services
	categoryRepo := NewGormRepository(testDB.DB)
	categoryService := NewService(categoryRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("List and filter categories", func(t *testing.T) {
		// Create categories with different properties
		categories := []CreateCategoryRequest{
			{
				Name:       "Electronics",
				IsFeatured: true,
				ShowInMenu: true,
				SortOrder:  1,
			},
			{
				Name:       "Clothing",
				IsFeatured: false,
				ShowInMenu: true,
				SortOrder:  2,
			},
			{
				Name:       "Books",
				IsFeatured: true,
				ShowInMenu: false,
				SortOrder:  3,
			},
			{
				Name:       "Archived Category",
				IsFeatured: false,
				ShowInMenu: false,
				SortOrder:  4,
			},
		}

		for _, req := range categories {
			_, err := categoryService.CreateCategory(ctx, tenantID, req)
			require.NoError(t, err)
		}

		// List all categories
		allCategories, total, err := categoryService.ListCategories(ctx, tenantID, CategoryFilter{}, 10, 0)
		require.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Len(t, allCategories, 4)

		// Filter by status - active only
		activeFilter := CategoryFilter{
			Status: StatusActive,
		}
		activeCategories, activeTotal, err := categoryService.ListCategories(ctx, tenantID, activeFilter, 10, 0)
		require.NoError(t, err)
		assert.Equal(t, int64(4), activeTotal)
		assert.Len(t, activeCategories, 4)

		// Debug: Check what was actually created
		allCategoriesDebug, _, err := categoryService.ListCategories(ctx, tenantID, CategoryFilter{}, 10, 0)
		require.NoError(t, err)
		for i, cat := range allCategoriesDebug {
			t.Logf("Category %d: Name=%s, IsFeatured=%v", i, cat.Name, cat.IsFeatured)
		}

		// Filter by featured categories
		featuredFilter := CategoryFilter{
			IsFeatured: &[]bool{true}[0],
		}
		featuredCategories, featuredTotal, err := categoryService.ListCategories(ctx, tenantID, featuredFilter, 10, 0)
		require.NoError(t, err)
		t.Logf("Featured categories found: %d", featuredTotal)
		for i, cat := range featuredCategories {
			t.Logf("Featured Category %d: Name=%s, IsFeatured=%v", i, cat.Name, cat.IsFeatured)
		}
		t.Logf("DEBUG: About to assert featuredCategories length. Current length: %d", len(featuredCategories))
		t.Logf("DEBUG: featuredCategories content: %+v", featuredCategories)
		assert.Equal(t, int64(2), featuredTotal)
		assert.Len(t, featuredCategories, 2)

		// Debug: Check all categories in database first
		allCategoriesFilter := CategoryFilter{}
		allCategories, allTotal, err := categoryService.ListCategories(ctx, tenantID, allCategoriesFilter, 10, 0)
		require.NoError(t, err)
		t.Logf("=== ALL CATEGORIES IN DATABASE (Total: %d) ===", allTotal)
		for i, cat := range allCategories {
			t.Logf("Category %d: Name=%s, ShowInMenu=%t, IsFeatured=%t", i, cat.Name, cat.ShowInMenu, cat.IsFeatured)
		}
		t.Logf("=== END ALL CATEGORIES ===")

		// Filter by show in menu
		menuFilter := CategoryFilter{
			ShowInMenu: &[]bool{true}[0],
		}
		menuCategories, menuTotal, err := categoryService.ListCategories(ctx, tenantID, menuFilter, 10, 0)
		require.NoError(t, err)
		t.Logf("DEBUG: Menu filter - menuTotal: %d, menuCategories length: %d", menuTotal, len(menuCategories))
		for i, cat := range menuCategories {
			t.Logf("DEBUG: Menu Category %d: Name=%s, ShowInMenu=%v", i, cat.Name, cat.ShowInMenu)
		}
		assert.Equal(t, int64(2), menuTotal)
		assert.Len(t, menuCategories, 2)

		// Test pagination
		firstPage, _, err := categoryService.ListCategories(ctx, tenantID, CategoryFilter{}, 2, 0)
		require.NoError(t, err)
		assert.Len(t, firstPage, 2)

		secondPage, _, err := categoryService.ListCategories(ctx, tenantID, CategoryFilter{}, 2, 2)
		require.NoError(t, err)
		assert.Len(t, secondPage, 2)

		// Verify different items on different pages
		assert.NotEqual(t, firstPage[0].ID, secondPage[0].ID)
	})
}

func TestCategoryService_CategoryOrdering(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&Category{})
	require.NoError(t, err)

	// Setup services
	categoryRepo := NewGormRepository(testDB.DB)
	categoryService := NewService(categoryRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Reorder categories", func(t *testing.T) {
		// Create categories with different sort orders
		categories := []CreateCategoryRequest{
			{Name: "Category A", SortOrder: 1},
			{Name: "Category B", SortOrder: 2},
			{Name: "Category C", SortOrder: 3},
		}

		createdCategories := make([]CategoryResponse, 0, len(categories))
		for _, req := range categories {
			category, err := categoryService.CreateCategory(ctx, tenantID, req)
			require.NoError(t, err)
			createdCategories = append(createdCategories, *category)
		}

		// Reorder categories
		newOrders := map[uuid.UUID]int{
			createdCategories[0].ID: 3, // A -> 3
			createdCategories[1].ID: 1, // B -> 1
			createdCategories[2].ID: 2, // C -> 2
		}

		err = categoryService.ReorderCategories(ctx, tenantID, newOrders)
		require.NoError(t, err)

		// Verify new order
		reorderedCategories, _, err := categoryService.ListCategories(ctx, tenantID, CategoryFilter{}, 10, 0)
		require.NoError(t, err)

		// Should be ordered: B(1), C(2), A(3)
		assert.Equal(t, "Category B", reorderedCategories[0].Name)
		assert.Equal(t, 1, reorderedCategories[0].SortOrder)
		assert.Equal(t, "Category C", reorderedCategories[1].Name)
		assert.Equal(t, 2, reorderedCategories[1].SortOrder)
		assert.Equal(t, "Category A", reorderedCategories[2].Name)
		assert.Equal(t, 3, reorderedCategories[2].SortOrder)
	})
}

func TestCategoryService_BulkOperations(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&Category{})
	require.NoError(t, err)

	// Setup services
	categoryRepo := NewGormRepository(testDB.DB)
	categoryService := NewService(categoryRepo)

	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("Bulk update status", func(t *testing.T) {
		// Create multiple categories
		categories := []CreateCategoryRequest{
			{Name: "Category 1"},
			{Name: "Category 2"},
			{Name: "Category 3"},
		}

		categoryIDs := make([]uuid.UUID, 0, len(categories))
		for _, req := range categories {
			category, err := categoryService.CreateCategory(ctx, tenantID, req)
			require.NoError(t, err)
			categoryIDs = append(categoryIDs, category.ID)
		}

		// Bulk update status to inactive
		err = categoryService.BulkUpdateStatus(ctx, tenantID, categoryIDs, StatusInactive)
		require.NoError(t, err)

		// Verify all categories are now inactive
		for _, categoryID := range categoryIDs {
			category, err := categoryService.GetCategory(ctx, tenantID, categoryID)
			require.NoError(t, err)
			assert.Equal(t, StatusInactive, category.Status)
		}

		// Bulk update back to active
		err = categoryService.BulkUpdateStatus(ctx, tenantID, categoryIDs, StatusActive)
		require.NoError(t, err)

		// Verify all categories are now active
		for _, categoryID := range categoryIDs {
			category, err := categoryService.GetCategory(ctx, tenantID, categoryID)
			require.NoError(t, err)
			assert.Equal(t, StatusActive, category.Status)
		}
	})
}

func TestCategoryService_MultiTenantIsolation(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&Category{})
	require.NoError(t, err)

	// Setup services
	categoryRepo := NewGormRepository(testDB.DB)
	categoryService := NewService(categoryRepo)

	ctx := context.Background()
	tenant1ID := uuid.New()
	tenant2ID := uuid.New()

	t.Run("Category isolation between tenants", func(t *testing.T) {
		// Create category for tenant 1
		tenant1Req := CreateCategoryRequest{
			Name:        "Tenant 1 Electronics",
			Description: "Electronics for tenant 1",
		}

		tenant1Category, err := categoryService.CreateCategory(ctx, tenant1ID, tenant1Req)
		require.NoError(t, err)

		// Create category for tenant 2
		tenant2Req := CreateCategoryRequest{
			Name:        "Tenant 2 Electronics",
			Description: "Electronics for tenant 2",
		}

		tenant2Category, err := categoryService.CreateCategory(ctx, tenant2ID, tenant2Req)
		require.NoError(t, err)

		// Verify tenant 1 can only see their categories
		tenant1Categories, total1, err := categoryService.ListCategories(ctx, tenant1ID, CategoryFilter{}, 10, 0)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total1)
		assert.Len(t, tenant1Categories, 1)
		assert.Equal(t, tenant1Category.ID, tenant1Categories[0].ID)
		assert.Contains(t, tenant1Categories[0].Name, "Tenant 1")

		// Verify tenant 2 can only see their categories
		tenant2Categories, total2, err := categoryService.ListCategories(ctx, tenant2ID, CategoryFilter{}, 10, 0)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total2)
		assert.Len(t, tenant2Categories, 1)
		assert.Equal(t, tenant2Category.ID, tenant2Categories[0].ID)
		assert.Contains(t, tenant2Categories[0].Name, "Tenant 2")

		// Verify cross-tenant access is blocked
		_, err = categoryService.GetCategory(ctx, tenant2ID, tenant1Category.ID)
		assert.Error(t, err) // Should not be able to access other tenant's category

		_, err = categoryService.GetCategory(ctx, tenant1ID, tenant2Category.ID)
		assert.Error(t, err) // Should not be able to access other tenant's category

		// Verify categories are truly isolated
		assert.NotEqual(t, tenant1Category.ID, tenant2Category.ID)
		assert.NotEqual(t, tenant1Categories[0].Name, tenant2Categories[0].Name)
	})
}