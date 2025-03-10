package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestShouldBeAbleToDisplayStatisticsForNoneData(t *testing.T) {
	useDb(t, func(saver model.Saver, query model.Queryer) {
		timesheet := model.TimesheetForDate(time.Now())

		statistics, err := query.Daily(context.Background(), timesheet)
		if err != nil {
			t.Errorf("Error getting daily statistics: %v", err)
		}
		if len(statistics) != 0 {
			t.Errorf("Expected none statistic, got %d", len(statistics))
		}
	})
}

func TestShouldBeAbleToPresentDailyStatisticsPendingThenSave(t *testing.T) {
	useDb(t, func(saver model.Saver, query model.Queryer) {
		timesheet := model.TimesheetForDate(time.Now())

		task := "work"
		entry := model.TimesheetEntry{Hours: 4, Minutes: 0, Category: "work", Comment: "work", Task: &task}
		if error := timesheet.AddEntry(entry); error != nil {
			t.Errorf("Error adding entry: %v", error)
		}
		saveError := saver.PendingSave(context.Background(), timesheet)
		if saveError != nil {
			t.Errorf("Error pending saving time sheet: %v", saveError)
		}
		statistics, err := query.Daily(context.Background(), timesheet)
		if err != nil {
			t.Errorf("Error getting daily statistics: %v", err)
		}
		assert.Equal(t, 1, len(statistics), "Expected 1 entry statistic, got %d", len(statistics))
		stat := statistics[0]
		assert.Equal(t, uint8(4), stat.Dirty.Hours, "Expected 4 hours, got %d", stat.Dirty.Hours)
		assert.Equal(t, uint8(0), stat.Daily.Hours, "Expected 0 hours, got %d", stat.Daily.Hours)

		saveError = saver.Save(context.Background(), timesheet)
		if saveError != nil {
			t.Errorf("Error saving time sheet: %v", saveError)
		}

		statistics, err = query.Daily(context.Background(), timesheet)
		if err != nil {
			t.Errorf("Error getting daily statistics: %v", err)
		}
		assert.Equal(t, 1, len(statistics), "Expected 1 entry statistic, got %d", len(statistics))
		stat = statistics[0]
		assert.Equal(t, uint8(4), stat.Dirty.Hours, "Expected 4 hours, got %d", stat.Dirty.Hours)
		assert.Equal(t, uint8(4), stat.Daily.Hours, "Expected 4 hours, got %d", stat.Daily.Hours)
	})
}

func TestShouldBeAbleToPresentDailyStatisticsSave(t *testing.T) {
	useDb(t, func(saver model.Saver, query model.Queryer) {
		timesheet := model.TimesheetForDate(time.Now())

		task := "work"
		entry := model.TimesheetEntry{Hours: 4, Minutes: 20, Category: "work", Comment: "work", Task: &task}
		if error := timesheet.AddEntry(entry); error != nil {
			t.Errorf("Error adding entry: %v", error)
		}
		saveError := saver.Save(context.Background(), timesheet)
		if saveError != nil {
			t.Errorf("Error pending saving time sheet: %v", saveError)
		}
		statistics, err := query.Daily(context.Background(), timesheet)
		if err != nil {
			t.Errorf("Error getting daily statistics: %v", err)
		}
		assert.Equal(t, 1, len(statistics), "Expected 1 entry statistic, got %d", len(statistics))
		stat := statistics[0]
		assert.Equal(t, uint8(4), stat.Dirty.Hours, "Expected 4 hours, got %d", stat.Dirty.Hours)
		assert.Equal(t, uint8(4), stat.Daily.Hours, "Expected 4 hours, got %d", stat.Daily.Hours)
		assert.Equal(t, uint8(20), stat.Dirty.Minutes, "Expected 20 minutes, got %d", stat.Dirty.Minutes)
		assert.Equal(t, uint8(20), stat.Daily.Minutes, "Expected 20 minutes, got %d", stat.Daily.Minutes)

	})
}

func TestShouldBeAbleToPresentDailyStatisticsOrderShouldNotMatter(t *testing.T) {
	useDb(t, func(saver model.Saver, query model.Queryer) {
		timesheet := model.TimesheetForDate(time.Now())

		task := "work"
		entry := model.TimesheetEntry{Hours: 4, Minutes: 20, Category: "work", Comment: "work", Task: &task}
		if error := timesheet.AddEntry(entry); error != nil {
			t.Errorf("Error adding entry: %v", error)
		}
		entry = model.TimesheetEntry{Hours: 4, Minutes: 30, Category: "work", Comment: "work", Task: &task}
		if error := timesheet.AddEntry(entry); error != nil {
			t.Errorf("Error adding entry: %v", error)
		}
		entry = model.TimesheetEntry{Hours: 4, Minutes: 30, Category: "work_reversed", Comment: "work", Task: &task}
		if error := timesheet.AddEntry(entry); error != nil {
			t.Errorf("Error adding entry: %v", error)
		}
		entry = model.TimesheetEntry{Hours: 4, Minutes: 20, Category: "work_reversed", Comment: "work", Task: &task}
		if error := timesheet.AddEntry(entry); error != nil {
			t.Errorf("Error adding entry: %v", error)
		}
		saveError := saver.Save(context.Background(), timesheet)
		if saveError != nil {
			t.Errorf("Error pending saving time sheet: %v", saveError)
		}
		statistics, err := query.Daily(context.Background(), timesheet)
		if err != nil {
			t.Errorf("Error getting daily statistics: %v", err)
		}
		for _, stat := range statistics {
			assert.Equal(t, uint8(8), stat.Daily.Hours, "Expected 8 hours, got %d", stat.Daily.Hours)
			assert.Equal(t, uint8(50), stat.Daily.Minutes, "Expected 50 minutes, got %d", stat.Daily.Minutes)
		}
	})
}
