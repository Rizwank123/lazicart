package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Rizwank123/LahziCart/order-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/events"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type OrderServiceImpl struct {
	repo   models.OrderRepository
	logger *zap.Logger
}

func NewOrderService(repo models.OrderRepository, logger *zap.Logger) models.OrderService {
	return &OrderServiceImpl{repo: repo, logger: logger}
}

func (s *OrderServiceImpl) CreateOrder(in models.CreateOrderInput, idempotencyKey string) (*models.Order, error) {
	if in.UserID == uuid.Nil || len(in.Items) == 0 {
		return nil, errors.New("user_id and items are required")
	}
	if in.Currency == "" {
		in.Currency = "INR"
	}
	requestData, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	requestHash := ""
	if idempotencyKey != "" {
		requestHash = fmt.Sprintf("%x", sha256.Sum256(requestData))
	}
	orderID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	order := &models.Order{UserID: in.UserID, Status: "pending", Currency: in.Currency, ShippingAddress: in.ShippingAddress}
	order.ID = orderID
	for _, item := range in.Items {
		if item.ProductID == uuid.Nil || item.Quantity <= 0 || item.UnitPrice < 0 {
			return nil, errors.New("invalid order item")
		}
		order.Total += float64(item.Quantity) * item.UnitPrice
		order.Items = append(order.Items, models.OrderItem{ProductID: item.ProductID, Quantity: item.Quantity, UnitPrice: item.UnitPrice, Currency: in.Currency})
	}
	event, err := newOutboxEvent(events.OrderCreated, order, s.orderCreatedPayload(order), in.CorrelationID, in.SagaID)
	if err != nil {
		return nil, err
	}
	created, err := s.repo.CreateOrder(context.Background(), order, event, idempotencyKey, requestHash)
	if err != nil {
		return nil, err
	}
	return created, nil
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
	event, err := newOutboxEvent(events.OrderCancelled, order, events.OrderCancelledEvent{
		EventType:   events.OrderCancelled,
		Version:     1,
		OrderID:     id.String(),
		UserID:      order.UserID.String(),
		CancelledAt: time.Now().UTC(),
	}, nil, nil)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateStatus(context.Background(), id, "cancelled", event); err != nil {
		return err
	}
	return nil
}

func (s *OrderServiceImpl) UpdateStatus(id uuid.UUID, status string) error {
	return s.repo.UpdateStatus(context.Background(), id, status, nil)
}

func (s *OrderServiceImpl) orderCreatedPayload(order *models.Order) events.OrderCreatedEvent {
	items := make([]events.OrderItemEvent, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, events.OrderItemEvent{ProductID: item.ProductID.String(), Quantity: item.Quantity, UnitPrice: item.UnitPrice})
	}
	return events.OrderCreatedEvent{
		EventType: events.OrderCreated,
		Version:   1,
		OrderID:   order.ID.String(),
		UserID:    order.UserID.String(),
		Items:     items,
		Total:     order.Total,
		CreatedAt: time.Now().UTC(),
	}
}

func newOutboxEvent(eventType string, order *models.Order, payload interface{}, correlationID, sagaID *uuid.UUID) (*models.OutboxEvent, error) {
	eventID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	envelope := events.NewEnvelope(eventType, "order", order.ID, "order-service", payload)
	envelope.EventID = eventID
	envelope.CorrelationID = correlationID
	envelope.SagaID = sagaID
	envelopeData, err := json.Marshal(envelope)
	if err != nil {
		return nil, err
	}
	return &models.OutboxEvent{
		ID:            uuid.Must(uuid.NewV4()),
		EventID:       eventID,
		AggregateType: "order",
		AggregateID:   order.ID,
		EventType:     eventType,
		EventVersion:  1,
		Payload:       envelopeData,
		Status:        "pending",
		AvailableAt:   time.Now().UTC(),
	}, nil
}
