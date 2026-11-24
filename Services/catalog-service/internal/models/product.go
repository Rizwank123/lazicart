package models

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type (
	// Product model for product
	Product struct {
		Base
		Name        string                 `json:"name" db:"name" example:"MacBook Pro M5"`
		Description *string                `json:"description,omitempty" db:"description" example:"MacBook Pro M5 16GB RAM 512GB SSD"`
		Price       float64                `json:"price" db:"price" example:"115290"`
		Currency    string                 `json:"currency" db:"currency" example:"INR"`
		CategoryID  *uuid.UUID             `json:"category_id,omitempty" db:"category_id" example:"a6c5a8d0-3c3c-4c5b-9a6c-0c9e7f4b9c5a"`
		Brand       *string                `json:"brand,omitempty" db:"brand" example:"Apple"`
		Attributes  map[string]interface{} `json:"attributes" db:"attributes" example:"{\"ram\":\"16GB\",\"storage\":\"512GB\",\"color\":\"Space Gray\"}"`
		IsActive    bool                   `json:"is_active" db:"is_active" example:"true"`
		Images      []ProductImage         `json:"images,omitempty" db:"images"`
		Review      []Review               `json:"review,omitempty" db:"review"`
		BaseAudit
	} // @name Product

	// ProductImage defines the model for product images
	ProductImage struct {
		Base
		ProductID    uuid.UUID `json:"product_id" db:"product_id" example:"d9a8f0d2-9e88-4c67-b7c2-8a4c78b76d01"`
		ImageURL     string    `json:"image_url" db:"image_url" example:"https://cdn.lahzicart.com/products/macbook-pro/main.jpg"`
		ThumbnailURL *string   `json:"thumbnail_url,omitempty" db:"thumbnail_url" example:"https://cdn.lahzicart.com/products/macbook-pro/thumb.jpg"`

		IsPrimary    bool `json:"is_primary" db:"is_primary" example:"true"`
		DisplayOrder int  `json:"display_order" db:"display_order" example:"1"`
		BaseAudit
	} // @name ProductImage

	// Review defines the model for Review
	Review struct {
		Base
		ProductID  uuid.UUID `json:"product_id" db:"product_id" example:"d9a8f0d2-9e88-4c67-b7c2-8a4c78b76d01"`
		UserID     uuid.UUID `json:"user_id" db:"user_id" example:"7c4c51c3-8f84-4d8f-a5d5-0e5f6d7b8a2e"`
		Rating     int       `json:"rating" db:"rating" example:"5"`
		Comment    *string   `json:"comment,omitempty" db:"comment" example:"Excellent performance and battery life."`
		IsApproved bool      `json:"is_approved" db:"is_approved" example:"true"`
		BaseAudit
	} // @name Review

	// ProductRating defines the model for product rating
	ProductRating struct {
		ProductID     uuid.UUID `json:"product_id" db:"product_id" example:"d9a8f0d2-9e88-4c67-b7c2-8a4c78b76d01"`
		AverageRating float64   `json:"average_rating" db:"average_rating" example:"4.7"`
		TotalReviews  int       `json:"total_reviews" db:"total_reviews" example:"132"`
	} // @name ProductRating
)

