package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/Rizwank123/LahziCart/catalog-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxProductRepository struct {
	db  *pgxpool.Pool
	sqt squirrel.StatementBuilderType
}

func NewProductRepository(db *pgxpool.Pool) models.ProductRepository {
	return &pgxProductRepository{
		db:  db,
		sqt: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *pgxProductRepository) SaveProduct(ctx context.Context, in *models.Product) error {
	if ctx == nil {
		ctx = context.Background()
	}

	attributes, err := json.Marshal(in.Attributes)
	if err != nil {
		return err
	}

	q := `
		INSERT INTO products (name, description, price, currency, category_id, brand, attributes, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, q,
		in.Name,
		in.Description,
		in.Price,
		in.Currency,
		in.CategoryID,
		in.Brand,
		attributes,
		in.IsActive,
	).Scan(&in.ID, &in.CreatedAt, &in.UpdatedAt)
}

func (r *pgxProductRepository) UpdateProduct(ctx context.Context, in *models.Product) error {
	if ctx == nil {
		ctx = context.Background()
	}

	attributes, err := json.Marshal(in.Attributes)
	if err != nil {
		return err
	}

	q := `
		UPDATE products
		SET name = $1,
			description = $2,
			price = $3,
			currency = $4,
			category_id = $5,
			brand = $6,
			attributes = $7,
			is_active = $8,
			updated_at = NOW()
		WHERE id = $9 AND deleted_at IS NULL`

	_, err = r.db.Exec(ctx, q,
		in.Name,
		in.Description,
		in.Price,
		in.Currency,
		in.CategoryID,
		in.Brand,
		attributes,
		in.IsActive,
		in.ID,
	)
	return err
}

func (r *pgxProductRepository) GetProductById(ctx context.Context, id uuid.UUID) (models.Product, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	q := `
		SELECT id, name, description, price, currency, category_id, brand, attributes, is_active,
			created_at, updated_at, deleted_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL`

	rows, err := r.db.Query(ctx, q, id)
	if err != nil {
		return models.Product{}, err
	}

	product, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Product])
	if err != nil {
		return models.Product{}, err
	}
	product.Images, _ = r.GetProductImages(ctx, product.ID)
	product.Review, _ = r.GetProductReview(ctx, product.ID)
	return product, nil
}

func (r *pgxProductRepository) GetAllProducts(ctx context.Context, limit, offset int64) (int64, []models.Product, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return 0, nil, err
	}

	q := `
		SELECT id, name, description, price, currency, category_id, brand, attributes, is_active,
			created_at, updated_at, deleted_at
		FROM products
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, q, limit, offset)
	if err != nil {
		return 0, nil, err
	}

	products, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Product])
	if err != nil {
		return 0, nil, err
	}

	return total, products, nil
}

func (r *pgxProductRepository) SaveProductReview(ctx context.Context, in *models.Review) error {
	if ctx == nil {
		ctx = context.Background()
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := `
		INSERT INTO reviews (product_id, user_id, rating, comment, is_approved)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	if err := tx.QueryRow(ctx, q, in.ProductID, in.UserID, in.Rating, in.Comment, in.IsApproved).
		Scan(&in.ID, &in.CreatedAt, &in.UpdatedAt); err != nil {
		return err
	}

	if err := refreshProductRating(ctx, tx, in.ProductID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgxProductRepository) GetReviewsByProductId(ctx context.Context, productID uuid.UUID) ([]models.Review, error) {
	return r.GetProductReview(ctx, productID)
}

func (r *pgxProductRepository) GetProductReview(ctx context.Context, productID uuid.UUID) ([]models.Review, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	q := `
		SELECT id, product_id, user_id, rating, comment, is_approved, created_at, updated_at, deleted_at
		FROM reviews
		WHERE product_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, q, productID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Review])
}

func (r *pgxProductRepository) FilterProduct(ctx context.Context, in models.FilterInput, pageNo, limit int64) ([]models.Product, int64, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if pageNo < 1 {
		pageNo = 1
	}
	if limit < 1 {
		limit = 20
	}

	countBuilder := r.sqt.Select("COUNT(p.id)").From("products p").Where(sq.Eq{"p.deleted_at": nil})
	dataBuilder := r.sqt.Select(`p.id, p.name, p.description, p.price, p.currency, p.category_id, p.brand,
			p.attributes, p.is_active, p.created_at, p.updated_at, p.deleted_at`).
		From("products p").
		Where(sq.Eq{"p.deleted_at": nil})
	if in.Brand != nil && *in.Brand != "" {
		countBuilder = countBuilder.Where(sq.Eq{"p.brand": *in.Brand})
		dataBuilder = dataBuilder.Where(sq.Eq{"p.brand": *in.Brand})
	}
	if len(in.CategoryIDs) > 0 {
		countBuilder = countBuilder.Where(sq.Eq{"p.category_id": in.CategoryIDs})
		dataBuilder = dataBuilder.Where(sq.Eq{"p.category_id": in.CategoryIDs})
	}
	if in.Price > 0 {
		countBuilder = countBuilder.Where(sq.LtOrEq{"p.price": in.Price})
		dataBuilder = dataBuilder.Where(sq.LtOrEq{"p.price": in.Price})
	}
	if in.Currency != "" {
		countBuilder = countBuilder.Where(sq.Eq{"p.currency": in.Currency})
		dataBuilder = dataBuilder.Where(sq.Eq{"p.currency": in.Currency})
	}

	countQuery, args, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count query failed: %w", err)
	}

	orderBy := "p.created_at DESC"
	if len(in.SortKey) > 0 && in.SortKey[0].Key != "" {
		direction := "ASC"
		if in.SortKey[0].Val == "desc" || in.SortKey[0].Val == "DESC" {
			direction = "DESC"
		}
		switch in.SortKey[0].Key {
		case "price", "created_at", "name":
			orderBy = "p." + in.SortKey[0].Key + " " + direction
		}
	}

	offset := (pageNo - 1) * limit
	dataQuery, args, err := dataBuilder.
		OrderBy(orderBy).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, total, err
	}

	products, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Product])
	return products, total, err
}

