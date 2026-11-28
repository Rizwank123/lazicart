package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Rizwank123/LahziCart/shared/identity"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// HeaderUserID carries the authenticated caller identity from the gateway to
// the downstream services.
//
// The gateway always removes any value supplied by the client before setting
// it, so services can trust it as long as they are only reachable through the
// gateway.
const HeaderUserID = identity.HeaderUserID

// IdentityMiddleware injects the authenticated caller identity and rejects
// requests that address another user's data.
//
// The gateway is the only place where the JWT is verified, so it is also the
// only place where the caller identity is known with certainty. Downstream
// services key their rows by the same user id that is carried in the JWT
// subject claim, which makes it possible to reject cross-user access here.
func IdentityMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		caller, ok := callerID(c)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing caller identity")
		}

		req := c.Request()
		// Never trust an identity supplied by the client.
		req.Header.Del(HeaderUserID)
		req.Header.Set(HeaderUserID, caller)

		if requested, ok := requestedUserID(c); ok && !strings.EqualFold(requested, caller) {
			return echo.NewHTTPError(http.StatusForbidden, "you cannot access another user's data")
		}

		if err := enforceBodyUserID(c, caller); err != nil {
			return err
		}

		return next(c)
	}
}

// callerID returns the user id carried in the JWT subject claim.
func callerID(c echo.Context) (string, bool) {
	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return "", false
	}
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", false
	}
	return sub, true
}

// requestedUserID extracts the user id addressed by the request path.
//
// The gateway registers its routes as wildcards, so the path has to be
// inspected directly instead of relying on the named Echo parameters.
func requestedUserID(c echo.Context) (string, bool) {
	segments := splitPath(c.Request().URL.Path)
	switch {
	case len(segments) == 2 && segments[0] == "users":
		// /api/users/<id>. "email" and "phone" are not ids.
		return asUserID(segments[1])
	case len(segments) == 2 && segments[0] == "carts":
		// /api/carts/<user_id>. "items" is a sub resource.
		return asUserID(segments[1])
	case len(segments) == 3 && segments[0] == "orders" && segments[1] == "users":
		return asUserID(segments[2])
	case len(segments) == 3 && segments[0] == "notifications" && segments[1] == "users":
		return asUserID(segments[2])
	case len(segments) == 3 && segments[0] == "deliveries" && segments[1] == "users":
		return asUserID(segments[2])
	}
	return "", false
}

func asUserID(segment string) (string, bool) {
	if _, err := uuid.FromString(segment); err != nil {
		return "", false
	}
	return segment, true
}

func splitPath(path string) []string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/api/"), "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

// enforceBodyUserID rejects a request whose body addresses another user.
func enforceBodyUserID(c echo.Context, caller string) error {
	req := c.Request()
	switch req.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
	default:
		return nil
	}
	if req.Body == nil || !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		return nil
	}

	body, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	// Restore the body so the downstream service can still read it.
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))

	if len(bytes.TrimSpace(body)) == 0 {
		return nil
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		// Not a JSON object; leave validation to the service.
		return nil
	}

	raw, ok := payload["user_id"]
	if !ok {
		return nil
	}
	requested, ok := raw.(string)
	if !ok || requested == "" {
		return nil
	}
	if !strings.EqualFold(requested, caller) {
		return echo.NewHTTPError(http.StatusForbidden, "you cannot act on another user's data")
	}
	return nil
}
