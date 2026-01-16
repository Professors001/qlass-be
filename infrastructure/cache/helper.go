package cache

import (
	"context"
	"encoding/json"
	"time"
)

// CacheHelper provides helper methods for common cache operations
type CacheHelper struct {
	cache *CacheService
}

// NewCacheHelper creates a new cache helper
func NewCacheHelper(cache *CacheService) *CacheHelper {
	return &CacheHelper{
		cache: cache,
	}
}

// SetJSON stores a JSON-serialized value in cache
func (ch *CacheHelper) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ch.cache.Set(ctx, key, data, ttl)
}

// GetJSON retrieves and deserializes a JSON value from cache
func (ch *CacheHelper) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := ch.cache.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}