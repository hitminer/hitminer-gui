package ui

import (
	"archive/tar"
	"archive/zip"
	"compress/lzw"
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/hitminer/hitminer-file-manager/util/multibar"
	"github.com/hitminer/hitminer-gui/bar"
	"github.com/hitminer/hitminer-gui/vars"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func UpgradeWindows() {
	a := fyne.CurrentApp()
	w := a.NewWindow("更新")

	ctx, cancel := context.WithCancel(context.Background())
	w.SetOnClosed(func() {
		cancel()
	})

	vbox := container.NewVBox()
	bars := bar.NewProgressBar(w, vbox)
	go func() {
		err := upgrade(ctx, bars)
		if err != nil {
			dialog.ShowError(err, w)
		}

		w.SetOnClosed(func() {
			cancel()
			a.Quit()
		})
		vbox.Add(widget.NewButton("确认", func() {
			a.Quit()
		}))
		canvas.Refresh(vbox)
	}()

	w.SetContent(vbox)
	w.Resize(fyne.NewSize(500, 200))
	w.Show()
}

func upgrade(ctx context.Context, bars multibar.MultiBar) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	nowDir := filepath.Dir(executable)
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		nowDir = filepath.Dir(filepath.Dir(filepath.Dir(nowDir)))
	}
	downloadName := filepath.Base(vars.DownloadUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, vars.DownloadUrl, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	tempDir, err := os.MkdirTemp("", "hitminer")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()
	downloadFile, err := os.CreateTemp(tempDir, downloadName)
	if err != nil {
		return err
	}

	b := bars.NewBarReader(resp.Body, resp.ContentLength, "更新下载中")
	_, err = io.Copy(downloadFile, b)
	if err != nil {
		_ = downloadFile.Close()
		return err
	}
	err = downloadFile.Close()
	if err != nil {
		return err
	}

	if strings.HasSuffix(downloadName, ".zip") {
		files, err := zip.OpenReader(downloadFile.Name())
		if err != nil {
			return err
		}
		for _, f := range files.File {
			filePath := filepath.Join(nowDir, f.Name)
			if f.FileInfo().IsDir() {
				err := os.MkdirAll(filePath, 0755)
				if err != nil {
					return err
				}
				continue
			}

			err := os.MkdirAll(filepath.Dir(filePath), 0755)
			if err != nil {
				return err
			}

			err = os.MkdirAll(filepath.Dir(filepath.Join(tempDir, f.Name)), 0755)
			if err != nil {
				return err
			}
			tempFile, err := os.OpenFile(filepath.Join(tempDir, f.Name), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
			if err != nil {
				return err
			}

			fileInArchive, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(tempFile, fileInArchive)
			if err != nil {
				_ = tempFile
				_ = fileInArchive.Close()
				return err
			}

			err = tempFile.Close()
			if err != nil {
				return err
			}
			err = fileInArchive.Close()
			if err != nil {
				return err
			}

			err = os.MkdirAll(filepath.Dir(filepath.Join(tempDir, filepath.Base(filePath)+".old")), 0755)
			if err != nil {
				return err
			}
			err = os.Rename(filePath, filepath.Join(tempDir, filepath.Base(filePath)+".old"))
			if err != nil {
				return err
			}

			err = os.Rename(tempFile.Name(), filePath)
			if err != nil {
				return err
			}
		}
	} else if strings.HasSuffix(downloadName, "tar.xz") {
		tarFile, err := os.Open(downloadFile.Name())
		if err != nil {
			return err
		}
		defer func() {
			_ = tarFile.Close()
		}()

		lzwReader := lzw.NewReader(tarFile, lzw.LSB, 8)
		defer func() {
			_ = lzwReader.Close()
		}()

		tarReader := tar.NewReader(lzwReader)
		for {
			f, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			filePath := filepath.Join(nowDir, f.Name)
			if f.FileInfo().IsDir() {
				err := os.MkdirAll(filePath, 0755)
				if err != nil {
					return err
				}
				continue
			}

			err = os.MkdirAll(filepath.Dir(filePath), 0755)
			if err != nil {
				return err
			}

			err = os.MkdirAll(filepath.Dir(filepath.Join(tempDir, f.Name)), 0755)
			if err != nil {
				return err
			}
			tempFile, err := os.OpenFile(filepath.Join(tempDir, f.Name), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
			if err != nil {
				return err
			}

			_, err = io.Copy(tempFile, tarReader)
			if err != nil {
				_ = tempFile.Close()
				return err
			}

			err = tempFile.Close()
			if err != nil {
				return err
			}

			err = os.MkdirAll(filepath.Dir(filepath.Join(tempDir, filepath.Base(filePath)+".old")), 0755)
			if err != nil {
				return err
			}
			err = os.Rename(filePath, filepath.Join(tempDir, filepath.Base(filePath)+".old"))
			if err != nil {
				return err
			}

			err = os.Rename(tempFile.Name(), filePath)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("unsupport file type")
	}
	return nil
}
