package gui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

// Colors the title bar: the bar becomes transparent and shows the window background,
// which gets the app color. The dark appearance makes the title white. The buttons,
// dragging and double click stay the ones of macOS.
static void colorTitleBar(uintptr_t ptr, double r, double g, double b) {
	NSWindow *w = (NSWindow *)ptr;
	dispatch_async(dispatch_get_main_queue(), ^{
		w.titlebarAppearsTransparent = YES;
		w.backgroundColor = [NSColor colorWithSRGBRed:r green:g blue:b alpha:1];
		w.appearance = [NSAppearance appearanceNamed:NSAppearanceNameDarkAqua];
	});
}
*/
import "C"

import (
	"image/color"

	"fyne.io/fyne/v2/driver"
)

func titleBarNative(ctx any, c color.Color) bool {
	mc, ok := ctx.(driver.MacWindowContext)
	if !ok || mc.NSWindow == 0 {
		return false
	}
	r, g, b, _ := c.RGBA()
	C.colorTitleBar(C.uintptr_t(mc.NSWindow), C.double(r)/0xffff, C.double(g)/0xffff, C.double(b)/0xffff)
	return true
}
