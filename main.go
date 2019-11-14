package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "go-api"
)

// NodeData is a basic type.
type NodeData struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Data string `json:"text" db:"data"`
}

var db *sql.DB

func main() {
	fmt.Println("Staring API")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
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

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM public.data_node")

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	dataNodes := []NodeData{}

	for rows.Next() {
		var u NodeData
		if err := rows.Scan(&u.ID, &u.Name, &u.Data); err != nil {
			return
		}
		dataNodes = append(dataNodes, u)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataNodes)
}

func repoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, %q", html.EscapeString(r.URL.Path))
}
