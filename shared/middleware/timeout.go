package middleware

import (
	"context"
	"time"

	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/labstack/echo/v4"
)

// TimeoutMiddleware enforces request timeout limits
type TimeoutMiddleware struct {
	timeout time.Duration
}

// NewTimeoutMiddleware creates a new timeout middleware
// timeout is the maximum duration for a request (e.g., 30*time.Second)
func NewTimeoutMiddleware(timeout time.Duration) *TimeoutMiddleware {
	if timeout <= 0 {
		timeout = 30 * time.Second // default to 30 seconds
	}
	return &TimeoutMiddleware{
		timeout: timeout,
	}
}

// Enforce is the middleware handler that enforces request timeouts
func (m *TimeoutMiddleware) Enforce(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request().Context(), m.timeout)
		defer cancel()

		// Replace the request context
		c.SetRequest(c.Request().WithContext(ctx))

		// Channel to capture handler response
		done := make(chan error, 1)

		// Execute handler in goroutine
		go func() {
			done <- next(c)
		}()

		// Wait for either completion or timeout
		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			// Timeout occurred
			if ctx.Err() == context.DeadlineExceeded {
				return c.JSON(408, response.Error("request timeout"))
			}
			return ctx.Err()
		}
	}
}
