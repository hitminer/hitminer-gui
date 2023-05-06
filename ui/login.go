package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/hitminer/hitminer-file-manager/login"
	"github.com/hitminer/hitminer-gui/vars"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func init() {
	viper.SetDefault("host", "www.hitminer.cn")
	viper.SetConfigName("file_manager")
	viper.SetConfigType("toml")
	home, _ := os.UserHomeDir()
	viper.AddConfigPath(filepath.Join(home, ".config", "hitminer"))
}

func LoginContainer(w fyne.Window) fyne.CanvasObject {
	_ = viper.ReadInConfig()
	defaultUser := viper.GetString("username")
	defaultPW := viper.GetString("password")
	defaultHost := viper.GetString("host")

	host := widget.NewEntry()
	host.SetPlaceHolder("集群")
	host.SetText(defaultHost)

	username := widget.NewEntry()
	username.SetPlaceHolder("用户名")
	username.SetText(defaultUser)

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("密码")
	password.SetText(defaultPW)

	form := widget.NewForm(
		widget.NewFormItem("集群", host),
		widget.NewFormItem("用户名", username),
		widget.NewFormItem("密码", password),
	)

	button := widget.NewButton("登入", func() {
		token, err := login.Login(host.Text, username.Text, password.Text)
		if err != nil {
			dialog.NewCustom("错误", "确定", widget.NewLabelWithStyle("登入失败", fyne.TextAlignCenter, fyne.TextStyle{}), w).Show()
			return
		}
		vars.Host = host.Text
		vars.Token = token
		home, _ := os.UserHomeDir()
		viper.Set("username", username.Text)
		viper.Set("password", password.Text)
		viper.Set("host", host.Text)
		path := filepath.Join(home, ".config", "hitminer")
		_ = os.MkdirAll(path, 0755)
		_ = viper.WriteConfigAs(filepath.Join(path, "file_manager.toml"))
		w.SetContent(DirectoryContainer("", w))
	})

	return container.New(
		layout.NewVBoxLayout(),
		form,
		button,
	)
}
