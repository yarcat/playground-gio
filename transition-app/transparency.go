package main

import (
	"image"
	"image/color"
	"image/draw"
)

func transparentImage(src image.Image, a uint8) *image.RGBA {
	mask := image.NewRGBA(src.Bounds())
	draw.Draw(mask, mask.Bounds(), &image.Uniform{color.RGBA{A: a}}, image.ZP, draw.Src)

	img := image.NewRGBA(src.Bounds())
	draw.DrawMask(img, img.Bounds(), src, image.ZP, mask, image.ZP, draw.Over)

	return img
}
