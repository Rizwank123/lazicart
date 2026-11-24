package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RecoveryMiddleware recovers from panics and returns 500 errors gracefully
type RecoveryMiddleware struct {
	logger *zap.Logger
}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware(logger *zap.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		logger: logger,
	}
}

// Recover is the middleware handler that recovers from panics
func (m *RecoveryMiddleware) Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				// Get stack trace
				stack := debug.Stack()

				// Get request ID for correlation
				requestID, _ := c.Get("request_id").(string)

				// Log the panic with full details
				m.logger.Error("panic recovered",
					zap.String("request_id", requestID),
					zap.String("method", c.Request().Method),
					zap.String("path", c.Request().URL.Path),
					zap.String("remote_ip", c.RealIP()),
					zap.Any("panic", r),
					zap.String("stack", string(stack)),
				)

				// Return internal server error response
				err := c.JSON(500, response.Error("internal server error"))
				if err != nil {
					m.logger.Error("failed to send panic recovery response", zap.Error(err))
				}
			}
		}()

		return next(c)
	}
}

// RecoverWithCustomHandler allows custom panic handling
func (m *RecoveryMiddleware) RecoverWithCustomHandler(handler func(c echo.Context, err interface{})) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					// Get stack trace
					stack := debug.Stack()

					// Get request ID for correlation
					requestID, _ := c.Get("request_id").(string)

					// Log the panic
					m.logger.Error("panic recovered",
						zap.String("request_id", requestID),
						zap.String("method", c.Request().Method),
						zap.String("path", c.Request().URL.Path),
						zap.String("remote_ip", c.RealIP()),
						zap.Any("panic", r),
						zap.String("stack", string(stack)),
					)

					// Call custom handler
					handler(c, r)
				}
			}()

			return next(c)
		}
	}
}

// PanicError creates a formatted error from panic recovery
func PanicError(r interface{}) error {
	return fmt.Errorf("panic: %v", r)
}
