package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DWethmar/go-api/pkg/contentitem"
	"github.com/gorilla/mux"
)

type ErrorResponds struct {
	error string
}

func (s *Server) HandleContentItemIndex() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data, err = s.contentItem.GetAll()

		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err == contentitem.ErrNotFound {
				json.NewEncoder(w).Encode([]string{})
			} else {
				json.NewEncoder(w).Encode(err.Error())
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	})
}

func (s *Server) HandleContentItemCreate() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		newContentItem := contentitem.AddContentItem{
			Attrs: make(contentitem.AttrsLocales),
		}
		err := decoder.Decode(&newContentItem)

		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponds{
				error: "something went wrong..",
			})
			return
		}

		v := contentitem.ValidateAddContentItem(newContentItem)
		if !v.IsValid() {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(v)
			return
		}

		contentItem, err := s.contentItem.Create(newContentItem)

		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err == contentitem.ErrNotFound {
				json.NewEncoder(w).Encode("Could not find contentItem.")
			} else {
				json.NewEncoder(w).Encode(err.Error())
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(contentItem)
	})
}

func (s *Server) HandleContentItemUpdate() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := contentitem.ParseId(vars["id"])

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

		err = s.contentItem.Update(*contentItem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contentItem)
	})
}

func (s *Server) HandleContentItemDelete() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := contentitem.ParseId(vars["id"])
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

func (s *Server) HandleContentItemSingle() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := contentitem.ParseId(vars["id"])

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
