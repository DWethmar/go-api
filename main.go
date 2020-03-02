package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB
var dbErr error

func main() {
	fmt.Println("Staring API")

	host := os.Getenv("pgHost")
	port := os.Getenv("pgPort")
	user := os.Getenv("pgUser")
	password := os.Getenv("pgPassword")
	dbName := os.Getenv("pgDbName")
	dbType := os.Getenv("dbType")

	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	// Make connection
	db, dbErr = sql.Open(dbType, dataSource)
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	// Test connection
	dbErr = db.Ping()
	if dbErr != nil {
		panic(dbErr)
	}

	fmt.Println("Successfully connected with Postgress db!")

	log.Fatal(http.ListenAndServe(":8080", NewRouter()))
}
