package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/tenant"
	"ecommerce-saas/internal/product"
	"ecommerce-saas/internal/order"
	"ecommerce-saas/internal/user"
)

// DB holds the database connection
var DB *gorm.DB

// Connect establishes database connection
func Connect(cfg *config.Config) (*gorm.DB, error) {
	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.App.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")

	// Store global reference
	DB = db

	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Migrate all models
	err := db.AutoMigrate(
		// Tenant models
		&tenant.Tenant{},
		
		// Product models
		&product.Product{},
		&product.ProductVariant{},
		&product.Category{},
		
		// Order models (TODO: define these)
		// &order.Order{},
		// &order.OrderItem{},
		
		// User models (TODO: define these)
		// &user.User{},
		// &user.Role{},
		// &user.Permission{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Health checks database health
func Health() error {
	if DB == nil {
		return fmt.Errorf("database not connected")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// Transaction executes a function within a database transaction
func Transaction(fn func(*gorm.DB) error) error {
	return DB.Transaction(fn)
}

// SetupTestDB sets up a test database (for testing)
func SetupTestDB() (*gorm.DB, error) {
	// Use SQLite for testing
	db, err := gorm.Open(postgres.Open("postgres://test:test@localhost/esass_test?sslmode=disable"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := AutoMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

// TearDownTestDB cleans up test database
func TearDownTestDB(db *gorm.DB) error {
	// Drop all tables
	tables := []string{
		"tenants",
		"products",
		"product_variants", 
		"categories",
		// Add more tables as needed
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			return err
		}
	}

	return nil
}

// Seed populates the database with initial data
func Seed(db *gorm.DB) error {
	log.Println("Seeding database...")

	// TODO: Add seed data
	// - Default categories
	// - Sample products
	// - Admin user
	// - Default settings

	log.Println("Database seeding completed")
	return nil
}

// Custom database helpers

// Paginate adds pagination to GORM query
func Paginate(offset, limit int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if offset < 0 {
			offset = 0
		}
		if limit <= 0 || limit > 100 {
			limit = 20
		}
		return db.Offset(offset).Limit(limit)
	}
}

// Search adds search conditions to GORM query
func Search(fields []string, query string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if query == "" || len(fields) == 0 {
			return db
		}

		var conditions []string
		var args []interface{}

		for _, field := range fields {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, "%"+query+"%")
		}

		whereClause := fmt.Sprintf("(%s)", fmt.Sprintf(conditions[0]))
		for i := 1; i < len(conditions); i++ {
			whereClause += fmt.Sprintf(" OR %s", conditions[i])
		}

		return db.Where(whereClause, args...)
	}
}

// TODO: Add more database helpers
// - BulkInsert(db *gorm.DB, models interface{}) error
// - BulkUpdate(db *gorm.DB, model interface{}, updates map[string]interface{}) error
// - GetTableStats(db *gorm.DB, tableName string) (map[string]interface{}, error)
// - BackupDatabase(db *gorm.DB, path string) error
