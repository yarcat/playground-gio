package main

import (
	"image"
	"image/color"
	"time"

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

	"golang.org/x/exp/shiny/materialdesign/icons"

	ycwidget "github.com/yarcat/playground-gio/transition-app/widget"
)

type transitionApp struct {
	win           *app.Window
	theme         *material.Theme
	thumbnails    layout.List
	thumbnailImgs []*ycwidget.Image
	animations    []*FrameSet
}

func newTransitionApp(imgs ...image.Image) *transitionApp {
	thumbnails := make([]*ycwidget.Image, 0, len(imgs))
	animations := make([]*FrameSet, 0, len(imgs))
	for i, img := range imgs {
		thumbnails = append(thumbnails, ycwidget.NewImage(img))
		frames := 25
		duration := 100 * time.Millisecond
		var opts []FrameSetOptionFunc
		if i == 0 {
			opts = append(opts, ReversePlayback)
		}
		animations = append(animations, ApplyTransparency(img, frames, duration, opts...))
	}
	return &transitionApp{
		win:           app.NewWindow(),
		theme:         material.NewTheme(gofont.Collection()),
		thumbnailImgs: thumbnails,
		animations:    animations,
	}
}

type avState int

const (
	avStatePaused avState = iota
	avStatePlaying
)

var avIcons = [2]*widget.Icon{mustNewIcon(icons.AVPlayArrow), mustNewIcon(icons.AVPause)}

func (state avState) icon() *widget.Icon {
	return avIcons[state]
}

func (state avState) change() avState { return 1 - state }

func (app *transitionApp) mainloop() error {
	ops := &op.Ops{}

	thumbs := make(thumbnails, 0, len(app.thumbnailImgs))
	for _, img := range app.thumbnailImgs {
		thumbs = append(thumbs, &clickable{widget: img.Layout})
	}

	selected := 0

	var (
		avCtrl  widget.Clickable
		avState avState
	)

	for e := range app.win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			for i, c := range thumbs {
				if c.button.Clicked() {
					selected = i
				}
			}

			if avCtrl.Clicked() {
				avState = avState.change()
			}

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					axis := layout.Horizontal
					if gtx.Constraints.Max.X < gtx.Constraints.Max.Y {
						axis = layout.Vertical
					}
					return layout.Flex{Axis: axis}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Center.Layout(gtx, app.thumbnailImgs[selected].Layout)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Stack{Alignment: layout.Center}.Layout(gtx,
								layout.Stacked(func(gtx layout.Context) layout.Dimensions {
									defer op.Record(gtx.Ops).Stop()
									return layout.Center.Layout(gtx, app.thumbnailImgs[selected].Layout)
								}),
								layout.Expanded(func(gtx layout.Context) layout.Dimensions {
									if avState == avStatePaused {
										l := material.H5(app.theme, "PREVIEW\nPLACEHOLDER")
										return layout.Center.Layout(gtx, l.Layout)
									}
									var d layout.Dimensions
									for _, anim := range app.animations {
										d = anim.Layout(gtx)
									}
									return d
								}),
								layout.Expanded(func(gtx layout.Context) layout.Dimensions {
									return layout.S.Layout(gtx, material.IconButton(app.theme, &avCtrl, avState.icon()).Layout)
								}),
							)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return thumbs.Layout(gtx, &app.thumbnails, selected)
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
