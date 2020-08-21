package models

import (
	"reflect"
	"testing"
)

func TestUnitFieldValue(t *testing.T) {
	a := FieldTranslations{
		"nl": Fields{
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

func TestUnitFieldScan(t *testing.T) {
	a := make(FieldTranslations)
	err := a.Scan([]byte("{\"nl\":{\"attrA\":1,\"attrB\":[\"a\",\"b\"]}}"))
	if err != nil {
		t.Errorf("A error occurred while performing a scan. %v", err)
	}
	if a["nl"] != nil && a["nl"]["attrA"] != float64(1) {
		t.Errorf("Expected attr1 to be 1 but got %v of type %v", a["attr1"], reflect.TypeOf(a["attr1"]))
	}
}


func TestUnitFieldValues(t *testing.T) {
	c := FieldTranslations{
		"nl": {
			"attr1": 1,
			"attr2": "attribute string value",
			"attr3": []int{1, 2, 3},
			"attr4": []string{"one", "two"},
			"attr5": float64(3) / float64(10),
		},
	}

	if errors := validateFields(c["nl"]); len(errors) > 0 {
		for key, err := range errors {
			fmt.Printf("Validation error on attr %v %+v\n", key, err)
		}
		t.Errorf("Encountered %v validation errors", len(errors))
	}
}

func TestUnitNameValidation(t *testing.T) {
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

func TestUnitInvalidFieldValues(t *testing.T) {
	var names []interface{}
	names = append(names, "test")
	names = append(names, make(map[string]string))

	c := FieldTranslations{
		"nl": {
			"attrX": nil,
			"attrY": names,
		},
	}

	if errors := validateFields(c["nl"]); len(errors) == 2 {
		for attr, err := range errors {
			if err != ErrUnsupportedFieldValue {
				if err != ErrUnsupportedFieldSliceValue {
					t.Errorf("Validation returned unexpected error on attr %v with error %v:", attr, err)
				}
			}
		}
	} else {
		t.Errorf("Expected %d errors but received %d errors.", 2, len(errors))

		for attr, err := range errors {
			fmt.Printf("%v: %v\n", attr, err)
		}
	}
}
