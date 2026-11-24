package events

import "context"

type Consumer interface {
	Subscribe(ctx context.Context, topic string, handler func(message []byte)) error
}
