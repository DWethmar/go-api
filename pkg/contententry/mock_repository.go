package contententry

type MockRepository struct {
	items []*Entry
}

func (repo *MockRepository) GetAll() ([]*Entry, error) {
	return repo.items, nil
}

func (repo *MockRepository) GetOne(id ID) (*Entry, error) {
	for _, n := range repo.items {
		if id == n.ID {
			return n, nil
		}
	}
	return nil, ErrNotFound
}

func (repo *MockRepository) Add(entry Entry) error {
	repo.items = append(repo.items, &entry)
	return nil
}

func (repo *MockRepository) Update(entry Entry) error {
	for i, n := range repo.items {
		if entry.ID == n.ID {
			repo.items = append(repo.items[:i], append([]*Entry{&entry}, repo.items[i:]...)...)
			return nil
		}
	}
	return nil
}

func (repo *MockRepository) Delete(id ID) error {
	for i, n := range repo.items {
		if id == n.ID {
			repo.items = append(repo.items[:i], repo.items[i+1:]...)
			return nil
		}
	}
	return nil
}

func CreateMockRepository() Repository {
	return &MockRepository{
		items: make([]*Entry, 0),
	}
}
