package server

import "net/http"

func (s *Server) routes() {
	s.router.Name("Index").Methods(http.MethodGet).Path("/").Handler(s.HandleIndex())
	s.router.Name("Single").Methods(http.MethodGet).Path("/{id}").Handler(s.HandleSingle())
	s.router.Name("Delete").Methods(http.MethodDelete).Path("/{id}").Handler(s.HandleDelete())
	s.router.Name("Update").Methods(http.MethodPost).Path("/{id}").Handler(s.HandleUpdate())
	s.router.Name("Create").Methods(http.MethodPost).Path("/").Handler(s.HandleCreate())
}
