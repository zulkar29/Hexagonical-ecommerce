package database

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedData populates the database with initial development/demo data
func SeedData(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Seed in order of dependencies
	if err := seedTenants(db); err != nil {
		return fmt.Errorf("failed to seed tenants: %w", err)
	}

	if err := seedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := seedCategories(db); err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}

	if err := seedProducts(db); err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}

	if err := seedOrders(db); err != nil {
		return fmt.Errorf("failed to seed orders: %w", err)
	}

	if err := seedProductVariants(db); err != nil {
		return fmt.Errorf("failed to seed product variants: %w", err)
	}

	if err := seedPermissions(db); err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	if err := seedAnalyticsData(db); err != nil {
		return fmt.Errorf("failed to seed analytics data: %w", err)
	}

	if err := seedPaymentData(db); err != nil {
		return fmt.Errorf("failed to seed payment data: %w", err)
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// seedTenants creates demo tenants
func seedTenants(db *gorm.DB) error {
	log.Println("Seeding tenants...")

	tenants := []map[string]interface{}{
		{
			"id":              "11111111-1111-1111-1111-111111111111",
			"name":            "TechHub Electronics",
			"subdomain":       "techhub",
			"custom_domain":   "techhub.example.com",
			"status":          "active",
			"plan":            "professional",
			"description":     "Leading electronics and gadgets retailer",
			"phone":           "+8801712345678",
			"email":           "info@techhub.com",
			"address":         "123 Tech Street, Gulshan, Dhaka, Bangladesh",
			"logo":            "https://example.com/logos/techhub.png",
			"currency":        "BDT",
			"language":        "en",
			"timezone":        "Asia/Dhaka",
			"product_limit":   1000,
			"storage_limit":   5120,
			"bandwidth_limit": 51200,
		},
		{
			"id":              "22222222-2222-2222-2222-222222222222",
			"name":            "Fashion Forward",
			"subdomain":       "fashionforward",
			"custom_domain":   "fashionforward.example.com",
			"status":          "active",
			"plan":            "premium",
			"description":     "Trendy fashion and lifestyle brand",
			"phone":           "+8801787654321",
			"email":           "contact@fashionforward.com",
			"address":         "456 Fashion Ave, Dhanmondi, Dhaka, Bangladesh",
			"logo":            "https://example.com/logos/fashionforward.png",
			"currency":        "BDT",
			"language":        "en",
			"timezone":        "Asia/Dhaka",
			"product_limit":   2000,
			"storage_limit":   10240,
			"bandwidth_limit": 102400,
		},
		{
			"id":              "33333333-3333-3333-3333-333333333333",
			"name":            "Home & Garden",
			"subdomain":       "homegarden",
			"custom_domain":   "",
			"status":          "active",
			"plan":            "starter",
			"description":     "Everything for your home and garden needs",
			"phone":           "+8801698765432",
			"email":           "hello@homegarden.com",
			"address":         "789 Garden Lane, Uttara, Dhaka, Bangladesh",
			"logo":            "https://example.com/logos/homegarden.png",
			"currency":        "BDT",
			"language":        "bn",
			"timezone":        "Asia/Dhaka",
			"product_limit":   500,
			"storage_limit":   2048,
			"bandwidth_limit": 20480,
		},
	}

	for _, tenant := range tenants {
		result := db.Exec(`
			INSERT INTO tenants (id, name, subdomain, custom_domain, status, plan, description, phone, email, address, logo, currency, language, timezone, product_limit, storage_limit, bandwidth_limit, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`, tenant["id"], tenant["name"], tenant["subdomain"], tenant["custom_domain"], tenant["status"], tenant["plan"], tenant["description"], tenant["phone"], tenant["email"], tenant["address"], tenant["logo"], tenant["currency"], tenant["language"], tenant["timezone"], tenant["product_limit"], tenant["storage_limit"], tenant["bandwidth_limit"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d tenants", len(tenants))
	return nil
}

// seedUsers creates demo users for each tenant
func seedUsers(db *gorm.DB) error {
	log.Println("Seeding users...")

	// Hash password for demo users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	users := []map[string]interface{}{
		// TechHub users
		{
			"id":                "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			"tenant_id":         "11111111-1111-1111-1111-111111111111",
			"email":             "admin@techhub.com",
			"password":          string(hashedPassword),
			"first_name":        "Ahmed",
			"last_name":         "Rahman",
			"phone":             "+8801712345678",
			"role":              "admin",
			"status":            "active",
			"email_verified":    true,
			"email_verified_at": time.Now(),
		},
		{
			"id":                "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			"tenant_id":         "11111111-1111-1111-1111-111111111111",
			"email":             "merchant@techhub.com",
			"password":          string(hashedPassword),
			"first_name":        "Fatima",
			"last_name":         "Khan",
			"phone":             "+8801787654321",
			"role":              "merchant",
			"status":            "active",
			"email_verified":    true,
			"email_verified_at": time.Now(),
		},
		{
			"id":                "cccccccc-cccc-cccc-cccc-cccccccccccc",
			"tenant_id":         "11111111-1111-1111-1111-111111111111",
			"email":             "customer1@example.com",
			"password":          string(hashedPassword),
			"first_name":        "Mohammad",
			"last_name":         "Islam",
			"phone":             "+8801698765432",
			"role":              "customer",
			"status":            "active",
			"email_verified":    true,
			"email_verified_at": time.Now(),
		},
		{
			"id":                "dddddddd-dddd-dddd-dddd-dddddddddddd",
			"tenant_id":         "11111111-1111-1111-1111-111111111111",
			"email":             "customer2@example.com",
			"password":          string(hashedPassword),
			"first_name":        "Rashida",
			"last_name":         "Begum",
			"phone":             "+8801512345678",
			"role":              "customer",
			"status":            "active",
			"email_verified":    true,
			"email_verified_at": time.Now(),
		},
		// Fashion Forward users
		{
			"id":                "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee",
			"tenant_id":         "22222222-2222-2222-2222-222222222222",
			"email":             "admin@fashionforward.com",
			"password":          string(hashedPassword),
			"first_name":        "Samira",
			"last_name":         "Ahmed",
			"phone":             "+8801612345678",
			"role":              "admin",
			"status":            "active",
			"email_verified":    true,
			"email_verified_at": time.Now(),
		},
		{
			"id":                "ffffffff-ffff-ffff-ffff-ffffffffffff",
			"tenant_id":         "22222222-2222-2222-2222-222222222222",
			"email":             "customer3@example.com",
			"password":          string(hashedPassword),
			"first_name":        "Karim",
			"last_name":         "Hassan",
			"phone":             "+8801812345678",
			"role":              "customer",
			"status":            "active",
			"email_verified":    true,
			"email_verified_at": time.Now(),
		},
		// Home & Garden users
		{
			"id":                "gggggggg-gggg-gggg-gggg-gggggggggggg",
			"tenant_id":         "33333333-3333-3333-3333-333333333333",
			"email":             "admin@homegarden.com",
			"password":          string(hashedPassword),
			"first_name":        "Nasir",
			"last_name":         "Uddin",
			"phone":             "+8801912345678",
			"role":              "admin",
			"status":            "active",
			"email_verified":    true,
			"email_verified_at": time.Now(),
		},
	}

	for _, user := range users {
		result := db.Exec(`
			INSERT INTO users (id, tenant_id, email, password, first_name, last_name, phone, role, status, email_verified, email_verified_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
			ON CONFLICT (email) DO NOTHING
		`, user["id"], user["tenant_id"], user["email"], user["password"], user["first_name"], user["last_name"], user["phone"], user["role"], user["status"], user["email_verified"], user["email_verified_at"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d users", len(users))
	return nil
}

// seedCategories creates demo product categories
func seedCategories(db *gorm.DB) error {
	log.Println("Seeding categories...")

	categories := []map[string]interface{}{
		// TechHub categories
		{
			"id":               "cat11111-1111-1111-1111-111111111111",
			"tenant_id":        "11111111-1111-1111-1111-111111111111",
			"name":             "Electronics",
			"slug":             "electronics",
			"description":      "Latest electronic gadgets and devices",
			"image":            "https://example.com/categories/electronics.jpg",
			"sort_order":       1,
			"is_active":        true,
			"meta_title":       "Electronics | TechHub",
			"meta_description": "Discover the latest electronics and gadgets at TechHub",
		},
		{
			"id":               "cat11112-1111-1111-1111-111111111111",
			"tenant_id":        "11111111-1111-1111-1111-111111111111",
			"name":             "Smartphones",
			"slug":             "smartphones",
			"description":      "Latest smartphones and mobile devices",
			"image":            "https://example.com/categories/smartphones.jpg",
			"parent_id":        "cat11111-1111-1111-1111-111111111111",
			"sort_order":       1,
			"is_active":        true,
			"meta_title":       "Smartphones | TechHub",
			"meta_description": "Browse our collection of latest smartphones",
		},
		{
			"id":               "cat11113-1111-1111-1111-111111111111",
			"tenant_id":        "11111111-1111-1111-1111-111111111111",
			"name":             "Laptops",
			"slug":             "laptops",
			"description":      "High-performance laptops and notebooks",
			"image":            "https://example.com/categories/laptops.jpg",
			"parent_id":        "cat11111-1111-1111-1111-111111111111",
			"sort_order":       2,
			"is_active":        true,
			"meta_title":       "Laptops | TechHub",
			"meta_description": "Find the perfect laptop for your needs",
		},
		{
			"id":               "cat11114-1111-1111-1111-111111111111",
			"tenant_id":        "11111111-1111-1111-1111-111111111111",
			"name":             "Accessories",
			"slug":             "accessories",
			"description":      "Tech accessories and peripherals",
			"image":            "https://example.com/categories/accessories.jpg",
			"sort_order":       2,
			"is_active":        true,
			"meta_title":       "Tech Accessories | TechHub",
			"meta_description": "Complete your tech setup with our accessories",
		},
		// Fashion Forward categories
		{
			"id":               "cat22221-2222-2222-2222-222222222222",
			"tenant_id":        "22222222-2222-2222-2222-222222222222",
			"name":             "Clothing",
			"slug":             "clothing",
			"description":      "Trendy clothing for all occasions",
			"image":            "https://example.com/categories/clothing.jpg",
			"sort_order":       1,
			"is_active":        true,
			"meta_title":       "Clothing | Fashion Forward",
			"meta_description": "Discover trendy clothing at Fashion Forward",
		},
		{
			"id":               "cat22222-2222-2222-2222-222222222222",
			"tenant_id":        "22222222-2222-2222-2222-222222222222",
			"name":             "Men's Fashion",
			"slug":             "mens-fashion",
			"description":      "Stylish clothing for men",
			"image":            "https://example.com/categories/mens-fashion.jpg",
			"parent_id":        "cat22221-2222-2222-2222-222222222222",
			"sort_order":       1,
			"is_active":        true,
			"meta_title":       "Men's Fashion | Fashion Forward",
			"meta_description": "Explore our men's fashion collection",
		},
		{
			"id":               "cat22223-2222-2222-2222-222222222222",
			"tenant_id":        "22222222-2222-2222-2222-222222222222",
			"name":             "Women's Fashion",
			"slug":             "womens-fashion",
			"description":      "Elegant clothing for women",
			"image":            "https://example.com/categories/womens-fashion.jpg",
			"parent_id":        "cat22221-2222-2222-2222-222222222222",
			"sort_order":       2,
			"is_active":        true,
			"meta_title":       "Women's Fashion | Fashion Forward",
			"meta_description": "Browse our women's fashion collection",
		},
		// Home & Garden categories
		{
			"id":               "cat33331-3333-3333-3333-333333333333",
			"tenant_id":        "33333333-3333-3333-3333-333333333333",
			"name":             "Home Decor",
			"slug":             "home-decor",
			"description":      "Beautiful home decoration items",
			"image":            "https://example.com/categories/home-decor.jpg",
			"sort_order":       1,
			"is_active":        true,
			"meta_title":       "Home Decor | Home & Garden",
			"meta_description": "Transform your home with our decor collection",
		},
		{
			"id":               "cat33332-3333-3333-3333-333333333333",
			"tenant_id":        "33333333-3333-3333-3333-333333333333",
			"name":             "Garden Tools",
			"slug":             "garden-tools",
			"description":      "Essential tools for your garden",
			"image":            "https://example.com/categories/garden-tools.jpg",
			"sort_order":       2,
			"is_active":        true,
			"meta_title":       "Garden Tools | Home & Garden",
			"meta_description": "Professional garden tools and equipment",
		},
	}

	for _, category := range categories {
		result := db.Exec(`
			INSERT INTO categories (id, tenant_id, name, slug, description, image, parent_id, sort_order, is_active, meta_title, meta_description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
			ON CONFLICT (tenant_id, slug) DO NOTHING
		`, category["id"], category["tenant_id"], category["name"], category["slug"], category["description"], category["image"], category["parent_id"], category["sort_order"], category["is_active"], category["meta_title"], category["meta_description"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d categories", len(categories))
	return nil
}

// seedProducts creates demo products
func seedProducts(db *gorm.DB) error {
	log.Println("Seeding products...")

	products := []map[string]interface{}{
		// TechHub products
		{
			"id":                 "prod1111-1111-1111-1111-111111111111",
			"tenant_id":          "11111111-1111-1111-1111-111111111111",
			"name":               "iPhone 15 Pro",
			"slug":               "iphone-15-pro",
			"description":        "Latest iPhone with advanced camera system and A17 Pro chip",
			"type":               "physical",
			"status":             "active",
			"price":              129999.00,
			"compare_price":      139999.00,
			"cost_price":         119999.00,
			"sku":                "APL-IP15P-256",
			"barcode":            "1234567890123",
			"inventory_quantity": 50,
			"track_quantity":     true,
			"allow_backorder":    false,
			"weight":             187.0,
			"length":             14.67,
			"width":              7.09,
			"height":             0.83,
			"featured_image":     "https://example.com/products/iphone-15-pro.jpg",
			"images":             `["https://example.com/products/iphone-15-pro-1.jpg", "https://example.com/products/iphone-15-pro-2.jpg"]`,
			"category_id":        "cat11112-1111-1111-1111-111111111111",
			"tags":               `["smartphone", "apple", "5g", "premium"]`,
			"meta_title":         "iPhone 15 Pro - Latest Apple Smartphone | TechHub",
			"meta_description":   "Get the latest iPhone 15 Pro with advanced features and A17 Pro chip",
		},
		{
			"id":                 "prod1112-1111-1111-1111-111111111111",
			"tenant_id":          "11111111-1111-1111-1111-111111111111",
			"name":               "MacBook Air M3",
			"slug":               "macbook-air-m3",
			"description":        "Powerful and lightweight laptop with M3 chip for ultimate performance",
			"type":               "physical",
			"status":             "active",
			"price":              154999.00,
			"compare_price":      164999.00,
			"cost_price":         144999.00,
			"sku":                "APL-MBA-M3-13",
			"barcode":            "1234567890124",
			"inventory_quantity": 25,
			"track_quantity":     true,
			"allow_backorder":    false,
			"weight":             1240.0,
			"length":             30.41,
			"width":              21.5,
			"height":             1.13,
			"featured_image":     "https://example.com/products/macbook-air-m3.jpg",
			"images":             `["https://example.com/products/macbook-air-m3-1.jpg", "https://example.com/products/macbook-air-m3-2.jpg"]`,
			"category_id":        "cat11113-1111-1111-1111-111111111111",
			"tags":               `["laptop", "apple", "m3", "ultrabook"]`,
			"meta_title":         "MacBook Air M3 - Powerful Laptop | TechHub",
			"meta_description":   "Experience ultimate performance with MacBook Air M3",
		},
		{
			"id":                 "prod1113-1111-1111-1111-111111111111",
			"tenant_id":          "11111111-1111-1111-1111-111111111111",
			"name":               "AirPods Pro (3rd Gen)",
			"slug":               "airpods-pro-3rd-gen",
			"description":        "Wireless earbuds with active noise cancellation and spatial audio",
			"type":               "physical",
			"status":             "active",
			"price":              24999.00,
			"compare_price":      27999.00,
			"cost_price":         22999.00,
			"sku":                "APL-APP3-WHT",
			"barcode":            "1234567890125",
			"inventory_quantity": 75,
			"track_quantity":     true,
			"allow_backorder":    true,
			"weight":             56.0,
			"length":             4.5,
			"width":              6.1,
			"height":             2.17,
			"featured_image":     "https://example.com/products/airpods-pro-3.jpg",
			"images":             `["https://example.com/products/airpods-pro-3-1.jpg", "https://example.com/products/airpods-pro-3-2.jpg"]`,
			"category_id":        "cat11114-1111-1111-1111-111111111111",
			"tags":               `["earbuds", "wireless", "apple", "noise-cancellation"]`,
			"meta_title":         "AirPods Pro 3rd Gen - Wireless Earbuds | TechHub",
			"meta_description":   "Premium wireless earbuds with active noise cancellation",
		},
		// Fashion Forward products
		{
			"id":                 "prod2221-2222-2222-2222-222222222222",
			"tenant_id":          "22222222-2222-2222-2222-222222222222",
			"name":               "Premium Cotton T-Shirt",
			"slug":               "premium-cotton-t-shirt",
			"description":        "Comfortable and stylish cotton t-shirt for everyday wear",
			"type":               "physical",
			"status":             "active",
			"price":              1299.00,
			"compare_price":      1599.00,
			"cost_price":         899.00,
			"sku":                "FF-TS-COT-001",
			"barcode":            "2234567890123",
			"inventory_quantity": 200,
			"track_quantity":     true,
			"allow_backorder":    true,
			"weight":             150.0,
			"length":             70.0,
			"width":              50.0,
			"height":             1.0,
			"featured_image":     "https://example.com/products/cotton-tshirt.jpg",
			"images":             `["https://example.com/products/cotton-tshirt-1.jpg", "https://example.com/products/cotton-tshirt-2.jpg"]`,
			"category_id":        "cat22222-2222-2222-2222-222222222222",
			"tags":               `["t-shirt", "cotton", "casual", "mens"]`,
			"meta_title":         "Premium Cotton T-Shirt | Fashion Forward",
			"meta_description":   "Comfortable cotton t-shirt perfect for casual wear",
		},
		{
			"id":                 "prod2222-2222-2222-2222-222222222222",
			"tenant_id":          "22222222-2222-2222-2222-222222222222",
			"name":               "Elegant Summer Dress",
			"slug":               "elegant-summer-dress",
			"description":        "Beautiful and comfortable summer dress for special occasions",
			"type":               "physical",
			"status":             "active",
			"price":              3499.00,
			"compare_price":      3999.00,
			"cost_price":         2799.00,
			"sku":                "FF-SD-ELE-001",
			"barcode":            "2234567890124",
			"inventory_quantity": 80,
			"track_quantity":     true,
			"allow_backorder":    false,
			"weight":             300.0,
			"length":             120.0,
			"width":              60.0,
			"height":             2.0,
			"featured_image":     "https://example.com/products/summer-dress.jpg",
			"images":             `["https://example.com/products/summer-dress-1.jpg", "https://example.com/products/summer-dress-2.jpg"]`,
			"category_id":        "cat22223-2222-2222-2222-222222222222",
			"tags":               `["dress", "summer", "elegant", "womens"]`,
			"meta_title":         "Elegant Summer Dress | Fashion Forward",
			"meta_description":   "Beautiful summer dress perfect for any occasion",
		},
		// Home & Garden products
		{
			"id":                 "prod3331-3333-3333-3333-333333333333",
			"tenant_id":          "33333333-3333-3333-3333-333333333333",
			"name":               "Ceramic Table Lamp",
			"slug":               "ceramic-table-lamp",
			"description":        "Elegant ceramic table lamp to brighten your living space",
			"type":               "physical",
			"status":             "active",
			"price":              2499.00,
			"compare_price":      2999.00,
			"cost_price":         1899.00,
			"sku":                "HG-LAMP-CER-001",
			"barcode":            "3234567890123",
			"inventory_quantity": 45,
			"track_quantity":     true,
			"allow_backorder":    false,
			"weight":             800.0,
			"length":             25.0,
			"width":              25.0,
			"height":             35.0,
			"featured_image":     "https://example.com/products/ceramic-lamp.jpg",
			"images":             `["https://example.com/products/ceramic-lamp-1.jpg", "https://example.com/products/ceramic-lamp-2.jpg"]`,
			"category_id":        "cat33331-3333-3333-3333-333333333333",
			"tags":               `["lamp", "ceramic", "decor", "lighting"]`,
			"meta_title":         "Ceramic Table Lamp | Home & Garden",
			"meta_description":   "Elegant ceramic table lamp for your home decor",
		},
		{
			"id":                 "prod3332-3333-3333-3333-333333333333",
			"tenant_id":          "33333333-3333-3333-3333-333333333333",
			"name":               "Garden Tool Set",
			"slug":               "garden-tool-set",
			"description":        "Complete set of essential garden tools for all your gardening needs",
			"type":               "physical",
			"status":             "active",
			"price":              4999.00,
			"compare_price":      5999.00,
			"cost_price":         3999.00,
			"sku":                "HG-TOOLS-SET-001",
			"barcode":            "3234567890124",
			"inventory_quantity": 30,
			"track_quantity":     true,
			"allow_backorder":    true,
			"weight":             2500.0,
			"length":             40.0,
			"width":              30.0,
			"height":             10.0,
			"featured_image":     "https://example.com/products/garden-tool-set.jpg",
			"images":             `["https://example.com/products/garden-tool-set-1.jpg", "https://example.com/products/garden-tool-set-2.jpg"]`,
			"category_id":        "cat33332-3333-3333-3333-333333333333",
			"tags":               `["tools", "garden", "set", "outdoor"]`,
			"meta_title":         "Garden Tool Set | Home & Garden",
			"meta_description":   "Complete garden tool set for all your gardening needs",
		},
	}

	for _, product := range products {
		result := db.Exec(`
			INSERT INTO products (id, tenant_id, name, slug, description, type, status, price, compare_price, cost_price, sku, barcode, inventory_quantity, track_quantity, allow_backorder, weight, length, width, height, featured_image, images, category_id, tags, meta_title, meta_description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, NOW(), NOW())
			ON CONFLICT (tenant_id, slug) DO NOTHING
		`, product["id"], product["tenant_id"], product["name"], product["slug"], product["description"], product["type"], product["status"], product["price"], product["compare_price"], product["cost_price"], product["sku"], product["barcode"], product["inventory_quantity"], product["track_quantity"], product["allow_backorder"], product["weight"], product["length"], product["width"], product["height"], product["featured_image"], product["images"], product["category_id"], product["tags"], product["meta_title"], product["meta_description"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d products", len(products))
	return nil
}

// seedOrders creates demo orders
func seedOrders(db *gorm.DB) error {
	log.Println("Seeding orders...")

	// First, create orders
	orders := []map[string]interface{}{
		{
			"id":                 "ord11111-1111-1111-1111-111111111111",
			"tenant_id":          "11111111-1111-1111-1111-111111111111",
			"user_id":            "cccccccc-cccc-cccc-cccc-cccccccccccc", // customer1
			"order_number":       "TH-2024-0001",
			"status":             "confirmed",
			"subtotal":           154998.00,
			"shipping_amount":    500.00,
			"discount_amount":    0.00,
			"total_amount":       155498.00,
			"currency":           "BDT",
			"payment_status":     "paid",
			"fulfillment_status": "fulfilled",
			"notes":              "Customer requested fast delivery",
		},
		{
			"id":                 "ord11112-1111-1111-1111-111111111111",
			"tenant_id":          "11111111-1111-1111-1111-111111111111",
			"user_id":            "dddddddd-dddd-dddd-dddd-dddddddddddd", // customer2
			"order_number":       "TH-2024-0002",
			"status":             "processing",
			"subtotal":           24999.00,
			"shipping_amount":    300.00,
			"discount_amount":    1000.00,
			"total_amount":       24299.00,
			"currency":           "BDT",
			"payment_status":     "paid",
			"fulfillment_status": "pending",
			"notes":              "",
		},
		{
			"id":                 "ord22221-2222-2222-2222-222222222222",
			"tenant_id":          "22222222-2222-2222-2222-222222222222",
			"user_id":            "ffffffff-ffff-ffff-ffff-ffffffffffff", // customer3
			"order_number":       "FF-2024-0001",
			"status":             "delivered",
			"subtotal":           4798.00,
			"shipping_amount":    200.00,
			"discount_amount":    300.00,
			"total_amount":       4698.00,
			"currency":           "BDT",
			"payment_status":     "paid",
			"fulfillment_status": "fulfilled",
			"notes":              "Gift wrap requested",
		},
	}

	for _, order := range orders {
		result := db.Exec(`
			INSERT INTO orders (id, tenant_id, user_id, order_number, status, subtotal, shipping_amount, discount_amount, total_amount, currency, payment_status, fulfillment_status, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days')
			ON CONFLICT (order_number) DO NOTHING
		`, order["id"], order["tenant_id"], order["user_id"], order["order_number"], order["status"], order["subtotal"], order["shipping_amount"], order["discount_amount"], order["total_amount"], order["currency"], order["payment_status"], order["fulfillment_status"], order["notes"])

		if result.Error != nil {
			return result.Error
		}
	}

	// Then, create order items
	orderItems := []map[string]interface{}{
		// Order 1 items
		{
			"tenant_id":     "11111111-1111-1111-1111-111111111111",
			"order_id":      "ord11111-1111-1111-1111-111111111111",
			"product_id":    "prod1111-1111-1111-1111-111111111111", // iPhone 15 Pro
			"quantity":      1,
			"unit_price":    129999.00,
			"total_price":   129999.00,
			"product_name":  "iPhone 15 Pro",
			"product_sku":   "APL-IP15P-256",
			"product_image": "https://example.com/products/iphone-15-pro.jpg",
		},
		{
			"tenant_id":     "11111111-1111-1111-1111-111111111111",
			"order_id":      "ord11111-1111-1111-1111-111111111111",
			"product_id":    "prod1113-1111-1111-1111-111111111111", // AirPods Pro
			"quantity":      1,
			"unit_price":    24999.00,
			"total_price":   24999.00,
			"product_name":  "AirPods Pro (3rd Gen)",
			"product_sku":   "APL-APP3-WHT",
			"product_image": "https://example.com/products/airpods-pro-3.jpg",
		},
		// Order 2 items
		{
			"tenant_id":     "11111111-1111-1111-1111-111111111111",
			"order_id":      "ord11112-1111-1111-1111-111111111111",
			"product_id":    "prod1113-1111-1111-1111-111111111111", // AirPods Pro
			"quantity":      1,
			"unit_price":    24999.00,
			"total_price":   24999.00,
			"product_name":  "AirPods Pro (3rd Gen)",
			"product_sku":   "APL-APP3-WHT",
			"product_image": "https://example.com/products/airpods-pro-3.jpg",
		},
		// Order 3 items
		{
			"tenant_id":     "22222222-2222-2222-2222-222222222222",
			"order_id":      "ord22221-2222-2222-2222-222222222222",
			"product_id":    "prod2221-2222-2222-2222-222222222222", // T-Shirt
			"quantity":      2,
			"unit_price":    1299.00,
			"total_price":   2598.00,
			"product_name":  "Premium Cotton T-Shirt",
			"product_sku":   "FF-TS-COT-001",
			"product_image": "https://example.com/products/cotton-tshirt.jpg",
		},
		{
			"tenant_id":     "22222222-2222-2222-2222-222222222222",
			"order_id":      "ord22221-2222-2222-2222-222222222222",
			"product_id":    "prod2222-2222-2222-2222-222222222222", // Summer Dress
			"quantity":      1,
			"unit_price":    3499.00,
			"total_price":   3499.00,
			"product_name":  "Elegant Summer Dress",
			"product_sku":   "FF-SD-ELE-001",
			"product_image": "https://example.com/products/summer-dress.jpg",
		},
	}

	for _, item := range orderItems {
		result := db.Exec(`
			INSERT INTO order_items (tenant_id, order_id, product_id, quantity, unit_price, total_price, product_name, product_sku, product_image, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days')
		`, item["tenant_id"], item["order_id"], item["product_id"], item["quantity"], item["unit_price"], item["total_price"], item["product_name"], item["product_sku"], item["product_image"])

		if result.Error != nil {
			return result.Error
		}
	}

	// Create shipping addresses
	shippingAddresses := []map[string]interface{}{
		{
			"tenant_id":      "11111111-1111-1111-1111-111111111111",
			"order_id":       "ord11111-1111-1111-1111-111111111111",
			"first_name":     "Mohammad",
			"last_name":      "Islam",
			"company":        "",
			"address_line_1": "House 123, Road 4",
			"address_line_2": "Block B, Bashundhara R/A",
			"city":           "Dhaka",
			"state":          "Dhaka",
			"postal_code":    "1229",
			"country":        "BD",
			"phone":          "+8801698765432",
		},
		{
			"tenant_id":      "11111111-1111-1111-1111-111111111111",
			"order_id":       "ord11112-1111-1111-1111-111111111111",
			"first_name":     "Rashida",
			"last_name":      "Begum",
			"company":        "",
			"address_line_1": "Flat 5B, Building 12",
			"address_line_2": "Dhanmondi 15",
			"city":           "Dhaka",
			"state":          "Dhaka",
			"postal_code":    "1209",
			"country":        "BD",
			"phone":          "+8801512345678",
		},
		{
			"tenant_id":      "22222222-2222-2222-2222-222222222222",
			"order_id":       "ord22221-2222-2222-2222-222222222222",
			"first_name":     "Karim",
			"last_name":      "Hassan",
			"company":        "",
			"address_line_1": "456 Fashion Street",
			"address_line_2": "Gulshan 2",
			"city":           "Dhaka",
			"state":          "Dhaka",
			"postal_code":    "1212",
			"country":        "BD",
			"phone":          "+8801812345678",
		},
	}

	for _, address := range shippingAddresses {
		result := db.Exec(`
			INSERT INTO shipping_addresses (tenant_id, order_id, first_name, last_name, company, address_line_1, address_line_2, city, state, postal_code, country, phone, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days')
		`, address["tenant_id"], address["order_id"], address["first_name"], address["last_name"], address["company"], address["address_line_1"], address["address_line_2"], address["city"], address["state"], address["postal_code"], address["country"], address["phone"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d orders with %d items and %d shipping addresses", len(orders), len(orderItems), len(shippingAddresses))
	return nil
}

// seedProductVariants creates demo product variants
func seedProductVariants(db *gorm.DB) error {
	log.Println("Seeding product variants...")

	variants := []map[string]interface{}{
		// iPhone 15 Pro variants (colors)
		{
			"product_id":         "prod1111-1111-1111-1111-111111111111",
			"name":               "iPhone 15 Pro - Natural Titanium 256GB",
			"sku":                "APL-IP15P-256-NAT",
			"price":              129999.00,
			"inventory_quantity": 25,
			"allow_backorder":    false,
			"options":            `{"color": "Natural Titanium", "storage": "256GB"}`,
			"images":             `["https://example.com/products/iphone-15-pro-natural.jpg"]`,
		},
		{
			"product_id":         "prod1111-1111-1111-1111-111111111111",
			"name":               "iPhone 15 Pro - Blue Titanium 256GB",
			"sku":                "APL-IP15P-256-BLU",
			"price":              129999.00,
			"inventory_quantity": 15,
			"allow_backorder":    false,
			"options":            `{"color": "Blue Titanium", "storage": "256GB"}`,
			"images":             `["https://example.com/products/iphone-15-pro-blue.jpg"]`,
		},
		{
			"product_id":         "prod1111-1111-1111-1111-111111111111",
			"name":               "iPhone 15 Pro - Natural Titanium 512GB",
			"sku":                "APL-IP15P-512-NAT",
			"price":              149999.00,
			"inventory_quantity": 10,
			"allow_backorder":    false,
			"options":            `{"color": "Natural Titanium", "storage": "512GB"}`,
			"images":             `["https://example.com/products/iphone-15-pro-natural.jpg"]`,
		},
		// T-Shirt variants (sizes and colors)
		{
			"product_id":         "prod2221-2222-2222-2222-222222222222",
			"name":               "Premium Cotton T-Shirt - White - Medium",
			"sku":                "FF-TS-COT-WHT-M",
			"price":              1299.00,
			"inventory_quantity": 50,
			"allow_backorder":    true,
			"options":            `{"color": "White", "size": "M"}`,
			"images":             `["https://example.com/products/cotton-tshirt-white.jpg"]`,
		},
		{
			"product_id":         "prod2221-2222-2222-2222-222222222222",
			"name":               "Premium Cotton T-Shirt - Black - Large",
			"sku":                "FF-TS-COT-BLK-L",
			"price":              1299.00,
			"inventory_quantity": 30,
			"allow_backorder":    true,
			"options":            `{"color": "Black", "size": "L"}`,
			"images":             `["https://example.com/products/cotton-tshirt-black.jpg"]`,
		},
		{
			"product_id":         "prod2221-2222-2222-2222-222222222222",
			"name":               "Premium Cotton T-Shirt - Navy - Small",
			"sku":                "FF-TS-COT-NAV-S",
			"price":              1299.00,
			"inventory_quantity": 25,
			"allow_backorder":    true,
			"options":            `{"color": "Navy", "size": "S"}`,
			"images":             `["https://example.com/products/cotton-tshirt-navy.jpg"]`,
		},
		// Summer Dress variants (sizes)
		{
			"product_id":         "prod2222-2222-2222-2222-222222222222",
			"name":               "Elegant Summer Dress - Medium",
			"sku":                "FF-SD-ELE-M",
			"price":              3499.00,
			"inventory_quantity": 20,
			"allow_backorder":    false,
			"options":            `{"size": "M"}`,
			"images":             `["https://example.com/products/summer-dress-m.jpg"]`,
		},
		{
			"product_id":         "prod2222-2222-2222-2222-222222222222",
			"name":               "Elegant Summer Dress - Large",
			"sku":                "FF-SD-ELE-L",
			"price":              3499.00,
			"inventory_quantity": 15,
			"allow_backorder":    false,
			"options":            `{"size": "L"}`,
			"images":             `["https://example.com/products/summer-dress-l.jpg"]`,
		},
	}

	for _, variant := range variants {
		result := db.Exec(`
			INSERT INTO product_variants (product_id, name, sku, price, inventory_quantity, allow_backorder, options, images, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			ON CONFLICT (sku) DO NOTHING
		`, variant["product_id"], variant["name"], variant["sku"], variant["price"], variant["inventory_quantity"], variant["allow_backorder"], variant["options"], variant["images"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d product variants", len(variants))
	return nil
}

// seedPermissions creates comprehensive permission system
func seedPermissions(db *gorm.DB) error {
	log.Println("Seeding permissions...")

	permissions := []map[string]interface{}{
		// Product permissions
		{"name": "products.create", "description": "Create products", "resource": "products", "action": "create"},
		{"name": "products.read", "description": "View products", "resource": "products", "action": "read"},
		{"name": "products.update", "description": "Update products", "resource": "products", "action": "update"},
		{"name": "products.delete", "description": "Delete products", "resource": "products", "action": "delete"},

		// Order permissions
		{"name": "orders.create", "description": "Create orders", "resource": "orders", "action": "create"},
		{"name": "orders.read", "description": "View orders", "resource": "orders", "action": "read"},
		{"name": "orders.update", "description": "Update orders", "resource": "orders", "action": "update"},
		{"name": "orders.delete", "description": "Delete orders", "resource": "orders", "action": "delete"},
		{"name": "orders.fulfill", "description": "Fulfill orders", "resource": "orders", "action": "fulfill"},

		// User permissions
		{"name": "users.create", "description": "Create users", "resource": "users", "action": "create"},
		{"name": "users.read", "description": "View users", "resource": "users", "action": "read"},
		{"name": "users.update", "description": "Update users", "resource": "users", "action": "update"},
		{"name": "users.delete", "description": "Delete users", "resource": "users", "action": "delete"},

		// Category permissions
		{"name": "categories.create", "description": "Create categories", "resource": "categories", "action": "create"},
		{"name": "categories.read", "description": "View categories", "resource": "categories", "action": "read"},
		{"name": "categories.update", "description": "Update categories", "resource": "categories", "action": "update"},
		{"name": "categories.delete", "description": "Delete categories", "resource": "categories", "action": "delete"},

		// Analytics permissions
		{"name": "analytics.read", "description": "View analytics", "resource": "analytics", "action": "read"},
		{"name": "analytics.export", "description": "Export analytics", "resource": "analytics", "action": "export"},

		// Settings permissions
		{"name": "settings.read", "description": "View settings", "resource": "settings", "action": "read"},
		{"name": "settings.update", "description": "Update settings", "resource": "settings", "action": "update"},

		// Payment permissions
		{"name": "payments.read", "description": "View payments", "resource": "payments", "action": "read"},
		{"name": "payments.process", "description": "Process payments", "resource": "payments", "action": "process"},
		{"name": "payments.refund", "description": "Refund payments", "resource": "payments", "action": "refund"},
	}

	for _, permission := range permissions {
		result := db.Exec(`
			INSERT INTO permissions (name, description, resource, action, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
			ON CONFLICT (resource, action) DO NOTHING
		`, permission["name"], permission["description"], permission["resource"], permission["action"])

		if result.Error != nil {
			return result.Error
		}
	}

	// Assign permissions to roles
	rolePermissions := []map[string]interface{}{
		// Super admin gets all permissions
		{"role": "super_admin", "resource": "products", "action": "create"},
		{"role": "super_admin", "resource": "products", "action": "read"},
		{"role": "super_admin", "resource": "products", "action": "update"},
		{"role": "super_admin", "resource": "products", "action": "delete"},
		{"role": "super_admin", "resource": "orders", "action": "create"},
		{"role": "super_admin", "resource": "orders", "action": "read"},
		{"role": "super_admin", "resource": "orders", "action": "update"},
		{"role": "super_admin", "resource": "orders", "action": "delete"},
		{"role": "super_admin", "resource": "orders", "action": "fulfill"},
		{"role": "super_admin", "resource": "users", "action": "create"},
		{"role": "super_admin", "resource": "users", "action": "read"},
		{"role": "super_admin", "resource": "users", "action": "update"},
		{"role": "super_admin", "resource": "users", "action": "delete"},
		{"role": "super_admin", "resource": "analytics", "action": "read"},
		{"role": "super_admin", "resource": "analytics", "action": "export"},
		{"role": "super_admin", "resource": "settings", "action": "read"},
		{"role": "super_admin", "resource": "settings", "action": "update"},
		{"role": "super_admin", "resource": "payments", "action": "read"},
		{"role": "super_admin", "resource": "payments", "action": "process"},
		{"role": "super_admin", "resource": "payments", "action": "refund"},

		// Admin gets most permissions except user management
		{"role": "admin", "resource": "products", "action": "create"},
		{"role": "admin", "resource": "products", "action": "read"},
		{"role": "admin", "resource": "products", "action": "update"},
		{"role": "admin", "resource": "products", "action": "delete"},
		{"role": "admin", "resource": "orders", "action": "create"},
		{"role": "admin", "resource": "orders", "action": "read"},
		{"role": "admin", "resource": "orders", "action": "update"},
		{"role": "admin", "resource": "orders", "action": "fulfill"},
		{"role": "admin", "resource": "categories", "action": "create"},
		{"role": "admin", "resource": "categories", "action": "read"},
		{"role": "admin", "resource": "categories", "action": "update"},
		{"role": "admin", "resource": "categories", "action": "delete"},
		{"role": "admin", "resource": "analytics", "action": "read"},
		{"role": "admin", "resource": "analytics", "action": "export"},
		{"role": "admin", "resource": "settings", "action": "read"},
		{"role": "admin", "resource": "settings", "action": "update"},
		{"role": "admin", "resource": "payments", "action": "read"},
		{"role": "admin", "resource": "payments", "action": "process"},
		{"role": "admin", "resource": "payments", "action": "refund"},

		// Merchant gets basic product and order management
		{"role": "merchant", "resource": "products", "action": "create"},
		{"role": "merchant", "resource": "products", "action": "read"},
		{"role": "merchant", "resource": "products", "action": "update"},
		{"role": "merchant", "resource": "orders", "action": "read"},
		{"role": "merchant", "resource": "orders", "action": "update"},
		{"role": "merchant", "resource": "orders", "action": "fulfill"},
		{"role": "merchant", "resource": "categories", "action": "read"},
		{"role": "merchant", "resource": "analytics", "action": "read"},
		{"role": "merchant", "resource": "payments", "action": "read"},

		// Customer gets basic read permissions
		{"role": "customer", "resource": "products", "action": "read"},
		{"role": "customer", "resource": "categories", "action": "read"},
		{"role": "customer", "resource": "orders", "action": "create"},
		{"role": "customer", "resource": "orders", "action": "read"}, // Only their own orders
	}

	for _, rolePermission := range rolePermissions {
		result := db.Exec(`
			INSERT INTO role_permissions (role, permission_id, created_at, updated_at)
			SELECT $1, p.id, NOW(), NOW()
			FROM permissions p
			WHERE p.resource = $2 AND p.action = $3
			ON CONFLICT (role, permission_id) DO NOTHING
		`, rolePermission["role"], rolePermission["resource"], rolePermission["action"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d permissions with role assignments", len(permissions))
	return nil
}

// seedAnalyticsData creates demo analytics data
func seedAnalyticsData(db *gorm.DB) error {
	log.Println("Seeding analytics data...")

	// Seed page views
	pageViews := []map[string]interface{}{
		{
			"tenant_id":  "11111111-1111-1111-1111-111111111111",
			"user_id":    "cccccccc-cccc-cccc-cccc-cccccccccccc",
			"session_id": "sess_001",
			"page_url":   "/products/iphone-15-pro",
			"page_title": "iPhone 15 Pro | TechHub",
			"referrer":   "https://google.com",
			"ip_address": "192.168.1.100",
			"user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"utm_source": "google",
			"utm_medium": "organic",
		},
		{
			"tenant_id":  "11111111-1111-1111-1111-111111111111",
			"user_id":    "dddddddd-dddd-dddd-dddd-dddddddddddd",
			"session_id": "sess_002",
			"page_url":   "/products/macbook-air-m3",
			"page_title": "MacBook Air M3 | TechHub",
			"referrer":   "https://facebook.com",
			"ip_address": "192.168.1.101",
			"user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
			"utm_source": "facebook",
			"utm_medium": "social",
		},
		{
			"tenant_id":  "22222222-2222-2222-2222-222222222222",
			"user_id":    "ffffffff-ffff-ffff-ffff-ffffffffffff",
			"session_id": "sess_003",
			"page_url":   "/products/premium-cotton-t-shirt",
			"page_title": "Premium Cotton T-Shirt | Fashion Forward",
			"referrer":   "",
			"ip_address": "192.168.1.102",
			"user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15",
			"utm_source": "direct",
			"utm_medium": "none",
		},
	}

	for _, view := range pageViews {
		result := db.Exec(`
			INSERT INTO page_views (tenant_id, user_id, session_id, page_url, page_title, referrer, ip_address, user_agent, utm_source, utm_medium, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW() - INTERVAL '2 days')
		`, view["tenant_id"], view["user_id"], view["session_id"], view["page_url"], view["page_title"], view["referrer"], view["ip_address"], view["user_agent"], view["utm_source"], view["utm_medium"])

		if result.Error != nil {
			return result.Error
		}
	}

	// Seed product views
	productViews := []map[string]interface{}{
		{
			"tenant_id":  "11111111-1111-1111-1111-111111111111",
			"product_id": "prod1111-1111-1111-1111-111111111111", // iPhone 15 Pro
			"user_id":    "cccccccc-cccc-cccc-cccc-cccccccccccc",
			"session_id": "sess_001",
			"referrer":   "https://google.com",
			"ip_address": "192.168.1.100",
			"user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		},
		{
			"tenant_id":  "11111111-1111-1111-1111-111111111111",
			"product_id": "prod1112-1111-1111-1111-111111111111", // MacBook Air M3
			"user_id":    "dddddddd-dddd-dddd-dddd-dddddddddddd",
			"session_id": "sess_002",
			"referrer":   "https://facebook.com",
			"ip_address": "192.168.1.101",
			"user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		},
		{
			"tenant_id":  "22222222-2222-2222-2222-222222222222",
			"product_id": "prod2221-2222-2222-2222-222222222222", // T-Shirt
			"user_id":    "ffffffff-ffff-ffff-ffff-ffffffffffff",
			"session_id": "sess_003",
			"referrer":   "",
			"ip_address": "192.168.1.102",
			"user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15",
		},
	}

	for _, view := range productViews {
		result := db.Exec(`
			INSERT INTO product_views (tenant_id, product_id, user_id, session_id, referrer, ip_address, user_agent, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW() - INTERVAL '2 days')
		`, view["tenant_id"], view["product_id"], view["user_id"], view["session_id"], view["referrer"], view["ip_address"], view["user_agent"])

		if result.Error != nil {
			return result.Error
		}
	}

	// Seed analytics events
	events := []map[string]interface{}{
		{
			"tenant_id":    "11111111-1111-1111-1111-111111111111",
			"event_type":   "user_action",
			"event_name":   "add_to_cart",
			"properties":   `{"product_id": "prod1111-1111-1111-1111-111111111111", "variant_id": null, "quantity": 1, "price": 129999.00}`,
			"user_id":      "cccccccc-cccc-cccc-cccc-cccccccccccc",
			"session_id":   "sess_001",
			"ip_address":   "192.168.1.100",
			"user_agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"utm_source":   "google",
			"utm_medium":   "organic",
			"utm_campaign": "tech_deals",
		},
		{
			"tenant_id":    "11111111-1111-1111-1111-111111111111",
			"event_type":   "ecommerce",
			"event_name":   "purchase",
			"properties":   `{"order_id": "ord11111-1111-1111-1111-111111111111", "revenue": 155498.00, "currency": "BDT", "items": [{"product_id": "prod1111-1111-1111-1111-111111111111", "quantity": 1, "price": 129999.00}]}`,
			"user_id":      "cccccccc-cccc-cccc-cccc-cccccccccccc",
			"session_id":   "sess_001",
			"ip_address":   "192.168.1.100",
			"user_agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"utm_source":   "google",
			"utm_medium":   "organic",
			"utm_campaign": "tech_deals",
		},
		{
			"tenant_id":  "22222222-2222-2222-2222-222222222222",
			"event_type": "user_action",
			"event_name": "wishlist_add",
			"properties": `{"product_id": "prod2222-2222-2222-2222-222222222222"}`,
			"user_id":    "ffffffff-ffff-ffff-ffff-ffffffffffff",
			"session_id": "sess_003",
			"ip_address": "192.168.1.102",
			"user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15",
			"utm_source": "direct",
			"utm_medium": "none",
		},
	}

	for _, event := range events {
		result := db.Exec(`
			INSERT INTO analytics_events (tenant_id, event_type, event_name, properties, user_id, session_id, ip_address, user_agent, referrer, utm_source, utm_medium, utm_campaign, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, '', $9, $10, $11, NOW() - INTERVAL '2 days')
		`, event["tenant_id"], event["event_type"], event["event_name"], event["properties"], event["user_id"], event["session_id"], event["ip_address"], event["user_agent"], event["utm_source"], event["utm_medium"], event["utm_campaign"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded analytics data: %d page views, %d product views, %d events", len(pageViews), len(productViews), len(events))
	return nil
}

// seedPaymentData creates demo payment data
func seedPaymentData(db *gorm.DB) error {
	log.Println("Seeding payment data...")

	// Seed payments for existing orders
	payments := []map[string]interface{}{
		{
			"tenant_id":         "11111111-1111-1111-1111-111111111111",
			"order_id":          "ord11111-1111-1111-1111-111111111111",
			"user_id":           "cccccccc-cccc-cccc-cccc-cccccccccccc",
			"payment_intent_id": "pi_tech001_stripe",
			"payment_method_id": "pm_card_visa_001",
			"amount":            155498.00,
			"currency":          "BDT",
			"status":            "succeeded",
			"gateway":           "sslcommerz",
			"gateway_response":  `{"transaction_id": "TXN_TECH_001", "status": "SUCCESS", "card_type": "VISA", "card_no": "4***********1234"}`,
			"processed_at":      "NOW() - INTERVAL '3 days'",
		},
		{
			"tenant_id":         "11111111-1111-1111-1111-111111111111",
			"order_id":          "ord11112-1111-1111-1111-111111111111",
			"user_id":           "dddddddd-dddd-dddd-dddd-dddddddddddd",
			"payment_intent_id": "pi_tech002_stripe",
			"payment_method_id": "pm_card_master_001",
			"amount":            24299.00,
			"currency":          "BDT",
			"status":            "succeeded",
			"gateway":           "sslcommerz",
			"gateway_response":  `{"transaction_id": "TXN_TECH_002", "status": "SUCCESS", "card_type": "MASTERCARD", "card_no": "5***********6789"}`,
			"processed_at":      "NOW() - INTERVAL '2 days'",
		},
		{
			"tenant_id":         "22222222-2222-2222-2222-222222222222",
			"order_id":          "ord22221-2222-2222-2222-222222222222",
			"user_id":           "ffffffff-ffff-ffff-ffff-ffffffffffff",
			"payment_intent_id": "pi_fashion001_stripe",
			"payment_method_id": "pm_mobile_bkash_001",
			"amount":            4698.00,
			"currency":          "BDT",
			"status":            "succeeded",
			"gateway":           "bkash",
			"gateway_response":  `{"transaction_id": "TXN_FASHION_001", "status": "SUCCESS", "payment_method": "bKash", "sender_msisdn": "01*********"}`,
			"processed_at":      "NOW() - INTERVAL '1 days'",
		},
	}

	for _, payment := range payments {
		result := db.Exec(`
			INSERT INTO payments (tenant_id, order_id, user_id, payment_intent_id, payment_method_id, amount, currency, status, gateway, gateway_response, processed_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days')
		`, payment["tenant_id"], payment["order_id"], payment["user_id"], payment["payment_intent_id"], payment["payment_method_id"], payment["amount"], payment["currency"], payment["status"], payment["gateway"], payment["gateway_response"])

		if result.Error != nil {
			return result.Error
		}
	}

	// Seed payment methods
	paymentMethods := []map[string]interface{}{
		{
			"tenant_id":    "11111111-1111-1111-1111-111111111111",
			"user_id":      "cccccccc-cccc-cccc-cccc-cccccccccccc",
			"type":         "card",
			"provider":     "sslcommerz",
			"provider_id":  "pm_card_visa_001",
			"last4":        "1234",
			"brand":        "visa",
			"expiry_month": 12,
			"expiry_year":  2026,
			"is_default":   true,
			"is_active":    true,
		},
		{
			"tenant_id":    "11111111-1111-1111-1111-111111111111",
			"user_id":      "dddddddd-dddd-dddd-dddd-dddddddddddd",
			"type":         "card",
			"provider":     "sslcommerz",
			"provider_id":  "pm_card_master_001",
			"last4":        "6789",
			"brand":        "mastercard",
			"expiry_month": 8,
			"expiry_year":  2025,
			"is_default":   true,
			"is_active":    true,
		},
		{
			"tenant_id":   "22222222-2222-2222-2222-222222222222",
			"user_id":     "ffffffff-ffff-ffff-ffff-ffffffffffff",
			"type":        "mobile_wallet",
			"provider":    "bkash",
			"provider_id": "pm_mobile_bkash_001",
			"last4":       "5678",
			"brand":       "bkash",
			"is_default":  true,
			"is_active":   true,
		},
	}

	for _, method := range paymentMethods {
		result := db.Exec(`
			INSERT INTO payment_methods (tenant_id, user_id, type, provider, provider_id, last4, brand, expiry_month, expiry_year, is_default, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		`, method["tenant_id"], method["user_id"], method["type"], method["provider"], method["provider_id"], method["last4"], method["brand"], method["expiry_month"], method["expiry_year"], method["is_default"], method["is_active"])

		if result.Error != nil {
			return result.Error
		}
	}

	log.Printf("Seeded %d payments and %d payment methods", len(payments), len(paymentMethods))
	return nil
}
