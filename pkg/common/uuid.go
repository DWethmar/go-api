package common

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ID type
type ID = uuid.UUID

// NewID creates new UUID
func NewID() ID {
	return uuid.New()
}

// StringToID parses UUID from string.
func StringToID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

type uuidCtxKeyType string

const uuidCtxKey uuidCtxKeyType = "uuid"

// WithID puts the request ID into the current context.
func WithID(ctx context.Context, id ID) context.Context {
	return context.WithValue(ctx, uuidCtxKey, id)
}

// UUIDFromContext returns the request ID from the context.
// A zero ID is returned if there are no idenfiers in the
// current context.
func UUIDFromContext(ctx context.Context) (ID, error) {
	v, ok := ctx.Value(uuidCtxKey).(ID)
	if !ok {
		return uuid.Nil, errors.New("Could not receive ID")
	}
	return v, nil
}
