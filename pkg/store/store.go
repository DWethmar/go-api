package store

import (
	"database/sql"

	"github.com/dwethmar/go-api/pkg/contententry"
)

type Store struct {
	Entries contententry.Service
}

type CreateStoreOption struct {
	EntryRepo contententry.Repository
}

func CreateStoreWithOption(options CreateStoreOption) *Store {
	return &Store{
		Entries: contententry.CreateService(options.EntryRepo),
	}
}

func CreateStore(db *sql.DB) *Store {
	return CreateStoreWithOption(CreateStoreOption{
		EntryRepo: contententry.CreatePostgresRepository(db),
	})
}
