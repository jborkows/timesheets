package model_test

import (
	"errors"
	"testing"

	"github.com/jborkows/timesheets/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestEmptyLineIsInvalid(t *testing.T) {
	t.Parallel()
	timesheet, err := model.ParseLine("", nil)

	assert.Nil(t, timesheet)
	var error *model.EmptyLine
	if errors.As(err, &error) {
	} else {
		t.Errorf("Expected error to be of type InvalidLine, got %v", err)
	}
}
