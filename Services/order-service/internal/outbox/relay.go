package outbox

import (
	"context"
	"time"

	"github.com/Rizwank123/LahziCart/order-service/internal/models"
	kafkabroker "github.com/Rizwank123/LahziCart/shared/messaging/kafkaBroker"
	"go.uber.org/zap"
)

type Repository interface {
	GetPendingOutboxEvents(ctx context.Context, limit int) ([]models.OutboxEvent, error)
}

type Relay struct {
	repo     models.OrderRepository
	producer *kafkabroker.Producer
	logger   *zap.Logger
}

func NewRelay(repo models.OrderRepository, producer *kafkabroker.Producer, logger *zap.Logger) *Relay {
	return &Relay{repo: repo, producer: producer, logger: logger}
}

func (r *Relay) Run(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.publishBatch(ctx)
		}
	}
}

func (r *Relay) publishBatch(ctx context.Context) {
	events, err := r.repo.GetPendingOutboxEvents(ctx, 100)
	if err != nil {
		r.logger.Error("failed to load order outbox events", zap.Error(err))
		return
	}

	for _, event := range events {
		if err := r.producer.Publish(ctx, event.EventType, event.Payload); err != nil {
			r.logger.Error("failed to publish order outbox event",
				zap.String("event_type", event.EventType),
				zap.String("event_id", event.EventID.String()),
				zap.Error(err))
			if markErr := r.repo.MarkOutboxFailed(ctx, event.ID, err.Error()); markErr != nil {
				r.logger.Error("failed to record outbox failure", zap.Error(markErr))
			}
			continue
		}
		if err := r.repo.MarkOutboxPublished(ctx, event.ID); err != nil {
			r.logger.Error("failed to mark outbox event published", zap.Error(err))
		}
	}
}
