package router

import (
	"github.com/Rizwank123/LahziCart/delivery-service/internal/handler"
	"github.com/labstack/echo/v4"
)

type ShipmentRouter struct{ handler *handler.ShipmentHandler }

func NewShipmentRouter(handler *handler.ShipmentHandler) *ShipmentRouter {
	return &ShipmentRouter{handler: handler}
}
func (r *ShipmentRouter) RegisterRoutes(e *echo.Echo) {
	e.POST("/deliveries", r.handler.Create)
	e.GET("/deliveries/:id", r.handler.Get)
	e.GET("/deliveries/orders/:order_id", r.handler.GetByOrder)
	e.GET("/deliveries/users/:user_id", r.handler.ListByUser)
	e.PUT("/deliveries/:id/status", r.handler.UpdateStatus)
}
