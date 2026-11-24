package router

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/Rizwank123/LahziCart/api-gateway/internal/config"
	"github.com/Rizwank123/LahziCart/api-gateway/internal/middleware"
	"github.com/Rizwank123/LahziCart/api-gateway/internal/proxy"
)

func Register(e *echo.Echo, cfg *config.ApiGateWayConfig) {

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	authProxy := proxy.NewProxy(cfg.AuthServiceURL, "/api/auth", "")
	userProxy := proxy.NewProxy(cfg.UserServiceURL, "/api/users", "/users")
	productProxy := proxy.NewProxy(cfg.ProductServiceURL, "/api/products", "/products")
	inventoryProxy := proxy.NewProxy(cfg.InventoryServiceURL, "/api/inventory", "/inventory")
	cartProxy := proxy.NewProxy(cfg.CartServiceURL, "/api/carts", "/carts")
	orderProxy := proxy.NewProxy(cfg.OrderServiceURL, "/api/orders", "/orders")
	paymentProxy := proxy.NewProxy(cfg.PaymentServiceURL, "/api/payments", "/payments")
	notificationProxy := proxy.NewProxy(cfg.NotificationServiceURL, "/api/notifications", "/notifications")
	deliveryProxy := proxy.NewProxy(cfg.DeliveryServiceURL, "/api/deliveries", "/deliveries")

	e.Any("/api/auth/*", authProxy)
	e.GET("/api/products", productProxy)
	e.GET("/api/products/*", productProxy)

	userGroup := e.Group("/api/users")
	userGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	userGroup.Any("/*", userProxy)

	productGroup := e.Group("/api/products")
	productGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	productGroup.POST("", productProxy)
	productGroup.PUT("/*", productProxy)
	productGroup.DELETE("/*", productProxy)
	productGroup.POST("/*", productProxy)

	inventoryGroup := e.Group("/api/inventory")
	inventoryGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	inventoryGroup.Any("", inventoryProxy)
	inventoryGroup.Any("/*", inventoryProxy)

	cartGroup := e.Group("/api/carts")
	cartGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	cartGroup.Any("", cartProxy)
	cartGroup.Any("/*", cartProxy)

	orderGroup := e.Group("/api/orders")
	orderGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	orderGroup.Any("", orderProxy)
	orderGroup.Any("/*", orderProxy)

	paymentGroup := e.Group("/api/payments")
	paymentGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	paymentGroup.Any("", paymentProxy)
	paymentGroup.Any("/*", paymentProxy)

	notificationGroup := e.Group("/api/notifications")
	notificationGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	notificationGroup.Any("", notificationProxy)
	notificationGroup.Any("/*", notificationProxy)

	deliveryGroup := e.Group("/api/deliveries")
	deliveryGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	deliveryGroup.Any("", deliveryProxy)
	deliveryGroup.Any("/*", deliveryProxy)
}
