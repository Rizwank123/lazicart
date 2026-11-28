package router

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	"github.com/Rizwank123/LahziCart/api-gateway/internal/config"
	"github.com/Rizwank123/LahziCart/api-gateway/internal/middleware"
	"github.com/Rizwank123/LahziCart/api-gateway/internal/proxy"
)

func Register(e *echo.Echo, cfg *config.ApiGateWayConfig, sessions middleware.SessionStore, logger *zap.Logger) {

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

	// Every protected route first verifies the token, then checks that the
	// session is still live, then resolves the caller identity.
	protected := func(group *echo.Group) *echo.Group {
		group.Use(middleware.JWTMiddleware(cfg.JWTSecret))
		group.Use(middleware.RevocationMiddleware(sessions, logger))
		group.Use(middleware.IdentityMiddleware)
		return group
	}

	userGroup := protected(e.Group("/api/users"))
	userGroup.Any("/*", userProxy)

	productGroup := protected(e.Group("/api/products"))
	productGroup.POST("", productProxy)
	productGroup.PUT("/*", productProxy)
	productGroup.DELETE("/*", productProxy)
	productGroup.POST("/*", productProxy)

	inventoryGroup := protected(e.Group("/api/inventory"))
	inventoryGroup.Any("", inventoryProxy)
	inventoryGroup.Any("/*", inventoryProxy)

	cartGroup := protected(e.Group("/api/carts"))
	cartGroup.Any("", cartProxy)
	cartGroup.Any("/*", cartProxy)

	orderGroup := protected(e.Group("/api/orders"))
	orderGroup.Any("", orderProxy)
	orderGroup.Any("/*", orderProxy)

	paymentGroup := protected(e.Group("/api/payments"))
	paymentGroup.Any("", paymentProxy)
	paymentGroup.Any("/*", paymentProxy)

	notificationGroup := protected(e.Group("/api/notifications"))
	notificationGroup.Any("", notificationProxy)
	notificationGroup.Any("/*", notificationProxy)

	deliveryGroup := protected(e.Group("/api/deliveries"))
	deliveryGroup.Any("", deliveryProxy)
	deliveryGroup.Any("/*", deliveryProxy)
}
