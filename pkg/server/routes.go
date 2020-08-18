package server

import (
	"github.com/dwethmar/go-api/pkg/handler"
	"github.com/dwethmar/go-api/pkg/store"

	"github.com/go-chi/chi"
)

func Router(store *store.Store) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.HandleEntryIndex(store))
	r.Get("/{id}", requireEntryId(handler.HandleEntrySingle(store)))
	r.Delete("/{id}", requireEntryId(handler.HandleEntryDelete(store)))
	r.Post("/{id}", requireEntryId(handler.HandleEntryUpdate(store)))
	r.Post("/", handler.HandleEntryCreate(store))

	return r
}
