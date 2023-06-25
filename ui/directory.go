package ui

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
	"github.com/hitminer/hitminer-file-manager/server/s3gateway"
	"github.com/hitminer/hitminer-gui/util/filetype"
	"github.com/hitminer/hitminer-gui/vars"
	"path/filepath"
	"strings"
	"time"
)

func DirectoryContainer(prefix string, w fyne.Window) fyne.CanvasObject {
	ctx, cancel := context.WithCancel(context.Background())
	w.SetOnClosed(func() {
		cancel()
	})
	svr := s3gateway.NewS3Server(ctx, vars.Host, vars.Token, nil)
	selectEntry := widget.NewSelect([]string{"上传文件", "上传文件夹", "上传ERO文件", "重传文件夹"}, func(selected string) {
		if selected == "上传文件" {
			dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if reader == nil {
					return
				}
				path := reader.URI().Path()
				_ = reader.Close()
				UploadWindows(path, prefix, false, false)
			}, w).Show()
		} else if selected == "上传文件夹" {
			dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if dir == nil {
					return
				}
				UploadWindows(dir.Path(), prefix, false, false)
			}, w).Show()
		} else if selected == "上传ERO文件" {
			dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if dir == nil {
					return
				}
				UploadWindows(dir.Path(), prefix, true, false)
			}, w).Show()
		} else if selected == "重传文件夹" {
			dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if dir == nil {
					return
				}
				UploadWindows(dir.Path(), prefix, false, true)
			}, w).Show()
		}
	})
	selectEntry.PlaceHolder = "上传       "
	selectEntry.Resize(fyne.NewSize(200, 20))

	uploadFileButton := widget.NewToolbarAction(theme.UploadIcon(), func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			path := reader.URI().Path()
			_ = reader.Close()
			UploadWindows(path, prefix, false, false)
		}, w).Show()
	})
	uploadFolderButton := widget.NewToolbarAction(theme.MoveUpIcon(), func() {
		dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if dir == nil {
				return
			}
			UploadWindows(dir.Path(), prefix, false, false)
		}, w).Show()
	})
	createFolderButton := widget.NewToolbarAction(theme.FolderNewIcon(), func() {
		dir := widget.NewEntry()
		dir.SetPlaceHolder("文件夹名")
		callback := func(response bool) {
			if response {
				dirName := dir.Text
				if dirName == "" {
					return
				}
				svr := s3gateway.NewS3Server(ctx, vars.Host, vars.Token, nil)
				err := svr.MakeDirectory(ctx, filepath.Join(prefix, dirName))
				if err != nil {
					dialog.ShowError(err, w)
				}
				cancel()
				w.SetContent(DirectoryContainer(prefix, w))
			}
		}

		form := dialog.NewForm("文件夹名:", "确认", "取消", []*widget.FormItem{widget.NewFormItem("文件夹名:", dir)}, callback, w)
		form.Resize(fyne.NewSize(400, 100))
		form.Show()
	})
	refreshButton := widget.NewToolbarAction(theme.MediaReplayIcon(), func() {
		cancel()
		w.SetContent(DirectoryContainer(prefix, w))
	})
	returnButton := widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
		if prefix != "" && prefix != "." {
			dir := filepath.Dir(strings.TrimSuffix(prefix, "/"))
			cancel()
			w.SetContent(DirectoryContainer(dir, w))
		}
	})
	logoutButton := widget.NewToolbarAction(theme.LogoutIcon(), func() {
		cancel()
		w.SetContent(LoginContainer(w))
	})
	updateButton := widget.NewToolbarAction(theme.HelpIcon(), func() {
		dialog.ShowCustomConfirm("信息", "确认", "更新", widget.NewLabel(fmt.Sprintf("当前版本 %s", vars.Version)), func(success bool) {
			if !success {
				UpgradeWindows()
			}
		}, w)
	})

	toolbar := widget.NewToolbar(uploadFileButton, uploadFolderButton, createFolderButton, refreshButton, returnButton, logoutButton, updateButton)
	up := container.NewHBox(selectEntry, toolbar)

	vbox := container.NewVBox()
	if prefix != "" && prefix != "." {
		dir := filepath.Dir(strings.TrimSuffix(prefix, "/"))
		hyperlink := widget.NewHyperlink("..", nil)
		hyperlink.OnTapped = func() {
			cancel()
			w.SetContent(DirectoryContainer(dir, w))
		}
		var item fyne.CanvasObject = container.NewHBox(
			widget.NewIcon(theme.FolderIcon()),
			hyperlink,
			layout.NewSpacer(),
		)
		vbox.Add(item)
	}

	go func() {
		vboxList := make([]fyne.CanvasObject, 0)
		for obj := range svr.ListObjects(ctx, prefix, "/") {
			object := obj
			var sizeLabel, nameLabel, icon fyne.CanvasObject
			if object.IsDirectory {
				sizeLabel = layout.NewSpacer()
				hyperlink := widget.NewHyperlink(object.Name, nil)
				hyperlink.OnTapped = func() {
					cancel()
					w.SetContent(DirectoryContainer(object.FullPath, w))
				}
				nameLabel = hyperlink
				icon = widget.NewIcon(theme.FolderIcon())
			} else {
				sizeLabel = widget.NewLabelWithStyle(humanize.IBytes(uint64(object.Size)), fyne.TextAlignTrailing, fyne.TextStyle{})
				nameLabel = widget.NewLabel(object.Name)
				icon = widget.NewIcon(filetype.GetFileType(object.Name))
			}
			deleteButton := widget.NewToolbarAction(theme.DeleteIcon(), func() {
				svr := s3gateway.NewS3Server(ctx, vars.Host, vars.Token, nil)
				err := svr.RemoveObjects(ctx, object.FullPath, object.IsDirectory)
				if err != nil {
					dialog.ShowError(err, w)
				}
				cancel()
				w.SetContent(DirectoryContainer(prefix, w))
			})
			downloadButton := widget.NewToolbarAction(theme.DownloadIcon(), func() {
				dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
					if err != nil {
						dialog.ShowError(err, w)
						return
					}
					if dir == nil {
						return
					}
					DownloadWindows(dir.Path(), object.FullPath)
				}, w).Show()
			})
			buttonBar := widget.NewToolbar(deleteButton, downloadButton)
			item := container.NewHBox(
				icon,
				nameLabel,
				layout.NewSpacer(),
				sizeLabel,
				buttonBar,
			)
			vboxList = append(vboxList, item)
			if len(vboxList)%5 == 0 {
				for _, o := range vboxList {
					vbox.Add(o)
				}
				canvas.Refresh(vbox)
				vboxList = vboxList[:0]
				time.Sleep(100 * time.Millisecond)
			}
		}
		if len(vboxList) != 0 {
			for _, o := range vboxList {
				vbox.Add(o)
			}
			canvas.Refresh(vbox)
			vboxList = vboxList[:0]
		}
	}()

	return container.NewBorder(up, nil, nil, nil, container.NewVScroll(vbox))
}
