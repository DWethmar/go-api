package api

import (
	"net/http"

	"github.com/dwethmar/go-api/pkg/middleware"
	"github.com/dwethmar/go-api/pkg/store"

	"github.com/go-chi/chi"
)

// ContentRoutes returns the api routes handler
func ContentRoutes(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Get("/", middleware.ContentIndex(store))
	r.Get("/{id}", RequireEntryID(middleware.GetSingleContent(store)))
	r.Delete("/{id}", RequireEntryID(middleware.DeleteContent(store)))
	r.Post("/{id}", RequireEntryID(middleware.UpdateContent(store)))
	r.Post("/", middleware.CreateContent(store))

	return r
}

// ContentModelsRoutes returns the api routes handler
func ContentModelsRoutes(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Get("/", middleware.ContentModelIndex(store))

	return r
}

// NewRouter creates a new api router.
func NewRouter(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Mount("/content", ContentRoutes(store))
	r.Mount("/model", ContentModelsRoutes(store))

	return r
}
