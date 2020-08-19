package contententry

import (
	"time"
)

type Service interface {
	GetOne(ID) (*Entry, error)
	GetAll() ([]*Entry, error)
	Update(Entry) error
	Create(AddEntry) (*Entry, error)
	Delete(ID) error
}

type service struct {
	repo Repository
}

func (s *service) GetOne(id ID) (*Entry, error) {
	item, err := s.repo.GetOne(id)
	return item, err
}

func (s *service) GetAll() ([]*Entry, error) {
	items, err := s.repo.GetAll()
	return items, err
}

func (s *service) Update(contentItem Entry) error {
	contentItem.UpdatedOn = time.Now()
	err := s.repo.Update(contentItem)
	return err
}

func (s *service) Create(entry AddEntry) (*Entry, error) {
	if entry.Fields == nil {
		entry.Fields = make(FieldTranslations)
	}
	var contentItem = Entry{
		ID:        createNewID(),
		Name:      entry.Name,
		Fields:    entry.Fields,
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
