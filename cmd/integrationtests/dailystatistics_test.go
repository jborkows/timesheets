package integrationtests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestShowDailyStatistics(t *testing.T) {
	config := model.NewConfig([]string{"aaa", "bbb", "ccc"}, "Task-")
	date, err := time.Parse("2006-01-02", "2025-03-06")
	if err != nil {
		log.Fatalf("Failed to parse date: %v", err)
	}

	useWorkspace(config, func(service *model.Service) {
		_, _ = service.ProcessForSave(`aaa 1.0 first
bbb 2.0 description
ccc 1.5 Task-123 some other 
aaa 1h30m second
aaa 1h45m third`, date)
		_, _ = service.ProcessForSave(`aaa 1.0 first`, date.AddDate(0, 0, 1))
		_, _ = service.ProcessForSave(`aaa 1.0 first`, date.AddDate(0, 0, 8))
		reportFile, err := service.ShowDailyStatistics(date)
		if err != nil {
			log.Fatalf("Failed to generate report: %v", err)
		}
		content, err := os.ReadFile(string(reportFile))
		if err != nil {
			log.Fatalf("Failed to read report: %v", err)
		}
		log.Printf("Report content: %s", content)
		desiredContent := `For 2025-03-06
Daily statistics (7:45)

aaa 4:15
1.0 first
1.5 second
1.75 third

bbb 2:00
2.0 description

ccc 1:30
1.5 Task-123 some other

####################

Weekly statistics (8:45)
aaa 5.25
bbb 2.0
ccc 1.5

####################

Monthly statistics (9:45/24)
aaa 6.25
bbb 2.0
ccc 1.5
`
		assert.Equal(t, desiredContent, string(content))
	})
}

func TestReportFileShouldBeReused(t *testing.T) {
	config := model.NewConfig([]string{"aaa"}, "Task-")
	date, err := time.Parse("2006-01-02", "2025-03-06")
	if err != nil {
		log.Fatalf("Failed to parse date: %v", err)
	}

	useWorkspace(config, func(service *model.Service) {
		_, _ = service.ProcessForSave(`aaa 1.0 first`, date)
		reportFile, err := service.ShowDailyStatistics(date)
		if err != nil {
			log.Fatalf("Failed to generate report: %v", err)
		}
		_, _ = service.ProcessForSave(`aaa 1.0 first`, date.AddDate(0, 0, 1))
		_, err = service.ShowDailyStatistics(date)
		if err != nil {
			log.Fatalf("Failed to generate report: %v", err)
		}
		content, err := os.ReadFile(string(reportFile))
		if err != nil {
			log.Fatalf("Failed to read report: %v", err)
		}
		log.Printf("Report content: %s", content)
		desiredContent := `For 2025-03-06
Daily statistics (1:00)

aaa 1:00
1.0 first

####################

Weekly statistics (2:00)
aaa 2.0

####################

Monthly statistics (2:00/16)
aaa 2.0
`
		assert.Equal(t, desiredContent, string(content))
	})
}
