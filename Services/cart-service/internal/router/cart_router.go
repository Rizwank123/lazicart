package router

import (
	"github.com/Rizwank123/LahziCart/cart-service/internal/handler"
	"github.com/labstack/echo/v4"
)

type CartRouter struct {
	handler *handler.CartHandler
}

func NewCartRouter(handler *handler.CartHandler) *CartRouter {
	return &CartRouter{handler: handler}
}

func (r *CartRouter) RegisterRoutes(e *echo.Echo) {
	e.GET("/carts/:user_id", r.handler.GetCart)
	e.POST("/carts/items", r.handler.AddItem)
	e.PUT("/carts/items/:item_id", r.handler.UpdateItem)
	e.DELETE("/carts/items/:item_id", r.handler.RemoveItem)
	e.DELETE("/carts/:user_id", r.handler.Clear)
}
