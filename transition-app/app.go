package main

import (
	"image"

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
	imgSource  []image.Image
	imgWidget  []*ycwidget.Image
	imgStates  []*ycwidget.AffineState
	win        *app.Window
	theme      *material.Theme
	thumbnails layout.List
}

func newTransitionApp(imgs ...image.Image) *transitionApp {
	a := &transitionApp{
		win:        app.NewWindow(),
		theme:      material.NewTheme(gofont.Collection()),
		thumbnails: layout.List{Axis: layout.Vertical},
	}

	var angle widget.Float
	for _, src := range imgs {
		img := transparentImage(src, 0xa0)
		a.imgSource = append(a.imgSource, src)
		a.imgWidget = append(a.imgWidget, ycwidget.NewImage(img))
		a.imgStates = append(a.imgStates, ycwidget.DragAndRotate(&angle))
	}
	return a
}

func (a *transitionApp) mainloop() error {
	ops := &op.Ops{}

	for e := range a.win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			for i, state := range a.imgStates {
				state.Layout(gtx, a.imgWidget[i].Layout)
			}

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
