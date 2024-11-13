package example_test

import (
	"github.com/jborkows/timesheets/internal/example"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CheckAdderTest(t *testing.T) {
	assert.Equal(t, example.Example(1, 2), 3, "1 + 2 = 3")
}
