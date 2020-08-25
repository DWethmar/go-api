package models

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
)

// FieldTypeDocumentString A document string type
var FieldTypeDocumentString = "document_string"

// FieldTypeDocumentStringList a document string list type
var FieldTypeDocumentStringList = "document_string_list"

// FieldTypeDocumentNumber A document string type
var FieldTypeDocumentNumber = "document_number"

// FieldTypeDocumentNumberList a document string list type
var FieldTypeDocumentNumberList = "document_number_list"

// FieldTypeDocumentInt A document string type
var FieldTypeDocumentInt = "document_int"

// FieldTypeDocumentIntList a document string list type
var FieldTypeDocumentIntList = "document_int_list"

// FieldTypes all field types
var FieldTypes = []string{
	FieldTypeDocumentString,
	FieldTypeDocumentStringList,
	FieldTypeDocumentNumber,
	FieldTypeDocumentNumberList,
	FieldTypeDocumentInt,
	FieldTypeDocumentIntList,
}

// ContentType model
type ContentType struct {
	ID        common.UUID         `json:"id" db:"id"`
	Name      string              `json:"name" db:"name"`
	CreatedOn time.Time           `json:"createdOn" db:"created_on"`
	UpdatedOn time.Time           `json:"updatedOn" db:"updated_on"`
	Fields    []*ContentTypeField `json:"fields"`
}

// ContentTypeField content model field model
type ContentTypeField struct {
	ID           common.UUID `json:"id" db:"id"`
	EntryModelID common.UUID `db:"entry_model_id"`
	Key          string      `json:"key" db:"key"`
	Name         string      `json:"name" db:"name"`
	FieldType    string      `json:"type" db:"type"`
	Length       int         `json:"length" db:"length"`
	CreatedOn    time.Time   `json:"createdOn" db:"created_on"`
	UpdatedOn    time.Time   `json:"updatedOn" db:"updated_on"`
}

// AddContentType model
type AddContentType struct {
	Name   string              `json:"name" db:"name"`
	Fields []*ContentTypeField `json:"fields"`
}

// NewContentType creates new EntryModel
func NewContentType() ContentType {
	return ContentType{
		ID:        common.CreateNewUUID(),
		Name:      "",
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Fields:    make([]*ContentTypeField, 0),
	}
}
