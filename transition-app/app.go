package main

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	ycwidget "github.com/yarcat/playground-gio/transition-app/widget"
)

type transitionApp struct {
	win           *app.Window
	theme         *material.Theme
	thumbnails    layout.List
	thumbnailImgs []*ycwidget.Image
}

func newTransitionApp(imgs ...image.Image) *transitionApp {
	thumbnails := make([]*ycwidget.Image, 0, len(imgs))
	for _, img := range imgs {
		thumbnails = append(thumbnails, ycwidget.NewImage(img))
	}
	return &transitionApp{
		win:           app.NewWindow(),
		theme:         material.NewTheme(gofont.Collection()),
		thumbnailImgs: thumbnails,
	}
}

func (a *transitionApp) mainloop() error {
	ops := &op.Ops{}

	thumbs := make(thumbnails, 0, len(a.thumbnailImgs))
	for _, img := range a.thumbnailImgs {
		thumbs = append(thumbs, &clickable{widget: img.Layout})
	}

	selected := 0
	for e := range a.win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			for i, c := range thumbs {
				if c.button.Clicked() {
					selected = i
				}
			}
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					axis := layout.Horizontal
					if gtx.Constraints.Max.X < gtx.Constraints.Max.Y {
						axis = layout.Vertical
					}
					return layout.Flex{Axis: axis}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Center.Layout(gtx, a.thumbnailImgs[selected].Layout)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							l := material.H5(a.theme, "PREVIEW\nPLACEHOLDER")
							return layout.Center.Layout(gtx, l.Layout)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return thumbs.Layout(gtx, &a.thumbnails, selected)
					})
				}),
			)

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}

const (
	thumbnailInsetDp = 10
	thumbnailSizePx  = 150
)

type thumbnails []*clickable

func (th thumbnails) Layout(gtx layout.Context, list *layout.List, selected int) layout.Dimensions {
	return list.Layout(gtx, len(th), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(thumbnailInsetDp)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return th[index].Layout(gtx, selected == index)
		})
	})
}

type clickable struct {
	widget layout.Widget
	button widget.Clickable
}

func (btn *clickable) Layout(gtx layout.Context, selected bool) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	cornerRad := unit.Dp(10)
	d := material.Clickable(gtx, &btn.button, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max = image.Pt(thumbnailSizePx, thumbnailSizePx)
		if !selected {
			return btn.widget(gtx)
		}
		return widget.Border{
			Color:        color.RGBA{A: 0xff, R: 0x1f, G: 0x1f, B: 0x1f},
			Width:        unit.Dp(1),
			CornerRadius: cornerRad,
		}.Layout(gtx, btn.widget)
	})
	call := macro.Stop()

	defer op.Push(gtx.Ops).Pop()
	clip.RRect{
		Rect: f32.Rectangle{Max: layout.FPt(d.Size)},
		NW:   cornerRad.V * gtx.Metric.PxPerDp,
		NE:   cornerRad.V * gtx.Metric.PxPerDp,
		SW:   cornerRad.V * gtx.Metric.PxPerDp,
		SE:   cornerRad.V * gtx.Metric.PxPerDp,
	}.Add(gtx.Ops)
	call.Add(gtx.Ops)

	return d
}
