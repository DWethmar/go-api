package contentitem

import (
	"fmt"
	"testing"
)

func TestAttrsValues(t *testing.T) {
	c := Attrs{
		"nl": {
			"attr1": 1,
			"attr2": "attribute string value",
			"attr3": []int{1, 2, 3},
			"attr4": []string{"one", "two"},
			"attr5": float64(3) / float64(10),
		},
	}

	if errors := ValidateAttr(c["nl"]); len(errors) > 0 {
		for key, err := range errors {
			fmt.Printf("Validation error on attr %v %+v\n", key, err)
		}
		t.Errorf("Encountered %v validation errors", len(errors))
	}
}

func TestNameValidation(t *testing.T) {
	c := ContentItem{
		Name: "This name is to loooooooOooooOOo0000000000000000000000000oooong",
	}

	if err := ValidateName(c.Name); err == nil {
		t.Errorf("Expected error")
	} else {
		if err != ErrNameLength {
			t.Error("Unexpected error.", err)
		}
	}
}

func TestInvalidAttrsValues(t *testing.T) {
	var names []interface{}

	c := Attrs{
		"nl": {
			"attr1": nil,
			"attr2": names,
		},
	}

	if errors := ValidateAttr(c["nl"]); len(errors) == 2 {
		for attr, err := range errors {
			if err != ErrUnsupportedAttrType {
				t.Errorf("Validation returned unexpected error on attr %v with error %v:", attr, err)
			}
		}
	} else {
		t.Errorf("Expected %v errors but got %v error", len(c), len(errors))
	}
}
