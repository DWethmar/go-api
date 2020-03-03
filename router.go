package main

import (
	"github.com/DWethmar/go-api/internal/store"
	"github.com/gorilla/mux"
)

func NewRouter(ds store.Datastore) *mux.Router {
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
