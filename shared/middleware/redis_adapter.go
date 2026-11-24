package middleware

// RedisAdapter adapts the auth-service RedisClient to work with rate limiter and other middleware
// This allows the existing Redis client to be used with all middleware components

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

// RedisClientAdapter wraps any Redis client that has Get/Set methods
type RedisClientAdapter struct {
	client interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, val string, ttl time.Duration) error
	}
}

// NewRedisClientAdapter creates an adapter for the existing Redis client
func NewRedisClientAdapter(client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string, ttl time.Duration) error
}) *RedisClientAdapter {
	return &RedisClientAdapter{client: client}
}

// Get implements RedisStore interface for rate limiter
func (r *RedisClientAdapter) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key)
}

// Set implements RedisStore interface for rate limiter
func (r *RedisClientAdapter) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration)
}

// Incr implements RedisStore interface for rate limiter
// Since the original client doesn't have Incr, we simulate it with Get + Set
func (r *RedisClientAdapter) Incr(ctx context.Context, key string) (int64, error) {
	// Try to get current value
	val, err := r.client.Get(ctx, key)
	if err != nil {
		// Key doesn't exist, set to 1
		if err := r.client.Set(ctx, key, "1", 0); err != nil {
			return 0, err
		}
		return 1, nil
	}

	// Parse and increment
	var count int64
	if _, err := fmt.Sscanf(val, "%d", &count); err != nil {
		count = 0
	}
	count++

	// Set new value
	if err := r.client.Set(ctx, key, fmt.Sprintf("%d", count), 0); err != nil {
		return 0, err
	}

	return count, nil
}

// Expire implements RedisStore interface for rate limiter
func (r *RedisClientAdapter) Expire(ctx context.Context, key string, expiration time.Duration) error {
	// Get current value, then set with expiration
	val, err := r.client.Get(ctx, key)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, val, expiration)
}

// RedisTokenStoreAdapter wraps Echo context-based Get for auth middleware compatibility
type RedisTokenStoreAdapter struct {
	client interface {
		Get(ctx context.Context, key string) (string, error)
	}
}

// NewRedisTokenStoreAdapter creates an adapter for auth middleware
func NewRedisTokenStoreAdapter(client interface {
	Get(ctx context.Context, key string) (string, error)
}) *RedisTokenStoreAdapter {
	return &RedisTokenStoreAdapter{client: client}
}

// Get implements RedisTokenStore interface for auth middleware
func (r *RedisTokenStoreAdapter) Get(c echo.Context, key string) (string, error) {
	return r.client.Get(c.Request().Context(), key)
}
