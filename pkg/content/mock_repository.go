package content

import (
	"sync"

	"github.com/dwethmar/go-api/pkg/common"
)

// MockRepository mock repository for operating on entry data.
type MockRepository struct {
	items []*Content
	mux   *sync.Mutex
}

// List entries.
func (repo *MockRepository) List() ([]*Content, error) {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	return repo.items, nil
}

// Get one entry.
func (repo *MockRepository) Get(id common.ID) (*Content, error) {
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
func (repo *MockRepository) Create(entry *Content) (common.ID, error) {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	repo.items = append(repo.items, entry)
	return entry.ID, nil
}

// Update Updates entry.
func (repo *MockRepository) Update(entry *Content) error {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	for i, n := range repo.items {
		if entry.ID == n.ID {
			repo.items = append(repo.items[:i], append([]*Content{entry}, repo.items[i:]...)...)
			return nil
		}
	}

	return nil
}

// Delete deletes entry.
func (repo *MockRepository) Delete(id common.ID) error {
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

// NewMockRepository creates new mockservice.
func NewMockRepository() Repository {
	return &MockRepository{
		items: make([]*Content, 0),
		mux:   &sync.Mutex{},
	}
}
