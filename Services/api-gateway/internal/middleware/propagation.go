package middleware

import (
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

// PropagationMiddleware establishes IDs at the system edge. Reverse proxies
// preserve these headers when forwarding requests to downstream services.
type PropagationMiddleware struct{}

func NewPropagationMiddleware() *PropagationMiddleware {
	return &PropagationMiddleware{}
}

func (m *PropagationMiddleware) Inject(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := validOrNew(c.Request().Header.Get("X-Request-ID"))
		correlationID := validOrNew(c.Request().Header.Get("X-Correlation-ID"))

		c.Request().Header.Set("X-Request-ID", requestID.String())
		c.Request().Header.Set("X-Correlation-ID", correlationID.String())
		c.Response().Header().Set("X-Request-ID", requestID.String())
		c.Response().Header().Set("X-Correlation-ID", correlationID.String())

		if sagaID, err := uuid.FromString(c.Request().Header.Get("X-Saga-ID")); err == nil {
			c.Request().Header.Set("X-Saga-ID", sagaID.String())
			c.Response().Header().Set("X-Saga-ID", sagaID.String())
		}

		c.Set("request_id", requestID)
		c.Set("correlation_id", correlationID)
		return next(c)
	}
}

func validOrNew(value string) uuid.UUID {
	if id, err := uuid.FromString(value); err == nil {
		return id
	}
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil
	}
	return id
}
