package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/DWethmar/go-api/pkg/contentitem"
	"github.com/gorilla/mux"
)

type ErrorResponds struct {
	error string
}

func (s *Server) HandleIndex() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data, err = s.contentItem.GetAll()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	})
}

func (s *Server) HandleCreate() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var newContentItem contentitem.NewContentItem
		err := decoder.Decode(&newContentItem)

		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponds{
				error: "bad request!",
			})
			return
		}

		contentItem, err := s.contentItem.Create(newContentItem)

		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(contentItem)
	})
}

func (s *Server) HandleUpdate() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var newContentItem contentitem.ContentItem
		err = decoder.Decode(&newContentItem)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		contentItem, err := s.contentItem.GetOne(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		contentItem.Name = newContentItem.Name
		contentItem.Attrs = newContentItem.Attrs
		contentItem.UpdatedOn = time.Now()

		err = s.contentItem.Update(contentItem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contentItem)
	})
}

func (s *Server) HandleDelete() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		contentItem, err := s.contentItem.GetOne(id)
		if err != nil {
			if err == contentitem.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = s.contentItem.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contentItem)
	})
}

func (s *Server) HandleSingle() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		contentItem, err := s.contentItem.GetOne(id)

		if err != nil {
			if err == contentitem.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contentItem)
	})
}
