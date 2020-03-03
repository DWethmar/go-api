package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DWethmar/go-api/models"
	"github.com/gorilla/mux"
)

func IndexHandler(ds models.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data, err = ds.GetAllContentItems()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	})
}

func CreateHandler(ds models.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var contentItem models.ContentItem
		err := decoder.Decode(&contentItem)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = ds.CreateContentItem(contentItem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contentItem)
	})
}

func SingleHandler(ds models.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, err := ds.GetOneContentItem(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	})
}
