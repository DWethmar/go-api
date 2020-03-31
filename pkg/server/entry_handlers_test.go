package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DWethmar/go-api/pkg/config"
	"github.com/DWethmar/go-api/pkg/contententry"
	"github.com/DWethmar/go-api/pkg/database"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func createTestServer(dbName string) (contententry.Entry, Server) {
	var repo contententry.Repository
	con := config.LoadEnv()

	if dbName != "" && con.DBHost != "" {
		con.DBName = dbName
		driver, cs := config.GetPostgresConnectionInfo(con)
		db, err := database.CreateDB(driver, cs)

		if err != nil {
			panic(err)
		}

		repo = contententry.CreatePostgresRepository(db)
	} else {
		repo = contententry.CreateMockRepository()
	}

	server := Server{
		entries: contententry.CreateService(repo),
		router:  mux.NewRouter().StrictSlash(true),
	}

	contentItem, err := server.entries.Create(contententry.AddEntry{
		Name: "Test",
		Fields: contententry.FieldTranslations{
			"nl": {
				"attr1": "test",
			},
		},
	})

	if err != nil {
		panic("Could not create contententry.")
	}

	server.routes()

	return *contentItem, server
}

func TestIntergrationInvalidIds(t *testing.T) {
	now := time.Now()
	_, server := createTestServer("test_two")

	addEntry := contententry.AddEntry{
		Name: "test",
		Fields: contententry.FieldTranslations{
			"nl": contententry.Fields{
				"attrA": "Value A",
			},
		},
	}

	body, _ := json.Marshal(addEntry)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusCreated, req.RequestURI)
		return
	}

	// Check the response content-type is what we expect.
	if rType := rr.Header().Get("Content-Type"); rType != "application/json" {
		t.Errorf("content type header does not match: got %v want %v", rType, "application/json")
		return
	}

	// Check the response body is what we expect.
	addedEntry := contententry.Entry{}
	if err := json.Unmarshal(rr.Body.Bytes(), &addedEntry); err != nil {
		t.Errorf("Error while parsing body for added content-item %v", err)
		return
	}

	addedEntry.Fields = contententry.FieldTranslations{
		"nl": contententry.Fields{
			"attrA": "Value B",
			"attrB": 1,
			"attrC": []string{"three", "four", "five"},
		},
		"en": contententry.Fields{
			"attrD": "Value C",
			"attrE": 1,
			"attrF": []string{"one", "two", "three"},
		},
	}

	now = time.Now()

	body, _ = json.Marshal(addedEntry)
	req = httptest.NewRequest("POST", fmt.Sprintf("/%v", addedEntry.ID), bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	updatedEntry := contententry.Entry{}
	if err := json.Unmarshal(rr.Body.Bytes(), &updatedEntry); err != nil {
		t.Errorf("Error while parsing body for updated content-item: %v", err)
		return
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusCreated, req.RequestURI)
		return
	}

	if !addedEntry.CreatedOn.Equal(updatedEntry.CreatedOn) {
		t.Errorf("handler returned invalid createdOn: got %v excepted CreatedOn to equal %v", updatedEntry.CreatedOn, addedEntry.CreatedOn)
	}

	if now.Before(addedEntry.UpdatedOn) {
		t.Errorf("handler returned invalid createdOn: got %v excepted UpdatedOn to be larger then %v", updatedEntry.CreatedOn, addedEntry.CreatedOn)
	}

	eAttr, err := json.Marshal(addedEntry.Fields)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	gAttr, err := json.Marshal(updatedEntry.Fields)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	if string(eAttr) != string(gAttr) {
		t.Errorf("handler returned unexpected fields: got %v want %v", addedEntry.Fields, addedEntry.Fields)
	}
}
func TestIntergrationEntryIndex(t *testing.T) {
	contentItem, server := createTestServer("test_one")
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusOK, req.RequestURI)
	}

	// Check the response content-type is what we expect.
	if rType := rr.Header().Get("Content-Type"); rType != "application/json" {
		t.Errorf("content type header does not match: got %v want %v", rType, "application/json")
	}

	// Check the response body is what we expect.
	c, err := json.Marshal(contentItem)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	if expected := fmt.Sprintf("[%s]\n", string(c)); rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \n excpected: %v \n received:%v", expected, strings.TrimSuffix(rr.Body.String(), "\n"))
	}
}

func TestIntergrationEntryCreate(t *testing.T) {
	now := time.Now()
	_, server := createTestServer("test_two")

	addEntry := contententry.AddEntry{
		Name: "name",
		Fields: contententry.FieldTranslations{
			"nl": contententry.Fields{
				"attrA": "Value",
				"attrB": 1,
				"attrC": []string{"one", "two", "three"},
			},
		},
	}

	body, _ := json.Marshal(addEntry)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusCreated, req.RequestURI)
	}

	// Check the response content-type is what we expect.
	if rType := rr.Header().Get("Content-Type"); rType != "application/json" {
		t.Errorf("content type header does not match: got %v want %v", rType, "application/json")
	}

	// Check the response body is what we expect.
	addedEntry := contententry.Entry{}
	err := json.Unmarshal(rr.Body.Bytes(), &addedEntry)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	if now.After(addedEntry.CreatedOn) {
		t.Errorf("handler returned invalid createdOn: got %v excepted CreatedOn to be smaller then %v", addedEntry.CreatedOn, now)
	}

	if now.After(addedEntry.UpdatedOn) {
		t.Errorf("handler returned invalid updatedOn: got %v excepted UpdatedOn to be smaller then  %v", addedEntry.UpdatedOn, now)
	}

	eAttr, err := json.Marshal(addEntry.Fields)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	gAttr, err := json.Marshal(addedEntry.Fields)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	if string(eAttr) != string(gAttr) {
		t.Errorf("handler returned unexpected fields: excpected %v received %v", addEntry.Fields, addedEntry.Fields)
	}
}

