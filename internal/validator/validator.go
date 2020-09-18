package validator

import v "github.com/go-playground/validator/v10"

// NewValidator creates a new validator.
func NewValidator() *v.Validate {
	validate := v.New()
	// validate.RegisterValidation("content-fields", validateMyVal)
	return validate
}
