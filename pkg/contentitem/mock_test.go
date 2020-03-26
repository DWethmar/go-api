package contentitem

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var contentItem1 = ContentItem{
	ID:   createNewId(),
	Name: "wow",
	Attrs: Attrs{
		"nl": {
			"attr1": 1,
			"attr2": "attribute string value",
			"attr3": []int{1, 2, 3},
			"attr4": []string{"one", "two"},
			"attr5": float64(3) / float64(10),
		},
	},
	CreatedOn: time.Now(),
	UpdatedOn: time.Now(),
}

func TestMockCreateAndGet(t *testing.T) {
	c := CreateMockRepository()
	c.Add(contentItem1)

	a, err := json.Marshal(contentItem1)
	if err != nil {
		fmt.Println(err)
		return
	}

	createdContentItem, err := c.GetOne(contentItem1.ID)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	b, err := json.Marshal(createdContentItem)
	if err != nil {
		fmt.Println(err)
		return
	}

	if string(a) != string(b) {
		t.Errorf("Expected: \n%v \ngot: \n%v", string(a), string(b))
	}
}

func TestMockGetAll(t *testing.T) {
	c := CreateMockRepository()
	c.Add(contentItem1)
	c.Add(contentItem1)

	items, _ := c.GetAll()

	if len(items) != 2 {
		t.Errorf("Expected: 2 items got: %v", len(items))
	}
}
