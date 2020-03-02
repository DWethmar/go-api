package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DWethmar/go-api/models"
	_ "github.com/lib/pq"
)

type Env struct {
	db models.Datastore
}

var env Env

func main() {
	fmt.Println("Staring API")

	host := os.Getenv("pgHost")
	port := os.Getenv("pgPort")
	user := os.Getenv("pgUser")
	password := os.Getenv("pgPassword")
	dbName := os.Getenv("pgDbName")
	driverName := os.Getenv("driverName")
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	// Make connection
	db, err := models.NewDB(driverName, dataSource)
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db}

	defer db.Close()

	fmt.Println("Successfully connected with Postgress db!")
	log.Fatal(http.ListenAndServe(":8080", NewRouter(env.db)))
}
