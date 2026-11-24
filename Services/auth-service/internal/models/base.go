package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

// Base struct for all models
type Base struct {
	ID uuid.UUID `json:"id" db:"id"`
} // @name Base

type BaseAudit struct {
	CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
} // @name BaseAudit

type Address struct {
	City        string    `json:"city,omitempty" db:"city" example:"Mumbai"`
	State       string    `json:"state,omitempty" db:"state" example:"Maharashtra"`
	Country     string    `json:"country,omitempty" db:"country" example:"India"`
	Pincode     string    `json:"pincode,omitempty" db:"pincode" example:"123456"`
	Coordinates []float64 `json:"coordinates,omitempty" db:"coordinates" example:"[12.9716, 77.5946]"`
}

func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Address) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
