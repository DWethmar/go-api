package store

import (
	"database/sql"

	"github.com/dwethmar/go-api/pkg/contententry"
)

// Store is a collection of services to create, read, update and delete data.
type Store struct {
	Entries contententry.Service
}

// CreateStoreOption are used to compose a store with provided services.
type CreateStoreOption struct {
	EntryRepo contententry.Repository
}

// CreateStoreWithOption creates a store with the provided options.
func CreateStoreWithOption(options CreateStoreOption) *Store {
	return &Store{
		Entries: contententry.CreateService(options.EntryRepo),
	}
}

// CreateMockStore creates a store with mock services.
func CreateMockStore() *Store {
	return CreateStoreWithOption(CreateStoreOption{
		EntryRepo: contententry.CreateMockRepository(),
	})
}

// CreateStore creates a store with services that use persistsent storage.
func CreateStore(db *sql.DB) *Store {
	return CreateStoreWithOption(CreateStoreOption{
		EntryRepo: contententry.CreatePostgresRepository(db),
	})
}
