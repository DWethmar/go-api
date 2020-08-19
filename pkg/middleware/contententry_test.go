package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwethmar/go-api/pkg/contententry"
	"github.com/dwethmar/go-api/pkg/store"
)

func TestHandleEntryIndex(t *testing.T) {
	addItems := []contententry.AddEntry{
		{
			Name: "Test1",
			Fields: contententry.FieldTranslations{
				"nl": contententry.Fields{
					"attrA": 1,
				},
			},
		},
		{
			Name: "Test2",
			Fields: contententry.FieldTranslations{
				"nl": contententry.Fields{
					"attrA": 1,
				},
			},
		},
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	entries := []*contententry.Entry{}
	rr := httptest.NewRecorder()

	withStore(func(store *store.Store) {
		for _, newEntry := range addItems {
			entry, _ := store.Entries.Create(newEntry)
			entries = append(entries, entry)
		}

		handler := http.HandlerFunc(HandleEntryIndex(store))
		handler.ServeHTTP(rr, req)
	})

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected, _ := json.Marshal(entries)

	if rr.Body.String() != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(expected))
	}
}
