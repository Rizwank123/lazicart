package repository

import (
	"context"
	"errors"

	"github.com/Rizwank123/LahziCart/inventory-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxInventoryRepository struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) models.InventoryRepository {
	return &pgxInventoryRepository{db: db}
}

func (r *pgxInventoryRepository) UpsertItem(ctx context.Context, item *models.InventoryItem) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `
		INSERT INTO inventory_items (product_id, sku, quantity_available, reorder_level, is_active)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (product_id)
		DO UPDATE SET sku = EXCLUDED.sku,
			quantity_available = EXCLUDED.quantity_available,
			reorder_level = EXCLUDED.reorder_level,
			is_active = EXCLUDED.is_active,
			updated_at = NOW()
		RETURNING id, quantity_reserved, created_at, updated_at, deleted_at`
	return r.db.QueryRow(ctx, q, item.ProductID, item.SKU, item.QuantityAvailable, item.ReorderLevel, item.IsActive).
		Scan(&item.ID, &item.QuantityReserved, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
}

func (r *pgxInventoryRepository) GetByProductID(ctx context.Context, productID uuid.UUID) (*models.InventoryItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `SELECT id, product_id, sku, quantity_available, quantity_reserved, reorder_level, is_active, created_at, updated_at, deleted_at
		FROM inventory_items WHERE product_id = $1 AND deleted_at IS NULL`
	rows, err := r.db.Query(ctx, q, productID)
	if err != nil {
		return nil, err
	}
	item, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.InventoryItem])
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *pgxInventoryRepository) ListItems(ctx context.Context, limit, offset int64) (int64, []models.InventoryItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM inventory_items WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return 0, nil, err
	}
	q := `SELECT id, product_id, sku, quantity_available, quantity_reserved, reorder_level, is_active, created_at, updated_at, deleted_at
		FROM inventory_items WHERE deleted_at IS NULL ORDER BY updated_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, q, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.InventoryItem])
	return total, items, err
}

func (r *pgxInventoryRepository) Adjust(ctx context.Context, productID uuid.UUID, delta int) (*models.InventoryItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `UPDATE inventory_items
		SET quantity_available = quantity_available + $1, updated_at = NOW()
		WHERE product_id = $2 AND deleted_at IS NULL AND quantity_available + $1 >= 0
		RETURNING id, product_id, sku, quantity_available, quantity_reserved, reorder_level, is_active, created_at, updated_at, deleted_at`
	rows, err := r.db.Query(ctx, q, delta, productID)
	if err != nil {
		return nil, err
	}
	item, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.InventoryItem])
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *pgxInventoryRepository) Reserve(ctx context.Context, orderID uuid.UUID, items []models.ReserveInventoryItem) ([]models.InventoryReservation, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	reservations := make([]models.InventoryReservation, 0, len(items))
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than zero")
		}
		tag, err := tx.Exec(ctx, `
			UPDATE inventory_items
			SET quantity_available = quantity_available - $1,
				quantity_reserved = quantity_reserved + $1,
				updated_at = NOW()
			WHERE product_id = $2 AND deleted_at IS NULL AND is_active = true AND quantity_available >= $1`,
			item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() == 0 {
			return nil, errors.New("insufficient stock")
		}

		var reservation models.InventoryReservation
		err = tx.QueryRow(ctx, `
			INSERT INTO inventory_reservations (order_id, product_id, quantity, status)
			VALUES ($1, $2, $3, 'reserved')
			ON CONFLICT (order_id, product_id)
			DO UPDATE SET quantity = inventory_reservations.quantity + EXCLUDED.quantity,
				status = 'reserved',
				updated_at = NOW()
			RETURNING id, order_id, product_id, quantity, status, created_at, updated_at, deleted_at`,
			orderID, item.ProductID, item.Quantity).
			Scan(&reservation.ID, &reservation.OrderID, &reservation.ProductID, &reservation.Quantity, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt, &reservation.DeletedAt)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *pgxInventoryRepository) Release(ctx context.Context, orderID uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `SELECT product_id, quantity FROM inventory_reservations WHERE order_id = $1 AND status = 'reserved'`, orderID)
	if err != nil {
		return err
	}
	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.ReserveInventoryItem])
	if err != nil {
		return err
	}
	for _, reservation := range reservations {
		if _, err := tx.Exec(ctx, `
			UPDATE inventory_items
			SET quantity_available = quantity_available + $1,
				quantity_reserved = quantity_reserved - $1,
				updated_at = NOW()
			WHERE product_id = $2`,
			reservation.Quantity, reservation.ProductID); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(ctx, `UPDATE inventory_reservations SET status = 'released', updated_at = NOW() WHERE order_id = $1 AND status = 'reserved'`, orderID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgxInventoryRepository) Commit(ctx context.Context, orderID uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `SELECT product_id, quantity FROM inventory_reservations WHERE order_id = $1 AND status = 'reserved'`, orderID)
	if err != nil {
		return err
	}
	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.ReserveInventoryItem])
	if err != nil {
		return err
	}
	for _, reservation := range reservations {
		if _, err := tx.Exec(ctx, `
			UPDATE inventory_items
			SET quantity_reserved = quantity_reserved - $1,
				updated_at = NOW()
			WHERE product_id = $2`,
			reservation.Quantity, reservation.ProductID); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(ctx, `UPDATE inventory_reservations SET status = 'committed', updated_at = NOW() WHERE order_id = $1 AND status = 'reserved'`, orderID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
