package router

import (
	"github.com/Rizwank123/LahziCart/payment-service/internal/handler"
	"github.com/labstack/echo/v4"
)

type PaymentRouter struct{ handler *handler.PaymentHandler }

func NewPaymentRouter(handler *handler.PaymentHandler) *PaymentRouter {
	return &PaymentRouter{handler: handler}
}

func (r *PaymentRouter) RegisterRoutes(e *echo.Echo) {
	e.POST("/payments", r.handler.Create)
	e.GET("/payments/:id", r.handler.Get)
	e.GET("/payments/orders/:order_id", r.handler.GetByOrder)
	e.PUT("/payments/:id/status", r.handler.UpdateStatus)
}
