// Package ssd1305 controls a 132x64 monochrome OLED display via an
// SSD1305 controller.
//
// Datasheet:
// https://cdn-shop.adafruit.com/datasheets/SSD1305.pdf
package ssd1305

import (
	"fmt"
	"image"
	"time"

	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/devices/v3/ssd1306/image1bit"
)

type SSD1305 struct {
	// Display size.  Defaults to 132x64 if not specified.
	Width, Height int

	// SPI port to use.
	Port spi.Port

	// Data/Command selection pin (high for data, low for command).
	DC gpio.PinOut

	// Reset pin, low active. May be nil.
	RST gpio.PinOut

	// Start column.  Defaults to 0.
	StartCol int

	conn conn.Conn
	rect image.Rectangle
}

// Open connects to an SSD1305 display controller.
func (d *SSD1305) Open() error {
	if d.conn != nil {
		return ErrConnected
	} else if d.Port == nil {
		panic("nil port")
	} else if d.DC == nil {
		panic("nil DC pin")
	}

	w := d.Width
	if w == 0 {
		w = 132
	} else if w < 1 || w > 132 {
		return fmt.Errorf("ssd1305: invalid width %d", w)
	}
	d.rect.Max.X = w

	h := d.Height
	if h == 0 {
		h = 64
	} else if h < 8 || h > 64 || h&7 != 0 {
		return fmt.Errorf("ssd1305: invalid height %d", h)
	}
	d.rect.Max.Y = h

	sc := d.StartCol
	if sc < 0 || sc+w > 132 {
		return fmt.Errorf("ssd1305: invalid start column %d", sc)
	}

	if c, err := d.Port.Connect(3300*physic.KiloHertz, spi.Mode0, 8); err != nil {
		return err
	} else {
		d.conn = c
	}

	return d.Reset()
}

// Close implements [io.Closer].
func (d *SSD1305) Close() error {
	if d.conn == nil {
		return ErrNotConnected
	}
	defer func() { d.conn = nil }()
	return d.sendCommand([]byte{
		0xAE, // Set Display OFF
	})
}

// String implements [fmt.Stringer].
func (d *SSD1305) String() string {
	return fmt.Sprintf("SSD1305{%v, %v, %v, %s}", d.conn, d.DC, d.RST, d.rect.Max)
}

// Reset resets an SSD1305 display controller.
func (d *SSD1305) Reset() error {
	if d.conn == nil {
		return ErrNotConnected
	}

	if rp := d.RST; rp != nil {
		if err := rp.Out(gpio.Low); err != nil {
			return err
		}
		time.Sleep(10 * time.Millisecond)

		if err := rp.Out(gpio.High); err != nil {
			return err
		}
	}

	return d.sendCommand([]byte{
		0xAE,       // Set Display OFF
		0x40,       // Set Display Start Line
		0x81, 0x80, // Set Contrast Control for BANK0
		0xA1,       // Set Segment Re-map
		0xA6,       // Set Normal/Inverse Display
		0xA8, 0x1F, // Set Multiplex Ratio
		0xC8,       // Set COM Output Scan Direction
		0xD3, 0x00, // Set Display Offset
		0xD5, 0xF0, // Set Display Clock Divide Ratio/ Oscillator Frequency
		0xD8, 0x05, // Set Area Color Mode ON/OFF & Low Power Display Mode
		0xD9, 0xC2, // Set pre-charge period
		0xDA, 0x12, // Set COM Pins Hardware Configuration
		0xDB, 0x08, // Set VCOMH Deselect Level
		0xAF, // Set Display ON (Normal Brightness)
	})
}

// Bounds implements [display.Drawer].
func (d *SSD1305) Bounds() image.Rectangle {
	return d.rect
}

// Draw implements [display.Drawer].
func (d *SSD1305) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
	var next []byte
	if img, ok := src.(*image1bit.VerticalLSB); ok && r == d.rect &&
		img.Rect == d.rect && sp.X == 0 && sp.Y == 0 {
		// Exact size, full frame, image1bit encoding: fast path!
		next = img.Pix
	} else {
		panic("not implemented")
	}

	sc := d.StartCol
	sclo := byte(sc & 15)
	schi := byte(sc>>4) | 0x10

	w := d.rect.Max.X
	for page := range byte(d.rect.Max.Y / 8) {
		if err := d.sendCommand([]byte{
			0xB0 + page, // Set Page Start Address for Page Addressing Mode
			sclo,        // Set Lower Column Start Address for Page Addressing Mode
			schi,        // Set Higher Column Start Address for Page Addressing Mode
		}); err != nil {
			return err
		}

		if err := d.DC.Out(gpio.High); err != nil {
			return err
		}

		if err := d.conn.Tx(next[:w], nil); err != nil {
			return err
		}

		next = next[w:]
	}

	return nil
}

func (d *SSD1305) sendCommand(b []byte) error {
	if err := d.DC.Out(gpio.Low); err != nil {
		return err
	}
	return d.conn.Tx(b, nil)
}
