package server

import "net/http"

func (s *Server) routes() {
	s.router.Name("Index").Methods(http.MethodGet).Path("/").Handler(s.HandleEntryIndex())
	s.router.Name("Single").Methods(http.MethodGet).Path("/{id}").Handler(s.HandleEntrySingle())
	s.router.Name("Delete").Methods(http.MethodDelete).Path("/{id}").Handler(s.HandleEntryDelete())
	s.router.Name("Update").Methods(http.MethodPost).Path("/{id}").Handler(s.HandleEntryUpdate())
	s.router.Name("Create").Methods(http.MethodPost).Path("/").Handler(s.HandleEntryCreate())
}
