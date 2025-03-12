package integrationtests

import (
	"testing"
	"time"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestShouldBeAbleToGetTokens(t *testing.T) {
	config := model.NewConfig([]string{"aaa", "bbb", "ccc"}, "Task-")
	date := time.Now()

	useWorkspace(config, func(service *model.Service) {
		tokens := service.SemanaticTokenFrom([]string{"aaa 1.0 first", "bbb 2.0 description haha"}, date)

		assert.Equal(t, 7, len(tokens))
		//aaa
		assert.Equal(t, 0, tokens[0].Line)
		assert.Equal(t, 0, tokens[0].Column)
		assert.Equal(t, model.ClassType, tokens[0].Type)
		assert.Equal(t, len("aaa"), tokens[0].Length)
		//1.0
		assert.Equal(t, 0, tokens[1].Line)
		assert.Equal(t, len("aaa "), tokens[1].Column)
		assert.Equal(t, model.PropertyType, tokens[1].Type)
		assert.Equal(t, len("1.0"), tokens[1].Length)
		//first
		assert.Equal(t, 0, tokens[2].Line)
		assert.Equal(t, len("aaa 1.0 "), tokens[2].Column)
		assert.Equal(t, model.StringType, tokens[2].Type)
		assert.Equal(t, len("first"), tokens[2].Length)
		//bbb
		assert.Equal(t, 1, tokens[3].Line)
		assert.Equal(t, 0, tokens[3].Column)
		assert.Equal(t, model.ClassType, tokens[3].Type)
		assert.Equal(t, len("bbb"), tokens[3].Length)
		//2.0
		assert.Equal(t, 1, tokens[4].Line)
		assert.Equal(t, len("bbb "), tokens[4].Column)
		assert.Equal(t, model.PropertyType, tokens[4].Type)
		assert.Equal(t, len("2.0"), tokens[4].Length)
		//description
		assert.Equal(t, 1, tokens[5].Line)
		assert.Equal(t, len("bbb 2.0 "), tokens[5].Column)
		assert.Equal(t, model.StringType, tokens[5].Type)
		assert.Equal(t, len("description"), tokens[5].Length)
		//haha
		assert.Equal(t, 1, tokens[6].Line)
		assert.Equal(t, len("bbb 2.0 description "), tokens[6].Column)
		assert.Equal(t, model.StringType, tokens[6].Type)
		assert.Equal(t, len("haha"), tokens[6].Length)

	})
}
