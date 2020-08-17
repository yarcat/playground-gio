package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

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

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(header(t)),
				layout.Flexed(1, body(t)),
				layout.Rigid(footer(t)),
			)

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}

var (
	progress = 40
	wb       = &widget.Bool{Value: false}
	wbtn     = &widget.Clickable{}
	wf       = &widget.Float{Value: 10}
	volt     = &Voltmeter{Value: 10}
)

func body(t *material.Theme) layout.Widget {
	if wf.Changed() {
		volt.Value = wf.Value
	}
	return func(gtx layout.Context) layout.Dimensions {
		d := layout.Dimensions{Size: gtx.Constraints.Min}
		// fmt.Println("body/H1:", layout.Center.Layout(gtx, material.H1(t, "Body").Layout))
		// fmt.Println("body/d :", d)

		widgets := []layout.Widget{
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return widget.Border{
							Width:        unit.Dp(2),
							Color:        color.RGBA{A: 0xff, G: 0xff},
							CornerRadius: unit.Dp(10),
						}.Layout(gtx,
							material.H4(t, "This is H4").Layout)
					}),
					layout.Flexed(1, material.Button(t, wbtn, "This is a button").Layout),
				)
			},
			material.ProgressBar(t, progress).Layout,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(material.Label(t, t.TextSize.Scale(2), "Loader:").Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !wb.Value {
							return layout.Dimensions{}
						}
						return layout.Inset{Left: unit.Dp(8)}.Layout(gtx,
							material.Loader(t).Layout,
						)
					}),
				)
			},
			material.Switch(t, wb).Layout,
			material.Slider(t, wf, 0, 20).Layout,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(material.Caption(t, fmt.Sprintf("V = %05.02f (V)", wf.Value)).Layout),
							layout.Rigid(volt.Layout),
						)
					},
				)
			},
		}

		list := &layout.List{Axis: layout.Vertical}
		list.Layout(gtx, len(widgets), func(gtx layout.Context, i int) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, widgets[i])
		})

		return d
	}
}

func header(t *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		w := material.Label(t, t.TextSize.Scale(2), "Header")
		w.Alignment = text.Middle
		return filled(gtx, red(0x4f), w.Layout)
	}
}

func footer(t *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		w := material.Label(t, t.TextSize.Scale(2), "Footer")
		w.Alignment = text.Middle
		return filled(gtx, green(0x4f), w.Layout)
	}
}
