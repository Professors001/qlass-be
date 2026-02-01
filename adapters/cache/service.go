package cache 

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService provides cache operations
type CacheService struct {
	client *redis.Client
}

// NewCacheService creates a new cache service
func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{
		client: client,
	}
}

// Set stores a value in cache with TTL
func (cs *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return cs.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value from cache
func (cs *CacheService) Get(ctx context.Context, key string) (string, error) {
	return cs.client.Get(ctx, key).Result()
}

// Delete removes a key from cache
func (cs *CacheService) Delete(ctx context.Context, keys ...string) error {
	return cs.client.Del(ctx, keys...).Err()
}

// Exists checks if key exists in cache
func (cs *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := cs.client.Exists(ctx, key).Result()
	return result > 0, err
}

// FlushAll clears all cache
func (cs *CacheService) FlushAll(ctx context.Context) error {
	return cs.client.FlushAll(ctx).Err()
}

// Close closes the Redis connection
func (cs *CacheService) Close() error {
	return cs.client.Close()
}