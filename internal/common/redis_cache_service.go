package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCacheService implements CacheInterface using Redis
type RedisCacheService struct {
	client *redis.Client
	ctx    context.Context
}

// Ensure RedisCacheService implements CacheInterface
var _ CacheInterface = (*RedisCacheService)(nil)

// NewRedisCacheService creates a new Redis-based cache service
func NewRedisCacheService() (*RedisCacheService, error) {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	// No password by default for local development

	redisDB := 0 // Default DB

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password:     redisPassword,
		DB:           redisDB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCacheService{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a value in Redis with the given key and duration
func (r *RedisCacheService) Set(key string, value interface{}, duration time.Duration) {
	// Serialize value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		// Log error but don't crash
		fmt.Printf("Redis cache: failed to marshal value for key %s: %v\n", key, err)
		return
	}

	if err := r.client.Set(r.ctx, key, data, duration).Err(); err != nil {
		fmt.Printf("Redis cache: failed to set key %s: %v\n", key, err)
	}
}

// Get retrieves a value from Redis by key
func (r *RedisCacheService) Get(key string) (interface{}, bool) {
	data, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		// Key not found
		return nil, false
	}
	if err != nil {
		fmt.Printf("Redis cache: failed to get key %s: %v\n", key, err)
		return nil, false
	}

	// Unmarshal JSON back to generic interface{}
	var result interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		fmt.Printf("Redis cache: failed to unmarshal value for key %s: %v\n", key, err)
		return nil, false
	}

	return result, true
}

// Delete removes a value from Redis by key
func (r *RedisCacheService) Delete(key string) {
	if err := r.client.Del(r.ctx, key).Err(); err != nil {
		fmt.Printf("Redis cache: failed to delete key %s: %v\n", key, err)
	}
}

// GetOrSet retrieves a value from cache, or loads it using the loader function if not found
func (r *RedisCacheService) GetOrSet(
	key string,
	duration time.Duration,
	loader func() (any, error),
) (interface{}, error) {
	// Try to get from cache first
	if val, found := r.Get(key); found {
		return val, nil
	}

	// Load value
	val, err := loader()
	if err != nil {
		return nil, err
	}

	// Store in cache
	r.Set(key, val, duration)

	return val, nil
}

// Close closes the Redis connection
func (r *RedisCacheService) Close() error {
	return r.client.Close()
}

// FlushAll clears all keys in the current Redis database (use with caution!)
func (r *RedisCacheService) FlushAll() error {
	return r.client.FlushDB(r.ctx).Err()
}

// Keys returns all keys matching a pattern
func (r *RedisCacheService) Keys(pattern string) ([]string, error) {
	return r.client.Keys(r.ctx, pattern).Result()
}

// TTL returns the remaining time to live of a key
func (r *RedisCacheService) TTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

// Exists checks if a key exists
func (r *RedisCacheService) Exists(keys ...string) (int64, error) {
	return r.client.Exists(r.ctx, keys...).Result()
}
