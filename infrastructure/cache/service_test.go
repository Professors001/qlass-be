package cache

import (
	"context"
	"testing"
	"time"

	"qlass-be/config"
)

// TestRedisIntegration checks if the Go app can talk to the Redis instance.
// Run this with: go test -v ./infrastructure/cache/...
func TestRedisIntegration(t *testing.T) {
	// 1. Load Config
	cfg := config.LoadConfig()

	// Force localhost for local testing.
	cfg.RedisHost = "localhost"

	// 2. Initialize Dependencies
	rdb := config.NewRedisClient(cfg)
	cache := NewCacheService(rdb)
	defer cache.Close()

	ctx := context.Background()
	key := "test_connection"
	value := "redis_is_working"

	// 3. Test SET
	err := cache.Set(ctx, key, value, 10*time.Second)
	if err != nil {
		t.Fatalf("❌ Failed to SET key: %v", err)
	}
	t.Log("✅ SET operation successful")

	// 4. Test GET
	result, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("❌ Failed to GET key: %v", err)
	}

	if result != value {
		t.Errorf("❌ Value mismatch! Expected '%s', got '%s'", value, result)
	} else {
		t.Logf("✅ GET operation successful: %s", result)
	}

	// 5. Cleanup
	_ = cache.Delete(ctx, key)
}