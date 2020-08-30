package response

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/content"
)

// Content model
type Content struct {
	ID        common.ID                 `json:"id"`
	Name      string                    `json:"name"`
	Fields    content.FieldTranslations `json:"fields"`
	CreatedOn time.Time                 `json:"createdOn"`
	UpdatedOn time.Time                 `json:"updatedOn"`
}
