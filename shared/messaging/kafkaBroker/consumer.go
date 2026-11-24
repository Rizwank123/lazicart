package kafkabroker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Rizwank123/LahziCart/shared/observability"
	"github.com/Rizwank123/LahziCart/shared/utils"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader     *kafka.Reader
	producer   *Producer
	topic      string
	logger     *zap.Logger
	maxRetries int
}

func NewConsumer(
	brokers []string,
	groupID string,
	producer *Producer,
	topic string,
	logger *zap.Logger,
) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		producer:   producer,
		topic:      topic,
		logger:     logger,
		maxRetries: 5,
	}
}

func (c *Consumer) Consume(
	ctx context.Context,
	handler func(message []byte) error,
) {
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.logger.Error("failed to read message", zap.Error(err))
			continue
		}

		start := time.Now()

		err = c.processWithRetry(ctx, msg.Value, handler)

		duration := time.Since(start).Seconds()
		observability.KafkaProcessingDration.
			WithLabelValues(c.topic).
			Observe(duration)

		if err != nil {
			observability.KafkaMessagesConsumed.
				WithLabelValues(c.topic, "failed").
				Inc()
			c.logger.Error("message processing failed", zap.Error(err))
		} else {
			observability.KafkaMessagesConsumed.
				WithLabelValues(c.topic, "success").
				Inc()
		}
	}
}

func (c *Consumer) processWithRetry(
	ctx context.Context,
	message []byte,
	handler func(message []byte) error,
) error {
	var err error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		err = handler(message)
		if err == nil {
			return nil
		}

		delay := utils.ExponentialBackoff(attempt)
		c.logger.Warn(
			"retrying message",
			zap.Int("attempt", attempt),
			zap.Duration("delay", delay),
			zap.Error(err),
		)

		time.Sleep(delay)
	}

	return c.SendToDLQ(ctx, message, err)
}

func (c *Consumer) SendToDLQ(
	ctx context.Context,
	message []byte,
	cause error,
) error {
	dlqTopic := c.topic + ".dlq"

	payload := map[string]interface{}{
		"original_topic": c.topic,
		"message":        json.RawMessage(message),
		"error":          cause.Error(),
		"failed_at":      time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	c.logger.Error(
		"sending message to DLQ",
		zap.String("topic", dlqTopic),
		zap.Error(cause),
	)

	return c.producer.Publish(ctx, dlqTopic, data)
}
