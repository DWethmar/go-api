package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DWethmar/go-api/internal/store"
	"github.com/DWethmar/go-api/pkg/contentitem"
	"github.com/gorilla/mux"
)

func IndexHandler(ds store.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data, err = ds.ContentItem.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	})
}

func CreateHandler(ds store.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var contentItem contentitem.ContentItem
		err := decoder.Decode(&contentItem)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		contentItem.CreatedOn = time.Now()
		contentItem.UpdatedOn = time.Now()

		err = ds.ContentItem.Create(&contentItem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

func UpdateHandler(ds store.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var newContentItem contentitem.ContentItem
		err = decoder.Decode(&newContentItem)

		fmt.Println(id)
		fmt.Println(newContentItem)
		fmt.Println(err)
		fmt.Println(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		contentItem, err := ds.ContentItem.GetOne(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		contentItem.Name = newContentItem.Name
		contentItem.Data = newContentItem.Data
		contentItem.UpdatedOn = time.Now()

		err = ds.ContentItem.Update(contentItem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contentItem)
	})
}

func DeleteHandler(ds store.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		contentItem, err := ds.ContentItem.GetOne(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if contentItem == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = ds.ContentItem.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contentItem)
	})
}

func SingleHandler(ds store.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		contentItem, err := ds.ContentItem.GetOne(id)

		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contentItem)
	})
}
