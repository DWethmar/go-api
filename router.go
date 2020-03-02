package main

import (
	"github.com/gorilla/mux"
	
	"github.com/DWethmar/go-api/models"
)

func NewRouter(ds models.Datastore) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc(ds))
	}
	return router
}
