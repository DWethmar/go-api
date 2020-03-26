package contentitem

import (
	"errors"
)

var (
	MaxNameLength          = 50
	ErrAttrsRequired       = errors.New("Attrs is required.")
	ErrNameLength          = errors.New("Name exceeded max length of 50.")
	ErrUnsupportedAttrType = errors.New("Unsupported value.")
)

type validationResult struct {
	Errors struct {
		Name  string                       `json:"name,omitempty"`
		Attrs map[string]map[string]string `json:"attrs,omitempty"`
	} `json:"errors"`
}

func (v *validationResult) IsValid() bool {
	if v.Errors.Name != "" || len(v.Errors.Attrs) > 0 {
		return false
	}
	return true
}

func ValidateAddContentItem(addContentItem AddContentItem) validationResult {
	e := CreateValidationResult()

	nameErr := ValidateName(addContentItem.Name)
	if nameErr != nil {
		e.Errors.Name = nameErr.Error()
	}

	for locale, attrs := range addContentItem.Attrs {
		attrErRors := ValidateAttr(attrs)
		for attr, error := range attrErRors {
			if e.Errors.Attrs[locale] == nil {
				e.Errors.Attrs[locale] = map[string]string{}
			}
			e.Errors.Attrs[locale][attr] = error.Error()
		}
	}
	return e
}

func ValidateContentItem(contentItem ContentItem) validationResult {
	e := CreateValidationResult()

	nameErr := ValidateName(contentItem.Name)
	if nameErr != nil {
		e.Errors.Name = nameErr.Error()
	}

	for locale, attrs := range contentItem.Attrs {
		attrErRors := ValidateAttr(attrs)
		for attr, error := range attrErRors {
			e.Errors.Attrs[locale][attr] = error.Error()
		}
	}
	return e
}

func ValidateName(name string) error {
	if len(name) > MaxNameLength {
		return ErrNameLength
	}
	return nil
}

func ValidateAttr(attr map[string]interface{}) map[string]error {
	var e = make(map[string]error)

	validateAttr := func(attrs map[string]interface{}) []error {
		ce := make([]error, 0)
		for key, value := range attrs {
			switch value.(type) {
			case int:
			case []int:
			case string:
			case []string:
			case float64:
			case bool:
			default:
				e[key] = ErrUnsupportedAttrType
			}
		}
		return ce
	}

	validateAttr(attr)
	return e
}

func CreateValidationResult() validationResult {
	v := validationResult{}
	v.Errors.Attrs = make(map[string]map[string]string)
	return v
}
