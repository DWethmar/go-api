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
	"github.com/DWethmar/go-api/pkg/contentitem"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func createTestServer(dbName string) (contentitem.ContentItem, Server) {
	var repo contentitem.Repository
	con := config.LoadEnv()

	if dbName != "" && con.DBHost != "" {
		con.DBName = dbName
		driver, cs := config.GetPostgresConnectionInfo(con)
		db, err := NewDB(driver, cs)
		if err != nil {
			panic(err)
		}
		repo = contentitem.CreatePostgresRepository(db)
	} else {
		repo = contentitem.CreateMockRepository()
	}

	server := Server{
		contentItem: contentitem.CreateService(repo),
		router:      mux.NewRouter().StrictSlash(true),
	}
	contentItem, err := server.contentItem.Create(contentitem.AddContentItem{
		Name: "Test",
	})
	if err != nil {
		panic("Could not create contentitem.")
	}
	server.routes()
	return *contentItem, server
}

func TestIntergrationHandleIndex(t *testing.T) {
	contentItem, server := createTestServer("test_one")
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v on %v", status, http.StatusOK, req.RequestURI)
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

func TestIntergrationHandleCreate(t *testing.T) {
	now := time.Now()
	_, server := createTestServer("")

	addContentItem := contentitem.AddContentItem{
		Name: "name",
		Attrs: contentitem.Attrs{
			"nl": {
				"attrA": "Value",
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
	}

	// Check the response content-type is what we expect.
	if rType := rr.Header().Get("Content-Type"); rType != "application/json" {
		t.Errorf("content type header does not match: got %v want %v", rType, "application/json")
	}

	// Check the response body is what we expect.
	addedContentitem := contentitem.ContentItem{}
	err := json.Unmarshal(rr.Body.Bytes(), &addedContentitem)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	if now.After(addedContentitem.CreatedOn) {
		t.Errorf("handler returned invalid createdOn: got %v excepted CreatedOn to be smaller then %v", addedContentitem.CreatedOn, now)
	}

	if now.After(addedContentitem.UpdatedOn) {
		t.Errorf("handler returned invalid updatedOn: got %v excepted UpdatedOn to be smaller then  %v", addedContentitem.UpdatedOn, now)
	}

	eAttr, err := json.Marshal(addContentItem.Attrs)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}
	gAttr, err := json.Marshal(addedContentitem.Attrs)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}
	if string(eAttr) != string(gAttr) {
		t.Errorf("handler returned unexpected Attrs: got %v want %v", addContentItem.Attrs, addedContentitem.Attrs)
	}
}
