package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Datastore interface {
	GetAllContentItems() ([]*ContentItem, error)
	GetOneContentItem(id int64) (*ContentItem, error)
}

type DB struct {
	*sql.DB
}

func NewDB(driverName string, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
