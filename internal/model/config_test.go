package model_test

import (
	"strings"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var fakingToml = ` 
[categories]
regular=["categoryA", "categoryB"]
overtime=["overtimeA"]
[holidays]
repeatable=["11-11","05-01","01-01", "12-25"]
addhoc=["2021-02-01"]
`

func TestReadToml(t *testing.T) {
	t.Parallel()
	config, err := model.ReadConfig(strings.NewReader(fakingToml))
	assert.Nil(t, err)
	assert.Equal(t, 2, len(config.Categories.Regular))
	assert.Equal(t, 1, len(config.Categories.Overtime))
	assert.Equal(t, 4, len(config.Holidays.Repeatable))
	assert.Equal(t, 1, len(config.Holidays.AddHoc))
}
