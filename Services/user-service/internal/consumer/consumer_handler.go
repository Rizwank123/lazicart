package consumer

import (
	"encoding/json"
	"log"

	"github.com/Rizwank123/lahziCart/user-service/internal/models"
)

type UserConsumerHandler struct {
	us models.UserService
}

func NewUserConsumerHandler(us models.UserService) *UserConsumerHandler {
	return &UserConsumerHandler{us: us}
}

func (h *UserConsumerHandler) CreateUserHandler() func([]byte) error {
	return func(message []byte) error {
		var req models.CreateUserRequest

		if err := json.Unmarshal(message, &req); err != nil {
			return err
		}
		log.Println(req)
		_, err := h.us.CreateUser(req)
		return err
	}
}
