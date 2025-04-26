//go:build darwin

package platform

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>

CGPoint currentMousePosition() {
	CGEventRef event = CGEventCreate(NULL);
	if (!event) return CGPointMake(0, 0);
	CGPoint point = CGEventGetLocation(event);
	CFRelease(event);
	return point;
}

CGEventType getMouseEventType(int button, int type) {
	switch (type) {
		case 0:
			switch (button) {
				case kCGMouseButtonLeft: return kCGEventLeftMouseDown;
				case kCGMouseButtonRight: return kCGEventRightMouseDown;
				case kCGMouseButtonCenter: return kCGEventOtherMouseDown;
				default: return kCGEventNull;
			}
		case 1:
			switch (button) {
				case kCGMouseButtonLeft: return kCGEventLeftMouseUp;
				case kCGMouseButtonRight: return kCGEventRightMouseUp;
				case kCGMouseButtonCenter: return kCGEventOtherMouseUp;
				default: return kCGEventNull;
			}
		case 2:
			switch (button) {
				case kCGMouseButtonLeft: return kCGEventLeftMouseDragged;
				case kCGMouseButtonRight: return kCGEventRightMouseDragged;
				case kCGMouseButtonCenter: return kCGEventOtherMouseDragged;
				default: return kCGEventNull;
			}
	}
	return kCGEventNull;
}


void moveMouseTo(int x, int y, int left_pressed, int right_pressed, int middle_pressed) {
	CGPoint newPoint = CGPointMake(x, y);
	CGMouseButton targetButton = kCGMouseButtonLeft;
	CGEventType dragEventType = kCGEventNull;

	if (left_pressed == 1) {
		dragEventType = kCGEventLeftMouseDragged;
		targetButton = kCGMouseButtonLeft;
	} else if (right_pressed == 1) {
		dragEventType = kCGEventRightMouseDragged;
		targetButton = kCGMouseButtonRight;
	} else if (middle_pressed == 1) {
		dragEventType = kCGEventOtherMouseDragged;
		targetButton = kCGMouseButtonCenter;
	}

	if (dragEventType != kCGEventNull) {
		CGEventRef dragEvent = CGEventCreateMouseEvent(NULL, dragEventType, newPoint, targetButton);
		if (!dragEvent) return;
		CGEventPost(kCGHIDEventTap, dragEvent);
		CFRelease(dragEvent);
	} else {
		CGWarpMouseCursorPosition(newPoint);
        CGEventRef moveEvent = CGEventCreateMouseEvent(NULL, kCGEventMouseMoved, newPoint, 0);
        if (!moveEvent) return;
		CGEventPost(kCGHIDEventTap, moveEvent);
		CFRelease(moveEvent);

	}
}

void mouseDown(int button) {
	CGEventType eventType = getMouseEventType(button, 0);
	if (eventType == kCGEventNull) return;

	CGPoint pos = currentMousePosition();
	CGEventRef event = CGEventCreateMouseEvent(NULL, eventType, pos, button);
	if (!event) return;
	CGEventPost(kCGHIDEventTap, event);
	CFRelease(event);
}

void mouseUp(int button) {
	CGEventType eventType = getMouseEventType(button, 1);
	if (eventType == kCGEventNull) return;

	CGPoint pos = currentMousePosition();
	CGEventRef event = CGEventCreateMouseEvent(NULL, eventType, pos, button);
	if (!event) return;
	CGEventPost(kCGHIDEventTap, event);
	CFRelease(event);
}

void mouseDoubleClick(int button) {
	CGEventType downEventType = getMouseEventType(button, 0);
	CGEventType upEventType = getMouseEventType(button, 1);
	if (downEventType == kCGEventNull || upEventType == kCGEventNull) return;

	CGPoint pos = currentMousePosition();

	CGEventRef down1 = CGEventCreateMouseEvent(NULL, downEventType, pos, button);
	if (!down1) return;
	CGEventSetIntegerValueField(down1, kCGMouseEventClickState, 1);
	CGEventPost(kCGHIDEventTap, down1);
	CFRelease(down1);

	CGEventRef up1 = CGEventCreateMouseEvent(NULL, upEventType, pos, button);
	if (!up1) {  return; }
CGEventSetIntegerValueField(up1, kCGMouseEventClickState, 1);
CGEventPost(kCGHIDEventTap, up1);
CFRelease(up1);

CGEventRef down2 = CGEventCreateMouseEvent(NULL, downEventType, pos, button);
if (!down2) return;
CGEventSetIntegerValueField(down2, kCGMouseEventClickState, 2);
CGEventPost(kCGHIDEventTap, down2);
CFRelease(down2);

CGEventRef up2 = CGEventCreateMouseEvent(NULL, upEventType, pos, button);
if (!up2) return;
CGEventSetIntegerValueField(up2, kCGMouseEventClickState, 2);
CGEventPost(kCGHIDEventTap, up2);
CFRelease(up2);
}

