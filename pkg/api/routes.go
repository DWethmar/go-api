package api

import (
	"net/http"

	"github.com/dwethmar/go-api/internal/validator"
	"github.com/dwethmar/go-api/pkg/api/handler"
	"github.com/dwethmar/go-api/pkg/api/middleware"
	"github.com/dwethmar/go-api/pkg/store"

	"github.com/go-chi/chi"
)

/**
TODO
GET /tickets - Retrieves a list of tickets
GET /tickets/12 - Retrieves a specific ticket
POST /tickets - Creates a new ticket
PUT /tickets/12 - Updates ticket #12
PATCH /tickets/12 - Partially updates ticket #12
DELETE /tickets/12 - Deletes ticket #12
**/

// ContentRoutes returns the api routes handler
func ContentRoutes(store *store.Store) http.Handler {
	r := chi.NewRouter()

	handler := handler.NewContentHandler(store.Content, validator.NewValidator())

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

	handler := handler.NewContentTypeHandler(store.ContentType)
	r.Get("/", handler.List)

	return r
}

// NewRouter creates a new api router.
func NewRouter(store *store.Store) http.Handler {
	r := chi.NewRouter()

	r.Mount("/content", ContentRoutes(store))
	r.Mount("/types", ContentTypesRoutes(store))

	return r
}
