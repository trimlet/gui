// Package gui holds what the trimlet apps share: the dark console look, the text area
// with mouse selection and the window that opens maximized.
package gui

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/theme"
)

const TextSize = 13

// Window size in characters when not maximized.
const (
	DefaultCols = 150
	DefaultRows = 20
)

// Smallest window size in characters.
const (
	MinCols = 60
	MinRows = 14
)

// Background is the console background: a fixed near-black, not one of the app colors
// (the black palette was dropped on 2026-11-18).
var Background color.Color = color.NRGBA{0x16, 0x17, 0x19, 0xFF}

// White is the bright text (e.g. the questions of the windows in Trimlet Files): fixed, like the icon's white body.
var White color.Color = color.NRGBA{0xF2, 0xF2, 0xF2, 0xFF}

// SelectedBg is the background of the selected row of a list: trimlet's green, in every app.
var SelectedBg color.Color = Palettes["green"][9]

// Colors shared by the apps, all from the trimlet palettes (palette.go), the same in every
// app; only the title bar takes the app color (UsePalette).
var (
	Accent  color.Color = Palettes["green"][4]
	Amber   color.Color = Palettes["copper"][4] // messages such as "✔ Row copied"
	Red     color.Color = Palettes["red"][4]    // errors
	Muted   color.Color = Palettes["grey"][5]
	ValueFg color.Color = Palettes["blue"][3] // values (a field's current value, JSON/XML values)
)

// IconStep is the index in Palettes of the icon color of name (also the title bar): 5 (500),
// except red, whose 500 was too bright on the title bar (site/data/iconstep.yaml).
func IconStep(name string) int {
	if name == "red" {
		return 6
	}
	return 5
}

// UsePalette gives the window's title bar the icon color of the trimlet color name (a key of
// Palettes), before Run. Inside the console every app otherwise keeps trimlet's green (the defaults above:
// accent green 400, selected row green 900, background fixed near-black, muted grey 500).
// It returns false for an unknown name.
func UsePalette(name string) bool {
	p, ok := Palettes[name]
	if !ok {
		return false
	}
	TitleBar = p[IconStep(name)]
	return true
}

// ConsoleTheme is a dark theme with our background color and text size.
type ConsoleTheme struct{ fyne.Theme }

func (t ConsoleTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if n == theme.ColorNameBackground {
		return Background
	}
	return t.Theme.Color(n, theme.VariantDark)
}

func (t ConsoleTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameText {
		return TextSize
	}
	return t.Theme.Size(n)
}

// CellSize is the size of one character of the console. It is rounded to whole pixels as
// Fyne's TextGrid does: otherwise the app counts more columns than the grid really shows,
// and the end of the longest lines is cut off.
func CellSize() fyne.Size {
	s := fyne.MeasureText("M", TextSize, fyne.TextStyle{Monospace: true})
	return fyne.NewSize(float32(math.Round(float64(s.Width))), float32(math.Round(float64(s.Height))))
}

// NewWindow creates the app and its window with the console theme and no outer margin:
// each screen sets its own, so its background reaches the edges.
func NewWindow(id, title string) (fyne.App, fyne.Window) {
	a := app.NewWithID(id)
	a.Settings().SetTheme(ConsoleTheme{theme.DefaultTheme()})
	w := a.NewWindow(title)
	w.SetPadded(false)
	return a, w
}

// Run shows the window maximized and runs the app until it is closed.
// started, if given, is called once the window is on screen.
func Run(a fyne.App, w fyne.Window, started ...func()) {
	// Size used if the window cannot be maximized.
	cell := CellSize()
	pad := theme.Padding() * 2
	w.Resize(fyne.NewSize(cell.Width*DefaultCols+pad, cell.Height*DefaultRows+pad))
	w.CenterOnScreen()
	a.Lifecycle().SetOnStarted(func() {
		maximize(w, 20)
		if TitleBar != nil {
			colorTitleBar(w, TitleBar, 20)
		}
		for _, f := range started {
			f()
		}
	})
	w.ShowAndRun()
}

// maximize enlarges the window to fill the screen, leaving the menu bar, Dock
// and taskbar visible. The native window may not exist yet right after start,
// so it retries a few times.
// TitleBar, if set before Run, is the color of the window's title bar (only on Mac
// for now; elsewhere the bar stays standard).
var TitleBar color.Color

// colorTitleBar gives the title bar the color c, retrying like maximize.
func colorTitleBar(w fyne.Window, c color.Color, tries int) {
	nw, ok := w.(driver.NativeWindow)
	if !ok {
		return
	}
	done := false
	nw.RunNative(func(ctx any) { done = titleBarNative(ctx, c) })
	if !done && tries > 0 {
		time.AfterFunc(50*time.Millisecond, func() { fyne.Do(func() { colorTitleBar(w, c, tries-1) }) })
	}
}

func maximize(w fyne.Window, tries int) {
	nw, ok := w.(driver.NativeWindow)
	if !ok {
		return
	}
	done := false
	nw.RunNative(func(ctx any) { done = maximizeNative(ctx) })
	if !done && tries > 0 {
		time.AfterFunc(50*time.Millisecond, func() { fyne.Do(func() { maximize(w, tries-1) }) })
	}
}
