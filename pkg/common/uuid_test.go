package common

import (
	"net/http/httptest"
	"testing"
)

func TestParseUUID(t *testing.T) {
	_, err := ParseUUID("1f7e2926-31a0-4a1a-bc01-88811ff60165")

	if err != nil {
		t.Error("Unexpected error.", err)
	}
}

func TestInvalidParseUUID(t *testing.T) {
	_, err := ParseUUID("abc")

	if err == nil {
		t.Error("Expected error.")
	}
}

func TestWithWithUUID(t *testing.T) {
	ID := CreateNewUUID()

	req := httptest.NewRequest("DELETE", "/", nil)
	ctx := req.Context()
	ctx = WithUUID(ctx, ID)

	IDFromContext, err := UUIDFromContext(ctx)

	if err != nil {
		t.Error("Unexpected error.", err)
	}

	if ID != IDFromContext {
		t.Errorf("IDs are not equal: %v %v", ID, IDFromContext)
	}
}
