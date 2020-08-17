package main

import (
	"image"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/yarcat/playground-gio/hello/res"
)

var (
	voltOp = paint.NewImageOp(mustDecodeJPG(res.VoltJPG))
)

// Voltmeter is analog voltmeter widget.
type Voltmeter struct {
	Value float32

	changed bool
}

// Layout lays out this widget within given context.
func (v *Voltmeter) Layout(gtx layout.Context) layout.Dimensions {
	p := voltOp.Size()
	imgx, imgy := p.X/2, p.Y/2
	widget.Image{Src: voltOp, Scale: 0.25}.Layout(gtx)

	defer op.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: red(0xff)}.Add(gtx.Ops)

	cx, cy := float32(imgx/2), float32(imgy/2)+float32(imgy)*0.035
	r := f32.Rect(cx, cy-1, cx-cx/3, cy+1)
	dv := float32(v.Value-8) / 8
	phi := float32(math.Pi*0.2 + dv*math.Pi*0.6)
	if phi < math.Pi/8 {
		phi = math.Pi / 8
	}
	if phi > math.Pi*7/8 {
		phi = math.Pi * 7 / 8
	}
	op.Affine(f32.Affine2D{}.Rotate(f32.Pt(cx, cy), phi)).Add(gtx.Ops)
	paint.PaintOp{Rect: r}.Add(gtx.Ops)
	return layout.Dimensions{Size: image.Pt(imgx, imgy)}
}

// Changed returns true if Value was changed since the last time it was checked.
func (v *Voltmeter) Changed() bool {
	changed := v.changed
	v.changed = false
	return changed
}
