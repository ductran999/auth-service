package request

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ParseAndValidateJSON[T any](ctx *gin.Context) (*T, error) {
	var payload T
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) && len(ve) > 0 {
			// You must inject a common error response function
			return nil, fmt.Errorf("validation error: %s", validationErrorMessage(ve[0]))
		}
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return &payload, nil
}

func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Field() {
	case "Email":
		switch fe.Tag() {
		case "required":
			return "Email is required"
		case "email":
			return "Email must be valid"
		}
	case "Password":
		switch fe.Tag() {
		case "required":
			return "Password is required"
		case "password":
			return "Password must include at least 1 uppercase, 1 lowercase, 1 number, and 1 special character"
		}
	}

	return "Invalid input"
}
