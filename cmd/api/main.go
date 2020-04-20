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

var httpVersion = flag.Int("http", 2, "HTTP version")

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

	srv := &http.Server{Addr: ":8080", Handler: &server}
	log.Printf("Serving on :8080 with http %v", *httpVersion)
	// log.Fatal(srv.ListenAndServeTLS("cert/server.crt", "cert/server.key"))
	log.Fatal(srv.ListenAndServe())
	// 	log.Fatal(http.ListenAndServeTLS(":8080", "cert/server.crt", "cert/server.key", &server))
}
