package db

import (
	"context"
	"fmt"

	"github.com/Rizwank123/LahziCart/catalog-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(cfg *config.CatalogConfig) (*pgxpool.Pool, error) {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	pool, err := pgxpool.New(context.Background(), dns)
	if err != nil {
		return nil, err
	}
	return pool, nil

}
