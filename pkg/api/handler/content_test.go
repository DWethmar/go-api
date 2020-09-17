package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dwethmar/go-api/pkg/api/output"
	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/content"
	"github.com/dwethmar/go-api/pkg/database"
	"github.com/stretchr/testify/assert"
)

var defaultLocale = "nl"

func TestListContent(t *testing.T) {
	addItems := []*content.Content{
		{
			ID:   common.NewID(),
			Name: "Test1",
			Fields: content.FieldTranslations{
				defaultLocale: content.Fields{
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
				defaultLocale: content.Fields{
					"attrA": 1,
				},
			},
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
		},
	}

	var p []*output.Content
	for _, d := range addItems {
		p = append(p, &output.Content{
			ID:        d.ID,
			Name:      d.Name,
			Fields:    d.Fields,
			CreatedOn: d.CreatedOn,
			UpdatedOn: d.UpdatedOn,
		})
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	entries := []*content.Content{}
	rr := httptest.NewRecorder()

	var service content.Service
	db, cleanup, err := database.NewTestDB()
	if err == nil {
		service = content.NewPostgresRepository(db)
		defer cleanup()
	} else {
		log.Print(fmt.Sprintf("Error : %v", err))
		service = content.NewInMemRepository()
	}

	h := NewContentHandler(service)

	for _, newEntry := range addItems {
		ID, _ := service.Create(newEntry)
		entry, err := service.Get(ID)
		if err != nil {
			t.Errorf("something went wrong %v", err)
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

	if rr.Body.String() != string(expected) {
		t.Errorf("handler returned unexpected body: received %v expected %v", rr.Body.String(), string(expected))
	}
}

// func TestCreateContent(t *testing.T) {
// 	now := time.Now()

// 	addEntry := input.AddContent{
// 		Name: "Test2",
// 		Fields: content.FieldTranslations{
// 			defaultLocale: content.Fields{
// 				"attrA": 1,
// 			},
// 		},
// 	}

// 	body, _ := json.Marshal(addEntry)
// 	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))

// 	rr := httptest.NewRecorder()

// 	var service content.Service
// 	db, cleanup, err := database.NewTestDB()

// 	if err != nil {
// 		service = content.NewPostgresRepository(db)
// 	} else {
// 		service = content.NewInMemRepository()
// 		cleanup = func() {
// 			fmt.Println("Nothing to cleanup")
// 		}
// 	}
// 	defer cleanup()
// 	h := NewContentHandler(service)

// 	handler := http.HandlerFunc(h.Create)
// 	handler.ServeHTTP(rr, req)

// 	status := rr.Code
// 	assert.Equal(t, status, http.StatusCreated, "Status code should be equal")

// 	rType := rr.Header().Get("Content-Type")
// 	assert.Equal(t, rType, "application/json", "Content-Type code should be equal")

// 	// Check the response body is what we expect.
// 	addedEntry := content.Content{}
// 	err = json.Unmarshal(rr.Body.Bytes(), &addedEntry)
// 	if err != nil {
// 		t.Errorf("Error while parsing body %v", err)
// 	}

// 	if now.After(addedEntry.CreatedOn) {
// 		t.Errorf("handler returned invalid createdOn: received %v, excepted CreatedOn to be smaller then %v", addedEntry.CreatedOn, now)
// 	}

// 	if now.After(addedEntry.UpdatedOn) {
// 		t.Errorf("handler returned invalid updatedOn: received %v, excepted UpdatedOn to be smaller then  %v", addedEntry.UpdatedOn, now)
// 	}

// 	if equal, err := areFieldsEqual(addEntry.Fields, addedEntry.Fields); !equal || err != nil {
// 		if err != nil {
// 			t.Error(err)
// 		} else {
// 			t.Errorf("Fields are not equal. left: %v right %v", addEntry.Fields, addedEntry.Fields)
// 		}
// 	}
// }

// func TestUpdateContent(t *testing.T) {
// 	var service content.Service
// 	db, cleanup, err := database.NewTestDB()

// 	if err != nil {
// 		service = content.NewPostgresRepository(db)
// 	} else {
// 		service = content.NewInMemRepository()
// 		cleanup = func() {
// 			fmt.Println("Nothing to cleanup")
// 		}
// 	}
// 	defer cleanup()
// 	h := NewContentHandler(service)

// 	ID, _ := service.Create(&content.Content{
// 		ID:   common.NewID(),
// 		Name: "name",
// 		Fields: content.FieldTranslations{
// 			"nl": content.Fields{
// 				"attrA": "Value",
// 				"attrB": 1,
// 				"attrC": []string{"one", "two", "three"},
// 			},
// 		},
// 	})

// 	addedEntry, err := service.Get(ID)
// 	assert.Nil(t, err)

// 	updateEntry := input.UpdateContent{
// 		Name: "updated name",
// 		Fields: content.FieldTranslations{
// 			"nl": content.Fields{
// 				"attrA": "Value2",
// 				"attrB": 2,
// 				"attrC": []string{"four", "five", "six"},
// 			},
// 		},
// 	}

// 	body, _ := json.Marshal(updateEntry)

// 	req := httptest.NewRequest("POST", fmt.Sprintf("/%v", addedEntry.ID), bytes.NewBuffer(body))
// 	rr := httptest.NewRecorder()

// 	ctx := req.Context()
// 	ctx = common.WithID(ctx, addedEntry.ID)
// 	req = req.WithContext(ctx)

// 	handler := http.HandlerFunc(h.Update)
// 	handler.ServeHTTP(rr, req)

// 	status := rr.Code
// 	assert.Equal(t, status, http.StatusOK, "Status code should be equal")

// 	rType := rr.Header().Get("Content-Type")
// 	assert.Equal(t, rType, "application/json", "Content-Type code should be equal")

// 	// Check the response body is what we expect.
// 	updatedEntry := content.Content{}
// 	err = json.Unmarshal(rr.Body.Bytes(), &updatedEntry)
// 	if err != nil {
// 		t.Errorf("Error while parsing body for added content-item %v", err)
// 		return
// 	}

// 	if equal, err := areFieldsEqual(updateEntry.Fields, updatedEntry.Fields); !equal || err != nil {
// 		if err != nil {
// 			t.Error(err)
// 		} else {
// 			t.Errorf("Fields are not equal. left: %v right %v", updateEntry.Fields, updatedEntry.Fields)
// 		}
// 	}

// }

// func TestDeleteContent(t *testing.T) {
// 	var service content.Service
// 	db, cleanup, err := database.NewTestDB()

// 	if err != nil {
// 		service = content.NewPostgresRepository(db)
// 	} else {
// 		service = content.NewInMemRepository()
// 		cleanup = func() {
// 			fmt.Println("Nothing to cleanup")
// 		}
// 	}
// 	defer cleanup()
// 	h := NewContentHandler(service)

// 	ID, _ := service.Create(&content.Content{
// 		ID:   common.NewID(),
// 		Name: "name",
// 		Fields: content.FieldTranslations{
// 			"nl": content.Fields{
// 				"attrA": "Value",
// 				"attrB": 1,
// 				"attrC": []string{"one", "two", "three"},
// 			},
// 		},
// 	})

// 	req := httptest.NewRequest("DELETE", fmt.Sprintf("/%s", ID), nil)

// 	ctx := req.Context()
// 	ctx = common.WithID(ctx, ID)
// 	req = req.WithContext(ctx)

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(h.Delete)
// 	handler.ServeHTTP(rr, req)

// 	status := rr.Code
// 	assert.Equal(t, status, http.StatusOK, "Status code should be equal")

// 	rType := rr.Header().Get("Content-Type")
// 	assert.Equal(t, rType, "application/json", "Content-Type code should be equal")

// 	deletedEntry := content.Content{}
// 	err = json.Unmarshal(rr.Body.Bytes(), &deletedEntry)
// 	if err != nil {
// 		t.Errorf("Error while parsing body for deleted entry %v", err)
// 		return
// 	}
// }
