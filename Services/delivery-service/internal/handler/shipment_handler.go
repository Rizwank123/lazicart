package handler

import (
	"net/http"
	"strconv"

	"github.com/Rizwank123/LahziCart/delivery-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/identity"
	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type ShipmentHandler struct{ service models.ShipmentService }

func NewShipmentHandler(service models.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{service: service}
}
func (h *ShipmentHandler) Create(c echo.Context) error {
	var in models.CreateShipmentInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	if err := identity.RequireOwnershipUUID(c, in.UserID); err != nil {
		return err
	}
	s, err := h.service.Create(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success("shipment created", s))
}
func (h *ShipmentHandler) Get(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid shipment id"))
	}
	s, err := h.service.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("shipment not found"))
	}
	if err := identity.RequireOwnershipUUID(c, s.UserID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.Success("shipment found", s))
}
func (h *ShipmentHandler) GetByOrder(c echo.Context) error {
	id, err := uuid.FromString(c.Param("order_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid order id"))
	}
	s, err := h.service.GetByOrder(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("shipment not found"))
	}
	if err := identity.RequireOwnershipUUID(c, s.UserID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.Success("shipment found", s))
}
func (h *ShipmentHandler) ListByUser(c echo.Context) error {
	id, err := uuid.FromString(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid user id"))
	}
	if err := identity.RequireOwnershipUUID(c, id); err != nil {
		return err
	}
	items, err := h.service.ListByUser(id, queryInt(c, "limit", 20), queryInt(c, "page", 1))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to list shipments"))
	}
	return c.JSON(http.StatusOK, response.Success("shipments found", items))
}
func (h *ShipmentHandler) UpdateStatus(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid shipment id"))
	}
	// The shipment is loaded first so that its owner can be checked before the
	// status is changed. Requests that did not come through the gateway, such
	// as an internal dispatch callback, are not affected.
	s, err := h.service.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("shipment not found"))
	}
	if err := identity.RequireOwnershipUUID(c, s.UserID); err != nil {
		return err
	}
	var in models.UpdateShipmentStatusInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	s, err = h.service.UpdateStatus(id, in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("shipment updated", s))
}
func queryInt(c echo.Context, name string, fallback int) int {
	value, err := strconv.Atoi(c.QueryParam(name))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
