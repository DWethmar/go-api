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
	"github.com/dwethmar/go-api/pkg/validator"
)

// ErrorResponds is the default error responds.
type ErrorResponds struct {
	error string
}

type contentHandler struct {
	content  content.Service
	validate *validator.Validation
}

// ContentHandler handle content requests.
type ContentHandler interface {
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

// List get list of entries
func (h *contentHandler) List(w http.ResponseWriter, r *http.Request) {
	entries, err := h.content.List()

	if err != nil {
		fmt.Printf("Error while getting entries: %v", err)
		common.SendServerError(w, r)
		return
	}

	var p = make([]*output.Content, len(entries))
	for i, d := range entries {
		p[i] = output.MapContent(d)
	}

	common.SendJSON(w, r, p, http.StatusOK)
}

// Create creates a new entry from post data.
func (h *contentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var add = &input.AddContent{}

	err := json.NewDecoder(r.Body).Decode(&add)
	if err != nil {
		fmt.Printf("Error while decoding entry: %v", err)
		common.SendServerError(w, r)
		return
	}

	err = h.validate.Struct(add)
	if err != nil {
		common.SendBadRequestError(w, r, validator.NewErrValidation(err))
		return
	}

	c := input.MapAddContent(add)

	c.ID, err = h.content.Create(c)

	if err != nil {
		fmt.Printf("Error while creating entry: %v", err)
		common.SendServerError(w, r)
		return
	}

	common.SendJSON(w, r, output.MapContent(c), http.StatusCreated)
}

// UpdateContent updates an existing entry from post data.
func (h *contentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := common.UUIDFromContext(r.Context())
	if err != nil {
		common.SendServerError(w, r)
	}

	var update = &input.UpdateContent{}

	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		fmt.Printf("Error while decoding entry: %v", err)
		common.SendServerError(w, r)
		return
	}

	c, err := h.content.Get(id)
	if err != nil {
		fmt.Printf("Error while getting entry: %v", err)
		common.SendServerError(w, r)
		return
	}

	u := input.MapUpdateContent(update)

	c.Name = u.Name
	c.Fields = u.Fields
	c.UpdatedOn = time.Now()

	err = h.content.Update(c)

	if err != nil {
		fmt.Printf("Error while updating entry: %v", err)
		common.SendServerError(w, r)
		return
	}

	common.SendJSON(w, r, output.MapContent(c), http.StatusOK)
}

// Delete deletes an entry by entry id.
func (h *contentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := common.UUIDFromContext(r.Context())

	if err != nil {
		fmt.Printf("Error while getting id: %v", err)
		common.SendServerError(w, r)
	}

	c, err := h.content.Get(id)

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

	err = h.content.Delete(id)

	if err != nil {
		fmt.Printf("Error while deleting entry: %v", err)
		common.SendServerError(w, r)
		return
	}

	common.SendJSON(w, r, output.MapContent(c), http.StatusOK)
}

// GetSingleContent gets an single entry by entry id.
func (h *contentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := common.UUIDFromContext(r.Context())
	if err != nil {
		log.Print(fmt.Sprintf("Error on retreiving ID: %v", err))
		common.SendServerError(w, r)
	}

	entry, err := h.content.Get(id)

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
}

// NewContentHandler creates new handler
func NewContentHandler(service content.Service, validate *validator.Validation) ContentHandler {
	return &contentHandler{
		content:  service,
		validate: validate,
	}
}
