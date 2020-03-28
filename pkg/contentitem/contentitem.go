package contentitem

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type ID = uuid.UUID

func createNewId() ID {
	return uuid.New()
}

func ParseId(val string) (ID, error) {
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

type ContentItem struct {
	ID        ID           `json:"id"   db:"id"`
	Name      string       `json:"name" db:"name"`
	CreatedOn time.Time    `json:"createdOn" db:"created_on"`
	UpdatedOn time.Time    `json:"updatedOn" db:"updated_on"`
	Attrs     AttrsLocales `json:"attrs" db:"attrs"`
}

type AddContentItem struct {
	Name  string       `json:"name" db:"name"`
	Attrs AttrsLocales `json:"attrs"`
}

// https://www.alexedwards.net/blog/using-postgresql-jsonb
type AttrsLocales map[string]Attrs

type Attrs map[string]interface{}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a AttrsLocales) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *AttrsLocales) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (a Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func CreateContentItem() ContentItem {
	return ContentItem{
		Attrs: AttrsLocales{},
	}
}
