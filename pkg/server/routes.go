package server

import "net/http"

func (s *Server) routes() {
	s.router.Name("Index").Methods(http.MethodGet).Path("/").Handler(s.HandleContentItemIndex())
	s.router.Name("Single").Methods(http.MethodGet).Path("/{id}").Handler(s.HandleContentItemSingle())
	s.router.Name("Delete").Methods(http.MethodDelete).Path("/{id}").Handler(s.HandleContentItemDelete())
	s.router.Name("Update").Methods(http.MethodPost).Path("/{id}").Handler(s.HandleContentItemUpdate())
	s.router.Name("Create").Methods(http.MethodPost).Path("/").Handler(s.HandleContentItemCreate())
}
