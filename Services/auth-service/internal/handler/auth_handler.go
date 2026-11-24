package handler

import (
	"log"
	"net/http"

	"github.com/Rizwank123/LahziCart/auth-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	as models.UserAuthService
}

func NewAuthHandler(as models.UserAuthService) AuthHandler {
	return AuthHandler{as: as}
}

// @Summary Register
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.RegisterInput true "Register Input"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(ctx echo.Context) error {
	in := models.RegisterInput{}
	if err := ctx.Bind(&in); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	user, err := h.as.RegisterUser(in)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.Success("User registered successfully", user))
}

// @Summary Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.LoginUser true "Login Input"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(ctx echo.Context) error {
	in := models.LoginUser{}
	if err := ctx.Bind(&in); err != nil {
		log.Println("Bind error", err, ctx.Request().Body)
		return ctx.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	user, err := h.as.LoginUser(&in)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.Success("User logged in successfully", user))
}

// @Summary Send OTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param phone query string true "Phone Number"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/auth/send-otp [post]
func (h *AuthHandler) SendOtp(ctx echo.Context) error {
	phone := ctx.QueryParam("phone")
	err := h.as.SendOtp(phone)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.Success("OTP sent successfully", nil))
}

// @Summary Login with OTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.OtpLogin true "Login Input"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/auth/login-with-otp [post]
func (h *AuthHandler) LoginWithPhone(ctx echo.Context) error {
	in := models.OtpLogin{}
	if err := ctx.Bind(&in); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	user, err := h.as.LoginWithPhone(&in)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.Success("User logged in successfully", user))
}

// @Summary Refresh Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh_token query string true "Refresh Token"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(ctx echo.Context) error {
	in := ctx.QueryParam("refresh_token")
	result, err := h.as.RefreshToken(in)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.Success("Token refreshed successfully", result))
}

// @Summary Logout
// @Tags Auth
// @Accept json
// @Produce json
// @Param jti query string true "JTI"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(ctx echo.Context) error {
	refreshToken := ctx.QueryParam("jti")
	err := h.as.Logout(refreshToken)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.Success("User logged out successfully", nil))
}
