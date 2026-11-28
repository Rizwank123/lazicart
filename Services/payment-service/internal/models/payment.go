package models

import (
	"context"
	"errors"
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

type Payment struct {
	Base
	OrderID       uuid.UUID `json:"order_id" db:"order_id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Amount        float64   `json:"amount" db:"amount"`
	Currency      string    `json:"currency" db:"currency"`
	Provider      string    `json:"provider" db:"provider"`
	Status        string    `json:"status" db:"status"`
	TransactionID *string   `json:"transaction_id,omitempty" db:"transaction_id"`
	FailureReason *string   `json:"failure_reason,omitempty" db:"failure_reason"`
	BaseAudit
}

type CreatePaymentInput struct {
	OrderID  uuid.UUID `json:"order_id"`
	UserID   uuid.UUID `json:"user_id"`
	Amount   float64   `json:"amount"`
	Currency string    `json:"currency"`
	Provider string    `json:"provider"`
}

type UpdatePaymentStatusInput struct {
	Status        string  `json:"status"`
	TransactionID *string `json:"transaction_id,omitempty"`
	FailureReason *string `json:"failure_reason,omitempty"`
}

// ErrIdempotencyKeyConflict is returned when an idempotency key is reused
// with a different request payload.
var ErrIdempotencyKeyConflict = errors.New("idempotency key was already used with a different request")

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment, idempotencyKey, requestHash string) (*Payment, error)
	Get(ctx context.Context, id uuid.UUID) (*Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*Payment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, transactionID, failureReason *string) error
}

type PaymentService interface {
	CreatePayment(in CreatePaymentInput, idempotencyKey string) (*Payment, error)
	GetPayment(id uuid.UUID) (*Payment, error)
	GetByOrderID(orderID uuid.UUID) (*Payment, error)
	UpdateStatus(id uuid.UUID, in UpdatePaymentStatusInput) (*Payment, error)
}
