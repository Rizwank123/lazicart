package events

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type UserCreatedEvent struct {
	EventType string    `json:"event_type"`
	Version   int       `json:"version"`
	UserID    string    `json:"user_id"`
	Phone     string    `json:"phone,omitempty"`
	Email     string    `json:"email,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Address   Address   `json:"address,omitempty"`
	Dob       time.Time `json:"dob,omitempty" format:"date"`
	Gender    string    `json:"gender,omitempty"`
	AuthId    uuid.UUID `json:"auth_id"`
	CreatedAt time.Time `json:"created_at"`
}

type UserUpdatedEvent struct {
	EventType string    `json:"event_type"`
	Version   int       `json:"version"`
	Phone     string    `json:"phone,omitempty"`
	Email     string    `json:"email,omitempty"`
	AuthId    uuid.UUID `json:"auth_id"`
}

type UserDeletedEvent struct {
	EventType string    `json:"event_type"`
	Version   int       `json:"version"`
	UserID    string    `json:"user_id"`
	DeletedAt time.Time `json:"deleted_at"`
	AuthId    uuid.UUID `json:"auth_id"`
}

type Address struct {
	City        string    `json:"city,omitempty" db:"city" example:"Mumbai"`
	State       string    `json:"state,omitempty" db:"state" example:"Maharashtra"`
	Country     string    `json:"country,omitempty" db:"country" example:"India"`
	Pincode     string    `json:"pincode,omitempty" db:"pincode" example:"123456"`
	Coordinates []float64 `json:"coordinates,omitempty" db:"coordinates" example:"[12.9716, 77.5946]"`
}
