package models

import (
	"testing"
)

func TestValidateName(t *testing.T) {
	c := Entry{
		Name: "This name is to loooooooOooooOOo0000000000000000000000000oooong",
	}

	if err := validateName(c.Name); err == nil {
		t.Errorf("Expected error")
	} else {
		if err != ErrNameLength {
			t.Error("Unexpected error.", err)
		}
	}
}
