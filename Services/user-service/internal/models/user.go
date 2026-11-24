package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	// User define the model for user
	User struct {
		Base
		FirstName string    `json:"first_name,omitempty" db:"first_name" example:"Mohammad"`
		LastName  string    `json:"last_name,omitempty" db:"last_name" example:"Rizwan"`
		Email     string    `json:"email,omitempty" db:"email" example:"rizwan@gmail.com"`
		AuthId    uuid.UUID `json:"auth_id" db:"auth_id"`
		Phone     string    `json:"phone,omitempty" db:"phone" example:"1234567890"`
		IsActive  bool      `json:"is_active,omitempty" db:"is_active" example:"true"`
		Dob       time.Time `json:"dob,omitempty" format:"date" db:"dob" example:"2000-01-01"`
		Gender    string    `json:"gender,omitempty" example:"male"`
		Address   Address   `json:"address,omitempty" db:"address"`
		BaseAudit
	} //	@name	User

	UpdateUser struct {
		FirstName string    `json:"first_name,omitempty" example:"Mohammad"`
		LastName  string    `json:"last_name,omitempty" example:"Rizwan"`
		Email     string    `json:"email,omitempty" example:"rizwan@gmail.com"`
		Phone     string    `json:"phone,omitempty" example:"1234567890"`
		IsActive  bool      `json:"is_active,omitempty" example:"true"`
		Dob       time.Time `json:"dob,omitempty" format:"date" example:"2000-01-01"`
		Address   Address   `json:"address,omitempty" db:"address"`
		Gender    string    `json:"gender,omitempty" example:"male"`
	} //	@name	UpdateUser

	CreateUserRequest struct {
		FirstName string    `json:"first_name,omitempty" example:"Mohammad"`
		LastName  string    `json:"last_name,omitempty" example:"Rizwan"`
		Email     string    `json:"email,omitempty" example:"rizwan@gmail.com"`
		Phone     string    `json:"phone,omitempty" example:"1234567890"`
		IsActive  bool      `json:"is_active,omitempty" example:"true"`
		AuthId    uuid.UUID `json:"auth_id" db:"auth_id"`
		Dob       time.Time `json:"dob,omitempty" format:"date" example:"2000-01-01"`
		Address   Address   `json:"address,omitempty" db:"address"`
		Gender    string    `json:"gender,omitempty" example:"male"`
	} //	@name	CreateUserRequest

)

type (
	// UserRepository define the interface for user repository
	UserRepository interface {
		CreateUser(ctx context.Context, user *User) error
		UpdateUser(ctx context.Context, user *User) error
		DeleteUser(ctx context.Context, user *User) error
		GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
		GetUserByEmail(ctx context.Context, email string) (*User, error)
		GetUserByPhone(ctx context.Context, phone string) (*User, error)
	}

	// UserService define the interface for user service
	UserService interface {
		CreateUser(in CreateUserRequest) (*User, error)
		UpdateUser(id uuid.UUID, in UpdateUser) (*User, error)
		DeleteUser(id uuid.UUID) (*User, error)
		GetUserByID(id uuid.UUID) (*User, error)
		GetUserByEmail(email string) (*User, error)
		GetUserByPhone(phone string) (*User, error)
	}
)
