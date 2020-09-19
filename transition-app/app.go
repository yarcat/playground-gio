package main

import (
	"image"
	"math"

	ycwidget "github.com/yarcat/playground-gio/transition-app/widget"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type transitionApp struct {
	images      []*ycwidget.Image
	states      []*ycwidget.AffineState
	angles      []widget.Float
	lastChanged *widget.Float
	win         *app.Window
	theme       *material.Theme
	// Real rotation is shiftedRotation - Pi
	shiftedRotation widget.Float
}

func newTransitionApp(imgs ...image.Image) *transitionApp {
	a := &transitionApp{
		win:             app.NewWindow(),
		theme:           material.NewTheme(gofont.Collection()),
		shiftedRotation: widget.Float{Value: math.Pi},
	}
	a.images = make([]*ycwidget.Image, len(imgs))
	a.states = make([]*ycwidget.AffineState, len(imgs))
	a.angles = make([]widget.Float, len(imgs))
	for i, img := range imgs {
		imgWidget := ycwidget.NewImage(img)
		a.images[i] = imgWidget
		w := ycwidget.DragAndRotate(imgWidget.Layout, &a.angles[i])
		a.states[i] = &w
	}
	return a
}

func (a *transitionApp) mainloop() error {
	ops := &op.Ops{}

	for e := range a.win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			for i, img := range a.states {
				if img.Changed() {
					a.lastChanged = &a.angles[i]
				}
				img.Layout(gtx)
			}

			a.layoutRotationSlider(gtx)

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}

func (a *transitionApp) layoutRotationSlider(gtx layout.Context) layout.Dimensions {
	if a.lastChanged == nil {
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
	return layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx,
			layout.Flexed(1, material.Slider(
				a.theme,
				a.lastChanged,
				0,
				2*math.Pi,
			).Layout),
		)
	})
}
