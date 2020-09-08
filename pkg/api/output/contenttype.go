package output

import (
	"time"

	"github.com/dwethmar/go-api/pkg/contenttype"
)

// ContentType model
type ContentType struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	CreatedOn time.Time            `json:"createdOn"`
	UpdatedOn time.Time            `json:"updatedOn"`
	Fields    []*contenttype.Field `json:"fields"` // TODO: Also cast this.
}
