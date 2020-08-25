package store

import (
	"database/sql"

	"github.com/dwethmar/go-api/pkg/content"
	"github.com/dwethmar/go-api/pkg/contenttype"
)

// Store is a collection of services to create, read, update and delete data.
type Store struct {
	Content     content.Service
	ContentType contenttype.Service
}

// CreateStoreOption are used to compose a store with provided services.
type CreateStoreOption struct {
	ContentRepo     content.Repository
	ContentTypeRepo contenttype.Repository
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
		ContentRepo:     content.NewMockRepository(),
		ContentTypeRepo: contenttype.NewMockRepository(),
	})
}

// NewStore creates a store with services that use persistsent storage.
func NewStore(db *sql.DB) *Store {
	return CreateStoreWithOption(CreateStoreOption{
		ContentRepo:     content.NewPostgresRepository(db),
		ContentTypeRepo: contenttype.NewPostgresRepository(db),
	})
}
