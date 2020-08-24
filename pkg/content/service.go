package content

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/models"
)

// Service entry service
type Service interface {
	GetOne(common.UUID) (*models.Content, error)
	GetAll() ([]*models.Content, error)
	Update(models.Content) error
	Create(models.AddContent) (*models.Content, error)
	Delete(common.UUID) error
}

type service struct {
	repo Repository
}

func (s *service) GetOne(id common.UUID) (*models.Content, error) {
	item, err := s.repo.GetOne(id)
	return item, err
}

func (s *service) GetAll() ([]*models.Content, error) {
	items, err := s.repo.GetAll()
	return items, err
}

func (s *service) Update(contentItem models.Content) error {
	contentItem.UpdatedOn = time.Now()
	err := s.repo.Update(contentItem)
	return err
}

func (s *service) Create(entry models.AddContent) (*models.Content, error) {
	if entry.Fields == nil {
		entry.Fields = make(models.FieldTranslations)
	}

	var contentItem = models.Content{
		ID:        common.CreateNewUUID(),
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

func (s *service) Delete(id common.UUID) error {
	err := s.repo.Delete(id)
	return err
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}
