package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	// AllowOrigins defines a list of origins that may access the resource
	// Default: ["*"]
	AllowOrigins []string

	// AllowMethods defines a list methods allowed when accessing the resource
	// Default: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"]
	AllowMethods []string

	// AllowHeaders defines a list of request headers that can be used
	// Default: ["Origin", "Content-Type", "Accept", "Authorization"]
	AllowHeaders []string

	// ExposeHeaders defines whitelist headers that clients are allowed to access
	// Default: []
	ExposeHeaders []string

	// AllowCredentials indicates whether the request can include user credentials
	// Default: false
	AllowCredentials bool

	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached
	// Default: 0
	MaxAge int
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{},
		AllowCredentials: false,
		MaxAge:           0,
	}
}

// NewCORSMiddleware creates a new CORS middleware with the given configuration
func NewCORSMiddleware(config CORSConfig) echo.MiddlewareFunc {
	// Set defaults if not provided
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = []string{"*"}
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}
	}
	if len(config.AllowHeaders) == 0 {
		config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}

	allowOrigins := strings.Join(config.AllowOrigins, ",")
	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			origin := req.Header.Get("Origin")

			// Check if origin is allowed
			allowOrigin := ""
			if len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*" {
				allowOrigin = "*"
			} else if origin != "" {
				for _, o := range config.AllowOrigins {
					if o == origin {
						allowOrigin = origin
						break
					}
				}
			}

			// Set CORS headers
			if allowOrigin != "" {
				res.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			} else if len(config.AllowOrigins) > 0 {
				res.Header().Set("Access-Control-Allow-Origin", allowOrigins)
			}

			if config.AllowCredentials {
				res.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight request
			if req.Method == http.MethodOptions {
				res.Header().Set("Access-Control-Allow-Methods", allowMethods)
				res.Header().Set("Access-Control-Allow-Headers", allowHeaders)
				if exposeHeaders != "" {
					res.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
				}
				if config.MaxAge > 0 {
					res.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
				}
				return c.NoContent(http.StatusNoContent)
			}

			// Set expose headers for actual request
			if exposeHeaders != "" {
				res.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
			}

			return next(c)
		}
	}
}
