package widget

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	ycf32 "github.com/yarcat/playground-gio/transition-app/f32"
	ycgest "github.com/yarcat/playground-gio/transition-app/gesture"
)

// DragImage implements a draggable image.
type DragImage struct {
	img  paint.ImageOp
	gest ycgest.Drag
	offs f32.Point
}

// NewDragImage returns new draggable image.
func NewDragImage(img image.Image) *DragImage {
	return &DragImage{img: paint.NewImageOp(img)}
}

// Layout lays out the image by taking the minimal space close the image size.
func (img *DragImage) Layout(gtx layout.Context) layout.Dimensions {
	if offs, ok := img.gest.Offset(gtx.Metric, gtx); ok {
		img.offs = img.offs.Add(offs)
	}

	stack := op.Push(gtx.Ops)
	op.Offset(img.offs).Add(gtx.Ops)
	d := layoutImg(gtx, img.img)
	stack.Pop()

	minOffs := image.Pt(int(img.offs.X), int(img.offs.Y))
	rect := image.Rectangle{Min: minOffs, Max: minOffs.Add(d.Size)}
	pointer.Rect(rect).Add(gtx.Ops)
	img.gest.Add(gtx.Ops)

	return d
}

func layoutImg(gtx layout.Context, img paint.ImageOp) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dim := func(gtx layout.Context) layout.Dimensions {
		cs := gtx.Constraints.Constrain(img.Size())
		gs := img.Size()
		k := ycf32.Minf32(1, float32(cs.X)/float32(gs.X), float32(cs.Y)/float32(gs.Y))
		k /= gtx.Metric.PxPerDp
		var d layout.Dimensions
		if cs.X < cs.Y {
			d.Size = image.Pt(cs.X, cs.X)
		} else {
			d.Size = image.Pt(cs.Y, cs.Y)
		}
		gtx.Constraints.Max = d.Size
		widget.Image{Src: img, Scale: k}.Layout(gtx)
		clip.Rect{Max: d.Size}.Add(gtx.Ops)
		return d
	}(gtx)
	call := macro.Stop()

	func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints = layout.Constraints{
			Min: dim.Size, Max: dim.Size,
		}
		return widget.Border{
			Color: color.RGBA{A: 0xff, R: 0xff},
			Width: unit.Dp(2),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			call.Add(gtx.Ops)
			return dim
		})
	}(gtx)

	return dim
}
