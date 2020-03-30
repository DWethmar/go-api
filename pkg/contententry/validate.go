package contententry

import (
	"errors"
	"reflect"
)

var (
	MaxNameLength                 = 50
	ErrFieldsRequired             = errors.New("Fields is required.")
	ErrNameLength                 = errors.New("Name exceeded max length of 50.")
	ErrUnsupportedFieldValue      = errors.New("Value is unsupported.")
	ErrUnsupportedFieldSliceValue = errors.New("This array contains one or more unsupported values.")
)

type validationResult struct {
	Errors struct {
		Name   string                       `json:"name,omitempty"`
		Fields map[string]map[string]string `json:"fields,omitempty"`
	} `json:"errors"`
}

func (v *validationResult) IsValid() bool {
	if v.Errors.Name != "" || len(v.Errors.Fields) > 0 {
		return false
	}

	return true
}

func ValidateAddEntry(addEntry AddEntry) validationResult {
	e := CreateValidationResult()

	nameErr := ValidateName(addEntry.Name)
	if nameErr != nil {
		e.Errors.Name = nameErr.Error()
	}

	for locale, fields := range addEntry.Fields {
		attrErRors := ValidateAttr(fields)

		for attr, error := range attrErRors {
			if e.Errors.Fields[locale] == nil {
				e.Errors.Fields[locale] = map[string]string{}
			}

			e.Errors.Fields[locale][attr] = error.Error()
		}

	}
	return e
}

func ValidateContentItem(entry Entry) validationResult {
	e := CreateValidationResult()

	nameErr := ValidateName(entry.Name)

	if nameErr != nil {
		e.Errors.Name = nameErr.Error()
	}

	for locale, fields := range entry.Fields {
		attrErRors := ValidateAttr(fields)

		for attr, error := range attrErRors {
			e.Errors.Fields[locale][attr] = error.Error()
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

func ValidateAttr(fields Fields) map[string]error {
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

	for attr, value := range fields {
		t := reflect.TypeOf(value)
		if t == nil {
			e[attr] = ErrUnsupportedFieldValue
			continue
		}

		kind := t.Kind()

		if kind == reflect.Invalid {
			e[attr] = ErrUnsupportedFieldValue
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
						e[attr] = ErrUnsupportedFieldSliceValue
						break
					}
				}
			} else {
				e[attr] = ErrUnsupportedFieldValue
			}
		}
	}

	return e
}

func CreateValidationResult() validationResult {
	v := validationResult{}
	v.Errors.Fields = make(map[string]map[string]string)
	return v
}
