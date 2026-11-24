// Package middleware demonstrates how to use shared middleware in your microservices
//
// This file contains documentation and example code patterns.
// Copy these examples into your service's main.go file.
package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Example1_BasicSetup shows minimal middleware setup
func Example1_BasicSetup() {
	logger, _ := zap.NewProduction()
	e := echo.New()

	// Recovery and Request ID are the bare minimum
	recovery := NewRecoveryMiddleware(logger)
	requestID := NewRequestIDMiddleware("X-Request-ID")

	e.Use(recovery.Recover)
	e.Use(requestID.Inject)

	// Your routes here...
	e.Start(":8080")
}

// Example2_ProductionSetup shows comprehensive middleware for production
func Example2_ProductionSetup(logger *zap.Logger, redisClient interface{}) {
	e := echo.New()

	// 1. Set validator
	e.Validator = NewValidator()

	// 2. Recovery - catch panics first
	recovery := NewRecoveryMiddleware(logger)
	e.Use(recovery.Recover)

	// 3. Request ID - generate unique IDs for tracing
	requestID := NewRequestIDMiddleware("X-Request-ID")
	e.Use(requestID.Inject)

	// 4. Logging - structured logs with request IDs
	logging := NewLoggingMiddleware(logger, "/health", "/metrics")
	e.Use(logging.Log)

	// 5. Metrics - Prometheus monitoring
	metrics := NewMetricsMiddleware("my-service")
	e.Use(metrics.Collect)

	// 6. CORS - for frontend integration
	corsConfig := CORSConfig{
		AllowOrigins:     []string{"https://example.com", "https://app.example.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           3600,
	}
	e.Use(NewCORSMiddleware(corsConfig))

	// 7. Timeout - prevent hung requests
	timeout := NewTimeoutMiddleware(30 * time.Second)
	e.Use(timeout.Enforce)

	// Your routes here...
	e.Start(":8080")
}

// Example3_WithAuthentication shows auth + rate limiting setup
func Example3_WithAuthentication(logger *zap.Logger, jwtManager interface{}, redisClient interface{}) {
	e := echo.New()

	// Basic middleware
	e.Use(NewRecoveryMiddleware(logger).Recover)
	e.Use(NewRequestIDMiddleware("X-Request-ID").Inject)
	e.Use(NewLoggingMiddleware(logger).Log)

	// Public routes
	e.POST("/register", nil) // your handler
	e.POST("/login", nil)    // your handler

	// Protected routes group
	api := e.Group("/api")

	// Authentication middleware (only on protected routes)
	// Note: Make sure your redisClient implements the RedisTokenStore interface
	// authMw := NewAuthMiddleware(jwtManager, redisClient, logger)
	// api.Use(authMw.RequireAuth)

	// Example protected endpoints
	api.GET("/profile", nil) // your handler
	api.PUT("/profile", nil) // your handler
	api.GET("/orders", nil)  // your handler
	api.POST("/orders", nil) // your handler

	e.Start(":8080")
}

// Example4_RateLimiting shows different rate limiting strategies
func Example4_RateLimiting(logger *zap.Logger, redisClient interface{}) {
	e := echo.New()

	e.Use(NewRecoveryMiddleware(logger).Recover)
	e.Use(NewRequestIDMiddleware("X-Request-ID").Inject)

	// Strategy 1: Global rate limit by IP
	// Note: Make sure your redisClient implements the RedisStore interface
	// globalRateLimitConfig := RateLimitConfig{
	// 	RequestsPerWindow: 100,
	// 	Window:           1 * time.Minute,
	// 	KeyExtractor:     DefaultKeyExtractor, // by IP
	// }
	// globalRateLimiter := NewRateLimiter(redisClient, globalRateLimitConfig, logger)
	// e.Use(globalRateLimiter.Limit)

	// Public routes - already rate limited by IP
	e.POST("/register", nil)
	e.POST("/login", nil)

	// Strategy 2: Stricter rate limit for authenticated users
	api := e.Group("/api")
	// api.Use(authMiddleware.RequireAuth) // authenticate first

	// Per-user rate limiting (requires auth middleware)
	// userRateLimitConfig := RateLimitConfig{
	// 	RequestsPerWindow: 1000,
	// 	Window:           1 * time.Minute,
	// 	KeyExtractor:     UserKeyExtractor, // by user_id
	// }
	// userRateLimiter := NewRateLimiter(redisClient, userRateLimitConfig, logger)
	// api.Use(userRateLimiter.Limit)

	api.GET("/data", nil)

	e.Start(":8080")
}

// Example5_ValidationUsage shows how to use validator middleware
func Example5_ValidationUsage() {
	e := echo.New()

	// Set the validator on Echo instance
	e.Validator = NewValidator()

	// Define request struct with validation tags
	type CreateUserRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		Name     string `json:"name" validate:"required,min=2,max=100"`
		Age      int    `json:"age" validate:"required,gte=18,lte=120"`
	}

	// In your handler, use BindAndValidate
	e.POST("/users", func(c echo.Context) error {
		var req CreateUserRequest

		// Bind and validate in one call
		if err := BindAndValidate(c, &req); err != nil {
			return err // Returns formatted 400 error with validation message
		}

		// Continue with validated request
		// ... your business logic ...

		return c.JSON(200, map[string]string{"message": "user created"})
	})

	e.Start(":8080")
}

// Example6_CustomErrorHandler shows custom error handling
func Example6_CustomErrorHandler(logger *zap.Logger) {
	e := echo.New()

	// Set custom error handler globally
	errorHandler := NewErrorHandler(logger, false) // false = hide errors in production
	e.HTTPErrorHandler = errorHandler.CustomErrorHandler()

	// Now all errors will be handled consistently
	e.GET("/test", func(c echo.Context) error {
		// Any error returned will be processed by the error handler
		return echo.NewHTTPError(400, "invalid request")
	})

	e.Start(":8080")
}

// Example7_MiddlewareOrdering shows the correct order for middleware
func Example7_MiddlewareOrdering(logger *zap.Logger) {
	e := echo.New()

	// IMPORTANT: Order matters!
	//
	// 1. Recovery - must be first to catch panics from all other middleware
	e.Use(NewRecoveryMiddleware(logger).Recover)

	// 2. Request ID - generate IDs before logging
	e.Use(NewRequestIDMiddleware("X-Request-ID").Inject)

	// 3. Logging - log with request IDs
	e.Use(NewLoggingMiddleware(logger, "/health").Log)

	// 4. Metrics - track all requests
	e.Use(NewMetricsMiddleware("my-service").Collect)

	// 5. CORS - handle OPTIONS early
	e.Use(NewCORSMiddleware(DefaultCORSConfig()))

	// 6. Timeout - enforce limits before processing
	e.Use(NewTimeoutMiddleware(30 * time.Second).Enforce)

	// 7. Rate Limiting - throttle before auth/business logic
	// e.Use(rateLimiter.Limit)

	// 8. Authentication - last, only on protected routes
	// protectedGroup.Use(authMiddleware.RequireAuth)

	e.Start(":8080")
}
