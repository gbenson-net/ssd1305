package ssd1305

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestString(t *testing.T) {
	d := &SSD1305{}
	assert.Equal(t, d.String(), "SSD1305{<nil>, <nil>, <nil>, (0,0)}")
}
