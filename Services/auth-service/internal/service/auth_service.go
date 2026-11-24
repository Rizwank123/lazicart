package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Rizwank123/LahziCart/auth-service/internal/cache"
	"github.com/Rizwank123/LahziCart/auth-service/internal/config"
	"github.com/Rizwank123/LahziCart/auth-service/internal/models"
	"github.com/Rizwank123/LahziCart/auth-service/internal/utils"
	"github.com/Rizwank123/LahziCart/shared/events"
	kafkabroker "github.com/Rizwank123/LahziCart/shared/messaging/kafkaBroker"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

type AuthServiceImpl struct {
	ar            models.AuthRepository
	am            utils.AuthManager
	rdc           cache.RedisClient
	kafkaProducer *kafkabroker.Producer
	cfg           *config.AuthConfig
}

func NewAuthServiceImpl(ar models.AuthRepository, am utils.AuthManager, rdc cache.RedisClient, cfg *config.AuthConfig, kafkaProducer *kafkabroker.Producer) models.UserAuthService {
	return &AuthServiceImpl{ar: ar, am: am, rdc: rdc, cfg: cfg, kafkaProducer: kafkaProducer}

}

// LoginUser implements [models.UserAuthService].
func (a *AuthServiceImpl) LoginUser(userAuth *models.LoginUser) (*models.LoginResponse, error) {
	res, err := a.ar.GetUserAuthByEmail(context.Background(), userAuth.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := a.am.ComparePassword(userAuth.Password, res.Password); err != nil {
		return nil, err
	}
	token, err := a.am.GenerateToken(res.ID.String(), "user")
	if err != nil {
		return nil, err
	}
	a.rdc.Set(
		context.Background(),
		fmt.Sprintf("access:%s", token.Jti),
		token.RefreshToken,
		time.Minute*time.Duration(a.cfg.ExpireIn),
	)
	a.rdc.Set(
		context.Background(),
		fmt.Sprintf("refresh:%s", token.Jti),
		token.RefreshToken,
		time.Minute*time.Duration(a.cfg.ExpireIn),
	)
	return token, nil
}

// LoginWithPhone implements [models.UserAuthService].
func (a *AuthServiceImpl) LoginWithPhone(userAuth *models.OtpLogin) (*models.LoginResponse, error) {
	otp, err := a.rdc.Get(context.Background(), fmt.Sprintf("otp:%s", userAuth.Phone))
	if err != nil {
		return nil, err
	}
	if otp != userAuth.Otp {
		return nil, errors.New("invalid otp")
	}
	token, err := a.am.GenerateToken(userAuth.Phone, "user")
	if err != nil {
		return nil, err
	}
	a.rdc.Set(
		context.Background(),
		fmt.Sprintf("access:%s", token.Jti),
		token.RefreshToken,
		time.Minute*time.Duration(a.cfg.ExpireIn),
	)
	a.rdc.Set(
		context.Background(),
		fmt.Sprintf("refresh:%s", token.Jti),
		token.RefreshToken,
		time.Minute*time.Duration(a.cfg.ExpireIn),
	)
	return token, nil
}

// Logout implements [models.UserAuthService].
func (a *AuthServiceImpl) Logout(jti string) error {
	ctx := context.Background()

	err := a.rdc.Delete(ctx, fmt.Sprintf("refresh:%s", jti))
	if err != nil {
		return err
	}
	err = a.rdc.Delete(ctx, fmt.Sprintf("access:%s", jti))
	if err != nil {
		return err
	}
	return nil
}

// RefreshToken implements [models.UserAuthService].
func (a *AuthServiceImpl) RefreshToken(refreshToken string) (*models.LoginResponse, error) {
	ctx := context.Background()

	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}
	jwtToken, err := a.am.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}
	// 2. Extract Claims as MapClaims
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	// 3. Extract specific fields with type assertion
	// In your GenerateToken, "sub" holds the userId
	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("userID not found in token")
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("jti not found in token")
	}

	token, err := a.am.GenerateToken(userID, "user")
	if err != nil {
		return nil, err
	}
	// Delete refresh token from cache
	err = a.rdc.Delete(ctx, fmt.Sprintf("refresh:%s", jti))
	if err != nil {
		return nil, err
	}
	// Update access token in cache
	accessTTl := time.Minute * time.Duration(a.cfg.ExpireIn)
	err = a.rdc.Set(ctx, fmt.Sprintf("access:%s", token.Jti), token.RefreshToken, accessTTl)
	if err != nil {
		return nil, err
	}
	// Update refresh token in cache
	refreshTTL := time.Hour * 24 * 7
	err = a.rdc.Set(ctx, fmt.Sprintf("refresh:%s", token.Jti), token.RefreshToken, refreshTTL)
	if err != nil {
		return nil, err
	}
	return token, nil

}

