package style

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestColor(t *testing.T) {
	result := Color("hello", tcell.ColorBlue)
	assert.Equal(t, "[#0000ff::]hello[-::]", result, "should be colorized text")
}
