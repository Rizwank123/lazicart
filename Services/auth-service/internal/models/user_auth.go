package models

import (
	"context"
	"time"

	"github.com/Rizwank123/LahziCart/shared/events"
	"github.com/gofrs/uuid/v5"
)

type (
	// UserAuth defines the structure of a user authentication record
	UserAuth struct {
		Base
		Email    string `json:"email" db:"email"`
		Phone    string `json:"phone" db:"phone"`
		Password string `json:"password" db:"password"`
		IsActive bool   `json:"is_active" db:"is_active"`
		BaseAudit
	}
	// RegisterInput defines the model for register new auth
	RegisterInput struct {
		Email     string    `json:"email"`
		Phone     string    `json:"phone"`
		Password  string    `json:"password"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Address   Address   `json:"address"`
		Dob       time.Time `json:"dob" format:"date"`
		Gender    string    `json:"gender"`
	}
	// LoginUser defines the structure of a user login record
	LoginUser struct {
		Email    string `json:"email" db:"email" example:"test@gmail.com"`
		Password string `json:"password" db:"password" example:"password"`
	} // @name LoginUser
	LoginResponse struct {
		AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
		RefreshToken string `json:"refresh_token" example:""`
		Jti          string `json:"jti" example:"1234567890"`
	} // @name LoginResponse
	OtpLogin struct {
		Phone string `json:"phone" example:"1234567890"`
		Otp   string `json:"otp" example:"123456"`
	} // @name OtpLogin

)

// AuthRepository defines the interface for user authentication repository
type AuthRepository interface {
	CreateUserAuth(ctx context.Context, userAuth *UserAuth) error
	GetUserAuthByEmail(ctx context.Context, email string) (*UserAuth, error)
	GetUserAuthByPhone(ctx context.Context, phone string) (*UserAuth, error)
	UpdateUserAuth(ctx context.Context, userAuth *UserAuth) error
	DeleteUserAuth(ctx context.Context, id uuid.UUID) error
	GetUserAuthById(ctx context.Context, id uuid.UUID) (*UserAuth, error)
}

// userAuthService defines the service for user authentication
type UserAuthService interface {
	RegisterUser(userAuth RegisterInput) (UserAuth, error)
	LoginUser(userAuth *LoginUser) (*LoginResponse, error)
	SendOtp(phone string) error
	LoginWithPhone(userAuth *OtpLogin) (*LoginResponse, error)
	RefreshToken(refreshToken string) (*LoginResponse, error)
	Logout(refreshToken string) error
	DeleteUserAuth(id uuid.UUID) error
	UpdateUserAuth(id uuid.UUID, in events.UserUpdatedEvent) error
	ValidateToken(token string) (bool, string, string, string, error)
}
