package model_test

import (
	"errors"
	"testing"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

func aDate() model.DateInfo {
	return model.DateInfo{Value: "2021-01-01"}
}

func TestEmptyLineIsInvalid(t *testing.T) {
	t.Parallel()
	timesheet, err := model.ParseLine(aDate(), func(a *model.DateInfo) bool { return false })("")

	assert.Nil(t, timesheet)
	var error *model.EmptyLine
	if errors.As(err, &error) {
	} else {
		t.Errorf("Expected error to be of type InvalidLine, got %v", err)
	}
}

func TestHolidayIsValid(t *testing.T) {
	t.Parallel()

	timesheet, err := model.ParseLine(aDate(), func(a *model.DateInfo) bool { return true })("aaaa")
	if err != nil {
		t.Fatalf("Error parsing line: %v", err)
	}
	assert.NotNil(t, timesheet)
	assert.True(t, timesheet.IsHoliday())
}
