package handler

import (
	"net/http"
	"strconv"

	"github.com/Rizwank123/LahziCart/notification-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct{ service models.NotificationService }

func NewNotificationHandler(service models.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) Create(c echo.Context) error {
	var in models.CreateNotificationInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	n, err := h.service.Create(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success("notification created", n))
}

func (h *NotificationHandler) ListByUser(c echo.Context) error {
	userID, err := uuid.FromString(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid user id"))
	}
	limit := queryInt(c, "limit", 20)
	page := queryInt(c, "page", 1)
	items, err := h.service.ListByUser(userID, limit, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to list notifications"))
	}
	return c.JSON(http.StatusOK, response.Success("notifications found", items))
}

func (h *NotificationHandler) MarkSent(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid notification id"))
	}
	if err := h.service.MarkSent(id); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to mark notification sent"))
	}
	return c.JSON(http.StatusOK, response.Success("notification sent", nil))
}

func queryInt(c echo.Context, name string, fallback int) int {
	value, err := strconv.Atoi(c.QueryParam(name))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
