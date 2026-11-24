package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type (

	// Base struct for all models
	Base struct {
		ID uuid.UUID `json:"id" db:"id"`
	} //	@name	Base

	BaseAudit struct {
		CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt time.Time  `json:"updated_at,omitempty" db:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	}

	SortKey struct {
		Key string `json:"key"`
		Val string `json:"val"`
	}
)
