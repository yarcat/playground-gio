package widget

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"

	ycgest "github.com/yarcat/playground-gio/transition-app/gesture"
)

// AffineState is for dragging and rotating widgets.
type AffineState struct {
	gest   ycgest.Drag
	offs   f32.Point
	float  *widget.Float
	widget layout.Widget

	changed bool
}

// Changed returns true if there was a state change since the last update.
// Calling his method resets the changed state to false.
func (state *AffineState) Changed() bool {
	c := state.changed
	state.changed = false
	return c
}

// DragAndRotate allows to drag and rotate widgets. Rotation requires an
// external widget that would modify the float. Dragging is applied on top
// of a widget.
func DragAndRotate(widget layout.Widget, float *widget.Float) AffineState {
	return AffineState{
		widget: widget,
		float:  float,
	}
}

// Layout handles events and lays out a widget.
func (state *AffineState) Layout(gtx layout.Context) layout.Dimensions {
	if state.float.Changed() {
		state.changed = true
	}

	if offs, ok := state.gest.Offset(gtx.Metric, gtx); ok {
		state.changed = true
		state.offs = state.offs.Add(offs)
	}

	macro := op.Record(gtx.Ops)
	d := state.widget(gtx)
	call := macro.Stop()

	defer op.Push(gtx.Ops).Pop()
	// Not translating pointer area to ensure its offset is always calculated
	// relatively to the same origin.
	minOffs := image.Pt(int(state.offs.X), int(state.offs.Y))
	rect := image.Rectangle{Min: minOffs, Max: minOffs.Add(d.Size)}
	pointer.Rect(rect).Add(gtx.Ops)
	state.gest.Add(gtx.Ops)

	op.Affine(f32.Affine2D{}.
		Rotate(layout.FPt(d.Size).Mul(0.5), state.float.Value).
		Offset(state.offs),
	).Add(gtx.Ops)
	call.Add(gtx.Ops)

	// AffineState represents a floating draggable widget, which doesn't really
	// fit any layout management.
	// TODO(yarcat): Don't allow to drag outside of the constraint area.
	return layout.Dimensions{Size: gtx.Constraints.Min}
}
