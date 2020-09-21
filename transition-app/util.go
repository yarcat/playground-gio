package main

import (
	"image"
	"image/png"
	"strings"

	"gioui.org/widget"
)

func mustDecodePNG(data string) image.Image {
	r := strings.NewReader(data)
	img, err := png.Decode(r)
	if err != nil {
		panic(err)
	}
	return img
}

func mustNewIcon(data []byte) *widget.Icon {
	icon, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return icon
}
