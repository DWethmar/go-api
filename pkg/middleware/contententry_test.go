package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/contententry"
	"github.com/dwethmar/go-api/pkg/store"
)

var defaultLocale = "nl"

func areFieldsEqual(a, b contententry.FieldTranslations) (bool, error) {
	ar, err := json.Marshal(a)
	if err != nil {
		return false, err
	}

	br, err := json.Marshal(b)
	if err != nil {
		return false, err
	}

	return string(ar) != string(br), nil
}

func areEntriesEqual(a, b contententry.Entry) (bool, error) {
	ar, err := json.Marshal(a)
	if err != nil {
		return false, err
	}

	br, err := json.Marshal(b)
	if err != nil {
		return false, err
	}

	return string(ar) != string(br), nil
}

func TestEntryIndex(t *testing.T) {
	addItems := []contententry.AddEntry{
		{
			Name: "Test1",
			Fields: contententry.FieldTranslations{
				defaultLocale: contententry.Fields{
					"attrA": 1,
				},
			},
		},
		{
			Name: "Test2",
			Fields: contententry.FieldTranslations{
				defaultLocale: contententry.Fields{
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

	store.WithStore(func(store *store.Store) {
		for _, newEntry := range addItems {
			entry, _ := store.Entries.Create(newEntry)
			entries = append(entries, entry)
		}

		handler := http.HandlerFunc(EntryIndex(store))
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

func TestCreateEntry(t *testing.T) {
	now := time.Now()

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

	store.WithStore(func(store *store.Store) {
		handler := http.HandlerFunc(CreateEntry(store))
		handler.ServeHTTP(rr, req)
	})

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: received %v, excepted %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	addedEntry := contententry.Entry{}
	err := json.Unmarshal(rr.Body.Bytes(), &addedEntry)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	if now.After(addedEntry.CreatedOn) {
		t.Errorf("handler returned invalid createdOn: received %v, excepted CreatedOn to be smaller then %v", addedEntry.CreatedOn, now)
	}

	if now.After(addedEntry.UpdatedOn) {
		t.Errorf("handler returned invalid updatedOn: received %v, excepted UpdatedOn to be smaller then  %v", addedEntry.UpdatedOn, now)
	}

	if equal, err := areFieldsEqual(addEntry.Fields, addedEntry.Fields); equal || err != nil {
		if err != nil {
			t.Error(err)
		} else {
			t.Errorf("Fields are not equal. left: %v right %v", addEntry.Fields, addedEntry.Fields)
		}
	}
}

func TestUpdateEntry(t *testing.T) {

	store.WithStore(func(store *store.Store) {
		addedEntry, _ := store.Entries.Create(contententry.AddEntry{
			Name: "name",
			Fields: contententry.FieldTranslations{
				"nl": contententry.Fields{
					"attrA": "Value",
					"attrB": 1,
					"attrC": []string{"one", "two", "three"},
				},
			},
		})

		updateEntry := contententry.Entry{
			ID:   addedEntry.ID,
			Name: "updated name",
			Fields: contententry.FieldTranslations{
				"nl": contententry.Fields{
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
		ctx = common.WithUUID(ctx, addedEntry.ID)
		req = req.WithContext(ctx)

		handler := http.HandlerFunc(UpdateEntry(store))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusOK, req.RequestURI)
			return
		}

		// Check the response content-type is what we expect.
		if rType := rr.Header().Get("Content-Type"); rType != "application/json" {
			t.Errorf("content type header does not match: received %v want %v", rType, "application/json")
			return
		}

		// Check the response body is what we expect.
		updatedEntry := contententry.Entry{}
		err := json.Unmarshal(rr.Body.Bytes(), &updatedEntry)
		if err != nil {
			t.Errorf("Error while parsing body for added content-item %v", err)
			return
		}

		if equal, err := areFieldsEqual(updateEntry.Fields, updatedEntry.Fields); equal || err != nil {
			if err != nil {
				t.Error(err)
			} else {
				t.Errorf("Fields are not equal. left: %v right %v", updateEntry.Fields, updatedEntry.Fields)
			}
		}
	})
}

// func TestDeleteEntry(t *testing.T) {

// 	store.WithStore(func(store *store.Store) {
// 		addedEntry, _ := store.Entries.Create(contententry.AddEntry{
// 			Name: "name",
// 			Fields: contententry.FieldTranslations{
// 				"nl": contententry.Fields{
// 					"attrA": "Value",
// 					"attrB": 1,
// 					"attrC": []string{"one", "two", "three"},
// 				},
// 			},
// 		})

// 		body, _ := json.Marshal(addedEntry)
// 		req := httptest.NewRequest("DELETE", fmt.Sprintf("/%s", addedEntry.ID), bytes.NewBuffer(body))
// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(DeleteEntry(store))
// 		handler.ServeHTTP(rr, req)

// 		if status := rr.Code; status != http.StatusCreated {
// 			t.Errorf("handler returned wrong status code: received %v expected %v on %v", status, http.StatusCreated, req.RequestURI)
// 			return
// 		}

// 		// Check the response content-type is what we expect.
// 		if rType := rr.Header().Get("Content-Type"); rType != "application/json" {
// 			t.Errorf("content type header does not match: received %v want %v", rType, "application/json")
// 			return
// 		}

// 		deletedEntry := contententry.Entry{}
// 		err := json.Unmarshal(rr.Body.Bytes(), &deletedEntry)
// 		if err != nil {
// 			t.Errorf("Error while parsing body for deleted entry %v", err)
// 			return
// 		}

// 		if equal, err := areEntriesEqual(*addedEntry, deletedEntry); equal || err != nil {
// 			if err != nil {
// 				t.Error(err)
// 			} else {
// 				t.Errorf("Entries are not equal. left: %v right %v", addedEntry, deletedEntry)
// 			}
// 		}
// 	})
// }
