package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB
var dbErr error

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

	// Make connection
	db, dbErr = sql.Open(dbType, psqlInfo)
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	// Test connection
	dbErr = db.Ping()
	if dbErr != nil {
		panic(dbErr)
	}

	fmt.Println("Successfully connected with Postgress db!")

	http.HandleFunc("/api/index", indexHandler)
	http.HandleFunc("/api/index/", singleHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func singleHandler(w http.ResponseWriter, r *http.Request) {
	var id, err = URLPathPartInt(r.URL.Path, 2)
	if err != nil {
		http.Error(w, "No ID found!", 400)
		return
	}
	var node NodeData
	row := db.QueryRow(`SELECT * FROM public.data_node WHERE data_node.id = $1`, id)
	err = row.Scan(&node.ID, &node.Name, &node.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No results!", 404)
			return
		}
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}