func (r *pgxProductRepository) GetProductImages(ctx context.Context, productID uuid.UUID) ([]models.ProductImage, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	q := `
		SELECT id, product_id, image_url, thumbnail_url, is_primary, display_order, created_at, updated_at, deleted_at
		FROM product_images
		WHERE product_id = $1 AND deleted_at IS NULL
		ORDER BY is_primary DESC, display_order ASC, created_at ASC`

	rows, err := r.db.Query(ctx, q, productID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.ProductImage])
}

func (r *pgxProductRepository) GetImageByID(ctx context.Context, id uuid.UUID) (models.ProductImage, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	q := `
		SELECT id, product_id, image_url, thumbnail_url, is_primary, display_order, created_at, updated_at, deleted_at
		FROM product_images
		WHERE id = $1 AND deleted_at IS NULL`

	rows, err := r.db.Query(ctx, q, id)
	if err != nil {
		return models.ProductImage{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.ProductImage])
}

func (r *pgxProductRepository) GetReviewByID(ctx context.Context, id uuid.UUID) (models.Review, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	q := `
		SELECT id, product_id, user_id, rating, comment, is_approved, created_at, updated_at, deleted_at
		FROM reviews
		WHERE id = $1 AND deleted_at IS NULL`

	rows, err := r.db.Query(ctx, q, id)
	if err != nil {
		return models.Review{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Review])
}

func (r *pgxProductRepository) DeleteProdduct(ctx context.Context, id uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `UPDATE products SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *pgxProductRepository) DeleteProductImage(ctx context.Context, id uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := r.db.Exec(ctx, `UPDATE product_images SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *pgxProductRepository) DeleteReview(ctx context.Context, id uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}

	review, err := r.GetReviewByID(ctx, id)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `UPDATE reviews SET deleted_at = NOW() WHERE id = $1`, id); err != nil {
		return err
	}
	if err := refreshProductRating(ctx, tx, review.ProductID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgxProductRepository) SaveProductImage(ctx context.Context, in *models.ProductImage) error {
	if ctx == nil {
		ctx = context.Background()
	}

	q := `
		INSERT INTO product_images (product_id, image_url, thumbnail_url, is_primary, display_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, q, in.ProductID, in.ImageURL, in.ThumbnailURL, in.IsPrimary, in.DisplayOrder).
		Scan(&in.ID, &in.CreatedAt, &in.UpdatedAt)
}

func (r *pgxProductRepository) UpdateProductImage(ctx context.Context, in *models.ProductImage) error {
	if ctx == nil {
		ctx = context.Background()
	}

	q := `
		UPDATE product_images
		SET image_url = $1, thumbnail_url = $2, is_primary = $3, display_order = $4, updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL`

	_, err := r.db.Exec(ctx, q, in.ImageURL, in.ThumbnailURL, in.IsPrimary, in.DisplayOrder, in.ID)
	return err
}

func (r *pgxProductRepository) UpdateReview(ctx context.Context, in *models.Review) error {
	if ctx == nil {
		ctx = context.Background()
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := `
		UPDATE reviews
		SET rating = $1, comment = $2, updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL`
	if _, err := tx.Exec(ctx, q, in.Rating, in.Comment, in.ID); err != nil {
		return err
	}
	if err := refreshProductRating(ctx, tx, in.ProductID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func refreshProductRating(ctx context.Context, tx pgx.Tx, productID uuid.UUID) error {
	q := `
		INSERT INTO product_ratings (product_id, average_rating, total_reviews)
		SELECT $1, COALESCE(AVG(rating), 0), COUNT(*)
		FROM reviews
		WHERE product_id = $1 AND deleted_at IS NULL AND is_approved = true
		ON CONFLICT (product_id)
		DO UPDATE SET average_rating = EXCLUDED.average_rating, total_reviews = EXCLUDED.total_reviews`
	_, err := tx.Exec(ctx, q, productID)
	return err
}
