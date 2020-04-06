package contententry

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

// https://medium.com/capital-one-tech/event-sourcing-with-aggregates-in-rust-4022af41cf67
type Aggregate interface {
	GetVersion() int
	Apply(event interface{})
}

type Entry struct {
	ID        ID                `json:"id"   db:"id"`
	Name      string            `json:"name" db:"name"`
	CreatedOn time.Time         `json:"createdOn" db:"created_on"`
	UpdatedOn time.Time         `json:"updatedOn" db:"updated_on"`
	Fields    FieldTranslations `json:"fields" db:"fields"`
}

type AddEntry struct {
	Name   string            `json:"name"`
	Fields FieldTranslations `json:"fields"`
}

// https://www.alexedwards.net/blog/using-postgresql-jsonb
type FieldTranslations map[string]Fields

type Fields map[string]interface{}

// Make the Fields struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a FieldTranslations) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Fields struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *FieldTranslations) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (a Fields) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func CreateEntry() Entry {
	return Entry{
		Fields: FieldTranslations{},
	}
}
