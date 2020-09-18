package input

import "github.com/dwethmar/go-api/pkg/content"

// AddContent model
type AddContentType struct {
	Name   string                    `json:"name"`
	Fields content.FieldTranslations `json:"fields"`
}

// UpdateContent model
type UpdateContentType struct {
	Name   string                    `json:"name"`
	Fields content.FieldTranslations `json:"fields"`
}
