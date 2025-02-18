package db_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

var testWorkHours = uint8(4)

// func TestShouldBeAbleToPresentWeeklyStatisticsPending(t *testing.T) {
// 	useDb(t, func(saver model.Saver, query model.Queryer) {
// 		monday := time.Date(2025, time.February, 17, 0, 0, 0, 0, time.UTC)
// 		saveWorkTimeSheet(t, monday, saver)
// 		saveWorkTimeSheet(t, time.Date(2025, time.February, 18, 0, 0, 0, 0, time.UTC), saver)
// 		saveWorkTimeSheet(t, time.Date(2025, time.February, 19, 0, 0, 0, 0, time.UTC), saver)
// 		sunday := time.Date(2025, time.February, 16, 0, 0, 0, 0, time.UTC)
// 		saveWorkTimeSheet(t, sunday, saver)
//
// 		statistics, err := query.Weekly(model.TimesheetForDate(sunday))
// 		if err != nil {
// 			t.Errorf("Error getting daily statistics: %v", err)
// 		}
// 		assert.Equal(t, 1, len(statistics), "Expected 1 entry statistic, got %d", len(statistics))
// 		stat := statistics[0]
// 		assert.Equal(t, testWorkHours, stat.Weekly.Hours, "Expected 4 hours, got %d", stat.Weekly.Hours)
//
// 		statistics, err = query.Weekly(model.TimesheetForDate(monday))
// 		if err != nil {
// 			t.Errorf("Error getting daily statistics: %v", err)
// 		}
// 		stat = statistics[0]
// 		assert.Equal(t, testWorkHours*3, stat.Weekly.Hours, "Expected 4 hours, got %d", stat.Weekly.Hours)
// 	})
// }

func TestShouldBeAbleToPresentMonthlyStatisticsPending(t *testing.T) {
	useDb(t, func(saver model.Saver, query model.Queryer) {
		sunday := time.Date(2025, time.February, 16, 0, 0, 0, 0, time.UTC)
		saveWorkTimeSheet(t, sunday, saver)
		monday := time.Date(2025, time.February, 17, 0, 0, 0, 0, time.UTC)
		saveWorkTimeSheet(t, monday, saver)
		saveWorkTimeSheet(t, time.Date(2025, time.February, 18, 0, 0, 0, 0, time.UTC), saver)
		saveWorkTimeSheet(t, time.Date(2025, time.February, 19, 0, 0, 0, 0, time.UTC), saver)

		statistics, err := query.Monthly(model.TimesheetForDate(monday))
		if err != nil {
			t.Errorf("Error getting daily statistics: %v", err)
		}
		assert.Equal(t, 1, len(statistics), "Expected 1 entry statistic, got %d", len(statistics))
		stat := statistics[0]
		assert.Equal(t, testWorkHours*uint8(4), stat.Monthly.Hours, "Expected 4 hours, got %d", stat.Monthly.Hours)

	})
}

func saveWorkTimeSheet(t *testing.T, time time.Time, saver model.Saver) *model.Timesheet {
	timesheet := model.TimesheetForDate(time)

	task := "work"
	entry := model.TimesheetEntry{Hours: testWorkHours, Minutes: 0, Category: "work", Comment: "work", Task: &task}
	if error := timesheet.AddEntry(entry); error != nil {
		assert.FailNow(t, fmt.Sprintf("Error adding entry: %v", error))
		log.Printf("Added entry failed: %v %v", entry, error)
	}
	saveError := saver.Save(timesheet)
	if saveError != nil {
		assert.FailNow(t, fmt.Sprintf("Error saving time sheet: %v", saveError))
		log.Printf("Added entry failed (failed): %v %v", entry, saveError)
	}
	log.Printf("Saved time sheet: %v", time)
	return timesheet
}
