package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Rizwank123/LahziCart/shared/cache"
	"github.com/Rizwank123/LahziCart/shared/session"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// sessionLookupTimeout bounds the time that a request waits for the session
// state. A slow Redis must not slow down every request.
const sessionLookupTimeout = 2 * time.Second

// SessionStore reads the server side session state.
//
// It is the part of the cache client that the gateway needs, kept as an
// interface so the middleware can be tested without a Redis instance.
type SessionStore interface {
	Get(ctx context.Context, key string) (string, error)
}

// RevocationMiddleware rejects access tokens that the auth-service has
// invalidated.
//
// The auth-service keeps one Redis entry per issued token pair and removes it
// when the session is refreshed or the user signs out. A JWT on its own stays
// valid until its expiry claim passes, so the signature cannot tell a live
// token apart from a revoked one. Reading the session state is what makes a
// sign out effective right away, and what stops a token that was already
// replaced by a refresh from being used a second time.
//
// The answer is a 401 so that the mobile client treats it the same way as an
// expired token: it refreshes the session and retries the request.
func RevocationMiddleware(store SessionStore, logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("user").(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token claims")
			}
			jti, ok := claims["jti"].(string)
			if !ok || jti == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token identifier")
			}

			if store == nil {
				// The session state cannot be read at all. The token
				// signature was already verified, so the request is allowed
				// and the problem is reported instead of taking the whole
				// API down.
				logger.Warn("session state is unavailable, revocation is not enforced", zap.String("jti", jti))
				return next(c)
			}

			ctx, cancel := context.WithTimeout(c.Request().Context(), sessionLookupTimeout)
			defer cancel()

			_, err := store.Get(ctx, session.AccessKey(jti))
			switch {
			case err == nil:
				return next(c)
			case errors.Is(err, cache.ErrKeyNotFound):
				// The session was revoked, it was already rotated or it
				// expired. In all three cases the token is no longer usable.
				logger.Info("access token rejected", zap.String("jti", jti), zap.String("reason", "revoked"))
				return echo.NewHTTPError(http.StatusUnauthorized, "token has been revoked")
			default:
				// The session state could not be read. The token signature
				// was already verified, so the request is allowed and the
				// problem is reported instead of taking the whole API down.
				logger.Warn("could not read the session state", zap.String("jti", jti), zap.Error(err))
				return next(c)
			}
		}
	}
}
