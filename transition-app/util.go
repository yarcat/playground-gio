package main

import (
	"image"
	"image/png"
	"strings"
)

func mustDecodePNG(data string) image.Image {
	r := strings.NewReader(data)
	img, err := png.Decode(r)
	if err != nil {
		panic(err)
	}
	return img
}

func minf32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
