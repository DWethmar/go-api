package contentitem

import "database/sql"

type MockRepository struct {
	items []ContentItem
}

func (repo *MockRepository) GetAll() ([]ContentItem, error) {
	return repo.items, nil
}

func (repo *MockRepository) GetOne(id int) (ContentItem, error) {
	for _, n := range repo.items {
		if id == n.ID {
			return n, nil
		}
	}
	return ContentItem{}, nil
}

func (repo *MockRepository) Create(contentItem ContentItem) error {
	repo.items = append(repo.items, contentItem)
	return nil
}

func (repo *MockRepository) Update(contentItem ContentItem) error {
	return nil
}

func (repo *MockRepository) Delete(id int) error {
	return nil
}

func CreateMockRepository(db *sql.DB) *MockRepository {
	items := make([]ContentItem, 0)
	return &MockRepository{
		items: items,
	}
}
