package model_test

import (
	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractDateFromFileName(t *testing.T) {
	t.Parallel()
	params := model.DateFromFileNameParams{

		URI:         "file:///mnt/ramdisk/test_project/2025/02/01.tsf",
		ProjectRoot: "/mnt/ramdisk/test_project",
	}
	date, err := model.DateFromFile(params)
	assert.Nil(t, err)
	assert.Equal(t, "2025-02-01", date.Format("2006-01-02"))
}
