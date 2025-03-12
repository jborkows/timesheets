package model_test

import (
	"log"
	"strings"
	"testing"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
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

func TestShouldTokenizeTextWithSingleSpace(t *testing.T) {
	t.Parallel()
	input := "aaa 1.0 first second"
	tokens := model.TokenizeFromIndex(input, len("aaa"))
	log.Printf("Tokens: %v", tokens)
	assert.Equal(t, 3, len(tokens))
	assert.Equal(t, strings.Index(input, "1.0"), tokens[0].Index)
	assert.Equal(t, "1.0", tokens[0].Word)
	assert.Equal(t, strings.Index(input, "first"), tokens[1].Index)
	assert.Equal(t, "first", tokens[1].Word)
	assert.Equal(t, strings.Index(input, "second"), tokens[2].Index)
	assert.Equal(t, "second", tokens[2].Word)
}

func TestShouldTokenizeTextWithMultipleSpace(t *testing.T) {
	t.Parallel()
	input := "aaa  1.0   first  second"
	tokens := model.TokenizeFromIndex(input, len("aaa"))
	log.Printf("Tokens: %v", tokens)
	assert.Equal(t, 3, len(tokens))
	assert.Equal(t, strings.Index(input, "1.0"), tokens[0].Index)
	assert.Equal(t, "1.0", tokens[0].Word)
	assert.Equal(t, strings.Index(input, "first"), tokens[1].Index)
	assert.Equal(t, "first", tokens[1].Word)
	assert.Equal(t, strings.Index(input, "second"), tokens[2].Index)
	assert.Equal(t, "second", tokens[2].Word)
}
