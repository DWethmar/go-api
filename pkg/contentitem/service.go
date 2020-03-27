package contentitem

import (
	"time"
)

type Service interface {
	GetOne(ID) (*ContentItem, error)
	GetAll() ([]*ContentItem, error)
	Update(ContentItem) error
	Create(AddContentItem) (*ContentItem, error)
	Delete(ID) error
}

type service struct {
	repo Repository
}

func (s *service) GetOne(id ID) (*ContentItem, error) {
	item, err := s.repo.GetOne(id)
	return item, err
}

func (s *service) GetAll() ([]*ContentItem, error) {
	items, err := s.repo.GetAll()
	return items, err
}

func (s *service) Update(contentItem ContentItem) error {
	contentItem.UpdatedOn = time.Now()
	err := s.repo.Update(contentItem)
	return err
}

func (s *service) Create(addContentItem AddContentItem) (*ContentItem, error) {
	if addContentItem.Attrs == nil {
		addContentItem.Attrs = make(map[string]map[string]interface{})
	}
	var contentItem = ContentItem{
		ID:        createNewId(),
		Name:      addContentItem.Name,
		Attrs:     addContentItem.Attrs,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}
	err := s.repo.Add(contentItem)
	if err != nil {
		return nil, err
	}
	addedContentItem, err := s.GetOne(contentItem.ID)
	return addedContentItem, err
}

func (s *service) Delete(id ID) error {
	err := s.repo.Delete(id)
	return err
}

// CreateService creates a listing service with the necessary dependencies
func CreateService(r Repository) Service {
	return &service{r}
}
