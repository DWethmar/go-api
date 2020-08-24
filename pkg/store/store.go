package store

import (
	"database/sql"

	"github.com/dwethmar/go-api/pkg/content"
)

// Store is a collection of services to create, read, update and delete data.
type Store struct {
	Content content.Service
}

// CreateStoreOption are used to compose a store with provided services.
type CreateStoreOption struct {
	ContentRepo content.Repository
}

// CreateStoreWithOption creates a store with the provided options.
func CreateStoreWithOption(options CreateStoreOption) *Store {
	return &Store{
		Content: content.NewService(options.ContentRepo),
	}
}

// CreateMockStore creates a store with mock services.
func CreateMockStore() *Store {
	return CreateStoreWithOption(CreateStoreOption{
		ContentRepo: content.NewMockRepository(),
	})
}

// CreateStore creates a store with services that use persistsent storage.
func CreateStore(db *sql.DB) *Store {
	return CreateStoreWithOption(CreateStoreOption{
		ContentRepo: content.NewPostgresRepository(db),
	})
}
