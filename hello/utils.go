package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"strings"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func red(r uint8) color.RGBA   { return color.RGBA{A: 0xff, R: r} }
func green(g uint8) color.RGBA { return color.RGBA{A: 0xff, G: g} }

func filled(gtx layout.Context, col color.RGBA, w layout.Widget) layout.Dimensions {
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			const rr = 1
			d := gtx.Constraints.Min
			clip.RRect{
				Rect: f32.Rectangle{Max: layout.FPt(d)},
				NE:   rr, NW: rr, SE: rr, SW: rr,
			}.Add(gtx.Ops)
			return fill(gtx, col)
		}),
		layout.Stacked(func(_ layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, w)
		}),
		// layout.Expanded(func(gtx layout.Context) layout.Dimensions {
		// 	d := w(gtx)
		// 	fmt.Println("filled/d  :", d)
		// 	fmt.Println("filled/gtx:", gtx.Constraints)
		// 	return d
		// }),
		// layout.Expanded(w),
	)
}

func fill(gtx layout.Context, col color.RGBA) layout.Dimensions {
	cs := gtx.Constraints
	d := cs.Min
	dr := f32.Rectangle{Max: layout.FPt(d)}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	return layout.Dimensions{Size: d}
}

func mustDecodeJPG(data string) image.Image {
	img, err := jpeg.Decode(strings.NewReader(data))
	if err != nil {
		log.Fatal("Unable to decode JPG:", err)
	}
	return img
}
