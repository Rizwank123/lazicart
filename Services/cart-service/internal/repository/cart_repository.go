package repository

import (
	"context"

	"github.com/Rizwank123/LahziCart/cart-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxCartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) models.CartRepository {
	return &pgxCartRepository{db: db}
}

func (r *pgxCartRepository) GetOrCreateCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `
		INSERT INTO carts (user_id, status)
		VALUES ($1, 'active')
		ON CONFLICT (user_id)
		DO UPDATE SET updated_at = carts.updated_at
		RETURNING id, user_id, status, created_at, updated_at, deleted_at`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	cart, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Cart])
	if err != nil {
		return nil, err
	}
	cart.Items, _ = r.ListItems(ctx, cart.ID)
	return &cart, nil
}

func (r *pgxCartRepository) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `SELECT id, user_id, status, created_at, updated_at, deleted_at FROM carts WHERE user_id = $1 AND deleted_at IS NULL`, userID)
	if err != nil {
		return nil, err
	}
	cart, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Cart])
	if err != nil {
		return nil, err
	}
	cart.Items, _ = r.ListItems(ctx, cart.ID)
	return &cart, nil
}

func (r *pgxCartRepository) AddItem(ctx context.Context, cartID uuid.UUID, item *models.CartItem) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `
		INSERT INTO cart_items (cart_id, product_id, quantity, unit_price, currency)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity,
			unit_price = EXCLUDED.unit_price,
			currency = EXCLUDED.currency,
			updated_at = NOW()
		RETURNING id, created_at, updated_at, deleted_at`
	return r.db.QueryRow(ctx, q, cartID, item.ProductID, item.Quantity, item.UnitPrice, item.Currency).
		Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
}

func (r *pgxCartRepository) UpdateItem(ctx context.Context, itemID uuid.UUID, quantity int) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `UPDATE cart_items SET quantity = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`, quantity, itemID)
	return err
}

func (r *pgxCartRepository) RemoveItem(ctx context.Context, itemID uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `UPDATE cart_items SET deleted_at = NOW() WHERE id = $1`, itemID)
	return err
}

func (r *pgxCartRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `
		UPDATE cart_items SET deleted_at = NOW()
		WHERE cart_id IN (SELECT id FROM carts WHERE user_id = $1) AND deleted_at IS NULL`, userID)
	return err
}

func (r *pgxCartRepository) ListItems(ctx context.Context, cartID uuid.UUID) ([]models.CartItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, cart_id, product_id, quantity, unit_price, currency, created_at, updated_at, deleted_at
		FROM cart_items
		WHERE cart_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`, cartID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.CartItem])
}
