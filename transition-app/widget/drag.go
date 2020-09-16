package widget

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"

	ycgest "github.com/yarcat/playground-gio/transition-app/gesture"
)

// Drag implements a draggable image.
type Drag struct {
	Widget layout.Widget
	gest   ycgest.Drag
	offs   f32.Point
}

// Layout lays out the underlying image and makes its area draggable.
func (drag *Drag) Layout(gtx layout.Context) layout.Dimensions {
	if offs, ok := drag.gest.Offset(gtx.Metric, gtx); ok {
		drag.offs = drag.offs.Add(offs)
	}

	stack := op.Push(gtx.Ops)
	op.Offset(drag.offs).Add(gtx.Ops)
	d := drag.Widget(gtx)
	stack.Pop()

	minOffs := image.Pt(int(drag.offs.X), int(drag.offs.Y))
	rect := image.Rectangle{Min: minOffs, Max: minOffs.Add(d.Size)}
	pointer.Rect(rect).Add(gtx.Ops)
	drag.gest.Add(gtx.Ops)

	return d
}
