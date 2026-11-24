package router

import (
	"github.com/Rizwank123/LahziCart/order-service/internal/handler"
	"github.com/labstack/echo/v4"
)

type OrderRouter struct {
	handler *handler.OrderHandler
}

func NewOrderRouter(handler *handler.OrderHandler) *OrderRouter {
	return &OrderRouter{handler: handler}
}

func (r *OrderRouter) RegisterRoutes(e *echo.Echo) {
	e.POST("/orders", r.handler.Create)
	e.GET("/orders/:id", r.handler.Get)
	e.GET("/orders/users/:user_id", r.handler.ListByUser)
	e.POST("/orders/:id/cancel", r.handler.Cancel)
}
