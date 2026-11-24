# Shared Middleware

Common middleware components for LahziCart microservices. All middleware is built for Echo v4 framework and can be used across any service without modification.

## Available Middleware

### 1. Authentication (`auth.go`)
**Existing middleware** - JWT authentication with Redis token validation.

```go
import (
    "github.com/Rizwank123/LahziCart/shared/middleware"
    "github.com/Rizwank123/LahziCart/shared/utils"
)

authMiddleware := middleware.NewAuthMiddleware(jwtManager, redisClient, logger)
e.Use(authMiddleware.RequireAuth)
```

### 2. Request ID (`request_id.go`)
Generates and propagates unique request IDs for distributed tracing.

```go
requestIDMiddleware := middleware.NewRequestIDMiddleware("X-Request-ID")
e.Use(requestIDMiddleware.Inject)
```

### 3. Logging (`logging.go`)
Structured HTTP request/response logging with zap.

```go
loggingMiddleware := middleware.NewLoggingMiddleware(logger, "/health", "/metrics")
e.Use(loggingMiddleware.Log)
```

### 4. CORS (`cors.go`)
Cross-origin resource sharing configuration.

```go
corsConfig := middleware.CORSConfig{
    AllowOrigins: []string{"https://example.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}
e.Use(middleware.NewCORSMiddleware(corsConfig))
```

### 5. Recovery (`recovery.go`)
Recovers from panics and returns graceful 500 errors.

```go
recoveryMiddleware := middleware.NewRecoveryMiddleware(logger)
e.Use(recoveryMiddleware.Recover)
```

### 6. Error Handling (`error_handler.go`)
Centralized error handling with consistent responses.

```go
errorHandler := middleware.NewErrorHandler(logger, false) // false = hide internal errors in production
e.Use(errorHandler.Handle)

// OR set as global error handler
e.HTTPErrorHandler = errorHandler.CustomErrorHandler()
```

### 7. Timeout (`timeout.go`)
Enforces request timeout limits.

```go
timeoutMiddleware := middleware.NewTimeoutMiddleware(30 * time.Second)
e.Use(timeoutMiddleware.Enforce)
```

### 8. Rate Limiting (`rate_limiter.go`)
Redis-backed rate limiting with configurable windows.

```go
rateLimitConfig := middleware.RateLimitConfig{
    RequestsPerWindow: 100,
    Window: 1 * time.Minute,
    KeyExtractor: middleware.DefaultKeyExtractor, // by IP
}
rateLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfig, logger)
e.Use(rateLimiter.Limit)

// For per-user rate limiting (requires auth middleware first)
rateLimitConfig.KeyExtractor = middleware.UserKeyExtractor
```

### 9. Metrics (`metrics.go`)
Prometheus HTTP metrics collection.

```go
metricsMiddleware := middleware.NewMetricsMiddleware("auth-service")
e.Use(metricsMiddleware.Collect)
```

### 10. Validation (`validator.go`)
Request validation using go-playground/validator.

```go
// Set validator on Echo instance
e.Validator = middleware.NewValidator()

// In handlers, use the helper function
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func (h *Handler) Register(c echo.Context) error {
    var req RegisterRequest
    if err := middleware.BindAndValidate(c, &req); err != nil {
        return err // Returns formatted 400 error
    }
    // Continue with validated request...
}
```

## Middleware Ordering

The order of middleware matters! Here's the recommended order:

```go
// 1. Recovery - catch panics first
e.Use(recoveryMiddleware.Recover)

// 2. Request ID - for tracing all subsequent middleware
e.Use(requestIDMiddleware.Inject)

// 3. Logging - log with request IDs
e.Use(loggingMiddleware.Log)

// 4. Metrics - track all requests
e.Use(metricsMiddleware.Collect)

// 5. CORS - handle preflight early
e.Use(middleware.NewCORSMiddleware(corsConfig))

// 6. Timeout - enforce time limits
e.Use(timeoutMiddleware.Enforce)

// 7. Rate Limiting - throttle before processing
e.Use(rateLimiter.Limit)

// 8. Authentication - only on protected routes
protectedGroup := e.Group("/api")
protectedGroup.Use(authMiddleware.RequireAuth)
```

## Complete Integration Example

```go
package main

import (
    "time"
    
    "github.com/labstack/echo/v4"
    "go.uber.org/zap"
    
    "github.com/Rizwank123/LahziCart/shared/middleware"
)

func main() {
    logger, _ := zap.NewProduction()
    e := echo.New()
    
    // Set validator
    e.Validator = middleware.NewValidator()
    
    // Global middleware (order matters!)
    recoveryMw := middleware.NewRecoveryMiddleware(logger)
    requestIDMw := middleware.NewRequestIDMiddleware("X-Request-ID")
    loggingMw := middleware.NewLoggingMiddleware(logger, "/health", "/metrics")
    metricsMw := middleware.NewMetricsMiddleware("my-service")
    timeoutMw := middleware.NewTimeoutMiddleware(30 * time.Second)
    
    e.Use(recoveryMw.Recover)
    e.Use(requestIDMw.Inject)
    e.Use(loggingMw.Log)
    e.Use(metricsMw.Collect)
    e.Use(timeoutMw.Enforce)
    
    // CORS for API routes
    corsConfig := middleware.DefaultCORSConfig()
    e.Use(middleware.NewCORSMiddleware(corsConfig))
    
    // Rate limiting (optional, requires Redis)
    if redisClient != nil {
        rateLimitConfig := middleware.RateLimitConfig{
            RequestsPerWindow: 100,
            Window: 1 * time.Minute,
        }
        rateLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfig, logger)
        e.Use(rateLimiter.Limit)
    }
    
    // Protected routes with authentication
    api := e.Group("/api")
    authMw := middleware.NewAuthMiddleware(jwtManager, redisClient, logger)
    api.Use(authMw.RequireAuth)
    
    e.Start(":8080")
}
```

## Features

- ✅ **Echo v4 Compatible** - Works seamlessly with Echo framework
- ✅ **Zero Dependencies on Services** - Completely independent, shareable across all microservices
- ✅ **Production Ready** - Includes error handling, logging, metrics
- ✅ **Composable** - Mix and match middleware as needed
- ✅ **Configurable** - Flexible configuration for different environments
- ✅ **Observable** - Built-in Prometheus metrics and structured logging

## Dependencies

Add these to your service's `go.mod` if not already present:

```bash
go get github.com/labstack/echo/v4
go get go.uber.org/zap
go get github.com/google/uuid
go get github.com/prometheus/client_golang/prometheus
go get github.com/go-playground/validator/v10
```

## Testing

Run tests for all middleware:

```bash
cd shared
go test ./middleware/... -v -cover
```

## Notes

- **Auth middleware** requires `utils.AuthManager` and `RedisTokenStore` - these are already implemented in the project
- **Rate limiter** requires a Redis client implementing the `RedisStore` interface
- **Metrics middleware** automatically registers Prometheus metrics
- **Validator** uses struct tags for validation rules (see go-playground/validator docs)
