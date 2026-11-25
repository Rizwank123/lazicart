package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rizwank123/LahziCart/shared/events"
	kafkabroker "github.com/Rizwank123/LahziCart/shared/messaging/kafkaBroker"
	"github.com/Rizwank123/lahziCart/user-service/internal/config"
	"github.com/Rizwank123/lahziCart/user-service/internal/consumer"
	"github.com/Rizwank123/lahziCart/user-service/internal/db"
	"github.com/Rizwank123/lahziCart/user-service/internal/handler"
	"github.com/Rizwank123/lahziCart/user-service/internal/repository"
	"github.com/Rizwank123/lahziCart/user-service/internal/router"
	"github.com/Rizwank123/lahziCart/user-service/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	db, err := db.NewDatabase(config)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		logger.Error("failed to ping database", zap.Error(err))
		panic(err)
	}
	kafkaProducer := kafkabroker.InitKafkaProducer()
	if kafkaProducer == nil {
		logger.Fatal("failed to init kafka producer")
	}
	userRepo := repository.NewPgxUserRepository(db)
	userService := service.NewUserService(userRepo, kafkaProducer, logger)
	userController := handler.NewUserController(userService)
	userRouter := router.NewUserRouter(userController)
	e := echo.New()
	userRouter.RegisterRoutes(e)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.GET("/ready", func(c echo.Context) error {
		checkContext, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(checkContext); err != nil {
			return c.String(http.StatusServiceUnavailable, "NOT_READY")
		}
		return c.String(http.StatusOK, "READY")
	})
	KafkaConsumer := kafkabroker.NewConsumer(
		config.KafkaBrokers,
		"user-service-group",
		kafkaProducer, // DLQ producer
		events.UserCreated,
		logger,
	)
	userConsumerHandler := consumer.NewUserConsumerHandler(userService)
	consumerContext, cancelConsumer := context.WithCancel(context.Background())
	consumerDone := make(chan struct{})
	go func() {
		defer close(consumerDone)
		KafkaConsumer.Consume(consumerContext, userConsumerHandler.CreateUserHandler())
	}()

	go func() {
		if err := e.Start(":" + config.AppPort); err != nil && err != http.ErrServerClosed {
			logger.Fatal("shutting down the server", zap.Error(err))
		}
	}()

	// 2. Wait for interrupt signal to gracefully shut down the server
	// Use a buffered channel as recommended by the signal package
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Block here until a signal is received

	logger.Info("Shutting down user-service...")
	cancelConsumer()
	KafkaConsumer.Close()
	select {
	case <-consumerDone:
	case <-time.After(5 * time.Second):
		logger.Warn("timed out waiting for Kafka consumer shutdown")
	}

	// 3. Create a timeout context for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Shutdown Echo (stops accepting new requests)
	if err := e.Shutdown(ctx); err != nil {
		logger.Error("Echo shutdown error", zap.Error(err))
	}

	// 6. Close DB connection
	db.Close()

	logger.Info("User service exited gracefully")
}
