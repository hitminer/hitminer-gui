package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/hitminer/hitminer-file-manager/server/s3gateway"
	"github.com/hitminer/hitminer-gui/bar"
	"github.com/hitminer/hitminer-gui/vars"
	"strings"
)

func UploadWindows(filePath, objectName string) {
	if objectName == "." {
		objectName = ""
	}
	if objectName != "" {
		if !strings.HasSuffix(objectName, "/") {
			objectName = objectName + "/"
		}
	}

	a := fyne.CurrentApp()
	w := a.NewWindow("上传")

	ctx, cancel := context.WithCancel(context.Background())
	w.SetOnClosed(func() {
		cancel()
	})

	vbox := container.NewVBox()
	svr := s3gateway.NewS3Server(ctx, vars.Host, vars.Token, bar.NewProgressBar(w, vbox))
	go func() {
		err := svr.PutObjects(ctx, filePath, objectName)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
	}()

	w.Resize(fyne.NewSize(500, 400))
	w.SetContent(container.NewVScroll(vbox))
	w.Show()
}
