package models

import (
	"context"
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

type InventoryItem struct {
	Base
	ProductID         uuid.UUID `json:"product_id" db:"product_id"`
	SKU               *string   `json:"sku,omitempty" db:"sku"`
	QuantityAvailable int       `json:"quantity_available" db:"quantity_available"`
	QuantityReserved  int       `json:"quantity_reserved" db:"quantity_reserved"`
	ReorderLevel      int       `json:"reorder_level" db:"reorder_level"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	BaseAudit
}

type InventoryReservation struct {
	Base
	OrderID   uuid.UUID `json:"order_id" db:"order_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Status    string    `json:"status" db:"status"`
	BaseAudit
}

type UpsertInventoryInput struct {
	ProductID    uuid.UUID `json:"product_id"`
	SKU          *string   `json:"sku,omitempty"`
	Quantity     int       `json:"quantity"`
	ReorderLevel int       `json:"reorder_level"`
	IsActive     bool      `json:"is_active"`
}

type AdjustInventoryInput struct {
	ProductID uuid.UUID `json:"product_id"`
	Delta     int       `json:"delta"`
}

type ReserveInventoryInput struct {
	OrderID uuid.UUID              `json:"order_id"`
	Items   []ReserveInventoryItem `json:"items"`
}

type ReserveInventoryItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type InventoryRepository interface {
	UpsertItem(ctx context.Context, item *InventoryItem) error
	GetByProductID(ctx context.Context, productID uuid.UUID) (*InventoryItem, error)
	ListItems(ctx context.Context, limit, offset int64) (int64, []InventoryItem, error)
	Adjust(ctx context.Context, productID uuid.UUID, delta int) (*InventoryItem, error)
	Reserve(ctx context.Context, orderID uuid.UUID, items []ReserveInventoryItem) ([]InventoryReservation, error)
	Release(ctx context.Context, orderID uuid.UUID) error
	Commit(ctx context.Context, orderID uuid.UUID) error
}

type InventoryService interface {
	UpsertItem(in UpsertInventoryInput) (*InventoryItem, error)
	GetByProductID(productID uuid.UUID) (*InventoryItem, error)
	ListItems(limit, page int) ([]InventoryItem, int64, error)
	Adjust(in AdjustInventoryInput) (*InventoryItem, error)
	Reserve(in ReserveInventoryInput) ([]InventoryReservation, error)
	Release(orderID uuid.UUID) error
	Commit(orderID uuid.UUID) error
}
