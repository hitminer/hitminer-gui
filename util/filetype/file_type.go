package filetype

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"path/filepath"
	"strings"
)

func GetFileType(filename string) fyne.Resource {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico":
		return theme.FileImageIcon()
	case ".txt":
		return theme.FileTextIcon()
	case ".doc", ".docx", ".odt", ".pdf", ".ppt", ".pptx", ".rtf", ".xls", ".xlsx":
		return theme.DocumentIcon()
	case ".apk", ".exe":
		return theme.FileApplicationIcon()
	case ".mp3", ".wav", ".wma":
		return theme.FileAudioIcon()
	case ".mp4", ".avi", ".flv", ".mkv", ".mov", ".wmv":
		return theme.FileVideoIcon()
	default:
		return theme.FileIcon()
	}
}
