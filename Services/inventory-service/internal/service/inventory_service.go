package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Rizwank123/LahziCart/inventory-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/events"
	kafkabroker "github.com/Rizwank123/LahziCart/shared/messaging/kafkaBroker"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type InventoryServiceImpl struct {
	repo   models.InventoryRepository
	kafka  *kafkabroker.Producer
	logger *zap.Logger
}

func NewInventoryService(repo models.InventoryRepository, kafka *kafkabroker.Producer, logger *zap.Logger) models.InventoryService {
	return &InventoryServiceImpl{repo: repo, kafka: kafka, logger: logger}
}

func (s *InventoryServiceImpl) UpsertItem(in models.UpsertInventoryInput) (*models.InventoryItem, error) {
	if in.ProductID == uuid.Nil {
		return nil, errors.New("product_id is required")
	}
	if in.Quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}
	item := &models.InventoryItem{
		ProductID:         in.ProductID,
		SKU:               in.SKU,
		QuantityAvailable: in.Quantity,
		ReorderLevel:      in.ReorderLevel,
		IsActive:          in.IsActive,
	}
	if err := s.repo.UpsertItem(context.Background(), item); err != nil {
		return nil, err
	}
	s.publish(events.InventoryAdjusted, map[string]interface{}{
		"event_type": events.InventoryAdjusted,
		"version":    1,
		"product_id": item.ProductID,
		"quantity":   item.QuantityAvailable,
		"updated_at": time.Now(),
	})
	return item, nil
}

func (s *InventoryServiceImpl) GetByProductID(productID uuid.UUID) (*models.InventoryItem, error) {
	return s.repo.GetByProductID(context.Background(), productID)
}

func (s *InventoryServiceImpl) ListItems(limit, page int) ([]models.InventoryItem, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	total, items, err := s.repo.ListItems(context.Background(), int64(limit), int64((page-1)*limit))
	return items, total, err
}

func (s *InventoryServiceImpl) Adjust(in models.AdjustInventoryInput) (*models.InventoryItem, error) {
	if in.ProductID == uuid.Nil {
		return nil, errors.New("product_id is required")
	}
	item, err := s.repo.Adjust(context.Background(), in.ProductID, in.Delta)
	if err != nil {
		return nil, err
	}
	s.publish(events.InventoryAdjusted, map[string]interface{}{
		"event_type": events.InventoryAdjusted,
		"version":    1,
		"product_id": item.ProductID,
		"delta":      in.Delta,
		"updated_at": time.Now(),
	})
	return item, nil
}

func (s *InventoryServiceImpl) Reserve(in models.ReserveInventoryInput) ([]models.InventoryReservation, error) {
	if in.OrderID == uuid.Nil {
		return nil, errors.New("order_id is required")
	}
	reservations, err := s.repo.Reserve(context.Background(), in.OrderID, in.Items)
	if err != nil {
		return nil, err
	}
	eventItems := make([]events.OrderItemEvent, 0, len(in.Items))
	for _, item := range in.Items {
		eventItems = append(eventItems, events.OrderItemEvent{ProductID: item.ProductID.String(), Quantity: item.Quantity})
	}
	s.publish(events.InventoryReserved, events.InventoryReservedEvent{
		EventType:  events.InventoryReserved,
		Version:    1,
		OrderID:    in.OrderID.String(),
		Items:      eventItems,
		ReservedAt: time.Now(),
	})
	return reservations, nil
}

func (s *InventoryServiceImpl) Release(orderID uuid.UUID) error {
	if err := s.repo.Release(context.Background(), orderID); err != nil {
		return err
	}
	s.publish(events.InventoryReleased, map[string]interface{}{
		"event_type":  events.InventoryReleased,
		"version":     1,
		"order_id":    orderID.String(),
		"released_at": time.Now(),
	})
	return nil
}

func (s *InventoryServiceImpl) Commit(orderID uuid.UUID) error {
	return s.repo.Commit(context.Background(), orderID)
}

func (s *InventoryServiceImpl) publish(topic string, payload interface{}) {
	if s.kafka == nil {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error("failed to marshal inventory event", zap.Error(err))
		return
	}
	if err := s.kafka.Publish(context.Background(), topic, data); err != nil {
		s.logger.Error("failed to publish inventory event", zap.String("topic", topic), zap.Error(err))
	}
}
