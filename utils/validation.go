package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// single shared validator instance
var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []ValidationError {
	var validationErrors []ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var ve ValidationError
			ve.Field = strings.ToLower(err.Field())
			ve.Message = formatValidationError(err)
			validationErrors = append(validationErrors, ve)
		}
	}

	return validationErrors
}

func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", strings.ToLower(err.Field()))
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", strings.ToLower(err.Field()), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", strings.ToLower(err.Field()), err.Param())
	case "future_date":
		return "dateTime must be in the future"
	default:
		return fmt.Sprintf("%s is invalid", strings.ToLower(err.Field()))
	}
}

// adds any custom rules beyond what the package provides
func RegisterCustomValidations() {
	validate.RegisterValidation("future_date", func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		if !ok {
			return false
		}
		return date.After(time.Now())
	})
}
