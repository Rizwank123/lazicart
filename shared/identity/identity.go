// Package identity carries the authenticated caller identity from the API
// gateway to the downstream services.
//
// The gateway verifies the JWT and removes any identity header supplied by the
// client before setting its own, so services can rely on the header as long as
// they are only reachable through the gateway.
package identity

import (
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

// HeaderUserID carries the authenticated caller identity.
const HeaderUserID = "X-User-Id"

// CallerID returns the authenticated caller id, or an empty string when the
// request did not pass through the gateway.
func CallerID(c echo.Context) string {
	if c == nil || c.Request() == nil {
		return ""
	}
	return strings.TrimSpace(c.Request().Header.Get(HeaderUserID))
}

// Owns reports whether the caller owns the given user id.
//
// An empty caller id means the request was not routed through the gateway, for
// example an internal call or a service reached directly. Such a request is
// not rejected here because the gateway is the component responsible for
// authentication.
func Owns(c echo.Context, userID string) bool {
	caller := CallerID(c)
	if caller == "" {
		return true
	}
	return strings.EqualFold(caller, userID)
}

// RequireOwnership rejects the request when the caller does not own the
// resource. It returns nil when the request did not come through the gateway.
func RequireOwnership(c echo.Context, userID string) error {
	if Owns(c, userID) {
		return nil
	}
	return echo.NewHTTPError(http.StatusForbidden, "you cannot access another user's data")
}

// RequireOwnershipUUID is [RequireOwnership] for a resource owner that is
// stored as a UUID.
func RequireOwnershipUUID(c echo.Context, userID uuid.UUID) error {
	return RequireOwnership(c, userID.String())
}

// RequireOwnershipOptional is [RequireOwnership] for a resource owner that may
// be absent, for example a notification that is not tied to a user.
//
// A resource without an owner is not owned by anybody, so it is only
// reachable by callers that are not identified at all.
func RequireOwnershipOptional(c echo.Context, userID *uuid.UUID) error {
	if userID == nil {
		if CallerID(c) == "" {
			return nil
		}
		return echo.NewHTTPError(http.StatusForbidden, "you cannot access another user's data")
	}
	return RequireOwnershipUUID(c, *userID)
}
