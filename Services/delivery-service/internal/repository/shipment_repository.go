package repository

import (
	"context"

	"github.com/Rizwank123/LahziCart/delivery-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxShipmentRepository struct{ db *pgxpool.Pool }

func NewShipmentRepository(db *pgxpool.Pool) models.ShipmentRepository {
	return &pgxShipmentRepository{db: db}
}

func (r *pgxShipmentRepository) Create(ctx context.Context, s *models.Shipment) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `INSERT INTO shipments (order_id,user_id,carrier,tracking_number,status,shipping_address,estimated_delivery_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id,created_at,updated_at`
	return r.db.QueryRow(ctx, q, s.OrderID, s.UserID, s.Carrier, s.TrackingNumber, s.Status, s.ShippingAddress, s.EstimatedDeliveryAt).
		Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *pgxShipmentRepository) Get(ctx context.Context, id uuid.UUID) (*models.Shipment, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `SELECT id,order_id,user_id,carrier,tracking_number,status,shipping_address,estimated_delivery_at,delivered_at,created_at,updated_at,deleted_at FROM shipments WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}
	shipment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Shipment])
	return &shipment, err
}

func (r *pgxShipmentRepository) GetByOrder(ctx context.Context, orderID uuid.UUID) (*models.Shipment, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `SELECT id,order_id,user_id,carrier,tracking_number,status,shipping_address,estimated_delivery_at,delivered_at,created_at,updated_at,deleted_at FROM shipments WHERE order_id=$1 AND deleted_at IS NULL`, orderID)
	if err != nil {
		return nil, err
	}
	shipment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Shipment])
	return &shipment, err
}

func (r *pgxShipmentRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]models.Shipment, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `SELECT id,order_id,user_id,carrier,tracking_number,status,shipping_address,estimated_delivery_at,delivered_at,created_at,updated_at,deleted_at FROM shipments WHERE user_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Shipment])
}

func (r *pgxShipmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, in models.UpdateShipmentStatusInput) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `UPDATE shipments SET status=$1, carrier=COALESCE($2, carrier), tracking_number=COALESCE($3, tracking_number), delivered_at=COALESCE($4, delivered_at), updated_at=NOW() WHERE id=$5 AND deleted_at IS NULL`, in.Status, in.Carrier, in.TrackingNumber, in.DeliveredAt, id)
	return err
}
