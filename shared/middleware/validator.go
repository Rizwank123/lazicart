package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator wraps go-playground validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate validates a struct based on tags
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// BindAndValidate binds request data and validates in one step
func BindAndValidate(c echo.Context, i interface{}) error {
	// Bind the request
	if err := c.Bind(i); err != nil {
		return echo.NewHTTPError(400, "invalid request format")
	}

	// Validate the bound data
	if err := c.Validate(i); err != nil {
		// Format validation errors nicely
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return echo.NewHTTPError(400, formatValidationErrors(validationErrors))
		}
		return echo.NewHTTPError(400, "validation failed")
	}

	return nil
}

// formatValidationErrors formats validator.ValidationErrors into readable messages
func formatValidationErrors(errs validator.ValidationErrors) string {
	if len(errs) == 0 {
		return "validation failed"
	}

	// Return first error message
	err := errs[0]
	field := err.Field()
	tag := err.Tag()

	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + err.Param() + " characters"
	case "max":
		return field + " must be at most " + err.Param() + " characters"
	case "len":
		return field + " must be exactly " + err.Param() + " characters"
	case "gt":
		return field + " must be greater than " + err.Param()
	case "gte":
		return field + " must be greater than or equal to " + err.Param()
	case "lt":
		return field + " must be less than " + err.Param()
	case "lte":
		return field + " must be less than or equal to " + err.Param()
	default:
		return field + " validation failed on " + tag
	}
}
