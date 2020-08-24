package contentmodels

import (
	"errors"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/models"
)

// ErrNotFound entries not found error.
var ErrNotFound = errors.New("Entry or entries not found")

// Repository entry repository
type Repository interface {
	GetAll() ([]*models.ContentModel, error)
	GetOne(id common.UUID) (*models.ContentModel, error)
	Add(contentItem models.ContentModel) error
	Update(contentItem models.ContentModel) error
	Delete(id common.UUID) error
}
