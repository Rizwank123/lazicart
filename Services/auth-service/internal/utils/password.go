package utils

import (
	"fmt"

	"github.com/Rizwank123/LahziCart/auth-service/internal/config"
	"github.com/Rizwank123/LahziCart/auth-service/internal/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthManager interface {
	HashPassword(password string) (string, error)
	ComparePassword(password string, hash string) error
	GenerateToken(userId string, role string) (*models.LoginResponse, error)
	ValidateToken(token string) (*jwt.Token, error)
	ValidateRefreshToken(token string) (*jwt.Token, error)
}
type authManager struct {
	config *config.AuthConfig
}

// ValidateToken implements [AuthManager].
func (a *authManager) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.config.SecretKey), nil
	})
}

// ValidateRefreshToken implements [AuthManager].
//
// Refresh tokens are signed with RefreshKey while access tokens are signed
// with SecretKey, so they must be verified with their own key. Validating a
// refresh token with SecretKey always fails the signature check and breaks
// the entire refresh flow.
func (a *authManager) ValidateRefreshToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.config.RefreshKey), nil
	})
}

// ComparePassword implements [AuthManager].
func (a *authManager) ComparePassword(password string, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}

// HashPassword implements [AuthManager].
func (a *authManager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func NewAuthManager(config *config.AuthConfig) AuthManager {
	return &authManager{config: config}
}
