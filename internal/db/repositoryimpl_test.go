package db_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	. "github.com/jborkows/timesheets/internal/db"
	. "github.com/jborkows/timesheets/internal/model"
)

func TestShouldBeAbleToDisplayStatisticsForNoneData(t *testing.T) {
	useDb(t, func(saver Saver, query Queryer) {
		timesheet := TimesheetForDate(time.Now())

		statistics, err := query.Daily(timesheet)
		if err != nil {
			t.Errorf("Error getting daily statistics: %v", err)
		}
		if len(statistics) != 0 {
			t.Errorf("Expected none statistic, got %d", len(statistics))
		}
	})
}

func useDb(t *testing.T, test func(saver Saver, querier Queryer)) {
	t.Parallel()
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	defer func() {
		tempFile.Close()

		err = os.Remove(tempFile.Name())
		if err != nil {
			fmt.Println("Error removing temporary file:", err)
		} else {
			fmt.Println("Temporary file removed.")
		}
	}()

	_, err = NewDatabase(tempFile.Name())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	repository := Repository()
	test(repository, repository)
}
