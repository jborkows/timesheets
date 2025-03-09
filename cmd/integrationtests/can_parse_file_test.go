package integrationtests

import (
	"log"
	"testing"
	"time"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCanAddParseFile(t *testing.T) {
	config := model.NewConfig([]string{"aaa", "bbb", "ccc"}, "Task-")
	date := time.Now()

	useWorkspace(config, func(service *model.Service) {
		parsed, lineErrors := service.ProcessForSave(`aaa 1.0 first
bbb 2.0 description
ccc 1.5 some other 
aaa 1h30m second
aaa 1h45m third`, date)
		log.Printf("Line errors: %v", lineErrors)
		assert.Equal(t, 0, len(lineErrors))
		assert.Equal(t, 5, len(parsed))
	})
}
