package contenttype

import (
	"errors"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/models"
)

// ErrNotFound entries not found error.
var ErrNotFound = errors.New("Entry or entries not found")

// Repository entry repository
type Repository interface {
	GetAll() ([]*models.ContentType, error)
	GetOne(ID common.UUID) (*models.ContentType, error)
	Add(contentItem models.ContentType) error
	Update(contentItem models.ContentType) error
	Delete(ID common.UUID) error
}
