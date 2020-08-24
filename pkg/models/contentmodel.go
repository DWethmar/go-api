package models

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
)

// ContentModel model
type ContentModel struct {
	ID        common.UUID          `json:"id"   db:"id"`
	Name      string               `json:"name" db:"name"`
	CreatedOn time.Time            `json:"createdOn" db:"created_on"`
	UpdatedOn time.Time            `json:"updatedOn" db:"updated_on"`
	Fields    []*ContentModelField `json:"fields"`
}

// ContentModelField content model field model
type ContentModelField struct {
	EntryModelID common.UUID `db:"entry_model_id"`
	Key          string      `json:"key" db:"key"`
	Name         string      `json:"name" db:"name"`
	FieldType    string      `json:"type" db:"type"`
	Length       int         `json:"length" db:"length"`
	CreatedOn    time.Time   `json:"createdOn" db:"created_on"`
	UpdatedOn    time.Time   `json:"updatedOn" db:"updated_on"`
}

// AddContentModel model
type AddContentModel struct {
	Name   string               `json:"name" db:"name"`
	Fields []*ContentModelField `json:"fields"`
}

// NewContentModel creates new EntryModel
func NewContentModel() ContentModel {
	return ContentModel{
		ID:        common.CreateNewUUID(),
		Name:      "",
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Fields:    make([]*ContentModelField, 0),
	}
}
