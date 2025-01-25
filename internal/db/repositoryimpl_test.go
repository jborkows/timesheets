package db_test

import (
	"testing"
	"time"

	. "github.com/jborkows/timesheets/internal/db"
	. "github.com/jborkows/timesheets/internal/model"
)

func TestShouldBeAbleToDisplayStatisticsForNoneData(t *testing.T) {
	t.Parallel()
	repository := Repository()
	var query Queryer = repository
	timesheet := TimesheetForDate(time.Now())

	statistics, err := query.Daily(timesheet)
	if err != nil {
		t.Errorf("Error getting daily statistics: %v", err)
	}
	if len(statistics) != 0 {
		t.Errorf("Expected none statistic, got %d", len(statistics))
	}
}
