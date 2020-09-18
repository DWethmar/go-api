package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwethmar/go-api/pkg/api/output"
	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/contenttype"
	"github.com/stretchr/testify/assert"
)

func TestContentTypeHandler_List(t *testing.T) {
	addItems := []*contenttype.ContentType{
		{
			ID:   common.NewID(),
			Name: "Test1",
			Fields: []*contenttype.Field{
				{
					ID:            common.NewID(),
					ContentTypeID: common.NewID(),
					Key:           "name",
					Name:          "Name",
					FieldType:     "string",
					Length:        10,
					CreatedOn:     common.Now(),
					UpdatedOn:     common.Now(),
				},
			},
			CreatedOn: common.Now(),
			UpdatedOn: common.Now(),
		},
		{
			ID:   common.NewID(),
			Name: "Test2",
			Fields: []*contenttype.Field{
				{
					ID:            common.NewID(),
					ContentTypeID: common.NewID(),
					Key:           "name",
					Name:          "Name",
					FieldType:     "string",
					Length:        10,
					CreatedOn:     common.Now(),
					UpdatedOn:     common.Now(),
				},
			},
			CreatedOn: common.Now(),
			UpdatedOn: common.Now(),
		},
	}

	var p []*output.ContentType
	for _, d := range addItems {
		p = append(p, output.MapContentType(d))
	}

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	entries := []*contenttype.ContentType{}

	rr := httptest.NewRecorder()
	service := contenttype.NewInMemRepository()
	h := NewContentTypeHandler(service)

	for _, newEntry := range addItems {
		ID, _ := service.Create(newEntry)
		entry, err := service.Get(ID)
		if err != nil {
			t.Error(err)
		}
		entries = append(entries, entry)
	}

	handler := http.HandlerFunc(h.List)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	assert.Equal(t, status, http.StatusOK, "Status code should be equal")

	rType := rr.Header().Get("Content-Type")
	assert.Equal(t, rType, "application/json", "Content-Type code should be equal")

	// Check the response body is what we expect.
	expected, _ := json.Marshal(p)

	assert.Equal(t, rr.Body.String(), string(expected), "Didn't expect value.")
}
