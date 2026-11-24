package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Rizwank123/LahziCart/auth-service/internal/cache"
	"github.com/Rizwank123/LahziCart/auth-service/internal/config"
	"github.com/Rizwank123/LahziCart/auth-service/internal/db"
	"github.com/Rizwank123/LahziCart/auth-service/internal/grpc_handler"
	"github.com/Rizwank123/LahziCart/auth-service/internal/handler"
	"github.com/Rizwank123/LahziCart/auth-service/internal/repository"
	"github.com/Rizwank123/LahziCart/auth-service/internal/router"
	"github.com/Rizwank123/LahziCart/auth-service/internal/service"
	"github.com/Rizwank123/LahziCart/auth-service/internal/utils"

	kafkabroker "github.com/Rizwank123/LahziCart/shared/messaging/kafkaBroker"
	"github.com/Rizwank123/LahziCart/shared/observability"
	authpb "github.com/Rizwank123/LahziCart/shared/proto"
)

func main() {

	// ------------------------------------------------
	// Logger
	// ------------------------------------------------
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// ------------------------------------------------
	// Config
	// ------------------------------------------------
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// ------------------------------------------------
	// Database
	// ------------------------------------------------
	dbConn, err := db.NewDb(cfg)
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}

	// ------------------------------------------------
	// Redis
	// ------------------------------------------------
	redisClient, err := cache.NewRedisClient(
		cfg.RedisHost,
		cfg.RedisPort,
	)
	if err != nil {
		logger.Fatal("failed to connect redis", zap.Error(err))
	}

	// ------------------------------------------------
	// Kafka
	// ------------------------------------------------
	kafkaProducer := kafkabroker.InitKafkaProducer()
	if kafkaProducer == nil {
		logger.Fatal("failed to init kafka producer")
	}

	// ------------------------------------------------
	// Observability
	// ------------------------------------------------
	observability.Register()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("metrics server started on :9090")
		_ = http.ListenAndServe(":9090", nil)
	}()

	// ------------------------------------------------
	// Repository
	// ------------------------------------------------
	authRepo := repository.NewPgxRepository(dbConn.Pool)

	// ------------------------------------------------
	// Utils
	// ------------------------------------------------
	authManager := utils.NewAuthManager(cfg)

	// ------------------------------------------------
	// Service
	// ------------------------------------------------
	authService := service.NewAuthServiceImpl(
		authRepo,
		authManager,
		redisClient,
		cfg,
		kafkaProducer,
	)

	// ------------------------------------------------
	// HTTP Server
	// ------------------------------------------------
	e := echo.New()

	handlers := handler.NewAuthHandler(authService)
	authRouter := router.NewAuthRouter(handlers)
	authRouter.RegisterRoutes(e)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	httpServer := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: e,
	}

	// ------------------------------------------------
	// gRPC Server
	// ------------------------------------------------
	grpcServer := grpc.NewServer()

	grpcAuthHandler := grpc_handler.NewAuthGRPCServer(authService)
	authpb.RegisterAuthServiceServer(grpcServer, grpcAuthHandler)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		logger.Fatal("failed to listen grpc", zap.Error(err))
	}

	// ------------------------------------------------
	// Start Servers
	// ------------------------------------------------
	go func() {
		logger.Info("HTTP server started", zap.String("port", cfg.AppPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server error", zap.Error(err))
		}
	}()

	go func() {
		logger.Info("gRPC server started", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("grpc server error", zap.Error(err))
		}
	}()

	// ------------------------------------------------
	// Graceful Shutdown
	// ------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop HTTP
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("http shutdown error", zap.Error(err))
	}

	// Stop gRPC
	grpcServer.GracefulStop()

	// Close DB
	dbConn.Pool.Close()

	// Close Redis
	if redisClient != nil {
		_ = redisClient.Close()
	}

	// Close Kafka
	if kafkaProducer != nil {
		kafkaProducer.Close()
	}

	logger.Info("auth-service stopped gracefully")
}
