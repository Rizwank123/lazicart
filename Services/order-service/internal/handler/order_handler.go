package handler

import (
	"net/http"
	"strconv"

	"github.com/Rizwank123/LahziCart/order-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	service models.OrderService
}

func NewOrderHandler(service models.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(c echo.Context) error {
	var in models.CreateOrderInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	order, err := h.service.CreateOrder(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success("order created", order))
}

func (h *OrderHandler) Get(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid order id"))
	}
	order, err := h.service.GetOrder(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("order not found"))
	}
	return c.JSON(http.StatusOK, response.Success("order found", order))
}

func (h *OrderHandler) ListByUser(c echo.Context) error {
	userID, err := uuid.FromString(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid user id"))
	}
	limit := queryInt(c, "limit", 20)
	page := queryInt(c, "page", 1)
	orders, err := h.service.ListByUser(userID, limit, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to list orders"))
	}
	return c.JSON(http.StatusOK, response.Success("orders found", orders))
}

func (h *OrderHandler) Cancel(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid order id"))
	}
	if err := h.service.CancelOrder(id); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("order cancelled", nil))
}

func queryInt(c echo.Context, name string, fallback int) int {
	value, err := strconv.Atoi(c.QueryParam(name))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