void scrollBoth(int deltaX, int deltaY) {

CGEventRef event = CGEventCreateScrollWheelEvent2(NULL, kCGScrollEventUnitPixel, 2, deltaY, deltaX, 0);
CGEventPost(kCGHIDEventTap, event);
CFRelease(event);
}
*/
import "C"

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var _ InputController = (*macosInputController)(nil)

type macosInputController struct {
	mu           sync.Mutex
	buttonStates map[string]bool

	lastClickTime        map[string]time.Time
	doubleClickThreshold time.Duration
}

func NewMacOSInputController() InputController {
	return &macosInputController{
		buttonStates:         make(map[string]bool),
		lastClickTime:        make(map[string]time.Time),
		doubleClickThreshold: time.Millisecond * 500,
	}
}

func mapButtonStrToC(button string) (C.int, error) {
	switch button {
	case "left":
		return C.kCGMouseButtonLeft, nil
	case "right":
		return C.kCGMouseButtonRight, nil
	case "middle":
		return C.kCGMouseButtonCenter, nil
	default:
		return -1, fmt.Errorf("unsupported button: %s", button)
	}
}

func (mic *macosInputController) getButtonStateAsCFlags() (C.int, C.int, C.int) {
	var left, right, middle C.int = 0, 0, 0
	if pressed, ok := mic.buttonStates["left"]; ok && pressed {
		left = 1
	}
	if pressed, ok := mic.buttonStates["right"]; ok && pressed {
		right = 1
	}
	if pressed, ok := mic.buttonStates["middle"]; ok && pressed {
		middle = 1
	}
	return left, right, middle
}

func (mic *macosInputController) MoveMouseBy(dx, dy int) error {
	mic.mu.Lock()
	defer mic.mu.Unlock()

	currentPosC := C.currentMousePosition()
	currentX := int(currentPosC.x)
	currentY := int(currentPosC.y)

	newX := currentX + dx
	newY := currentY + dy

	left, right, middle := mic.getButtonStateAsCFlags()

	C.moveMouseTo(C.int(newX), C.int(newY), left, right, middle)

	return nil
}

func (mic *macosInputController) SetButtonState(button string, pressed bool) error {
	mic.mu.Lock()
	defer mic.mu.Unlock()

	log.Printf("[SetButtonState] Received: button=%s, pressed=%t", button, pressed)

	currentPressed, exists := mic.buttonStates[button]
	if exists && currentPressed == pressed {
		log.Printf("[SetButtonState] State for button '%s' is already %t, skipping.", button, pressed)
		return nil
	}

	cButton, err := mapButtonStrToC(button)
	if err != nil {
		return err
	}

	if pressed {
		log.Printf("[SetButtonState] Calling mouseDown for button '%s' (%d)", button, cButton)
		C.mouseDown(cButton)
	} else {
		log.Printf("[SetButtonState] Calling mouseUp for button '%s' (%d)", button, cButton)
		C.mouseUp(cButton)
	}

	mic.buttonStates[button] = pressed
	log.Printf("[SetButtonState] State for button '%s' updated to %t", button, pressed)
	return nil
}

func (mic *macosInputController) ScrollBy(dx, dy int) error {
	mic.mu.Lock()
	defer mic.mu.Unlock()
	return mic.scrollByInternal(-dx*30, -dy*30)
}

func (mic *macosInputController) scrollByInternal(dx, dy int) error {
	if dx == 0 && dy == 0 {
		return nil
	}
	C.scrollBoth(C.int(dx), C.int(dy))
	return nil
}

func (mic *macosInputController) DoubleClick(button string) error {
	mic.mu.Lock()
	defer mic.mu.Unlock()

	log.Printf("[DoubleClick] Received: button=%s", button)

	cButton, err := mapButtonStrToC(button)
	if err != nil {
		return err
	}

	C.mouseDoubleClick(cButton)

	return nil
}

func (mic *macosInputController) Close() error {
	log.Println("macos controller: cleaning up button states...")
	mic.mu.Lock()
	defer mic.mu.Unlock()

	for btn, pressed := range mic.buttonStates {
		if pressed {
			log.Printf("macos controller: Releasing potentially stuck button %s", btn)
			mic.setButtonStateInternal(btn, false)
		}
	}
	log.Println("macos controller: Button cleanup finished.")
	return nil
}

func (mic *macosInputController) setButtonStateInternal(button string, pressed bool) error {
	cButton, err := mapButtonStrToC(button)
	if err != nil {
		log.Printf("Error mapping button '%s' in internal set state: %v", button, err)
		return err
	}

	if pressed {
		C.mouseDown(cButton)
	} else {
		C.mouseUp(cButton)
	}

	mic.buttonStates[button] = pressed
	return nil
}
