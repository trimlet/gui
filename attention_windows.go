package gui

import (
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2/driver"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	flashWindowEx        = user32.NewProc("FlashWindowEx")
	setForegroundWindow  = user32.NewProc("SetForegroundWindow")
	isIconic             = user32.NewProc("IsIconic")
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

type flashInfo struct {
	size    uint32
	hwnd    uintptr
	flags   uint32
	count   uint32
	timeout uint32
}

const (
	flashTray        = 0x2 // the taskbar button
	flashUntilFront  = 0xC // until the window comes to the front
	swRestore        = 9
	spiAreaAnimation = 0x1042 // SPI_GETCLIENTAREAANIMATION
)

func attentionNative(ctx any) {
	c, ok := ctx.(driver.WindowsWindowContext)
	if !ok || c.HWND == 0 {
		return
	}
	fi := flashInfo{hwnd: c.HWND, flags: flashTray | flashUntilFront}
	fi.size = uint32(unsafe.Sizeof(fi))
	flashWindowEx.Call(uintptr(unsafe.Pointer(&fi)))
}

func frontNative(ctx any) {
	c, ok := ctx.(driver.WindowsWindowContext)
	if !ok || c.HWND == 0 {
		return
	}
	if r, _, _ := isIconic.Call(c.HWND); r != 0 {
		showWindow.Call(c.HWND, swRestore)
	}
	setForegroundWindow.Call(c.HWND)
}

func reduceMotionNative() bool {
	var on int32 = 1
	r, _, _ := systemParametersInfo.Call(spiAreaAnimation, 0, uintptr(unsafe.Pointer(&on)), 0)
	return r != 0 && on == 0
}

func askToNotifyNative() {}
