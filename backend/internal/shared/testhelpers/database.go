package testhelpers

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ecommerce-saas/internal/shared/database"
)

// TestDatabase represents a test database instance
type TestDatabase struct {
	Container testcontainers.Container
	DB        *gorm.DB
	DSN       string
}

// TestDatabaseOptions configures test database setup
type TestDatabaseOptions struct {
	RunMigrations bool // Whether to run file-based migrations
}

// DefaultTestDatabaseOptions returns default configuration
func DefaultTestDatabaseOptions() TestDatabaseOptions {
	return TestDatabaseOptions{
		RunMigrations: true, // Default to true for backward compatibility
	}
}

// SimpleTestDatabaseOptions returns configuration for simple tests without migrations
func SimpleTestDatabaseOptions() TestDatabaseOptions {
	return TestDatabaseOptions{
		RunMigrations: false, // Skip migrations for simple integration tests
	}
}

// SetupTestDatabase creates a test database using testcontainers with default options
func SetupTestDatabase(t *testing.T) *TestDatabase {
	return SetupTestDatabaseWithOptions(t, DefaultTestDatabaseOptions())
}

// SetupSimpleTestDatabase creates a test database without migrations for backward compatibility
func SetupSimpleTestDatabase(t *testing.T) *TestDatabase {
	return SetupTestDatabaseWithOptions(t, SimpleTestDatabaseOptions())
}

// SetupTestDatabaseWithOptions creates a test database using testcontainers with custom options
func SetupTestDatabaseWithOptions(t *testing.T, options TestDatabaseOptions) *TestDatabase {
	t.Helper()
	ctx := context.Background()

	// Create PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Get connection details
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Construct DSN
	dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable",
		host, port.Port())

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable logging in tests
	})
	require.NoError(t, err)

	// Conditionally run migrations based on options
	if options.RunMigrations {
		err = database.RunMigrations(db)
		require.NoError(t, err)
	}

	return &TestDatabase{
		Container: container,
		DB:        db,
		DSN:       dsn,
	}
}

// TeardownTestDatabase cleans up the test database
func (td *TestDatabase) TeardownTestDatabase(t *testing.T) {
	ctx := context.Background()

	// Close database connection
	if sqlDB, err := td.DB.DB(); err == nil {
		sqlDB.Close()
	}

	// Remove container
	if td.Container != nil {
		if err := td.Container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}
}

// TeardownSimpleTestDatabase is an alias for backward compatibility
func (td *TestDatabase) TeardownSimpleTestDatabase(t *testing.T) {
	td.TeardownTestDatabase(t)
}

// CleanupTables truncates all tables for test isolation
func (td *TestDatabase) CleanupTables(t *testing.T) {
	tables := []string{
		"role_permissions",
		"permissions",
		"analytics_events",
		"page_views",
		"product_views",
		"payments",
		"payment_methods",
		"refunds",
		"order_items",
		"shipping_addresses",
		"orders",
		"product_variants",
		"products",
		"categories",
		"notification_logs",
		"notification_preferences",
		"notifications",
		"notification_templates",
		"user_sessions",
		"users",
		"tenants",
	}

	// Disable foreign key checks temporarily
	td.DB.Exec("SET session_replication_role = replica;")

	for _, table := range tables {
		result := td.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table))
		if result.Error != nil {
			t.Logf("Warning: Could not truncate table %s: %v", table, result.Error)
		}
	}

	// Re-enable foreign key checks
	td.DB.Exec("SET session_replication_role = DEFAULT;")
}

// InMemoryTestDB creates an in-memory SQLite database for unit tests
func InMemoryTestDB(t *testing.T) *gorm.DB {
	// For simple unit tests, we might want to use SQLite
	// However, since we're using PostgreSQL-specific features, we'll use the container approach
	return SetupTestDatabase(t).DB
}

// BeginTx starts a database transaction for isolated testing
func (td *TestDatabase) BeginTx(t *testing.T) *gorm.DB {
	tx := td.DB.Begin()
	require.NoError(t, tx.Error)

	// Rollback transaction in cleanup
	t.Cleanup(func() {
		tx.Rollback()
	})

	return tx
}

// AssertDatabaseClean verifies that test tables are empty
func (td *TestDatabase) AssertDatabaseClean(t *testing.T) {
	tables := []string{"users", "tenants", "products", "orders"}

	for _, table := range tables {
		var count int64
		err := td.DB.Table(table).Count(&count).Error
		require.NoError(t, err)
		require.Equal(t, int64(0), count, "Table %s should be empty after cleanup", table)
	}
}

// SeedTestData populates database with minimal test data
func (td *TestDatabase) SeedTestData(t *testing.T) {
	// Create test tenant
	result := td.DB.Exec(`
		INSERT INTO tenants (id, name, subdomain, status, plan, created_at, updated_at)
		VALUES ('11111111-1111-1111-1111-111111111111', 'Test Store', 'test', 'active', 'starter', NOW(), NOW())
	`)
	require.NoError(t, result.Error)

	// Create test user
	result = td.DB.Exec(`
		INSERT INTO users (id, tenant_id, email, password, first_name, last_name, role, status, email_verified, created_at, updated_at)
		VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 'test@example.com', '$2a$10$hash', 'Test', 'User', 'admin', 'active', true, NOW(), NOW())
	`)
	require.NoError(t, result.Error)
}

// ExecuteInTransaction runs a function within a database transaction
func (td *TestDatabase) ExecuteInTransaction(t *testing.T, fn func(*gorm.DB) error) {
	tx := td.DB.Begin()
	require.NoError(t, tx.Error)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			t.Fatalf("Transaction panicked: %v", r)
		}
	}()

	err := fn(tx)
	if err != nil {
		tx.Rollback()
		require.NoError(t, err)
	} else {
		require.NoError(t, tx.Commit().Error)
	}
}

// WaitForDatabase waits for database to be ready
func WaitForDatabase(dsn string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		db, err := sql.Open("postgres", dsn)
		if err == nil {
			if err = db.Ping(); err == nil {
				db.Close()
				return nil
			}
			db.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("database not ready within %v", timeout)
}
