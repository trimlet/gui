//go:build !darwin && !windows

package gui

func attentionNative(any)      {}
func frontNative(any)          {}
func reduceMotionNative() bool { return false }
func askToNotifyNative()       {}
