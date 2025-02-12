package model_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestShouldBeAbleToCreateTimeSheetOnlyIfIsYYYYMMDD(t *testing.T) {
	t.Parallel()
	timesheet, err := model.NewTimesheet("2021-01-01")
	assert.NotNil(t, timesheet)
	assert.Nil(t, err)
}

func TestShouldNotBeAbleToCreateTimeSheetOnlyIfIsSomethingElse(t *testing.T) {
	t.Parallel()
	cases := []string{"aaa", "2021-11-55", "a2002-12-12"}
	for _, c := range cases {
		t.Run(fmt.Sprintf("Should not be able to create time sheet for %s", c), func(t *testing.T) {
			timesheet, err := model.NewTimesheet(c)
			assert.Nil(t, timesheet)
			var validateDateError *model.ValidateDateError
			if errors.As(err, &validateDateError) {
			} else {
				t.Errorf("Expected error to be of type ValidateDateError, got %v", err)
			}
		})
	}
}

func TestShouldNotBePossibleToAddHolidayNotToSameTimesheet(t *testing.T) {
	t.Parallel()
	holiday, error := model.NewHoliday("2021-01-01", "description")
	if error != nil {
		t.Errorf("Error creating holiday: %v", error)
	}
	timesheet, error := model.NewTimesheet("2021-01-02")
	if error != nil {
		t.Errorf("Error creating time sheet: %v", error)
	}
	error = timesheet.AddHoliday(holiday)
	assert.NotNil(t, error)
}

func TestShouldBePossibleToAddHolidayToSameTimesheet(t *testing.T) {
	t.Parallel()
	holiday, error := model.NewHoliday("2021-01-01", "description")
	if error != nil {
		t.Errorf("Error creating holiday: %v", error)
	}
	timesheet, error := model.NewTimesheet("2021-01-01")
	if error != nil {
		t.Errorf("Error creating time sheet: %v", error)
	}
	error = timesheet.AddHoliday(holiday)
	assert.Nil(t, error)
}

func TestIfOneHolidayThenTimesheetShouldBeSpecific(t *testing.T) {
	t.Parallel()
	holiday, error := model.NewHoliday("2021-01-01", "description")
	if error != nil {
		t.Errorf("Error creating holiday: %v", error)
	}
	timesheet, error := model.NewTimesheet("2021-01-01")
	if error != nil {
		t.Errorf("Error creating time sheet: %v", error)
	}
	error = timesheet.AddHoliday(holiday)
	if error != nil {
		t.Errorf("Error adding holiday: %v", error)
	}
	assert.Equal(t, uint8(0), timesheet.PotentialWorkingTime())
	assert.Equal(t, uint8(8), timesheet.PotentialTotalTime())
	assert.Equal(t, float32(0), timesheet.WorkingTime())
}

func aTimesheet() *model.Timesheet {
	timesheet, error := model.NewTimesheet("2021-01-01")
	if error != nil {
		panic("Error creating time sheet")
	}
	return timesheet
}

func aTimeSheetEntry() model.TimesheetEntry {
	task := "work"
	return model.TimesheetEntry{Hours: 8, Minutes: 0, Category: "work", Comment: "work", Task: &task}
}

func TestShouldNotAllowCreatingEntryWithNotCorrectHours(t *testing.T) {
	t.Parallel()
	entry := aTimeSheetEntry()
	entry.Hours = 25
	err := aTimesheet().AddEntry(entry)

	var invalidEntry *model.InvalidTimesheetEntry
	if errors.As(err, &invalidEntry) {
	} else {
		t.Errorf("Expected error to be of type InvalidTimesheetEntry, got %v", err)
	}
}

func TestShouldNotAllowCreatingEntryWithNotCorrectMinutes(t *testing.T) {
	t.Parallel()
	entry := aTimeSheetEntry()
	entry.Minutes = 61
	err := aTimesheet().AddEntry(entry)

	var invalidEntry *model.InvalidTimesheetEntry
	if errors.As(err, &invalidEntry) {
	} else {
		t.Errorf("Expected error to be of type InvalidTimesheetEntry, got %v", err)
	}
}

