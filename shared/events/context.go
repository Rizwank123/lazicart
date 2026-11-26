package events

import (
	"context"
	"encoding/json"

	"github.com/gofrs/uuid/v5"
)

type contextKey string

const (
	correlationIDKey contextKey = "correlation_id"
	sagaIDKey        contextKey = "saga_id"
)

func WithCorrelationID(ctx context.Context, id *uuid.UUID) context.Context {
	if id == nil {
		return ctx
	}
	return context.WithValue(ctx, correlationIDKey, *id)
}

func WithSagaID(ctx context.Context, id *uuid.UUID) context.Context {
	if id == nil {
		return ctx
	}
	return context.WithValue(ctx, sagaIDKey, *id)
}

func CorrelationID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(correlationIDKey).(uuid.UUID)
	return id, ok
}

func SagaID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(sagaIDKey).(uuid.UUID)
	return id, ok
}

// ContextFromMessage extracts propagation metadata from a versioned event
// envelope. Legacy payloads simply return the original context.
func ContextFromMessage(ctx context.Context, message []byte) context.Context {
	var envelope struct {
		CorrelationID *uuid.UUID `json:"correlation_id"`
		SagaID        *uuid.UUID `json:"saga_id"`
	}
	if json.Unmarshal(message, &envelope) != nil {
		return ctx
	}
	ctx = WithCorrelationID(ctx, envelope.CorrelationID)
	return WithSagaID(ctx, envelope.SagaID)
}
