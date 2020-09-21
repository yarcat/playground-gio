package main

import (
	"image"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
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
		var opts []FrameSetOptionFunc
		if i%2 == 0 {
			opts = append(opts, ReversePlayback)
		}
		fs := ApplyTransparency(src, 10, 200*time.Millisecond, opts...)
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

	for e := range app.win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			for _, widget := range app.animations {
				widget.Layout(gtx)
			}

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
