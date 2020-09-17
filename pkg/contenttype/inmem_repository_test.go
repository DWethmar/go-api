package contenttype

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"gotest.tools/v3/assert"
)

var mock = ContentType{
	ID:   common.NewID(),
	Name: "wow",
	Fields: []*Field{
		{
			ID:           common.NewID(),
			EntryModelID: common.NewID(),
			Key:          "name",
			Name:         "Name",
			FieldType:    "string",
			Length:       10,
			CreatedOn:    time.Now(),
			UpdatedOn:    time.Now(),
		},
	},
	CreatedOn: time.Now(),
	UpdatedOn: time.Now(),
}

func TestInMemGetAll(t *testing.T) {
	c := NewInMemRepository()
	c.Create(&mock)
	c.Create(&mock)

	items, _ := c.List()

	if len(items) != 2 {
		t.Errorf("Expected 2 items but received %v", len(items))
	}
}

func TestInMemGetOne(t *testing.T) {
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

func TestInMemAdd(t *testing.T) {
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

	if string(a) != string(b) {
		t.Errorf("Expected %v but received %v", string(a), string(b))
	}
}

func TestInMemUpdate(t *testing.T) {
	c := NewInMemRepository()
	c.Create(&mock)
	c.Create(&mock)

	createdEntry, err := c.Get(mock.ID)
	if err != nil {
		t.Errorf("Unexpected error")
		return
	}

	updateEntry := *createdEntry
	updateEntry.Fields[0].Name = "new name!"

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

	assert.Equal(t, updatedEntry.Fields[0].Name, "new name!", "Name is not equal")
}

func TestInMemDelete(t *testing.T) {
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
