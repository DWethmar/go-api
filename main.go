package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// NodeData is a basic type.
type NodeData struct {
	ID   int    `json:"id"   db:"id"`
	Name string `json:"name" db:"name"`
	Data string `json:"text" db:"data"`
}

var db *sql.DB
var err error

func main() {
	fmt.Println("Staring API")

	host := os.Getenv("pgHost")
	port := os.Getenv("pgPort")
	user := os.Getenv("pgUser")
	password := os.Getenv("pgPassword")
	dbname := os.Getenv("pgDbName")
	dbType := os.Getenv("dbType")

	fmt.Println("host = ", host)
	fmt.Println("port = ", port)
	fmt.Println("user = ", user)
	fmt.Println("password = ", password)
	fmt.Println("dbname = ", dbname)
	fmt.Println("dbType = ", dbType)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println(psqlInfo)

	db, err = sql.Open(dbType, psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected with Postgress db!")

	http.HandleFunc("/api/index", indexHandler)
	http.HandleFunc("/api/repo/", repoHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	rows, err := db.Query("SELECT * FROM public.data_node")

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var nodes []NodeData

	defer rows.Close()

	for rows.Next() {
		var node NodeData
		err := rows.Scan(&node.ID, &node.Name, &node.Data)
		if err != nil {
			panic(err.Error())
		}
		nodes = append(nodes, node)
	}

	json.NewEncoder(w).Encode(nodes)
}

func repoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, %q", html.EscapeString(r.URL.Path))
}
