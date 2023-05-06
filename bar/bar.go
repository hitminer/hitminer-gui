package bar

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
	"io"
)

type ProgressBar struct {
	w fyne.Window
	c *fyne.Container
}

func NewProgressBar(w fyne.Window, c *fyne.Container) *ProgressBar {
	return &ProgressBar{
		w: w,
		c: c,
	}
}

func (b *ProgressBar) NewBarReader(reader io.Reader, size int64, description string) io.Reader {
	name := widget.NewLabel(description)
	progress := widget.NewProgressBar()
	info := fmt.Sprintf("%s/%s", humanize.IBytes(0), humanize.IBytes(uint64(size)))
	completion := widget.NewLabelWithStyle(fmt.Sprintf("%s\t", info), fyne.TextAlignLeading, fyne.TextStyle{TabWidth: 14})
	progressReader := NewProgressReader(reader, size, progress, completion)
	b.c.Add(container.NewHBox(
		name,
		layout.NewSpacer(),
		progress,
		completion,
	))
	canvas.Refresh(b.c)
	return progressReader
}

func (b *ProgressBar) Wait() {
	closeButton := widget.NewButton("确认", func() {
		b.w.Close()
	})
	b.c.Add(closeButton)
	b.c.Refresh()
}
