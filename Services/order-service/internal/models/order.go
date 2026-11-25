package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
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

type Address map[string]interface{}

func (a Address) Value() (driver.Value, error) { return json.Marshal(a) }
func (a *Address) Scan(value interface{}) error {
	if value == nil {
		*a = Address{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, a)
}

type Order struct {
	Base
	UserID          uuid.UUID   `json:"user_id" db:"user_id"`
	Status          string      `json:"status" db:"status"`
	Total           float64     `json:"total" db:"total"`
	Currency        string      `json:"currency" db:"currency"`
	ShippingAddress Address     `json:"shipping_address" db:"shipping_address"`
	Items           []OrderItem `json:"items,omitempty"`
	BaseAudit
}

type OrderItem struct {
	Base
	OrderID   uuid.UUID `json:"order_id" db:"order_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	UnitPrice float64   `json:"unit_price" db:"unit_price"`
	Currency  string    `json:"currency" db:"currency"`
	BaseAudit
}

type OutboxEvent struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	EventID       uuid.UUID  `json:"event_id" db:"event_id"`
	AggregateType string     `json:"aggregate_type" db:"aggregate_type"`
	AggregateID   uuid.UUID  `json:"aggregate_id" db:"aggregate_id"`
	EventType     string     `json:"event_type" db:"event_type"`
	EventVersion  int        `json:"event_version" db:"event_version"`
	Payload       []byte     `json:"payload" db:"payload"`
	Status        string     `json:"status" db:"status"`
	Attempts      int        `json:"attempts" db:"attempts"`
	AvailableAt   time.Time  `json:"available_at" db:"available_at"`
	PublishedAt   *time.Time `json:"published_at,omitempty" db:"published_at"`
	LastError     *string    `json:"last_error,omitempty" db:"last_error"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

type CreateOrderInput struct {
	UserID          uuid.UUID         `json:"user_id"`
	Currency        string            `json:"currency"`
	ShippingAddress Address           `json:"shipping_address"`
	Items           []CreateOrderItem `json:"items"`
	CorrelationID   *uuid.UUID        `json:"-"`
	SagaID          *uuid.UUID        `json:"-"`
}

var ErrIdempotencyKeyConflict = errors.New("idempotency key was already used with a different request")

type CreateOrderItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order, event *OutboxEvent, idempotencyKey, requestHash string) (*Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*Order, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, event *OutboxEvent) error
	GetPendingOutboxEvents(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkOutboxPublished(ctx context.Context, id uuid.UUID) error
	MarkOutboxFailed(ctx context.Context, id uuid.UUID, cause string) error
}

type OrderService interface {
	CreateOrder(in CreateOrderInput, idempotencyKey string) (*Order, error)
	GetOrder(id uuid.UUID) (*Order, error)
	ListByUser(userID uuid.UUID, limit, page int) ([]Order, error)
	CancelOrder(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, status string) error
}
