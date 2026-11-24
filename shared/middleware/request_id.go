package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RequestIDMiddleware generates and injects unique request IDs for distributed tracing
type RequestIDMiddleware struct {
	headerName string
}

// NewRequestIDMiddleware creates a new request ID middleware
// headerName is the HTTP header name for the request ID (e.g., "X-Request-ID")
func NewRequestIDMiddleware(headerName string) *RequestIDMiddleware {
	if headerName == "" {
		headerName = "X-Request-ID"
	}
	return &RequestIDMiddleware{
		headerName: headerName,
	}
}

// Inject is the middleware handler that generates or propagates request IDs
func (m *RequestIDMiddleware) Inject(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		// Check if request ID already exists in header (from gateway/upstream)
		requestID := req.Header.Get(m.headerName)
		if requestID == "" {
			// Generate new UUID v4 if not present
			requestID = uuid.New().String()
		}

		// Inject into context for downstream use
		c.Set("request_id", requestID)

		// Add to response header for client tracking
		res.Header().Set(m.headerName, requestID)

		return next(c)
	}
}
