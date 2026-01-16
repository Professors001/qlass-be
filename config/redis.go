package config

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient initializes and returns a Redis client
func NewRedisClient(cfg *Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	log.Printf("✅ Redis client initialized: %s:%s", cfg.RedisHost, cfg.RedisPort)
	return client
}
