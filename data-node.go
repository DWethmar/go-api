package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// NodeData is a basic type.
type NodeData struct {
	ID   int    `json:"id"   db:"id"`
	Name string `json:"name" db:"name"`
	Data Attrs  `json:"text" db:"data"`
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
