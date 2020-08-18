package server

import (
	"github.com/dwethmar/go-api/pkg/handler"
	"github.com/dwethmar/go-api/pkg/store"

	"github.com/go-chi/chi"
)

func Router(store *store.Store) chi.Router {
	r := chi.NewRouter()

	r.Get("/content", handler.HandleEntryIndex(store))

	return r
}
