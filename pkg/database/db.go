package database

import (
	"database/sql"
	"fmt"
)

type DB struct {
	*sql.DB
}

func ConnectDB(driverName string, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected with Postgress db!")
	return db, nil
}
