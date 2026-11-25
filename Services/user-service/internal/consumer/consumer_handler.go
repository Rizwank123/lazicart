package consumer

import (
	"encoding/json"

	"github.com/Rizwank123/LahziCart/shared/events"
	"github.com/Rizwank123/lahziCart/user-service/internal/models"
	"github.com/gofrs/uuid/v5"
)

type UserConsumerHandler struct {
	us models.UserService
}

func NewUserConsumerHandler(us models.UserService) *UserConsumerHandler {
	return &UserConsumerHandler{us: us}
}

func (h *UserConsumerHandler) CreateUserHandler() func([]byte) error {
	return func(message []byte) error {
		var envelope struct {
			EventID   uuid.UUID       `json:"event_id"`
			EventType string          `json:"event_type"`
			Payload   json.RawMessage `json:"payload"`
		}
		if err := json.Unmarshal(message, &envelope); err != nil {
			return err
		}
		if envelope.EventType != events.UserCreated || envelope.EventID == uuid.Nil {
			return &invalidEventError{message: "invalid user.created event envelope"}
		}
		var event events.UserCreatedEvent
		if err := json.Unmarshal(envelope.Payload, &event); err != nil {
			return err
		}
		return h.us.CreateUserFromEvent(envelope.EventID, models.CreateUserRequest{
			FirstName: event.FirstName,
			LastName:  event.LastName,
			Email:     event.Email,
			Phone:     event.Phone,
			AuthId:    event.AuthId,
			Dob:       event.Dob,
			Address: models.Address{
				City:        event.Address.City,
				State:       event.Address.State,
				Country:     event.Address.Country,
				Pincode:     event.Address.Pincode,
				Coordinates: event.Address.Coordinates,
			},
			Gender: event.Gender,
		})
	}
}

type invalidEventError struct{ message string }

func (e *invalidEventError) Error() string { return e.message }
