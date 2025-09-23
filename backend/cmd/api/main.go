package main

import (
	"log"

	"ecommerce-saas/internal/billing"
	"ecommerce-saas/internal/components"
	"ecommerce-saas/internal/observability"
	"ecommerce-saas/internal/order"
	"ecommerce-saas/internal/security"
	"ecommerce-saas/internal/server"
	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/shared/database"
	"ecommerce-saas/internal/shared/utils"
	"ecommerce-saas/internal/user"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	if err = database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Run module migrations
	log.Println("Running module migrations...")
	if err = runModuleMigrations(db); err != nil {
		log.Fatalf("Failed to run module migrations: %v", err)
	}
	log.Println("Module migrations completed successfully")

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.App.Name)

	// Create server
	srv := server.New(cfg, db, jwtManager)

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// runModuleMigrations runs GORM AutoMigrate for all modules
func runModuleMigrations(db *gorm.DB) error {
	// Initialize modules and run their migrations
	modules := []interface {
		Migrate() error
	}{
		components.NewModule(db),
		observability.NewModule(db),
		billing.NewModule(db, nil, nil, nil, nil), // Pass nil for dependencies during migration
	}

	// Run migrations for modules with Migrate() method
	for _, module := range modules {
		if err := module.Migrate(); err != nil {
			return err
		}
	}

	// Run migrations for modules with Migrate(db) method
	orderModule := order.NewModule(db, nil, nil, nil, nil) // Pass nil for dependencies during migration
	if err := orderModule.Migrate(db); err != nil {
		return err
	}

	// Run migrations for security module (requires user repository)
	userModule := user.NewModule(db, nil) // Pass nil for JWT manager during migration
	securityModule := security.NewModule(db, userModule.GetSecurityRepository())
	if err := securityModule.Migrate(); err != nil {
		return err
	}

	return nil
}
