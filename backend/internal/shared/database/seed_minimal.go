package database

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedMinimalData populates the database with essential demo data only
func SeedMinimalData(db *gorm.DB) error {
	log.Println("Starting minimal database seeding...")

	// Seed basic data in dependency order
	if err := seedMinimalTenants(db); err != nil {
		return fmt.Errorf("failed to seed tenants: %w", err)
	}

	if err := seedMinimalUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := seedMinimalCategories(db); err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}

	if err := seedMinimalProducts(db); err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}

	if err := seedEssentialPermissions(db); err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	log.Println("Minimal database seeding completed successfully!")
	return nil
}

func seedMinimalTenants(db *gorm.DB) error {
	tenants := []map[string]interface{}{
		{
			"id":           "123e4567-e89b-12d3-a456-426614174000",
			"name":         "Demo Store",
			"subdomain":    "demo",
			"status":       "active",
			"plan":         "starter",
			"currency":     "BDT",
			"language":     "bn",
			"timezone":     "Asia/Dhaka",
		},
		{
			"id":           "123e4567-e89b-12d3-a456-426614174001",
			"name":         "Test Shop",
			"subdomain":    "test",
			"status":       "active",
			"plan":         "professional",
			"currency":     "BDT",
			"language":     "en",
			"timezone":     "Asia/Dhaka",
		},
	}

	for _, tenant := range tenants {
		result := db.Exec(`
			INSERT INTO tenants (id, name, subdomain, status, plan, currency, language, timezone, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`, tenant["id"], tenant["name"], tenant["subdomain"], tenant["status"], tenant["plan"], tenant["currency"], tenant["language"], tenant["timezone"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d tenants", len(tenants))
	return nil
}

func seedMinimalUsers(db *gorm.DB) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []map[string]interface{}{
		{
			"id":         "user-admin-demo",
			"tenant_id":  "123e4567-e89b-12d3-a456-426614174000",
			"email":      "admin@demo.com",
			"password":   string(hashedPassword),
			"first_name": "Demo",
			"last_name":  "Admin",
			"role":       "admin",
			"status":     "active",
		},
		{
			"id":         "user-customer-demo",
			"tenant_id":  "123e4567-e89b-12d3-a456-426614174000",
			"email":      "customer@demo.com",
			"password":   string(hashedPassword),
			"first_name": "Demo",
			"last_name":  "Customer",
			"role":       "customer",
			"status":     "active",
		},
	}

	for _, user := range users {
		result := db.Exec(`
			INSERT INTO users (id, tenant_id, email, password, first_name, last_name, role, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`, user["id"], user["tenant_id"], user["email"], user["password"], user["first_name"], user["last_name"], user["role"], user["status"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d users", len(users))
	return nil
}

func seedMinimalCategories(db *gorm.DB) error {
	categories := []map[string]interface{}{
		{
			"id":        "cat-electronics",
			"tenant_id": "123e4567-e89b-12d3-a456-426614174000",
			"name":      "Electronics",
			"slug":      "electronics",
			"is_active": true,
		},
		{
			"id":        "cat-clothing",
			"tenant_id": "123e4567-e89b-12d3-a456-426614174000",
			"name":      "Clothing",
			"slug":      "clothing",
			"is_active": true,
		},
		{
			"id":        "cat-books",
			"tenant_id": "123e4567-e89b-12d3-a456-426614174001",
			"name":      "Books",
			"slug":      "books",
			"is_active": true,
		},
	}

	for _, category := range categories {
		result := db.Exec(`
			INSERT INTO categories (id, tenant_id, name, slug, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`, category["id"], category["tenant_id"], category["name"], category["slug"], category["is_active"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d categories", len(categories))
	return nil
}

func seedMinimalProducts(db *gorm.DB) error {
	products := []map[string]interface{}{
		{
			"id":          "prod-smartphone",
			"tenant_id":   "123e4567-e89b-12d3-a456-426614174000",
			"category_id": "cat-electronics",
			"name":        "Demo Smartphone",
			"slug":        "demo-smartphone",
			"description": "A sample smartphone for demonstration",
			"price":       25000.00,
			"status":      "active",
		},
		{
			"id":          "prod-tshirt",
			"tenant_id":   "123e4567-e89b-12d3-a456-426614174000",
			"category_id": "cat-clothing",
			"name":        "Cotton T-Shirt",
			"slug":        "cotton-tshirt",
			"description": "Comfortable cotton t-shirt",
			"price":       800.00,
			"status":      "active",
		},
		{
			"id":          "prod-novel",
			"tenant_id":   "123e4567-e89b-12d3-a456-426614174001",
			"category_id": "cat-books",
			"name":        "Programming Book",
			"slug":        "programming-book",
			"description": "Learn programming fundamentals",
			"price":       1200.00,
			"status":      "active",
		},
	}

	for _, product := range products {
		result := db.Exec(`
			INSERT INTO products (id, tenant_id, category_id, name, slug, description, price, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`, product["id"], product["tenant_id"], product["category_id"], product["name"], product["slug"], product["description"], product["price"], product["status"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d products", len(products))
	return nil
}

func seedEssentialPermissions(db *gorm.DB) error {
	permissions := []map[string]interface{}{
		{"name": "products.create", "resource": "products", "action": "create"},
		{"name": "products.read", "resource": "products", "action": "read"},
		{"name": "products.update", "resource": "products", "action": "update"},
		{"name": "products.delete", "resource": "products", "action": "delete"},
		{"name": "orders.create", "resource": "orders", "action": "create"},
		{"name": "orders.read", "resource": "orders", "action": "read"},
		{"name": "orders.update", "resource": "orders", "action": "update"},
		{"name": "users.create", "resource": "users", "action": "create"},
		{"name": "users.read", "resource": "users", "action": "read"},
		{"name": "users.update", "resource": "users", "action": "update"},
	}

	for _, permission := range permissions {
		result := db.Exec(`
			INSERT INTO permissions (name, resource, action, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
			ON CONFLICT (name) DO NOTHING
		`, permission["name"], permission["resource"], permission["action"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d permissions", len(permissions))
	return nil
}