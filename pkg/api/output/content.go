package output

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/content"
)

// Content model
type Content struct {
	ID        common.ID         `json:"id"`
	Name      string            `json:"name"`
	Fields    FieldTranslations `json:"fields"`
	CreatedOn time.Time         `json:"createdOn"`
	UpdatedOn time.Time         `json:"updatedOn"`
}

// FieldTranslations model
type FieldTranslations map[string]Fields

// Fields model
type Fields map[string]interface{}

// ContentOut maps to output model.
func ContentOut(c *content.Content) *Content {
	ct := &Content{
		ID:        c.ID,
		Name:      c.Name,
		Fields:    FieldTranslations{},
		CreatedOn: c.CreatedOn,
		UpdatedOn: c.UpdatedOn,
	}
	for locale, fields := range c.Fields {
		ct.Fields[locale] = Fields{}
		for key, value := range fields {
			ct.Fields[locale][key] = value
		}
	}
	return ct
}
