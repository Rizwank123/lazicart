package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
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
	for topic, handle := range map[string]func([]byte) error{
		events.OrderCreated:     consumerHandler.OrderCreated(),
		events.PaymentSucceeded: consumerHandler.PaymentEvent(events.PaymentSucceeded),
		events.PaymentFailed:    consumerHandler.PaymentEvent(events.PaymentFailed),
	} {
		c := kafkabroker.NewConsumer([]string{"localhost:9092"}, "notification-service-group-"+topic, kafkaProducer, topic, logger)
		go c.Consume(context.Background(), handle)
	}
	r := router.NewNotificationRouter(handler.NewNotificationHandler(svc))
	e := echo.New()
	e.Use(sharedmw.NewRecoveryMiddleware(logger).Recover)
	e.Use(sharedmw.NewRequestIDMiddleware("X-Request-ID").Inject)
	e.Use(sharedmw.NewLoggingMiddleware(logger, "/health").Log)
	e.Use(sharedmw.NewMetricsMiddleware("notification-service").Collect)
	e.Use(sharedmw.NewCORSMiddleware(sharedmw.DefaultCORSConfig()))
	r.RegisterRoutes(e)
	e.GET("/health", func(c echo.Context) error { return c.String(http.StatusOK, "OK") })
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
	_ = e.Shutdown(ctx)
}
