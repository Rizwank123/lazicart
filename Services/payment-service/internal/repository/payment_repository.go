package repository

import (
	"context"

	"github.com/Rizwank123/LahziCart/payment-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxPaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) models.PaymentRepository {
	return &pgxPaymentRepository{db: db}
}

func (r *pgxPaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `INSERT INTO payments (order_id, user_id, amount, currency, provider, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, q, payment.OrderID, payment.UserID, payment.Amount, payment.Currency, payment.Provider, payment.Status).
		Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
}

func (r *pgxPaymentRepository) Get(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `SELECT id, order_id, user_id, amount, currency, provider, status, transaction_id, failure_reason, created_at, updated_at, deleted_at FROM payments WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}
	payment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Payment])
	return &payment, err
}

func (r *pgxPaymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `SELECT id, order_id, user_id, amount, currency, provider, status, transaction_id, failure_reason, created_at, updated_at, deleted_at FROM payments WHERE order_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1`, orderID)
	if err != nil {
		return nil, err
	}
	payment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Payment])
	return &payment, err
}

func (r *pgxPaymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, transactionID, failureReason *string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `UPDATE payments SET status = $1, transaction_id = $2, failure_reason = $3, updated_at = NOW() WHERE id = $4 AND deleted_at IS NULL`, status, transactionID, failureReason, id)
	return err
}
