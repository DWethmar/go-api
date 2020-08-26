package api

import (
	"net/http"

	"github.com/dwethmar/go-api/pkg/api/handler"
	"github.com/dwethmar/go-api/pkg/api/middleware"
	"github.com/dwethmar/go-api/pkg/store"

	"github.com/go-chi/chi"
)

// ContentRoutes returns the api routes handler
func ContentRoutes(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Get("/", handler.ContentIndex(store))
	r.Get("/{id}", middleware.RequireID(handler.GetSingleContent(store)))
	r.Delete("/{id}", middleware.RequireID(handler.DeleteContent(store)))
	r.Post("/{id}", middleware.RequireID(handler.UpdateContent(store)))
	r.Post("/", handler.CreateContent(store))

	return r
}

// ContentTypesRoutes returns the api routes handler
func ContentTypesRoutes(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Get("/", handler.ContentTypeIndex(store))

	return r
}

// NewRouter creates a new api router.
func NewRouter(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Mount("/content", ContentRoutes(store))
	r.Mount("/type", ContentTypesRoutes(store))

	return r
}
