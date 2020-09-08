package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dwethmar/go-api/pkg/api"
	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/database"
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

	files, err := ioutil.ReadDir("/app")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	fmt.Println("-----------------")

	files, err = ioutil.ReadDir("../")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	err = database.RunMigrations(db, c.DBName, "file:///app/migrations")
	if err != nil {
		// panic(err)
		fmt.Println(err)
	}

	server := api.CreateServer(api.NewRouter(store.NewStore(db)))
	srv := &http.Server{Addr: ":8080", Handler: &server}
	log.Printf("Serving on :8080")
	log.Fatal(srv.ListenAndServe())
}
