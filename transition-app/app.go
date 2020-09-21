package main

import (
	"image"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"

	ycwidget "github.com/yarcat/playground-gio/transition-app/widget"
)

type transitionApp struct {
	imgSource  []image.Image
	animations []*FrameSet
	win        *app.Window
	theme      *material.Theme
}

func newTransitionApp(imgs ...image.Image) *transitionApp {
	imgSource := make([]image.Image, 0, len(imgs))
	animations := make([]*FrameSet, 0, len(imgs))
	for i, src := range imgs {
		imgSource = append(imgSource, src)
		var frames int
		if i > 0 {
			frames = 50
		}
		fs := ApplyTransparency(src, frames, 50*time.Millisecond)
		animations = append(animations, fs)
	}
	return &transitionApp{
		win:        app.NewWindow(),
		theme:      material.NewTheme(gofont.Collection()),
		imgSource:  imgSource,
		animations: animations,
	}
}

func (app *transitionApp) mainloop() error {
	ops := &op.Ops{}

	var opaqueFirstImage *ycwidget.Image
	if len(app.imgSource) > 0 {
		opaqueFirstImage = ycwidget.NewImage(app.imgSource[0])
	}

	th := material.NewTheme(gofont.Collection())
	var opaque widget.Bool
	opaqueCheckBox := material.CheckBox(th, &opaque, "Use opaque bottom image")

	for e := range app.win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			for i, widget := range app.animations {
				// Experimenting to understand whether there is an output
				// difference if we don't make our first image transparent.
				// My theory is that it would make the mid-transition state
				// more colorful.
				if i == 0 && opaque.Value {
					opaqueFirstImage.Layout(gtx)
				} else {
					widget.Layout(gtx)
				}
			}

			layout.NW.Layout(gtx, opaqueCheckBox.Layout)

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
