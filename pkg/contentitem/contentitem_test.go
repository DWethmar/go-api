package contentitem

import (
	"reflect"
	"testing"
)

func TestUnitAttrsValue(t *testing.T) {
	a := Attrs{
		"nl": {
			"attrA": 1,
		},
	}

	value, err := a.Value()
	if err != nil {
		t.Errorf("json encoding failed.")
	}

	valueBytes, ok := value.([]byte)
	if !ok {
		t.Errorf("type of Value() assertion to []byte failed.")
	}

	expected := "{\"nl\":{\"attrA\":1}}"
	if expected != string(valueBytes) {
		t.Errorf("Encountered %v expected %v", string(valueBytes), expected)
	}
}

func TestUnitAttrsScan(t *testing.T) {
	a := make(Attrs)
	err := a.Scan([]byte("{\"nl\":{\"attrA\":1,\"attrB\":[\"a\",\"b\"]}}"))
	if err != nil {
		t.Errorf("A error occurred while performing a scan. %v", err)
	}
	if a["nl"] != nil && a["nl"]["attrA"] != float64(1) {
		t.Errorf("Expected attr1 to be 1 but got %v of type %v", a["attr1"], reflect.TypeOf(a["attr1"]))
	}
}
