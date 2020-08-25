package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
)

// FieldTranslations model
type FieldTranslations map[string]Fields

// Fields model
type Fields map[string]interface{}

// Value make the Fields struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a FieldTranslations) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan makes the Fields struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *FieldTranslations) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// Value make the Fields struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a Fields) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func validateFields(fields Fields) map[string]error {
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

				slice := reflect.ValueOf(value)

				for i := 0; i < slice.Len(); i++ {
					interf := slice.Index(i).Interface()
					kind := reflect.TypeOf(interf).Kind()

					if !isValidType(kind) {
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
