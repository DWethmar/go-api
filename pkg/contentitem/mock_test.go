package contentitem

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var contentItem1 = ContentItem{
	ID:   1,
	Name: "wow",
	Attrs: Attrs{
		"attr1": 1,
		"attr2": "attribute string value",
		"attr3": []int{1, 2, 3},
		"attr4": []string{"one", "two"},
		"attr5": float64(3) / float64(10),
	},
	CreatedOn: time.Now(),
	UpdatedOn: time.Now(),
}

func TestMockCreateAndGet(t *testing.T) {
	c := CreateMockRepository()
	c.Create(contentItem1)
	a, err := json.Marshal(contentItem1)

	if err != nil {
		fmt.Println(err)
		return
	}

	createdContentItem, _ := c.GetOne(1)

	b, err := json.Marshal(createdContentItem)
	if err != nil {
		fmt.Println(err)
		return
	}

	if string(a) != string(b) {
		t.Errorf("Expected: %v got: %v", a, b)
	}
}

func TestMockGetAll(t *testing.T) {
	c := CreateMockRepository(contentItem1, contentItem1)
	items, _ := c.GetAll()

	if len(items) != 2 {
		t.Errorf("Expected: 2 items got: %v", len(items))
	}
}
