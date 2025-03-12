package lspserver_test

import (
	"testing"

	"github.com/jborkows/timesheets/internal/lspserver"
	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestTranslatorSameLine(t *testing.T) {
	t.Parallel()
	tokens := []model.TokenReady{
		{
			Line:   1,
			Column: 0,
			Length: 3,
			Type:   model.ClassType,
		},
		{
			Line:   1,
			Column: 5,
			Length: 6,
			Type:   model.PropertyType,
		},
		{
			Line:   1,
			Column: 8,
			Length: 9,
			Type:   model.StringType,
		},
	}
	tokensToSend := lspserver.TranslateSemanticTokens(tokens)
	assert.Equal(t, []int{1, 0, 3, int(model.ClassType), 0, 0, 5, 6, int(model.PropertyType), 0, 0, 3, 9, int(model.StringType), 0}, tokensToSend)
}

func TestTranslatorDifferentLines(t *testing.T) {
	t.Parallel()
	tokens := []model.TokenReady{
		{
			Line:   1,
			Column: 0,
			Length: 3,
			Type:   model.ClassType,
		},
		{
			Line:   2,
			Column: 5,
			Length: 6,
			Type:   model.PropertyType,
		},
	}
	tokensToSend := lspserver.TranslateSemanticTokens(tokens)
	assert.Equal(t, []int{1, 0, 3, int(model.ClassType), 0, 1, 5, 6, int(model.PropertyType), 0}, tokensToSend)
}
