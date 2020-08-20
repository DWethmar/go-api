package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dwethmar/go-api/pkg/api"
	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/database"
	"github.com/dwethmar/go-api/pkg/store"

	_ "github.com/lib/pq"
)

var port = flag.Int("port", 8080, "Run on port")

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

	server := api.CreateServer(api.Routes(store.CreateStore(db)))

	// srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: &server}
	srv := &http.Server{Addr: ":8080", Handler: &server}

	log.Printf("Serving on :%d", port)
	log.Fatal(srv.ListenAndServe())
}
