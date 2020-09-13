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

func minf32(v ...float32) (min float32) {
	if len(v) == 0 {
		return
	}
	min = v[0]
	for _, x := range v[1:] {
		if x < min {
			min = x
		}
	}
	return
}
