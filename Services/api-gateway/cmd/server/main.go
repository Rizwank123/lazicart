package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Rizwank123/LahziCart/api-gateway/internal/config"
	gatewaymw "github.com/Rizwank123/LahziCart/api-gateway/internal/middleware"
	"github.com/Rizwank123/LahziCart/api-gateway/internal/router"
	"github.com/Rizwank123/LahziCart/shared/observability"
)

// @title API Gateway
// @version 1.0
// @description API Gateway
// @host localhost:7070
// @BasePath /

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}
	shutdownTracer, err := observability.InitTracer(context.Background(), "api-gateway")
	if err != nil {
		logger.Fatal("failed to initialize tracing", zap.Error(err))
	}
	defer shutdownTracer(context.Background())

	e := echo.New()
	e.Use(observability.TracingMiddleware("api-gateway"))
	e.Use(gatewaymw.NewPropagationMiddleware().Inject)

	router.Register(e, cfg)

	go func() {
		if err := e.Start(":" + cfg.AppPort); err != nil && err != http.ErrServerClosed {
			logger.Fatal("shutting down server", zap.Error(err))
		}
	}()

	logger.Info("API Gateway Started", zap.String("port", cfg.AppPort))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal(err.Error())
	}
}
