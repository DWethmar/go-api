package handler

import (
	"encoding/json"

	"github.com/dwethmar/go-api/pkg/content"
)

var defaultLocale = "nl"

func areFieldsEqual(a, b content.FieldTranslations) (bool, error) {
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

func areEntriesEqual(a, b content.Content) (bool, error) {
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

// func timeAfter(t *testing.T) {
// 	if now.After(addedEntry.CreatedOn) {
// 		t.Errorf("handler returned invalid createdOn: received %v, excepted CreatedOn to be smaller then %v", addedEntry.CreatedOn, now)
// 	}
// }
