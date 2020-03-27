package contentitem

type MockRepository struct {
	items []*ContentItem
}

func (repo *MockRepository) GetAll() ([]*ContentItem, error) {
	return repo.items, nil
}

func (repo *MockRepository) GetOne(id ID) (*ContentItem, error) {
	for _, n := range repo.items {
		if id == n.ID {
			return n, nil
		}
	}
	return nil, ErrNotFound
}

func (repo *MockRepository) Add(contentItem ContentItem) error {
	repo.items = append(repo.items, &contentItem)
	return nil
}

func (repo *MockRepository) Update(contentItem ContentItem) error {
	for i, n := range repo.items {
		if contentItem.ID == n.ID {
			repo.items = append(repo.items[:i], append([]*ContentItem{&contentItem}, repo.items[i:]...)...)
			return nil
		}
	}
	return nil
}

func (repo *MockRepository) Delete(id ID) error {
	for i, n := range repo.items {
		if id == n.ID {
			repo.items = append(repo.items[:i], repo.items[i:]...)
			return nil
		}
	}
	return nil
}

func CreateMockRepository() Repository {
	return &MockRepository{
		items: make([]*ContentItem, 0),
	}
}
