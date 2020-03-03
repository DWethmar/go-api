package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DWethmar/go-api/internal/store"
	"github.com/DWethmar/go-api/pkg/contentitem"
	_ "github.com/lib/pq"
)

var dataStore store.Datastore

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
	db, err := store.NewDB(driverName, dataSource)
	if err != nil {
		log.Panic(err)
	}

	store := store.Datastore{
		ContentItem: contentitem.CreatePostgresRepository(db),
	}

	defer db.Close()

	fmt.Println("Successfully connected with Postgress db!")
	log.Fatal(http.ListenAndServe(":8080", NewRouter(store)))
}
