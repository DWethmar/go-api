package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dwethmar/go-api/pkg/api/input"
	"github.com/dwethmar/go-api/pkg/api/output"
	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/content"
	"github.com/dwethmar/go-api/pkg/validator"
	"github.com/stretchr/testify/assert"
)

// go test ./pkg/api/handler/ -run TestContentHandler_Create

func TestContentHandler_List(t *testing.T) {
	addItems := []*content.Content{
		{
			ID:   common.NewID(),
			Name: "Test1",
			Fields: content.FieldTranslations{
				defaultLocale: {
					"attrA": 1,
				},
			},
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
		},
		{
			ID:   common.NewID(),
			Name: "Test2",
			Fields: content.FieldTranslations{
				defaultLocale: {
					"attrA": 1,
				},
			},
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
		},
	}

	var p = make([]*output.Content, len(addItems))
	for i, d := range addItems {
		p[i] = output.MapContent(d)
	}

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	entries := []*content.Content{}
	rr := httptest.NewRecorder()
	service := content.NewInMemRepository()
	h := NewContentHandler(service, validator.NewValidator())

	for _, newEntry := range addItems {
		ID, _ := service.Create(newEntry)
		entry, err := service.Get(ID)
		assert.Nil(t, err)
		entries = append(entries, entry)
	}

	handler := http.HandlerFunc(h.List)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be equal")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Content-Type code should be equal")

	// Check the response body is what we expect.
	expected, _ := json.Marshal(p)
	assert.Equal(t, string(expected), rr.Body.String(), "handler returned unexpected body")
}

func TestContentHandler_Create(t *testing.T) {
	now := time.Now()

	addEntry := input.AddContent{
		Name: "Test2",
		Fields: input.FieldTranslations{
			defaultLocale: {
				"attrA": 1,
			},
		},
	}

	body, _ := json.Marshal(addEntry)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()
	service := content.NewInMemRepository()
	h := NewContentHandler(service, validator.NewValidator())

	handler := http.HandlerFunc(h.Create)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Status code should be equal")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Content-Type code should be equal")

	// Check the response body is what we expect.
	addedEntry := content.Content{}
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &addedEntry))

	if now.After(addedEntry.CreatedOn) {
		t.Errorf("handler returned invalid createdOn: received %v, excepted CreatedOn to be smaller then %v", addedEntry.CreatedOn, now)
	}

	if now.After(addedEntry.UpdatedOn) {
		t.Errorf("handler returned invalid updatedOn: received %v, excepted UpdatedOn to be smaller then  %v", addedEntry.UpdatedOn, now)
	}

	if equal, err := areFieldsEqual(addEntry.Fields, addedEntry.Fields); !equal || err != nil {
		if err != nil {
			t.Error(err)
		} else {
			t.Errorf("Fields are not equal. left: %v right %v", addEntry.Fields, addedEntry.Fields)
		}
	}
}

func TestContentHandler_InvalidCreate(t *testing.T) {
	body, _ := json.Marshal(input.AddContent{
		Name: "abcdefghijklmnokqrstuvwxyz-abcdefghijklmnokqrstuvwxyz",
		Fields: input.FieldTranslations{
			defaultLocale: {
				"attrA": 1,
			},
		},
	})
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()
	service := content.NewInMemRepository()
	h := NewContentHandler(service, validator.NewValidator())

	handler := http.HandlerFunc(h.Create)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be equal")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Content-Type code should be equal")

	fmt.Printf(rr.Body.String())
}

func TestContentHandler_Update(t *testing.T) {
	service := content.NewInMemRepository()
	h := NewContentHandler(service, validator.NewValidator())

	ID, _ := service.Create(&content.Content{
		ID:   common.NewID(),
		Name: "name",
		Fields: content.FieldTranslations{
			"nl": content.Fields{
				"attrA": "Value",
				"attrB": 1,
				"attrC": []string{"one", "two", "three"},
			},
		},
	})

	addedEntry, err := service.Get(ID)
	assert.Nil(t, err)

	updateEntry := input.UpdateContent{
		Name: "updated name",
		Fields: input.FieldTranslations{
			"nl": {
				"attrA": "Value2",
				"attrB": 2,
				"attrC": []string{"four", "five", "six"},
			},
		},
	}

	body, _ := json.Marshal(updateEntry)

	req := httptest.NewRequest("POST", fmt.Sprintf("/%v", addedEntry.ID), bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = common.WithID(ctx, addedEntry.ID)
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(h.Update)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be equal")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Content-Type code should be equal")

	// Check the response body is what we expect.
	updatedEntry := content.Content{}
	err = json.Unmarshal(rr.Body.Bytes(), &updatedEntry)
	assert.Nil(t, err)

	if equal, err := areFieldsEqual(updateEntry.Fields, updatedEntry.Fields); !equal || err != nil {
		if err != nil {
			t.Error(err)
		} else {
			t.Errorf("Fields are not equal. left: %v right %v", updateEntry.Fields, updatedEntry.Fields)
		}
	}

}

func TestContentHandler_Delete(t *testing.T) {
	service := content.NewInMemRepository()
	h := NewContentHandler(service, validator.NewValidator())

	ID, _ := service.Create(&content.Content{
		ID:   common.NewID(),
		Name: "name",
		Fields: content.FieldTranslations{
			"nl": content.Fields{
				"attrA": "Value",
				"attrB": 1,
				"attrC": []string{"one", "two", "three"},
			},
		},
	})

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/%s", ID), nil)

	ctx := req.Context()
	ctx = common.WithID(ctx, ID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Delete)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be equal")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Content-Type code should be equal")

	deletedEntry := content.Content{}
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &deletedEntry))
}
