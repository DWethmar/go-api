package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dwethmar/go-api/pkg/api/input"
	"github.com/dwethmar/go-api/pkg/api/output"
	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/content"
	"github.com/dwethmar/go-api/pkg/store"
)

// ErrorResponds is the default error responds.
type ErrorResponds struct {
	error string
}

/**
GET /tickets - Retrieves a list of tickets
GET /tickets/12 - Retrieves a specific ticket
POST /tickets - Creates a new ticket
PUT /tickets/12 - Updates ticket #12
PATCH /tickets/12 - Partially updates ticket #12
DELETE /tickets/12 - Deletes ticket #12
**/

// ListContent get list of entries
func ListContent(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		entries, err := s.Content.List()

		if err != nil {
			fmt.Printf("Error while getting entries: %v", err)
			common.SendServerError(w, r)
			return
		}

		var p = make([]*output.Content, 0)
		for _, d := range entries {
			p = append(p, &output.Content{
				ID:        d.ID,
				Name:      d.Name,
				Fields:    d.Fields,
				CreatedOn: d.CreatedOn,
				UpdatedOn: d.UpdatedOn,
			})
		}

		common.SendJSON(w, r, p, http.StatusOK)
	})
}

// CreateContent creates a new entry from post data.
func CreateContent(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input = &input.AddContent{}

		err := json.NewDecoder(r.Body).Decode(&input)

		if err != nil {
			fmt.Printf("Error while decoding entry: %v", err)
			common.SendServerError(w, r)
			return
		}

		c := &content.Content{
			ID:        common.NewID(),
			Name:      input.Name,
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
			Fields:    input.Fields,
		}
		c.ID, err = s.Content.Create(c)

		if err != nil {
			fmt.Printf("Error while creating entry: %v", err)
			common.SendServerError(w, r)
			return
		}

		p := &output.Content{
			ID:        c.ID,
			Name:      c.Name,
			Fields:    c.Fields,
			CreatedOn: c.CreatedOn,
			UpdatedOn: c.UpdatedOn,
		}

		common.SendJSON(w, r, p, http.StatusCreated)
	})
}

// UpdateContent updates an existing entry from post data.
func UpdateContent(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())
		if err != nil {
			common.SendServerError(w, r)
		}

		var input = &input.UpdateContent{}

		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			fmt.Printf("Error while decoding entry: %v", err)
			common.SendServerError(w, r)
			return
		}

		c, err := s.Content.Get(id)
		if err != nil {
			fmt.Printf("Error while getting entry: %v", err)
			common.SendServerError(w, r)
			return
		}

		c.Name = input.Name
		c.Fields = input.Fields
		c.UpdatedOn = time.Now()

		err = s.Content.Update(c)

		if err != nil {
			fmt.Printf("Error while updating entry: %v", err)
			common.SendServerError(w, r)
			return
		}

		p := &output.Content{
			ID:        c.ID,
			Name:      c.Name,
			Fields:    c.Fields,
			CreatedOn: c.CreatedOn,
			UpdatedOn: c.UpdatedOn,
		}

		common.SendJSON(w, r, p, http.StatusOK)
	})
}

// DeleteContent deletes an entry by entry id.
func DeleteContent(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())

		if err != nil {
			fmt.Printf("Error while getting id: %v", err)
			common.SendServerError(w, r)
		}

		c, err := s.Content.Get(id)

		if err != nil {
			if err == content.ErrNotFound {
				fmt.Printf("Could not find entry: %v", err)
				common.SendNotFoundError(w, r)
				return
			}
			fmt.Printf("Error while Getting entry: %v %v", err, c)
			common.SendServerError(w, r)
			return
		}

		err = s.Content.Delete(id)

		if err != nil {
			fmt.Printf("Error while deleting entry: %v", err)
			common.SendServerError(w, r)
			return
		}

		p := &output.Content{
			ID:        c.ID,
			Name:      c.Name,
			Fields:    c.Fields,
			CreatedOn: c.CreatedOn,
			UpdatedOn: c.UpdatedOn,
		}

		common.SendJSON(w, r, p, http.StatusOK)
	})
}

// GetSingleContent gets an single entry by entry id.
func GetSingleContent(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.UUIDFromContext(r.Context())
		if err != nil {
			log.Print(fmt.Sprintf("Error on retreiving ID: %v", err))
			common.SendServerError(w, r)
		}

		entry, err := s.Content.Get(id)

		if err != nil {
			if err == content.ErrNotFound {
				log.Print(fmt.Sprintf("Entry not found: %v", err))
				common.SendNotFoundError(w, r)
				return
			}
			log.Print(fmt.Sprintf("Somthing went wrong: %v", err))
			common.SendServerError(w, r)
			return
		}

		common.SendJSON(w, r, entry, http.StatusOK)
	})
}
