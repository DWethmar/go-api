package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dwethmar/go-api/pkg/api"
	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/database"
	"github.com/dwethmar/go-api/pkg/services"

	_ "github.com/lib/pq"
)

func main() {
	flag.Parse()

	fmt.Println("Staring API")

	env := config.LoadEnvFile()
	driverName, dataSource := config.GetPostgresConnectionInfo(env)

	db, err := database.ConnectDB(driverName, dataSource)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	server := api.CreateServer(api.EntryRoutes(services.CreateStore(db)))
	srv := &http.Server{Addr: ":8080", Handler: &server}
	log.Printf("Serving on :8080")
	log.Fatal(srv.ListenAndServe())
}
