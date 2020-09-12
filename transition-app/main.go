package main

import (
	"log"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var gopherOp paint.ImageOp

func init() {
	img := mustDecodePNG(strings.NewReader(gopherPNG))
	gopherOp = paint.NewImageOp(img)
}

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

			gs := gopherOp.Size()
			cs := gtx.Constraints.Max

			k := minf32(float32(cs.X)/float32(gs.X), float32(cs.Y)/float32(gs.Y))
			layout.S.Layout(gtx, widget.Image{Src: gopherOp, Scale: k / 2}.Layout)

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
