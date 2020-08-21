package entries

import "errors"

var ErrNotFound = errors.New("Entry or entries not found")

type Repository interface {
	GetAll() ([]*Entry, error)
	GetOne(id ID) (*Entry, error)
	Add(contentItem Entry) error
	Update(contentItem Entry) error
	Delete(id ID) error
}
