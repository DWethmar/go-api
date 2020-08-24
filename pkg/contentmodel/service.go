package contentmodels

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/models"
)

// Service entry service
type Service interface {
	GetOne(common.UUID) (*models.ContentModel, error)
	GetAll() ([]*models.ContentModel, error)
	Update(models.ContentModel) error
	Create(models.AddContentModel) (*models.ContentModel, error)
	Delete(common.UUID) error
}

type service struct {
	repo Repository
}

func (s *service) GetOne(id common.UUID) (*models.ContentModel, error) {
	item, err := s.repo.GetOne(id)
	return item, err
}

func (s *service) GetAll() ([]*models.ContentModel, error) {
	items, err := s.repo.GetAll()
	return items, err
}

func (s *service) Update(contentItem models.ContentModel) error {
	contentItem.UpdatedOn = time.Now()
	err := s.repo.Update(contentItem)
	return err
}

func (s *service) Create(entry models.AddContentModel) (*models.ContentModel, error) {
	if entry.Fields == nil {
		entry.Fields = make([]*models.ContentModelField, 0)
	}

	var contentItem = models.ContentModel{
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
