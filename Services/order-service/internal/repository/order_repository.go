package repository

import (
	"context"

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

func (r *pgxOrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	if ctx == nil {
		ctx = context.Background()
	}
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := `INSERT INTO orders (user_id, status, total, currency, shipping_address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	if err := tx.QueryRow(ctx, q, order.UserID, order.Status, order.Total, order.Currency, order.ShippingAddress).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return err
	}
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		err := tx.QueryRow(ctx, `INSERT INTO order_items (order_id, product_id, quantity, unit_price, currency)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at`,
			order.ID, order.Items[i].ProductID, order.Items[i].Quantity, order.Items[i].UnitPrice, order.Items[i].Currency).
			Scan(&order.Items[i].ID, &order.Items[i].CreatedAt, &order.Items[i].UpdatedAt)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
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

func (r *pgxOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`, status, id)
	return err
}

func (r *pgxOrderRepository) listItems(ctx context.Context, orderID uuid.UUID) ([]models.OrderItem, error) {
	rows, err := r.db.Query(ctx, `SELECT id, order_id, product_id, quantity, unit_price, currency, created_at, updated_at, deleted_at
		FROM order_items WHERE order_id = $1 AND deleted_at IS NULL`, orderID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.OrderItem])
}
