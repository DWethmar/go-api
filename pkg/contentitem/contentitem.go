package contentitem

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ContentItem struct {
	ID        int       `json:"id"   db:"id"`
	Name      string    `json:"name" db:"name"`
	Attrs     Attrs     `json:"attrs" db:"attrs"`
	CreatedOn time.Time `json:"createdOn" db:"created_on"`
	UpdatedOn time.Time `json:"updatedOn" db:"updated_on"`
}

type NewContentItem struct {
	Name  string `json:"name" db:"name"`
	Attrs Attrs  `json:"data" db:"attrs"`
}

// https://www.alexedwards.net/blog/using-postgresql-jsonb
type Attrs map[string]interface{}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
