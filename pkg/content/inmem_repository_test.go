package content

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"gotest.tools/v3/assert"
)

var mock = Content{
	ID:   common.NewID(),
	Name: "wow",
	Fields: FieldTranslations{
		"nl": Fields{
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

func TestInMemRepository_List(t *testing.T) {
	c := NewInMemRepository()
	c.Create(&mock)
	c.Create(&mock)

	items, _ := c.List()

	if len(items) != 2 {
		t.Errorf("Expected 2 items but received %v", len(items))
	}
}

func TestInMemRepository_Get(t *testing.T) {
	c := NewInMemRepository()
	c.Create(&mock)
	c.Create(&mock)

	item, err := c.Get(mock.ID)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if item == nil {
		t.Errorf("Excepted item")
	}

	if item.ID != mock.ID {
		t.Errorf("Expected item with id %v but received %v", mock.ID, item.ID)
	}
}

func TestInMemRepository_Create(t *testing.T) {
	c := NewInMemRepository()
	c.Create(&mock)
	c.Create(&mock)

	a, err := json.Marshal(mock)
	if err != nil {
		fmt.Println(err)
		return
	}

	createdEntry, err := c.Get(mock.ID)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	b, err := json.Marshal(createdEntry)
	if err != nil {
		t.Errorf("Unexpected error")
		return
	}

	assert.Equal(t, string(a), string(b), "Entries not the same")
}

func TestInMemRepository_Update(t *testing.T) {
	c := NewInMemRepository()
	c.Create(&mock)
	c.Create(&mock)

	createdEntry, err := c.Get(mock.ID)
	if err != nil {
		t.Errorf("Unexpected error")
		return
	}

	updateEntry := *createdEntry
	updateEntry.Fields["nl"]["attr1"] = "new value!"

	err = c.Update(&updateEntry)
	if err != nil {
		t.Errorf("Unexpected error while updating entry")
		return
	}

	updatedEntry, err := c.Get(mock.ID)
	if err != nil {
		t.Errorf("Unexpected error")
		return
	}

	if updatedEntry.Fields["nl"]["attr1"] != "new value!" {
		t.Errorf("Expected fields.nl.attr1 to be \"new value!\"")
	}
}

func TestInMemRepository_Delete(t *testing.T) {
	c := NewInMemRepository()
	c.Create(&mock)

	createdEntry, err := c.Get(mock.ID)
	if err != nil {
		t.Errorf("Unexpected error")
		return
	}

	err = c.Delete(createdEntry.ID)
	if err != nil {
		t.Errorf("Unexpected error while deleting entry")
		return
	}

	_, err = c.Get(mock.ID)
	if err != nil {
		if err != ErrNotFound {
			t.Errorf("Unexpected entry")
		}
		return
	}
}
