package resoure

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

//go:embed static/font/PingFang.ttf
var pingFang []byte

var fontResource = &fyne.StaticResource{
	StaticName:    "font",
	StaticContent: pingFang,
}

type FontTheme struct{}

var _ fyne.Theme = (*FontTheme)(nil)

func (*FontTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return theme.DefaultTheme().Font(s)
	}
	if s.Bold {
		if s.Italic {
			return theme.DefaultTheme().Font(s)
		}
		return fontResource
	}
	if s.Italic {
		return theme.DefaultTheme().Font(s)
	}
	return fontResource

}

func (*FontTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*FontTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*FontTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
