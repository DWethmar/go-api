package contententry

import "testing"

func TestIntergrationX(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
