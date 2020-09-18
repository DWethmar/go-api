package input

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/content"
)

// FieldTranslations model
type FieldTranslations map[string]Fields

// Fields model
type Fields map[string]interface{}

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

// MapAddContent to content model
func MapAddContent(add *AddContent) *content.Content {
	c := &content.Content{
		ID:        common.NewID(),
		Name:      add.Name,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Fields:    content.FieldTranslations{},
	}

	for locale, fields := range add.Fields {
		c.Fields[locale] = content.Fields{}
		for key, value := range fields {
			c.Fields[locale][key] = value
		}
	}

	return c
}

// MapUpdateContent to content model
func MapUpdateContent(add *UpdateContent) *content.Content {
	c := &content.Content{
		ID:        common.NewID(),
		Name:      add.Name,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Fields:    content.FieldTranslations{},
	}

	for locale, fields := range add.Fields {
		c.Fields[locale] = content.Fields{}
		for key, value := range fields {
			c.Fields[locale][key] = value
		}
	}

	return c
}
