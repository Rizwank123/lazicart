package consumers

import (
	"encoding/json"
	"log"

	"github.com/Rizwank123/LahziCart/auth-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/events"
)

type UserConsumerHandler struct {
	as models.UserAuthService
}

func NewUserConsumerHandler(as models.UserAuthService) *UserConsumerHandler {
	return &UserConsumerHandler{as: as}
}

func (h *UserConsumerHandler) DeleteUserHandler() func([]byte) error {
	return func(message []byte) error {
		var req events.UserDeletedEvent
		if err := json.Unmarshal(message, &req); err != nil {
			return err
		}
		err := h.as.DeleteUserAuth(req.AuthId)
		if err != nil {
			return err
		}
		log.Println("User deleted successfully")
		return nil
	}
}

func (h *UserConsumerHandler) UpdateUserHandler() func([]byte) error {
	return func(message []byte) error {
		var req events.UserUpdatedEvent
		if err := json.Unmarshal(message, &req); err != nil {
			return err
		}
		err := h.as.UpdateUserAuth(req.AuthId, req)
		if err != nil {
			return err
		}
		log.Println("User updated successfully")
		return nil
	}
}
