package main

import (
	"log"

	"ecommerce-saas/internal/server"
	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/shared/database"
	"ecommerce-saas/internal/shared/utils"
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

	// Run raw SQL migrations
	if err = database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Module migrations are now handled by raw SQL migrations
	log.Println("All migrations completed via raw SQL files")

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.App.Name)

	// Create server
	srv := server.New(cfg, db, jwtManager)

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Note: Module migrations are now handled by raw SQL migration files
// in the /migrations directory. No programmatic migrations needed.
