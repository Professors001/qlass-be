package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient initializes and returns a Redis client
func NewRedisClient(cfg *Config) *redis.Client {
	if strings.TrimSpace(cfg.RedisURL) != "" {
		opts, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			log.Fatalf("❌ Invalid REDIS_URL: %v", err)
		}
		client := redis.NewClient(opts)
		log.Printf("✅ Redis client initialized via REDIS_URL")
		return client
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	log.Printf("✅ Redis client initialized: %s:%s", cfg.RedisHost, cfg.RedisPort)
	return client
}
