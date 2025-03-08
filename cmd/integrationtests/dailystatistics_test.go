package integrationtests

import (
	"log"
	"os"
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

func TestShowDailyStatistics(t *testing.T) {
	config := model.NewConfig([]string{"aaa", "bbb", "ccc"}, "Task-")
	date, err := time.Parse("2006-01-02", "2022-02-02")
	if err != nil {
		log.Fatalf("Failed to parse date: %v", err)
	}

	useWorkspace(config, func(service *model.Service) {
		_, _ = service.ProcessForSave(`aaa 1.0 first
bbb 2.0 description
ccc 1.5 Task-123 some other 
aaa 1h30m  second
aaa 1h45m third`, date)
		reportFile, err := service.ShowDailyStatistics(date)
		if err != nil {
			log.Fatalf("Failed to generate report: %v", err)
		}
		content, err := os.ReadFile(string(reportFile))
		if err != nil {
			log.Fatalf("Failed to read report: %v", err)
		}
		log.Printf("Report content: %s", content)
		desiredContent := `Daily statistics for 2022-02-02
aaa 3.75
1.0 first
1.5 second
1.75 third
bbb 2.0
2.0 description
ccc 1.5
1.5 Task-123 some other`
		assert.Equal(t, desiredContent, string(content))
	})
}
