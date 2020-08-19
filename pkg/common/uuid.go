package common

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// UUID type
type UUID = uuid.UUID

// CreateNewUUID creates new UUID.
func CreateNewUUID() UUID {
	return uuid.New()
}

// ParseUUID parses UUID from string.
func ParseUUID(val string) (UUID, error) {
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

type uuidCtxKeyType string

const uuidCtxKey uuidCtxKeyType = "uuid"

// WithUUID puts the request ID into the current context.
func WithUUID(ctx context.Context, id UUID) context.Context {
	return context.WithValue(ctx, uuidCtxKey, id)
}

// UUIDFromContext returns the request ID from the context.
// A zero ID is returned if there are no idenfiers in the
// current context.
func UUIDFromContext(ctx context.Context) (UUID, error) {
	v, ok := ctx.Value(uuidCtxKey).(UUID)
	if !ok {
		return uuid.Nil, errors.New("Could not receive ID")
	}
	return v, nil
}
