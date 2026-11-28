package repository

import (
	"context"

	"github.com/Rizwank123/LahziCart/notification-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxNotificationRepository struct{ db *pgxpool.Pool }

func NewNotificationRepository(db *pgxpool.Pool) models.NotificationRepository {
	return &pgxNotificationRepository{db: db}
}

func (r *pgxNotificationRepository) Create(ctx context.Context, n *models.Notification) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `INSERT INTO notifications (user_id, channel, type, recipient, subject, body, status, metadata)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, q, n.UserID, n.Channel, n.Type, n.Recipient, n.Subject, n.Body, n.Status, n.Metadata).
		Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)
}

func (r *pgxNotificationRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]models.Notification, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `SELECT id,user_id,channel,type,recipient,subject,body,status,metadata,created_at,updated_at,deleted_at FROM notifications WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Notification])
}

func (r *pgxNotificationRepository) MarkSent(ctx context.Context, id uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `UPDATE notifications SET status = 'sent', updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (r *pgxNotificationRepository) Get(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rows, err := r.db.Query(ctx, `SELECT id,user_id,channel,type,recipient,subject,body,status,metadata,created_at,updated_at,deleted_at FROM notifications WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}
	n, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Notification])
	if err != nil {
		return nil, err
	}
	return &n, nil
}
