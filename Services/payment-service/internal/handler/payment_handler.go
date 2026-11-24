package handler

import (
	"net/http"

	"github.com/Rizwank123/LahziCart/payment-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type PaymentHandler struct{ service models.PaymentService }

func NewPaymentHandler(service models.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) Create(c echo.Context) error {
	var in models.CreatePaymentInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	payment, err := h.service.CreatePayment(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success("payment created", payment))
}

func (h *PaymentHandler) Get(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid payment id"))
	}
	payment, err := h.service.GetPayment(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("payment not found"))
	}
	return c.JSON(http.StatusOK, response.Success("payment found", payment))
}

func (h *PaymentHandler) GetByOrder(c echo.Context) error {
	orderID, err := uuid.FromString(c.Param("order_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid order id"))
	}
	payment, err := h.service.GetByOrderID(orderID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("payment not found"))
	}
	return c.JSON(http.StatusOK, response.Success("payment found", payment))
}

func (h *PaymentHandler) UpdateStatus(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid payment id"))
	}
	var in models.UpdatePaymentStatusInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	payment, err := h.service.UpdateStatus(id, in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("payment updated", payment))
}
