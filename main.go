package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/hitminer/hitminer-gui/resoure"
	"github.com/hitminer/hitminer-gui/ui"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&resoure.FontTheme{})

	w := a.NewWindow("hitminer")
	w.SetContent(ui.LoginContainer(w))

	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
