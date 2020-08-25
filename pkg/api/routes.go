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

// ContentTypesRoutes returns the api routes handler
func ContentTypesRoutes(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Get("/", middleware.ContentTypeIndex(store))

	return r
}

// NewRouter creates a new api router.
func NewRouter(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Mount("/content", ContentRoutes(store))
	r.Mount("/model", ContentTypesRoutes(store))

	return r
}
