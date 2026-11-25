package middleware

import (
	gofrsuuid "github.com/gofrs/uuid/v5"
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

		correlationID := parseOrNewID(req.Header.Get("X-Correlation-ID"))
		c.Set("correlation_id", correlationID)
		res.Header().Set("X-Correlation-ID", correlationID.String())

		if sagaHeader := req.Header.Get("X-Saga-ID"); sagaHeader != "" {
			if sagaID, err := gofrsuuid.FromString(sagaHeader); err == nil {
				c.Set("saga_id", sagaID)
				res.Header().Set("X-Saga-ID", sagaID.String())
			}
		}

		return next(c)
	}
}

func parseOrNewID(value string) gofrsuuid.UUID {
	if id, err := gofrsuuid.FromString(value); err == nil {
		return id
	}
	id, err := gofrsuuid.NewV4()
	if err != nil {
		return gofrsuuid.Nil
	}
	return id
}
