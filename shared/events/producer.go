package events

import "context"

type Producer interface {
	Publish(ctx context.Context, topic string, message []byte) error
}
