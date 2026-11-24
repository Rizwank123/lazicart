package router

import (
	"github.com/Rizwank123/LahziCart/catalog-service/internal/handler"
	"github.com/labstack/echo/v4"
)

type ProductRouter struct {
	handler *handler.ProductHandler
}

func NewProductRouter(handler *handler.ProductHandler) *ProductRouter {
	return &ProductRouter{handler: handler}
}

func (r *ProductRouter) RegisterRoutes(e *echo.Echo) {
	e.POST("/products", r.handler.CreateProduct)
	e.GET("/products", r.handler.ListProducts)
	e.GET("/products/:id", r.handler.GetProduct)
	e.PUT("/products/:id", r.handler.UpdateProduct)
	e.DELETE("/products/:id", r.handler.DeleteProduct)

	e.POST("/products/reviews", r.handler.CreateReview)
	e.GET("/products/:id/reviews", r.handler.ListReviews)
	e.PUT("/products/reviews/:id", r.handler.UpdateReview)
	e.DELETE("/products/reviews/:id", r.handler.DeleteReview)

	e.POST("/products/images", r.handler.CreateImage)
	e.GET("/products/:id/images", r.handler.ListImages)
	e.PUT("/products/images/:id", r.handler.UpdateImage)
	e.DELETE("/products/images/:id", r.handler.DeleteImage)
}
