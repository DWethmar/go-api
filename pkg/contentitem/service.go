package contentitem

import (
	"time"
)

type Service interface {
	GetOne(int) (ContentItem, error)
	GetAll() ([]ContentItem, error)
	Update(ContentItem) error
	Create(NewContentItem) (ContentItem, error)
	Delete(int) error
}

type service struct {
	r Repository
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetOne(id int) (ContentItem, error) {
	item, err := s.r.GetOne(id)
	return item, err
}

func (s *service) GetAll() ([]ContentItem, error) {
	item, err := s.r.GetAll()
	return item, err
}

func (s *service) Update(contentItem ContentItem) error {
	contentItem.UpdatedOn = time.Now()
	err := s.r.Update(contentItem)
	return err
}

func (s *service) Create(newContentItem NewContentItem) (ContentItem, error) {
	var contentItem = ContentItem{
		Name:      newContentItem.Name,
		Attrs:     newContentItem.Attrs,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}
	id, err := s.r.Create(contentItem)
	contentItem, err = s.GetOne(id)
	return contentItem, err
}

func (s *service) Delete(id int) error {
	err := s.r.Delete(id)
	return err
}
