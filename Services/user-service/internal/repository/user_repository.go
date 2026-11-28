package repository

import (
	"context"
	"errors"

	"github.com/Rizwank123/lahziCart/user-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
)

type pgxUserRepository struct {
	db *pgxpool.Pool
}

func NewPgxUserRepository(db *pgxpool.Pool) models.UserRepository {
	return &pgxUserRepository{db: db}
}

// CreateUser implements [models.UserRepository].
func (r *pgxUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if ctx == nil {
		ctx = context.Background()
	}
	// The profile id is the auth-service id so that one UUID identifies the
	// user everywhere: in the JWT subject claim, in the users table and in the
	// cart, order, payment, notification and delivery services.
	q := `INSERT INTO users (id, first_name, last_name, email, phone, is_active, dob, address, auth_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`
	args := []interface{}{
		user.AuthId,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.IsActive,
		user.Dob,
		user.Address,
		user.AuthId,
	}
	return r.db.QueryRow(ctx, q, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *pgxUserRepository) ProcessUserCreated(ctx context.Context, eventID uuid.UUID, in models.CreateUserRequest) error {
	if ctx == nil {
		ctx = context.Background()
	}
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var claimed uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO inbox_messages (event_id, event_type)
		VALUES ($1, $2)
		ON CONFLICT (event_id) DO NOTHING
		RETURNING event_id`, eventID, "user.created").Scan(&claimed)
	if errors.Is(err, pgx.ErrNoRows) {
		return tx.Commit(ctx)
	}
	if err != nil {
		return err
	}

	// The profile id is the auth-service id so that one UUID identifies the
	// user everywhere: in the JWT subject claim, in the users table and in the
	// cart, order, payment, notification and delivery services.
	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, first_name, last_name, email, phone, is_active, dob, address, auth_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		in.AuthId, in.FirstName, in.LastName, in.Email, in.Phone, true, in.Dob, in.Address, in.AuthId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE inbox_messages
		SET status = 'processed', processed_at = NOW(), updated_at = NOW()
		WHERE event_id = $1`, eventID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// DeleteUser implements [models.UserRepository].
func (r *pgxUserRepository) DeleteUser(ctx context.Context, user *models.User) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `UPDATE users SET deleted_at = NOW() WHERE id = $1`
	args := []interface{}{
		user.ID,
	}
	_, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByEmail implements [models.UserRepository].
func (r *pgxUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL`
	args := []interface{}{
		email,
	}
	var row pgx.Rows
	row, err := r.db.Query(ctx, q, args...)
	resutlt, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByNameLax[models.User])
	if err != nil {
		return nil, err
	}
	return &resutlt, nil
}

// GetUserByID implements [models.UserRepository].
//
// The identifier may be either the users primary key or the auth_id issued by
// the auth-service. Clients receive the auth-service identifier in the JWT
// subject claim, while the cart, order, payment, notification and delivery
// services are all keyed by the users primary key. Resolving both here keeps
// those services consistent with each other.
func (r *pgxUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`
	args := []interface{}{
		id,
	}
	var row pgx.Rows
	row, err := r.db.Query(ctx, q, args...)
	result, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByNameLax[models.User])
	if err == nil {
		return &result, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Fall back to the auth-service identifier.
	q = `SELECT * FROM users WHERE auth_id = $1 AND deleted_at IS NULL`
	args = []interface{}{
		id,
	}
	row, err = r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	result, err = pgx.CollectExactlyOneRow(row, pgx.RowToStructByNameLax[models.User])
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUserByPhone implements [models.UserRepository].
func (r *pgxUserRepository) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	// 1. Idiomatic context handling
	if ctx == nil {
		ctx = context.TODO()
	}

	// 2. Explicitly select columns instead of * (Best Practice)
	q := `SELECT id, first_name, last_name, email, phone, is_active, dob, address, created_at, updated_at, auth_id 
          FROM users 
          WHERE phone = $1 AND deleted_at IS NULL 
          LIMIT 1`

	// 3. Use Query (as you were) but handle the row collection properly
	rows, err := r.db.Query(ctx, q, phone)
	if err != nil {
		log.Errorf("database query error: %v", err)
		return nil, err
	}
	// CollectExactlyOneRow will automatically close the 'rows'

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warnf("no user found with phone: %s", phone)
			return nil, nil // Or return a custom ErrNotFound
		}
		return nil, err
	}

	return &result, nil
}

// UpdateUser implements [models.UserRepository].
func (r *pgxUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	if ctx == nil {
		ctx = context.Background()
	}
	q := `UPDATE users SET first_name = $1, last_name = $2, email = $3, phone = $4, is_active = $5, dob = $6, address = $7, updated_at = NOW() WHERE id = $8`
	args := []interface{}{
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.IsActive,
		user.Dob,
		user.Address,
		user.ID,
	}
	_, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}
