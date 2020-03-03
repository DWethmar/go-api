package contentitem

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Repository interface {
	GetAll() ([]*ContentItem, error)
	GetOne(id int64) (*ContentItem, error)
	Create(contentItem *ContentItem) error
	Update(contentItem *ContentItem) error
	Delete(id int64) error
}

type ContentItem struct {
	ID        int       `json:"id"   db:"id"`
	Name      string    `json:"name" db:"name"`
	Data      Attrs     `json:"data" db:"data"`
	CreatedOn time.Time `json:"createdOn" db:"created_on"`
	UpdatedOn time.Time `json:"updatedOn" db:"updated_on"`
}

type Attrs map[string]interface{}

func (a Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
