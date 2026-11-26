package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rizwank123/LahziCart/delivery-service/internal/config"
	"github.com/Rizwank123/LahziCart/delivery-service/internal/db"
	"github.com/Rizwank123/LahziCart/delivery-service/internal/handler"
	"github.com/Rizwank123/LahziCart/delivery-service/internal/repository"
	"github.com/Rizwank123/LahziCart/delivery-service/internal/router"
	"github.com/Rizwank123/LahziCart/delivery-service/internal/service"
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
	shutdownTracer, err := observability.InitTracer(context.Background(), "delivery-service")
	if err != nil {
		logger.Fatal("failed to initialize tracing", zap.Error(err))
	}
	defer shutdownTracer(context.Background())
	pool, err := db.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}
	defer pool.Close()
	repo := repository.NewShipmentRepository(pool)
	svc := service.NewShipmentService(repo)
	r := router.NewShipmentRouter(handler.NewShipmentHandler(svc))
	e := echo.New()
	e.Use(observability.TracingMiddleware("delivery-service"))
	e.Use(sharedmw.NewRecoveryMiddleware(logger).Recover)
	e.Use(sharedmw.NewRequestIDMiddleware("X-Request-ID").Inject)
	e.Use(sharedmw.NewLoggingMiddleware(logger, "/health").Log)
	e.Use(sharedmw.NewMetricsMiddleware("delivery-service").Collect)
	e.Use(sharedmw.NewCORSMiddleware(sharedmw.DefaultCORSConfig()))
	r.RegisterRoutes(e)
	e.GET("/health", func(c echo.Context) error { return c.String(http.StatusOK, "OK") })
	go func() {
		logger.Info("delivery-service started", zap.String("port", cfg.AppPort))
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
