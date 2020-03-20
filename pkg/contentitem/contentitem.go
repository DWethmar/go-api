package contentitem

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var MaxNameLength = 50

type UnsupportedAttrTypeError struct {
	Attr string
}

func (e *UnsupportedAttrTypeError) Error() string {
	return fmt.Sprintf("Attribute type of %v is not supported.", e.Attr)
}

type NameLengthError struct{}

func (e *NameLengthError) Error() string {
	return fmt.Sprintf("Name exceeded max length of %d.", MaxNameLength)
}

var (
	AttrsRequiredError = errors.New("Attrs is required.")
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

func (c *ContentItem) Validate() []error {
	var e []error

	if len(c.Name) > MaxNameLength {
		e = append(e, &NameLengthError{})
	}

	ct := func(attrs map[string]interface{}) []error {
		ce := make([]error, 0)
		for key, value := range attrs {
			switch t := value.(type) {
			case float64:
				fmt.Printf("%v is of type float64: %v\n", key, t)
			case int:
				fmt.Printf("%v is of type int: %v\n", key, t)
			case []int:
				fmt.Printf("%v is of type []int: %v\n", key, t)
			case string:
				fmt.Printf("%v is of type string: %v\n", key, t)
			case []string:
				fmt.Printf("%v is of type []string: %v\n", key, t)
			case bool:
				fmt.Printf("%v is of type bool: %v\n", key, t)
			default:
				e = append(e, &UnsupportedAttrTypeError{
					Attr: key,
				})
			}
		}
		return ce
	}
	ct(c.Attrs)
	return e
}
