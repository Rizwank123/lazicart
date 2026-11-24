module github.com/Rizwank123/LahziCart/cart-service

go 1.24.12

replace github.com/Rizwank123/LahziCart/shared => ../../shared

require (
	github.com/Rizwank123/LahziCart/shared v0.0.0-00010101000000-000000000000
	github.com/gofrs/uuid/v5 v5.4.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/labstack/echo/v4 v4.15.0
	github.com/spf13/viper v1.21.0
	go.uber.org/zap v1.27.1
)
