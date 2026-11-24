package handler

import (
	"net/http"

	"github.com/Rizwank123/LahziCart/cart-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	service models.CartService
}

func NewCartHandler(service models.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) GetCart(c echo.Context) error {
	userID, err := uuid.FromString(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid user id"))
	}
	cart, err := h.service.GetCart(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to get cart"))
	}
	return c.JSON(http.StatusOK, response.Success("cart found", cart))
}

func (h *CartHandler) AddItem(c echo.Context) error {
	var in models.AddCartItemInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	cart, err := h.service.AddItem(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("cart item added", cart))
}

func (h *CartHandler) UpdateItem(c echo.Context) error {
	itemID, err := uuid.FromString(c.Param("item_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid item id"))
	}
	var in models.UpdateCartItemInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	if err := h.service.UpdateItem(itemID, in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("cart item updated", nil))
}

func (h *CartHandler) RemoveItem(c echo.Context) error {
	itemID, err := uuid.FromString(c.Param("item_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid item id"))
	}
	if err := h.service.RemoveItem(itemID); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to remove item"))
	}
	return c.JSON(http.StatusOK, response.Success("cart item removed", nil))
}

func (h *CartHandler) Clear(c echo.Context) error {
	userID, err := uuid.FromString(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid user id"))
	}
	if err := h.service.ClearCart(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to clear cart"))
	}
	return c.JSON(http.StatusOK, response.Success("cart cleared", nil))
}
