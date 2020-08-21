package entries

import (
	"time"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/models"
)

// Service entry service
type Service interface {
	GetOne(common.UUID) (*models.Entry, error)
	GetAll() ([]*models.Entry, error)
	Update(models.Entry) error
	Create(models.AddEntry) (*models.Entry, error)
	Delete(common.UUID) error
}

type service struct {
	repo Repository
}

func (s *service) GetOne(id common.UUID) (*models.Entry, error) {
	item, err := s.repo.GetOne(id)
	return item, err
}

func (s *service) GetAll() ([]*models.Entry, error) {
	items, err := s.repo.GetAll()
	return items, err
}

func (s *service) Update(contentItem models.Entry) error {
	contentItem.UpdatedOn = time.Now()
	err := s.repo.Update(contentItem)
	return err
}

func (s *service) Create(entry models.AddEntry) (*models.Entry, error) {
	if entry.Fields == nil {
		entry.Fields = make(models.FieldTranslations)
	}

	var contentItem = models.Entry{
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

// CreateService creates a listing service with the necessary dependencies
func CreateService(r Repository) Service {
	return &service{r}
}
