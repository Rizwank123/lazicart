package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rizwank123/LahziCart/payment-service/internal/config"
	"github.com/Rizwank123/LahziCart/payment-service/internal/db"
	"github.com/Rizwank123/LahziCart/payment-service/internal/handler"
	"github.com/Rizwank123/LahziCart/payment-service/internal/repository"
	"github.com/Rizwank123/LahziCart/payment-service/internal/router"
	"github.com/Rizwank123/LahziCart/payment-service/internal/service"
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
	repo := repository.NewPaymentRepository(pool)
	svc := service.NewPaymentService(repo, kafkaProducer, logger)
	r := router.NewPaymentRouter(handler.NewPaymentHandler(svc))
	e := echo.New()
	e.Use(sharedmw.NewRecoveryMiddleware(logger).Recover)
	e.Use(sharedmw.NewRequestIDMiddleware("X-Request-ID").Inject)
	e.Use(sharedmw.NewLoggingMiddleware(logger, "/health").Log)
	e.Use(sharedmw.NewMetricsMiddleware("payment-service").Collect)
	e.Use(sharedmw.NewCORSMiddleware(sharedmw.DefaultCORSConfig()))
	r.RegisterRoutes(e)
	e.GET("/health", func(c echo.Context) error { return c.String(http.StatusOK, "OK") })
	go func() {
		logger.Info("payment-service started", zap.String("port", cfg.AppPort))
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
