package kafkabroker

import (
	"context"
	"log"

	"github.com/Rizwank123/LahziCart/shared/config"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer  *kafka.Writer
	Logger  *zap.Logger
	Brokers []string
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
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, message []byte) error {
	// ensure topic exists
	// err := events.EnsureTopic(p.Brokers, topic, 1, 1)
	// if err != nil {
	// 	p.Logger.Error("failed to ensure topic", zap.Error(err))
	// 	return err
	// }
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: message,
	})
	if err != nil {
		p.Logger.Error("failed to publish message", zap.Error(err))
		return err
	}
	p.Logger.Info("message published successfully", zap.String("topic", topic), zap.ByteString("message", message))
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
