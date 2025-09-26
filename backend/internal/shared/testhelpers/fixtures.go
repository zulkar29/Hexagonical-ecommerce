package testhelpers

import (
	"testing"
	"time"

	"gorm.io/gorm"
)

// TestTenant creates a test tenant
func CreateTestTenant() map[string]interface{} {
	return map[string]interface{}{
		"id":              "11111111-1111-1111-1111-111111111111",
		"name":            "Test Store",
		"subdomain":       "test",
		"status":          "active",
		"plan":            "starter",
		"currency":        "BDT",
		"language":        "en",
		"timezone":        "Asia/Dhaka",
		"product_limit":   100,
		"storage_limit":   1024,
		"bandwidth_limit": 10240,
	}
}

// CreateTestUser creates a test user
func CreateTestUser(tenantID string) map[string]interface{} {
	return map[string]interface{}{
		"id":                "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"tenant_id":         tenantID,
		"email":             "test@example.com",
		"password":          "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		"first_name":        "Test",
		"last_name":         "User",
		"role":              "admin",
		"status":            "active",
		"email_verified":    true,
		"email_verified_at": time.Now(),
	}
}

// CreateTestProduct creates a test product
func CreateTestProduct(tenantID string) map[string]interface{} {
	return map[string]interface{}{
		"id":                 "prod1111-1111-1111-1111-111111111111",
		"tenant_id":          tenantID,
		"name":               "Test Product",
		"slug":               "test-product",
		"description":        "A test product for testing",
		"type":               "physical",
		"status":             "active",
		"price":              99.99,
		"cost_price":         50.00,
		"sku":                "TEST-001",
		"inventory_quantity": 100,
		"track_quantity":     true,
		"allow_backorder":    false,
	}
}

// CreateTestCategory creates a test category
func CreateTestCategory(tenantID string) map[string]interface{} {
	return map[string]interface{}{
		"id":          "cat11111-1111-1111-1111-111111111111",
		"tenant_id":   tenantID,
		"name":        "Test Category",
		"slug":        "test-category",
		"description": "A test category for testing",
		"sort_order":  1,
		"is_active":   true,
	}
}

// CreateTestOrder creates a test order
func CreateTestOrder(tenantID, userID string) map[string]interface{} {
	return map[string]interface{}{
		"id":                 "ord11111-1111-1111-1111-111111111111",
		"tenant_id":          tenantID,
		"user_id":            userID,
		"order_number":       "TEST-001",
		"status":             "pending",
		"subtotal":           99.99,
		"shipping_amount":    10.00,
		"discount_amount":    0.00,
		"total_amount":       109.99,
		"currency":           "BDT",
		"payment_status":     "pending",
		"fulfillment_status": "unfulfilled",
	}
}

// CleanupTables truncates test tables
func CleanupTables(t *testing.T, db *gorm.DB) {
	tables := []string{
		"order_items",
		"shipping_addresses",
		"orders",
		"product_variants",
		"products",
		"categories",
		"user_sessions",
		"users",
		"tenants",
	}

	for _, table := range tables {
		result := db.Exec("DELETE FROM " + table)
		if result.Error != nil {
			t.Logf("Warning: Could not clean table %s: %v", table, result.Error)
		}
	}
}

// SeedMinimalTestData seeds minimal data for testing
func SeedMinimalTestData(t *testing.T, db *gorm.DB) {
	// Create test tenant
	tenant := CreateTestTenant()
	result := db.Exec(`
		INSERT INTO tenants (id, name, subdomain, status, plan, currency, language, timezone, product_limit, storage_limit, bandwidth_limit, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
	`, tenant["id"], tenant["name"], tenant["subdomain"], tenant["status"], tenant["plan"], tenant["currency"], tenant["language"], tenant["timezone"], tenant["product_limit"], tenant["storage_limit"], tenant["bandwidth_limit"])

	if result.Error != nil {
		t.Fatalf("Failed to seed test tenant: %v", result.Error)
	}

	// Create test user
	user := CreateTestUser(tenant["id"].(string))
	result = db.Exec(`
		INSERT INTO users (id, tenant_id, email, password, first_name, last_name, role, status, email_verified, email_verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`, user["id"], user["tenant_id"], user["email"], user["password"], user["first_name"], user["last_name"], user["role"], user["status"], user["email_verified"], user["email_verified_at"])

	if result.Error != nil {
		t.Fatalf("Failed to seed test user: %v", result.Error)
	}
}

// Use utils.GenerateUUID() instead - removed duplicate

// AssertNoError is a helper for asserting no error
func AssertNoError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// AssertEqual is a helper for asserting equality
func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	if expected != actual {
		t.Fatalf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// AssertNotNil is a helper for asserting not nil
func AssertNotNil(t *testing.T, value interface{}, msg string) {
	if value == nil {
		t.Fatalf("%s: expected not nil", msg)
	}
}
