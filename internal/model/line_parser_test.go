package model_test

import (
	"errors"
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
	var error *model.EmptyLine
	if errors.As(err, &error) {
	} else {
		t.Errorf("Expected error to be of type InvalidLine, got %v", err)
	}
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
	timesheet, err := parser.ParseLine(aDate())("Category 1.5h Task-123 description")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	assert.NotNil(t, timesheet)
	assert.False(t, timesheet.IsHoliday())
	assert.IsTypef(t, &model.TimesheetEntry{}, timesheet, "Expected type to be TimesheetEntry, got %T", timesheet)
	valued := timesheet.(*model.TimesheetEntry)
	assert.Equal(t, uint8(1), valued.Hours)
	assert.Equal(t, uint8(30), valued.Minutes)
	assert.Equal(t, "description", valued.Comment)
	assert.Equal(t, "Category", valued.Category)
	assert.Equal(t, "Task-123", *valued.Task)

}