func TestIntergrationEntryUpdate(t *testing.T) {
	now := time.Now()
	_, server := createTestServer("test_two")

	addEntry := contententry.AddEntry{
		Name: "test",
		Fields: contententry.FieldTranslations{
			"nl": contententry.Fields{
				"attrA": "Value A",
				"attrB": 1,
				"attrC": []string{"one", "two", "three"},
			},
		},
	}

	body, _ := json.Marshal(addEntry)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusCreated, req.RequestURI)
		return
	}

	// Check the response content-type is what we expect.
	if rType := rr.Header().Get("Content-Type"); rType != "application/json" {
		t.Errorf("content type header does not match: got %v want %v", rType, "application/json")
		return
	}

	// Check the response body is what we expect.
	addedEntry := contententry.Entry{}
	err := json.Unmarshal(rr.Body.Bytes(), &addedEntry)
	if err != nil {
		t.Errorf("Error while parsing body for added content-item %v", err)
		return
	}

	addedEntry.Fields = contententry.FieldTranslations{
		"nl": contententry.Fields{
			"attrA": "Value B",
			"attrB": 1,
			"attrC": []string{"three", "four", "five"},
		},
		"en": contententry.Fields{
			"attrD": "Value C",
			"attrE": 1,
			"attrF": []string{"one", "two", "three"},
		},
	}

	now = time.Now()

	body, _ = json.Marshal(addedEntry)
	req = httptest.NewRequest("POST", fmt.Sprintf("/%v", addedEntry.ID), bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	updatedEntry := contententry.Entry{}
	err = json.Unmarshal(rr.Body.Bytes(), &updatedEntry)
	if err != nil {
		t.Errorf("Error while parsing body for updated content-item: %v", err)
		return
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusCreated, req.RequestURI)
		return
	}

	if !addedEntry.CreatedOn.Equal(updatedEntry.CreatedOn) {
		t.Errorf("handler returned invalid createdOn: got %v excepted CreatedOn to equal %v", updatedEntry.CreatedOn, addedEntry.CreatedOn)
	}

	if now.Before(addedEntry.UpdatedOn) {
		t.Errorf("handler returned invalid createdOn: got %v excepted UpdatedOn to be larger then %v", updatedEntry.CreatedOn, addedEntry.CreatedOn)
	}

	eAttr, err := json.Marshal(addedEntry.Fields)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	gAttr, err := json.Marshal(updatedEntry.Fields)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	if string(eAttr) != string(gAttr) {
		t.Errorf("handler returned unexpected fields: got %v want %v", addedEntry.Fields, addedEntry.Fields)
	}
}

func TestIntergrationEntryDelete(t *testing.T) {
	_, server := createTestServer("test_two")

	addContentItem := contententry.AddEntry{
		Name: "test",
		Fields: contententry.FieldTranslations{
			"nl": contententry.Fields{
				"attrA": "Value A",
				"attrB": 1,
				"attrC": []string{"one", "two", "three"},
			},
		},
	}

	body, _ := json.Marshal(addContentItem)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v on %v", status, http.StatusCreated, req.RequestURI)
		return
	}

	// Check the response content-type is what we expect.
	if rType := rr.Header().Get("Content-Type"); rType != "application/json" {
		t.Errorf("content type header does not match: got %v want %v", rType, "application/json")
		return
	}

	addEntry := contententry.Entry{}
	err := json.Unmarshal(rr.Body.Bytes(), &addEntry)
	if err != nil {
		t.Errorf("Error while parsing body for added content-item %v", err)
		return
	}

	body, _ = json.Marshal(addContentItem)
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/%v", addEntry.ID), bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusOK, req.RequestURI)
		return
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%v", addEntry.ID), nil)
	rr = httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusNotFound, req.RequestURI)
		return
	}
}

func TestIntergrationEntrySingle(t *testing.T) {
	_, server := createTestServer("test_two")

	addEntry := contententry.AddEntry{
		Name: "test",
		Fields: contententry.FieldTranslations{
			"nl": contententry.Fields{
				"attrA": "Value A",
				"attrB": 1,
				"attrC": []string{"one", "two", "three"},
			},
		},
	}

	body, _ := json.Marshal(addEntry)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v on %v", status, http.StatusCreated, req.RequestURI)
		return
	}

	// Check the response content-type is what we expect.
	if rType := rr.Header().Get("Content-Type"); rType != "application/json" {
		t.Errorf("content type header does not match: got %v want %v", rType, "application/json")
		return
	}

	addedEntry := contententry.Entry{}
	err := json.Unmarshal(rr.Body.Bytes(), &addedEntry)
	if err != nil {
		t.Errorf("Error while parsing body for added content-item %v", err)
		return
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%v", addedEntry.ID), nil)
	rr = httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusNotFound, req.RequestURI)
		return
	}

	eAttr, err := json.Marshal(addEntry.Fields)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	gAttr, err := json.Marshal(addedEntry.Fields)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	if string(eAttr) != string(gAttr) {
		t.Errorf("handler returned unexpected fields: got %v want %v", addEntry.Fields, addedEntry.Fields)
	}
}
