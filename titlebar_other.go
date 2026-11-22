//go:build !darwin

package gui

import "image/color"

// Other systems keep the standard title bar.
func titleBarNative(any, color.Color) bool { return true }
