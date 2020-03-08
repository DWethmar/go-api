package contentitem

import (
	"fmt"
	"testing"
)

func TestValidation(t *testing.T) {
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
