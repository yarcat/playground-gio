package main

import (
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/unit"
)

type drag struct {
	dragging bool
	start    f32.Point
	pid      pointer.ID
}

func (d *drag) Offset(cfg unit.Metric, q event.Queue) (offs f32.Point, ok bool) {
	var pe *pointer.Event
	for _, e := range q.Events(d) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}

		switch e.Type {
		case pointer.Press:
			if e.Buttons != pointer.ButtonLeft && e.Source != pointer.Touch {
				continue
			}
			if d.dragging {
				continue
			}
			d.dragging = true
			d.pid = e.PointerID
			d.start = e.Position
		case pointer.Release, pointer.Cancel:
			if !d.dragging || e.PointerID != d.pid {
				continue
			}
			d.dragging = false
		case pointer.Drag:
			if !d.dragging || e.PointerID != d.pid {
				continue
			}
			pe = &e
		}
	}
	if pe == nil {
		return
	}

	offs = pe.Position.Sub(d.start)
	d.start = pe.Position

	return offs, true
}

func (d *drag) Add(ops *op.Ops) {
	op := pointer.InputOp{
		Tag:   d,
		Grab:  d.dragging,
		Types: pointer.Press | pointer.Drag | pointer.Release,
	}
	op.Add(ops)
}
