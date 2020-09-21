package main

import (
	"image"
	"time"

	"gioui.org/layout"
	"gioui.org/op"

	ycwidget "github.com/yarcat/playground-gio/transition-app/widget"
)

// FrameSet caches image transformations.
type FrameSet struct {
	frames   []image.Image
	images   []*ycwidget.Image
	curFrame int // Current frame.
	dirFrame int // Next frame direction (+1 or -1).
	duration time.Duration
	nextAt   time.Time
}

// Layout animates the frameset.
func (fs *FrameSet) Layout(gtx layout.Context) layout.Dimensions {
	img := fs.images[fs.curFrame]
	now := time.Now()
	if now.After(fs.nextAt) {
		if len(fs.frames) > 1 {
			fs.curFrame += fs.dirFrame
			if fs.curFrame == 0 || fs.curFrame == len(fs.frames)-1 {
				fs.dirFrame = -fs.dirFrame
			}
		}
		fs.nextAt = now.Add(fs.duration)
	}
	op.InvalidateOp{At: fs.nextAt}.Add(gtx.Ops)
	return img.Layout(gtx)
}

// ApplyTransparency makes an image transparent during specified amount of
// frames, allowing animation with the specified frame duration.
func ApplyTransparency(img image.Image, frames int, duration time.Duration, opts ...FrameSetOptionFunc) *FrameSet {
	if frames < 1 {
		frames = 1
	}
	fsC := make(chan *FrameSet, 1)
	fsC <- &FrameSet{
		frames:   make([]image.Image, frames),
		images:   make([]*ycwidget.Image, frames),
		dirFrame: 1,
		duration: duration,
		nextAt:   time.Now(),
	}
	var da float64
	if frames > 1 {
		da = 0xff / float64(frames-1)
	}
	type token struct{}
	done := make(chan token)
	for i := 0; i < frames; i++ {
		i := i
		go func() {
			defer func() { done <- token{} }()
			img := transparentImage(img, uint8(0xff-da*float64(i)))
			widget := ycwidget.NewImage(img)
			fs := <-fsC
			fs.frames[i] = img
			fs.images[i] = widget
			fsC <- fs
		}()
	}
	for i := 0; i < frames; i++ {
		<-done
	}
	fs := <-fsC
	for _, opt := range opts {
		opt(fs)
	}
	return fs
}

// FrameSetOptionFunc is a FrameSet option.
type FrameSetOptionFunc func(*FrameSet)

// ReversePlayback makes FrameSet playback frames in a reverse direction.
func ReversePlayback(fs *FrameSet) {
	fs.dirFrame = -fs.dirFrame
	fs.curFrame = len(fs.frames) - fs.curFrame - 1
}
