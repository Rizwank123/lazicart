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

type Cart struct {
	Base
	UserID uuid.UUID  `json:"user_id" db:"user_id"`
	Status string     `json:"status" db:"status"`
	Items  []CartItem `json:"items,omitempty"`
	BaseAudit
}

type CartItem struct {
	Base
	CartID    uuid.UUID `json:"cart_id" db:"cart_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	UnitPrice float64   `json:"unit_price" db:"unit_price"`
	Currency  string    `json:"currency" db:"currency"`
	BaseAudit
}

type AddCartItemInput struct {
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	Currency  string    `json:"currency"`
}

type UpdateCartItemInput struct {
	Quantity int `json:"quantity"`
}

type CartRepository interface {
	GetOrCreateCart(ctx context.Context, userID uuid.UUID) (*Cart, error)
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*Cart, error)
	AddItem(ctx context.Context, cartID uuid.UUID, item *CartItem) error
	UpdateItem(ctx context.Context, itemID uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, itemID uuid.UUID) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
	ListItems(ctx context.Context, cartID uuid.UUID) ([]CartItem, error)
	// GetItemOwnerUserID resolves the user that owns a cart item.
	GetItemOwnerUserID(ctx context.Context, itemID uuid.UUID) (uuid.UUID, error)
}

type CartService interface {
	GetCart(userID uuid.UUID) (*Cart, error)
	AddItem(in AddCartItemInput) (*Cart, error)
	UpdateItem(itemID uuid.UUID, in UpdateCartItemInput) error
	RemoveItem(itemID uuid.UUID) error
	ClearCart(userID uuid.UUID) error
	// GetItemOwner resolves the user that owns a cart item.
	GetItemOwner(itemID uuid.UUID) (uuid.UUID, error)
}
