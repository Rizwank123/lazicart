package repository

import (
	"context"
	"log"
	"time"

	"github.com/Rizwank123/LahziCart/auth-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxAuthRepository struct {
	pgxPool *pgxpool.Pool
}

func NewPgxRepository(pgxPool *pgxpool.Pool) models.AuthRepository {
	return &PgxAuthRepository{pgxPool: pgxPool}
}

// CreateUserAuth implements [models.AuthRepository].
func (p *PgxAuthRepository) CreateUserAuth(ctx context.Context, userAuth *models.UserAuth) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `INSERT INTO user_auth (email, phone, password, is_active) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	tx, err := p.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	args := []interface{}{
		userAuth.Email,
		userAuth.Phone,
		userAuth.Password,
		userAuth.IsActive,
	}
	tx.QueryRow(ctx, q, args...).Scan(&userAuth.ID, &userAuth.CreatedAt, &userAuth.UpdatedAt)
	tx.Commit(ctx)

	return err

}

// DeleteUserAuth implements [models.AuthRepository].
func (p *PgxAuthRepository) DeleteUserAuth(ctx context.Context, id uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `UPDATE user_auth SET deleted_at = $1 WHERE id = $2`
	_, err := p.pgxPool.Exec(ctx, q, time.Now(), id)
	return err
}

// GetUserAuthByEmail implements [models.AuthRepository].
func (p *PgxAuthRepository) GetUserAuthByEmail(ctx context.Context, email string) (*models.UserAuth, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `SELECT * FROM user_auth WHERE email = $1`
	log.Println("q", q, email)
	var rows pgx.Rows
	rows, err := p.pgxPool.Query(ctx, q, email)
	if err != nil {
		log.Println("err", err)
		return nil, err
	}
	defer rows.Close()
	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.UserAuth])
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUserAuthById implements [models.AuthRepository].
func (p *PgxAuthRepository) GetUserAuthById(ctx context.Context, id uuid.UUID) (*models.UserAuth, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `SELECT * FROM user_auth WHERE id = $1`

	var rows pgx.Rows
	rows, err := p.pgxPool.Query(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// CollectExactlyOneRow advances the cursor before scanning. Scanning
	// pgx.Rows directly without Next() fails with "number of field
	// descriptions must equal number of values".
	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.UserAuth])
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUserAuthByPhone implements [models.AuthRepository].
func (p *PgxAuthRepository) GetUserAuthByPhone(ctx context.Context, phone string) (*models.UserAuth, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `SELECT id, email, phone, password, is_active, created_at, updated_at FROM user_auth WHERE phone = $1`

	var rows pgx.Rows
	rows, err := p.pgxPool.Query(ctx, q, phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// CollectExactlyOneRow advances the cursor before scanning. Scanning
	// pgx.Rows directly without Next() fails with "number of field
	// descriptions must equal number of values".
	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.UserAuth])
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateUserAuth implements [models.AuthRepository].
func (p *PgxAuthRepository) UpdateUserAuth(ctx context.Context, userAuth *models.UserAuth) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `UPDATE user_auth SET email = $1, phone = $2, password = $3, is_active = $4 WHERE id = $5`
	_, err := p.pgxPool.Exec(ctx, q, userAuth.Email, userAuth.Phone, userAuth.Password, userAuth.IsActive, userAuth.ID)
	return err
}

// DisableUserAuth implements [models.AuthRepository].
func (p *PgxAuthRepository) DisableUserAuth(ctx context.Context, id string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `UPDATE user_auth SET is_active = $1 WHERE id = $2`
	_, err := p.pgxPool.Exec(ctx, q, false, id)
	return err
}
