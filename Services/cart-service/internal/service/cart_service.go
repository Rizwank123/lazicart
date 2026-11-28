package service

import (
	"context"
	"errors"

	"github.com/Rizwank123/LahziCart/cart-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type CartServiceImpl struct {
	repo   models.CartRepository
	logger *zap.Logger
}

func NewCartService(repo models.CartRepository, logger *zap.Logger) models.CartService {
	return &CartServiceImpl{repo: repo, logger: logger}
}

func (s *CartServiceImpl) GetCart(userID uuid.UUID) (*models.Cart, error) {
	return s.repo.GetOrCreateCart(context.Background(), userID)
}

func (s *CartServiceImpl) AddItem(in models.AddCartItemInput) (*models.Cart, error) {
	if in.UserID == uuid.Nil || in.ProductID == uuid.Nil {
		return nil, errors.New("user_id and product_id are required")
	}
	if in.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	if in.Currency == "" {
		in.Currency = "INR"
	}
	cart, err := s.repo.GetOrCreateCart(context.Background(), in.UserID)
	if err != nil {
		return nil, err
	}
	item := &models.CartItem{
		CartID:    cart.ID,
		ProductID: in.ProductID,
		Quantity:  in.Quantity,
		UnitPrice: in.UnitPrice,
		Currency:  in.Currency,
	}
	if err := s.repo.AddItem(context.Background(), cart.ID, item); err != nil {
		s.logger.Error("failed to add cart item", zap.Error(err))
		return nil, err
	}
	return s.repo.GetCartByUserID(context.Background(), in.UserID)
}

func (s *CartServiceImpl) UpdateItem(itemID uuid.UUID, in models.UpdateCartItemInput) error {
	if in.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	return s.repo.UpdateItem(context.Background(), itemID, in.Quantity)
}

func (s *CartServiceImpl) RemoveItem(itemID uuid.UUID) error {
	return s.repo.RemoveItem(context.Background(), itemID)
}

func (s *CartServiceImpl) ClearCart(userID uuid.UUID) error {
	return s.repo.ClearCart(context.Background(), userID)
}

func (s *CartServiceImpl) GetItemOwner(itemID uuid.UUID) (uuid.UUID, error) {
	return s.repo.GetItemOwnerUserID(context.Background(), itemID)
}
