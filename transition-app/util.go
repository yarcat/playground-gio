package main

import (
	"image"
	"image/png"
	"io"
)

func mustDecodePNG(r io.Reader) image.Image {
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