// RegisterUser implements [models.UserAuthService].
func (a *AuthServiceImpl) RegisterUser(userAuth models.RegisterInput) (models.UserAuth, error) {
	res := models.UserAuth{
		Email:    userAuth.Email,
		Phone:    userAuth.Phone,
		Password: userAuth.Password,
		IsActive: true,
	}
	hashedPassword, err := a.am.HashPassword(userAuth.Password)
	if err != nil {
		return res, err
	}
	res.Password = hashedPassword
	err = a.ar.CreateUserAuth(context.Background(), &res)
	if err != nil {
		return res, err
	}
	userCreatedEvent := events.UserCreatedEvent{
		EventType: events.UserCreated,
		Version:   1,
		UserID:    res.ID.String(),
		Phone:     userAuth.Phone,
		Email:     userAuth.Email,
		FirstName: userAuth.FirstName,
		LastName:  userAuth.LastName,
		Address: events.Address{
			City:        userAuth.Address.City,
			State:       userAuth.Address.State,
			Country:     userAuth.Address.Country,
			Pincode:     userAuth.Address.Pincode,
			Coordinates: userAuth.Address.Coordinates,
		},
		Dob:       userAuth.Dob,
		Gender:    userAuth.Gender,
		AuthId:    res.ID,
		CreatedAt: time.Now(),
	}
	data, err := json.Marshal(userCreatedEvent)
	if err != nil {
		return res, err
	}
	a.kafkaProducer.Publish(context.Background(), events.UserCreated, data)
	return res, nil
}

// SendOtp implements [models.UserAuthService].
func (a *AuthServiceImpl) SendOtp(phone string) error {
	authUser, err := a.ar.GetUserAuthByPhone(context.Background(), phone)
	if err != nil {
		return err
	}
	if authUser == nil {
		return errors.New("user not found")
	}
	otp := a.GenerateOtp()

	key := fmt.Sprintf("otp:%s", phone)
	err = a.rdc.Set(context.Background(), key, otp, 5*time.Minute)
	if err != nil {
		return err
	}
	log.Println("OTP is ", otp)
	log.Println("OTP sent successfully")
	return nil
}
func (a *AuthServiceImpl) DeleteUserAuth(id uuid.UUID) error {
	return a.ar.DeleteUserAuth(context.Background(), id)
}
func (a *AuthServiceImpl) UpdateUserAuth(id uuid.UUID, in events.UserUpdatedEvent) error {
	userAuth, err := a.ar.GetUserAuthById(context.Background(), id)
	if err != nil {
		return err
	}
	if in.Email != "" {
		userAuth.Email = in.Email
	}
	if in.Phone != "" {
		userAuth.Phone = in.Phone
	}
	return a.ar.UpdateUserAuth(context.Background(), userAuth)
}

func (a *AuthServiceImpl) ValidateToken(token string) (bool, string, string, string, error) {
	jwtToken, err := a.am.ValidateToken(token)
	if err != nil {
		return false, "", "", "", err
	}

	// Check token validity
	if !jwtToken.Valid {
		return false, "", "", "", fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, "", "", "", fmt.Errorf("invalid claims")
	}

	userID, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)
	jti, _ := claims["jti"].(string)

	return true, userID, jti, role, nil
}

func (a *AuthServiceImpl) GenerateOtp() string {
	return fmt.Sprintf("%06d", rand.Intn(900000)+100000)
}
