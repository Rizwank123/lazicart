package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Rizwank123/LahziCart/shared/events"
	kafkabroker "github.com/Rizwank123/LahziCart/shared/messaging/kafkaBroker"
	"github.com/Rizwank123/lahziCart/user-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type UserServiceImpl struct {
	ur     models.UserRepository
	kafka  *kafkabroker.Producer
	logger *zap.Logger
}

func NewUserService(ur models.UserRepository, kafka *kafkabroker.Producer, logger *zap.Logger) models.UserService {
	return &UserServiceImpl{ur: ur, kafka: kafka, logger: logger}
}

// CreateUser implements [models.UserService].
func (u *UserServiceImpl) CreateUser(in models.CreateUserRequest) (*models.User, error) {
	user := models.User{

		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Phone:     in.Phone,
		IsActive:  true,
		Address:   in.Address,
		AuthId:    in.AuthId,
		Dob:       in.Dob,
		Gender:    in.Gender,
	}
	err := u.ur.CreateUser(context.Background(), &user)
	if err != nil {
		u.logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (u *UserServiceImpl) CreateUserFromEvent(eventID uuid.UUID, in models.CreateUserRequest) error {
	if eventID == uuid.Nil {
		return errors.New("event_id is required")
	}
	if err := u.ur.ProcessUserCreated(context.Background(), eventID, in); err != nil {
		u.logger.Error("failed to process user.created event", zap.Error(err), zap.String("event_id", eventID.String()))
		return err
	}
	return nil
}

// DeleteUser implements [models.UserService].
func (u *UserServiceImpl) DeleteUser(id uuid.UUID) (*models.User, error) {
	user, err := u.ur.GetUserByID(context.Background(), id)
	if err != nil {
		u.logger.Error("failed to get user", zap.Error(err))
		return nil, err
	}
	// The profile is marked as deleted first. The event is published
	// afterwards, so a broker problem cannot leave the profile intact while
	// the caller is told that the delete failed.
	if err := u.ur.DeleteUser(context.Background(), user); err != nil {
		u.logger.Error("failed to delete user", zap.Error(err))
		return nil, err
	}
	event := events.UserDeletedEvent{
		EventType: events.UserDeleted,
		Version:   1,
		UserID:    user.ID.String(),
		AuthId:    user.AuthId,
	}
	eventByte, err := json.Marshal(event)
	if err != nil {
		u.logger.Error("failed to marshal user deleted event", zap.Error(err))
		return user, nil
	}
	if err := u.kafka.Publish(context.Background(), events.UserDeleted, eventByte); err != nil {
		// The delete already happened, so the event failure is logged
		// instead of being reported to the caller.
		u.logger.Error("failed to publish user deleted event", zap.Error(err))
	}
	return user, nil
}

// GetUserByEmail implements [models.UserService].
func (u *UserServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	user, err := u.ur.GetUserByEmail(context.Background(), email)
	if err != nil {
		u.logger.Error("failed to get user", zap.Error(err))
		return nil, err
	}
	return user, nil
}

// GetUserByID implements [models.UserService].
func (u *UserServiceImpl) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := u.ur.GetUserByID(context.Background(), id)
	if err != nil {
		u.logger.Error("failed to get user", zap.Error(err))
		return nil, err
	}
	return user, nil
}

// GetUserByPhone implements [models.UserService].
func (u *UserServiceImpl) GetUserByPhone(phone string) (*models.User, error) {
	user, err := u.ur.GetUserByPhone(context.Background(), phone)
	if err != nil {
		u.logger.Error("failed to get user", zap.Error(err))
		return nil, err
	}
	return user, nil
}

// UpdateUser implements [models.UserService].
func (u *UserServiceImpl) UpdateUser(id uuid.UUID, in models.UpdateUser) (*models.User, error) {
	user, err := u.ur.GetUserByID(context.Background(), id)
	if err != nil {
		u.logger.Error("failed to get user", zap.Error(err))
		return nil, err
	}
	if in.FirstName != "" {
		user.FirstName = in.FirstName
	}
	if in.LastName != "" {
		user.LastName = in.LastName
	}
	if in.Email != "" {
		user.Email = in.Email
	}
	if in.Phone != "" {
		user.Phone = in.Phone
	}
	if in.Address != nil {
		user.Address = *in.Address
	}
	if !in.Dob.IsZero() {
		user.Dob = in.Dob
	}
	if in.Gender != "" {
		user.Gender = in.Gender
	}
	err = u.ur.UpdateUser(context.Background(), user)
	if err != nil {
		u.logger.Error("failed to update user", zap.Error(err))
		return nil, err
	}
	event := events.UserUpdatedEvent{
		EventType: events.UserUpdated,
		Version:   1,
		AuthId:    user.AuthId,
		Email:     user.Email,
		Phone:     user.Phone,
	}
	eventByte, err := json.Marshal(event)
	if err != nil {
		u.logger.Error("failed to marshal user updated event", zap.Error(err))
		return user, nil
	}
	if err := u.kafka.Publish(context.Background(), events.UserUpdated, eventByte); err != nil {
		// The profile was already updated, so the event failure is logged
		// instead of being reported to the caller. Reporting it would tell
		// the client that the update failed even though it was stored.
		u.logger.Error("failed to publish user updated event", zap.Error(err))
	}
	return user, nil
}
