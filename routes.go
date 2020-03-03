package main

import (
	"net/http"

	"github.com/DWethmar/go-api/internal/store"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(ds store.Datastore) http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		IndexHandler,
	},
	Route{
		"Create",
		"POST",
		"/",
		CreateHandler,
	},
	Route{
		"Update",
		"POST",
		"/{id}",
		UpdateHandler,
	},
	Route{
		"Delete",
		"DELETE",
		"/{id}",
		DeleteHandler,
	},
	Route{
		"Single",
		"GET",
		"/{id}",
		SingleHandler,
	},
}
