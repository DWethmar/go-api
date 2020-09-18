package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dwethmar/go-api/internal/database"
	"github.com/dwethmar/go-api/pkg/api"
	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/store"

	_ "github.com/lib/pq"
)

func main() {
	flag.Parse()

	fmt.Println("Staring API")

	c := config.LoadEnvFile()
	driverName, dataSource := config.GetPostgresConnectionInfo(c)

	db, err := database.NewDB(driverName, dataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = database.RunMigrations(db, c.DBName, c.MigrationFiles, c.DBMigrationVersion)
	if err != nil {
		fmt.Println(err)
	}

	server := api.NewServer(api.NewRouter(store.NewStore(db)))
	srv := &http.Server{Addr: ":8080", Handler: &server}
	log.Printf("Serving on :8080")
	log.Fatal(srv.ListenAndServe())
}
