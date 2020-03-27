package contentitem

import (
	"errors"
	"reflect"
)

var (
	MaxNameLength                = 50
	ErrAttrsRequired             = errors.New("Attrs is required.")
	ErrNameLength                = errors.New("Name exceeded max length of 50.")
	ErrUnsupportedAttrValue      = errors.New("Value is unsupported.")
	ErrUnsupportedAttrSliceValue = errors.New("This array contains one or more unsupported values.")
)

type validationResult struct {
	Errors struct {
		Name  string                       `json:"name,omitempty"`
		Attrs map[string]map[string]string `json:"attrs,omitempty"`
	} `json:"errors"`
}

func (v *validationResult) IsValid() bool {
	if v.Errors.Name != "" || len(v.Errors.Attrs) > 0 {
		return false
	}
	return true
}

func ValidateAddContentItem(addContentItem AddContentItem) validationResult {
	e := CreateValidationResult()

	nameErr := ValidateName(addContentItem.Name)
	if nameErr != nil {
		e.Errors.Name = nameErr.Error()
	}

	for locale, attrs := range addContentItem.Attrs {
		attrErRors := ValidateAttr(attrs)
		for attr, error := range attrErRors {
			if e.Errors.Attrs[locale] == nil {
				e.Errors.Attrs[locale] = map[string]string{}
			}
			e.Errors.Attrs[locale][attr] = error.Error()
		}
	}
	return e
}

func ValidateContentItem(contentItem ContentItem) validationResult {
	e := CreateValidationResult()

	nameErr := ValidateName(contentItem.Name)
	if nameErr != nil {
		e.Errors.Name = nameErr.Error()
	}

	for locale, attrs := range contentItem.Attrs {
		attrErRors := ValidateAttr(attrs)
		for attr, error := range attrErRors {
			e.Errors.Attrs[locale][attr] = error.Error()
		}
	}
	return e
}

func ValidateName(name string) error {
	if len(name) > MaxNameLength {
		return ErrNameLength
	}
	return nil
}

func ValidateAttr(attrs map[string]interface{}) map[string]error {
	var e = make(map[string]error)

	validTypes := []reflect.Kind{
		reflect.Int,
		reflect.String,
		reflect.Bool,
		reflect.Float64,
	}

	isValidType := func(kind reflect.Kind) bool {
		for _, item := range validTypes {
			if item == kind {
				return true
			}
		}
		return false
	}

	for attr, value := range attrs {
		t := reflect.TypeOf(value)
		if t == nil {
			e[attr] = ErrUnsupportedAttrValue
			continue
		}

		kind := t.Kind()

		if kind == reflect.Invalid {
			e[attr] = ErrUnsupportedAttrValue
			continue
		}

		valid := isValidType(kind)
		if !valid {
			if kind == reflect.Slice {

				s := reflect.ValueOf(value)

				for i := 0; i < s.Len(); i++ {

					z := s.Index(i).Interface()
					k := reflect.TypeOf(z).Kind()

					if !isValidType(k) {
						e[attr] = ErrUnsupportedAttrSliceValue
						break
					}
				}
			} else {
				e[attr] = ErrUnsupportedAttrValue
			}
		}
	}

	return e
}

func CreateValidationResult() validationResult {
	v := validationResult{}
	v.Errors.Attrs = make(map[string]map[string]string)
	return v
}
