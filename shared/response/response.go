package response

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

func Success(message string, data interface{}) Response {
	return Response{
		Message: message,
		Data:    data,
	}
}

func Error(message string) Response {
	return Response{
		Message: message,
	}
}

func SuccessWithPagination(message string, data interface{}, pagination interface{}) Response {
	return Response{
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}

// Unauthorized returns a 401 unauthorized error
func Unauthorized(c echo.Context, message string) error {
	return c.JSON(401, Error(message))
}

// BadRequest returns a 400 bad request error
func BadRequest(c echo.Context, message string) error {
	return c.JSON(400, Error(message))
}

// OK returns a 200 success response
func OK(c echo.Context, message string, data interface{}) error {
	return c.JSON(200, Success(message, data))
}
