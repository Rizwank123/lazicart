package handler

import (
	"net/http"
	"strconv"

	"github.com/Rizwank123/LahziCart/inventory-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type InventoryHandler struct {
	service models.InventoryService
}

func NewInventoryHandler(service models.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) Upsert(c echo.Context) error {
	var in models.UpsertInventoryInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	item, err := h.service.UpsertItem(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("inventory saved", item))
}

func (h *InventoryHandler) GetByProduct(c echo.Context) error {
	productID, err := uuid.FromString(c.Param("product_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid product id"))
	}
	item, err := h.service.GetByProductID(productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("inventory not found"))
	}
	return c.JSON(http.StatusOK, response.Success("inventory found", item))
}

func (h *InventoryHandler) List(c echo.Context) error {
	limit := queryInt(c, "limit", 20)
	page := queryInt(c, "page", 1)
	items, total, err := h.service.ListItems(limit, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to list inventory"))
	}
	return c.JSON(http.StatusOK, response.SuccessWithPagination("inventory found", items, map[string]interface{}{"total": total, "page": page, "limit": limit}))
}

func (h *InventoryHandler) Adjust(c echo.Context) error {
	var in models.AdjustInventoryInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	item, err := h.service.Adjust(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("inventory adjusted", item))
}

func (h *InventoryHandler) Reserve(c echo.Context) error {
	var in models.ReserveInventoryInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	reservations, err := h.service.Reserve(in)
	if err != nil {
		return c.JSON(http.StatusConflict, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("inventory reserved", reservations))
}

func (h *InventoryHandler) Release(c echo.Context) error {
	orderID, err := uuid.FromString(c.Param("order_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid order id"))
	}
	if err := h.service.Release(orderID); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to release inventory"))
	}
	return c.JSON(http.StatusOK, response.Success("inventory released", nil))
}

func (h *InventoryHandler) Commit(c echo.Context) error {
	orderID, err := uuid.FromString(c.Param("order_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid order id"))
	}
	if err := h.service.Commit(orderID); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to commit inventory"))
	}
	return c.JSON(http.StatusOK, response.Success("inventory committed", nil))
}

func queryInt(c echo.Context, name string, fallback int) int {
	value, err := strconv.Atoi(c.QueryParam(name))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
