package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToDatetime(t *testing.T) {
	assert.Equal(t, "2023-01-06 15:09:27", ToDatetime(1672988967))
}
