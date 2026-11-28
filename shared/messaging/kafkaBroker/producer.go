package kafkabroker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Rizwank123/LahziCart/shared/config"
	"github.com/Rizwank123/LahziCart/shared/events"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	// publishAttempts is the number of times a message is written before the
	// publish is reported as failed.
	publishAttempts = 3
	// publishBackoff is the base delay between two attempts. It grows with
	// the attempt number.
	publishBackoff = 250 * time.Millisecond
)

type Producer struct {
	writer  *kafka.Writer
	Logger  *zap.Logger
	Brokers []string

	ensureMu sync.Mutex
	ensured  map[string]bool
}

func NewProducer(brokers []string, logger *zap.Logger) *Producer {
	return &Producer{
		writer: kafka.NewWriter(
			kafka.WriterConfig{
				Brokers: brokers,
			},
		),
		Logger:  logger,
		Brokers: brokers,
		ensured: map[string]bool{},
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, message []byte) error {
	tracer := otel.Tracer("shared.kafka")
	ctx, span := tracer.Start(ctx, "kafka.publish "+topic, trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	carrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, carrier)
	headers := make([]kafka.Header, 0, len(carrier))
	for key, value := range carrier {
		headers = append(headers, kafka.Header{Key: key, Value: []byte(value)})
	}

	// The broker does not create topics on demand, so a publish to a topic
	// that does not exist yet would fail.
	if err := p.ensureTopic(topic); err != nil {
		p.Logger.Error("failed to ensure topic", zap.String("topic", topic), zap.Error(err))
		return err
	}

	msg := kafka.Message{
		Topic:   topic,
		Value:   message,
		Headers: headers,
	}

	var err error
	for attempt := 1; attempt <= publishAttempts; attempt++ {
		err = p.writer.WriteMessages(ctx, msg)
		if err == nil {
			p.Logger.Info("message published successfully", zap.String("topic", topic), zap.ByteString("message", message))
			return nil
		}
		p.Logger.Warn(
			"failed to publish message",
			zap.String("topic", topic),
			zap.Int("attempt", attempt),
			zap.Error(err),
		)
		if attempt == publishAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(publishBackoff * time.Duration(attempt)):
		}
	}
	return err
}

// ensureTopic creates the topic once per producer instance.
func (p *Producer) ensureTopic(topic string) error {
	p.ensureMu.Lock()
	defer p.ensureMu.Unlock()
	if p.ensured[topic] {
		return nil
	}
	if err := events.EnsureTopic(p.Brokers, topic, 1, 1); err != nil {
		return err
	}
	p.ensured[topic] = true
	return nil
}

func (p *Producer) Close() {
	if p.writer != nil {
		_ = p.writer.Close()
	}
}

func InitKafkaProducer() *Producer {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
		return nil
	}
	log.Println("config loaded successfully", zap.Any("config", cfg))
	logger, _ := zap.NewProduction()
	return NewProducer(cfg.KafkaBrokers, logger)
}
