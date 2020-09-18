package output

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/contenttype"
)

// ContentType model
type ContentType struct {
	ID        common.ID           `json:"id"`
	Name      string              `json:"name"`
	CreatedOn time.Time           `json:"createdOn"`
	UpdatedOn time.Time           `json:"updatedOn"`
	Fields    []*ContentTypeField `json:"fields"`
}

// ContentTypeField model
type ContentTypeField struct {
	ID            common.ID `json:"id"`
	ContentTypeID common.ID `json:"ContentTypeI"`
	Key           string    `json:"key"`
	Name          string    `json:"name"`
	FieldType     string    `json:"type"`
	Length        int       `json:"length"`
	CreatedOn     time.Time `json:"createdOn"`
	UpdatedOn     time.Time `json:"updatedOn"`
}

// ContentTypeOut maps to output model.
func ContentTypeOut(c *contenttype.ContentType) *ContentType {
	ct := &ContentType{
		ID:        c.ID,
		Name:      c.Name,
		CreatedOn: c.CreatedOn,
		UpdatedOn: c.UpdatedOn,
		Fields:    []*ContentTypeField{},
	}
	for _, field := range c.Fields {
		ct.Fields = append(ct.Fields, &ContentTypeField{
			ID:            field.ID,
			ContentTypeID: field.ContentTypeID,
			Key:           field.Key,
			Name:          field.Name,
			FieldType:     field.FieldType,
			Length:        field.Length,
			CreatedOn:     field.CreatedOn,
			UpdatedOn:     field.UpdatedOn,
		})
	}
	return ct
}
