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

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	switch *action {
	case "up":
		log.Println("Running migrations...")
		if err := database.AutoMigrate(db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")

	case "seed":
		log.Println("Seeding database...")
		if err := database.Seed(db); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		log.Println("Database seeded successfully")

	case "reset":
		log.Println("Resetting database...")
		// TODO: Implement database reset
		log.Println("Database reset - TODO: implement")

	default:
		log.Printf("Unknown action: %s", *action)
		log.Println("Available actions: up, seed, reset")
	}
}