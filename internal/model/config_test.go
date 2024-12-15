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

func TestShouldBeAddHocHoliday(t *testing.T) {
	t.Parallel()
	config, _ := model.ReadConfig(strings.NewReader(fakingToml))
	date := model.DateInfo{Value: "2021-02-01"}
	assert.True(t, config.IsHoliday(&date))
}

func TestShouldBeFindRegularHoliday(t *testing.T) {
	t.Parallel()
	config, _ := model.ReadConfig(strings.NewReader(fakingToml))
	date := model.DateInfo{Value: "2021-11-11"}
	assert.True(t, config.IsHoliday(&date))
	date = model.DateInfo{Value: "2022-11-11"}
	assert.True(t, config.IsHoliday(&date))
	date = model.DateInfo{Value: "2024-12-25"}
	assert.True(t, config.IsHoliday(&date))
}
func TestShouldBeNotFindHoliday(t *testing.T) {
	t.Parallel()
	config, _ := model.ReadConfig(strings.NewReader(fakingToml))
	date := model.DateInfo{Value: "2021-11-12"}
	assert.False(t, config.IsHoliday(&date))
}
