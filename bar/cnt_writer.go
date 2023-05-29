package bar

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"sync/atomic"
)

type CntWriter struct {
	progress   *widget.ProgressBar
	completion *widget.Label
	now        int64
	size       int64
}

func NewCntWriter(size int64, progress *widget.ProgressBar, completion *widget.Label) *CntWriter {
	return &CntWriter{
		progress:   progress,
		completion: completion,
		now:        0,
		size:       size,
	}
}

func (r *CntWriter) Write(p []byte) (int, error) {
	atomic.AddInt64(&r.now, 1)
	if r.size == 0 {
		r.progress.SetValue(1)
	} else {
		r.progress.SetValue(float64(r.now) / float64(r.size))
	}
	info := fmt.Sprintf("%d/%d", r.now, r.size)
	r.completion.SetText(fmt.Sprintf("%s\t", info))
	return 0, nil
}
