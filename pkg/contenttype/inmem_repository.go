package contenttype

import (
	"sync"

	"github.com/dwethmar/go-api/pkg/common"
)

// InMemRepository mock repository for operating on entry data.
type InMemRepository struct {
	items []*ContentType
	mux   *sync.Mutex
}

// List get all content content.
func (repo *InMemRepository) List() ([]*ContentType, error) {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	return repo.items, nil
}

// Get entry.
func (repo *InMemRepository) Get(id common.ID) (*ContentType, error) {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	for _, n := range repo.items {
		if id == n.ID {
			return n, nil
		}
	}
	return nil, ErrNotFound
}

// Create adds new entry.
func (repo *InMemRepository) Create(entry *ContentType) (common.ID, error) {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	repo.items = append(repo.items, entry)
	return entry.ID, nil
}

// Update Updates entry.
func (repo *InMemRepository) Update(entry *ContentType) error {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	for i, n := range repo.items {
		if entry.ID == n.ID {
			repo.items = append(repo.items[:i], append([]*ContentType{entry}, repo.items[i:]...)...)
			return nil
		}
	}
	return nil
}

// Delete deletes entry.
func (repo *InMemRepository) Delete(id common.ID) error {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	for i, n := range repo.items {
		if id == n.ID {
			repo.items = append(repo.items[:i], repo.items[i+1:]...)
			return nil
		}
	}
	return nil
}

// NewInMemRepository creates new mockservice.
func NewInMemRepository() Repository {
	return &InMemRepository{
		items: make([]*ContentType, 0),
		mux:   &sync.Mutex{},
	}
}
