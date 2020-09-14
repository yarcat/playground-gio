package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	ycwidget "github.com/yarcat/playground-gio/transition-app/widget"
)

var img = ycwidget.NewDragImage(mustDecodePNG(gopherPNG))

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

			img.Layout(gtx)

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
