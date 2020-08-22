package content

import (
	"errors"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/models"
)

// ErrNotFound entries not found error.
var ErrNotFound = errors.New("Entry or entries not found")

// Repository entry repository
type Repository interface {
	GetAll() ([]*models.Content, error)
	GetOne(id common.UUID) (*models.Content, error)
	Add(contentItem models.Content) error
	Update(contentItem models.Content) error
	Delete(id common.UUID) error
}
