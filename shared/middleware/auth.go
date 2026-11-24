package middleware

import (
	"strings"

	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AuthManager defines the interface for JWT operations
// Services should implement this interface with their own JWT logic
type AuthManager interface {
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// RedisTokenStore defines the interface for Redis token storage
type RedisTokenStore interface {
	Get(ctx echo.Context, key string) (string, error)
}

type AuthMiddleware struct {
	jwt    AuthManager
	redis  RedisTokenStore
	logger *zap.Logger
}

func NewAuthMiddleware(
	jwt AuthManager,
	redis RedisTokenStore,
	logger *zap.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		jwt:    jwt,
		redis:  redis,
		logger: logger,
	}
}

func (a *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "authorization header missing")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "invalid authorization format")
		}

		tokenStr := parts[1]

		// 1️⃣ Validate JWT
		token, err := a.jwt.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			a.logger.Warn("invalid jwt", zap.Error(err))
			return response.Unauthorized(c, "invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return response.Unauthorized(c, "invalid token claims")
		}

		userID, _ := claims["user_id"].(string)
		role, _ := claims["role"].(string)

		// 2️⃣ Redis validation (logout / revoke)
		redisKey := "access_token:" + tokenStr
		if _, err := a.redis.Get(c, redisKey); err != nil {
			a.logger.Warn("token revoked", zap.String("user_id", userID))
			return response.Unauthorized(c, "token revoked")
		}

		// 3️⃣ Inject context
		c.Set("user_id", userID)
		c.Set("role", role)

		return next(c)
	}
}
