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

	// Run migrations
	// Temporarily disabled to avoid GORM schema conflicts
	// if parseErr := database.AutoMigrate(db); parseErr != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.App.Name)

	// Create server
	srv := server.New(cfg, db, jwtManager)

	// Start server
	if parseErr := srv.Start(); parseErr != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
