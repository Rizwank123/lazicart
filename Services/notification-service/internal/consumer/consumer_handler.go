package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/Rizwank123/LahziCart/notification-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/events"
	"github.com/gofrs/uuid/v5"
)

type Handler struct{ service models.NotificationService }

func NewHandler(service models.NotificationService) *Handler { return &Handler{service: service} }

func (h *Handler) OrderCreated() func([]byte) error {
	return func(message []byte) error {
		var event events.OrderCreatedEvent
		if err := json.Unmarshal(message, &event); err != nil {
			return err
		}
		userID, _ := uuid.FromString(event.UserID)
		_, err := h.service.Create(models.CreateNotificationInput{
			UserID: &userID, Channel: "email", Type: events.OrderCreated,
			Body:     fmt.Sprintf("Your order %s was created successfully.", event.OrderID),
			Metadata: models.Metadata{"order_id": event.OrderID},
		})
		return err
	}
}

func (h *Handler) PaymentEvent(topic string) func([]byte) error {
	return func(message []byte) error {
		var event events.PaymentEvent
		if err := json.Unmarshal(message, &event); err != nil {
			return err
		}
		userID, _ := uuid.FromString(event.UserID)
		_, err := h.service.Create(models.CreateNotificationInput{
			UserID: &userID, Channel: "email", Type: topic,
			Body:     fmt.Sprintf("Payment for order %s is %s.", event.OrderID, topic),
			Metadata: models.Metadata{"order_id": event.OrderID, "payment_id": event.PaymentID},
		})
		return err
	}
}
