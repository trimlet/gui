package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
)

// The window asks for attention (Trimlet Chat, when something happens while it is behind): the icon
// bounces in the Dock on Mac, flashes in the taskbar on Windows. Nothing elsewhere.
func RequestAttention(w fyne.Window) {
	if nw, ok := w.(driver.NativeWindow); ok {
		nw.RunNative(func(ctx any) { attentionNative(ctx) })
	}
}

// BringToFront raises the window, also from behind other apps or minimized (a second copy of the app
// asks the first to come forward).
func BringToFront(w fyne.Window) {
	w.Show()
	w.RequestFocus()
	if nw, ok := w.(driver.NativeWindow); ok {
		nw.RunNative(func(ctx any) { frontNative(ctx) })
	}
}

// ReduceMotion says whether the system asks for less motion (macOS: Accessibility → Display → Reduce
// motion; Windows: Settings → Accessibility → Visual effects → Animation effects off).
func ReduceMotion() bool { return reduceMotionNative() }

// AskToNotify asks the system, once, the permission to show notifications, at a moment the user is
// looking (macOS; the app must be a bundle). Elsewhere nothing is asked.
func AskToNotify() { askToNotifyNative() }
