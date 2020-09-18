package common

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUUID(t *testing.T) {
	_, err := StringToID("1f7e2926-31a0-4a1a-bc01-88811ff60165")
	assert.Nil(t, err)
}

func TestInvalidParseUUID(t *testing.T) {
	_, err := StringToID("abc")
	assert.NotNil(t, err)
}

func TestWithWithUUID(t *testing.T) {
	ID := NewID()

	req := httptest.NewRequest("DELETE", "/", nil)
	ctx := req.Context()
	ctx = WithID(ctx, ID)

	IDFromContext, err := UUIDFromContext(ctx)

	assert.Nil(t, err)

	if ID != IDFromContext {
		t.Errorf("IDs are not equal: %v %v", ID, IDFromContext)
	}
}
