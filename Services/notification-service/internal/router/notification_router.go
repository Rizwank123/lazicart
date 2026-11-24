package router

import (
	"github.com/Rizwank123/LahziCart/notification-service/internal/handler"
	"github.com/labstack/echo/v4"
)

type NotificationRouter struct{ handler *handler.NotificationHandler }

func NewNotificationRouter(handler *handler.NotificationHandler) *NotificationRouter {
	return &NotificationRouter{handler: handler}
}

func (r *NotificationRouter) RegisterRoutes(e *echo.Echo) {
	e.POST("/notifications", r.handler.Create)
	e.GET("/notifications/users/:user_id", r.handler.ListByUser)
	e.POST("/notifications/:id/sent", r.handler.MarkSent)
}
