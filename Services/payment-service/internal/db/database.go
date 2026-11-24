package db

import (
	"context"
	"fmt"

	"github.com/Rizwank123/LahziCart/payment-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(cfg *config.PaymentConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	return pgxpool.New(context.Background(), dsn)
}
