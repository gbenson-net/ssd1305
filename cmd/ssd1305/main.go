//go:build ignore

package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"gbenson.net/go/ssd1305"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/ssd1306/image1bit"
	"periph.io/x/host/v3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ssd1305:", err)
		os.Exit(1)
	}
}

func run() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	port, err := spireg.Open("")
	if err != nil {
		return err
	}
	defer port.Close()

	dev := ssd1305.SSD1305{
		Port:     port,
		Width:    128,
		Height:   32,
		DC:       gpioreg.ByName("GPIO24"),
		RST:      gpioreg.ByName("GPIO25"),
		StartCol: 4,
	}
	if err = dev.Open(); err != nil {
		return err
	}
	defer dev.Close()
	fmt.Println("Opened", dev.String())

	img := image1bit.NewVerticalLSB(dev.Bounds())
	drawer := font.Drawer{
		Src:  &image.Uniform{C: image1bit.On},
		Dst:  img,
		Face: basicfont.Face7x13,
	}

	for i := range 4 {
		drawer.Dot.X = fixed.I((i & 1) * 16)
		drawer.Dot.Y = fixed.I(i*6 + 12)
		drawer.DrawString("Hello world!")

		if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)

		if err := dev.Halt(); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	for i := range 5 {
		drawer.Dot.X = fixed.I(128 - 35 + i*7)
		drawer.Dot.Y = fixed.I(i*4 + 12)
		drawer.DrawString(fmt.Sprintf("%d", 5-i))

		if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}
	return nil
}
