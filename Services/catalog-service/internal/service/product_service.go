package service

import (
	"context"
	"errors"

	"github.com/Rizwank123/LahziCart/catalog-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type ProductServiceImpl struct {
	repo   models.ProductRepository
	logger *zap.Logger
}

func NewProductService(repo models.ProductRepository, logger *zap.Logger) models.ProductService {
	return &ProductServiceImpl{repo: repo, logger: logger}
}

func (s *ProductServiceImpl) SaveProduct(in models.CreateProductInput) (models.Product, error) {
	if in.Name == "" {
		return models.Product{}, errors.New("name is required")
	}
	if in.Price <= 0 {
		return models.Product{}, errors.New("price must be greater than zero")
	}
	if in.Currency == "" {
		in.Currency = "INR"
	}
	if in.Attributes == nil {
		in.Attributes = map[string]interface{}{}
	}

	product := models.Product{
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Currency:    in.Currency,
		CategoryID:  in.CategoryID,
		Brand:       in.Brand,
		Attributes:  in.Attributes,
		IsActive:    true,
	}
	if err := s.repo.SaveProduct(context.Background(), &product); err != nil {
		s.logger.Error("failed to save product", zap.Error(err))
		return models.Product{}, err
	}
	return product, nil
}

func (s *ProductServiceImpl) GetProductByID(id uuid.UUID) (models.Product, error) {
	return s.repo.GetProductById(context.Background(), id)
}

func (s *ProductServiceImpl) SaveReview(in models.CreateReviewInput) (models.Product, error) {
	if in.Rating < 1 || in.Rating > 5 {
		return models.Product{}, errors.New("rating must be between 1 and 5")
	}

	review := models.Review{
		ProductID:  in.ProductID,
		UserID:     in.UserID,
		Rating:     in.Rating,
		Comment:    in.Comment,
		IsApproved: in.IsApproved,
	}
	if err := s.repo.SaveProductReview(context.Background(), &review); err != nil {
		s.logger.Error("failed to save review", zap.Error(err))
		return models.Product{}, err
	}
	return s.repo.GetProductById(context.Background(), in.ProductID)
}

func (s *ProductServiceImpl) SaveProductImage(in models.CreateProductImageInput) (models.ProductImage, error) {
	if in.ImageURL == "" {
		return models.ProductImage{}, errors.New("image_url is required")
	}

	image := models.ProductImage{
		ProductID:    in.ProductID,
		ImageURL:     in.ImageURL,
		ThumbnailURL: in.ThumbnailURL,
		IsPrimary:    in.IsPrimary,
		DisplayOrder: in.DisplayOrder,
	}
	if err := s.repo.SaveProductImage(context.Background(), &image); err != nil {
		s.logger.Error("failed to save product image", zap.Error(err))
		return models.ProductImage{}, err
	}
	return image, nil
}

func (s *ProductServiceImpl) GetReviewForProduct(productID uuid.UUID) ([]models.Product, error) {
	product, err := s.repo.GetProductById(context.Background(), productID)
	if err != nil {
		return nil, err
	}
	return []models.Product{product}, nil
}

func (s *ProductServiceImpl) GetProducts(limit, pageNo int) ([]models.Product, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if pageNo <= 0 {
		pageNo = 1
	}
	offset := int64((pageNo - 1) * limit)
	total, products, err := s.repo.GetAllProducts(context.Background(), int64(limit), offset)
	return products, total, err
}

func (s *ProductServiceImpl) UpdateProduct(id uuid.UUID, in models.UpdateProductInput) (models.Product, error) {
	product, err := s.repo.GetProductById(context.Background(), id)
	if err != nil {
		return models.Product{}, err
	}

	if in.Name != "" {
		product.Name = in.Name
	}
	if in.Description != nil {
		product.Description = in.Description
	}
	if in.Price > 0 {
		product.Price = in.Price
	}
	if in.Currency != "" {
		product.Currency = in.Currency
	}
	if in.CategoryID != nil {
		product.CategoryID = in.CategoryID
	}
	if in.Brand != nil {
		product.Brand = in.Brand
	}
	if in.Attributes != nil {
		product.Attributes = in.Attributes
	}
	product.IsActive = in.IsActive

	if err := s.repo.UpdateProduct(context.Background(), &product); err != nil {
		return models.Product{}, err
	}
	return s.repo.GetProductById(context.Background(), id)
}

func (s *ProductServiceImpl) UpdateReview(id uuid.UUID, in models.UpdateReviewInput) (models.Review, error) {
	review, err := s.repo.GetReviewByID(context.Background(), id)
	if err != nil {
		return models.Review{}, err
	}
	if in.Rating < 1 || in.Rating > 5 {
		return models.Review{}, errors.New("rating must be between 1 and 5")
	}
	review.Rating = in.Rating
	review.Comment = in.Comment
	if err := s.repo.UpdateReview(context.Background(), &review); err != nil {
		return models.Review{}, err
	}
	return s.repo.GetReviewByID(context.Background(), id)
}

func (s *ProductServiceImpl) UpdateProductImage(id uuid.UUID, in models.UpdateProductImageInput) (models.ProductImage, error) {
	image, err := s.repo.GetImageByID(context.Background(), id)
	if err != nil {
		return models.ProductImage{}, err
	}
	if in.ImageURL != "" {
		image.ImageURL = in.ImageURL
	}
	if in.ThumbnailURL != nil {
		image.ThumbnailURL = in.ThumbnailURL
	}
	image.IsPrimary = in.IsPrimary
	image.DisplayOrder = in.DisplayOrder
	if err := s.repo.UpdateProductImage(context.Background(), &image); err != nil {
		return models.ProductImage{}, err
	}
	return s.repo.GetImageByID(context.Background(), id)
}

func (s *ProductServiceImpl) FilterProducts(in models.FilterInput, pageNo, limit int64) ([]models.Product, int64, error) {
	return s.repo.FilterProduct(context.Background(), in, pageNo, limit)
}

func (s *ProductServiceImpl) FindProductById(id uuid.UUID) (models.Product, error) {
	return s.repo.GetProductById(context.Background(), id)
}

func (s *ProductServiceImpl) FindProductImagesByProductId(productID uuid.UUID) ([]models.ProductImage, error) {
	return s.repo.GetProductImages(context.Background(), productID)
}

func (s *ProductServiceImpl) FindReviewForProduct(productID uuid.UUID) ([]models.Review, error) {
	return s.repo.GetProductReview(context.Background(), productID)
}

func (s *ProductServiceImpl) DeleteProdduct(id uuid.UUID) error {
	return s.repo.DeleteProdduct(context.Background(), id)
}

func (s *ProductServiceImpl) DeleteProductImage(id uuid.UUID) error {
	return s.repo.DeleteProductImage(context.Background(), id)
}

func (s *ProductServiceImpl) DeleteReview(id uuid.UUID) error {
	return s.repo.DeleteReview(context.Background(), id)
}
