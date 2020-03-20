package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/DWethmar/go-api/pkg/contentitem"
	"github.com/gorilla/mux"
)

type Server struct {
	contentItem contentitem.Service
	router      *mux.Router
}

func CreateServer(db *sql.DB) Server {
	s := Server{
		contentItem: contentitem.NewService(contentitem.CreatePostgresRepository(db)),
		router:      mux.NewRouter().StrictSlash(true),
	}
	s.routes()
	s.router.Use(logging)
	return s
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%s %s", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(res, req)
}
