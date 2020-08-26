package content

import (
	"errors"

	"github.com/dwethmar/go-api/pkg/common"
)

// ErrNotFound entries not found error.
var ErrNotFound = errors.New("Entry or entries not found")

// Repository entry repository
type Repository interface {
	List() ([]*Content, error)
	Get(ID common.ID) (*Content, error)
	Create(c *Content) (common.ID, error)
	Update(c *Content) error
	Delete(ID common.ID) error
}
