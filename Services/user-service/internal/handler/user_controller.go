package handler

import (
	"net/http"

	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/Rizwank123/lahziCart/user-service/internal/models"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	us models.UserService
}

func NewUserController(us models.UserService) *UserController {
	return &UserController{us: us}
}

// @Summary Find User By Email
// @Tags User
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/users/email [get]
func (c *UserController) FindUserByEmail(ctx echo.Context) error {
	email := ctx.QueryParam("email")
	if email == "" {

		return ctx.JSON(http.StatusBadRequest, response.Error("email is required"))
	}
	user, err := c.us.GetUserByEmail(email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.Error("failed to get user"))
	}
	return ctx.JSON(http.StatusOK, response.Success("user found", user))

}

// @Summary Find User By ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/users/{id} [get]
func (c *UserController) FindUserByID(ctx echo.Context) error {
	id, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error("id is required"))
	}
	user, err := c.us.GetUserByID(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.Error("failed to get user"))
	}
	return ctx.JSON(http.StatusOK, response.Success("user found", user))

}

// @Summary Find User By Phone
// @Tags User
// @Accept json
// @Produce json
// @Param phone query string true "Phone"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/users/phone [get]
func (c *UserController) FindUserByPhone(ctx echo.Context) error {
	phone := ctx.QueryParam("phone")
	if phone == "" {
		return ctx.JSON(http.StatusBadRequest, response.Error("phone is required"))
	}
	user, err := c.us.GetUserByPhone(phone)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.Error("failed to get user"))
	}
	return ctx.JSON(http.StatusOK, response.Success("user found", user))

}

// @Summary Update User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param input body models.UpdateUser true "Update User"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/users/{id} [put]
func (c *UserController) UpdateUser(ctx echo.Context) error {
	id, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error("id is required"))
	}
	var in models.UpdateUser
	if err := ctx.Bind(&in); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	user, err := c.us.UpdateUser(id, in)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.Error("failed to update user"))
	}
	return ctx.JSON(http.StatusOK, response.Success("user updated", user))

}

// @Summary Delete User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/users/{id} [delete]
func (c *UserController) DeleteUser(ctx echo.Context) error {
	id, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error("id is required"))
	}
	user, err := c.us.DeleteUser(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.Error("failed to delete user"))
	}
	return ctx.JSON(http.StatusOK, response.Success("user deleted", user))

}
