package contententry

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ID is the type that is used as the identifier for content entries.
type ID = uuid.UUID

func createNewID() ID {
	return uuid.New()
}

// ParseID test and returns a ID if a string is a valid value.
func ParseID(val string) (ID, error) {
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// Entry model
type Entry struct {
	ID        ID                `json:"id"   db:"id"`
	Name      string            `json:"name" db:"name"`
	CreatedOn time.Time         `json:"createdOn" db:"created_on"`
	UpdatedOn time.Time         `json:"updatedOn" db:"updated_on"`
	Fields    FieldTranslations `json:"fields" db:"fields"`
}

// AddEntry model
type AddEntry struct {
	Name   string            `json:"name"`
	Fields FieldTranslations `json:"fields"`
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

// CreateEntry creates a new entry.
func CreateEntry() Entry {
	return Entry{
		ID:        createNewID(),
		Name:      "",
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Fields:    FieldTranslations{},
	}
}
