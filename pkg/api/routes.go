package api

import (
	"net/http"

	"github.com/dwethmar/go-api/pkg/middleware"
	"github.com/dwethmar/go-api/pkg/services"

	"github.com/go-chi/chi"
)

// EntryRoutes returns the api routes handler
func EntryRoutes(store *services.Store) http.Handler {
	r := chi.NewRouter()

	r.Get("/", middleware.EntryIndex(store))
	r.Get("/{id}", RequireEntryID(middleware.HandleEntrySingle(store)))
	r.Delete("/{id}", RequireEntryID(middleware.DeleteEntry(store)))
	r.Post("/{id}", RequireEntryID(middleware.UpdateEntry(store)))
	r.Post("/", middleware.CreateEntry(store))

	return r
}
