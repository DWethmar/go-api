package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM public.data_node")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var data []Data
	defer rows.Close()
	for rows.Next() {
		var node Data
		err := rows.Scan(&node.ID, &node.Name, &node.Data)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, node)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func singleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		http.Error(w, "No ID found!", 400)
		return
	}

	var data Data
	row := db.QueryRow(`SELECT * FROM public.data_node WHERE data_node.id = $1`, id)
	err := row.Scan(&data.ID, &data.Name, &data.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No results!", 404)
			return
		}
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
