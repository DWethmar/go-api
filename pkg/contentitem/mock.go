package contentitem

type MockRepository struct {
	items  []ContentItem
	lastId int
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

func (repo *MockRepository) Create(contentItem ContentItem) (int, error) {
	c := ContentItem{
		ID:        repo.lastId,
		Name:      contentItem.Name,
		Attrs:     contentItem.Attrs,
		CreatedOn: contentItem.CreatedOn,
		UpdatedOn: contentItem.UpdatedOn,
	}
	repo.items = append(repo.items, c)
	repo.lastId++
	return c.ID, nil
}

func (repo *MockRepository) Update(contentItem ContentItem) error {
	for i, n := range repo.items {
		if contentItem.ID == n.ID {
			repo.items = append(repo.items[:i], append([]ContentItem{contentItem}, repo.items[i:]...)...)
			return nil
		}
	}
	return nil
}

func (repo *MockRepository) Delete(id int) error {
	for i, n := range repo.items {
		if id == n.ID {
			repo.items = append(repo.items[:i], repo.items[i:]...)
			return nil
		}
	}
	return nil
}

func CreateMockRepository(contentItems ...ContentItem) *MockRepository {
	items := make([]ContentItem, 0)
	if len(contentItems) > 0 {
		items = append(items, contentItems...)
	}
	return &MockRepository{
		lastId: 1,
		items:  items,
	}
}
