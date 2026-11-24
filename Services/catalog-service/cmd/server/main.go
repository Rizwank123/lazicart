package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rizwank123/LahziCart/catalog-service/internal/config"
	"github.com/Rizwank123/LahziCart/catalog-service/internal/db"
	"github.com/Rizwank123/LahziCart/catalog-service/internal/handler"
	"github.com/Rizwank123/LahziCart/catalog-service/internal/repository"
	"github.com/Rizwank123/LahziCart/catalog-service/internal/router"
	"github.com/Rizwank123/LahziCart/catalog-service/internal/service"
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

	productRepo := repository.NewProductRepository(pool)
	productService := service.NewProductService(productRepo, logger)
	productHandler := handler.NewProductHandler(productService)
	productRouter := router.NewProductRouter(productHandler)

	e := echo.New()
	e.Validator = sharedmw.NewValidator()
	e.Use(sharedmw.NewRecoveryMiddleware(logger).Recover)
	e.Use(sharedmw.NewRequestIDMiddleware("X-Request-ID").Inject)
	e.Use(sharedmw.NewLoggingMiddleware(logger, "/health").Log)
	e.Use(sharedmw.NewMetricsMiddleware("catalog-service").Collect)
	e.Use(sharedmw.NewCORSMiddleware(sharedmw.DefaultCORSConfig()))

	productRouter.RegisterRoutes(e)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	go func() {
		logger.Info("catalog-service started", zap.String("port", cfg.AppPort))
		if err := e.Start(":" + cfg.AppPort); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", zap.Error(err))
	}
}
