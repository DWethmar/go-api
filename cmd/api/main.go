package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/database"
	"github.com/dwethmar/go-api/pkg/server"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Staring API")
	env := config.LoadEnv()

	driverName, dataSource := config.GetPostgresConnectionInfo(env)
	db, err := database.ConnectDB(driverName, dataSource)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	server := server.CreateServer(db)
	log.Fatal(http.ListenAndServe(":8080", &server))
}
