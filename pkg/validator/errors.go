package validator

import (
	"errors"
	"fmt"
	"strings"

	v "github.com/go-playground/validator/v10"
)

// ErrValidation validation errors.
type ErrValidation struct {
	Errors []string `json:"errors"`
}

func (e ErrValidation) Error() string {
	return strings.Join(e.Errors[:], ",")
}

// NewErrValidation creates a err validation.
func NewErrValidation(err error) error {
	if _, ok := err.(*v.InvalidValidationError); ok {
		return errors.New("not a valid validation error")
	}

	var errors = []string{}
	for _, err := range err.(v.ValidationErrors) {
		// fmt.Println(err.Namespace()) // can differ when a custom TagNameFunc is registered or
		// fmt.Println(err.Field())     // by passing alt name to ReportError like below
		// fmt.Println(err.StructNamespace())
		// fmt.Println(err.StructField())
		// fmt.Println(err.Tag())
		// fmt.Println(err.ActualTag())
		// fmt.Println(err.Kind())
		// fmt.Println(err.Type())
		// fmt.Println(err.Value())
		// fmt.Println(err.Param())
		// fmt.Println()
		errors = append(errors, fmt.Sprintf("field %s fails contraint: %s %s", err.Field(), err.ActualTag(), err.Param()))
	}

	return ErrValidation{
		Errors: errors,
	}
}
