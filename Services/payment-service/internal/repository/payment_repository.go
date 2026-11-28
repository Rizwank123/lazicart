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

func (r *pgxPaymentRepository) Create(ctx context.Context, payment *models.Payment, idempotencyKey, requestHash string) (*models.Payment, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if idempotencyKey == "" {
		q := `INSERT INTO payments (order_id, user_id, amount, currency, provider, status)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at`
		err := r.db.QueryRow(ctx, q, payment.OrderID, payment.UserID, payment.Amount, payment.Currency, payment.Provider, payment.Status).
			Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
		return payment, err
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := `INSERT INTO payments (order_id, user_id, amount, currency, provider, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	if err := tx.QueryRow(ctx, q, payment.OrderID, payment.UserID, payment.Amount, payment.Currency, payment.Provider, payment.Status).
		Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt); err != nil {
		return nil, err
	}

	var insertedKey string
	err = tx.QueryRow(ctx, `
		INSERT INTO payment_idempotency_keys (idempotency_key, request_hash, payment_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (idempotency_key) DO NOTHING
		RETURNING idempotency_key`, idempotencyKey, requestHash, payment.ID).Scan(&insertedKey)
	if err == pgx.ErrNoRows {
		// The key was already used. The same request returns the original
		// payment, a different request is a conflict.
		var existingHash string
		var existingPaymentID uuid.UUID
		if lookupErr := tx.QueryRow(ctx, `
			SELECT request_hash, payment_id
			FROM payment_idempotency_keys
			WHERE idempotency_key = $1`, idempotencyKey).Scan(&existingHash, &existingPaymentID); lookupErr != nil {
			return nil, lookupErr
		}
		if existingHash != requestHash {
			return nil, models.ErrIdempotencyKeyConflict
		}
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != pgx.ErrTxClosed {
			return nil, rollbackErr
		}
		return r.Get(ctx, existingPaymentID)
	}
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return payment, nil
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
