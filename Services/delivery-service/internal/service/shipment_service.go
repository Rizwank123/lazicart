package service

import (
	"context"
	"errors"

	"github.com/Rizwank123/LahziCart/delivery-service/internal/models"
	"github.com/gofrs/uuid/v5"
)

type ShipmentServiceImpl struct{ repo models.ShipmentRepository }

func NewShipmentService(repo models.ShipmentRepository) models.ShipmentService {
	return &ShipmentServiceImpl{repo: repo}
}

func (s *ShipmentServiceImpl) Create(in models.CreateShipmentInput) (*models.Shipment, error) {
	if in.OrderID == uuid.Nil || in.UserID == uuid.Nil {
		return nil, errors.New("order_id and user_id are required")
	}
	if in.ShippingAddress == nil {
		in.ShippingAddress = models.Address{}
	}
	shipment := &models.Shipment{OrderID: in.OrderID, UserID: in.UserID, Carrier: in.Carrier, TrackingNumber: in.TrackingNumber, Status: "pending", ShippingAddress: in.ShippingAddress, EstimatedDeliveryAt: in.EstimatedDeliveryAt}
	if err := s.repo.Create(context.Background(), shipment); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *ShipmentServiceImpl) Get(id uuid.UUID) (*models.Shipment, error) {
	return s.repo.Get(context.Background(), id)
}
func (s *ShipmentServiceImpl) GetByOrder(orderID uuid.UUID) (*models.Shipment, error) {
	return s.repo.GetByOrder(context.Background(), orderID)
}
func (s *ShipmentServiceImpl) ListByUser(userID uuid.UUID, limit, page int) ([]models.Shipment, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	return s.repo.ListByUser(context.Background(), userID, int64(limit), int64((page-1)*limit))
}
func (s *ShipmentServiceImpl) UpdateStatus(id uuid.UUID, in models.UpdateShipmentStatusInput) (*models.Shipment, error) {
	if in.Status == "" {
		return nil, errors.New("status is required")
	}
	if err := s.repo.UpdateStatus(context.Background(), id, in); err != nil {
		return nil, err
	}
	return s.repo.Get(context.Background(), id)
}
