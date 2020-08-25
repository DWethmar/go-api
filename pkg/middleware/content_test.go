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
	"github.com/dwethmar/go-api/pkg/models"
	"github.com/dwethmar/go-api/pkg/store"
)

var defaultLocale = "nl"

func areEntryValidationErrorsEqual(a, b models.ErrContentValidation) (bool, error) {
	ar, err := json.Marshal(a)
	if err != nil {
		return false, err
	}

	br, err := json.Marshal(b)
	if err != nil {
		return false, err
	}

	return string(ar) == string(br), nil
}

func areFieldsEqual(a, b models.FieldTranslations) (bool, error) {
	ar, err := json.Marshal(a)
	if err != nil {
		return false, err
	}

	br, err := json.Marshal(b)
	if err != nil {
		return false, err
	}

	return string(ar) == string(br), nil
}

func areEntriesEqual(a, b models.Content) (bool, error) {
	ar, err := json.Marshal(a)
	if err != nil {
		return false, err
	}

	br, err := json.Marshal(b)
	if err != nil {
		return false, err
	}

	return string(ar) == string(br), nil
}

func TestContentIndex(t *testing.T) {
	addItems := []models.AddContent{
		{
			Name: "Test1",
			Fields: models.FieldTranslations{
				defaultLocale: models.Fields{
					"attrA": 1,
				},
			},
		},
		{
			Name: "Test2",
			Fields: models.FieldTranslations{
				defaultLocale: models.Fields{
					"attrA": 1,
				},
			},
		},
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	entries := []*models.Content{}
	rr := httptest.NewRecorder()

	store.WithTestStore(func(store *store.Store) {
		for _, newEntry := range addItems {
			entry, _ := store.Content.Create(newEntry)
			entries = append(entries, entry)
		}

		handler := http.HandlerFunc(ContentIndex(store))
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

func TestCreateContent(t *testing.T) {
	now := time.Now()

	addEntry := models.AddContent{
		Name: "name",
		Fields: models.FieldTranslations{
			"nl": models.Fields{
				"attrA": "Value",
				"attrB": 1,
				"attrC": []string{"one", "two", "three"},
			},
		},
	}

	body, _ := json.Marshal(addEntry)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()

	store.WithTestStore(func(store *store.Store) {
		handler := http.HandlerFunc(CreateContent(store))
		handler.ServeHTTP(rr, req)
	})

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: received %v, excepted %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	addedEntry := models.Content{}
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

	if equal, err := areFieldsEqual(addEntry.Fields, addedEntry.Fields); !equal || err != nil {
		if err != nil {
			t.Error(err)
		} else {
			t.Errorf("Fields are not equal. left: %v right %v", addEntry.Fields, addedEntry.Fields)
		}
	}
}

func TestCreateInvalidContent(t *testing.T) {
	addEntry := models.AddContent{
		Name: "name",
		Fields: models.FieldTranslations{
			"nl": models.Fields{
				"attrA": "Value",
				"attrB": 1,
				"attrC": nil,
			},
		},
	}

	validationErr := models.CreateContentValidationError()
	validationErr.Errors.Fields["nl"] = make(map[string]string)
	validationErr.Errors.Fields["nl"]["attrC"] = models.ErrUnsupportedFieldValue.Error()

	body, _ := json.Marshal(addEntry)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()

	store.WithTestStore(func(store *store.Store) {
		handler := http.HandlerFunc(CreateContent(store))
		handler.ServeHTTP(rr, req)
	})

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: received %v, excepted %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect.
	receivedValidationError := models.ErrContentValidation{}
	err := json.Unmarshal(rr.Body.Bytes(), &receivedValidationError)
	if err != nil {
		t.Errorf("Error while parsing body %v", err)
	}

	if equal, err := areEntryValidationErrorsEqual(receivedValidationError, validationErr); !equal || err != nil {
		if err != nil {
			t.Error(err)
		} else {
			t.Errorf("Fields are not equal. left: %v right %v", receivedValidationError, validationErr)
		}
	}
}

func TestUpdateContent(t *testing.T) {
	store.WithTestStore(func(store *store.Store) {
		addedEntry, _ := store.Content.Create(models.AddContent{
			Name: "name",
			Fields: models.FieldTranslations{
				"nl": models.Fields{
					"attrA": "Value",
					"attrB": 1,
					"attrC": []string{"one", "two", "three"},
				},
			},
		})

		updateEntry := models.Content{
			ID:   addedEntry.ID,
			Name: "updated name",
			Fields: models.FieldTranslations{
				"nl": models.Fields{
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

		handler := http.HandlerFunc(UpdateContent(store))
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
		updatedEntry := models.Content{}
		err := json.Unmarshal(rr.Body.Bytes(), &updatedEntry)
		if err != nil {
			t.Errorf("Error while parsing body for added content-item %v", err)
			return
		}

		if equal, err := areFieldsEqual(updateEntry.Fields, updatedEntry.Fields); !equal || err != nil {
			if err != nil {
				t.Error(err)
			} else {
				t.Errorf("Fields are not equal. left: %v right %v", updateEntry.Fields, updatedEntry.Fields)
			}
		}
	})
}

func TestDeleteContent(t *testing.T) {
	store.WithTestStore(func(store *store.Store) {
		addedEntry, _ := store.Content.Create(models.AddContent{
			Name: "name",
			Fields: models.FieldTranslations{
				"nl": models.Fields{
					"attrA": "Value",
					"attrB": 1,
					"attrC": []string{"one", "two", "three"},
				},
			},
		})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/%s", addedEntry.ID), nil)

		ctx := req.Context()
		ctx = common.WithUUID(ctx, addedEntry.ID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(DeleteContent(store))
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

		deletedEntry := models.Content{}
		err := json.Unmarshal(rr.Body.Bytes(), &deletedEntry)
		if err != nil {
			t.Errorf("Error while parsing body for deleted entry %v", err)
			return
		}

		if equal, err := areEntriesEqual(*addedEntry, deletedEntry); !equal || err != nil {
			if err != nil {
				t.Error(err)
			} else {
				t.Errorf("Entries are not equal. left: %v right %v", addedEntry, deletedEntry)
			}
		}
	})
}