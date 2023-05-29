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
	w      fyne.Window
	c      *fyne.Container
	print  bool
	cntBar *CntWriter
}

func NewProgressBar(w fyne.Window, c *fyne.Container) *ProgressBar {
	return &ProgressBar{
		w:      w,
		c:      c,
		print:  true,
		cntBar: nil,
	}
}

func (b *ProgressBar) Write(p []byte) (n int, err error) {
	if b.cntBar != nil {
		_, _ = b.cntBar.Write(nil)
	}
	return 0, nil
}

func (b *ProgressBar) NewCntBar(size int64, description string) {
	name := widget.NewLabel(description)
	progress := widget.NewProgressBar()
	info := fmt.Sprintf("%d/%d", 0, size)
	completion := widget.NewLabelWithStyle(fmt.Sprintf("%s\t", info), fyne.TextAlignLeading, fyne.TextStyle{TabWidth: 14})
	progressReader := NewCntWriter(size, progress, completion)
	b.c.Add(container.NewHBox(
		name,
		layout.NewSpacer(),
		progress,
		completion,
	))
	canvas.Refresh(b.c)
	b.cntBar = progressReader
}

func (b *ProgressBar) SetPrint(print bool) {
	b.print = print
}

func (b *ProgressBar) NewBarReader(reader io.Reader, size int64, description string) io.Reader {
	if !b.print {
		return reader
	}
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
