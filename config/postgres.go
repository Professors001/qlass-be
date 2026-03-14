package config

import (
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *Config) *gorm.DB {
	// 1. Build Data Source Name (DSN) from SUPABASE_URL only
	dsn := strings.TrimSpace(cfg.SupabaseURL)
	if dsn == "" {
		log.Fatal("❌ SUPABASE_URL is required")
	}

	// Supabase requires SSL. Add it automatically if missing.
	if !strings.Contains(dsn, "sslmode=") {
		if strings.Contains(dsn, "?") {
			dsn += "&sslmode=require"
		} else {
			dsn += "?sslmode=require"
		}
	}

	// 2. Configure GORM Logger (Show SQL logs in Dev, hide in Prod)
	var gormLogger logger.Interface
	if cfg.AppEnv == "development" {
		gormLogger = logger.Default.LogMode(logger.Info) // Show all SQL queries
	} else {
		gormLogger = logger.Default.LogMode(logger.Error) // Only show errors
	}

	// 3. Connect
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		if strings.Contains(err.Error(), "no route to host") {
			log.Fatalf("❌ Failed to connect to Supabase: %v\nHint: your network cannot reach the direct DB host (often IPv6 routing issue). Use the Supabase Session Pooler URL (IPv4) from Project Settings > Database > Connection string.", err)
		}
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to Postgres Database successfully!")
	return db
}
