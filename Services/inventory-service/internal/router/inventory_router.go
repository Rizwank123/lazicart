package router

import (
	"github.com/Rizwank123/LahziCart/inventory-service/internal/handler"
	"github.com/labstack/echo/v4"
)

type InventoryRouter struct {
	handler *handler.InventoryHandler
}

func NewInventoryRouter(handler *handler.InventoryHandler) *InventoryRouter {
	return &InventoryRouter{handler: handler}
}

func (r *InventoryRouter) RegisterRoutes(e *echo.Echo) {
	e.GET("/inventory", r.handler.List)
	e.PUT("/inventory", r.handler.Upsert)
	e.POST("/inventory/adjust", r.handler.Adjust)
	e.GET("/inventory/products/:product_id", r.handler.GetByProduct)
	e.POST("/inventory/reservations", r.handler.Reserve)
	e.POST("/inventory/reservations/:order_id/release", r.handler.Release)
	e.POST("/inventory/reservations/:order_id/commit", r.handler.Commit)
}
