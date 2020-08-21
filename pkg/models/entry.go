package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	MaxNameLength                 = 50
	ErrFieldsRequired             = errors.New("Fields is required.")
	ErrNameLength                 = errors.New("Name exceeded max length of 50.")
	ErrUnsupportedFieldValue      = errors.New("Value is unsupported.")
	ErrUnsupportedFieldSliceValue = errors.New("This array contains one or more unsupported values.")
)

// ID is the type that is used as the identifier for content entries.
type ID = uuid.UUID

// Entry model
type Entry struct {
	ID        ID                `json:"id"   db:"id"`
	Name      string            `json:"name" db:"name"`
	CreatedOn time.Time         `json:"createdOn" db:"created_on"`
	UpdatedOn time.Time         `json:"updatedOn" db:"updated_on"`
	Fields    FieldTranslations `json:"fields" db:"fields"`
}

// AddEntry model
type AddEntry struct {
	Name   string            `json:"name"`
	Fields FieldTranslations `json:"fields"`
}

// UpdateEntry model
type UpdateEntry struct {
	Name   string            `json:"name"`
	Fields FieldTranslations `json:"fields"`
}

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

// CreateEntry creates a new entry.
func CreateEntry() Entry {
	return Entry{
		ID:        createNewID(),
		Name:      "",
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Fields:    FieldTranslations{},
	}
}

type validationError struct {
	Errors struct {
		ID     string                       `json:"id,omitempty"`
		Name   string                       `json:"name,omitempty"`
		Fields map[string]map[string]string `json:"fields,omitempty"`
	} `json:"errors"`
}

func (v validationError) Error() string {
	return ""
}

func isValid(v *validationError) bool {
	if v.Errors.Name != "" || len(v.Errors.Fields) > 0 {
		return false
	}

	return true
}

// Validate addEntry
func (addEntry *AddEntry) Validate() error {
	return validate(&UpdateEntry {
		Name: addEntry.Name,
		Fields: addEntry.Fields,
	});
}

// Validate UpdateEntry
func (updateEntry *UpdateEntry) Validate() error {
	return validate(updateEntry);
}

func validate(entry *UpdateEntry) error {
	vErr := CreateValidationError()
	nameErr := validateName(entry.Name)

	if nameErr != nil {
		vErr.Errors.Name = nameErr.Error()
	}

	for locale, fields := range entry.Fields {
		fieldErrors := validateFields(fields)

		for fieldName, err := range fieldErrors {
			vErr.Errors.Fields[locale][fieldName] = err.Error()
		}
	}

	if (isValid(&vErr)) {
		return nil
	}

	return vErr
}

func validateName(name string) error {
	if len(name) > MaxNameLength {
		return ErrNameLength
	}

	return nil
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

func CreateValidationError() validationError {
	vErr := validationError{}
	vErr.Errors.Fields = make(map[string]map[string]string)
	return vErr
}
