package gui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework UserNotifications
#import <Cocoa/Cocoa.h>
#import <UserNotifications/UserNotifications.h>

static void requestAttention(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
		[NSApp requestUserAttention:NSInformationalRequest];
	});
}

static void bringToFront(uintptr_t ptr) {
	NSWindow *w = (NSWindow *)ptr;
	dispatch_async(dispatch_get_main_queue(), ^{
		if ([w isMiniaturized]) {
			[w deminiaturize:nil];
		}
		[NSApp activateIgnoringOtherApps:YES];
		[w makeKeyAndOrderFront:nil];
	});
}

static int reduceMotion(void) {
	return [[NSWorkspace sharedWorkspace] accessibilityDisplayShouldReduceMotion] ? 1 : 0;
}

static void askToNotify(void) {
	if ([[NSBundle mainBundle] bundleIdentifier] == nil) {
		return;
	}
	[[UNUserNotificationCenter currentNotificationCenter]
		requestAuthorizationWithOptions:UNAuthorizationOptionAlert
		completionHandler:^(BOOL granted, NSError *_Nullable error) {}];
}
*/
import "C"

import "fyne.io/fyne/v2/driver"

func attentionNative(any) { C.requestAttention() }

func frontNative(ctx any) {
	if c, ok := ctx.(driver.MacWindowContext); ok && c.NSWindow != 0 {
		C.bringToFront(C.uintptr_t(c.NSWindow))
	}
}

func reduceMotionNative() bool { return C.reduceMotion() == 1 }

func askToNotifyNative() { C.askToNotify() }
