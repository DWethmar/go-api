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

	handler := handler.NewContentHandler(store.Content)

	r.Get("/", handler.List)
	r.Get("/{id}", middleware.RequireID(handler.Get))
	r.Delete("/{id}", middleware.RequireID(handler.Delete))
	r.Post("/{id}", middleware.RequireID(handler.Update))
	r.Post("/", handler.Create)

	return r
}

// ContentTypesRoutes returns the api routes handler
func ContentTypesRoutes(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Get("/", handler.ListContentType(store))

	return r
}

// NewRouter creates a new api router.
func NewRouter(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Mount("/content", ContentRoutes(store))
	r.Mount("/types", ContentTypesRoutes(store))

	return r
}
