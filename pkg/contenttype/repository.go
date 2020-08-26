package contenttype

import (
	"errors"

	"github.com/dwethmar/go-api/pkg/common"
)

// ErrNotFound entries not found error.
var ErrNotFound = errors.New("Entry or entries not found")

// Repository entry repository
type Repository interface {
	List() ([]*ContentType, error)
	Get(ID common.ID) (*ContentType, error)
	Create(c *ContentType) (common.ID, error)
	Update(c *ContentType) error
	Delete(ID common.ID) error
}
