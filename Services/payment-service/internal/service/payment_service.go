package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Rizwank123/LahziCart/payment-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/events"
	kafkabroker "github.com/Rizwank123/LahziCart/shared/messaging/kafkaBroker"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type PaymentServiceImpl struct {
	repo   models.PaymentRepository
	kafka  *kafkabroker.Producer
	logger *zap.Logger
}

func NewPaymentService(repo models.PaymentRepository, kafka *kafkabroker.Producer, logger *zap.Logger) models.PaymentService {
	return &PaymentServiceImpl{repo: repo, kafka: kafka, logger: logger}
}

func (s *PaymentServiceImpl) CreatePayment(in models.CreatePaymentInput) (*models.Payment, error) {
	if in.OrderID == uuid.Nil || in.UserID == uuid.Nil || in.Amount <= 0 {
		return nil, errors.New("order_id, user_id and positive amount are required")
	}
	if in.Currency == "" {
		in.Currency = "INR"
	}
	payment := &models.Payment{OrderID: in.OrderID, UserID: in.UserID, Amount: in.Amount, Currency: in.Currency, Provider: in.Provider, Status: "pending"}
	if payment.Provider == "" {
		payment.Provider = "manual"
	}
	if err := s.repo.Create(context.Background(), payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *PaymentServiceImpl) GetPayment(id uuid.UUID) (*models.Payment, error) {
	return s.repo.Get(context.Background(), id)
}

func (s *PaymentServiceImpl) GetByOrderID(orderID uuid.UUID) (*models.Payment, error) {
	return s.repo.GetByOrderID(context.Background(), orderID)
}

func (s *PaymentServiceImpl) UpdateStatus(id uuid.UUID, in models.UpdatePaymentStatusInput) (*models.Payment, error) {
	if in.Status == "" {
		return nil, errors.New("status is required")
	}
	if err := s.repo.UpdateStatus(context.Background(), id, in.Status, in.TransactionID, in.FailureReason); err != nil {
		return nil, err
	}
	payment, err := s.repo.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}
	switch in.Status {
	case "succeeded":
		s.publish(events.PaymentSucceeded, payment)
	case "failed":
		s.publish(events.PaymentFailed, payment)
	}
	return payment, nil
}

func (s *PaymentServiceImpl) publish(topic string, payment *models.Payment) {
	if s.kafka == nil {
		return
	}
	tx := ""
	if payment.TransactionID != nil {
		tx = *payment.TransactionID
	}
	data, err := json.Marshal(events.PaymentEvent{
		EventType:     topic,
		Version:       1,
		PaymentID:     payment.ID.String(),
		OrderID:       payment.OrderID.String(),
		UserID:        payment.UserID.String(),
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		TransactionID: tx,
		OccurredAt:    time.Now(),
	})
	if err != nil {
		s.logger.Error("failed to marshal payment event", zap.Error(err))
		return
	}
	if err := s.kafka.Publish(context.Background(), topic, data); err != nil {
		s.logger.Error("failed to publish payment event", zap.String("topic", topic), zap.Error(err))
	}
}
