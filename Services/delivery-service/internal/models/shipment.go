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

type Shipment struct {
	Base
	OrderID             uuid.UUID  `json:"order_id" db:"order_id"`
	UserID              uuid.UUID  `json:"user_id" db:"user_id"`
	Carrier             *string    `json:"carrier,omitempty" db:"carrier"`
	TrackingNumber      *string    `json:"tracking_number,omitempty" db:"tracking_number"`
	Status              string     `json:"status" db:"status"`
	ShippingAddress     Address    `json:"shipping_address" db:"shipping_address"`
	EstimatedDeliveryAt *time.Time `json:"estimated_delivery_at,omitempty" db:"estimated_delivery_at"`
	DeliveredAt         *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	BaseAudit
}

type CreateShipmentInput struct {
	OrderID             uuid.UUID  `json:"order_id"`
	UserID              uuid.UUID  `json:"user_id"`
	Carrier             *string    `json:"carrier,omitempty"`
	TrackingNumber      *string    `json:"tracking_number,omitempty"`
	ShippingAddress     Address    `json:"shipping_address"`
	EstimatedDeliveryAt *time.Time `json:"estimated_delivery_at,omitempty"`
}

type UpdateShipmentStatusInput struct {
	Status         string     `json:"status"`
	Carrier        *string    `json:"carrier,omitempty"`
	TrackingNumber *string    `json:"tracking_number,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
}

type ShipmentRepository interface {
	Create(ctx context.Context, shipment *Shipment) error
	Get(ctx context.Context, id uuid.UUID) (*Shipment, error)
	GetByOrder(ctx context.Context, orderID uuid.UUID) (*Shipment, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]Shipment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, in UpdateShipmentStatusInput) error
}

type ShipmentService interface {
	Create(in CreateShipmentInput) (*Shipment, error)
	Get(id uuid.UUID) (*Shipment, error)
	GetByOrder(orderID uuid.UUID) (*Shipment, error)
	ListByUser(userID uuid.UUID, limit, page int) ([]Shipment, error)
	UpdateStatus(id uuid.UUID, in UpdateShipmentStatusInput) (*Shipment, error)
}
