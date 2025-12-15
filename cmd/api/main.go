package main

import (
	"log"
	"qlass-be/config"
	"qlass-be/pkg/database"
	"qlass-be/internal/domain"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Connect to Database
	db := database.NewPostgresDB(cfg)

	// 3. AUTO MIGRATION (Crucial Step)
	// This will create the 'users' table in Postgres based on the struct
	log.Println("Runnning Auto Migration for Users...")
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration successful! Table 'users' created.")

	// ... later we will start the server here
}