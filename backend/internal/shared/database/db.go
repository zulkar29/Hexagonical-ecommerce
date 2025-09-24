package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ecommerce-saas/internal/shared/config"
)

// DB holds the database connection
var DB *gorm.DB

// Note: NoOpMigrator removed - no longer needed since we use raw SQL migrations exclusively

// Connect establishes database connection
func Connect(cfg *config.Config) (*gorm.DB, error) {
	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.App.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Configure GORM to prevent any schema modifications
	gormConfig := &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		DisableAutomaticPing: true,
		SkipDefaultTransaction: true,
		CreateBatchSize: 0,
		QueryFields: true,
		DryRun: false,
		DisableNestedTransaction: true,
		AllowGlobalUpdate: false,
		PrepareStmt: false,
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// NoOpMigrator is defined above but not used here since we use
// raw SQL migrations exclusively for schema management

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

// Note: Database migrations are handled exclusively by raw SQL files
// in the /migrations directory via RunMigrations function

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

// Note: Test database setup is handled by testhelpers package
// which uses PostgreSQL containers for consistent testing environment

// Seed populates the database with initial data
func Seed(db *gorm.DB) error {
	log.Println("Seeding database...")

	// Use the comprehensive seed data implementation
	if err := SeedData(db); err != nil {
		return fmt.Errorf("failed to seed database: %w", err)
	}

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

// ResetDatabase drops all tables and recreates the schema
func ResetDatabase(db *gorm.DB) error {
	log.Println("Dropping all tables...")
	
	// Get the underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	// Drop all tables in the public schema
	dropTablesSQL := `
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
				EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`
	
	if _, err := sqlDB.Exec(dropTablesSQL); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}
	
	log.Println("All tables dropped successfully")
	
	// Run migrations to recreate schema
	log.Println("Recreating schema with migrations...")
	if err := RunMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations after reset: %w", err)
	}
	
	log.Println("Database reset completed successfully")
	return nil
}

// BulkInsert performs bulk insert operation for better performance
func BulkInsert(db *gorm.DB, models interface{}) error {
	return db.CreateInBatches(models, 100).Error
}

// BulkUpdate performs bulk update operation
func BulkUpdate(db *gorm.DB, model interface{}, updates map[string]interface{}) error {
	return db.Model(model).Updates(updates).Error
}

// GetTableStats returns statistics for a given table
func GetTableStats(db *gorm.DB, tableName string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Get row count
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM ?", tableName).Scan(&count).Error; err != nil {
		return nil, err
	}
	stats["row_count"] = count
	
	// Get table size
	var size string
	if err := db.Raw("SELECT pg_size_pretty(pg_total_relation_size(?)) as size", tableName).Scan(&size).Error; err != nil {
		return nil, err
	}
	stats["table_size"] = size
	
	return stats, nil
}

// BackupPostgreSQL creates a database backup using pg_dump
func BackupPostgreSQL(db *gorm.DB, path string) error {
	// This would require executing pg_dump command
	// Implementation depends on system requirements
	return fmt.Errorf("backup functionality not implemented - use pg_dump directly")
}
