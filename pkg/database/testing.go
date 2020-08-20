package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
)

func CreateDatabase(db *sql.DB, name string) {
	fmt.Printf("Creating database: %v \n", name)
	_, err := db.Exec("CREATE DATABASE " + name + ";")

	if err != nil {
		panic(err)
	}
}

func DropDatabase(db *sql.DB, name string) {
	fmt.Printf("Dropping database: %v \n", name)
	_, err := db.Exec("DROP DATABASE " + name)

	if err != nil {
		panic(err)
	}
}

func ExecSQLFile(db *sql.DB, sqlFile string) {

	fmt.Printf("Reading SQL file: %v \n", sqlFile)

	b, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatal(err)
	}

	sql := string(b)

	fmt.Printf("Excecuting SQL file: %v \n", sqlFile)

	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}
}
