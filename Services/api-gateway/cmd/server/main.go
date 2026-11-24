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
	"github.com/Rizwank123/LahziCart/api-gateway/internal/router"
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

	e := echo.New()

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
