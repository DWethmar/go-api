package contenttype

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
)

// Service entry service
type Service interface {
	Get(common.ID) (*ContentType, error)
	List() ([]*ContentType, error)
	Update(*ContentType) error
	Create(*ContentType) (common.ID, error)
	Delete(common.ID) error
}

type service struct {
	repo Repository
}

func (s *service) Get(id common.ID) (*ContentType, error) {
	item, err := s.repo.Get(id)
	return item, err
}

func (s *service) List() ([]*ContentType, error) {
	items, err := s.repo.List()
	return items, err
}

func (s *service) Update(contentItem *ContentType) error {
	contentItem.UpdatedOn = time.Now()
	return s.repo.Update(contentItem)
}

func (s *service) Create(entry *ContentType) (common.ID, error) {
	if entry.Fields == nil {
		entry.Fields = make([]*Field, 0)
	}

	var contentItem = &ContentType{
		ID:        common.NewID(),
		Name:      entry.Name,
		Fields:    entry.Fields,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}

	return s.repo.Create(contentItem)
}

func (s *service) Delete(id common.ID) error {
	err := s.repo.Delete(id)
	return err
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}
