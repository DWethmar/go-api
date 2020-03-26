package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DWethmar/go-api/pkg/contentitem"
	"github.com/gorilla/mux"
)

func createTestServer() (contentitem.ContentItem, Server) {
	server := Server{
		contentItem: contentitem.CreateService(contentitem.CreateMockRepository()),
		router:      mux.NewRouter().StrictSlash(true),
	}
	c, _ := server.contentItem.Create(contentitem.AddContentItem{
		Name: "Test",
	})
	server.routes()
	return *c, server
}

var addContentItem = contentitem.AddContentItem{
	Name: "name",
	Attrs: contentitem.Attrs{
		"nl": {
			"attrA": "Value",
			"attrB": 1,
			"attrC": []string{"one", "two", "three"},
			"attrD": []int{1, 2, 3, 4},
			"attrE": float64(100) / float64(3),
			"attrF": math.MaxFloat64,
		},
	},
}

func TestHandleIndex(t *testing.T) {
	contentItem, server := createTestServer()
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
		t.Errorf("handler returned unexpected body: got %v want %v", strings.TrimSuffix(rr.Body.String(), "\n"), expected)
	}
}

func TestHandleCreate(t *testing.T) {
	now := time.Now()
	_, server := createTestServer()
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
	newContentitem := contentitem.ContentItem{}
	err := json.Unmarshal(rr.Body.Bytes(), &newContentitem)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	if now.After(newContentitem.CreatedOn) {
		t.Errorf("handler returned invalid createdOn: got %v excepted %v", newContentitem.CreatedOn, now)
	}

	if now.After(newContentitem.UpdatedOn) {
		t.Errorf("handler returned invalid updatedOn: got %v excepted %v", newContentitem.UpdatedOn, now)
	}

	eAttr, err := json.Marshal(addContentItem.Attrs)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}
	gAttr, err := json.Marshal(newContentitem.Attrs)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}
	if string(eAttr) != string(gAttr) {
		t.Errorf("handler returned unexpected Attrs: got %v want %v", addContentItem.Attrs, newContentitem.Attrs)
	}
}
