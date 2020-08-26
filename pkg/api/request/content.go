package request

import "github.com/dwethmar/go-api/pkg/content"

// AddContent model
type AddContent struct {
	Name   string                    `json:"name"`
	Fields content.FieldTranslations `json:"fields"`
}

// UpdateContent model
type UpdateContent struct {
	Name   string                    `json:"name"`
	Fields content.FieldTranslations `json:"fields"`
}