type (
	// CreateProductInput defines the model for CreateProduct
	CreateProductInput struct {
		Name        string                 `json:"name" example:"MacBook Pro M5"`
		Description *string                `json:"description,omitempty" example:"MacBook Pro M5 16GB RAM 512GB SSD"`
		Price       float64                `json:"price" example:"115290"`
		Currency    string                 `json:"currency" example:"INR"`
		CategoryID  *uuid.UUID             `json:"category_id,omitempty" example:"a6c5a8d0-3c3c-4c5b-9a6c-0c9e7f4b9c5a"`
		Brand       *string                `json:"brand,omitempty" example:"Apple"`
		Attributes  map[string]interface{} `json:"attributes" example:"{\"ram\":\"16GB\",\"storage\":\"512GB\"}"`
		IsActive    bool                   `json:"is_active" example:"true"`
	} // @name CreateProductInput

	// UpdateProductInput defines the model for update product
	UpdateProductInput struct {
		Name        string                 `json:"name,omitempty" example:"MacBook Pro M5"`
		Description *string                `json:"description,omitempty" example:"MacBook Pro M5 Updated Model"`
		Price       float64                `json:"price,omitempty" example:"119999"`
		Currency    string                 `json:"currency,omitempty" example:"INR"`
		CategoryID  *uuid.UUID             `json:"category_id,omitempty" example:"a6c5a8d0-3c3c-4c5b-9a6c-0c9e7f4b9c5a"`
		Brand       *string                `json:"brand,omitempty" example:"Apple"`
		Attributes  map[string]interface{} `json:"attributes,omitempty" example:"{\"ram\":\"32GB\",\"storage\":\"1TB\"}"`
		IsActive    bool                   `json:"is_active,omitempty" example:"true"`
	} // @name UpdateProductInput

	// CreateReviewInput defines the model for CreateReviewInput
	CreateReviewInput struct {
		ProductID  uuid.UUID `json:"product_id" example:"d9a8f0d2-9e88-4c67-b7c2-8a4c78b76d01"`
		UserID     uuid.UUID `json:"user_id" example:"7c4c51c3-8f84-4d8f-a5d5-0e5f6d7b8a2e"`
		Rating     int       `json:"rating" example:"5"`
		Comment    *string   `json:"comment,omitempty" example:"Amazing laptop with blazing fast performance."`
		IsApproved bool      `json:"is_approved" example:"true"`
	} // @name CreateReviewInput
	// UpdateReviewInput defines the model for Update the Review
	UpdateReviewInput struct {
		Rating  int     `json:"rating" example:"5"`
		Comment *string `json:"comment,omitempty" example:"Amazing laptop with blazing fast performance."`
	} // @name UpdateReviewInput
	// CreateProductImageInput defines the model for CreateProductImage
	CreateProductImageInput struct {
		ProductID    uuid.UUID `json:"product_id" example:"d9a8f0d2-9e88-4c67-b7c2-8a4c78b76d01"`
		ImageURL     string    `json:"image_url" example:"https://cdn.lahzicart.com/products/macbook-pro/main.jpg"`
		ThumbnailURL *string   `json:"thumbnail_url,omitempty" example:"https://cdn.lahzicart.com/products/macbook-pro/thumb.jpg"`
		IsPrimary    bool      `json:"is_primary" example:"true"`
		DisplayOrder int       `json:"display_order" example:"1"`
	} // @name CreateProductImageInput

	// UpdateProductImageInput Defines the model for UpdateProductImageInput
	UpdateProductImageInput struct {
		ImageURL     string  `json:"image_url" example:"https://cdn.lahzicart.com/products/macbook-pro/main.jpg"`
		ThumbnailURL *string `json:"thumbnail_url,omitempty" example:"https://cdn.lahzicart.com/products/macbook-pro/thumb.jpg"`
		IsPrimary    bool    `json:"is_primary" example:"true"`
		DisplayOrder int     `json:"display_order" example:"1"`
	} // @name UpdateProductImageInput
	// FilterInput models  defines the fields that can use to filter the products
	FilterInput struct {
		Price       float64     `json:"price,omitempty"  example:"115290"`
		Currency    string      `json:"currency,omitempty" example:"INR"`
		CategoryIDs []uuid.UUID `json:"category_ids,omitempty" example:"a6c5a8d0-3c3c-4c5b-9a6c-0c9e7f4b9c5a"`
		Brand       *string     `json:"brand,omitempty" example:"Apple"`
		SortKey     []SortKey   `json:"sortkey,omitempty"`
	} // @name FilterInput
)
type (
	// ProductRepository defines the methods that any product Repository implements
	ProductRepository interface {
		// Save Product save the records into database
		SaveProduct(ctx context.Context, in *Product) (err error)
		// Update updates the existing product
		UpdateProduct(ctx context.Context, in *Product) (err error)
		// GetProductById returns product by id
		GetProductById(ctx context.Context, id uuid.UUID) (result Product, err error)
		// Get All Prodducts return product with pagination
		GetAllProducts(ctx context.Context, limit, offset int64) (total int64, result []Product, err error)
		// SaveReview save the review crosponding to Product
		SaveProductReview(ctx context.Context, in *Review) (err error)
		// SaveProductImage saves product image records
		SaveProductImage(ctx context.Context, in *ProductImage) (err error)
		// GetReviewByProduct id returs the review for the product
		GetReviewsByProductId(ctx context.Context, productID uuid.UUID) (result []Review, err error)
		// FilterProduct returns the filterd products  with sorting
		FilterProduct(ctx context.Context, in FilterInput, pageno, limit int64) (result []Product, total int64, err error)
		// GetProductReview return the review for the product
		GetProductReview(ctx context.Context, productID uuid.UUID) (resutl []Review, err error)
		// GetProductImages return the images of product
		GetProductImages(ctx context.Context, productID uuid.UUID) (result []ProductImage, err error)
		// GetImageByID return the image By id
		GetImageByID(ctx context.Context, id uuid.UUID) (result ProductImage, err error)
		// GetReviewByID  return review by id
		GetReviewByID(ctx context.Context, id uuid.UUID) (result Review, err error)
		// DeleteProduct delete the product by ID
		DeleteProdduct(ctx context.Context, id uuid.UUID) (err error)
		// DeleteReview delete the Review by id
		DeleteReview(ctx context.Context, id uuid.UUID) (err error)
		// DeleteProductImage delete the product image by id
		DeleteProductImage(ctx context.Context, id uuid.UUID) (err error)
		// UpdateReview updates an existing review
		UpdateReview(ctx context.Context, in *Review) (err error)
		// UpdateProductImage updates an existing product image
		UpdateProductImage(ctx context.Context, in *ProductImage) (err error)
	}
	// ProductService defines the methods that any product service implements
	ProductService interface {
		// SaveProduct  save the product
		SaveProduct(in CreateProductInput) (result Product, err error)
		// GetProductByID  return the product by id
		GetProductByID(id uuid.UUID) (result Product, err error)
		// SaveReview save the
		SaveReview(in CreateReviewInput) (result Product, err error)
		// SaveProductImage saves an image for a product
		SaveProductImage(in CreateProductImageInput) (result ProductImage, err error)
		// GetRevicwForProduct return reviews of  product
		GetReviewForProduct(productId uuid.UUID) (result []Product, err error)
		// GetProducts restuns the all products
		GetProducts(limit, pageNo int) (result []Product, total int64, err error)
		// UpdateProduct update the exisitng recod
		UpdateProduct(id uuid.UUID, in UpdateProductInput) (result Product, err error)
		//UpdateReview update existing Review
		UpdateReview(id uuid.UUID, in UpdateReviewInput) (result Review, err error)
		// UpdateProductImage  update image
		UpdateProductImage(id uuid.UUID, in UpdateProductImageInput) (result ProductImage, err error)
		// FilterProducts  returns the products for
		FilterProducts(in FilterInput, pageNo, limt int64) (result []Product, total int64, err error)
		// FindProductById return the product by id
		FindProductById(id uuid.UUID) (result Product, err error)
		// FindProductImagesByProductId return the images of product
		FindProductImagesByProductId(productID uuid.UUID) (result []ProductImage, err error)
		// FindReviewForProduct return the review of the product
		FindReviewForProduct(productID uuid.UUID) (result []Review, err error)
		// DeleteProdduct delete the product by id
		DeleteProdduct(id uuid.UUID) (err error)
		// DeleteProductImage delete the spacific image
		DeleteProductImage(id uuid.UUID) (err error)
		// DeleteReview Delete the review by id
		DeleteReview(id uuid.UUID) (err error)
	}
)
