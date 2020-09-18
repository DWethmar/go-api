package content

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/database"
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

	addItems := []*Content{
		{
			ID:   common.NewID(),
			Name: "Test1",
			Fields: FieldTranslations{
				defaultLocale: Fields{
					"attrA": 1,
				},
			},
			CreatedOn: time.Now().Truncate(time.Microsecond),
			UpdatedOn: time.Now().Truncate(time.Microsecond),
		},
		{
			ID:   common.NewID(),
			Name: "Test2",
			Fields: FieldTranslations{
				defaultLocale: Fields{
					"attrA": 1,
				},
			},
			CreatedOn: time.Now().Truncate(time.Microsecond),
			UpdatedOn: time.Now().Truncate(time.Microsecond),
		},
	}

	entries := []*Content{}
	for _, newEntry := range addItems {
		ID, _ := repo.Create(newEntry)
		entry, err := repo.Get(ID)
		if err != nil {
			t.Errorf("something went wrong %v", err)
		}
		entries = append(entries, entry)
	}

	received, _ := json.Marshal(entries)
	expected, _ := json.Marshal(addItems)

	if string(received) != string(expected) {
		t.Errorf("handler returned unexpected body: received %v expected %v", string(received), string(expected))
	}
}

func TestPostgresRepository_Get(t *testing.T) {

}

func TestPostgresRepository_Create(t *testing.T) {

}

func TestPostgresRepository_Update(t *testing.T) {

}

func TestPostgresRepository_Delete(t *testing.T) {

}
