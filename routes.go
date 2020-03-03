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
		IndexHandler,
	},
	Route{
		"Create",
		"POST",
		"/",
		CreateHandler,
	},
	Route{
		"Single",
		"GET",
		"/{id}",
		SingleHandler,
	},
}
