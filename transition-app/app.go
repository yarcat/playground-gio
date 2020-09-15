package main

import (
	"image"

	ycwidget "github.com/yarcat/playground-gio/transition-app/widget"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

type transitionApp struct {
	imgs  []*ycwidget.DragImage
	win   *app.Window
	theme *material.Theme
}

func newTransitionApp(imgs ...image.Image) *transitionApp {
	a := &transitionApp{
		win:   app.NewWindow(),
		theme: material.NewTheme(gofont.Collection()),
	}
	for _, img := range imgs {
		a.imgs = append(a.imgs, ycwidget.NewDragImage(img))
	}
	return a
}

func (a *transitionApp) mainloop() error {
	ops := &op.Ops{}

	for e := range a.win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			for _, img := range a.imgs {
				img.Layout(gtx)
			}

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
