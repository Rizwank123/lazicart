package proxy

import (
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

func NewProxy(target string, prefix string, serviceBase string) echo.HandlerFunc {
	targetURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

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
