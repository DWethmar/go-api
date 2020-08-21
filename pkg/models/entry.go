package models

import (
	"errors"
	"time"

	"github.com/dwethmar/go-api/pkg/common"
)

var (
	MaxNameLength                 = 50
	ErrFieldsRequired             = errors.New("Fields is required.")
	ErrNameLength                 = errors.New("Name exceeded max length of 50.")
	ErrUnsupportedFieldValue      = errors.New("Value is unsupported.")
	ErrUnsupportedFieldSliceValue = errors.New("This array contains one or more unsupported values.")
)

type ID = common.UUID

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

// CreateEntry creates a new entry.
func CreateEntry() Entry {
	return Entry{
		ID:        common.CreateNewUUID(),
		Name:      "",
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Fields:    FieldTranslations{},
	}
}

type EntryValidationError struct {
	Errors struct {
		ID     string                       `json:"id,omitempty"`
		Name   string                       `json:"name,omitempty"`
		Fields map[string]map[string]string `json:"fields,omitempty"`
	} `json:"errors"`
}

func (v EntryValidationError) Error() string {
	return ""
}

func isValid(v *EntryValidationError) bool {
	if v.Errors.Name != "" || len(v.Errors.Fields) > 0 {
		return false
	}

	return true
}

// Validate addEntry
func (addEntry *AddEntry) Validate() error {
	return validate(&UpdateEntry{
		Name:   addEntry.Name,
		Fields: addEntry.Fields,
	})
}

// Validate UpdateEntry
func (updateEntry *UpdateEntry) Validate() error {
	return validate(updateEntry)
}

func validate(entry *UpdateEntry) error {
	vErr := CreateEntryValidationError()
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

	if isValid(&vErr) {
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

func CreateEntryValidationError() EntryValidationError {
	vErr := EntryValidationError{}
	vErr.Errors.Fields = make(map[string]map[string]string)
	return vErr
}
