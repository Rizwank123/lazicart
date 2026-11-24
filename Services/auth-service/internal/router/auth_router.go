package router

import (
	"github.com/Rizwank123/LahziCart/auth-service/internal/handler"
	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
	h handler.AuthHandler
}

func NewAuthRouter(h handler.AuthHandler) AuthRouter {
	return AuthRouter{h: h}
}

func (r *AuthRouter) RegisterRoutes(e *echo.Echo) {
	e.POST("/register", r.h.Register)
	e.POST("/login", r.h.Login)
	e.POST("/send-otp", r.h.SendOtp)
	e.POST("/login-with-phone", r.h.LoginWithPhone)
	e.POST("/refresh-token", r.h.RefreshToken)
	e.GET("/logout", r.h.Logout)
}
