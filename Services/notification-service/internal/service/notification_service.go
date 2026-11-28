package service

import (
	"context"
	"errors"

	"github.com/Rizwank123/LahziCart/notification-service/internal/models"
	"github.com/gofrs/uuid/v5"
)

type NotificationServiceImpl struct{ repo models.NotificationRepository }

func NewNotificationService(repo models.NotificationRepository) models.NotificationService {
	return &NotificationServiceImpl{repo: repo}
}

func (s *NotificationServiceImpl) Create(in models.CreateNotificationInput) (*models.Notification, error) {
	if in.Body == "" || in.Type == "" {
		return nil, errors.New("type and body are required")
	}
	if in.Channel == "" {
		in.Channel = "email"
	}
	if in.Metadata == nil {
		in.Metadata = models.Metadata{}
	}
	n := &models.Notification{UserID: in.UserID, Channel: in.Channel, Type: in.Type, Recipient: in.Recipient, Subject: in.Subject, Body: in.Body, Status: "pending", Metadata: in.Metadata}
	if err := s.repo.Create(context.Background(), n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *NotificationServiceImpl) ListByUser(userID uuid.UUID, limit, page int) ([]models.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	return s.repo.ListByUser(context.Background(), userID, int64(limit), int64((page-1)*limit))
}

func (s *NotificationServiceImpl) MarkSent(id uuid.UUID) error {
	return s.repo.MarkSent(context.Background(), id)
}

func (s *NotificationServiceImpl) Get(id uuid.UUID) (*models.Notification, error) {
	return s.repo.Get(context.Background(), id)
}
