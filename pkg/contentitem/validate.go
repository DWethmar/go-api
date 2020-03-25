package contentitem

import (
	"errors"
	"fmt"
)

var MaxNameLength = 50

type ErrUnsupportedAttrType struct {
	Attr string
}

func (e *ErrUnsupportedAttrType) Error() string {
	return fmt.Sprintf("Attribute type of %v is not supported.", e.Attr)
}

type NameLengthError struct{}

func (e *NameLengthError) Error() string {
	return fmt.Sprintf("Name exceeded max length of %d.", MaxNameLength)
}

var (
	ErrAttrsRequired = errors.New("Attrs is required.")
)

func ValidateName(name string) []error {
	var e []error

	if len(name) > MaxNameLength {
		e = append(e, &NameLengthError{})
	}

	return e
}

func ValidateAttr(attr Attrs) []error {
	var e []error

	validateAttr := func(attrs map[string]interface{}) []error {
		ce := make([]error, 0)
		for key, value := range attrs {
			switch value.(type) {
			case int:
			case []int:
			case string:
			case []string:
			case float64:
			case bool:
			default:
				e = append(e, &ErrUnsupportedAttrType{
					Attr: key,
				})
			}
		}
		return ce
	}

	validateAttr(attr)
	return e
}
