package main

import (
	"net/http"

	"github.com/DWethmar/go-api/models"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(ds models.Datastore) http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		createIndexHandler,
	},
	Route{
		"Single",
		"GET",
		"/{id}",
		createSingleHandler,
	},
}
