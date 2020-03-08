package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DWethmar/go-api/config"
	"github.com/DWethmar/go-api/server"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Staring API")
	driverName, dataSource := config.GetPostgresConnectionInfo(config.LoadEnv())
	db, _ := server.NewDB(driverName, dataSource)
	defer db.Close()
	server := server.CreateServer(db)
	log.Fatal(http.ListenAndServe(":8080", &server))
}
