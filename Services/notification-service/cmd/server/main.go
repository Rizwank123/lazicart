package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Rizwank123/LahziCart/notification-service/internal/config"
	"github.com/Rizwank123/LahziCart/notification-service/internal/consumer"
	"github.com/Rizwank123/LahziCart/notification-service/internal/db"
	"github.com/Rizwank123/LahziCart/notification-service/internal/handler"
	"github.com/Rizwank123/LahziCart/notification-service/internal/repository"
	"github.com/Rizwank123/LahziCart/notification-service/internal/router"
	"github.com/Rizwank123/LahziCart/notification-service/internal/service"
	"github.com/Rizwank123/LahziCart/shared/events"
	kafkabroker "github.com/Rizwank123/LahziCart/shared/messaging/kafkaBroker"
	sharedmw "github.com/Rizwank123/LahziCart/shared/middleware"
	"github.com/Rizwank123/LahziCart/shared/observability"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}
	shutdownTracer, err := observability.InitTracer(context.Background(), "notification-service")
	if err != nil {
		logger.Fatal("failed to initialize tracing", zap.Error(err))
	}
	defer shutdownTracer(context.Background())
	pool, err := db.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}
	defer pool.Close()
	kafkaProducer := kafkabroker.InitKafkaProducer()
	defer func() {
		if kafkaProducer != nil {
			kafkaProducer.Close()
		}
	}()
	repo := repository.NewNotificationRepository(pool)
	svc := service.NewNotificationService(repo)
	consumerHandler := consumer.NewHandler(svc)
	consumerContext, cancelConsumers := context.WithCancel(context.Background())
	var consumerWG sync.WaitGroup
	var consumers []*kafkabroker.Consumer
	for topic, handle := range map[string]func([]byte) error{
		events.OrderCreated:     consumerHandler.OrderCreated(),
		events.PaymentSucceeded: consumerHandler.PaymentEvent(events.PaymentSucceeded),
		events.PaymentFailed:    consumerHandler.PaymentEvent(events.PaymentFailed),
	} {
		c := kafkabroker.NewConsumer(cfg.KafkaBrokers, "notification-service-group-"+topic, kafkaProducer, topic, logger)
		consumers = append(consumers, c)
		consumerWG.Add(1)
		go func() {
			defer consumerWG.Done()
			c.ConsumeWithContext(consumerContext, func(messageContext context.Context, message []byte) error {
				if correlationID, ok := events.CorrelationID(messageContext); ok {
					logger.Debug("processing notification event", zap.String("topic", topic), zap.String("correlation_id", correlationID.String()))
				}
				if sagaID, ok := events.SagaID(messageContext); ok {
					logger.Debug("notification event saga", zap.String("topic", topic), zap.String("saga_id", sagaID.String()))
				}
				return handle(message)
			})
		}()
	}
	r := router.NewNotificationRouter(handler.NewNotificationHandler(svc))
	e := echo.New()
	e.Use(observability.TracingMiddleware("notification-service"))
	e.Use(sharedmw.NewRecoveryMiddleware(logger).Recover)
	e.Use(sharedmw.NewRequestIDMiddleware("X-Request-ID").Inject)
	e.Use(sharedmw.NewLoggingMiddleware(logger, "/health").Log)
	e.Use(sharedmw.NewMetricsMiddleware("notification-service").Collect)
	e.Use(sharedmw.NewCORSMiddleware(sharedmw.DefaultCORSConfig()))
	r.RegisterRoutes(e)
	e.GET("/health", func(c echo.Context) error { return c.String(http.StatusOK, "OK") })
	e.GET("/ready", func(c echo.Context) error {
		checkContext, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(checkContext); err != nil {
			return c.String(http.StatusServiceUnavailable, "NOT_READY")
		}
		return c.String(http.StatusOK, "READY")
	})
	go func() {
		logger.Info("notification-service started", zap.String("port", cfg.AppPort))
		if err := e.Start(":" + cfg.AppPort); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cancelConsumers()
	for _, c := range consumers {
		_ = c.Close()
	}
	consumerDone := make(chan struct{})
	go func() {
		consumerWG.Wait()
		close(consumerDone)
	}()
	select {
	case <-consumerDone:
	case <-time.After(5 * time.Second):
		logger.Warn("timed out waiting for Kafka consumers shutdown")
	}
	_ = e.Shutdown(ctx)
}
