package content

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/dwethmar/go-api/pkg/common"
)

// Content model
type Content struct {
	ID        common.ID
	Name      string
	CreatedOn time.Time
	UpdatedOn time.Time
	Fields    FieldTranslations
}

// FieldTranslations model
type FieldTranslations map[string]Fields

// Fields model
type Fields map[string]interface{}

// Value make the Fields struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a FieldTranslations) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan makes the Fields struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *FieldTranslations) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// Value make the Fields struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a Fields) Value() (driver.Value, error) {
	return json.Marshal(a)
}
