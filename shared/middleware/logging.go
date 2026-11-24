package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggingMiddleware provides structured HTTP request/response logging
type LoggingMiddleware struct {
	logger    *zap.Logger
	skipPaths map[string]bool
}

// NewLoggingMiddleware creates a new logging middleware
// skipPaths is an optional list of paths to exclude from logging (e.g., "/health", "/metrics")
func NewLoggingMiddleware(logger *zap.Logger, skipPaths ...string) *LoggingMiddleware {
	skipMap := make(map[string]bool)
	for _, path := range skipPaths {
		skipMap[path] = true
	}
	return &LoggingMiddleware{
		logger:    logger,
		skipPaths: skipMap,
	}
}

// Log is the middleware handler that logs HTTP requests and responses
func (m *LoggingMiddleware) Log(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		path := req.URL.Path

		// Skip logging for configured paths
		if m.skipPaths[path] {
			return next(c)
		}

		start := time.Now()

		// Get request ID from context if available
		requestID, _ := c.Get("request_id").(string)

		// Execute the request
		err := next(c)

		// Calculate duration
		duration := time.Since(start)
		res := c.Response()

		// Log the request/response
		fields := []zap.Field{
			zap.String("method", req.Method),
			zap.String("path", path),
			zap.Int("status", res.Status),
			zap.Duration("duration", duration),
			zap.String("remote_ip", c.RealIP()),
		}

		// Add request ID if available
		if requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		// Add user info if available (from auth middleware)
		if userID, ok := c.Get("user_id").(string); ok && userID != "" {
			fields = append(fields, zap.String("user_id", userID))
		}

		// Log based on status code
		if res.Status >= 500 {
			m.logger.Error("request completed with error", fields...)
		} else if res.Status >= 400 {
			m.logger.Warn("request completed with client error", fields...)
		} else {
			m.logger.Info("request completed", fields...)
		}

		return err
	}
}
