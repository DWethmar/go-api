package output

import (
	"testing"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/content"
	"gotest.tools/v3/assert"
)

// ContentOut maps to output model.
func TestContentOut(t *testing.T) {
	now := common.Now()

	c := ContentOut(&content.Content{
		ID:   common.NewID(),
		Name: "Test1",
		Fields: content.FieldTranslations{
			"nl": content.Fields{
				"attrA": 1,
			},
		},
		CreatedOn: now,
		UpdatedOn: now,
	})

	assert.Equal(t, c.Name, "Test1", "Name not the same")
	assert.Equal(t, c.CreatedOn, now, "CreatedOn not the same")
	assert.Equal(t, c.UpdatedOn, now, "UpdatedOn not the same")

	for locale, fields := range c.Fields {
		assert.Equal(t, locale, "nl", "Locale not the same")

		for key, value := range fields {
			assert.Equal(t, key, "attrA", "Field key not the same")
			assert.Equal(t, value, 1, "Field value not the same")
		}
	}
}
