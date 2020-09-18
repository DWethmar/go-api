package content

import (
	"encoding/json"
	"testing"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/database"
	"gotest.tools/v3/assert"
)

var defaultLocale = "nl"

func TestPostgresRepository_List(t *testing.T) {
	c := config.Load()
	if !c.TestWithDB {
		t.Skip("skipping test case without db")
	}
	db, err := database.NewTestDB(c)
	defer db.Close()

	if err != nil {
		t.Error(err)
	}

	repo := NewPostgresRepository(db)

	// Test empty repo
	entries, err := repo.List()
	received, _ := json.Marshal(entries)
	assert.Equal(t, string(received), "[]", "Didn't expect value")

	addItems := []*Content{
		{
			ID:   common.NewID(),
			Name: "Test1",
			Fields: FieldTranslations{
				defaultLocale: Fields{
					"attrA": 1,
				},
			},
			CreatedOn: common.Now(),
			UpdatedOn: common.Now(),
		},
		{
			ID:   common.NewID(),
			Name: "Test2",
			Fields: FieldTranslations{
				defaultLocale: Fields{
					"attrA": 1,
				},
			},
			CreatedOn: common.Now(),
			UpdatedOn: common.Now(),
		},
	}

	for _, newEntry := range addItems {
		repo.Create(newEntry)
	}

	entries, err = repo.List()
	if err != nil {
		t.Error(err)
	}

	received, err = json.Marshal(entries)
	if err != nil {
		t.Error(err)
	}

	expected, err := json.Marshal(addItems)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(received), string(expected), "Didn't expect value")
}

func TestPostgresRepository_Get(t *testing.T) {
	c := config.Load()
	if !c.TestWithDB {
		t.Skip("skipping test case without db")
	}
	db, err := database.NewTestDB(c)
	defer db.Close()

	if err != nil {
		t.Error(err)
	}

	repo := NewPostgresRepository(db)

	addEntry := &Content{
		ID:   common.NewID(),
		Name: "Test1",
		Fields: FieldTranslations{
			defaultLocale: Fields{
				"attrA": 1,
			},
		},
		CreatedOn: common.Now(),
		UpdatedOn: common.Now(),
	}

	ID, err := repo.Create(addEntry)
	if err != nil {
		t.Error(err)
	}

	createdEntry, err := repo.Get(ID)
	if err != nil {
		t.Error(err)
	}

	received, err := json.Marshal(createdEntry)
	if err != nil {
		t.Error(err)
	}

	expected, err := json.Marshal(addEntry)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(received), string(expected), "Didn't expect value.")
}

func TestPostgresRepository_Create(t *testing.T) {
	c := config.Load()
	if !c.TestWithDB {
		t.Skip("skipping test case without db")
	}
	db, err := database.NewTestDB(c)
	defer db.Close()

	if err != nil {
		t.Error(err)
	}

	repo := NewPostgresRepository(db)

	addEntry := &Content{
		ID:   common.NewID(),
		Name: "Test1",
		Fields: FieldTranslations{
			defaultLocale: Fields{
				"attrA": 1,
			},
		},
		CreatedOn: common.Now(),
		UpdatedOn: common.Now(),
	}

	ID, err := repo.Create(addEntry)
	if err != nil {
		t.Error(err)
	}

	createdEntry, err := repo.Get(ID)
	if err != nil {
		t.Error(err)
	}

	received, err := json.Marshal(createdEntry)
	if err != nil {
		t.Error(err)
	}

	expected, err := json.Marshal(addEntry)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(received), string(expected), "Didn't expect value.")
}

func TestPostgresRepository_Update(t *testing.T) {

}

func TestPostgresRepository_Delete(t *testing.T) {

}
