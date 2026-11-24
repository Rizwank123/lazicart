package db

import (
	"context"
	"fmt"

	"github.com/Rizwank123/LahziCart/auth-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	Pool *pgxpool.Pool
}

func NewDb(c *config.AuthConfig) (*Connection, error) {
	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
	pool, err := pgxpool.New(context.Background(), dns)
	if err != nil {
		return nil, err
	}
	return &Connection{Pool: pool}, nil

}
