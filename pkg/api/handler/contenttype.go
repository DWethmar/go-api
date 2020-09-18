package handler

import (
	"fmt"
	"net/http"

	"github.com/dwethmar/go-api/pkg/api/output"
	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/contenttype"
)

type contentTypeHandler struct {
	types contenttype.Service
}

// ContentTypeHandler handle content requests.
type ContentTypeHandler interface {
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

func (h *contentTypeHandler) List(w http.ResponseWriter, r *http.Request) {
	var types, err = h.types.List()

	if err != nil {
		fmt.Printf("Error while getting entries: %v", err)
		common.SendServerError(w, r)
		return
	}

	var o = []*output.ContentType{}
	for _, t := range types {
		o = append(o, output.ContentTypeOut(t))
	}

	common.SendJSON(w, r, o, http.StatusOK)
}

func (h *contentTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	common.SendJSON(w, r, nil, http.StatusOK)
}

func (h *contentTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	common.SendJSON(w, r, nil, http.StatusOK)
}

func (h *contentTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	common.SendJSON(w, r, nil, http.StatusOK)
}

func (h *contentTypeHandler) Get(w http.ResponseWriter, r *http.Request) {
	common.SendJSON(w, r, nil, http.StatusOK)
}

// NewContentTypeHandler creates new handler
func NewContentTypeHandler(service contenttype.Service) ContentTypeHandler {
	return &contentTypeHandler{
		types: service,
	}
}
