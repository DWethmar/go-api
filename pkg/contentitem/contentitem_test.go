package contentitem

import (
	"fmt"
	"testing"
)

func TestAttrsValueValidation(t *testing.T) {
	c := ContentItem{
		Name: "Test Name",
		Attrs: Attrs{
			"attr1": 1,
			"attr2": "attribute string value",
			"attr3": []int{1, 2, 3},
			"attr4": []string{"one", "two"},
			"attr5": float64(3) / float64(10),
		},
	}

	if errors := c.Validate(); len(errors) > 0 {
		for _, err := range errors {
			fmt.Printf("Validation error: %+v\n", err)
		}
		t.Errorf("Encountered %v validation errors", len(errors))
	}
}

func TestNameValidation(t *testing.T) {
	c := ContentItem{
		Name: "This name is to loooooooOooooOOo0000000000000000000000000oooong",
	}

	if errors := c.Validate(); len(errors) == 1 {
		for _, e := range errors {
			if e, ok := e.(*NameLengthError); !ok {
				t.Error("Expected name to fail.", e)
			}
		}
	} else {
		t.Errorf("Expected %v errors but got %v error", len(c.Attrs), len(errors))
	}
}

func TestInvalidAttrsValueValidation(t *testing.T) {
	var names []interface{}

	c := ContentItem{
		Name: "Test Name",
		Attrs: Attrs{
			"attr1": nil,
			"attr2": names,
		},
	}

	if errors := c.Validate(); len(errors) == 2 {
		for _, e := range errors {
			if e, ok := e.(*UnsupportedAttrTypeError); !ok {
				t.Error("Validation returned unexpected error:", e)
			}
		}
	} else {
		t.Errorf("Expected %v errors but got %v error", len(c.Attrs), len(errors))
	}
}

func TestInvalidAttrsKeyValidation(t *testing.T) {
	// var names []interface{}

	// c := ContentItem{
	// 	Name: "Test Name",
	// 	Attrs: Attrs{
	// 		"attr 1": nil,
	// 		"attr2":  names,
	// 	},
	// }

	// if errors := c.Validate(); len(errors) == 2 {
	// 	for _, e := range errors {
	// 		if e, ok := e.(*UnsupportedAttrTypeError); !ok {
	// 			t.Error("Validation returned unexpected error:", e)
	// 		}
	// 	}
	// } else {
	// 	t.Errorf("Expected %v errors but got %v error", len(c.Attrs), len(errors))
	// }
}
