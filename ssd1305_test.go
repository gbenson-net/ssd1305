package ssd1305

import (
	"errors"
	"image"
	"testing"

	"gotest.tools/v3/assert"
	"periph.io/x/conn/v3/display"
	"periph.io/x/devices/v3/ssd1306/image1bit"
)

func TestString(t *testing.T) {
	d := &SSD1305{}
	assert.Equal(t, d.String(), "SSD1305{<nil>, <nil>, <nil>, (0,0)}")
}

func TestDrawer(t *testing.T) {
	d := display.Drawer(&SSD1305{
		rect: image.Rect(1, 2, 3, 4),
	})

	assert.Equal(t, d.ColorModel(), image1bit.BitModel)
	assert.Equal(t, d.Bounds(), image.Rect(1, 2, 3, 4))
	assert.Check(t, errors.Is(d.Halt(), ErrNotConnected))
}
