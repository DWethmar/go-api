package contentitem

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var MaxNameLength = 50

var (
	NameRequiredError    = errors.New("Name is required.")
	NameLengthError      = errors.New(fmt.Sprintf("Name exceeded max length of %d.", MaxNameLength))
	AttrsRequiredError   = errors.New("Attrs is required.")
	UnknownAttrTypeError = errors.New("Attrs has an unknown type.")
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

func (c *ContentItem) Validate() map[string][]error {
	e := map[string][]error{}

	addErr := func(attr string, err error) {
		if e[attr] == nil {
			e[attr] = []error{}
		}
		e[attr] = append(e[attr], err)
	}

	if c.Name == "" {
		addErr("Name", NameRequiredError)
	}
	if len(c.Name) > MaxNameLength {
		addErr("Name", NameLengthError)
	}

	ct := func(attrs map[string]interface{}) []error {
		ce := make([]error, 0)
		for _, value := range attrs {
			switch v := value.(type) {
			case float64:
				fmt.Println("float64:", v)
			case int:
				fmt.Println("int:", v)
			case []int:
				fmt.Println("[]int:", v)
			case string:
				fmt.Println("string:", v)
			case []string:
				fmt.Println("[]string:", v)
			case bool:
				fmt.Println("bool:", v)
			default:
				fmt.Println("unknown")
				addErr("Name", UnknownAttrTypeError)
			}
		}
		return ce
	}
	ct(c.Attrs)
	return e
}
