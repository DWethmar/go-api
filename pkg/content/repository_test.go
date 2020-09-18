package content

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/database"
	"gotest.tools/v3/assert"
)

var defaultLocale = "nl"

func TestRepository_List(t *testing.T) {

	var repos []Repository

	c := config.Load()
	if c.TestWithDB {
		db, err := database.NewTestDB(c)
		defer db.Close()

		if err != nil {
			t.Error(err)
		}
		repos = append(repos, NewPostgresRepository(db))
	} else {
		fmt.Printf("Test without PostgresRepository")
	}

	repos = append(repos, NewInMemRepository())

	test := func(repo Repository) {
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
			_, err := repo.Create(newEntry)
			if err != nil {
				t.Error(err)
			}
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

	for _, repo := range repos {
		test(repo)
	}
}

func TestRepository_Get(t *testing.T) {
	var repos []Repository

	c := config.Load()
	if c.TestWithDB {
		db, err := database.NewTestDB(c)
		defer db.Close()

		if err != nil {
			t.Error(err)
		}
		repos = append(repos, NewPostgresRepository(db))
	} else {
		fmt.Printf("Test without PostgresRepository")
	}

	repos = append(repos, NewInMemRepository())

	test := func(repo Repository) {
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

		// ErrNotFound
		// Test get on non existing thing
		_, err = repo.Get(common.NewID())
		if err != ErrNotFound {
			t.Errorf("Excpected not found error.")
		}
	}

	for _, repo := range repos {
		test(repo)
	}
}

func TestRepository_Create(t *testing.T) {
	var repos []Repository

	c := config.Load()
	if c.TestWithDB {
		db, err := database.NewTestDB(c)
		defer db.Close()

		if err != nil {
			t.Error(err)
		}
		repos = append(repos, NewPostgresRepository(db))
	} else {
		fmt.Printf("Test without PostgresRepository")
	}

	repos = append(repos, NewInMemRepository())

	test := func(repo Repository) {
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

	for _, repo := range repos {
		test(repo)
	}
}

func TestRepository_Update(t *testing.T) {
	var repos []Repository

	c := config.Load()
	if c.TestWithDB {
		db, err := database.NewTestDB(c)
		defer db.Close()

		if err != nil {
			t.Error(err)
		}
		repos = append(repos, NewPostgresRepository(db))
	} else {
		fmt.Printf("Test without PostgresRepository")
	}

	repos = append(repos, NewInMemRepository())

	test := func(repo Repository) {

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
			_, err := repo.Create(newEntry)
			if err != nil {
				t.Error(err)
			}
		}

		addItems[0].Name = "Updated name"
		addItems[0].Fields[defaultLocale]["attrA"] = 2

		err := repo.Update(addItems[0])
		if err != nil {
			t.Error(err)
		}

		updatedEntry, err := repo.Get(addItems[0].ID)
		if err != nil {
			t.Error(err)
		}

		received, err := json.Marshal(updatedEntry)
		if err != nil {
			t.Error(err)
		}

		expected, err := json.Marshal(addItems[0])
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, string(received), string(expected), "Didn't expect value.")

		// Check if the other entry is unaffected
		otherEntry, err := repo.Get(addItems[1].ID)
		if err != nil {
			t.Error(err)
		}

		received, err = json.Marshal(otherEntry)
		if err != nil {
			t.Error(err)
		}

		expected, err = json.Marshal(addItems[1])
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, string(received), string(expected), "Didn't expect value.")
	}

	for _, repo := range repos {
		test(repo)
	}
}

func TestRepository_Delete(t *testing.T) {
	var repos []Repository

	c := config.Load()
	if c.TestWithDB {
		db, err := database.NewTestDB(c)
		defer db.Close()

		if err != nil {
			t.Error(err)
		}
		repos = append(repos, NewPostgresRepository(db))
	} else {
		fmt.Printf("Test without PostgresRepository")
	}

	repos = append(repos, NewInMemRepository())

	test := func(repo Repository) {

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

		err := repo.Delete(addItems[0].ID)
		if err != nil {
			t.Error(err)
		}

		_, err = repo.Get(addItems[1].ID)
		if err != nil {
			t.Error(err)
		}

		// Delete entry that doesn't exists anymore.
		err = repo.Delete(addItems[0].ID)
		if err != nil {
			t.Error(err)
		}
	}

	for _, repo := range repos {
		test(repo)
	}
}
