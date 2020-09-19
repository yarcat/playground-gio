package widget

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	ycf32 "github.com/yarcat/playground-gio/transition-app/f32"
)

// Image implements a widget for drawing images.
type Image struct {
	op paint.ImageOp
}

// NewImage returns a widget for drawing images.
func NewImage(img image.Image) *Image {
	return &Image{
		op: paint.NewImageOp(img),
	}
}

// Layout lays out the image by taking the minimal space close the image size.
func (w *Image) Layout(gtx layout.Context) layout.Dimensions {
	img := w.op

	cs := gtx.Constraints.Constrain(img.Size())
	gs := img.Size()
	k := ycf32.Minf32(1, float32(cs.X)/float32(gs.X), float32(cs.Y)/float32(gs.Y))
	d := layout.Dimensions{Size: image.Pt(int(float32(gs.X)*k), int(float32(gs.X)*k))}
	k /= gtx.Metric.PxPerDp
	gtx.Constraints.Max = d.Size
	widget.Image{Src: img, Scale: k}.Layout(gtx)
	clip.Rect{Max: d.Size}.Add(gtx.Ops)
	return d
}
