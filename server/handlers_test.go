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

	"github.com/DWethmar/go-api/contentitem"
	"github.com/gorilla/mux"
)

func TestHandleIndex(t *testing.T) {
	server := Server{
		contentItem: contentitem.NewService(contentitem.CreateMockRepository(contentitem.ContentItem{
			ID:   1,
			Name: "Name",
			Attrs: contentitem.Attrs{
				"Test": "Value",
			},
		})),
		router: mux.NewRouter().StrictSlash(true),
	}
	server.routes()

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
	contentItem, _ := server.contentItem.GetOne(1)
	c, err := json.Marshal(contentItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	if expected := fmt.Sprintf("[%s]\n", string(c)); rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", strings.TrimSuffix(rr.Body.String(), "\n"), expected)
	}
}

func TestHandleCreate(t *testing.T) {
	now := time.Now()
	contentItem := contentitem.ContentItem{
		ID:   1,
		Name: "Name",
		Attrs: contentitem.Attrs{
			"test1": "Value",
			"test2": 1,
			"test3": []string{"one", "two", "three"},
			"test4": []int{1, 2, 3, 4},
			"test5": float64(100) / float64(3),
			"test6": math.MaxFloat64,
		},
	}
	server := Server{
		contentItem: contentitem.NewService(contentitem.CreateMockRepository()),
		router:      mux.NewRouter().StrictSlash(true),
	}
	server.routes()
	body, _ := json.Marshal(contentitem.NewContentItem{
		Name:  contentItem.Name,
		Attrs: contentItem.Attrs,
	})
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
		fmt.Println(err)
		return
	}

	if newContentitem.ID != 1 {
		t.Errorf("handler returned unexpected ID: got %v want %v", newContentitem.ID, 1)
	}

	if now.After(newContentitem.CreatedOn) {
		t.Errorf("handler returned unexpected createdOn: got %v want %v", newContentitem.CreatedOn, now)
	}

	if now.After(newContentitem.UpdatedOn) {
		t.Errorf("handler returned unexpected updatedOn: got %v want %v", newContentitem.UpdatedOn, now)
	}

	eAttr, err := json.Marshal(contentItem.Attrs)
	if err != nil {
		panic(err)
	}
	gAttr, err := json.Marshal(newContentitem.Attrs)
	if err != nil {
		panic(err)
	}
	if string(eAttr) != string(gAttr) {
		t.Errorf("handler returned unexpected Attrs: got %v want %v", contentItem.Attrs, newContentitem.Attrs)
	}
}
