package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

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

func (db *DB) GetAllContentItems() ([]*ContentItem, error) {
	rows, err := db.Query("SELECT * FROM public.content_item")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contentItems := make([]*ContentItem, 0)
	for rows.Next() {
		contentItem := new(ContentItem)
		err := rows.Scan(&contentItem.ID, &contentItem.Name, &contentItem.Data, &contentItem.CreatedOn, &contentItem.UpdatedOn)
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
	err := row.Scan(&contentItem.ID, &contentItem.Name, &contentItem.Data, &contentItem.CreatedOn, &contentItem.UpdatedOn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		panic(err)
	}
	return &contentItem, nil
}

func (db *DB) CreateContentItem(contentItem ContentItem) error {
	sqlStatement := `
	INSERT INTO public.content_item (name, data)
	VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, contentItem.Name, contentItem.Data)
	return err
}
