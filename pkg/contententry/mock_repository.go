package contententry

// MockRepository mock repository for operating on entry data.
type MockRepository struct {
	items []*Entry
}

// GetAll get all entries.
func (repo *MockRepository) GetAll() ([]*Entry, error) {
	return repo.items, nil
}

// GetOne get one entry.
func (repo *MockRepository) GetOne(id ID) (*Entry, error) {
	for _, n := range repo.items {
		if id == n.ID {
			return n, nil
		}
	}
	return nil, ErrNotFound
}

// Add add new entry.
func (repo *MockRepository) Add(entry Entry) error {
	repo.items = append(repo.items, &entry)
	return nil
}

// Update Updates entry.
func (repo *MockRepository) Update(entry Entry) error {
	for i, n := range repo.items {
		if entry.ID == n.ID {
			repo.items = append(repo.items[:i], append([]*Entry{&entry}, repo.items[i:]...)...)
			return nil
		}
	}
	return nil
}

// Delete deletes entry.
func (repo *MockRepository) Delete(id ID) error {
	for i, n := range repo.items {
		if id == n.ID {
			repo.items = append(repo.items[:i], repo.items[i+1:]...)
			return nil
		}
	}
	return nil
}

// CreateMockRepository creates new mockservice.
func CreateMockRepository() Repository {
	return &MockRepository{
		items: make([]*Entry, 0),
	}
}
