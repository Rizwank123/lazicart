package router

import (
	"github.com/Rizwank123/lahziCart/user-service/internal/handler"
	"github.com/labstack/echo/v4"
)

type UserRouter struct {
	h *handler.UserController
}

func NewUserRouter(h *handler.UserController) *UserRouter {
	return &UserRouter{h: h}
}

func (h *UserRouter) RegisterRoutes(e *echo.Echo) {

	e.GET("/users/:id", h.h.FindUserByID)
	e.GET("/users/email", h.h.FindUserByEmail)
	e.GET("/users/phone", h.h.FindUserByPhone)
	e.PUT("/users/:id", h.h.UpdateUser)
	e.DELETE("/users/:id", h.h.DeleteUser)
}
