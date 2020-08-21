package entries

import (
	"errors"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/models"
)

var ErrNotFound = errors.New("Entry or entries not found")

type Repository interface {
	GetAll() ([]*models.Entry, error)
	GetOne(id common.UUID) (*models.Entry, error)
	Add(contentItem models.Entry) error
	Update(contentItem models.Entry) error
	Delete(id common.UUID) error
}
