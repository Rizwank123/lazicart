package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RedisStore interface for rate limiter storage
type RedisStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	// RequestsPerWindow is the maximum number of requests allowed in the time window
	RequestsPerWindow int

	// Window is the time window for rate limiting (e.g., 1 minute)
	Window time.Duration

	// KeyExtractor is a function to extract the rate limit key from the context
	// Default: uses client IP address
	KeyExtractor func(echo.Context) string
}

// RateLimiter provides rate limiting using token bucket algorithm with Redis
type RateLimiter struct {
	redis  RedisStore
	config RateLimitConfig
	logger *zap.Logger
}

// DefaultKeyExtractor extracts rate limit key from client IP
func DefaultKeyExtractor(c echo.Context) string {
	return c.RealIP()
}

// NewRateLimiter creates a new rate limiter middleware
func NewRateLimiter(redis RedisStore, config RateLimitConfig, logger *zap.Logger) *RateLimiter {
	// Set defaults
	if config.RequestsPerWindow <= 0 {
		config.RequestsPerWindow = 100 // default: 100 requests
	}
	if config.Window <= 0 {
		config.Window = 1 * time.Minute // default: per minute
	}
	if config.KeyExtractor == nil {
		config.KeyExtractor = DefaultKeyExtractor
	}

	return &RateLimiter{
		redis:  redis,
		config: config,
		logger: logger,
	}
}

// Limit is the middleware handler that enforces rate limits
func (r *RateLimiter) Limit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract rate limit key
		key := r.config.KeyExtractor(c)
		redisKey := fmt.Sprintf("rate_limit:%s", key)

		ctx := c.Request().Context()

		// Increment the counter
		count, err := r.redis.Incr(ctx, redisKey)
		if err != nil {
			r.logger.Warn("rate limit redis error, allowing request",
				zap.String("key", key),
				zap.Error(err),
			)
			// On Redis error, allow the request (fail open)
			return next(c)
		}

		// Set expiration on first request
		if count == 1 {
			if err := r.redis.Expire(ctx, redisKey, r.config.Window); err != nil {
				r.logger.Warn("failed to set rate limit expiration",
					zap.String("key", key),
					zap.Error(err),
				)
			}
		}

		// Check if limit exceeded
		if count > int64(r.config.RequestsPerWindow) {
			r.logger.Warn("rate limit exceeded",
				zap.String("key", key),
				zap.Int64("count", count),
				zap.Int("limit", r.config.RequestsPerWindow),
			)

			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(r.config.RequestsPerWindow))
			c.Response().Header().Set("X-RateLimit-Remaining", "0")
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(r.config.Window).Unix(), 10))

			return c.JSON(429, response.Error("rate limit exceeded, please try again later"))
		}

		// Set rate limit headers
		remaining := r.config.RequestsPerWindow - int(count)
		c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(r.config.RequestsPerWindow))
		c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(r.config.Window).Unix(), 10))

		return next(c)
	}
}

// UserKeyExtractor extracts rate limit key from authenticated user ID
// Use this with auth middleware to rate limit per user instead of per IP
func UserKeyExtractor(c echo.Context) string {
	userID, ok := c.Get("user_id").(string)
	if ok && userID != "" {
		return fmt.Sprintf("user:%s", userID)
	}
	// Fallback to IP if user not authenticated
	return c.RealIP()
}
