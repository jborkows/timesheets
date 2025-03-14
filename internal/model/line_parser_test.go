package model_test

import (
	"strings"
	"testing"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

func aDate() model.DateInfo {
	return model.DateInfo{Value: "2021-01-01"}
}

func TestEmptyLineIsInvalid(t *testing.T) {
	t.Parallel()
	parser := model.Parser{
		HolidayClassifier: func(a *model.DateInfo) bool { return false },
	}
	timesheet, err := parser.ParseLine(aDate())("")

	assert.Nil(t, timesheet)
	assert.Equal(t, model.ErrEmptyLine, err)
}

func TestHolidayIsValid(t *testing.T) {
	t.Parallel()

	parser := model.Parser{
		HolidayClassifier: func(a *model.DateInfo) bool { return true },
	}
	timesheet, err := parser.ParseLine(aDate())("aaaa")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	assert.NotNil(t, timesheet)
	assert.True(t, timesheet.IsHoliday())
}

func workingDayParser() model.Parser {
	return model.Parser{
		HolidayClassifier: func(a *model.DateInfo) bool { return false },
		IsCategory:        func(text string) bool { return text == "Category" },
		IsTask:            func(text string) bool { return strings.HasPrefix(text, "Task-") },
	}
}

func TestShouldParseTextWithDecimalTime(t *testing.T) {
	t.Parallel()

	parser := workingDayParser()
	timesheet, err := parser.ParseLine(aDate())("Category 1.5 Task-123 description")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	assert.NotNil(t, timesheet)
	assert.False(t, timesheet.IsHoliday())
	assert.IsTypef(t, &model.TimesheetEntry{}, timesheet, "Expected type to be TimesheetEntry, got %T", timesheet)
	valued := timesheet.(*model.TimesheetEntry)
	assert.Equal(t, "Category", valued.Category)
	assert.Equal(t, uint8(1), valued.Hours)
	assert.Equal(t, uint8(30), valued.Minutes)
	assert.Equal(t, "Task-123", *valued.Task)
	assert.Equal(t, "description", valued.Comment)
}

func TestOnlyWrongCategory(t *testing.T) {
	t.Parallel()

	parser := workingDayParser()
	timesheet, err := parser.ParseLine(aDate())("Cate ")
	assert.Nil(t, timesheet)
	assert.Equal(t, model.ErrInvalidCategory, err)
}

func TestValidHourShouldNoBeReportedAsBothZeros(t *testing.T) {
	t.Parallel()

	parser := workingDayParser()
	entry, err := parser.ParseLine(aDate())("Category 1.0")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	time, ok := entry.(*model.TimesheetEntry)
	if !ok {
		t.Fatalf("Expected TimesheetEntry, got %T", entry)
	}
	assert.Equal(t, uint8(1), time.Hours)
}

func TestShouldParseTextWithDecimalTwoPlacesTime(t *testing.T) {
	t.Parallel()

	parser := workingDayParser()
	timesheet, err := parser.ParseLine(aDate())("Category 1.75 Task-123 description")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	assert.NotNil(t, timesheet)
	assert.False(t, timesheet.IsHoliday())
	valued := timesheet.(*model.TimesheetEntry)
	assert.Equal(t, uint8(1), valued.Hours)
	assert.Equal(t, uint8(45), valued.Minutes)

}

func TestShouldAllowSingleHour(t *testing.T) {
	t.Parallel()

	parser := workingDayParser()
	timesheet, err := parser.ParseLine(aDate())("Category 1h Task-123 description")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	valued := timesheet.(*model.TimesheetEntry)
	assert.Equal(t, uint8(1), valued.Hours)
	assert.Equal(t, uint8(0), valued.Minutes)
}

func TestShouldAllowSingleMinutes(t *testing.T) {
	t.Parallel()

	parser := workingDayParser()
	timesheet, err := parser.ParseLine(aDate())("Category 30m meeting")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	valued := timesheet.(*model.TimesheetEntry)
	assert.Equal(t, uint8(0), valued.Hours)
	assert.Equal(t, uint8(30), valued.Minutes)
	assert.Nil(t, valued.Task)
}

func TestShouldAllowHourWithMinutes(t *testing.T) {
	t.Parallel()

	parser := workingDayParser()
	timesheet, err := parser.ParseLine(aDate())("Category 1h30m Task-123 description")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	valued := timesheet.(*model.TimesheetEntry)
	assert.Equal(t, uint8(1), valued.Hours)
	assert.Equal(t, uint8(30), valued.Minutes)
}

func TestShouldAllowHourWithMinutesButTheyShouldNotGoToDescription(t *testing.T) {
	t.Parallel()

	parser := workingDayParser()
	timesheet, err := parser.ParseLine(aDate())("Category 1h30m description")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	valued := timesheet.(*model.TimesheetEntry)
	assert.Equal(t, "description", valued.Comment)
}

func TestShouldParseTextWithDecimalAboveTwoPlaces(t *testing.T) {
	t.Parallel()

	parser := workingDayParser()
	_, err := parser.ParseLine(aDate())("Category 1.753 Task-123 description")
	assert.Equal(t, model.ErrInvalidTime, err)
}

func TestIfTaskCouldNottBeMatchedItBecameComment(t *testing.T) {
	t.Parallel()
	parser := workingDayParser()
	timesheet, err := parser.ParseLine(aDate())("Category 1.5 Txsk-123 description")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	valued := timesheet.(*model.TimesheetEntry)
	assert.Nil(t, valued.Task)
	assert.Equal(t, "Txsk-123 description", valued.Comment)

}
