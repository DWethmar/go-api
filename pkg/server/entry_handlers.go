package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/contententry"
	"github.com/dwethmar/go-api/pkg/request"
)

type ErrorResponds struct {
	error string
}

func (s *Server) HandleEntryIndex() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if (*r).Method == "OPTIONS" {
			return
		}

		var entries, err = s.entries.GetAll()

		if err != nil {
			fmt.Printf("Error while getting entries: %v", err)
			request.SendServerError(w, r)
			return
		}

		request.SendJSON(w, r, entries, http.StatusOK)
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
			fmt.Printf("Error while decoding entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		v := contententry.ValidateAddEntry(newEntry)

		if !v.IsValid() {
			request.SendBadRequestError(w, r, v)
			return
		}

		entry, err := s.entries.Create(newEntry)

		if err != nil {
			fmt.Printf("Error while creating entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		request.SendJSON(w, r, entry, http.StatusCreated)
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
			fmt.Printf("Error while decoding entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		entry, err := s.entries.GetOne(id)

		if err != nil {
			fmt.Printf("Error while getting entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		entry.Name = newEntry.Name
		entry.Fields = newEntry.Fields
		entry.UpdatedOn = time.Now()

		err = s.entries.Update(*entry)

		if err != nil {
			fmt.Printf("Error while updating entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		request.SendJSON(w, r, entry, http.StatusCreated)
	})
}

func (s *Server) HandleEntryDelete() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())

		if err != nil {
			fmt.Printf("Error while getting id: %v", err)
			request.SendServerError(w, r)
		}

		entry, err := s.entries.GetOne(id)

		if err != nil {
			if err == contententry.ErrNotFound {
				fmt.Printf("Could not find entry: %v", err)
				request.SendNotFoundError(w, r)
				return
			}
			fmt.Printf("Error while Getting entry: %v %v", err, entry)
			request.SendServerError(w, r)
			return
		}

		err = s.entries.Delete(id)

		if err != nil {
			fmt.Printf("Error while deleting entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		request.SendJSON(w, r, entry, http.StatusOK)
	})
}

func (s *Server) HandleEntrySingle() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())
		if err != nil {
			log.Print(fmt.Sprintf("Error on retreiving ID: %v", err))
			request.SendServerError(w, r)
		}

		entry, err := s.entries.GetOne(id)

		if err != nil {
			if err == contententry.ErrNotFound {
				log.Print(fmt.Sprintf("Entry not found: %v", err))
				request.SendNotFoundError(w, r)
				return
			}
			log.Print(fmt.Sprintf("Somthing went wrong: %v", err))
			request.SendServerError(w, r)
			return
		}

		request.SendJSON(w, r, entry, http.StatusOK)
	})
}
