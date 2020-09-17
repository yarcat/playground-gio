package main

import (
	"image"
	"math"

	ycwidget "github.com/yarcat/playground-gio/transition-app/widget"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type transitionApp struct {
	images []*ycwidget.Drag
	win    *app.Window
	theme  *material.Theme
	// Real rotation is shiftedRotation - Pi
	shiftedRotation widget.Float
}

func newTransitionApp(imgs ...image.Image) *transitionApp {
	a := &transitionApp{
		win:             app.NewWindow(),
		theme:           material.NewTheme(gofont.Collection()),
		shiftedRotation: widget.Float{Value: math.Pi},
	}
	for _, img := range imgs {
		a.images = append(a.images, &ycwidget.Drag{
			Widget: ycwidget.NewImage(img).Layout,
		})
	}
	return a
}

func (a *transitionApp) mainloop() error {
	ops := &op.Ops{}

	for e := range a.win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			l := len(a.images) - 1
			for _, img := range a.images[:1] {
				img.Layout(gtx)
			}
			a.rotated(gtx, a.images[l].Layout)

			layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(1, material.Slider(
						a.theme,
						&a.shiftedRotation,
						0,
						2*math.Pi,
					).Layout),
				)
			})
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}

func (a *transitionApp) rotated(gtx layout.Context, widget layout.Widget) {
	macro := op.Record(gtx.Ops)
	d := widget(gtx)
	call := macro.Stop()

	defer op.Push(gtx.Ops).Pop()
	o := layout.FPt(d.Size).Mul(0.5)
	op.Affine(f32.Affine2D{}.Rotate(o, a.shiftedRotation.Value-math.Pi)).Add(gtx.Ops)
	call.Add(gtx.Ops)
}
