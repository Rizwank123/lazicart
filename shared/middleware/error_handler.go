package middleware

import (
	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ErrorHandler provides centralized error handling with consistent responses
type ErrorHandler struct {
	logger             *zap.Logger
	showInternalErrors bool
}

// NewErrorHandler creates a new error handler middleware
// showInternalErrors: if true, includes detailed error messages (for dev/staging)
func NewErrorHandler(logger *zap.Logger, showInternalErrors bool) *ErrorHandler {
	return &ErrorHandler{
		logger:             logger,
		showInternalErrors: showInternalErrors,
	}
}

// Handle is the middleware handler for centralized error handling
func (e *ErrorHandler) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		// Get request ID for correlation
		requestID, _ := c.Get("request_id").(string)

		// Log the error
		e.logger.Error("request error",
			zap.String("request_id", requestID),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().URL.Path),
			zap.Error(err),
		)

		// Check if it's an Echo HTTPError
		if he, ok := err.(*echo.HTTPError); ok {
			// Echo already handled the status code
			message := "an error occurred"
			if msg, ok := he.Message.(string); ok {
				message = msg
			}
			return c.JSON(he.Code, response.Error(message))
		}

		// For other errors, return internal server error
		message := "internal server error"
		if e.showInternalErrors {
			message = err.Error()
		}

		return c.JSON(500, response.Error(message))
	}
}

// CustomErrorHandler returns an Echo error handler function
// This can be set on the Echo instance: e.HTTPErrorHandler = ErrorHandler.CustomErrorHandler()
func (e *ErrorHandler) CustomErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		// Get request ID for correlation
		requestID, _ := c.Get("request_id").(string)

		code := 500
		message := "internal server error"

		// Check if it's an Echo HTTPError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if msg, ok := he.Message.(string); ok {
				message = msg
			} else {
				message = "an error occurred"
			}
		} else if e.showInternalErrors {
			message = err.Error()
		}

		// Log the error
		e.logger.Error("http error",
			zap.String("request_id", requestID),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().URL.Path),
			zap.Int("status", code),
			zap.Error(err),
		)

		// Send response
		if !c.Response().Committed {
			if err := c.JSON(code, response.Error(message)); err != nil {
				e.logger.Error("failed to send error response", zap.Error(err))
			}
		}
	}
}
