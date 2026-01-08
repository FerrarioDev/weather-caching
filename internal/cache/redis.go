// Package cache handles the cache storage with redis
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/FerrarioDev/weathercache/config"
	"github.com/FerrarioDev/weathercache/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

// NewCache creates a new Redis cache client configured from environment variables
func NewCache() *RedisCache {
	// Load configuration from environment
	addr := config.RedisAddr.GetValue()
	if addr == "" {
		addr = "localhost:6379" // default value
	}

	password := config.RedisPassword.GetValue()

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &RedisCache{client: rdb}
}

// Shutdown gracefully closes the Redis connection
func (r *RedisCache) Shutdown() error {
	return r.client.Close()
}

// SetValue stores weather data in cache with the given key and 1-hour TTL
func (r *RedisCache) SetValue(ctx context.Context, key string, weather *domain.TransformedWeatherResponse) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	data, err := json.Marshal(weather)
	if err != nil {
		return fmt.Errorf("failed to marshal weather: %v", err)
	}
	err = r.client.Set(ctx, key, data, 1*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in cache: %v", err)
	}

	return nil
}

func (r *RedisCache) GetValue(ctx context.Context, key string) (*domain.TransformedWeatherResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get value for key %s: %v", key, err)
	}

	var weather domain.TransformedWeatherResponse
	if err := json.Unmarshal([]byte(data), &weather); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal value: %v", err)
	}

	return &weather, nil
}
