package utils

import (
	"time"

	"github.com/Rizwank123/LahziCart/auth-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (a *authManager) GenerateToken(userId string, role string) (*models.LoginResponse, error) {
	jti := uuid.NewString()

	claims := jwt.MapClaims{
		"sub":  userId,
		"role": role,
		"jti":  jti,
		"exp":  time.Now().Add(time.Minute * time.Duration(a.config.ExpireIn)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(a.config.SecretKey))
	if err != nil {
		return nil, err
	}
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(a.config.RefreshKey))
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Jti:          jti,
	}, nil
}
