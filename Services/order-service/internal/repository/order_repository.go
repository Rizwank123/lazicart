package repository

import (
	"context"
	"time"

	"github.com/Rizwank123/LahziCart/order-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxOrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) models.OrderRepository {
	return &pgxOrderRepository{db: db}
}

func (r *pgxOrderRepository) CreateOrder(ctx context.Context, order *models.Order, event *models.OutboxEvent, idempotencyKey, requestHash string) (*models.Order, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := `INSERT INTO orders (id, user_id, status, total, currency, shipping_address)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	if err := tx.QueryRow(ctx, q, order.ID, order.UserID, order.Status, order.Total, order.Currency, order.ShippingAddress).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return nil, err
	}
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		err := tx.QueryRow(ctx, `INSERT INTO order_items (order_id, product_id, quantity, unit_price, currency)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at`,
			order.ID, order.Items[i].ProductID, order.Items[i].Quantity, order.Items[i].UnitPrice, order.Items[i].Currency).
			Scan(&order.Items[i].ID, &order.Items[i].CreatedAt, &order.Items[i].UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	if err := insertOutboxEvent(ctx, tx, event); err != nil {
		return nil, err
	}
	if idempotencyKey != "" {
		var insertedKey string
		err := tx.QueryRow(ctx, `
			INSERT INTO order_idempotency_keys (idempotency_key, request_hash, order_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (idempotency_key) DO NOTHING
			RETURNING idempotency_key`, idempotencyKey, requestHash, order.ID).Scan(&insertedKey)
		if err == pgx.ErrNoRows {
			var existingHash string
			var existingOrderID uuid.UUID
			if lookupErr := tx.QueryRow(ctx, `
				SELECT request_hash, order_id
				FROM order_idempotency_keys
				WHERE idempotency_key = $1`, idempotencyKey).Scan(&existingHash, &existingOrderID); lookupErr != nil {
				return nil, lookupErr
			}
			if existingHash != requestHash {
				return nil, models.ErrIdempotencyKeyConflict
			}
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != pgx.ErrTxClosed {
				return nil, rollbackErr
			}
			existingOrder, getErr := r.GetOrder(ctx, existingOrderID)
			return existingOrder, getErr
		}
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return order, nil
}

func (r *pgxOrderRepository) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `SELECT id, user_id, status, total, currency, shipping_address, created_at, updated_at, deleted_at FROM orders WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}
	order, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Order])
	if err != nil {
		return nil, err
	}
	order.Items, _ = r.listItems(ctx, order.ID)
	return &order, nil
}

func (r *pgxOrderRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]models.Order, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `SELECT id, user_id, status, total, currency, shipping_address, created_at, updated_at, deleted_at
		FROM orders WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Order])
	if err != nil {
		return nil, err
	}
	for i := range orders {
		orders[i].Items, _ = r.listItems(ctx, orders[i].ID)
	}
	return orders, nil
}

func (r *pgxOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, event *models.OutboxEvent) error {
	if ctx == nil {
		ctx = context.Background()
	}
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`, status, id); err != nil {
		return err
	}
	if err := insertOutboxEvent(ctx, tx, event); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgxOrderRepository) listItems(ctx context.Context, orderID uuid.UUID) ([]models.OrderItem, error) {
	rows, err := r.db.Query(ctx, `SELECT id, order_id, product_id, quantity, unit_price, currency, created_at, updated_at, deleted_at
		FROM order_items WHERE order_id = $1 AND deleted_at IS NULL`, orderID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.OrderItem])
}

func (r *pgxOrderRepository) GetPendingOutboxEvents(ctx context.Context, limit int) ([]models.OutboxEvent, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, event_id, aggregate_type, aggregate_id, event_type, event_version,
			payload, status, attempts, available_at, published_at, last_error, created_at
		FROM outbox_events
		WHERE status = 'pending' AND available_at <= NOW()
		ORDER BY created_at
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.OutboxEvent])
}

func (r *pgxOrderRepository) MarkOutboxPublished(ctx context.Context, id uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `
		UPDATE outbox_events
		SET status = 'published', published_at = NOW()
		WHERE id = $1 AND status = 'pending'`, id)
	return err
}

func (r *pgxOrderRepository) MarkOutboxFailed(ctx context.Context, id uuid.UUID, cause string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `
		UPDATE outbox_events
		SET attempts = attempts + 1,
			last_error = $2,
			available_at = NOW() + (INTERVAL '1 second' * LEAST(300, POWER(2, attempts + 1)))
		WHERE id = $1 AND status = 'pending'`, id, cause)
	return err
}

func insertOutboxEvent(ctx context.Context, tx pgx.Tx, event *models.OutboxEvent) error {
	if event == nil {
		return nil
	}
	if event.AvailableAt.IsZero() {
		event.AvailableAt = time.Now().UTC()
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO outbox_events
			(id, event_id, aggregate_type, aggregate_id, event_type, event_version, payload, status, available_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', $8)`,
		event.ID, event.EventID, event.AggregateType, event.AggregateID,
		event.EventType, event.EventVersion, event.Payload, event.AvailableAt)
	return err
}
