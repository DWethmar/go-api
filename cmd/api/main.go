package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DWethmar/go-api/pkg/config"
	"github.com/DWethmar/go-api/pkg/database"
	"github.com/DWethmar/go-api/pkg/server"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Staring API")
	driverName, dataSource := config.GetPostgresConnectionInfo(config.LoadEnv())
	db, _ := database.ConnectDB(driverName, dataSource)
	defer db.Close()
	server := server.CreateServer(db)
	log.Fatal(http.ListenAndServe(":8080", &server))
}
