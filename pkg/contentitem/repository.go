package contentitem

import "errors"

var ErrNotFound = errors.New("ContentItem(s) not found")

type Repository interface {
	GetAll() ([]ContentItem, error)
	GetOne(id int) (ContentItem, error)
	Create(contentItem ContentItem) error
	Update(contentItem ContentItem) error
	Delete(id int) error
}
