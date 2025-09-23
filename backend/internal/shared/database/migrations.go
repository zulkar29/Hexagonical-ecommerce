package database

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// RunMigrations executes all pending database migrations from files
// This handles:
// - File-based migration execution
// - Migration tracking and versioning
// - Proper error handling and rollback
func RunMigrations(db *gorm.DB) error {
	// Create migration tracking table if it doesn't exist
	if err := createMigrationTable(db); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Get list of migration files
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Run each migration file
	for _, file := range migrationFiles {
		if err := runMigrationFile(db, file); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}
	}

	return nil
}

// MigrationRecord tracks applied migrations
type MigrationRecord struct {
	ID        uint      `gorm:"primarykey"`
	Version   string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null"`
	AppliedAt time.Time `gorm:"not null"`
}

// createMigrationTable creates the migration tracking table
func createMigrationTable(db *gorm.DB) error {
	// Create migration_records table with raw SQL to ensure it works
	sql := `
		CREATE TABLE IF NOT EXISTS migration_records (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_migration_records_version ON migration_records(version);
	`
	return db.Exec(sql).Error
}

// getMigrationFiles returns a sorted list of migration files
func getMigrationFiles() ([]string, error) {
	// Try multiple possible migration directory paths
	migrationPaths := []string{
		"./migrations",                    // From backend root
		"../migrations",                  // From subdirectory
		"../../migrations",               // From deeper subdirectory
		"../../../migrations",            // From even deeper subdirectory
		"../../../../migrations",         // From very deep subdirectory
	}
	
	var migrationDir string
	for _, path := range migrationPaths {
		if _, err := os.Stat(path); err == nil {
			migrationDir = path
			break
		}
	}
	
	if migrationDir == "" {
		return nil, fmt.Errorf("migration directory not found")
	}

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, filepath.Join(migrationDir, file.Name()))
		}
	}

	// Sort files to ensure proper execution order
	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

// runMigrationFile executes a single migration file if not already applied
func runMigrationFile(db *gorm.DB, filePath string) error {
	// Extract migration ID from filename
	filename := filepath.Base(filePath)
	migrationID := strings.TrimSuffix(filename, ".sql")

	// Check if migration already applied
	var count int64
	db.Model(&MigrationRecord{}).Where("version = ?", migrationID).Count(&count)
	if count > 0 {
		return nil // Already applied
	}

	// Read migration file
	sqlContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", filePath, err)
	}

	// Execute migration in transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// Execute the SQL
		if err := tx.Exec(string(sqlContent)).Error; err != nil {
			return fmt.Errorf("failed to execute migration SQL from %s: %w", filename, err)
		}

		// Record the migration
		record := MigrationRecord{
			Version:   migrationID,
			Name:      filename,
			AppliedAt: time.Now(),
		}
		return tx.Create(&record).Error
	})
}

// Note: All SQL migrations are now handled through .sql files in the migrations/ directory
// This ensures single source of truth and prevents duplication
