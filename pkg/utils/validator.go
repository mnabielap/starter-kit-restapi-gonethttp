package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

// ValidateStruct validates a struct and returns a list of errors if any
func ValidateStruct(data interface{}) map[string]string {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = "Failed validation: " + err.Tag()
	}
	return errors
}