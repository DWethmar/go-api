package api

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/go-chi/chi"
	cors "github.com/go-chi/cors"
)

// Server provides an entry point for http requests.
type Server struct {
	router *chi.Mux
}

// NewServer creates a new server.
func NewServer(router http.Handler) Server {
	s := Server{
		router: chi.NewRouter(),
	}

	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // you can add routes here www.example.com
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	s.router.Use(logging)
	s.router.Use(cors.Handler)

	s.router.Mount("/", router)

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
