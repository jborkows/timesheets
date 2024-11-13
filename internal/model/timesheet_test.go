package model_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestShouldBeAbleToCreateTimeSheetOnlyIfIsYYYYMMDD(t *testing.T) {
	timesheet, err := model.NewTimesheet("2021-01-01")
	assert.NotNil(t, timesheet)
	assert.Nil(t, err)
}

func TestShouldNotBeAbleToCreateTimeSheetOnlyIfIsSomethingElse(t *testing.T) {
	cases := []string{"aaa", "2021-11-55", "a2002-12-12"}
	for _, c := range cases {
		t.Run(fmt.Sprintf("Should not be able to create timesheet for %s", c), func(t *testing.T) {
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
