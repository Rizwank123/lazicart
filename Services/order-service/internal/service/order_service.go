package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Rizwank123/LahziCart/order-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/events"
	kafkabroker "github.com/Rizwank123/LahziCart/shared/messaging/kafkaBroker"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type OrderServiceImpl struct {
	repo   models.OrderRepository
	kafka  *kafkabroker.Producer
	logger *zap.Logger
}

func NewOrderService(repo models.OrderRepository, kafka *kafkabroker.Producer, logger *zap.Logger) models.OrderService {
	return &OrderServiceImpl{repo: repo, kafka: kafka, logger: logger}
}

func (s *OrderServiceImpl) CreateOrder(in models.CreateOrderInput) (*models.Order, error) {
	if in.UserID == uuid.Nil || len(in.Items) == 0 {
		return nil, errors.New("user_id and items are required")
	}
	if in.Currency == "" {
		in.Currency = "INR"
	}
	order := &models.Order{UserID: in.UserID, Status: "pending", Currency: in.Currency, ShippingAddress: in.ShippingAddress}
	for _, item := range in.Items {
		if item.ProductID == uuid.Nil || item.Quantity <= 0 || item.UnitPrice < 0 {
			return nil, errors.New("invalid order item")
		}
		order.Total += float64(item.Quantity) * item.UnitPrice
		order.Items = append(order.Items, models.OrderItem{ProductID: item.ProductID, Quantity: item.Quantity, UnitPrice: item.UnitPrice, Currency: in.Currency})
	}
	if err := s.repo.CreateOrder(context.Background(), order); err != nil {
		return nil, err
	}
	s.publishOrderCreated(order)
	return order, nil
}

func (s *OrderServiceImpl) GetOrder(id uuid.UUID) (*models.Order, error) {
	return s.repo.GetOrder(context.Background(), id)
}

func (s *OrderServiceImpl) ListByUser(userID uuid.UUID, limit, page int) ([]models.Order, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	return s.repo.ListByUser(context.Background(), userID, int64(limit), int64((page-1)*limit))
}

func (s *OrderServiceImpl) CancelOrder(id uuid.UUID) error {
	order, err := s.repo.GetOrder(context.Background(), id)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateStatus(context.Background(), id, "cancelled"); err != nil {
		return err
	}
	s.publish(events.OrderCancelled, events.OrderCancelledEvent{
		EventType:   events.OrderCancelled,
		Version:     1,
		OrderID:     id.String(),
		UserID:      order.UserID.String(),
		CancelledAt: time.Now(),
	})
	return nil
}

func (s *OrderServiceImpl) UpdateStatus(id uuid.UUID, status string) error {
	return s.repo.UpdateStatus(context.Background(), id, status)
}

func (s *OrderServiceImpl) publishOrderCreated(order *models.Order) {
	items := make([]events.OrderItemEvent, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, events.OrderItemEvent{ProductID: item.ProductID.String(), Quantity: item.Quantity, UnitPrice: item.UnitPrice})
	}
	s.publish(events.OrderCreated, events.OrderCreatedEvent{
		EventType: events.OrderCreated,
		Version:   1,
		OrderID:   order.ID.String(),
		UserID:    order.UserID.String(),
		Items:     items,
		Total:     order.Total,
		CreatedAt: time.Now(),
	})
}

func (s *OrderServiceImpl) publish(topic string, payload interface{}) {
	if s.kafka == nil {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error("failed to marshal order event", zap.Error(err))
		return
	}
	if err := s.kafka.Publish(context.Background(), topic, data); err != nil {
		s.logger.Error("failed to publish order event", zap.String("topic", topic), zap.Error(err))
	}
}
