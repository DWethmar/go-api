package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/models"
	"github.com/dwethmar/go-api/pkg/request"
	"github.com/dwethmar/go-api/pkg/services"
	"github.com/dwethmar/go-api/pkg/services/content"
)

// ErrorResponds is the default error responds.
type ErrorResponds struct {
	error string
}

// EntryIndex get list of entries
func EntryIndex(s *services.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if (*r).Method == "OPTIONS" {
			return
		}

		var entries, err = s.Content.GetAll()

		if err != nil {
			fmt.Printf("Error while getting entries: %v", err)
			request.SendServerError(w, r)
			return
		}

		request.SendJSON(w, r, entries, http.StatusOK)
	})
}

// CreateEntry creates a new entry from post data.
func CreateEntry(s *services.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		newEntry := models.AddContent{
			Fields: make(models.FieldTranslations),
		}

		err := decoder.Decode(&newEntry)

		if err != nil {
			fmt.Printf("Error while decoding entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		if err := newEntry.Validate(); err != nil {
			request.SendBadRequestError(w, r, err)
			return
		}

		entry, err := s.Content.Create(newEntry)

		if err != nil {
			fmt.Printf("Error while creating entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		request.SendJSON(w, r, entry, http.StatusCreated)
	})
}

// UpdateEntry updates an existing entry from post data.
func UpdateEntry(s *services.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())
		if err != nil {
			request.SendServerError(w, r)
		}

		decoder := json.NewDecoder(r.Body)
		var updateEntry models.UpdateContent
		err = decoder.Decode(&updateEntry)

		if err != nil {
			fmt.Printf("Error while decoding entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		if err := updateEntry.Validate(); err != nil {
			request.SendBadRequestError(w, r, err)
			return
		}

		entry, err := s.Content.GetOne(id)

		if err != nil {
			fmt.Printf("Error while getting entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		entry.Name = updateEntry.Name
		entry.Fields = updateEntry.Fields
		entry.UpdatedOn = time.Now()

		err = s.Content.Update(*entry)

		if err != nil {
			fmt.Printf("Error while updating entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		request.SendJSON(w, r, entry, http.StatusOK)
	})
}

// DeleteEntry deletes an entry by entry id.
func DeleteEntry(s *services.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())

		if err != nil {
			fmt.Printf("Error while getting id: %v", err)
			request.SendServerError(w, r)
		}

		entry, err := s.Content.GetOne(id)

		if err != nil {
			if err == content.ErrNotFound {
				fmt.Printf("Could not find entry: %v", err)
				request.SendNotFoundError(w, r)
				return
			}
			fmt.Printf("Error while Getting entry: %v %v", err, entry)
			request.SendServerError(w, r)
			return
		}

		err = s.Content.Delete(id)

		if err != nil {
			fmt.Printf("Error while deleting entry: %v", err)
			request.SendServerError(w, r)
			return
		}

		request.SendJSON(w, r, entry, http.StatusOK)
	})
}

// HandleEntrySingle gets an single entry by entry id.
func HandleEntrySingle(s *services.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())
		if err != nil {
			log.Print(fmt.Sprintf("Error on retreiving ID: %v", err))
			request.SendServerError(w, r)
		}

		entry, err := s.Content.GetOne(id)

		if err != nil {
			if err == content.ErrNotFound {
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
