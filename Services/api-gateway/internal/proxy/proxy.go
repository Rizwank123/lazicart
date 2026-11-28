package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

// NewProxy builds a reverse proxy for a downstream service.
//
// target must be an absolute URL such as http://localhost:7702. An empty or
// malformed target previously produced an opaque 502 for every request, hiding
// the real cause (a missing service URL in the gateway .env). It is now
// rejected up front and reported as a 503 that names the affected route.
func NewProxy(target string, prefix string, serviceBase string) echo.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
		return func(c echo.Context) error {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"message": "service unavailable",
				"error":   "invalid upstream target configured for " + prefix,
			})
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Surface upstream connection failures as a structured JSON error instead
	// of Go's default plain-text 502 body.
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "service unavailable",
			"error":   err.Error(),
		})
	}

	return func(c echo.Context) error {
		req := c.Request()

		// Remove gateway prefix
		path := strings.TrimPrefix(req.URL.Path, prefix)

		// Add service base path
		req.URL.Path = serviceBase + path

		proxy.ServeHTTP(c.Response(), req)
		return nil
	}
}
