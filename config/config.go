package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// App Settings
	AppPort string `mapstructure:"APP_PORT"`
	AppEnv  string `mapstructure:"APP_ENV"`

	// Database Settings
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	SupabaseURL string `mapstructure:"SUPABASE_URL"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	DBSSLMode   string `mapstructure:"DB_SSLMODE"`

	// Redis Settings
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisUsername string `mapstructure:"REDIS_USERNAME"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	// Security
	JWTSecret string `mapstructure:"JWT_SECRET"`

	// MinIO Settings
	MinioEndpoint   string `mapstructure:"MINIO_ENDPOINT"`
	MinioRegion     string `mapstructure:"MINIO_REGION"`
	MinioAccessKey  string `mapstructure:"MINIO_ROOT_USER"`
	MinioSecretKey  string `mapstructure:"MINIO_ROOT_PASSWORD"`
	MinioBucketName string `mapstructure:"MINIO_BUCKET_NAME"`
	MinioUseSSL     bool   `mapstructure:"MINIO_USE_SSL"`

	SMTPHost string `mapstructure:"SMTP_HOST"`
	SMTPPort string `mapstructure:"SMTP_PORT"`
	SMTPUser string `mapstructure:"SMTP_USER"`
	SMTPPass string `mapstructure:"SMTP_PASS"`
}

// LoadConfig reads configuration from .env file or environment variables
func LoadConfig() *Config {
	viper.AutomaticEnv() // Automatically read environment variables (docker)

	// 1. Try to load .env from current directory
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		// 2. If failed, try loading from parent directory (for tests)
		viper.SetConfigFile("../.env")
		if err := viper.ReadInConfig(); err != nil {
			// 3. If failed, try loading from two levels up (for nested tests like infrastructure/cache)
			viper.SetConfigFile("../../.env")
			if err := viper.ReadInConfig(); err != nil {
				log.Println("⚠️  No .env file found, relying on system environment variables")
			}
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("❌ Unable to decode into struct, %v", err)
	}

	// Basic validation
	hasDBURL := config.DatabaseURL != "" || config.SupabaseURL != ""
	if !hasDBURL && (config.DBHost == "" || config.DBPort == "") {
		log.Fatal("❌ Database configuration is missing. Set DATABASE_URL/SUPABASE_URL or DB_HOST and DB_PORT.")
	}

	if config.MinioEndpoint == "" {
		log.Fatal("❌ MinIO configuration is missing. Check MINIO_ENDPOINT in your .env file.")
	}

	if config.SMTPHost == "" {
		log.Fatal("❌ SMTP configuration is missing. Check SMTP_HOST in your .env file.")
	}

	// Fix: Redis default user does not need a username.
	// If "root" is set (common mistake), clear it to avoid WRONGPASS error.
	if config.RedisUsername == "root" {
		config.RedisUsername = ""
	}

	return &config
}
