package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DWethmar/go-api/pkg/common"
	"github.com/DWethmar/go-api/pkg/contententry"
	"github.com/DWethmar/go-api/pkg/request"
)

type ErrorResponds struct {
	error string
}

func (s *Server) HandleEntryIndex() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data, err = s.entries.GetAll()

		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err == contententry.ErrNotFound {
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

func (s *Server) HandleEntryCreate() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		newEntry := contententry.AddEntry{
			Fields: make(contententry.FieldTranslations),
		}
		err := decoder.Decode(&newEntry)

		if err != nil {
			request.SendServerError(w, r)
			return
		}

		v := contententry.ValidateAddEntry(newEntry)
		if !v.IsValid() {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(v)
			return
		}

		entry, err := s.entries.Create(newEntry)

		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err == contententry.ErrNotFound {
				json.NewEncoder(w).Encode("Could not find Entry.")
			} else {
				json.NewEncoder(w).Encode(err.Error())
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(entry)
	})
}

func (s *Server) HandleEntryUpdate() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())
		if err != nil {
			request.SendServerError(w, r)
		}

		decoder := json.NewDecoder(r.Body)
		var newEntry contententry.Entry
		err = decoder.Decode(&newEntry)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		entry, err := s.entries.GetOne(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		entry.Name = newEntry.Name
		entry.Fields = newEntry.Fields
		entry.UpdatedOn = time.Now()

		err = s.entries.Update(*entry)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entry)
	})
}

func (s *Server) HandleEntryDelete() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())
		if err != nil {
			request.SendServerError(w, r)
		}
		Entry, err := s.entries.GetOne(id)
		if err != nil {
			if err == contententry.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = s.entries.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Entry)
	})
}

func (s *Server) HandleEntrySingle() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())
		if err != nil {
			log.Print(fmt.Sprintf("Error on retreiving ID: %v", err))
			request.SendServerError(w, r)
		}

		Entry, err := s.entries.GetOne(id)

		if err != nil {
			if err == contententry.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Entry)
	})
}
