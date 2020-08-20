package api

import (
	"net/http"

	"github.com/dwethmar/go-api/pkg/middleware"
	"github.com/dwethmar/go-api/pkg/store"

	"github.com/go-chi/chi"
)

// Routes returns the api routes handler
func Routes(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Get("/", middleware.EntryIndex(store))
	r.Get("/{id}", middleware.RequireEntryID(middleware.HandleEntrySingle(store)))
	r.Delete("/{id}", middleware.RequireEntryID(middleware.DeleteEntry(store)))
	r.Post("/{id}", middleware.RequireEntryID(middleware.UpdateEntry(store)))
	r.Post("/", middleware.CreateEntry(store))

	return r
}
