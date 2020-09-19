package validator

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/dwethmar/go-api/pkg/api/input"
	v "github.com/go-playground/validator/v10"
)

// Validation type
type Validation = v.Validate

// NewValidator creates a new validator.
func NewValidator() *Validation {
	validate := v.New()

	validate.RegisterCustomTypeFunc(ValidateValuer, input.FieldTranslations{})

	return validate
}

// ValidateValuer implements validator.CustomTypeFunc
func ValidateValuer(field reflect.Value) interface{} {

	fmt.Println(":D:D:D:D:")

	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}

	return nil
}
