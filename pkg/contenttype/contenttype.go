package contenttype

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
)

// ContentType model
type ContentType struct {
	ID        common.ID
	Name      string
	CreatedOn time.Time
	UpdatedOn time.Time
	Fields    []*Field
}

// Field content model field model
type Field struct {
	ID           common.ID
	EntryModelID common.ID
	Key          string
	Name         string
	FieldType    string
	Length       int
	CreatedOn    time.Time
	UpdatedOn    time.Time
}
