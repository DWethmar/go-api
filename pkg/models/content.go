package models

import (
	"errors"
	"time"

	"github.com/dwethmar/go-api/pkg/common"
)

var (
	// MaxNameLength max length of the name field.
	MaxNameLength = 50
	// ErrFieldsRequired error caused by a field that is not added.
	ErrFieldsRequired = errors.New("fields is required")
	// ErrNameLength error caused by a name that is to long.
	ErrNameLength = errors.New("name exceeded max length of 50")
	// ErrUnsupportedFieldValue error caused by an field value that is not supported.
	ErrUnsupportedFieldValue = errors.New("value is unsupported")
	// ErrUnsupportedFieldSliceValue error caused by an invalid value in slice.
	ErrUnsupportedFieldSliceValue = errors.New("this array contains one or more unsupported values")
)

// Content model
type Content struct {
	ID        common.UUID       `json:"id"   db:"id"`
	Name      string            `json:"name" db:"name"`
	CreatedOn time.Time         `json:"createdOn" db:"created_on"`
	UpdatedOn time.Time         `json:"updatedOn" db:"updated_on"`
	Fields    FieldTranslations `json:"fields" db:"fields"`
}

// AddContent model
type AddContent struct {
	Name   string            `json:"name"`
	Fields FieldTranslations `json:"fields"`
}

// UpdateContent model
type UpdateContent struct {
	Name   string            `json:"name"`
	Fields FieldTranslations `json:"fields"`
}

// NewContent creates a new entry.
func NewContent() *Content {
	return &Content{
		ID:        common.CreateNewUUID(),
		Name:      "",
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Fields:    FieldTranslations{},
	}
}

// ErrContentValidation error containing infomation about validation errors
type ErrContentValidation struct {
	Errors struct {
		ID     string                       `json:"id,omitempty"`
		Name   string                       `json:"name,omitempty"`
		Fields map[string]map[string]string `json:"fields,omitempty"`
	} `json:"errors"`
}

func (v ErrContentValidation) Error() string {
	return ""
}

func isValid(v *ErrContentValidation) bool {
	if v.Errors.Name != "" || len(v.Errors.Fields) > 0 {
		return false
	}

	return true
}

// Validate AddContent
func (addContent *AddContent) Validate() error {
	return validate(&UpdateContent{
		Name:   addContent.Name,
		Fields: addContent.Fields,
	})
}

// Validate UpdateContent
func (content *UpdateContent) Validate() error {
	return validate(content)
}

func validate(content *UpdateContent) error {
	vErr := CreateContentValidationError()
	nameErr := validateName(content.Name)

	if nameErr != nil {
		vErr.Errors.Name = nameErr.Error()
	}

	for locale, fields := range content.Fields {
		fieldErrors := validateFields(fields)

		for fieldName, err := range fieldErrors {
			if vErr.Errors.Fields[locale] == nil {
				vErr.Errors.Fields[locale] = make(map[string]string)
			}
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

// CreateContentValidationError create new validation error.
func CreateContentValidationError() ErrContentValidation {
	vErr := ErrContentValidation{}
	vErr.Errors.Fields = make(map[string]map[string]string)
	return vErr
}
