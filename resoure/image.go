package resoure

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed static/image/folder.png
var folder []byte

var Folder = &fyne.StaticResource{
	StaticName:    "folder",
	StaticContent: folder,
}

//go:embed static/image/file.png
var file []byte

var File = &fyne.StaticResource{
	StaticName:    "file",
	StaticContent: file,
}
