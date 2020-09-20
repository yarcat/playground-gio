package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"

	"golang.org/x/image/math/f64"

	"golang.org/x/image/draw"

	"gioui.org/app"
	"github.com/yarcat/playground-gio/transition-app/res"
)

func main() {
	go func() {
		images := make(chan image.Image)
		for _, d := range []struct {
			res string
			r   float64
		}{
			{res.GopherSimplePNG, math.Pi / 6},
			{res.GopherPNG, -math.Pi / 6},
		} {
			d := d
			go func() { images <- rotate(mustDecodePNG(d.res), d.r) }()
		}
		app := newTransitionApp(<-images, <-images)
		if err := app.mainloop(); err != nil {
			log.Fatal("mainloop failed:", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func rotate(src image.Image, r float64) image.Image {
	sw, sh := float64(src.Bounds().Dx()), float64(src.Bounds().Dy())
	dst := image.NewRGBA(image.Rectangle{
		Max: image.Pt(int(sw*1.45), int(sh*1.45)),
	})
	// draw.Copy(dst, image.Point{}, src, dst.Bounds(), draw.Src, nil)
	c, s := math.Cos(r), math.Sin(r)
	x0, y0 := sw/2, sh/2
	dw, dh := float64(dst.Bounds().Dx()), float64(dst.Bounds().Dy())
	fmt.Println(sw, sh)
	fmt.Println(dw, dh)
	dx, dy := (dw-sw)/2, (dh-sh)/2
	fmt.Println(dx, dy)
	m := f64.Aff3{
		c, -s, (-x0)*c - (-y0)*s + x0 + dx,
		s, c, (-x0)*s + (-y0)*c + y0 + dy,
	}
	draw.NearestNeighbor.Transform(dst, m, src, dst.Rect, draw.Src, nil)
	return dst
}
