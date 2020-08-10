package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/database"
	"github.com/dwethmar/go-api/pkg/server"

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

	server := server.CreateServer(db)

	srv := &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: &server}
	log.Printf("Serving on :%v", port)
	log.Fatal(srv.ListenAndServe())
}
