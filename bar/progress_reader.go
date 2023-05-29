package bar

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
	"io"
	"sync/atomic"
)

type ProgressReader struct {
	reader     io.Reader
	progress   *widget.ProgressBar
	completion *widget.Label
	now        int64
	size       int64
}

func NewProgressReader(reader io.Reader, size int64, progress *widget.ProgressBar, completion *widget.Label) *ProgressReader {
	return &ProgressReader{
		reader:     reader,
		progress:   progress,
		completion: completion,
		now:        0,
		size:       size,
	}
}

func (r *ProgressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	atomic.AddInt64(&r.now, int64(n))
	if r.size == 0 {
		r.progress.SetValue(1)
	} else {
		r.progress.SetValue(float64(r.now) / float64(r.size))
	}
	info := fmt.Sprintf("%s/%s", humanize.IBytes(uint64(r.now)), humanize.IBytes(uint64(r.size)))
	r.completion.SetText(fmt.Sprintf("%s\t", info))
	return n, err
}
