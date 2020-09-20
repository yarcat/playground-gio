package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	ycwidget "github.com/yarcat/playground-gio/transition-app/widget"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type transitionApp struct {
	rawImg      []image.Image
	images      []*ycwidget.Image
	states      []*ycwidget.AffineState
	angles      []widget.Float
	lastChanged *widget.Float
	win         *app.Window
	theme       *material.Theme
	thumbnails  layout.List
	// Real rotation is shiftedRotation - Pi
	shiftedRotation widget.Float
}

func newTransitionApp(imgs ...image.Image) *transitionApp {
	a := &transitionApp{
		win:             app.NewWindow(),
		theme:           material.NewTheme(gofont.Collection()),
		shiftedRotation: widget.Float{Value: math.Pi},
		thumbnails:      layout.List{Axis: layout.Vertical},
	}
	halfTransparent := image.Uniform{color.RGBA{A: 0x80}}

	a.images = make([]*ycwidget.Image, len(imgs))
	a.states = make([]*ycwidget.AffineState, len(imgs))
	a.angles = make([]widget.Float, len(imgs))
	a.rawImg = make([]image.Image, len(imgs))
	for i, raw := range imgs {
		a.rawImg[i] = raw

		mask := image.NewRGBA(raw.Bounds())
		draw.Draw(mask, mask.Bounds(), &halfTransparent, image.ZP, draw.Src)

		img := image.NewRGBA(raw.Bounds())
		draw.DrawMask(img, img.Bounds(), raw, image.ZP, mask, image.ZP, draw.Over)
		imgWidget := ycwidget.NewImage(img)
		a.images[i] = imgWidget
		w := ycwidget.DragAndRotate(&a.angles[i])
		a.states[i] = &w
	}
	return a
}

func (a *transitionApp) mainloop() error {
	ops := &op.Ops{}

	// a := [2]int{0xff, 0x00}
	// da := [2]int{-0x0a, 0x0a}

	for e := range a.win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			for i, state := range a.states {
				if state.Changed() {
					a.lastChanged = &a.angles[i]
				}
				state.Layout(gtx, a.images[i].Layout)
			}

			a.layoutRotationSlider(gtx)
			a.layoutThumbnails(gtx)

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

const thumbnailSize = 150

var thumbnailInset = unit.Dp(2)

func (a *transitionApp) layoutThumbnails(gtx layout.Context) layout.Dimensions {
	return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return a.thumbnails.Layout(gtx, len(a.images), func(gtx layout.Context, index int) layout.Dimensions {
			return layout.UniformInset(thumbnailInset).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max = image.Pt(thumbnailSize, thumbnailSize)
				return widget.Border{
					Color:        color.RGBA{A: 0xff, R: 0x1f, G: 0x1f, B: 0x1f},
					Width:        unit.Dp(1),
					CornerRadius: unit.Dp(10),
				}.Layout(gtx, a.images[index].Layout)
			})
		})
	})
}
