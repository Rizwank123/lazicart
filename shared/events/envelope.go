package events

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Envelope struct {
	EventID       uuid.UUID   `json:"event_id"`
	EventType     string      `json:"event_type"`
	Version       int         `json:"version"`
	AggregateType string      `json:"aggregate_type"`
	AggregateID   uuid.UUID   `json:"aggregate_id"`
	SagaID        *uuid.UUID  `json:"saga_id,omitempty"`
	CorrelationID *uuid.UUID  `json:"correlation_id,omitempty"`
	Producer      string      `json:"producer"`
	OccurredAt    time.Time   `json:"occurred_at"`
	Payload       interface{} `json:"payload"`
}

func NewEnvelope(eventType, aggregateType string, aggregateID uuid.UUID, producer string, payload interface{}) Envelope {
	return Envelope{
		EventID:       uuid.Must(uuid.NewV4()),
		EventType:     eventType,
		Version:       1,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Producer:      producer,
		OccurredAt:    time.Now().UTC(),
		Payload:       payload,
	}
}
