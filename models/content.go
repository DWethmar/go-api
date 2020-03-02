package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ContentItem struct {
	ID   int    `json:"id"   db:"id"`
	Name string `json:"name" db:"name"`
	Data Attrs  `json:"data" db:"data"`
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

func (db *DB) GetAllContentItems() ([]*ContentItem, error) {
	rows, err := db.Query("SELECT * FROM public.content_item")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contentItems := make([]*ContentItem, 0)
	for rows.Next() {
		contentItem := new(ContentItem)
		err := rows.Scan(&contentItem.ID, &contentItem.Name, &contentItem.Data)
		if err != nil {
			return nil, err
		}
		contentItems = append(contentItems, contentItem)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return contentItems, nil
}

func (db *DB) GetOneContentItem(id int64) (*ContentItem, error) {
	var contentItem ContentItem
	row := db.QueryRow(`SELECT * FROM public.content_item WHERE content_item.id = $1`, id)
	err := row.Scan(&contentItem.ID, &contentItem.Name, &contentItem.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		panic(err)
	}
	return &contentItem, nil
}
