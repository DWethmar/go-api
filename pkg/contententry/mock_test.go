package contententry

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var mock = Entry{
	ID:   createNewId(),
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

func TestUnitMockGetAll(t *testing.T) {
	c := CreateMockRepository()
	c.Add(mock)
	c.Add(mock)

	items, _ := c.GetAll()

	if len(items) != 2 {
		t.Errorf("Expected 2 items but received %v", len(items))
	}
}

func TestUnitMockGetOne(t *testing.T) {
	c := CreateMockRepository()
	c.Add(mock)
	c.Add(mock)

	item, err := c.GetOne(mock.ID)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if item == nil {
		t.Errorf("Excepted item")
	}

	if item.ID != mock.ID {
		t.Errorf("Expected item with id %v but recieved %v", mock.ID, item.ID)
	}
}

func TestUnitMockAdd(t *testing.T) {
	c := CreateMockRepository()
	c.Add(mock)

	a, err := json.Marshal(mock)
	if err != nil {
		fmt.Println(err)
		return
	}

	createdContentEntry, err := c.GetOne(mock.ID)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	b, err := json.Marshal(createdContentEntry)
	if err != nil {
		t.Errorf("Unexpected error")
		return
	}

	if string(a) != string(b) {
		t.Errorf("Expected %v but received %v", string(a), string(b))
	}
}

func TestUnitMockUpdate(t *testing.T) {
	c := CreateMockRepository()
	c.Add(mock)

	createdEntry, err := c.GetOne(mock.ID)
	if err != nil {
		t.Errorf("Unexpected error")
		return
	}

	updateEntry := *createdEntry
	updateEntry.Fields["nl"]["attr1"] = "new value!"

	err = c.Update(updateEntry)
	if err != nil {
		t.Errorf("Unexpected error while updating entry")
		return
	}

	updatedEntry, err := c.GetOne(mock.ID)
	if err != nil {
		t.Errorf("Unexpected error")
		return
	}

	if updatedEntry.Fields["nl"]["attr1"] != "new value!" {
		t.Errorf("Expected fields.nl.attr1 to be \"new value!\"")
	}
}