func TestShouldNotAllowAddingTimesheetWithBothHoursAndMinutesEmpty(t *testing.T) {
	t.Parallel()
	entry := aTimeSheetEntry()
	entry.Minutes = 0
	entry.Hours = 0
	err := aTimesheet().AddEntry(entry)

	var invalidEntry *model.InvalidTimesheetEntry
	if errors.As(err, &invalidEntry) {
	} else {
		t.Errorf("Expected error to be of type InvalidTimesheetEntry, got %v", err)
	}
}

func TestShouldAllowItemWithoutTask(t *testing.T) {
	t.Parallel()
	entry := aTimeSheetEntry()
	entry.Task = nil
	err := aTimesheet().AddEntry(entry)
	assert.Nil(t, err)
}

func TestShouldSumUpAllEntries(t *testing.T) {
	t.Parallel()
	timesheet := aTimesheet()
	entry := aTimeSheetEntry()
	entry.Hours = 4
	entry.Minutes = 0
	err := timesheet.AddEntry(entry)
	assert.Nil(t, err)

	entry = aTimeSheetEntry()
	entry.Hours = 2
	entry.Minutes = 30
	err = timesheet.AddEntry(entry)
	assert.Nil(t, err)

	assert.Equal(t, uint8(8), timesheet.PotentialWorkingTime())
	assert.Equal(t, uint8(8), timesheet.PotentialTotalTime())
	assert.Equal(t, float32(6.5), timesheet.WorkingTime())

}

func TestShouldAllEntriesCouldSumUpAbove8hours(t *testing.T) {
	t.Parallel()
	timesheet := aTimesheet()
	entry := aTimeSheetEntry()
	entry.Hours = 4
	entry.Minutes = 0
	err := timesheet.AddEntry(entry)
	assert.Nil(t, err)

	entry = aTimeSheetEntry()
	entry.Hours = 4
	entry.Minutes = 30
	err = timesheet.AddEntry(entry)
	assert.Nil(t, err)

	assert.Equal(t, uint8(8), timesheet.PotentialWorkingTime())
	assert.Equal(t, uint8(8), timesheet.PotentialTotalTime())
	assert.Equal(t, float32(8.5), timesheet.WorkingTime())

}

func TestShouldBeAbleToConvertTimesheetToDays(t *testing.T) {
	t.Parallel()
	timesheet, error := model.NewTimesheet("2025-01-25")
	if error != nil {
		t.Errorf("Error creating time sheet: %v", error)
	}
	aDay := timesheet.Day()
	assert.Equal(t, "2025-01-25", aDay.String())
}

func TestShouldBeAbleToConvertTimesheetToWeek(t *testing.T) {
	t.Parallel()
	timesheet, error := model.NewTimesheet("2025-01-25")
	if error != nil {
		t.Errorf("Error creating time sheet: %v", error)
	}
	aWeek := timesheet.Week()
	assert.Equal(t, "2025-01-20 - 2025-01-26", aWeek.String())
}

func TestShouldBeAbleToConvertTimesheetToWeekOnMonday(t *testing.T) {
	t.Parallel()
	timesheet, error := model.NewTimesheet("2025-01-20")
	if error != nil {
		t.Errorf("Error creating time sheet: %v", error)
	}
	aWeek := timesheet.Week()
	assert.Equal(t, "2025-01-20 - 2025-01-26", aWeek.String())
}

func TestShouldBeAbleToConvertTimesheetToWeekOnSunday(t *testing.T) {
	t.Parallel()
	timesheet, error := model.NewTimesheet("2025-01-26")
	if error != nil {
		t.Errorf("Error creating time sheet: %v", error)
	}
	aWeek := timesheet.Week()
	assert.Equal(t, "2025-01-20 - 2025-01-26", aWeek.String())
}

func TestShouldBeAbleToConvertTimesheetToEvenIfNewMonthStartsInMiddle(t *testing.T) {
	t.Parallel()
	timesheet, error := model.NewTimesheet("2025-04-01")
	if error != nil {
		t.Errorf("Error creating time sheet: %v", error)
	}
	aWeek := timesheet.Week()
	assert.Equal(t, "2025-04-01 - 2025-04-06", aWeek.String())
}

func TestShouldBeAbleToConvertTimesheetToEvenIfNewMonthEndsInMiddle(t *testing.T) {
	t.Parallel()
	timesheet, error := model.NewTimesheet("2025-03-31")
	if error != nil {
		t.Errorf("Error creating time sheet: %v", error)
	}
	aWeek := timesheet.Week()
	assert.Equal(t, "2025-03-31 - 2025-03-31", aWeek.String())
}
