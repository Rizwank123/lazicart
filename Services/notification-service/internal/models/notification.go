package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Base struct {
	ID uuid.UUID `json:"id" db:"id"`
}
type BaseAudit struct {
	CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Metadata map[string]interface{}

func (m Metadata) Value() (driver.Value, error) { return json.Marshal(m) }
func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		*m = Metadata{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, m)
}

type Notification struct {
	Base
	UserID    *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Channel   string     `json:"channel" db:"channel"`
	Type      string     `json:"type" db:"type"`
	Recipient *string    `json:"recipient,omitempty" db:"recipient"`
	Subject   *string    `json:"subject,omitempty" db:"subject"`
	Body      string     `json:"body" db:"body"`
	Status    string     `json:"status" db:"status"`
	Metadata  Metadata   `json:"metadata" db:"metadata"`
	BaseAudit
}

type CreateNotificationInput struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Channel   string     `json:"channel"`
	Type      string     `json:"type"`
	Recipient *string    `json:"recipient,omitempty"`
	Subject   *string    `json:"subject,omitempty"`
	Body      string     `json:"body"`
	Metadata  Metadata   `json:"metadata"`
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]Notification, error)
	MarkSent(ctx context.Context, id uuid.UUID) error
	// Get returns a notification by id.
	Get(ctx context.Context, id uuid.UUID) (*Notification, error)
}

type NotificationService interface {
	Create(in CreateNotificationInput) (*Notification, error)
	ListByUser(userID uuid.UUID, limit, page int) ([]Notification, error)
	MarkSent(id uuid.UUID) error
	// Get returns a notification by id.
	Get(id uuid.UUID) (*Notification, error)
}
