package main

import (
	"flag"
	"log"

	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/shared/database"
)

func main() {
	var action = flag.String("action", "up", "Migration action: up, down, or reset")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to PostgreSQL database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	switch *action {
	case "up":
		log.Println("Running raw SQL migrations...")
		if parseErr := database.RunMigrations(db); parseErr != nil {
			log.Fatalf("Failed to run migrations: %v", parseErr)
		}
		log.Println("Raw SQL migrations completed successfully")

	case "seed":
		log.Println("Seeding database...")
		if parseErr := database.Seed(db); parseErr != nil {
			log.Fatalf("Failed to seed database: %v", parseErr)
		}
		log.Println("Database seeded successfully")

	case "reset":
		log.Println("Resetting database...")
		if parseErr := database.ResetDatabase(db); parseErr != nil {
			log.Fatalf("Failed to reset database: %v", parseErr)
		}
		log.Println("Database reset completed successfully")

	default:
		log.Printf("Unknown action: %s", *action)
		log.Println("Available actions: up, seed, reset")
	}
}
