package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var gopherOp = paint.NewImageOp(mustDecodePNG(gopherPNG))

func main() {
	go func() {
		w := app.NewWindow()
		t := material.NewTheme(gofont.Collection())
		if err := mainloop(w, t); err != nil {
			log.Fatal("mainloop failed:", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func mainloop(w *app.Window, t *material.Theme) error {
	ops := &op.Ops{}

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return widget.Border{
					Color: color.RGBA{A: 0xff, R: 0xff},
					Width: unit.Dp(2),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					cs := gtx.Constraints.Constrain(gopherOp.Size())
					gs := gopherOp.Size()
					k := minf32(1, float32(cs.X)/float32(gs.X), float32(cs.Y)/float32(gs.Y))
					k /= gtx.Metric.PxPerDp
					return widget.Image{Src: gopherOp, Scale: k}.Layout(gtx)
				})
			})

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
