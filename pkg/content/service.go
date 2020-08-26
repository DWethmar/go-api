package content

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
)

// Service entry service
type Service interface {
	Get(common.ID) (*Content, error)
	List() ([]*Content, error)
	Update(*Content) error
	Create(*Content) (common.ID, error)
	Delete(common.ID) error
}

type service struct {
	repo Repository
}

func (s *service) Get(id common.ID) (*Content, error) {
	return s.repo.Get(id)
}

func (s *service) List() ([]*Content, error) {
	return s.repo.List()
}

func (s *service) Update(entry *Content) error {
	entry.UpdatedOn = time.Now()
	return s.repo.Update(entry)
}

func (s *service) Create(entry *Content) (common.ID, error) {
	if entry.Fields == nil {
		entry.Fields = make(FieldTranslations)
	}

	var contentItem = &Content{
		ID:        common.NewID(),
		Name:      entry.Name,
		Fields:    entry.Fields,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}

	return s.repo.Create(contentItem)
}

func (s *service) Delete(id common.ID) error {
	return s.repo.Delete(id)
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}
