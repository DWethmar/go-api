package content

import (
	"sync"

	"github.com/dwethmar/go-api/pkg/common"
)

// InMemRepository mock repository for operating on entry data.
type InMemRepository struct {
	items []*Content
	mux   *sync.Mutex
}

// List entries.
func (repo *InMemRepository) List() ([]*Content, error) {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	return repo.items, nil
}

// Get one entry.
func (repo *InMemRepository) Get(id common.ID) (*Content, error) {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	for _, n := range repo.items {
		if id == n.ID {
			return n, nil
		}
	}
	return nil, ErrNotFound
}

// Create new entry.
func (repo *InMemRepository) Create(entry *Content) (common.ID, error) {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	repo.items = append(repo.items, entry)
	return entry.ID, nil
}

// Update Updates entry.
func (repo *InMemRepository) Update(entry *Content) error {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	for i, n := range repo.items {
		if entry.ID == n.ID {
			repo.items = append(repo.items[:i], append([]*Content{entry}, repo.items[i:]...)...)
			return nil
		}
	}

	return ErrNotFound
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

	return ErrNotFound
}

// NewInMemRepository creates new mockservice.
func NewInMemRepository() Repository {
	return &InMemRepository{
		items: make([]*Content, 0),
		mux:   &sync.Mutex{},
	}
}
