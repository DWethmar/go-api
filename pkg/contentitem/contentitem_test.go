package contentitem

import (
	"reflect"
	"testing"
)

func TestAttrsValue(t *testing.T) {
	a := Attrs{
		"attrA": 1,
	}

	value, err := a.Value()
	if err != nil {
		t.Errorf("json encoding failed.")
	}

	valueBytes, ok := value.([]byte)
	if !ok {
		t.Errorf("type of Value() assertion to []byte failed.")
	}

	expected := "{\"attrA\":1}"
	if expected != string(valueBytes) {
		t.Errorf("Encountered %v expected %v", string(valueBytes), expected)
	}
}

func TestAttrsScan(t *testing.T) {
	a := make(Attrs)
	err := a.Scan([]byte("{\"attrA\": 1}"))
	if err != nil {
		t.Errorf("A error occurred while performing a scan. %v", err)
	}
	if a["attrA"] != float64(1) {
		t.Errorf("Expected attr1 to be 1 but got %v of type %v", a["attr1"], reflect.TypeOf(a["attr1"]))
	}
}
