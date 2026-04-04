package ssd1305

import "errors"

var (
	ErrConnected    = errors.New("already connected")
	ErrNotConnected = errors.New("not connected")
)
