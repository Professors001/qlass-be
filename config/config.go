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
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	// Security
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

// LoadConfig reads configuration from .env file or environment variables
func LoadConfig() *Config {
	viper.AddConfigPath(".")    // Look for config in the root directory
	viper.SetConfigFile(".env") // Specifically look for a file named .env

	viper.AutomaticEnv() // Automatically read environment variables (docker)

	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️  No .env file found, relying on system environment variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("❌ Unable to decode into struct, %v", err)
	}

	// Basic validation
	if config.DBHost == "" || config.DBPort == "" {
		log.Fatal("❌ Database configuration is missing. Check your .env file.")
	}

	return &config
}
